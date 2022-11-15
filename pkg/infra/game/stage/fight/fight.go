package fight

import (
	"github.com/benbjohnson/immutable"
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/state"
	"math"
	"sync"
)

func DealDamage(attack uint, agentMap map[commons.ID]agent.Agent, globalState *state.State) {
	splitDamage := attack / uint(len(agentMap))
	for id, agentState := range globalState.AgentState {
		newHp := commons.SaturatingSub(agentState.Hp, splitDamage)
		if newHp == 0 {
			// kill agent
			// todo: prune peer channels somehow...
			delete(globalState.AgentState, id)
			delete(agentMap, id)
		} else {
			globalState.AgentState[id] = state.AgentState{
				Hp:           newHp,
				Attack:       agentState.Attack,
				Defense:      agentState.Defense,
				BonusAttack:  agentState.BonusAttack,
				BonusDefense: agentState.BonusDefense,
			}
		}
	}
}

func HandleFightRound(state state.State, agents map[commons.ID]agent.Agent, baseHealth uint, previousDecisions immutable.Map[commons.ID, decision.FightAction], channelsMap map[commons.ID]chan message.TaggedMessage) (uint, uint, uint, map[commons.ID]decision.FightAction) {
	decisionMap := make(map[commons.ID]decision.FightAction)
	channels := make(map[commons.ID]chan decision.FightAction)

	view := state.ToView()
	//todo: unify to single channel to avoid mutex
	for i, a := range agents {
		a := a
		channels[i] = startAgentFightHandlers(*view, &a, previousDecisions)
	}

	for _, messages := range channelsMap {
		messages <- message.TaggedMessage{
			Sender:  "server",
			Message: *message.NewMessage(message.Something, nil),
		}
	}
	var mutex = &sync.RWMutex{}
	var wg sync.WaitGroup

	for i, dChan := range channels {
		wg.Add(1)
		go func(i commons.ID, c chan decision.FightAction) {
			defer wg.Done()
			mutex.Lock()
			decisionMap[i] = <-c
			mutex.Unlock()
			close(c)
		}(i, dChan)
	}

	var coweringAgents uint
	var attackSum uint
	var shieldSum uint
	wg.Wait()
	for agentID, d := range decisionMap {
		agentState := state.AgentState[agentID]

		switch d {
		case decision.Attack:
			if agentState.Stamina > agentState.BonusAttack {
				attackSum += agentState.TotalAttack()
				agentState.Stamina = commons.SaturatingSub(agentState.Stamina, agentState.BonusAttack)
			} else {
				coweringAgents++
				decisionMap[agentID] = decision.Cower
				agentState.Hp += uint(math.Ceil(0.05 * float64(baseHealth)))
				agentState.Stamina += 1
			}
		case decision.Defend:
			if agentState.Stamina > agentState.BonusDefense {
				shieldSum += agentState.TotalDefense()
				agentState.Stamina = commons.SaturatingSub(agentState.Stamina, agentState.BonusDefense)
			} else {
				coweringAgents++
				decisionMap[agentID] = decision.Cower
				agentState.Hp += uint(math.Ceil(0.05 * float64(baseHealth)))
				agentState.Stamina += 1
			}
		case decision.Cower:
			coweringAgents++
			agentState.Hp += uint(math.Ceil(0.05 * float64(baseHealth)))
			agentState.Stamina += 1
		}
		state.AgentState[agentID] = agentState
	}

	return coweringAgents, attackSum, shieldSum, decisionMap
}

func startAgentFightHandlers(view state.View, a *agent.Agent, decisionLog immutable.Map[commons.ID, decision.FightAction]) chan decision.FightAction {
	decisionChan := make(chan decision.FightAction)
	go a.HandleFight(view, decisionLog, decisionChan)
	return decisionChan
}
