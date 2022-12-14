package team5

import "infra/game/commons"

type AgentTrusts struct {
	StrategyScore float32
	GoodwillScore float32
}

type Personality struct {
	Strategy string
	Goodwill string
}

type AgentProfile struct {
	AgentID     commons.ID
	Trusts      AgentTrusts
	Personality Personality
}

//provisional enum method, maybe helps with screen personalities in the applicaitons
type strategy float32
type goodwill float32

const (
	Lawful strategy = iota
	StrategyNeutral
	Chaotic
)

const (
	Good goodwill = iota
	GoodwillNeutral
	Evil
)

type SocialNetwork struct {
	AgentProfile map[commons.ID]AgentProfile
	LawfullMin   float32
	ChaoticMax   float32
	GoodMin      float32
	EvilMax      float32
}

func (sn *SocialNetwork) updatePersonality(agentID commons.ID, extraStrategeScore float32, extraGoodwillScore float32) {

	agentProfile := sn.AgentProfile[agentID]
	agentProfile.Trusts.StrategyScore += extraStrategeScore
	agentProfile.Trusts.GoodwillScore += extraGoodwillScore
	sn.AgentProfile[agentID] = agentProfile

	sn.normaliseTrust()

	var goodwillPersonality, strategyPersonality string
	if sn.AgentProfile[agentID].Trusts.StrategyScore <= sn.ChaoticMax {
		goodwillPersonality = "Evil"
	} else if sn.AgentProfile[agentID].Trusts.StrategyScore < sn.LawfullMin {
		goodwillPersonality = "GoodwillNeutral"
	} else {
		goodwillPersonality = "Good"
	}

	if sn.AgentProfile[agentID].Trusts.GoodwillScore <= sn.ChaoticMax {
		strategyPersonality = "Chaotic"
	} else if sn.AgentProfile[agentID].Trusts.GoodwillScore < sn.LawfullMin {
		strategyPersonality = "StrategyNeutral"
	} else {
		strategyPersonality = "Lawful"
	}

	agentProfile.Personality = Personality{strategyPersonality, goodwillPersonality}
	sn.AgentProfile[agentID] = agentProfile
}

func (sn *SocialNetwork) normaliseTrust() {
	var minSTG float32 = 0.5
	var maxSTG float32 = 0.5
	var minGW float32 = 0.5
	var maxGW float32 = 0.5
	var id commons.ID

	for id = range sn.AgentProfile {
		if sn.AgentProfile[id].Trusts.GoodwillScore < minGW {
			minGW = sn.AgentProfile[id].Trusts.GoodwillScore
		}
		if sn.AgentProfile[id].Trusts.GoodwillScore > maxGW {
			maxGW = sn.AgentProfile[id].Trusts.GoodwillScore
		}
		if sn.AgentProfile[id].Trusts.GoodwillScore < minSTG {
			minSTG = sn.AgentProfile[id].Trusts.StrategyScore
		}
		if sn.AgentProfile[id].Trusts.GoodwillScore > maxSTG {
			maxSTG = sn.AgentProfile[id].Trusts.StrategyScore
		}
	}
	distanceGW := maxGW - minGW
	distanceSTG := maxSTG - minSTG
	for id = range sn.AgentProfile {
		agentProfile := sn.AgentProfile[id]
		agentProfile.Trusts.GoodwillScore = (sn.AgentProfile[id].Trusts.GoodwillScore - minGW) / distanceGW
		agentProfile.Trusts.StrategyScore = (sn.AgentProfile[id].Trusts.StrategyScore - minSTG) / distanceSTG
		sn.AgentProfile[id] = agentProfile
	}
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
