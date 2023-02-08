package team3

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"

	"github.com/benbjohnson/immutable"
)

// Resource (trading) utility
func (a *AgentThree) UpdateUtilityResource(baseAgent agent.BaseAgent, msg message.TradeResponse, id commons.ID) {
	// Check if msg type is correct

	// Punishment depending on previous score:
	switch msg.(type) {
	case message.TradeReject:
		if a.uR[id] >= 12 {
			a.uR[id] -= 1
		} else if a.uR[id] < 12 && a.uR[id] >= 7 {
			a.uR[id] -= 2
		} else {
			a.uR[id] -= 3
		}

		// Rewarding the Agents (successful trade):
	case message.TradeAccept:
		a.uR[id] += 3
	}
	a.uR[id] = limitScore(a.uR[id])
}

// Proposal utility
func (a *AgentThree) UpdateUtilityProposal(baseAgent agent.BaseAgent, propAction decision.FightAction, id commons.ID) {
	// calculate propAction from Agent[id] before calling this function! (based on his thresholds?)

	// Check if agent is in accordance with us and asked us in the last round
	// TODO: update contacts when receiving message
	if a.contactsLastRound[id] {
		if a.FightActionNoProposal(baseAgent) != propAction {
			a.proposalTolerance[id] += 1
			// Tolerance
			if a.proposalTolerance[id] >= 3 {
				a.uP[id] -= 1
				a.proposalTolerance[id] = 0
			}
		} else {
			a.uP[id] += 1
		}
	}
	a.uP[id] = limitScore(a.uP[id])
}

// Chair utility
func (a *AgentThree) UpdateUtilityChair(baseAgent agent.BaseAgent, prop immutable.Map[commons.ID, decision.FightAction]) {
	//should initiate a utilityChair to 7 (score= map[id]=utilityResource)

	view := baseAgent.View()
	chairID := view.CurrentLeader()

	propAction, _ := prop.Get(baseAgent.ID())

	// Check if chair is in accordance with us
	if a.FightActionNoProposal(baseAgent) != propAction {
		a.chairTolerance += 1
		// Tolerance
		if a.chairTolerance >= 4 {
			a.uC[chairID] -= 1
			a.chairTolerance = 0
		}
	} else {
		a.uC[chairID] += 1
	}
	a.uC[chairID] = limitScore(a.uC[chairID])
}

func (a *AgentThree) UpdateTotalUtility(baseAgent agent.BaseAgent) {
	// Weighted Average
	for i := range a.utilityScore {
		a.utilityScore[i] = int(0.45*float64(a.uC[i]) + 0.45*float64(a.uR[i]) + 0.10*float64(a.uP[i]))
	}
}

// Add agent to TSN (if he doesn't exist already)
func (a *AgentThree) AddToTSN(id commons.ID) {
	exists := false
	for _, TSNid := range a.TSN {
		// check if agent exists
		if TSNid == id {
			exists = true
			break
		}
	}
	if !exists {
		a.TSN = append(a.TSN, id)
	}
}

// Remove agent from TSN
func (a *AgentThree) RemoveFromTSN(id commons.ID) []commons.ID {
	for i, TSNid := range a.TSN {
		// Put last element on deleted space and return N-1
		if TSNid == id {
			a.TSN[i] = a.TSN[len(a.TSN)-1]
			a.TSN = a.TSN[:len(a.TSN)-1]
			break
		}
	}
	return a.TSN
}

// Update TSN
func (a *AgentThree) UpdateTSN(baseAgent agent.BaseAgent) {
	for id, rep := range a.reputationMap {
		if rep >= 75 {
			a.AddToTSN(id)
		} else {
			a.RemoveFromTSN(id)
		}
	}
}

func clampFloat(val, low, high float64) float64 {
	if val > high {
		return high
	}
	if val < low {
		return low
	}
	return val
}

func clampInt(val, low, high int) int {
	if val > high {
		return high
	}
	if val < low {
		return low
	}
	return val
}
