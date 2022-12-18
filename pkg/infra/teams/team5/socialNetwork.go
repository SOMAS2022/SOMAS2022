package team5

import (
	"fmt"
	"infra/game/agent"
	"infra/game/commons"
	"infra/logging"
)

type AgentTrusts struct {
	StrategyScore float32
	GoodwillScore float32
}

func initTrust() AgentTrusts {
	return AgentTrusts{
		StrategyScore: 0.5,
		GoodwillScore: 0.5,
	}
}

type Strategy uint
type Goodwill uint

type AgentProfile struct {
	AgentID  commons.ID
	Trusts   AgentTrusts
	Strategy Strategy
	Goodwill Goodwill
}

func (ap *AgentProfile) JSON() string {
	return fmt.Sprintf("AgentProfile: {ID: %s, Strategy: %f, Goodwill: %f}", ap.AgentID, ap.Trusts.StrategyScore, ap.Trusts.GoodwillScore)
}

func initAgentProfile(AgentID commons.ID) AgentProfile {
	return AgentProfile{
		AgentID:  AgentID,
		Trusts:   initTrust(),
		Strategy: StrategyNeutral,
		Goodwill: GoodwillNeutral,
	}
}

const (
	Chaotic         Strategy = iota
	StrategyNeutral Strategy = iota
	Lawful          Strategy = iota
)

const (
	Evil            Goodwill = iota
	GoodwillNeutral Goodwill = iota
	Good            Goodwill = iota
)

type SocialNetwork struct {
	AgentProfile map[commons.ID]AgentProfile
	LawfullMin   float32
	ChaoticMax   float32
	GoodMin      float32
	EvilMax      float32
}

func (sn *SocialNetwork) Log(id commons.ID, level uint) {
	logs := logging.LogField{"ID": id, "level": level}
	for _, ap := range sn.AgentProfile {
		logs[ap.AgentID] = ap.JSON()
	}
	logging.Log(logging.Trace, logs, "TEAM5.SocialNetwork")
}

func InitSocialNetwork(ba agent.BaseAgent) SocialNetwork {
	view := ba.View()
	agentState := view.AgentState()

	agentprofileMAP := make(map[commons.ID]AgentProfile)
	itr := agentState.Iterator()
	for !itr.Done() {
		id, _, _ := itr.Next()
		agentprofileMAP[id] = initAgentProfile(id)
	}

	return SocialNetwork{
		AgentProfile: agentprofileMAP,
		LawfullMin:   0.8,
		ChaoticMax:   0.2,
		GoodMin:      0.8,
		EvilMax:      0.2,
	}
}

func (sn *SocialNetwork) UpdatePersonality(agentID commons.ID, extraStrategeScore float32, extraGoodwillScore float32) {
	agentProfile := sn.AgentProfile[agentID]
	agentProfile.Trusts.StrategyScore += extraStrategeScore
	agentProfile.Trusts.GoodwillScore += extraGoodwillScore
	sn.AgentProfile[agentID] = agentProfile

	sn.normaliseTrust()

	if sn.AgentProfile[agentID].Trusts.StrategyScore <= sn.ChaoticMax {
		agentProfile = AgentProfile{Strategy: Chaotic}
	} else if sn.AgentProfile[agentID].Trusts.StrategyScore >= sn.LawfullMin {
		agentProfile = AgentProfile{Strategy: Lawful}
	} else {
		agentProfile = AgentProfile{Strategy: StrategyNeutral}
	}

	if sn.AgentProfile[agentID].Trusts.GoodwillScore <= sn.EvilMax {
		agentProfile = AgentProfile{Goodwill: Evil}
	} else if sn.AgentProfile[agentID].Trusts.GoodwillScore >= sn.GoodMin {
		agentProfile = AgentProfile{Goodwill: Good}
	} else {
		agentProfile = AgentProfile{Goodwill: GoodwillNeutral}
	}

	sn.AgentProfile[agentID] = agentProfile
}

func (sn *SocialNetwork) normaliseTrust() {
	var minSTG float32 = 0.5
	var maxSTG float32 = 0.5
	var minGW float32 = 0.5
	var maxGW float32 = 0.5
	//var id commons.ID

	for id := range sn.AgentProfile {
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

	if distanceGW > 1 {
		for id := range sn.AgentProfile {
			agentProfile := sn.AgentProfile[id]
			agentProfile.Trusts.GoodwillScore = (sn.AgentProfile[id].Trusts.GoodwillScore - minGW) / distanceGW
		}
	}
	if distanceSTG > 1 {
		for id := range sn.AgentProfile {
			agentProfile := sn.AgentProfile[id]
			agentProfile.Trusts.StrategyScore = (sn.AgentProfile[id].Trusts.StrategyScore - minSTG) / distanceSTG
		}
	}
}
