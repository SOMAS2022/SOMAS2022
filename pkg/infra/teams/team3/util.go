package team3

import (
	"infra/game/agent"
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


//implemented the 3 functions resource, proposal, chair: considering the utility is between 0 to 15
func (a *AgentThree) UtilityResource(baseAgent agent.BaseAgent){
	//should initiate a utilityResource to 7 (score= map[id]=utilityResource)
    //punishment:
	for id, score := range agentscoreResource{
		if score.id reject trade{
			if score>=12{
				score:= score + 1
			}
			else if score<=12 && score>=7{
				score := score - 2
			}
			else{
				score:= score + 3
			}
		}
	}
	//Update scoreresource map with agent's new scores

    //Rewarding the Agents
        if agent.id agrees trade // the trade has been made
            score+=3

    //Implement Boundaries
    if score > 15{
        score :=15
	}
    if score < 0{
        score :=0
	}
}

func (a *AgentThree) UtilityProposal(baseAgent agent.BaseAgent){
		//should initiate a utilityProposal to 7 (score= map[id]=utilityResource)

    var Lp=0

    if agent has a proposal AND asked us about it AND proposal /= our decision{ //message boolean
        Lp += 1
	}
    if Lp > 4{
        score := score + 1
        Lp := 0
	}
    if agent proposal AND asked us about it AND proposal = our decision{ //message boolean
        score:= score + 1
	}
	If score > 15{
        score :=15
	}
    if score < 0{
        score :=0
	}
}

        
func (a *AgentThree)   UtilityChair (baseAgent agent.BaseAgent){
		//should initiate a utilityChair to 7 (score= map[id]=utilityResource)

	var Lc:=0
    if fight decision =/ proposal.our fight decision{ //message changed to bool
        Lc := Lc + 1
	}
    if Lc > 4:
        score:= score 1 
		Lc = 0
    if ourfight == prop.ourfightdecision:
        score = score + 1

	If score > 15{
		score :=15
	}
	if score < 0{
		score :=0
	}
}

func (a *AgentThree)   TotalUtility (baseAgent agent.BaseAgent){  //average of the 3 score
	    utilityscore:= (UtilityChair() + UtilityProposal() + UtilityResource()) / 3
	}

//trusted social network, after each level
TSN := make(map[int]int) 
for id, util := range utilitymap {
	if util>=10{
		TSN[id]=util
	}
}