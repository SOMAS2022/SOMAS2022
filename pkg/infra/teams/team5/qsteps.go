package team5

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/state"
	"math/rand"
)

func (fiv *FivAgent) CurrentHPState(mystate state.AgentState) string {
	initHealth := 1000.0
	myHealth := ""
	switch {
	case mystate.Hp < uint(0.3*initHealth):
		myHealth = "Low"
	case uint(0.3*initHealth) <= mystate.Hp && mystate.Hp < uint(0.6*initHealth):
		myHealth = "Mid"
	case uint(0.6*initHealth) <= mystate.Hp:
		myHealth = "Hih"
	}
	return myHealth
}

func (fiv *FivAgent) CurrentAPState(mystate state.AgentState) string {
	initStamina := 2000.0
	myStamina := ""
	switch {
	case mystate.Stamina < uint(0.3*initStamina):
		myStamina = "Low"
	case uint(0.3*initStamina) <= mystate.Stamina && mystate.Stamina < uint(0.6*initStamina):
		myStamina = "Mid"
	case uint(0.6*initStamina) <= mystate.Stamina:
		myStamina = "Hih"
	}
	return myStamina
}

func (fiv *FivAgent) CurrentATState(popATGreaterToCount float32, numAlive float32) string {
	relativeAT := ""
	switch {
	case popATGreaterToCount < 0.25*numAlive:
		relativeAT = "Weakee"
	case 0.25*numAlive <= popATGreaterToCount && popATGreaterToCount < 0.75*numAlive:
		relativeAT = "Ordina"
	case 0.75 <= popATGreaterToCount && popATGreaterToCount <= numAlive:
		relativeAT = "Master"
	}
	return relativeAT
}

func (fiv *FivAgent) CurrentSHState(popSHGreaterToCount float32, numAlive float32) string {
	relativeSH := ""
	switch {
	case popSHGreaterToCount < 0.25*numAlive:
		relativeSH = "Weakee"
	case 0.25*numAlive <= popSHGreaterToCount && popSHGreaterToCount < 0.75*numAlive:
		relativeSH = "Ordina"
	case 0.75 <= popSHGreaterToCount && popSHGreaterToCount <= numAlive:
		relativeSH = "Master"
	}
	return relativeSH
}

func (fiv *FivAgent) CurrentQState(baseAgent agent.BaseAgent) string {
	mystate := baseAgent.AgentState()
	myview := baseAgent.View()

	var numAlive float32 = 0.0
	var popATGreaterToCount float32 = 0.0
	var popSHGreaterToCount float32 = 0.0

	globalStates := myview.AgentState()
	for _, id := range commons.ImmutableMapKeys(globalStates) {
		agState, _ := globalStates.Get(id)
		numAlive += 1
		if id != baseAgent.ID() {
			if agState.Attack+agState.BonusAttack < mystate.TotalAttack() {
				popATGreaterToCount += 1
			}
			if agState.Defense+agState.BonusDefense < mystate.TotalDefense() {
				popSHGreaterToCount += 1
			}
		}
	}

	return fiv.CurrentHPState(mystate) + "-" + fiv.CurrentAPState(mystate) + "-" + fiv.CurrentATState(popATGreaterToCount, numAlive) + "-" + fiv.CurrentSHState(popSHGreaterToCount, numAlive)
}

func (fiv *FivAgent) Explore(qstate string) decision.FightAction {
	var sa SaPair
	var fightDecision decision.FightAction
	fight := rand.Intn(3)
	switch fight {
	case 0:
		sa = SaPair{state: qstate, action: "Cower"}
		fightDecision = decision.Cower
	case 1:
		sa = SaPair{state: qstate, action: "Attck"}
		fightDecision = decision.Attack
	default:
		sa = SaPair{state: qstate, action: "Defnd"}
		fightDecision = decision.Defend
	}
	_, exist := fiv.qtable.table[sa]
	if !exist {
		fiv.qtable.table[sa] = 0
	}
	fiv.qtable.saTaken = sa
	return fightDecision
}

func (fiv *FivAgent) Exploit(qstate string) decision.FightAction {
	maxQAction := fiv.qtable.GetMaxQAction(qstate)
	var sa SaPair
	var fightDecision decision.FightAction
	switch maxQAction {
	case "NoSaPairAvailable":
		return fiv.Explore(qstate)
	case "Cower":
		sa = SaPair{state: qstate, action: "Cower"}
		fightDecision = decision.Cower
	case "Attck":
		sa = SaPair{state: qstate, action: "Attck"}
		fightDecision = decision.Attack
	case "Defnd":
		sa = SaPair{state: qstate, action: "Defnd"}
		fightDecision = decision.Defend
	}
	fiv.qtable.saTaken = sa
	return fightDecision
}

func (fiv *FivAgent) UpdateQ(baseAgent agent.BaseAgent) {
	// impact from self
	percentHealthLoss := (float32(fiv.preHealth) - float32(baseAgent.AgentState().Hp)) / float32(fiv.preHealth)
	// impact from group
	myview := baseAgent.View()
	globalStates := myview.AgentState()
	var guilt float32 = -1.0
	if fiv.qtable.saTaken.action == "Cower" {
		guilt = float32(globalStates.Len()) - float32(fiv.prePopNum)
	}
	// combined reward
	reward := -10*percentHealthLoss - 5*guilt

	cqState := fiv.CurrentQState(baseAgent)
	fSas := []SaPair{{state: cqState, action: "Cower"}, {state: cqState, action: "Attck"}, {state: cqState, action: "Defnd"}}
	fiv.qtable.Learn(reward, fiv.qtable.GetMaxFR(fSas))
	// fiv.qtable.Print()
}
