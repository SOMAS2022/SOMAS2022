package team6

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/logging"

	"github.com/benbjohnson/immutable"
)

var startingHP = uint(1000)
var startingAT = uint(20)
var startingSH = uint(20)
var startingST = uint(2000)

type Team6Agent struct {
	bravery    map[commons.ID]uint
	generosity map[commons.ID]uint
	similarity map[commons.ID]uint
	trust      map[commons.ID]uint
	leadership map[commons.ID]uint

	HPThreshold float32
	ATThreshold float32
	SHThreshold float32
	STThreshold float32

	currentLevel             uint
	lastFightRound           uint
	lastHPPoolDonationAmount uint
}

func NewTeam6Agent() agent.Strategy {
	return &Team6Agent{
		bravery:     make(map[commons.ID]uint),
		generosity:  make(map[commons.ID]uint),
		similarity:  make(map[commons.ID]uint),
		trust:       make(map[commons.ID]uint),
		leadership:  make(map[commons.ID]uint),
		HPThreshold: 0.1,
		ATThreshold: 0.1,
		SHThreshold: 0.1,
		STThreshold: 0.1,
	}
}

func (a *Team6Agent) HandleUpdateWeapon(_ agent.BaseAgent) decision.ItemIdx {
	// weapons := b.AgentState().weapons
	// return decision.ItemIdx(rand.Intn(weapons.Len() + 1))

	// 0th weapon has the greatest attack points
	return decision.ItemIdx(0)
}

func (a *Team6Agent) HandleUpdateShield(_ agent.BaseAgent) decision.ItemIdx {
	// shields := b.AgentState().Shields
	// return decision.ItemIdx(rand.Intn(shields.Len() + 1))
	return decision.ItemIdx(0)
}

func (a *Team6Agent) UpdateInternalState(ba agent.BaseAgent, _ *commons.ImmutableList[decision.ImmutableFightResult], _ *immutable.Map[decision.Intent, uint], log chan<- logging.AgentLog) {
	log <- logging.AgentLog{
		Name: ba.Name(),
		ID:   ba.ID(),
		Properties: map[string]float32{
			"hp pool donation": float32(a.lastHPPoolDonationAmount),
			"hp":               float32(ba.AgentState().Hp),
			"lastFightRound":   float32(a.lastFightRound),
		},
	}
	a.lastFightRound++
	a.currentLevel++
}
