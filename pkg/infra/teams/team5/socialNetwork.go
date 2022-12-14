package socialNetwork

import "infra/game/commons"

type agentTrusts struct {
	StrategyScore uint
	GoodwillScore uint
}

type agentProfile struct {
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
	AgentProfile map[commons.ID]agentProfile
	LawfullMin   uint
	ChaoticMax   uint
	GoodMin      uint
	EvilMax      uint
}

func updatePersonality(sn SocialNetwork, agentID commons.ID, extraStrategeScore uint, extraGoodwillScore uint) SocialNetwork {
	nsn := sn
	newStrategyScore := sn.AgentProfile[agentID].Trusts.StrategyScore + extraStrategeScore
	nsn.AgentProfile[agentID].Trusts.StrategyScore = newStrategyScore
	newGoodwillScore := sn.AgentProfile[agentID].Trusts.GoodwillScore + extraGoodwillScore
	nsn.AgentProfile[agentID].Trusts.GoodwillScore = newGoodwillScore

	nsn = updatePersonalityBoundaries(nsn)

	if nsn.AgentProfile[agentID].Trusts.StrategyScore <= nsn.ChaoticMax {
		goodwillPersonality := "Evil"
	} else if nsn.AgentProfile[agentID].Trusts.StrategyScore < nsn.LawfullMin {
		goodwillPersonality := "GoodwillNeutral"
	} else {
		goodwillPersonality := "Good"
	}

	if nsn.AgentProfile[agentID].Trusts.GoodwillScore <= nsn.ChaoticMax {
		strategyPersonality := "Chaotic"
	} else if nsn.AgentProfile[agentID].Trusts.GoodwillScore < nsn.LawfullMin {
		strategyPersonality := "StrategyNeutral"
	} else {
		strategyPersonality := "Lawful"
	}

	nsn.AgentProfile[agentID].Personality = strategyPersonality + goodwillPersonality
	return nsn
}

func updatePersonalityBoundaries(sn SocialNetwork) SocialNetwork {
	nsn := sn

	return nsn
}

//Initialise social network

// func initialiseSocialNetwork() socialNetwork {
// 	var initialTrust uint
// 	initialTrust = 50
// 	sn.AgentProfile[...commmons.ID].StrategyScore = initialTrust
// 	sn.AgentProfile.[...commmons.ID].AgentProfile.GoodwillScore = initialTrust
// 	sn := socialNetwork{
// 		AgentProfile: make(map[commons.ID]agentProfile),
// 		LawfullMin:       75,
// 		ChaoticMax:       25,
// 		GoodMin:          75,
// 		EvilMax:          25,
// 	}
// 	return sn
// }
