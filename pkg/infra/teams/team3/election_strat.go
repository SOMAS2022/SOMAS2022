package team3

import (
	"infra/game/agent"

	"infra/game/commons"
	"math"

	"infra/game/decision"
	"infra/game/state"
	"infra/logging"
	"math/rand"

	"github.com/benbjohnson/immutable"
)

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
	if len(candidates) > 0 && makeVote < a.personality {
		// Create Ballot
		var ballot decision.Ballot
		// number of manifesto preferences we are allowed
		numCandidate := int(param.NumberOfPreferences())
		for i := 0; i < numCandidate; i++ {
			// look at TSN... if any agents in it, extract their manifestos in reputation order and make decision
			// if no TSN... take manifesto of high reputation agent and evalutate manifesto
			// evaluation is:
			// high reputation + low social capital (that doesn't make any sense.....)
			randomIdx := rand.Intn(len(candidates))
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

func (a *AgentThree) calcW1(state state.HiddenAgentState, id commons.ID) float64 {
	w1 := a.w1Map[id]
	currentHP := state.Hp
	currentStamina := state.Stamina
	prevHP := a.pastHPMap[id]
	prevStamina := a.pastStaminaMap[id]

	// extract and normalise personality (range[0,100]), use to dictate update step size
	personalityMod := float64(a.personality / 100)

	HP := prevHP - int(currentHP)
	stamina := prevStamina - int(currentStamina)

	// alg 6
	if stamina > 0 {
		w1 += 0.5 * personalityMod
	} else if stamina < 0 {
		w1 -= 0.5 * personalityMod
	}
	if HP > 0 {
		w1 += 0.5 * personalityMod
	} else if HP < 0 {
		w1 -= 0.5 * personalityMod
	}

	w1 = clampFloat(w1, 0, 10)

	return w1
}

func (a *AgentThree) calcW2(id commons.ID) float64 {
	w2 := a.w2Map[id]
	agentFought := false
	agentDefended := false
	nFD := 0
	numRounds := a.fightRoundsHistory.Len()

	// extract and normalise personality (range[0,100]), use to dictate update step size
	personalityMod := float64(a.personality / 100)

	// iterate over rounds of last level
	itr := a.fightRoundsHistory.Iterator()
	for !itr.Done() {
		res, _ := itr.Next()

		// search for agent in fight list and assign action
		agentFought = findAgentAction(res.AttackingAgents(), id)
		agentDefended = findAgentAction(res.ShieldingAgents(), id)

		if agentFought || agentDefended {
			nFD++
		}
	}
	// shifted to [-0.5, 0.5]
	w2 += (float64(nFD/numRounds) - 0.5) * personalityMod

	w2 = clampFloat(w2, 0, 10)

	return w2
}

func (a *AgentThree) Reputation(baseAgent agent.BaseAgent) {
	view := baseAgent.View()
	vAS := view.AgentState()
	ids := commons.ImmutableMapKeys(vAS)

	productivity := 5.0
	needs := 5.0
	// fairness := make(map[commons.ID]float64)

	// Number of agents to sample for KA (fixed)
	intendedSample := float64(a.numAgents) * a.samplePercent
	maxLength := float64(vAS.Len())
	sampleLength := int(math.Min(intendedSample, maxLength))
	cnt := 0

	// Use randomness of maps in go to take n random samples
	for _, id := range ids {
		if cnt == sampleLength {
			// Unsorted array
			return
		} else {
			hiddenState, _ := vAS.Get(id)

			// Init values on first access
			if _, ok := a.reputationMap[id]; !ok {
				// init weights to middle value
				a.InitRepWeights(baseAgent, id)
			}

			// Update values according to previous state
			a.w1Map[id] = a.calcW1(hiddenState, id)
			a.w2Map[id] = a.calcW2(id)

			a.reputationMap[id] = a.w1Map[id]*needs + a.w2Map[id]*productivity

			// Store this rounds values for the next one
			a.pastHPMap[id] = int(hiddenState.Hp)
			a.pastStaminaMap[id] = int(hiddenState.Stamina)
		}
		cnt++
	}
}

func (a *AgentThree) InitRepWeights(baseAgent agent.BaseAgent, id commons.ID) {
	view := baseAgent.View()
	vAS := view.AgentState()
	hiddenState, _ := vAS.Get(id)

	a.w1Map[id] = 5.0
	a.w2Map[id] = 5.0
	a.pastHPMap[id] = int(hiddenState.Hp)
	a.pastStaminaMap[id] = int(hiddenState.Stamina)
}

func (a *AgentThree) InitSocialCapital(baseAgent agent.BaseAgent) {
	view := baseAgent.View()
	agentState := view.AgentState()
	itr := agentState.Iterator()
	for !itr.Done() {
		id, _, _ := itr.Next()

		a.socialCap[id] = 25
	}
}

func findAgentAction(agentIDsMap immutable.List[commons.ID], ID commons.ID) bool {
	itr := agentIDsMap.Iterator()
	for !itr.Done() {
		_, actionAgentID := itr.Next()
		if actionAgentID == ID {
			return true
		}
	}
	return false
}

// alg 5
// func (a *AgentThree) CalcReputation(baseAgent agent.BaseAgent) map[commons.ID]float64 {
// 	view := baseAgent.View()
// 	agentState := view.AgentState()

// 	sample := rand.Intn(int(math.Ceil(float64(agentState.Len())*a.samplePercent))-1) + 1
// 	counter := 1

// 	currentLevel := int(view.CurrentLevel())
// 	// init  history
// 	if currentLevel == a.fightRoundsHistory.Len()-1 {
// 		w1 = 0.0
// 		w2 = 0.0

// 		itr := agentState.Iterator()
// 		for !itr.Done() {
// 			if counter%sample == 0 {
// 				id, _, _ := itr.Next()
// 				//idN, _ := strconv.Atoi(id)
// 				//fmt.Println(idN)
// 				pastHP[id] = GetStartingHP()
// 				pastStamina[id] = GetStartingStamina()
// 			} else {
// 				itr.Next()
// 				counter++
// 			}

// 		}
// 	}
// 	productivity := 5.0
// 	needs := 5.0
// 	fairness := make(map[commons.ID]float64)
// 	sortedfairness := make(map[commons.ID]float64)

// 	itr := agentState.Iterator()
// 	counter1 := 1

// 	for i := 1; i <= int(math.Ceil(float64(agentState.Len())*a.samplePercent)); {
// 		if counter1%sample == 0 {
// 			id, hiddenState, _ := itr.Next()

// 			w1 = a.calcW1(hiddenState, w1, pastHP[id], pastStamina[id])
// 			w2 = a.calcW2(baseAgent, w2)

// 			score := w1*needs + w2*productivity

// 			fairness[id] = score

// 			pastHP[id] = int(hiddenState.Hp)
// 			pastStamina[id] = int(hiddenState.Stamina)
// 			i++
// 		} else {
// 			itr.Next()
// 			counter1++
// 		}

// 	}

// 	// get keys for sorting
// 	keys := make([]string, 0, len(fairness))
// 	for key := range fairness {
// 		keys = append(keys, key)
// 	}

// 	// sort map by decreasing score using keys
// 	sort.SliceStable(keys, func(i, j int) bool {
// 		return fairness[keys[i]] > fairness[keys[j]]
// 	})

// 	// reshuffle array
// 	for _, key := range keys {
// 		sortedfairness[key] = fairness[key]
// 	}

// 	return sortedfairness
// }

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

// 	borda := a.Reputation(baseAgent)
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
