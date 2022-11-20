package stages

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/stage/fight"
	"infra/game/stage/loot"
	"infra/game/state"

	"github.com/benbjohnson/immutable"

	//? Add you team folder like this:
	t0 "infra/teams/team0"
)

// ? Changed at compile time. eg run with `make TEAM=0` to set this to '0'
var mode = "default"

// TODO: Change to using views
func AgentLootDecisions(globalState state.State, agents map[commons.ID]agent.Agent, weaponLoot []uint, shieldLoot []uint) (allocatedState state.State) {
	switch mode {
	case "0":
		return t0.AllocateLoot(globalState, weaponLoot, shieldLoot)
	default:
		return loot.AllocateLoot(globalState, weaponLoot, shieldLoot)
	}
}

func AgentFightDecisions(state *state.View, agents map[commons.ID]agent.Agent, previousDecisions immutable.Map[commons.ID, decision.FightAction], channelsMap map[commons.ID]chan message.TaggedMessage) map[string]decision.FightAction {
	switch mode {
	case "0":
		//? Not necessary to use all function arguments
		return t0.AllDefend(agents)
	default:
		return fight.AgentFightDecisions(state, agents, previousDecisions, channelsMap)
	}
}
