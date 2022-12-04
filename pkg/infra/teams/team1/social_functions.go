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
	"infra/game/state"
	"infra/teams/team1/utils"
	"sort"

	"github.com/benbjohnson/immutable"
)

func (r *SocialAgent) initSocialCapital(selfID string, sci map[string]SocialCapInfo) {
	r.selfID = selfID
	// Create empty map
	r.socialCapital = sci

	// Delete the agents own id from the socialCapital array
	delete(r.socialCapital, r.selfID)

	// Set the lastLevelUpdated variable
	r.lastLevelUpdated = 0
}

// Called any time a message is received, initialises or updates the socialCapital map
func (r *SocialAgent) updateSocialCapital(view state.View, log *immutable.Map[commons.ID, decision.FightAction]) {
	// Ensure that socialCapital map is initialised
	agentState := view.AgentState()
	updatedSocCapInfo := SocialCapInfo{}
	if r.lastLevelUpdated < view.CurrentLevel() {
		for key := range r.socialCapital {
			// Remove any agents that have died from socialCapital map (Might be unnecessary as it adds a lot of computation)
			_, exists := agentState.Get(key)
			if !exists {
				delete(r.socialCapital, key)
			}

			updatedSocCapInfo.ID = key
			// Decay socialCapital values
			updatedSocCapInfo.arr = utils.DecayArray(r.socialCapital[key].arr)

			// TODO: Update of socialCaptial should be dependent on the agents own action (especially for favours)
			// Update socialCapital based on log
			action, exists := log.Get(key)
			if exists {
				updatedSocCapInfo.arr = utils.BoundArray(utils.AddArrays(updatedSocCapInfo.arr, utils.BoundArray(utils.ActionSentiment(action))))
			}

			r.socialCapital[key] = updatedSocCapInfo
		}

		// Set lastLevelUpdated to current level
		r.lastLevelUpdated = view.CurrentLevel()
	}
}

// Calculate utility value of different decisions
func (r *SocialAgent) utilityValue(action decision.FightAction, _ *state.View, agent agent.BaseAgent) float64 {
	// Utility of each action is dependent on relationship with others. If agent hates all other agents, then
	// will only act in its own interest.

	agentStats := agent.AgentState()

	switch action {
	case decision.Cower:
		// Goes down with health, and down with stamina and down with high social capital of others
		return 0.005 * float64(1000-int(agentStats.Hp))
	case decision.Attack:
		// Goes up with health, and up with stamina and up with high social capital of others
		return 0.005 * float64(int(agentStats.Hp))
	case decision.Defend:
		return 0.005 * float64(int(agentStats.Hp))
	default:
		return 1
	}
}

/**
 * Tell other trusted agents above a threshold T about the A agents they admire the most,
 * and the H agents they hate the most.
 * Currently messages may be sent to recipients praising or denouncing the recipient
 */
func (r *SocialAgent) sendGossip(agent agent.BaseAgent) {
	selfID := agent.Name()

	sortedSCTrustworthiness := make([]SocialCapInfo, 0, len(r.socialCapital))
	for _, sci := range r.socialCapital {
		if sci.ID == selfID { // Exclude self
			continue
		}
		sortedSCTrustworthiness = append(sortedSCTrustworthiness, sci)
	}

	sort.Slice(sortedSCTrustworthiness, func(i int, j int) bool {
		return sortedSCTrustworthiness[i].arr[2] > sortedSCTrustworthiness[j].arr[2]
	})

	numAdmire := int(r.propAdmire * float64(len(sortedSCTrustworthiness)))
	numHate := int(r.propHate * float64(len(sortedSCTrustworthiness)))

	admiredAgents := make([]string, 0, numAdmire)
	hatedAgents := make([]string, 0, numHate)
	for i, scit := range sortedSCTrustworthiness {
		if i >= numAdmire {
			break
		}

		admiredAgents = append(admiredAgents, scit.ID)
	}

	for i := len(sortedSCTrustworthiness) - 1; i >= 0; i-- {
		if len(sortedSCTrustworthiness)-1-i >= numHate {
			break
		}

		hatedAgents = append(hatedAgents, sortedSCTrustworthiness[i].ID)
	}

	for _, sci := range sortedSCTrustworthiness {
		if sci.arr[1] < r.gossipThreshold {
			continue
		}
		utils.Gossip(agent, sci.ID, utils.Praise, admiredAgents)
		utils.Gossip(agent, sci.ID, utils.Denounce, hatedAgents)
	}
}

/**
 * On receiving gossip, scale the network value up/down by a constant and
 * the senders overall perception to this agent. This means that someone
 * with 0.2 network would become 0.22, and someone with 0.7 network would
 * become 0.77, with a 10% increase
 *
 * TODO: More advanced but consider reversing perception if the agent
 * hates the sender
 */
func (r *SocialAgent) receiveGossip(m message.ArrayInfo, sender string) {
	senderPerception := utils.OverallPerception(r.socialCapital[sender].arr)
	var sign float64
	switch m.Num {
	case utils.Praise:
		sign = 1.0
	case utils.Denounce:
		sign = -1.0
	}

	for _, about := range m.StringArr {
		sci := r.socialCapital[about]
		sci.arr[1] += sign * senderPerception * 0.1 * sci.arr[1]
		sci.arr = utils.BoundArray(sci.arr)
		r.socialCapital[about] = sci
	}
}
