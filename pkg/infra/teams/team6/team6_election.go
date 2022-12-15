package team6

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"math"
)

func (a *Team6Agent) CreateManifesto(ba agent.BaseAgent) *decision.Manifesto {
	view := ba.View()

	fightDecisionPower, lootDecisionPower := false, false

	a.fightDecisionPowerOpinion = a.newPowerOpinion(a.fightDecisionPowerOpinion, a.leadership[view.CurrentLeader()], view.LeaderManifesto().FightDecisionPower())
	a.lootDecisionPowerOpinion = a.newPowerOpinion(a.lootDecisionPowerOpinion, a.leadership[view.CurrentLeader()], view.LeaderManifesto().LootDecisionPower())
	a.termLengthOpinion += (float32(a.leadership[view.CurrentLeader()]) - 50.) * (float32(view.LeaderManifesto().TermLength()) - a.termLengthOpinion) / 100.
	a.termLengthOpinion = clamp(a.termLengthOpinion, 1, 10)
	a.overthrowTHOpinion += (float32(a.leadership[view.CurrentLeader()]) - 50.) * (float32(view.LeaderManifesto().OverthrowThreshold()) - a.overthrowTHOpinion) / 100.
	a.overthrowTHOpinion = clamp(a.overthrowTHOpinion, 5, 75)

	if a.fightDecisionPowerOpinion > 50 {
		fightDecisionPower = true
	}

	if a.lootDecisionPowerOpinion > 50 {
		lootDecisionPower = true
	}

	a.proposedManifesto = *decision.NewManifesto(fightDecisionPower, lootDecisionPower, uint(a.termLengthOpinion), uint(a.overthrowTHOpinion))

	return &a.proposedManifesto
}

func (a *Team6Agent) HandleConfidencePoll(b agent.BaseAgent) decision.Intent {
	view := b.View()
	fightDecisionPower, lootDecisionPower := view.LeaderManifesto().FightDecisionPower(), view.LeaderManifesto().LootDecisionPower()
	fightDecisionValue, lootDecisionValue := 0, 0

	if fightDecisionPower {
		fightDecisionValue = 100
	}
	if lootDecisionPower {
		lootDecisionValue = 100
	}

	score := 1 / float32((4+math.Abs(float64(fightDecisionValue)-float64(a.fightDecisionPowerOpinion)))+
		1/(4+math.Abs(float64(lootDecisionValue)-float64(a.lootDecisionPowerOpinion)))+
		1/(4+math.Abs(float64(view.LeaderManifesto().TermLength())-float64(a.termLengthOpinion)))+
		1/math.Abs(float64(view.LeaderManifesto().OverthrowThreshold())-float64(a.overthrowTHOpinion)))

	length := len(a.agentsRemaining)
	score *= 1 / (1 + 0.25*float32(a.agentsRemaining[length-1]-a.agentsRemaining[length-2]))
	if score < 45 {
		return decision.Negative
	} else if score > 55 {
		return decision.Positive
	} else {
		return decision.Abstain
	}
}

func (a *Team6Agent) HandleElectionBallot(b agent.BaseAgent, params *decision.ElectionParams) decision.Ballot {
	// Extract ID of alive agents
	a.CreateManifesto(b)

	potentialCandidates := make(map[commons.ID]float32)

	itr := params.CandidateList().Iterator()
	for !itr.Done() {
		id, manifesto, _ := itr.Next()
		if manifesto.TermLength() > 10 || manifesto.OverthrowThreshold() > 75 {
			continue
		}

		fightDecisionPower, lootDecisionPower := manifesto.FightDecisionPower(), manifesto.LootDecisionPower()
		fightDecisionValue, lootDecisionValue := 0, 0

		if fightDecisionPower {
			fightDecisionValue = 100
		}
		if lootDecisionPower {
			lootDecisionValue = 100
		}

		potentialCandidates[id] = 1 / float32((4+math.Abs(float64(fightDecisionValue)-float64(a.fightDecisionPowerOpinion)))+
			1/(4+math.Abs(float64(lootDecisionValue)-float64(a.lootDecisionPowerOpinion)))+
			1/(4+math.Abs(float64(manifesto.TermLength())-float64(a.termLengthOpinion)))+
			1/math.Abs(float64(manifesto.OverthrowThreshold())-float64(a.overthrowTHOpinion)))
	}

	for id, score := range potentialCandidates {
		if SafeMapReadOrDefault(a.leadership, id, 50) < 40 {
			potentialCandidates[id] = 0
			continue
		}

		potentialCandidates[id] = score * (0.5*float32(SafeMapReadOrDefault(a.leadership, id, 50)) +
			0.2*float32(SafeMapReadOrDefault(a.similarity, id, 50)) +
			0.15*float32(SafeMapReadOrDefault(a.trust, id, 50)) +
			0.1*float32(SafeMapReadOrDefault(a.bravery, id, 50)) +
			0.05*float32(SafeMapReadOrDefault(a.generosity, id, 50)))
	}
	ballot := make([]commons.ID, 0)

	if params.Strategy() == decision.BordaCount {
		for id, score := range potentialCandidates {
			if score > 50 {
				ballot = append(ballot, id)
			}
		}
	} else {
		k, _ := FindMaxAgentInMap(potentialCandidates)
		ballot = append(ballot, k)
	}

	if len(ballot) == 0 {
		ballot = append(ballot, b.ID())
	}

	return ballot
}

func (a *Team6Agent) newPowerOpinion(initial uint, leadership uint, power bool) uint {
	powerValue := 0
	if power {
		powerValue = 100
	}
	fightDecisionDiff := (float32(leadership) - 50.) * ((float32(powerValue) - float32(a.fightDecisionPowerOpinion)) / 100.)
	if float32(initial)+fightDecisionDiff >= 100. {
		return uint(100)
	} else if float32(initial)+fightDecisionDiff <= 0. {
		return uint(0)
	} else {
		return uint(float32(initial) + fightDecisionDiff)
	}
}
