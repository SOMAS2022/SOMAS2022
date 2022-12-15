package internal

import (
	"infra/game/state"
	"math"
)

func UtilityOfState(state state.AgentState) float64 {
	// Take geometric mean of state
	return math.Sqrt(float64(state.Hp)) * math.Sqrt(float64(state.Stamina)) * math.Sqrt(float64(state.TotalAttack())) * math.Sqrt(float64(state.TotalDefense()))
}

func UtilityOfHiddenState(state state.HiddenAgentState) float64 {
	// Take geometric mean of state
	return math.Sqrt(float64(state.Hp)) * math.Sqrt(float64(state.Stamina)) * math.Sqrt(float64(state.Attack+state.BonusAttack)) * math.Sqrt(float64(state.Defense+state.BonusDefense))
}
