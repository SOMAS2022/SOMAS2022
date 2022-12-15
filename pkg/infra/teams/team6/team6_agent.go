package team6

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message/proposal"
	"infra/logging"

	"github.com/benbjohnson/immutable"
)

var startingHP = uint(1000)
var startingAT = uint(20)
var startingSH = uint(20)
var startingST = uint(2000)
var startingAgents = uint(100)

type Team6Agent struct {
	bravery    map[commons.ID]uint
	generosity map[commons.ID]uint
	similarity map[commons.ID]uint
	trust      map[commons.ID]uint
	leadership map[commons.ID]uint

	timesCowered map[commons.ID]uint

	HPThreshold float32
	ATThreshold float32
	SHThreshold float32
	STThreshold float32

	proposedManifesto         decision.Manifesto
	fightDecisionPowerOpinion uint
	lootDecisionPowerOpinion  uint
	termLengthOpinion         float32
	overthrowTHOpinion        float32

	currentLevel             uint
	lastFightRound           uint
	lastHPPoolDonationAmount uint
	totalFightRounds         uint
	agentsRemaining          []uint

	fightProposal commons.ImmutableList[proposal.Rule[decision.FightAction]]
	lootProposal  commons.ImmutableList[proposal.Rule[decision.LootAction]]
}

func NewTeam6Agent() agent.Strategy {
	return &Team6Agent{
		bravery:    make(map[commons.ID]uint),
		generosity: make(map[commons.ID]uint),
		similarity: make(map[commons.ID]uint),
		trust:      make(map[commons.ID]uint),
		leadership: make(map[commons.ID]uint),

		timesCowered: make(map[commons.ID]uint),

		HPThreshold: 0.1,
		ATThreshold: 0.1,
		SHThreshold: 0.1,
		STThreshold: 0.1,

		proposedManifesto:         *decision.NewManifesto(false, false, 3, 51),
		fightDecisionPowerOpinion: 25,
		lootDecisionPowerOpinion:  25,
		termLengthOpinion:         3.,
		overthrowTHOpinion:        51.,

		agentsRemaining: []uint{startingAgents, startingAgents},
	}
}

func (a *Team6Agent) HandleUpdateWeapon(ba agent.BaseAgent) decision.ItemIdx {
	state := ba.AgentState()
	weapons := state.Weapons

	itr := weapons.Iterator()
	for !itr.Done() {
		idx, value := itr.Next()
		if float32(value.Value()) < a.STThreshold*float32(startingST) {
			return decision.ItemIdx(idx)
		}
	}

	return decision.ItemIdx(Max(0, weapons.Len()-1))
}

func (a *Team6Agent) HandleUpdateShield(ba agent.BaseAgent) decision.ItemIdx {
	state := ba.AgentState()
	shields := state.Shields

	itr := shields.Iterator()
	for !itr.Done() {
		idx, value := itr.Next()
		if float32(value.Value()) < a.STThreshold*float32(startingST) {
			return decision.ItemIdx(idx)
		}
	}

	return decision.ItemIdx(Max(0, shields.Len()-1))
}

func (a *Team6Agent) UpdateInternalState(ba agent.BaseAgent, fightRounds *commons.ImmutableList[decision.ImmutableFightResult], _ *immutable.Map[decision.Intent, uint], log chan<- logging.AgentLog) {
	log <- logging.AgentLog{
		Name: ba.Name(),
		ID:   ba.ID(),
		Properties: map[string]float32{
			"hp pool donation": float32(a.lastHPPoolDonationAmount),
			"hp":               float32(ba.AgentState().Hp),
			"lastFightRound":   float32(a.lastFightRound),
		},
	}

	// Update number of times each agent has cowered
	itr2 := fightRounds.Iterator()
	for !itr2.Done() {
		a.totalFightRounds++
		fightResult, _ := itr2.Next()
		coweringAgents := fightResult.CoweringAgents()
		itr3 := coweringAgents.Iterator()
		for !itr3.Done() {
			_, id := itr3.Next()
			_, ok := a.timesCowered[id]
			if ok {
				a.timesCowered[id]++
			} else {
				a.trust[id] = 1
			}
		}
	}

	view := ba.View()
	agentStates := view.AgentState()
	itr := agentStates.Iterator()
	for !itr.Done() {
		id, as, _ := itr.Next()
		// Update trust
		if as.Defector.IsDefector() {
			_, ok := a.trust[id]
			if ok {
				a.trust[id] = commons.SaturatingSub(a.trust[id], 10)
			} else {
				a.trust[id] = 40
			}
		} else {
			_, ok := a.trust[id]
			if ok {
				a.trust[id] = SCSaturatingAdd(a.trust[id], 2, 100)
			} else {
				a.trust[id] = 50
			}
		}
		// Update bravery
		_, ok := a.timesCowered[id]
		if ok {
			a.bravery[id] = 100 - uint(float32(a.timesCowered[id])*100/float32(a.totalFightRounds))
		} else {
			a.bravery[id] = 100
		}
	}

	a.agentsRemaining = append(a.agentsRemaining, uint(agentStates.Len()))
	a.lastFightRound++ // No. of levels since last fight
	a.currentLevel++
}
