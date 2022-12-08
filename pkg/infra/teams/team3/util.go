package team3

import (
	"infra/game/agent"
	"math/rand"
)

// Update Utility
func (a *AgentThree) UpdateUtility(baseAgent agent.BaseAgent) {
	view := baseAgent.View()
	agentState := view.AgentState()
	itr := agentState.Iterator()
	for !itr.Done() {
		id, _, ok := itr.Next()
		if !ok {
			break
		}

		// We are already cool, don't need the utility score for ourselves
		if id != baseAgent.ID() {
			a.utilityScore[id] = rand.Intn(10)
		}
	}
	// Sort utility map
	// sort.Sort(a.utilityScore)
}
