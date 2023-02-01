package team3

import (
	"infra/game/agent"
	"infra/game/commons"
	"sort"

	// "infra/game/commons"
	"infra/game/decision"
	"infra/game/state"
	"infra/logging"
	"math/rand"
	"strconv"
)

var (
	w1, w2 float64
)

var pastHP = make(map[commons.ID]int)
var pastStamina = make(map[commons.ID]int)

// Handle No Confidence vote
func (a *AgentThree) HandleConfidencePoll(baseAgent agent.BaseAgent) decision.Intent {
	// decide whether to vote in the no-confidence vote based on personality
	toVote := rand.Intn(100)

	if toVote < a.personality {
		view := baseAgent.View()
		AS := baseAgent.AgentState()
		// vote for leader if they have a high reputation
		baseAgent.Log(logging.Trace, logging.LogField{"hp": AS.Hp, "util": a.utilityScore[view.CurrentLeader()]}, "Util")
		if a.utilityScore[view.CurrentLeader()] > 5 {
			return decision.Positive
		} else {
			// perform no-confidence calculation
			// return answer
			return decision.Positive
		}
	} else {
		return decision.Abstain
	}
}

func (a *AgentThree) HandleElectionBallot(baseAgent agent.BaseAgent, param *decision.ElectionParams) decision.Ballot {

	// Extract ID of alive agents
	// view := baseAgent.View()
	// agentState := view.AgentState()
	// aliveAgentIDs := commons.ImmutableMapKeys(agentState)

	// extract the name of the agents who have submitted manifestos
	candidates := make([]string, param.CandidateList().Len())
	i := 0
	iterator := param.CandidateList().Iterator()
	for !iterator.Done() {
		id, _, _ := iterator.Next()
		candidates[i] = id
		i++
	}

	// should we vote?
	makeVote := rand.Intn(100)
	// if makeVote is lower than personality, then vote.
	if makeVote < a.personality {
		// Randomly fill the ballot
		var ballot decision.Ballot
		numAliveAgents := param.CandidateList().Len()
		numCandidate := int(param.NumberOfPreferences())
		for i := 0; i < numCandidate; i++ {
			randomIdx := rand.Intn(numAliveAgents)
			randomCandidate := candidates[uint(randomIdx)]
			ballot = append(ballot, randomCandidate)
		}

		return ballot
	} else {
		// return an empty ballot (don't vote)
		var ballot decision.Ballot
		return ballot
	}
}

//Agent 3 Voting Strategy

// Effectivness code
// var effective bool

// var initial_monster_attack int = 1

// var //get monster attack

// var prevLevel int = 0
// func  (a *AgentThree) Effectivness_measure(baseAgent agent.BaseAgent){
//     effective := "True"
// 	monster_attack := globalState.MonsterAttack
//     num_agents_alive := len(aliveAgentIDs)
//     percentage_change_w_monster := 1-((monster_attack - initial_monster_attack)/initial_monster_attack)
//     if num_agents_alive>(prevLevel*percentage_change_w_monster){
//         effective := "True"
//     } else{
//         effective := "False"
//     }
//     prevLevel:=num_agents_alive
//     initial_monster_attack:=monster_attack
//     return effective
// }

func (a *AgentThree) calcW1(state state.HiddenAgentState, w1 float64, initHP int, initStamina int) float64 {
	currentHP := state.Hp
	currentStamina := state.Stamina

	HP := initHP - int(currentHP)
	stamina := initStamina - int(currentStamina)
	// alg 6
	if stamina > 0 {
		w1 += 0.2
	} else if stamina < 0 {
		w1 -= 0.2
	}
	if HP > 0 {
		w1 += 0.2
	} else if HP < 0 {
		w1 -= 0.2
	}

	if w1 > 10 {
		w1 = 10
	}

	if w1 < 0 {
		w1 = 0
	}

	return w1
}

func (a *AgentThree) calcW2(baseAgent agent.BaseAgent, w2 float64) float64 {
	var agentFought bool = false
	var agentShielded bool = false
	// iterate until we get most recent history
	i := 0
	itr := a.fightDecisionsHistory.Iterator()
	for !itr.Done() {
		res, _ := itr.Next()
		i = i + 1

		if i == a.fightDecisionsHistory.Len()-1 {
			agents_attack := res.AttackingAgents()
			agents_defended := res.ShieldingAgents()
			itr2 := agents_attack.Iterator()
			itr3 := agents_defended.Iterator()
			// search for our agent in fight list
			for !itr.Done() {
				_, attackingAgentID := itr2.Next()
				if attackingAgentID == baseAgent.ID() {
					agentFought = true
				}
			}
			for !itr.Done() {
				_, defendAgentID := itr3.Next()
				if defendAgentID == baseAgent.ID() {
					agentShielded = true
				}
			}

		}
	}
	if agentFought || agentShielded {
		w2 += 0.2
	} else {
		w2 -= 0.2
	}

	if w2 < 0 {
		w2 = 0
	}

	if w2 > 10 {
		w2 = 10
	}
	return w2
}

