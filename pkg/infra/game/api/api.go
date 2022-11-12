package api

import (
	"infra/game/agent"
	"infra/game/decision"
	"infra/game/stage/fight"
	"infra/game/stage/loot"
	"infra/game/state"

	"github.com/benbjohnson/immutable"

	t0 "infra/teams/team0"
)

var mode = "default"

// TODO: Change to using views
func AllocateLoot(globalState state.State, weaponLoot []uint, shieldLoot []uint) (allocatedState state.State) {
	switch mode {
	case "0":
		return t0.AllocateLoot(globalState, weaponLoot, shieldLoot)
	default:
		return loot.AllocateLoot(globalState, weaponLoot, shieldLoot)
	}

}
func AgentFightDecisions(state *state.View, agents map[string]agent.Agent, previousDecisions *immutable.Map[string, decision.FightAction]) map[string]decision.FightAction {
	switch mode {
	case "0":
		return t0.AllDefend(agents)
	default:
		return fight.AgentFightDecisions(state, agents, previousDecisions)
	}
}
