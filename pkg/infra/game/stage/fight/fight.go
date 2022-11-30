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

func AgentFightDecisions(state state.State, agents map[commons.ID]agent.Agent, previousDecisions immutable.Map[commons.ID, decision.FightAction], channelsMap map[commons.ID]chan message.TaggedMessage) *tally.Tally[decision.FightAction] {
	proposalVotes := make(chan commons.ProposalID)
	proposalSubmission := make(chan tally.Proposal[decision.FightAction])
	closure := make(chan struct{})

	propTally := tally.NewTally(proposalVotes, proposalSubmission, closure)
	go propTally.HandleMessages()

	for _, a := range agents {
		a := a
		agentState := state.AgentState[a.BaseAgent.Id()]
		if a.BaseAgent.Id() == state.CurrentLeader {
			go (&a).HandleFight(agentState, previousDecisions, proposalVotes, proposalSubmission)
		} else {
			go (&a).HandleFight(agentState, previousDecisions, proposalVotes, nil)
		}
	}
	mId, _ := uuid.NewUUID()

	for _, messages := range channelsMap {
		messages <- *message.NewTaggedMessage("server", *message.NewMessage(message.Inform, nil), mId)
	}
	time.Sleep(200 * time.Millisecond)
	for _, c := range channelsMap {
		c <- *message.NewTaggedMessage("server", *message.NewMessage(message.Close, nil), mId)
		go func(recv <-chan message.TaggedMessage) {
			for m := range recv {
				switch m.Message().MType() {
				case message.Request:
					//todo: respond with nil thing here as we're closing! Or do we need to?
					// maybe because we're closing there's no point...
				default:
				}
			}
		}(c)
	}

	for _, c := range channelsMap {
		close(c)
	}

	closure <- struct{}{}
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
