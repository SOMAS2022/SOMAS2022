package fight

import (
	"math"
	"time"

	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/state"
	"infra/game/tally"

	"github.com/benbjohnson/immutable"
	"github.com/google/uuid"
)

func DealDamage(damageToDeal uint, agentsFighting []string, agentMap map[commons.ID]agent.Agent, globalState *state.State) {
	splitDamage := damageToDeal / uint(len(agentsFighting))
	for _, id := range agentsFighting {
		agentState := globalState.AgentState[id]
		newHP := commons.SaturatingSub(agentState.Hp, splitDamage)
		if newHP == 0 {
			// kill agent
			removeItems(globalState, globalState.AgentState[id])

			delete(globalState.AgentState, id)
			delete(agentMap, id)
		} else {
			globalState.AgentState[id] = state.AgentState{
				Hp:          newHP,
				Attack:      agentState.Attack,
				Defense:     agentState.Defense,
				Stamina:     agentState.Stamina,
				Weapons:     agentState.Weapons,
				Shields:     agentState.Shields,
				WeaponInUse: agentState.WeaponInUse,
				ShieldInUse: agentState.ShieldInUse,
			}
		}
	}
}

func removeItems(globalState *state.State, agentState state.AgentState) {
	removeItemsFromMap(globalState.InventoryMap.Weapons, agentState.Weapons)
	removeItemsFromMap(globalState.InventoryMap.Shields, agentState.Shields)
}

func removeItemsFromMap(m map[commons.ID]uint, l immutable.List[state.Item]) {
	iterator := l.Iterator()
	for !iterator.Done() {
		_, v := iterator.Next()
		delete(m, v.Id())
	}
}

func AgentFightDecisions(state state.State, agents map[commons.ID]agent.Agent, previousDecisions immutable.Map[commons.ID, decision.FightAction], channelsMap map[commons.ID]chan message.TaggedMessage) *tally.Tally[decision.FightAction] {
	proposalVotes := make(chan commons.ProposalID)
	proposalSubmission := make(chan message.Proposal[decision.FightAction])
	tallyClosure := make(chan struct{})

	propTally := tally.NewTally(proposalVotes, proposalSubmission, tallyClosure)
	go propTally.HandleMessages()
	closures := make(map[commons.ID]chan<- struct{})
	for id, a := range agents {
		a := a
		closure := make(chan struct{})
		closures[id] = closure
		agentState := state.AgentState[a.BaseAgent.ID()]
		if a.BaseAgent.ID() == state.CurrentLeader {
			go (&a).HandleFight(agentState, previousDecisions, proposalVotes, proposalSubmission, closure)
		} else {
			go (&a).HandleFight(agentState, previousDecisions, proposalVotes, nil, closure)
		}
	}
	mID := uuid.Nil

	for _, messages := range channelsMap {
		messages <- *message.NewTaggedMessage("server", &message.StartFight{}, mID)
	}
	time.Sleep(25 * time.Millisecond)
	for id, c := range channelsMap {
		closures[id] <- struct{}{}
		go func(recv <-chan message.TaggedMessage) {
			for m := range recv {
				switch m.Message().(type) {
				case message.Request:
					// todo: respond with nil thing here as we're closing! Or do we need to?
					// maybe because we're closing there's no point...
				default:
				}
			}
		}(c)
	}

	for _, c := range channelsMap {
		close(c)
	}

	tallyClosure <- struct{}{}
	close(tallyClosure)
	return propTally
}

func HandleFightRound(state state.State, baseStamina uint, baseHealth uint, fightResult *decision.FightResult) *state.State {
	var attackSum uint
	var shieldSum uint

	for agentID, d := range fightResult.Choices {
		agentState := state.AgentState[agentID]

		const scalingFactor = 0.02
		switch d {
		case decision.Attack:
			if agentState.Stamina > agentState.BonusAttack() {
				fightResult.AttackingAgents = append(fightResult.AttackingAgents, agentID)
				attackSum += agentState.TotalAttack()
				agentState.Stamina = commons.SaturatingSub(agentState.Stamina, agentState.TotalAttack())
			} else {
				fightResult.CoweringAgents = append(fightResult.CoweringAgents, agentID)
				fightResult.Choices[agentID] = decision.Cower
				agentState.Hp += uint(math.Ceil(scalingFactor * float64(baseHealth)))
				agentState.Stamina += uint(math.Ceil(scalingFactor * float64(baseStamina)))
			}
		case decision.Defend:
			if agentState.Stamina > agentState.BonusDefense() {
				fightResult.ShieldingAgents = append(fightResult.ShieldingAgents, agentID)
				shieldSum += agentState.TotalDefense()
				agentState.Stamina = commons.SaturatingSub(agentState.Stamina, agentState.TotalDefense())
			} else {
				fightResult.CoweringAgents = append(fightResult.CoweringAgents, agentID)
				fightResult.Choices[agentID] = decision.Cower
				agentState.Hp += uint(math.Ceil(scalingFactor * float64(baseHealth)))
				agentState.Stamina += uint(math.Ceil(scalingFactor * float64(baseStamina)))
			}
		case decision.Cower:
			fightResult.CoweringAgents = append(fightResult.CoweringAgents, agentID)
			agentState.Hp += uint(math.Ceil(scalingFactor * float64(baseHealth)))
			agentState.Stamina += uint(math.Ceil(scalingFactor * float64(baseStamina)))
		}
		state.AgentState[agentID] = agentState
	}

	fightResult.AttackSum = attackSum
	fightResult.ShieldSum = shieldSum
	return &state
}
