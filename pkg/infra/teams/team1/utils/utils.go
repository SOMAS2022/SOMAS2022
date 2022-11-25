package utils

import (
	"infra/config"
	"infra/game/state"
)

var Config config.GameConfig

/**
 * This function converts agents stats into a single 'usefullness-to-battle' value
 *
 * Calculated based on a(curr)
 */
func AgentBattleUtility(agentState state.AgentState) float32 {
	attackPortion := float32(agentState.TotalAttack() / Config.StartingAttackStrength)
	defensePortion := float32(agentState.TotalDefense() / Config.StartingShieldStrength)
	healthPortion := float32(agentState.Hp / Config.StartingHealthPoints)
	staminaPortion := float32(agentState.Stamina / Config.Stamina)

	return 0.25*attackPortion + 0.25*defensePortion + 0.25*healthPortion + 0.25*staminaPortion
}
