package message

import (
	"fmt"
	"infra/game/commons"
	"infra/logging"

	"github.com/benbjohnson/immutable"
	"github.com/google/uuid"
)

/*
Each message should either be an info message ort a request messages,
request messages require the receiver to give a response in a certain amount of time

Each message should implement either a ProcessRequestMessage or ProcessInfoMessage method
These methods should call a method in the strategy interface on the strategy passed to them
This means that message types do not have to be switched by the agent implementations
*/

// TODO need to add a mechanism for agents to send messages

type Message interface {
	GenerateUUID()
	GetUUID() uuid.UUID
	SetUUID(UUID uuid.UUID)
}

// Base Message

type BaseMessage struct {
	UUID uuid.UUID
}

func (b *BaseMessage) GenerateUUID() {
	b.UUID = uuid.New()
}

func (b *BaseMessage) GetUUID() uuid.UUID {
	return b.UUID
}

func (b *BaseMessage) SetUUID(UUID uuid.UUID) {
	b.UUID = UUID
}

// Request Message force the receiver to respond in a given amount of time

/*
Message Structs
*/

/*
Associated messaging structs
*/

type TaggedMessage struct {
	Sender  commons.ID
	Message Message
}

type BaseAgent struct {
	Communication Communication
	Id            commons.ID
	AgentName     string
}

func NewBaseAgent(communication Communication, id commons.ID, agentName string) BaseAgent {
	return BaseAgent{Communication: communication, Id: id, AgentName: agentName}
}

type Communication struct {
	Receipt <-chan TaggedMessage
	peer    immutable.Map[commons.ID, chan<- TaggedMessage]
}

func NewCommunication(receipt <-chan TaggedMessage, peer immutable.Map[commons.ID, chan<- TaggedMessage]) Communication {
	return Communication{Receipt: receipt, peer: peer}
}

func (b *BaseAgent) broadcastBlockingMessage(m Message) {
	iterator := b.Communication.peer.Iterator()
	tm := TaggedMessage{
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

func SendBlockingMessageImp(b *BaseAgent, id commons.ID, m Message) (e error) {
	defer func() {
		if r := recover(); r != nil {
			e = fmt.Errorf("agent %s not available for messaging, submitted", id)
		}
	}()

	value, ok := b.Communication.peer.Get(id)
	if ok {
		value <- TaggedMessage{
			Sender:  b.Id,
			Message: m,
		}
	} else {
		e = fmt.Errorf("agent %s not available for messaging, dead", id)
	}
	return
}

/*
Available to the base agent implementations
*/
func (b *BaseAgent) sendBlockingMessage(id commons.ID, m Message) (e error) {
	m.SetUUID(uuid.New())
	return SendBlockingMessageImp(b, id, m)
}

func (b *BaseAgent) createLog(lvl logging.Level, fields logging.LogField, msg string) {
	agentFields := logging.LogField{
		"agentName": b.AgentName,
		"agentID":   b.Id,
	}

	logging.Log(lvl, logging.CombineFields(agentFields, fields), msg)
}
