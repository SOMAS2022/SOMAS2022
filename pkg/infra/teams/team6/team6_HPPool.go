package team6

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
)

func (a *Team6Agent) DonateToHpPool(baseAgent agent.BaseAgent) uint {
	/*
		Essentially want to minimise risk of over-donating, whilst recognising
		that donating is for the good of the game as donating enough could see
		a whole level being skipped.

		We know that if this agent thinks it should cower, it will...
		If it cowers, it can afford to give more HP - bearing in mind that when
		it attacks it will lose HP proportional to bonus attack
	*/
	donationPercentage := 25
	state := baseAgent.AgentState()

	if a.FightActionNoProposal(baseAgent) == decision.Cower {
		return state.Hp * uint(donationPercentage) / 100
	} else {
		return commons.SaturatingSub(state.Hp*uint(donationPercentage)/100, Max(state.Stamina, Max(state.TotalAttack(), state.TotalDefense())))
	}
}
