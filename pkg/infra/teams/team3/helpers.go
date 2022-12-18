package team3

import (
	"infra/game/agent"
	"infra/game/commons"
	"math/rand"
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
func GetAttackAllAgents(baseAgent agent.BaseAgent) []float64 {
	view := baseAgent.View()
	agentState := view.AgentState()
	var agentStaminaMap []float64

	itr := agentState.Iterator()
	for !itr.Done() {
		_, state, _ := itr.Next()

		agentStaminaMap = append(agentStaminaMap, float64(state.Attack))
	}
	return agentStaminaMap
}

func GetDefenceAllAgents(baseAgent agent.BaseAgent) []float64 {
	view := baseAgent.View()
	agentState := view.AgentState()
	var agentStaminaMap []float64

	itr := agentState.Iterator()
	for !itr.Done() {
		_, state, _ := itr.Next()

		agentStaminaMap = append(agentStaminaMap, float64(state.Defense))
	}
	return agentStaminaMap
}

// func (a *AgentThree) FightTSN(agentMap *immutable.Map[commons.ID, state.HiddenAgentState]) {

// 	// for i, id := range a.TSN {

// 	// 	agentMap.Get

// 	// }

// }

func BordaPercentage(baseAgent agent.BaseAgent, borda [][]int) int {
	for i, v := range borda {
		if strconv.FormatInt(int64(v[0]), 10) == baseAgent.ID() {
			return (i / len(borda)) * 100
		}
	}
	return 100
}

func BoolToInt(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
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

func CreateUtility() map[commons.ID]int {
	u := make(map[commons.ID]int, 7)
	return u
}

func (a *AgentThree) InitUtility(baseAgent agent.BaseAgent) map[commons.ID]int {

	view := baseAgent.View()
	agentState := view.AgentState()
	itr := agentState.Iterator()

	u := make(map[commons.ID]int, 7)

	for !itr.Done() {
		id, _, _ := itr.Next()

		u[id] = rand.Intn(10)
	}
	return u
}

func (a *AgentThree) ResetContacts() {
	for i := range a.contactsLastRound {
		a.contactsLastRound[i] = false
	}
}

// Utility scores in range [0, 15]
func limitScore(score int) int {
	// Implement Boundaries
	if score > 15 {
		score = 15
	}
	if score < 0 {
		score = 0
	}
	return score
}