// alg 5
func (a *AgentThree) CalcBordaScore(baseAgent agent.BaseAgent) map[commons.ID]float64 {
	view := baseAgent.View()
	agentState := view.AgentState()

	currentLevel := int(view.CurrentLevel())
	// init  history
	if currentLevel == a.fightDecisionsHistory.Len()-1 {
		w1 = 0.0
		w2 = 0.0

		itr := agentState.Iterator()
		for !itr.Done() {
			id, _, _ := itr.Next()
			// idN, _ := strconv.Atoi(id)
			// fmt.Println(idN)
			pastHP[id] = GetStartingHP()
			pastStamina[id] = GetStartingStamina()
		}
	}
	productivity := 5.0
	needs := 5.0
	fairness := make(map[commons.ID]float64)
	sortedfairness := make(map[commons.ID]float64)

	itr := agentState.Iterator()

	for !itr.Done() {
		id, hiddenState, _ := itr.Next()

		w1 = a.calcW1(hiddenState, w1, pastHP[id], pastStamina[id])
		w2 = a.calcW2(baseAgent, w2)

		score := w1*needs + w2*productivity

		fairness[id] = score

		pastHP[id] = int(hiddenState.Hp)
		pastStamina[id] = int(hiddenState.Stamina)
	}

	// get keys for sorting
	keys := make([]string, 0, len(fairness))
	for key := range fairness {
		keys = append(keys, key)
	}

	// sort map by decreasing score using keys
	sort.SliceStable(keys, func(i, j int) bool {
		return fairness[keys[i]] > fairness[keys[j]]
	})

	// reshuffle array
	for _, key := range keys {
		sortedfairness[key] = fairness[key]
	}

	return sortedfairness
}

func (a *AgentThree) SocialCapital(baseAgent agent.BaseAgent) [][]string {
	view := baseAgent.View()
	agentState := view.AgentState()
	itr := agentState.Iterator()
	disobedienceMap := [][]string{}
	for !itr.Done() {
		id, hiddenState, _ := itr.Next()

		if hiddenState.Defector.IsDefector() {
			a.Soc_cap++
		}
		temp := []string{id, strconv.Itoa(a.Soc_cap)}
		disobedienceMap = append(disobedienceMap, temp)
	}
	return disobedienceMap
}

// func (a *AgentThree) Disobedience(baseAgent agent.BaseAgent) {
// 	view := baseAgent.View()
// 	agentState := view.AgentState()
// 	disobedienceMap := make([]int, agentState.Len())
// 	var agentDefected bool
// 	itr := agentState.Iterator()
// 	i := 0

// 	for !itr.Done() {
// 		id, hiddenState, _ := itr.Next()

// 		disobedienceMap[i] += BoolToInt(hiddenState.Defector.IsDefector())
// 		// did we disobey
// 		if id == baseAgent.ID() {
// 			if hiddenState.Defector.IsDefector() {
// 				agentDefected = true
// 			} else {
// 				agentDefected = false
// 			}
// 		}
// 	}

// 	borda := a.CalcBordaScore(baseAgent)
// 	bordaPerCent := BordaPercentage(baseAgent, borda)
// 	for i := range disobedienceMap {
// 		if disobedienceMap[i] >= 5 {
// 			if bordaPerCent < 25 {
// 				// a.utilityScore[baseAgent.ID()] =a.utilityScore[baseAgent.ID()]
// 			} else if bordaPerCent > 25 && bordaPerCent < 50 {
// 				a.utilityScore[baseAgent.ID()] -= 1
// 			} else if bordaPerCent > 50 {
// 				a.utilityScore[baseAgent.ID()] -= 2
// 			} else if agentDefected {
// 				a.utilityScore[baseAgent.ID()] -= 4
// 			}
// 		}
// 	}
// }

// quick function to check if a is in list b
// func in_list(a string, list []string) bool {
// 	for _, b := range list {
// 		if b == a {
// 			return true
// 		}
// 	}
// 	return false
// }

//Now the functions for the voting

// var trusted_network bool := "False"
// var Agent_benefit bool := "False"
// var action_done string // Pleae ping from the thread what our agent did in the last level (fought, defeneded, or cowered)

// //Please fill in the blanks for id_Agent_given_to parameter (get from the thread the ID of the agent that received my resource inestead of me)
// //Also make sure fightDecision() is used correctly
// //Please fill in for CliqueList
// func  (a *AgentThree) update_Agent_benefit(baseAgent agent.BaseAgent, action_done){
//     Common_resource_given:= // please get from the thread if the common pool resource the agent asked for was given to the agent or not
//     Agent_benefit="True"
//     if action_done!=fightDecision(){
//         Agent_benefit="False"
//     }
//     if Common_resource_given="False" && Sorted_agents[id_Agent_given_to]<Sorted_agents[my_agent_ID] && in_list(id_given_to, cliqueList) {
//         Agent_benefit="True"
//     }
//     else if Common_resource_given="False" && Sorted_agents[id_Agent_given_to]<Sorted_agents[my_agent_ID]{
//         Agent_benefit="False"
//     }
// }

// //Please fill in the blanks for chair.ID (id of current chair) and ListClique(list of our trusted Network)
// //Please full for Utility function (Get current Utility score of the chair)
// func  (a *AgentThree) update_Trusted_network(baseAgent agent.BaseAgent, action_done){
//     if in_list(chair.ID, CliqueList) && Utility(chair.ID)>=8 {
//         Trusted_network="True"
//     }
// }

// //The following Code is for confidence vote
// //Fill in Chair.ID
// var vote int
// var counter_not_effective int = 0//to count the levels a chair is effective in.
// var counter_not_benefit_agent int = 0
// //update counter

// func  (a *AgentThree) Confidence_vote(baseAgent agent.BaseAgent){
//     vote:=1
//     Effectivness_measure()
//     update_Agent_benefit()
//     if !effective {
//         counter_not_effective+=1
//     }else {
//         counter_not_effective:=0
//     }
//     if !Agent_benefit{
//         counter_not_benefit_agent +=1
//     }else{
//         counter_not_benefit_agent := 0
//     }
//     if Trusted_network && counter_not_effective>2 {
//         return vote:=0
//     }
//     if Trusted_network && counter_not_benefit_agent>2{
//         return vote:=0
//     }
//     if  counter_not_effective>1 {
//         return vote:= 0
//     }
//     if counter_not_benefit_agent>1 {
//         return vote := 0
//     }

// }

//Now to vote for next chair if chair is deposed
