package agent

import (
	"errors"
	"fmt"
	"infra/game/decision"
	"infra/game/message/proposal"

	"infra/game/commons"
	"infra/game/message"
	"infra/game/state"
	"infra/logging"

	"github.com/google/uuid"
)

var errCommunication = errors.New("communicationError")

func communicationError(msg string) error {
	return fmt.Errorf("%w: %s", errCommunication, msg)
}

type BaseAgent struct {
	communication *Communication
	id            commons.ID
	name          string
	latestState   state.AgentState
	view          *state.View
	loot          state.LootPool
}

func (ba *BaseAgent) Loot() state.LootPool {
	return ba.loot
}

func (ba *BaseAgent) setCommunication(communication *Communication) {
	ba.communication = communication
}

func (ba *BaseAgent) View() state.View {
	return *ba.view
}

func (ba *BaseAgent) ID() commons.ID {
	return ba.id
}

func (ba *BaseAgent) Name() string {
	return ba.name
}

func NewBaseAgent(communication *Communication, id commons.ID, agentName string, ptr *state.View) *BaseAgent {
	return &BaseAgent{communication: communication, id: id, name: agentName, view: ptr}
}

func (ba *BaseAgent) BroadcastBlockingMessage(m message.Message) {
	iterator := ba.communication.peer.Iterator()
	tm := message.NewTaggedMessage(ba.id, m, uuid.New())

	for !iterator.Done() {
		_, c, ok := iterator.Next()
		if ok {
			c <- *tm
		}
	}
}

func (ba *BaseAgent) SendBlockingMessage(id commons.ID, m message.Message) (e error) {
	switch m.(type) {
	case message.Proposal[decision.FightAction]:
		return communicationError("Illegal attempt to send proposal - use SendFightProposalToLeader() instead")
	case message.Proposal[decision.LootAction]:
		return communicationError("Illegal attempt to send proposal - use SendLootProposalToLeader() instead")
	default:
		channel, ok := ba.communication.peer.Get(id)
		if ok {
			channel <- *message.NewTaggedMessage(ba.id, m, uuid.New())
		} else {
			return communicationError(fmt.Sprintf("agent %s not available for messaging", id))
		}
	}
	return nil
}

func (ba *BaseAgent) SendFightProposalToLeader(rules commons.ImmutableList[proposal.Rule[decision.FightAction]]) error {
	channel, ok := ba.communication.peer.Get(ba.view.CurrentLeader())
	if ok {
		channel <- *message.NewTaggedMessage(ba.id, *message.NewProposal(rules, ba.ID()), uuid.New())
		return nil
	}
	return communicationError("Leader not available for messaging, dead or bad!")
}

func (ba *BaseAgent) SendLootProposalToLeader(rules commons.ImmutableList[proposal.Rule[decision.LootAction]]) error {
	channel, ok := ba.communication.peer.Get(ba.view.CurrentLeader())
	if ok {
		channel <- *message.NewTaggedMessage(ba.id, *message.NewProposal(rules, ba.ID()), uuid.New())
		return nil
	}
	return communicationError("Leader not available for messaging, dead or bad!")
}

func (ba *BaseAgent) Log(lvl logging.Level, fields logging.LogField, msg string) {
	agentFields := logging.LogField{
		"agentName": ba.name,
		"agentID":   ba.id,
	}

	logging.Log(lvl, logging.CombineFields(agentFields, fields), msg)
}

func (ba *BaseAgent) AgentState() state.AgentState {
	return ba.latestState
}
