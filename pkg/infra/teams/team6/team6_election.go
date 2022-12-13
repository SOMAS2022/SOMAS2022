package team6

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"math/rand"
)

func (a *Team6Agent) CreateManifesto(ba agent.BaseAgent) *decision.Manifesto {
	view := ba.View()

	fightDecisionPower, lootDecisionPower := false, false

	var leaderFightDecision, leaderLootDecision uint

	if view.LeaderManifesto().FightDecisionPower() {
		leaderFightDecision = 100
	} else {
		leaderFightDecision = 0
	}

	if view.LeaderManifesto().LootDecisionPower() {
		leaderLootDecision = 100
	} else {
		leaderLootDecision = 0
	}

	a.fightDecisionPowerOpinion += Min(commons.SaturatingSub(a.leadership[view.CurrentLeader()], 50)*uint(a.fightDecisionPowerOpinion*leaderFightDecision/100), 100)
	a.lootDecisionPowerOpinion += Min(commons.SaturatingSub(a.leadership[view.CurrentLeader()], 50)*uint(a.lootDecisionPowerOpinion*leaderLootDecision/100), 100)
	a.termLengthOpinion += Min(float32(commons.SaturatingSub(a.leadership[view.CurrentLeader()], 50)*uint(a.termLengthOpinion*float32(view.LeaderManifesto().TermLength())/100)), 100)
	a.overthrowTHOpinion += Min(float32(commons.SaturatingSub(a.leadership[view.CurrentLeader()], 50)*uint(a.overthrowTHOpinion*float32(view.LeaderManifesto().OverthrowThreshold())/100)), 100)

	if a.fightDecisionPowerOpinion > 50 {
		fightDecisionPower = true
	}

	if a.lootDecisionPowerOpinion > 50 {
		lootDecisionPower = true
	}

	a.proposedManifesto = *decision.NewManifesto(fightDecisionPower, lootDecisionPower, uint(a.termLengthOpinion), uint(a.overthrowTHOpinion))

	return &a.proposedManifesto
}

func (a *Team6Agent) HandleConfidencePoll(_ agent.BaseAgent) decision.Intent {
	switch rand.Intn(3) {
	case 0:
		return decision.Abstain
	case 1:
		return decision.Negative
	default:
		return decision.Positive
	}
}

func (a *Team6Agent) HandleElectionBallot(b agent.BaseAgent, _ *decision.ElectionParams) decision.Ballot {
	// Extract ID of alive agents
	view := b.View()
	agentState := view.AgentState()
	aliveAgentIDs := make([]string, agentState.Len())
	i := 0
	itr := agentState.Iterator()
	for !itr.Done() {
		id, a, ok := itr.Next()
		if ok && a.Hp > 0 {
			aliveAgentIDs[i] = id
			i++
		}
	}

	// Randomly fill the ballot
	var ballot decision.Ballot
	numAliveAgents := len(aliveAgentIDs)
	numCandidate := rand.Intn(numAliveAgents)
	for i := 0; i < numCandidate; i++ {
		randomIdx := rand.Intn(numAliveAgents)
		randomCandidate := aliveAgentIDs[uint(randomIdx)]
		ballot = append(ballot, randomCandidate)
	}

	return ballot
}
