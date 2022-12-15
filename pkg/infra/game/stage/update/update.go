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

func UpdateInternalStates(agentMap map[commons.ID]agent.Agent, globalState *state.State, immutableFightRounds *commons.ImmutableList[decision.ImmutableFightResult], votesResult *immutable.Map[decision.Intent, uint]) map[commons.ID]logging.AgentLog {
	var wg sync.WaitGroup
	agentLogChan := make(chan logging.AgentLog)
	for id, a := range agentMap {
		id := id
		a := a
		wg.Add(1)
		go func(wait *sync.WaitGroup) {
			a.HandleUpdateInternalState(globalState.AgentState[id], immutableFightRounds, votesResult, agentLogChan)
			wait.Done()
		}(&wg)
	}

	agentLogs := make(map[commons.ID]logging.AgentLog)
	go func(agentLogChan chan logging.AgentLog, agentLogs map[commons.ID]logging.AgentLog) {
		for log := range agentLogChan {
			agentLogs[log.ID] = log
		}
	}(agentLogChan, agentLogs)
	wg.Wait()
	// fmt.Println(agentLogs)
	close(agentLogChan)
	// fmt.Println(agentLogs)
	return agentLogs
}
