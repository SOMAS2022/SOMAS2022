package update

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/state"
	"infra/logging"
	"sync"

	"github.com/benbjohnson/immutable"
)

func UpdateInternalStates(agentMap map[commons.ID]agent.Agent, globalState *state.State, immutableFightRounds *commons.ImmutableList[decision.ImmutableFightResult], votesResult *immutable.Map[decision.Intent, uint], logChan chan<- logging.AgentLog) {
	var wg sync.WaitGroup
	for id, a := range agentMap {
		id := id
		a := a
		wg.Add(1)
		go func(wait *sync.WaitGroup) {
			a.HandleUpdateInternalState(globalState.AgentState[id], immutableFightRounds, votesResult, logChan)
			wait.Done()
		}(&wg)
	}
	wg.Wait()
}
