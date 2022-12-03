/*******************************************************
* Copyright (C) 2022 Team 1 @ SOMAS2022
*
* This file is part of SOMAS2022.
*
* This file or its contents can not be copied and/or used
* without the express permission of Team 1, SOMAS2022
*******************************************************/

package team1

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/teams/team1/utils"
	"math/rand"

	"github.com/benbjohnson/immutable"
	"github.com/google/uuid"
)

type SocialAgent struct {
	// probability of chosing collaborative strategy. p(selfish_strategy) is 1 - this
	pCollaborate float64

	// current fightDecision
	currentDecision decision.FightAction

	// Metadata
	battleUtility float64

	socialCapital    map[string]SocialCapInfo
	lastLevelUpdated uint // The level at which the socialCapital was last updated

	// Will gosip to all agents who's network value is above this
	gossipThreshold float64
	// Proportion of agents to talk badly about
	propHate float64
	// Proportion of agents to talk well about
	propAdmire float64

	selfID string

	// Four metrics for each agent's perception of other agents. Three metrics are borrowed from Ostrom-Ahn
	// social capital model: (The ordering below is the same as the ordering in the array)
	// 1. Institutions: Influenced by institutional actions such as voting. Also influenced by any elected
	//positions the agent currently holds.
	// 2. Networks: Messages received about an agent from other agents.
	// The impact which a message from another agent has on the network metric depends on the overall perception of
	// the agent sending the message.
	// 3. Trustworthiness: Measures how likely the agent is to cooperate based on historic decisions.
	// Finally, our addition to the model is
	// 4. Favours: A metric of how many favours an agent has done for another agent, compared to the amount of
	//favours the other agent has done for them. This metric is supposed to capture social expectations, such
	//as if I went out with the trash last time, then I feel that next time someone else should do it. The value
	// of the  favours metric will over time impact the trustworthiness metric of another agent. In many ways,
	// favours is a short time horizon equivalent of trustworthiness.

	// Forgiveness: In order to enable even the bitterest of enemies to start cooperating again in the long term,
	// the metrics will decay with time. However, negative metrics will decay far faster than positive metrics.
	// This is similar to how you can have a massive fight and hate someone one day, and then be friends again the
	// next. On the other hand, you can go years without seeing a friend and still have a very positive impression
	// of them.

	// Messages we should have: Denounce and Praise
}

func (r *SocialAgent) Default() decision.FightAction {
	//TODO implement me
	panic("implement me")
}

/**
 * Create agent with given probability of collaborating
 */
func NewSocialAgent(pCollaborate float64) *SocialAgent {
	return &SocialAgent{pCollaborate: pCollaborate, gossipThreshold: 0.8, propHate: 0.1, propAdmire: 0.1}
}

// TODO: Currently default
func (r *SocialAgent) DonateToHpPool(baseAgent agent.BaseAgent) uint {
	return uint(rand.Intn(int(baseAgent.AgentState().Hp)))
}

// TODO: Currently default
func (r *SocialAgent) UpdateInternalState(agent agent.BaseAgent, _ *commons.ImmutableList[decision.ImmutableFightResult], _ *immutable.Map[decision.Intent, uint]) {
}

// TODO: Currently default
func (r *SocialAgent) FightResolution(baseAgent agent.BaseAgent) message.MapProposal[decision.FightAction] {
	actions := make(map[commons.ID]decision.FightAction)
	view := baseAgent.View()
	agentState := view.AgentState()
	itr := agentState.Iterator()
	for !itr.Done() {
		id, _, ok := itr.Next()
		if !ok {
			break
		}

		switch rand.Intn(3) {
		case 0:
			actions[id] = decision.Attack
		case 1:
			actions[id] = decision.Defend
		default:
			actions[id] = decision.Cower
		}
	}
	newUUID, _ := uuid.NewUUID()
	prop := message.NewProposal(newUUID.String(), commons.MapToImmutable(actions))
	return *prop
}

// TODO: Currently default
func (r *SocialAgent) CreateManifesto(_ agent.BaseAgent) *decision.Manifesto {
	manifesto := decision.NewManifesto(false, false, 10, 50)
	return manifesto
}

// TODO: Currently default
func (r *SocialAgent) HandleConfidencePoll(_ agent.BaseAgent) decision.Intent {
	switch rand.Intn(3) {
	case 0:
		return decision.Abstain
	case 1:
		return decision.Negative
	default:
		return decision.Positive
	}
}

