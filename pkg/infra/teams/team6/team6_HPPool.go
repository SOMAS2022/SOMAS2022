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
	donationPercentage := uint(5)
	donationMaximum := uint(0.10 * float32(startingHP))
	donationHPThreshold := uint(0.25 * float32(startingHP))
	state := baseAgent.AgentState()

	if a.lastFightRound == 0 {
		donationPercentage = 15
	} else {
		donationPercentage = 1
	}

	if state.Hp <= donationHPThreshold {
		a.lastHPPoolDonationAmount = 0
		return 0
	} else if a.FightActionNoProposal(baseAgent) == decision.Cower {
		a.lastHPPoolDonationAmount = uint(state.Hp * donationPercentage / 100)
		return a.lastHPPoolDonationAmount
	} else {
		if state.Stamina < Max(state.TotalAttack(), state.TotalDefense()) {
			expectedHPRemaining := commons.SaturatingSub(state.Hp, Max(state.Stamina, Max(state.TotalAttack(), state.TotalDefense())))
			a.lastHPPoolDonationAmount = Min(donationMaximum, expectedHPRemaining*uint(donationPercentage)/100)
			return a.lastHPPoolDonationAmount
		} else {
			expectedHPRemaining := commons.SaturatingSub(state.Hp, Max(state.TotalAttack(), state.TotalDefense()))
			a.lastHPPoolDonationAmount = Min(donationMaximum, expectedHPRemaining*uint(donationPercentage)/100)
			return a.lastHPPoolDonationAmount
		}
	}
}
