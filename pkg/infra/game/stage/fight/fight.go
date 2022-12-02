package fight

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/state"
	"infra/game/tally"
	"math"
	"time"

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
			// todo: prune peer channels somehow...
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

func AgentFightDecisions(state state.State, agents map[commons.ID]agent.Agent, previousDecisions immutable.Map[commons.ID, decision.FightAction], channelsMap map[commons.ID]chan message.TaggedMessage) *tally.Tally[decision.FightAction] {
	proposalVotes := make(chan commons.ProposalID)
	proposalSubmission := make(chan message.MapProposal[decision.FightAction])
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
	mID, _ := uuid.NewUUID()

	for _, messages := range channelsMap {
		messages <- *message.NewTaggedMessage("server", nil, mID)
	}
	time.Sleep(100 * time.Millisecond)
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
	return propTally
}

func HandleFightRound(state state.State, baseHealth uint, fightResult *decision.FightResult) state.State {
	var attackSum uint
	var shieldSum uint

	for agentID, d := range fightResult.Choices {
		agentState := state.AgentState[agentID]

		const scalingFactor = 0.01
		switch d {
		case decision.Attack:
			if agentState.Stamina > agentState.BonusAttack(state) {
				fightResult.AttackingAgents = append(fightResult.AttackingAgents, agentID)
				attackSum += agentState.TotalAttack(state)
				agentState.Stamina = commons.SaturatingSub(agentState.Stamina, agentState.BonusAttack(state))
			} else {
				fightResult.CoweringAgents = append(fightResult.CoweringAgents, agentID)
				fightResult.Choices[agentID] = decision.Cower
				agentState.Hp += uint(math.Ceil(scalingFactor * float64(baseHealth)))
				agentState.Stamina += 1
			}
		case decision.Defend:
			if agentState.Stamina > agentState.BonusDefense(state) {
				fightResult.ShieldingAgents = append(fightResult.ShieldingAgents, agentID)
				shieldSum += agentState.TotalDefense(state)
				agentState.Stamina = commons.SaturatingSub(agentState.Stamina, agentState.BonusDefense(state))
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
