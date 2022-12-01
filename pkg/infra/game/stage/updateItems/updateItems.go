package updateItems

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/state"
)

func UpdateItems(state state.State, agents map[commons.ID]agent.Agent) state.State {
	updatedState := state
	for _, agent := range agents {
		agentState := updatedState.AgentState[agent.BaseAgent.ID()]
		weaponId := agent.HandleUpdateWeapon(agentState, state.ToView())
		shieldId := agent.HandleUpdateShield(agentState, state.ToView())
		agentState.ChangeWeaponInUse(weaponId)
		agentState.ChangeShieldInUse(shieldId)
		updatedState.AgentState[agent.BaseAgent.ID()] = agentState
	}
	return updatedState
}
