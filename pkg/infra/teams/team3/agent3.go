package team3

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/logging"

	"github.com/benbjohnson/immutable"
)

const PERCENTAGE = 500

type Utility struct {
	ID    commons.ID
	score int
}

type UtilityMap []Utility

func (u UtilityMap) Len() int           { return len(u) }
func (u UtilityMap) Swap(i, j int)      { u[i], u[j] = u[j], u[i] }
func (u UtilityMap) Less(i, j int) bool { return u[i].score < u[j].score }

type AgentThree struct {
	HP                    int
	ST                    int
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
}

// Update internal parameters at the end of each stage
func (a *AgentThree) UpdateInternalState(baseAgent agent.BaseAgent, history *commons.ImmutableList[decision.ImmutableFightResult], votes *immutable.Map[decision.Intent, uint], log chan<- logging.AgentLog) {

	AS := baseAgent.AgentState()

	a.HP = int(AS.Hp)
	a.ST = int(AS.Stamina)
	a.AT = int(AS.Attack + AS.BonusAttack())
	a.SH = int(AS.Defense + AS.BonusDefense())

	a.fightDecisionsHistory = *history

	a.UpdateTotalUtility(baseAgent)
	a.ResetContacts()
	a.UpdateTSN(baseAgent)
}

func CreateUtility() map[commons.ID]int {
	u := make(map[commons.ID]int, 7)
	return u
}

func (a *AgentThree) ResetContacts() {
	for i := range a.contactsLastRound {
		a.contactsLastRound[i] = false
	}
}

func NewAgentThree() agent.Strategy {
	return &AgentThree{
		utilityScore:      CreateUtility(),
		uR:                CreateUtility(),
		uP:                CreateUtility(),
		uC:                CreateUtility(),
		chairTolerance:    0,
		proposalTolerance: make(map[commons.ID]int, 0),
	}
}
