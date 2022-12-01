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
	"infra/game/state"
	"infra/game/tally"
	t0 "infra/teams/team0"

	"github.com/benbjohnson/immutable"
)

// ? Changed at compile time. eg change in .env to `MODE=0` to set this to '0'.
var Mode string

func ChooseDefaultStrategyMap(defaultStrategyMap map[commons.ID]func() agent.Strategy) map[commons.ID]func() agent.Strategy {
	switch Mode {
	case "0":
		return t0.InitAgentMap
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
	case "0":
		return t0.InitAgents(defaultStrategyMap, gameConfig, ptr)
	default:
		return initialise.InitAgents(defaultStrategyMap, gameConfig, ptr)
	}
}

// TODO: Change to using views.
func AgentLootDecisions(globalState state.State, agents map[commons.ID]agent.Agent, weaponLoot []uint, shieldLoot []uint) (allocatedState state.State) {
	switch Mode {
	case "0":
		return t0.AllocateLoot(globalState, weaponLoot, shieldLoot)
	default:
		return loot.AllocateLoot(globalState, weaponLoot, shieldLoot)
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
