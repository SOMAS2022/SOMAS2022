package team3

import (
	"infra/game/agent"
	"os"
	"strconv"
)

func GetHealthAllAgents(baseAgent agent.BaseAgent) []float64 {
	view := baseAgent.View()
	agentState := view.AgentState()
	var agentHealthMap []float64

	itr := agentState.Iterator()
	for !itr.Done() {
		_, state, _ := itr.Next()

		agentHealthMap = append(agentHealthMap, float64(state.Hp))

	}

	return agentHealthMap
}

func GetStaminaAllAgents(baseAgent agent.BaseAgent) []float64 {
	view := baseAgent.View()
	agentState := view.AgentState()
	var agentStaminaMap []float64

	itr := agentState.Iterator()
	for !itr.Done() {
		_, state, _ := itr.Next()

		agentStaminaMap = append(agentStaminaMap, float64(state.Stamina))

	}

	return agentStaminaMap

}

func AverageArray(in []float64) float64 {

	var total float64 = 0
	for _, value := range in {
		total += value
	}
	return total / float64(len(in))
}

func GetStartingHP() int {
	n, _ := strconv.ParseUint(os.Getenv("STARTING_HP"), 10, 0)
	return int(n)
}
func GetStartingStamina() int {
	n, _ := strconv.ParseUint(os.Getenv("BASE_STAMINA"), 10, 0)
	return int(n)
}
