package team6

import (
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message/proposal"
)

// TODO: Come up with a good proposal
func (a *Team6Agent) generateFightProposal() {
	rules := make([]proposal.Rule[decision.FightAction], 0)

	rules = append(rules, *proposal.NewRule(decision.Attack,
		proposal.NewAndCondition(*proposal.NewComparativeCondition(proposal.Health, proposal.GreaterThan, startingHP*uint(a.HPThreshold)-1),
			*proposal.NewComparativeCondition(proposal.Stamina, proposal.GreaterThan, startingST*uint(a.STThreshold)-1))))

	rules = append(rules, *proposal.NewRule(decision.Defend,
		proposal.NewComparativeCondition(proposal.TotalDefence, proposal.GreaterThan, 20)))

	rules = append(rules, *proposal.NewRule(decision.Cower,
		proposal.NewComparativeCondition(proposal.Health, proposal.LessThan, startingHP*uint(a.HPThreshold)+1)))

	a.fightProposal = *commons.NewImmutableList(rules)
}

// TODO: Come up with a good proposal
func (a *Team6Agent) generateLootProposal() {
	rules := make([]proposal.Rule[decision.LootAction], 0)

	rules = append(rules, *proposal.NewRule(decision.HealthPotion,
		proposal.NewComparativeCondition(proposal.Health, proposal.LessThan, startingHP*uint(a.HPThreshold)-1)))

	rules = append(rules, *proposal.NewRule(decision.StaminaPotion,
		proposal.NewComparativeCondition(proposal.Stamina, proposal.LessThan, startingHP*uint(a.HPThreshold)-1)))

	rules = append(rules, *proposal.NewRule(decision.Weapon,
		proposal.NewComparativeCondition(proposal.Health, proposal.LessThan, startingAT+1)))

	rules = append(rules, *proposal.NewRule(decision.Shield,
		proposal.NewComparativeCondition(proposal.Health, proposal.LessThan, startingSH+1)))

	a.lootProposal = *commons.NewImmutableList(rules)
}
