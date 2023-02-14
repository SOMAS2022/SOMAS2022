package stages

import (
	"infra/config"
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/stage/fight"
	"infra/game/stage/initialise"
	"infra/game/stage/loot"
	"infra/game/stage/update"
	"infra/game/state"
	"infra/game/tally"
	"infra/logging"
	"time"

	"github.com/benbjohnson/immutable"
	//? Add you team folder like this:
	// t0 "infra/teams/team0"
	// t1 "infra/teams/team1"
)

// Mode ? Changed at compile time. eg change in .env to `MODE=0` to set this to '0'.
var Mode string

func ChooseDefaultStrategyMap(defaultStrategyMap map[commons.ID]func() agent.Strategy) map[commons.ID]func() agent.Strategy {
	switch Mode {
	// case "0":
	// 	return t0.InitAgentMap
	// case "1":
	// 	return t1.InitAgentMap
	default:
		return defaultStrategyMap
	}
}

func InitGameConfig() config.GameConfig {
	switch Mode {
	case "0":
		return initialise.InitGameConfig() // ? Can choose to just call the default function
	default:
		return initialise.InitGameConfig()
	}
}

func InitAgents(defaultStrategyMap map[commons.ID]func() agent.Strategy, gameConfig config.GameConfig, ptr *state.View) (numAgents uint, agentMap map[commons.ID]agent.Agent, agentStateMap map[commons.ID]state.AgentState, inventoryMap state.InventoryMap) {
	switch Mode {
	// case "0":
	// 	return t0.InitAgents(defaultStrategyMap, gameConfig, ptr)
	// case "1":
	// 	return t1.InitAgents(defaultStrategyMap, gameConfig, ptr)
	default:
		return initialise.InitAgents(defaultStrategyMap, gameConfig, ptr)
	}
}

func AgentLootDecisions(globalState state.State, availableLoot state.LootPool, agents map[commons.ID]agent.Agent, channelsMap map[commons.ID]chan message.TaggedMessage) *tally.Tally[decision.LootAction] {
	switch Mode {
	default:
		return loot.AgentLootDecisions(globalState, availableLoot, agents, channelsMap)
	}
}

func AgentFightDecisions(state state.State, agents map[commons.ID]agent.Agent, previousDecisions immutable.Map[commons.ID, decision.FightAction], channelsMap map[commons.ID]chan message.TaggedMessage) *tally.Tally[decision.FightAction] {
	switch Mode {
	// case "0":
	// 	//? Not necessary to use all function arguments
	// 	return t0.AllDefend(agents)
	default:
		return fight.AgentFightDecisions(state, agents, previousDecisions, channelsMap)
	}
}

func UpdateInternalStates(agentMap map[commons.ID]agent.Agent, globalState *state.State, immutableFightRounds *commons.ImmutableList[decision.ImmutableFightResult], votesResult *immutable.Map[decision.Intent, uint]) map[commons.ID]logging.AgentLog {
	switch Mode {
	// case "1":
	// 	return t1.UpdateInternalStates(agentMap, globalState, immutableFightRounds, votesResult)
	default:
		return update.UpdateInternalStates(agentMap, globalState, immutableFightRounds, votesResult)
	}
}

func HandleTrustStage(agentMap map[commons.ID]agent.Agent, channelsMap map[commons.ID]chan message.TaggedMessage) {
	closures := make(map[commons.ID]chan<- struct{})

	// SEND ALL MESSAGES OUT
	for _, a := range agentMap {
		msg := a.Strategy.CompileTrustMessage(agentMap)
		senderList := msg.Recipients

		for _, ag := range senderList {
			// fmt.Println("SENDING:")
			if a.ID() == ag {
				continue
			}
			a.SendBlockingMessage(ag, msg)
		}
	}

	for id, a := range agentMap {
		a := a
		closure := make(chan struct{})
		closures[id] = closure

		go (&a).HandleTrust(closure)
	}

	// timeout for agents to respond
	time.Sleep(25 * time.Millisecond)
	for _, closure := range closures {
		closure <- struct{}{}
		close(closure)
	}

	for _, c := range channelsMap {
		close(c)
	}
}

func AgentPruneMapping(agentMap map[commons.ID]agent.Agent, globalState *state.State) map[commons.ID]agent.Agent {
	leaderId := globalState.CurrentLeader
	leader, leaderIsAlive := agentMap[leaderId]

	if leaderIsAlive {
		prunedMap := leader.PruneAgentList(agentMap)
		prunedMap[leaderId] = leader

		return prunedMap
	}
	// leader has died, hence no sanctioning
	return agentMap

}
