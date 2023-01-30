package team3

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/message"
)

// This is where you must compile your trust message. My example implementation takes ALL agents from the agent map **
func (a *AgentThree) CompileTrustMessage(agentMap map[commons.ID]agent.Agent) message.Trust {

	// fmt.Println("AGENT 3 COMPOSED: message.Trust")

	keys := make([]commons.ID, len(agentMap))

	// ** it extracts the keys (i.e., the IDs) **
	i := 0
	for k := range agentMap {
		keys[i] = k
		i++
	}

	// declare new trust message
	trustMsg := new(message.Trust)

	// ** and puts stuff inside
	trustMsg.MakeNewTrust(keys[:1], make(map[string]int))

	// send off
	return *trustMsg
}

// You will receive a message of type "TaggedMessage"
func (a *AgentThree) HandleTrustMessage(m message.TaggedMessage) {
	// Receive the message.Trust type using m.Message()
	// fmt.Println("AGENT 3 RECEIVED: ", reflect.TypeOf(m.Message()))
	// This function is type void - you can do whatever you want with it. I would suggest keeping a local dictionary
}