/**
 * Agents dont talk to each other about fight decisions, they decide based on the Q-Table
 */
func (r *SocialAgent) HandleFightInformation(m message.TaggedInformMessage[message.FightInform], baseAgent agent.BaseAgent, log *immutable.Map[commons.ID, decision.FightAction]) {
	switch m.Message().(type) {
	case message.StartFight:
		r.selfID = baseAgent.ID()
		view := baseAgent.View()
		r.battleUtility = utils.AgentBattleUtility(view.AgentState(), r.selfID)
		r.updateSocialCapital(baseAgent.View(), log)
		r.sendGossip(baseAgent)

		// Calculate utility value of each action
		// utilCower := r.utilityValue(decision.Cower, view, agent)
		// utilAttack := r.utilityValue(decision.Attack, view, agent)
		// utilDefend := r.utilityValue(decision.Defend, view, agent)

		// Apply softmax to get probabilities
		//softArray := utils.Softmax([3]float64{utilCower, utilAttack, utilDefend})

		// Make number representation incremental
		//probArray := utils.MakeIncremental(softArray)

		/*it := view.AgentState().Iterator()
		nextId, _, _ := it.Next()
		if agent.Id == nextId {
			fmt.Println(utilCower)
			fmt.Println([3]float64{utilCower, utilAttack, utilDefend})
			fmt.Println(softArray)
			fmt.Println(probArray)
		}*/
	case message.ArrayInfo:
		r.receiveGossip(m.Message().(message.ArrayInfo), m.Sender())
	}

	random := rand.Float64()
	if random < r.pCollaborate {
		r.currentDecision = utils.CollaborativeFightDecision()
	} else {
		r.currentDecision = utils.SelfishFightDecision()
	}

	return
}

// TODO: Currently default
func (r *SocialAgent) HandleFightRequest(m message.TaggedRequestMessage[message.FightRequest], log *immutable.Map[commons.ID, decision.FightAction]) message.FightInform {
	return nil
}

func (r *SocialAgent) CurrentAction() decision.FightAction {
	return r.currentDecision
}

// TODO: Currently default
func (r *SocialAgent) HandleElectionBallot(b agent.BaseAgent, _ *decision.ElectionParams) decision.Ballot {
	// Extract ID of alive agents
	view := b.View()
	agentState := view.AgentState()
	aliveAgentIds := make([]string, agentState.Len())
	i := 0
	itr := agentState.Iterator()
	for !itr.Done() {
		id, a, ok := itr.Next()
		if ok && a.Hp > 0 {
			aliveAgentIds[i] = id
			i++
		}
	}

	// Randomly fill the ballot
	var ballot decision.Ballot
	numAliveAgents := len(aliveAgentIds)
	numCandidate := 2
	for i := 0; i < numCandidate; i++ {
		randomIdx := rand.Intn(numAliveAgents)
		randomCandidate := aliveAgentIds[uint(randomIdx)]
		ballot = append(ballot, randomCandidate)
	}

	return ballot
}

// TODO: Currently default
func (r *SocialAgent) HandleFightProposal(_ message.FightProposalMessage, _ agent.BaseAgent) decision.Intent {
	intent := rand.Intn(2)
	if intent == 0 {
		return decision.Positive
	} else {
		return decision.Negative
	}
}

// TODO: Currently default
func (r *SocialAgent) HandleFightProposalRequest(_ message.FightProposalMessage, _ agent.BaseAgent, _ *immutable.Map[commons.ID, decision.FightAction]) bool {
	switch rand.Intn(2) {
	case 0:
		return true
	default:
		return false
	}
}

// TODO: Currently default
func (r *SocialAgent) HandleUpdateWeapon(_ agent.BaseAgent) decision.ItemIdx {
	// weapons := b.AgentState().Weapons
	// return decision.ItemIdx(rand.Intn(weapons.Len() + 1))

	// 0th weapon has the greatest attack points
	return decision.ItemIdx(0)
}

// TODO: Currently default
func (r *SocialAgent) HandleUpdateShield(_ agent.BaseAgent) decision.ItemIdx {
	// shields := b.AgentState().Shields
	// return decision.ItemIdx(rand.Intn(shields.Len() + 1))
	return decision.ItemIdx(0)
}
