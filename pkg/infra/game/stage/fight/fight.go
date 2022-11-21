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

func DealDamage(damageToDeal uint, agentsFighting []string, agentMap map[commons.ID]agent.Agent, globalState *state.State) {
	splitDamage := damageToDeal / uint(len(agentsFighting))
	for _, id := range agentsFighting {
		agentState := globalState.AgentState[id]
		newHP := commons.SaturatingSub(agentState.Hp, splitDamage)
		if newHP == 0 {
			// kill agent
			// todo: prune peer channels somehow...
			delete(globalState.AgentState, id)
			delete(agentMap, id)
		} else {
			globalState.AgentState[id] = state.AgentState{
				Hp:           newHP,
				Attack:       agentState.Attack,
				Defense:      agentState.Defense,
				BonusAttack:  agentState.BonusAttack,
				BonusDefense: agentState.BonusDefense,
				Stamina:      agentState.Stamina,
			}
		}
	}
}

func AgentFightDecisions(state *state.View, agents map[commons.ID]agent.Agent, previousDecisions immutable.Map[commons.ID, decision.FightAction], channelsMap map[commons.ID]chan message.TaggedMessage) map[commons.ID]decision.FightAction {
	decisionMap := make(map[commons.ID]decision.FightAction)
	channel := make(chan message.ActionMessage, 100)

	var wg sync.WaitGroup

	for _, a := range agents {
		a := a
		wg.Add(1)
		startAgentFightHandlers(*state, &a, previousDecisions, channel, &wg)
	}

	for _, messages := range channelsMap {
		messages <- message.TaggedMessage{
			Sender:  "server",
			Message: *message.NewMessage(message.Something, nil),
		}
	}

	go func(group *sync.WaitGroup) {
		group.Wait()
		for _, messages := range channelsMap {
			close(messages)
		}
		close(channel)
	}(&wg)

	for actionMessage := range channel {
		decisionMap[actionMessage.Sender] = actionMessage.Action
	}

	wg.Wait()

	return decisionMap
}

func HandleFightRound(state *state.State, baseHealth uint, fightResult *decision.FightResult) {
	var attackSum uint
	var shieldSum uint

	for agentID, d := range fightResult.Choices {
		agentState := state.AgentState[agentID]

		const scalingFactor = 0.01
		switch d {
		case decision.Attack:
			if agentState.Stamina > agentState.BonusAttack {
				fightResult.AttackingAgents = append(fightResult.AttackingAgents, agentID)
				attackSum += agentState.TotalAttack()
				agentState.Stamina = commons.SaturatingSub(agentState.Stamina, agentState.BonusAttack)
			} else {
				fightResult.CoweringAgents = append(fightResult.CoweringAgents, agentID)
				fightResult.Choices[agentID] = decision.Cower
				agentState.Hp += uint(math.Ceil(scalingFactor * float64(baseHealth)))
				agentState.Stamina += 1
			}
		case decision.Defend:
			if agentState.Stamina > agentState.BonusDefense {
				fightResult.ShieldingAgents = append(fightResult.ShieldingAgents, agentID)
				shieldSum += agentState.TotalDefense()
				agentState.Stamina = commons.SaturatingSub(agentState.Stamina, agentState.BonusDefense)
			} else {
				fightResult.CoweringAgents = append(fightResult.CoweringAgents, agentID)
				fightResult.Choices[agentID] = decision.Cower
				agentState.Hp += uint(math.Ceil(scalingFactor * float64(baseHealth)))
				agentState.Stamina += 1
			}
		case decision.Cower:
			fightResult.CoweringAgents = append(fightResult.CoweringAgents, agentID)
			agentState.Hp += uint(math.Ceil(scalingFactor * float64(baseHealth)))
			agentState.Stamina += 1
		}
		state.AgentState[agentID] = agentState
	}

	fightResult.AttackSum = attackSum
	fightResult.ShieldSum = shieldSum
}

func startAgentFightHandlers(view state.View, a *agent.Agent, decisionLog immutable.Map[commons.ID, decision.FightAction], channel chan message.ActionMessage, wg *sync.WaitGroup) {
	go a.HandleFight(view, decisionLog, channel, wg)
}
