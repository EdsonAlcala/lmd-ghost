package dag

type ScoreChange struct {
	Target *DagNode
	ScoreDelta int64
}

type ForkChoice interface {
	OnNewNode(node *DagNode)
	ApplyScoreChanges(changes []ScoreChange)
	OnStartChange()
	HeadFn() *DagNode
}

type InitForkChoice func(dag *BeaconDag) ForkChoice
