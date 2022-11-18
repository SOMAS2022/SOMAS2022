package agent

import (
	//"fmt"
	"github.com/benbjohnson/immutable"
	"infra/game/decision"
	"infra/game/state"
)

type Team1MVP struct {
}

func (Team1MVP) HandleFight(gameState state.State, _ BaseAgent, decisionC chan<- decision.FightAction, log *immutable.Map[uint, decision.FightAction]) {

	//fmt.Println(gameState.AgentState[0].Hp)

	// Weights for linear regressors, l2-norm of each array should be 1
	attackWeights := [8]float64{0, 20, 50, 100, -20, 100, -20, 80}
	defendWeights := [8]float64{0, 20, 50, -20, 100, -20, 100, 80}
	cowerWeights := [8]float64{50, -80, -80, -30, -30, -30, -30, 100}

	// Reformat state into an array
	stateArray := [8]uint{
		gameState.CurrentLevel,
		gameState.AgentState[0].Hp,
		gameState.AgentState[0].Stamina,
		gameState.AgentState[0].Attack,
		gameState.AgentState[0].Defense,
		gameState.AgentState[0].BonusAttack,
		gameState.AgentState[0].BonusDefense,
		1, // Bias
	}

	// Multiply state with weights
	attack, defend, cower := 0.0, 0.0, 0.0
	for i := 0; i < len(attackWeights); i++ {
		attack += float64(attackWeights[i]) * float64(stateArray[i])
		defend += float64(defendWeights[i]) * float64(stateArray[i])
		cower += float64(cowerWeights[i]) * float64(stateArray[i])
	}

	// Choose action with highest regression value
	if attack >= defend && attack >= cower {
		decisionC <- decision.Attack
	} else if defend >= attack && defend >= cower {
		decisionC <- decision.Attack // decision.Defend
	} else {
		decisionC <- decision.Cower
	}

}
