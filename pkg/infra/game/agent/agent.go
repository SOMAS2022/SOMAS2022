package agent

import (
	"fmt"
	"github.com/google/uuid"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"infra/game/state"
	"infra/logging"
	"sync"

	"github.com/benbjohnson/immutable"
)

// type Strategy interface {
// 	// HandleFightMessage(m message.TaggedMessage,
// 	// 	view *state.View, agent BaseAgent,
// 	// 	log *immutable.Map[commons.ID, decision.FightAction],
// 	// )
// 	HandleInfoFightMessage(
// 		m message.TaggedMessage,
// 		view *state.View,
// 		agent BaseAgent,
// 		log *immutable.Map[commons.ID, decision.FightAction],
// 	)
// 	HandleResponseFightMessage(m message.TaggedMessage,
// 		view *state.View,
// 		agent BaseAgent,
// 		log *immutable.Map[commons.ID, decision.FightAction],
// 	) message.Message

// 	Default() decision.FightAction
// }

type Agent struct {
	BaseAgent BaseAgent
	Strategy  message.Strategy
}

func (a *Agent) HandleFight(
	agentMap map[commons.ID]Agent,
	view state.View,
	log immutable.Map[commons.ID, decision.FightAction],
	decisionChan chan message.ActionDecision,
	wg *sync.WaitGroup) {

	// Give the agent the current state to process
	a.Strategy.ProcessStartOfRound(&view, &log)

	// Process any messages that the agent currently has in a loop
	// TODO ? keep looping this for loop until all agents are completed
	for m := range a.BaseAgent.communication.receipt {
		a.handleMessage(agentMap, &view, &log, m)
		action := a.Strategy.GenerateActionDecision()
		if action != decision.Undecided {
			go func() {
				// TODO add switch here on message type
				<-a.BaseAgent.communication.receipt
			}()
			decisionChan <- message.ActionDecision{Action: action, Sender: a.BaseAgent.Id}
			wg.Done()
			return
		}
	}

	// Ask the agent for its final decision
	decisionChan <- message.ActionDecision{Action: a.Strategy.GenerateActionDecision(), Sender: a.BaseAgent.Id}
}

func (a *Agent) handleMessage(agentMap map[commons.ID]Agent, view *state.View, log *immutable.Map[commons.ID, decision.FightAction], m message.TaggedMessage) {
	switch m.Message.(type) {
	case message.RequestMessageInterface:
		// TODO add timeout in which agent must reply
		response := m.Message.(message.RequestMessageInterface).ProcessRequestMessage(a.Strategy, view, log)
		agentMap[m.Sender].BaseAgent._sendBlockingResponseMessage(a.BaseAgent.Id, response, m)
	case message.InfoMessageInterface:
		m.Message.(message.InfoMessageInterface).ProcessInfoMessage(a.Strategy, view, log)
	default:
	}
}

func (ba *BaseAgent) log(lvl logging.Level, fields logging.LogField, msg string) {
	agentFields := logging.LogField{
		"agentName": ba.AgentName,
		"agentID":   ba.Id,
	}

	logging.Log(lvl, logging.CombineFields(agentFields, fields), msg)
}

type BaseAgent struct {
	communication Communication
	Id            commons.ID
	AgentName     string
}

func NewBaseAgent(communication Communication, id commons.ID, agentName string) BaseAgent {
	return BaseAgent{communication: communication, Id: id, AgentName: agentName}
}

type Communication struct {
	receipt <-chan message.TaggedMessage
	peer    immutable.Map[commons.ID, chan<- message.TaggedMessage]
}

func NewCommunication(receipt <-chan message.TaggedMessage, peer immutable.Map[commons.ID, chan<- message.TaggedMessage]) Communication {
	return Communication{receipt: receipt, peer: peer}
}

func (b BaseAgent) broadcastBlockingMessage(m message.Message) {
	iterator := b.communication.peer.Iterator()
	tm := message.TaggedMessage{
		Sender:  b.Id,
		Message: m,
	}
	for !iterator.Done() {
		_, c, ok := iterator.Next()
		if ok {
			c <- tm
		}
	}
}

func (b BaseAgent) _sendBlockingResponseMessage(id commons.ID, response message.Message, request message.TaggedMessage) (e error) {
	defer func() {
		if r := recover(); r != nil {
			e = fmt.Errorf("agent %s not available for messaging, submitted", id)
		}
	}()

	value, ok := b.communication.peer.Get(id)
	if ok {
		value <- message.TaggedMessage{
			Sender:  b.Id,
			Message: response,
			UUID:    request.UUID,
		}
	} else {
		e = fmt.Errorf("agent %s not available for messaging, dead", id)
	}
	return
}

func (b BaseAgent) sendBlockingMessage(id commons.ID, m message.Message) (e error) {
	defer func() {
		if r := recover(); r != nil {
			e = fmt.Errorf("agent %s not available for messaging, submitted", id)
		}
	}()

	value, ok := b.communication.peer.Get(id)
	if ok {
		value <- message.TaggedMessage{
			Sender:  b.Id,
			Message: m,
			UUID:    uuid.New(),
		}
	} else {
		e = fmt.Errorf("agent %s not available for messaging, dead", id)
	}
	return
}
