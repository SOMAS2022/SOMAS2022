package team3

import (
	"infra/game/agent"
	"infra/game/decision"
	"math/rand"
)

// Update Utility
func (a *AgentThree) UpdateUtility(baseAgent agent.BaseAgent) {
	view := baseAgent.View()
	agentState := view.AgentState()
	itr := agentState.Iterator()
	for !itr.Done() {
		id, _, ok := itr.Next()
		if !ok {
			break
		}

		// We are already cool, don't need the utility score for ourselves
		if id != baseAgent.ID() {
			a.utilityScore[id] = rand.Intn(10)
		}
	}
	// Sort utility map
	// sort.Sort(a.utilityScore)
}

func (a *AgentThree) limitScore(score int) int {
	// Implement Boundaries
    if score > 15 {
        score = 15
	}
    if score < 0 {
        score = 0
	}
	return score
}

//implemented the 3 functions resource, proposal, chair: considering the utility is between 0 to 15
func (a *AgentThree) UtilityResource(baseAgent agent.BaseAgent) {
	//should initiate a utilityResource to 7 (score= map[id]=utilityResource)
    //punishment:
	for id, score := range a.uR {
		if trade[score.id] == Reject {
			if score >= 12 {
				score = score + 1
			} else if score <= 12 && score >= 7 {
				score = score - 2
			} else {
				score = score + 3
			}
		}
	}
	//Update scoreresource map with agent's new scores

    // Rewarding the Agents
	if trade[agent.id] == Agree {// the trade has been made
		score+=3
	}
    score = a.limitScore(score)
}

func (a *AgentThree) UtilityProposal(baseAgent agent.BaseAgent) {
	//should initiate a utilityProposal to 7 (score= map[id]=utilityResource)
    Lp := 0

    if agent has a proposal AND asked us about it AND proposal != our decision { //message boolean
        Lp += 1
	}
    if Lp > 4 {
        score := score + 1
        Lp := 0
	}
    if agent proposal AND asked us about it AND proposal == our decision { //message boolean
        score:= score + 1
	}
	score := a.limitScore(score)
}

        
func (a *AgentThree) UtilityChair(baseAgent agent.BaseAgent) {
	//should initiate a utilityChair to 7 (score= map[id]=utilityResource)
	Lc := 0

    if fight_decision != proposal.our_fight_decision { //message changed to bool
        Lc = Lc + 1
	}
    if Lc > 4 {
        score := score + 1 
		Lc = 0
	}
    if our_fight == prop.our_fight_decision {
        score := score + 1
	}
	score := a.limitScore(score)
}

func (a *AgentThree) TotalUtility(baseAgent agent.BaseAgent) map[string]int {  
	//average of the 3 score
	return (UtilityChair() + UtilityProposal() + UtilityResource()) / 3
}

func (a *AgentThree) UpdateTSN(basebaseAgent agent.BaseAgent) {
	//trusted social network, after each level
	for id, util := range a.utilityScore {
		if util >= 10 {
			a.TSN[id] = util
		}
	}
}
