package team5

import (
	"fmt"
	"infra/game/agent"
	"infra/game/commons"
	"infra/logging"
	"math"
	"strings"
)

type AgentTrusts struct {
	StrategyScore float64
	GoodwillScore float64
}

func InitTrust() AgentTrusts {
	return AgentTrusts{
		StrategyScore: 0.5,
		GoodwillScore: 0.5,
	}
}

type Strategy uint
type Goodwill uint

type AgentProfile struct {
	Trusts   AgentTrusts
	Strategy Strategy
	Goodwill Goodwill
}

func (ap *AgentProfile) JSON() string {
	return fmt.Sprintf("{\"StrategyScore\": %f, \"GoodwillScore\": %f, \"Strategy\": %v, \"Goodwill\": %v}", ap.Trusts.StrategyScore, ap.Trusts.GoodwillScore, ap.Strategy, ap.Goodwill)
}

func InitAgentProfile() AgentProfile {
	return AgentProfile{
		Trusts:   InitTrust(),
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
	Initilised   bool
	AgentProfile map[commons.ID]AgentProfile
	LawfullMin   float64
	ChaoticMax   float64
	GoodMin      float64
	EvilMax      float64
}

func (sn *SocialNetwork) Log(id commons.ID, level uint) {
	logs := logging.LogField{}
	logs["ID"] = id
	logs["LEVEL"] = level
	network := []string{}
	for id, ap := range sn.AgentProfile {
		network = append(
			network,
			fmt.Sprintf("\"%s\": {\"StrategyScore\": %f, \"GoodwillScore\": %f, \"Strategy\": %v, \"Goodwill\": %v}",
				id, ap.Trusts.StrategyScore, ap.Trusts.GoodwillScore, ap.Strategy, ap.Goodwill),
		)
	}
	logs["SocialNetwork"] = "{" + strings.Join(network, ",") + "}"
	logging.Log(logging.Trace, logs, "TEAM5.SocialNetwork")
}

func (sn *SocialNetwork) InitSocialNetwork(ba agent.BaseAgent) {
	if sn.Initilised {
		return
	}
	view := ba.View()
	agentState := view.AgentState()

	agentprofileMAP := make(map[commons.ID]AgentProfile)
	itr := agentState.Iterator()
	for !itr.Done() {
		id, _, _ := itr.Next()
		agentprofileMAP[id] = InitAgentProfile()
	}

	sn.Initilised = true
	sn.AgentProfile = agentprofileMAP
	sn.LawfullMin = 0.8
	sn.ChaoticMax = 0.2
	sn.GoodMin = 0.8
	sn.EvilMax = 0.2
}

func (sn *SocialNetwork) UpdatePersonality(agentID commons.ID, extraStrategeScore float64, extraGoodwillScore float64) {
	ap := sn.AgentProfile[agentID]

	ap.Trusts.StrategyScore += extraStrategeScore
	ap.Trusts.GoodwillScore += extraGoodwillScore
	sn.normaliseTrust()

	s := ap.Trusts.StrategyScore
	switch {
	case s <= sn.ChaoticMax:
		ap.Strategy = Chaotic
	case s >= sn.LawfullMin:
		ap.Strategy = Lawful
	default:
		ap.Strategy = StrategyNeutral
	}

	g := ap.Trusts.GoodwillScore
	switch {
	case g <= sn.EvilMax:
		ap.Goodwill = Evil
	case g >= sn.GoodMin:
		ap.Goodwill = Good
	default:
		ap.Goodwill = GoodwillNeutral
	}

	sn.AgentProfile[agentID] = ap
}

func (sn *SocialNetwork) normaliseTrust() {
	minSTG := 0.5
	maxSTG := 0.5
	minGW := 0.5
	maxGW := 0.5

	for id := range sn.AgentProfile {
		g := sn.AgentProfile[id].Trusts.GoodwillScore
		minGW = math.Min(g, minGW)
		maxGW = math.Max(g, maxGW)

		s := sn.AgentProfile[id].Trusts.StrategyScore
		minSTG = math.Min(s, minSTG)
		maxSTG = math.Max(s, maxSTG)
	}

	distanceGW := maxGW - minGW
	distanceSTG := maxSTG - minSTG

	if distanceGW > 1 {
		for id, ap := range sn.AgentProfile {
			ap.Trusts.GoodwillScore = (ap.Trusts.GoodwillScore - minGW) / distanceGW
			sn.AgentProfile[id] = ap
		}
	}
	if distanceSTG > 1 {
		for id, ap := range sn.AgentProfile {
			ap.Trusts.StrategyScore = (ap.Trusts.StrategyScore - minSTG) / distanceSTG
			sn.AgentProfile[id] = ap
		}
	}
}
