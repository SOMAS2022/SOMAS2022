package team3

import (
	"math/rand"

	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/logging"

	"github.com/benbjohnson/immutable"
)

const PERCENTAGE = 500

type Utility struct {
	ID    commons.ID
	score int
}

type UtilityMap []Utility

func (u UtilityMap) Len() int           { return len(u) }
func (u UtilityMap) Swap(i, j int)      { u[i], u[j] = u[j], u[i] }
func (u UtilityMap) Less(i, j int) bool { return u[i].score < u[j].score }

type AgentThree struct {
	HP           int
	ST           int
	AT           int
	bravery      int
	uR           map[commons.ID]int
	uP           map[commons.ID]int
	uC           map[commons.ID]int
	utilityScore map[commons.ID]int
	TSN          []commons.ID
}

// Update internal parameters at the end of each lvl!?
func (a *AgentThree) UpdateInternalState(baseAgent agent.BaseAgent, _ *commons.ImmutableList[decision.ImmutableFightResult], _ *immutable.Map[decision.Intent, uint], log chan<- logging.AgentLog) {
	a.UpdateUtility(baseAgent)
	a.HP = int(baseAgent.AgentState().Hp)
	a.ST = int(baseAgent.AgentState().Stamina)
	a.AT = int(baseAgent.AgentState().Attack)
}

// Create proposal for fight decisions
// func (a *AgentThree) FightResolution(baseAgent agent.BaseAgent) message.MapProposal[decision.FightAction] {
// 	actions := make(map[commons.ID]decision.FightAction)
// 	view := baseAgent.View()
// 	agentState := view.AgentState()
// 	itr := agentState.Iterator()
// 	for !itr.Done() {
// 		id, _, ok := itr.Next()
// 		if !ok {
// 			break
// 		}

// 		// Check for our agent and assign what we want to do
// 		if id == baseAgent.ID() {
// 			actions[id] = a.CurrentAction()
// 			baseAgent.Log(logging.Trace, logging.LogField{"bravery": a.bravery, "hp": a.HP, "choice": a.CurrentAction(), "util": a.utilityScore[view.CurrentLeader()]}, "Intent")
// 		} else {
// 			// Send some messages to other agents
// 			// send := rand.Intn(5)
// 			// if send == 0 {
// 			// 	m := message.FightInform()
// 			// 	_ = baseAgent.SendBlockingMessage(id, m)
// 			// }
// 			rNum := rand.Intn(3)
// 			switch rNum {
// 			case 0:
// 				actions[id] = decision.Attack
// 			case 1:
// 				actions[id] = decision.Defend
// 			default:
// 				actions[id] = decision.Cower
// 			}
// 		}
// 	}

// 	prop := message.NewProposal(uuid.NewString(), commons.MapToImmutable(actions))
// 	return *prop
// }

func NewAgentThree() agent.Strategy {
	return &AgentThree{
		bravery:      rand.Intn(10),
		utilityScore: make(map[string]int),
	}
}
