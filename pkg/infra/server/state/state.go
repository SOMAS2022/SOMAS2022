package state

type AgentState struct {
	hp     uint
	attack uint
	shield uint
}

type State struct {
	currentLevel uint
	hpPool       uint
	agentState   map[uint]AgentState
}
