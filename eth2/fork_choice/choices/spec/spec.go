package spec

import (
	"lmd-ghost/eth2/dag"
	"lmd-ghost/eth2/fork_choice"
)

/// The naive, but readable, spec implementation
type SpecLMDGhost struct {

	dag *dag.BeaconDag

	LatestScores map[*dag.DagNode]int64
}

func NewSpecLMDGhost() fork_choice.ForkChoice {
	return new(SpecLMDGhost)
}

func (gh *SpecLMDGhost) SetDag(dag *dag.BeaconDag) {
	gh.dag = dag
}

func (gh *SpecLMDGhost) ApplyScoreChanges(changes []fork_choice.ScoreChange) {
	for _, v := range changes {
		gh.LatestScores[v.Target] += v.ScoreDelta
	}
	// delete targets that have a 0 score
	for k, v := range gh.LatestScores {
		if v == 0 {
			// deletion during map iteration, safe in Go
			delete(gh.LatestScores, k)
		}
	}
}

func (gh *SpecLMDGhost) OnNewNode(node *dag.DagNode) {
	// free, at cost of head-function
}

func (gh *SpecLMDGhost) OnStartChange(newStart *dag.DagNode) {
	// nothing to do when the start changes
}

/// Retrieves the head by *recursively* looking for the highest voted block
//   at *every* block in the path from start to head.
func (gh *SpecLMDGhost) HeadFn() *dag.DagNode {
	// Minor difference:
	// Normally you would have to filter for the active validators, and get their targets.
	// We can just iterate over the values in the common-chain.
	// This difference only really matters when there's many validators inactive,
	//  and the client implementation doesn't store them separately.

	head := gh.dag.Start
	for {
		if len(head.Children) == 0 {
			return head
		}
		bestItem := head.Children[0]
		var bestScore int64 = 0
		for _, child := range head.Children {
			childVotes := gh.getVoteCount(child)
			if childVotes > bestScore {
				bestScore = childVotes
				bestItem = child
			}
		}
		head = bestItem
	}
}

func (gh *SpecLMDGhost) getVoteCount(block *dag.DagNode) int64 {
	totalWeight := int64(0)
	for target, weight := range gh.LatestScores {
		if anc := gh.getAncestor(target, block.Slot); anc != nil && anc == target {
			totalWeight += weight
		}
	}
	return totalWeight
}

/// Gets the ancestor of `node` at `slot`
func (gh *SpecLMDGhost) getAncestor(block *dag.DagNode, slot uint64) *dag.DagNode {
	if block.Slot == slot {
		return block
	} else if block.Slot < slot {
		return nil
	} else {
		return gh.getAncestor(block.Parent, slot)
	}
}
