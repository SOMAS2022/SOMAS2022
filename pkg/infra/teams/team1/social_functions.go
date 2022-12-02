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

// Called any time a message is received, initialises or updates the socialCapital map
func (r *SocialAgent) updateSocialCapital(_ message.TaggedMessage, view *state.View, agent agent.BaseAgent, log *immutable.Map[commons.ID, decision.FightAction]) {
	// Ensure that socialCapital map is initialised
	agentState := view.AgentState()
	agentStateLength := agentState.Len()
	updatedSocCapInfo := SocialCapInfo{}
	if len(r.socialCapital) == 0 && agentStateLength > 1 {
		// Create empty map
		r.socialCapital = map[string]SocialCapInfo{}

		// Populate map with every currently living agent, and calculate socialCapital based on log
		itr := agentState.Iterator()
		for !itr.Done() {
			key, _, _ := itr.Next()

			action, exists := log.Get(key)
			updatedSocCapInfo.ID = key
			if exists { // If agent exists in log, calculate socialCapital
				updatedSocCapInfo.arr = utils.BoundArray(utils.ActionSentiment(action))
			} else { // Else initialize socialCapital to 0
				updatedSocCapInfo.arr = [4]float64{0.0, 0.0, 0.0, 0.0}
			}
			r.socialCapital[key] = updatedSocCapInfo
		}

		// Delete the agents own id from the socialCapital array
		delete(r.socialCapital, agent.Name())

		// Set the lastLevelUpdated variable
		r.lastLevelUpdated = view.CurrentLevel()
	} else if r.lastLevelUpdated < view.CurrentLevel() { // socialCapital variable already exists
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

	agentStats := agent.ViewState()

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
 * Currently messages may be sent to recepients praising or denouncing the recepient
 */
func (r *SocialAgent) Gossip(agent agent.BaseAgent) {
	selfID := agent.Name()

	sortedSCTrustworthiness := make([]SocialCapInfo, 0, len(r.socialCapital))
	for _, sci := range r.socialCapital {
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

		if scit.ID == selfID { // Dont include self
			numAdmire++
			continue
		}

		admiredAgents = append(admiredAgents, scit.ID)
	}

	for i := len(sortedSCTrustworthiness) - 1; i >= 0; i-- {
		if len(sortedSCTrustworthiness)-1-i >= numHate {
			break
		}

		if sortedSCTrustworthiness[i].ID == selfID { // Dont include self
			numAdmire++
			continue
		}

		hatedAgents = append(hatedAgents, sortedSCTrustworthiness[i].ID)
	}

	for _, sci := range sortedSCTrustworthiness {
		if sci.arr[1] < r.gossipThreshold || sci.ID == selfID {
			continue
		}
		utils.Gossip(agent, sci.ID, utils.Praise, admiredAgents)
		utils.Gossip(agent, sci.ID, utils.Denounce, hatedAgents)

	}

}
