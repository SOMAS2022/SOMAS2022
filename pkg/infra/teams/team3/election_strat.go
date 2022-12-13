package team3

import (
	"infra/game/agent"
	"infra/game/decision"
	"math/rand"
	"sort"
)

// Handle No Confidence vote
func (a *AgentThree) HandleConfidencePoll(baseAgent agent.BaseAgent) decision.Intent {
	view := baseAgent.View()
	// Vote for leader to stay if he's our friend :)
	if a.utilityScore[view.CurrentLeader()] > 5 {
		return decision.Positive
	} else {
		switch rand.Intn(2) {
		case 0:
			return decision.Abstain
		default:
			return decision.Negative
		}
	}
}

func (a *AgentThree) HandleElectionBallot(baseAgent agent.BaseAgent, _ *decision.ElectionParams) decision.Ballot {
	// Extract ID of alive agents
	view := baseAgent.View()
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
	numCandidate := 2
	for i := 0; i < numCandidate; i++ {
		randomIdx := rand.Intn(numAliveAgents)
		randomCandidate := aliveAgentIDs[uint(randomIdx)]
		ballot = append(ballot, randomCandidate)
	}

	return ballot
}
//Agent 3 Voting Strategy

//Effectivness code
var effective bool

var initial_monster_attack int = 1
var //get monster attack
var prevLevel int = 0 
func  (a *AgentThree) Effectivness_measure(baseAgent agent.BaseAgent){
    effective := "True"
    monster_attack := //get monster attack
    num_agents_alive := len(aliveAgentIDs) 
    percentage_change_w_monster := 1-((monster_attack - initial_monster_attack)/initial_monster_attack)
    if num_agents_alive>(prevLevel*percentage_change_w_monster){
        effective := "True"
    } else{
        effective := "False"
    }
    prevLevel:=num_agents_alive
    initial_monster_attack:=monster_attack
    return effective 
}


//borda score 
var w1 map[string]float64
var w2 map[string]float64
var healthnow map[string]float64 //please fill this map with agents current health
var staminaNow map [string]float64 //please fill this map with agents current stamina

var scores map[string]float64

func (a *AgentThree) updateW1(baseAgent agent.BaseAgent){
    for i:=0; i == len(AgentMap);i++ {
        hp:=healthInitial[commons.ID]-healthNow[commons.ID]
        if hp>0{
            w1[commons.ID]+=0.2
        } else if hp==0 {
            w1[commons.ID]=w1[commons.ID]
        } else{
            w1[commons.ID]-=0.2
        }
        if w1[commons.ID]>10{
            w1[commons.ID]=10
        }
        if w1[commons.ID]<0{
            w1[commons.ID]=0
        }
}

func (a *AgentThree) updateW2(baseAgent agent.BaseAgent){
    for i:=0; i == len(AgentMap);i++ {
        st:=staminaInitial[commons.ID]-staminaNow[commons.ID]
        if st>0{
            w2[commons.ID]+=0.2
        } else{
            w1[commons.ID]-=0.2
        }
        if w2[commons.ID]>10{
            w2[commons.ID]=10
        }
        if w2[commons.ID]<0{
            w2[commons.ID]=0
        }
}

// Please fill in the blanks for the AgentMap parameter(used to iteratre over the agents)
func  (a *AgentThree) Borda_score(baseAgent agent.BaseAgent, w1, w2){
    needs:= 5
    productivity:= 5
    Sorted_Agents:= make([]string,0,len(AgentMap))
    for i:=0; i == len(AgentMap);i++ {
        Agent_Score:= w1[commons.ID]*needs+ w2[commons.ID]*productivity
        scores[commons.ID]=Agent_Score
    }
    for key := range scores{
        Sorted_Agents = append(Sorted_Agents, key)
    }
    sort.SliceStable(Sorted_Agents, func(i, j int)) bool{
        return scores[Sorting[i]] < scores[keys[j]]
    }
}

//quick function to check if a is in list b
func in_list(a string, list []string) bool{
    for _, b:= range list {
        if b==a{
            return true
        }
    }
    return false
}

//Now the functions for the voting

var trusted_network bool := "False"
var Agent_benefit bool := "False"
var action_done string // Pleae ping from the thread what our agent did in the last level (fought, defeneded, or cowered)

//Please fill in the blanks for id_Agent_given_to paramter (get from the thread the ID of the agent that received my resource inestead of me)
//Also make sure fightDecision() is used correctly
//Please fill in for CliqueList
func  (a *AgentThree) update_Agent_benefit(baseAgent agent.BaseAgent, action_done){
    Common_resource_given:= // please get from the thread if the common pool resource the agent asked for was given to the agent or not
    Agent_benefit="True"
    if action_done!=fightDecision(){
        Agent_benefit="False"
    }
    if Common_resource_given="False" && Sorted_agents[id_Agent_given_to]<Sorted_agents[my_agent_ID] && in_list(id_given_to, cliqueList) {
        Agent_benefit="True"
    }
    else if Common_resource_given="False" && Sorted_agents[id_Agent_given_to]<Sorted_agents[my_agent_ID]{
        Agent_benefit="False"
    } 
}

//Please fill in the blanks for chair.ID (id of current chair) and ListClique(list of our trusted Network)
//Please full for Utility funtion (Get current Utility score of the chair)
func  (a *AgentThree) update_Trusted_network(baseAgent agent.BaseAgent, action_done){
    if in_list(chair.ID, CliqueList) && Utility(chair.ID)>=8 {
        Trusted_network="True"
    }
}

//The following Code is for confidence vote
//Fill in Chair.ID
var vote int
var counter_not_effective int = 0//to count the levels a chair is effective in.
var counter_not_benefit_agent int = 0
//update counter

func  (a *AgentThree) Confidence_vote(baseAgent agent.BaseAgent){
    vote:=1
    Effectivness_measure()
    update_Agent_benefit()
    if !effective {
        counter_not_effective+=1
    }else {
        counter_not_effective:=0
    }
    if !Agent_benefit{
        counter_not_benefit_agent +=1
    }else{
        counter_not_benefit_agent := 0
    }
    if Trusted_network && counter_not_effective>2 {
        return vote:=0
    } 
    if Trusted_network && counter_not_benefit_agent>2{
        return vote:=0
    }
    if  counter_not_effective>1 {
        return vote:= 0
    }
    if counter_not_benefit_agent>1 {
        return vote := 0
    }

}

//Now to vote for next chair if chair is deposed




