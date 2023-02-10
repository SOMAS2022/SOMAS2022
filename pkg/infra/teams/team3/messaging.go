package team3

import (
	"fmt"
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/message"
)

// This is where you must compile your trust message. My example implementation takes ALL agents from the agent map **
func (a *AgentThree) CompileTrustMessage(agentMap map[commons.ID]agent.Agent) message.Trust {
	fmt.Println("AGENT 3 COMPOSED: message.Trust")

	//faireness = the function ids --> reputation number: which is the gossip
	//send to everyone
	keys := make([]commons.ID, len(agentMap))

	// ** it extracts the keys (i.e., the IDs) **
	i := 0
	for _, k := range a.TSN {
		keys[i] = k
		i++
	}

	for k := range agentMap {
		keys[i] = k
		i++
	}
	//fmt.Println("print keys: ", keys[:1])
	// declare new trust message
	trustMsg := new(message.Trust)

	// ** and puts stuff inside
	//trustMsg.MakeNewTrust(keys[:1], make(map[string]int))
	num := int(a.samplePercent * float64(len(agentMap)))
	trustMsg.MakeNewTrust(keys[0:num], a.reputationMap) //change the :1

	// // send off
	return *trustMsg
}

// You will receive a message of type "TaggedMessage"
func (a *AgentThree) HandleTrustMessage(m message.TaggedMessage) {
	// Receive the message.Trust type using m.Message()
	//fmt.Println("AGENT 3 RECEIVED: ", reflect.TypeOf(m))
	mes := m.Message()
	t := mes.(message.Trust)
	for key, value := range t.Gossip {
		//fmt.Println("lol: ", key)
		_, exists := a.reputationMap[key]
		if exists {
			diff := a.reputationMap[key] - t.Gossip[key]
			norm := diff * (a.reputationMap[m.Sender()] / 100)
			a.reputationMap[key] = a.reputationMap[key] + norm
		} else {
			a.reputationMap[key] = value
		}

	}
	a.socialCap[m.Sender()] += 1
	//fmt.Println("sender is", t.Recipients, m.Sender(), a.socialCap[m.Sender()])
	// This function is type void - you can do whatever you want with it. I would suggest keeping a local dictionary

}
