package team3

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/logging"
	"os"
	"strconv"

	"github.com/benbjohnson/immutable"
)

type AgentThree struct {
	AT                    int
	SH                    int
	uR                    map[commons.ID]int
	uP                    map[commons.ID]int
	uC                    map[commons.ID]int
	utilityScore          map[commons.ID]int
	TSN                   []commons.ID
	contactsLastRound     map[commons.ID]bool
	chairTolerance        int
	proposalTolerance     map[commons.ID]int
	fightDecisionsHistory commons.ImmutableList[decision.ImmutableFightResult]
	disobedience          int
}

// Update internal parameters at the end of each stage
func (a *AgentThree) UpdateInternalState(baseAgent agent.BaseAgent, _ *commons.ImmutableList[decision.ImmutableFightResult], votes *immutable.Map[decision.Intent, uint], log chan<- logging.AgentLog) {
	AS := baseAgent.AgentState()
	view := baseAgent.View()
	// Initialise utils
	if view.CurrentLevel() == 1 {
		a.utilityScore = a.InitUtility(baseAgent)
		a.uR = a.InitUtility(baseAgent)
		a.uP = a.InitUtility(baseAgent)
		a.uC = a.InitUtility(baseAgent)
	}

	a.AT = int(AS.Attack + AS.BonusAttack())
	a.SH = int(AS.Defense + AS.BonusDefense())

	// a.fightDecisionsHistory = *history

	a.UpdateTotalUtility(baseAgent)
	a.ResetContacts()
	a.UpdateTSN(baseAgent)

	// Potions
	if int(AS.Hp) < int(0.5*float64(GetStartingHP())) {
		// Drink HP Potion
	}
	if int(AS.Stamina) < int(0.5*float64(GetStartingStamina())) {
		// Drink Stamina Potion
	}
}

func NewAgentThreeNeutral() agent.Strategy {
	dis, _ := strconv.ParseInt(os.Getenv("NEUT_DIS"), 10, 0)
	return &AgentThree{
		utilityScore:      CreateUtility(),
		uR:                CreateUtility(),
		uP:                CreateUtility(),
		uC:                CreateUtility(),
		chairTolerance:    0,
		proposalTolerance: make(map[commons.ID]int, 0),
		disobedience:      int(dis),
	}
}
func NewAgentThreePassive() agent.Strategy {
	dis, _ := strconv.ParseInt(os.Getenv("PAS_DIS"), 10, 0)
	return &AgentThree{
		utilityScore:      CreateUtility(),
		uR:                CreateUtility(),
		uP:                CreateUtility(),
		uC:                CreateUtility(),
		chairTolerance:    0,
		proposalTolerance: make(map[commons.ID]int, 0),
		disobedience:      int(dis),
	}
}
func NewAgentThreeAggressive() agent.Strategy {
	dis, _ := strconv.ParseInt(os.Getenv("AGR_DIS"), 10, 0)
	return &AgentThree{
		utilityScore:      CreateUtility(),
		uR:                CreateUtility(),
		uP:                CreateUtility(),
		uC:                CreateUtility(),
		chairTolerance:    0,
		proposalTolerance: make(map[commons.ID]int, 0),
		disobedience:      int(dis),
	}
}
