/*******************************************************
* Copyright (C) 2022 Team 1 @ SOMAS2022
*
* This file is part of SOMAS2022.
*
* This file or its contents can not be copied and/or used
* without the express permission of Team 1, SOMAS2022
*******************************************************/
package internal

import (
	"bufio"
	"infra/game/state"
	"log"
	"os"
	"strings"
)

type QState struct {
	Hp            float64
	Stamina       float64
	TotalAttack   float64
	TotalDefense  float64
	CurrLevel     float64
	MonsterHealth float64
	MonsterAttack float64
}

type ActionStrategy struct {
	LinRegWeights [8]float64
}

// Global variables for strategies?
// For each action + coop/self have a set of 8 weights
// var CoopStrategies = [3]ActionStrategy{
// 	{[8]float64{0, 0, 0, 0, 0, 0, 0, 0}}, // Defend
// 	{[8]float64{0, 0, 0, 0, 0, 0, 0, 0}}, // Cower
// 	{[8]float64{0, 0, 0, 0, 0, 0, 0, 0}}, // Attack
// }

// var SelfishStrategies = [3]ActionStrategy{
// 	{[8]float64{0, 0, 0, 0, 0, 0, 0, 0}}, // Defend
// 	{[8]float64{0, 0, 0, 0, 0, 0, 0, 0}}, // Cower
// 	{[8]float64{0, 0, 0, 0, 0, 0, 0, 0}}, // Attack
// }

func ReadStrategy() ([3]ActionStrategy, [3]ActionStrategy) {
	// Read strategy from final_weights.csv file
	var coopstrat [3]ActionStrategy
	var selfstrat [3]ActionStrategy

	// Check if logging is turned on
	if os.Getenv("QLOGGING") != "true" {
		// Logging is not turned on so return
		return coopstrat, selfstrat
	}

	// Open log file
	logFile, err := os.Open("/mnt/c/Users/alexp/OneDrive - Imperial College London/Year 4/Y4_SOMAS/SOMAS2022/pkg/infra/teams/team1/final_weights.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	// Create a new scanner on logFile
	scanner := bufio.NewScanner(logFile)

	// Read in strategy weights

	for scanner.Scan() {
		snapshotString := scanner.Text()

		// Read cooperative weights
		for action := 0; action < 3; action++ {
			snapshot := strings.Split(snapshotString, ",")
			if len(snapshot) == 8 {
				for w := 0; w < 8; w++ {
					coopstrat[action].LinRegWeights[w] = StringToFloat(snapshot[w])
				}
			}
		}

		// Read selfish weights
		for action := 0; action < 3; action++ {
			snapshot := strings.Split(snapshotString, ",")
			if len(snapshot) == 8 {
				for w := 0; w < 8; w++ {
					selfstrat[action].LinRegWeights[w] = StringToFloat(snapshot[w])
				}
			}
		}

	}
	logFile.Close()

	return coopstrat, selfstrat

}

// Dot product between weights and array
func computeReward(weights [8]float64, array [8]float64) float64 {
	res := 0.0
	for idx, w := range weights {
		res += w * array[idx]
	}
	return res
}

// Get QState of agent
func getQState(state state.AgentState) [8]float64 {
	// TODO currently using dummy variables - need world state as well
	return [8]float64{
		1, float64(state.Hp), float64(state.Stamina), float64(state.Attack), float64(state.Defense),
		1, 2, 3,
	}
}

// Get QState of other agent given their hidden state
func getQStateOther(state state.HiddenAgentState) [8]float64 {
	return [8]float64{
		float64(state.Hp), float64(state.Stamina), float64(state.Attack), float64(state.Defense),
		1, 2, float64(state.BonusAttack), float64(state.BonusDefense),
	}
}

// Output state -> reward (given strategy)
func QFunction(qstate QState, strategy [3]ActionStrategy) [3]float64 {
	// Update strategy
	var reward [3]float64
	for i := 0; i < 3; i++ {
		reward[i] = computeReward(strategy[i].LinRegWeights, QStateToArray(qstate))
	}
	return reward
}

// Learn weights of strategies given training data
// Called outside of games to pre-train weights
func QFunctionTrain(data [][]float64, obs []float64) []float64 {
	// Choose observation column depending on coop/self strategy
	// Process entire table of experiences to separate actions and rewards
	// Instead of separating data here, could just store these experiences into different logs?
	// Need to decide on best way to transfer information

	// Train all strategies

	// Cooperative strategies get mean number of levels left as reward

	// Selfish strategies get number of levels it lived

	return FitLinReg(data, obs)
}
