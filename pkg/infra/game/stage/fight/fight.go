package fight

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/state"
	"math"
	"sync"

	"github.com/google/uuid"

	"github.com/benbjohnson/immutable"
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

func AgentFightDecisions(state state.State, agents map[commons.ID]agent.Agent, previousDecisions immutable.Map[commons.ID, decision.FightAction], channelsMap map[commons.ID]chan message.TaggedMessage) map[commons.ID]decision.FightAction {
	decisionMap := make(map[commons.ID]decision.FightAction)
	channel := make(chan message.ActionMessage, 100)

	var wg sync.WaitGroup

	for _, a := range agents {
		a := a
		wg.Add(1)
		agentState := state.AgentState[a.BaseAgent.Id()]
		startAgentFightHandlers(agentState, &a, previousDecisions, channel, &wg)
	}

	for _, messages := range channelsMap {
		mId, _ := uuid.NewUUID()
		messages <- *message.NewTaggedMessage("server", *message.NewMessage(message.Inform, nil), mId)
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

func HandleFightRound(state state.State, baseHealth uint, fightResult *decision.FightResult) state.State {
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
	return state
}

func startAgentFightHandlers(agentState state.AgentState, a *agent.Agent, decisionLog immutable.Map[commons.ID, decision.FightAction], channel chan message.ActionMessage, wg *sync.WaitGroup) {
	go a.HandleFight(agentState, decisionLog, channel, wg)
}
