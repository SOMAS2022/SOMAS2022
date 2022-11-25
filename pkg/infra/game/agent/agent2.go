package agent

import (
	"github.com/benbjohnson/immutable"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/state"
	"math/rand"
)

type Agent2 struct {
	BaseAgent BaseAgent
	Strategy  Strategy
	bravery   int
}

func (a Agent2) Default() decision.FightAction {
	//TODO implement me
	panic("implement me")
}

func NewAgent2() *Agent2 {
	return &Agent2{bravery: 0}
}

func (a Agent2) HandleFightMessage(m message.TaggedMessage, view *state.View, agent BaseAgent, log *immutable.Map[commons.ID, decision.FightAction]) decision.FightAction {

	/*
		// assuming a decision was already made before, the following is ran inside a fight round loop to drink potions and change our decision if necessary

		agentState := view.AgentState()[a.BaseAgent.Id]; // our agent's current state
		// WARNING : not sure how go works exactly, this should be a reference so that we always point to the current value in the following statements


		// Potions
		if agentState.Hp < 25 && hasHPpotion {
			// drink health potion
		}

		=====================TASKS=======================
		!!!FUNCTIONS TO IMPLEMENT!!!
========================
		Tim:
		
		// assume the following function that estimates if the amount of other agents is enough to kill the monsters (so we can cower if we're low)
		func enoughOtherAgentsFighting(...)
		input args : global state map
		output : bool

		for example, remember the amount of agents that fought last round and how many damage was dealt, to get a mean damage/agent
		then estimate this round's damage based on the number of fighting agents WITH HEALTH NOT LOW (otherwise assume they'll be cowering to survive instead)

		
		// cower if low health and not critically needed
		if agentState.Hp < 15 {
			if enoughOtherAgentsFighting(...) {
				return decision.Cower
			}
		}
========================
		
========================
		Giorgos:
		
		1. Overleaf Document
		2. func History() {}
		3. func FreuencyOfDecisions() {}
=========================

=========================
		Vasos:
		
		// drink potions if attacking or defending
		if decision != decision.Cower {
			if hasAttackPotion && decision == decision.Attack {
				// drink it
			}
			if hasShieldPotion {
				// drink it
			}
		}
==========================
		
		
	*/

	/*
	   Fighting is similar to Random Agent for MVP.
	   - Agent attacks by default (50% of the time)
	   - Agent defends 40% of the time
	   - Agent cowers 10% of time
	*/
	fight := rand.Intn(10)
	switch {
	case fight == 0:
		return decision.Cower
	case (fight <= 4) && (fight > 0):
		return decision.Defend
	default:
		return decision.Attack
	}
}


