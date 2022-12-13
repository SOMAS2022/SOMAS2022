package socialNetwork

import "infra/game/commons"

type agentTrusts struct {
	StrategyScore uint
	GoodwillScore uint
}

type agentPersonality struct {
	AgentID     commons.ID
	Trusts      agentTrusts
	Personality string
}


//provisional enum method, maybe helps with screen personalities in the applicaitons
//  type strategy uint
//  type goodwill uint

// const (
// 	Lawful strategy = iota
// 	StrategyNeutral
// 	Chaotic
// )

// const (
// 	Good goodwill = iota
// 	GoodwillNeutral
// 	Evil
// )

type SocialNetwork struct {
	AgentPersonality map[commons.ID]agentPersonality
	LawfullMin       uint
	ChaoticMax       uint
	GoodMin          uint
	EvilMax          uint
}

func updatePersonality(sn &SocialNetwork, agentID commons.ID, extraStrategeScore uint, extraGoodwillScore uint) {
	newStrategyScore := sn.AgentPersonality[agentID].Trusts.StrategyScore + extraStrategeScore
	sn.AgentPersonality[agentID].Trusts.StrategyScore = newStrategyScore
	newGoodwillScore := sn.AgentPersonality[agentID].Trusts.GoodwillScore + extraGoodwillScore
	sn.AgentPersonality[agentID].Trusts.GoodwillScore = newGoodwillScore
	
	updatePersonalityBoundaries(sn &SocialNetwork)

	if sn.AgentPersonality[agentID].Trusts.StrategyScore <= ChaoticMax {
		goodwillPersonality := "Evil"
	} else if sn.AgentPersonality[agentID].Trusts.StrategyScore < LawfullMin {
		goodwillPersonality := "GoodwillNeutral"
	} else{
		goodwillPersonality := "Good"
	}

	if sn.AgentPersonality[agentID].Trusts.GoodwillScore <= ChaoticMax{
		strategyPersonality := "Chaotic"
	} else if (sn.AgentPersonality[agentID].Trusts.GoodwillScore < LawfullMin) {
		strategyPersonality := "StrategyNeutral"
	} else {
		strategyPersonality := "Lawful"
	}

	sn.AgentPersonality[agentID].Personality = strategyPersonality+goodwillPersonality
}

func updatePersonalityBoundaries(sn &SocialNetwork) {

}