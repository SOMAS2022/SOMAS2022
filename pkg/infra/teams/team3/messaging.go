package team3

import (
	"fmt"
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/message"
	"sync"
)

// This is where you must compile your trust message. My example implementation takes ALL agents from the agent map **
func (a *AgentThree) CompileTrustMessage(agentMap map[commons.ID]agent.Agent) message.Trust {
	fmt.Println("AGENT 3 COMPOSED: message.Trust")

	//faireness = the function ids --> reputation number: which is the gossip
	//send to everyone
	keys := make([]commons.ID, len(agentMap)+len(a.TSN))

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

	// avoid concurrent write to/read from trust.Gossip
	repMapShallowCopy := make(map[commons.ID]float64)

	for k, v := range a.reputationMap {
		repMapShallowCopy[k] = v
	}

	trustMsg.MakeNewTrust(keys[0:num], repMapShallowCopy) //change the :1

	// // send off
	return *trustMsg
}

// You will receive a message of type "TaggedMessage"
func (a *AgentThree) HandleTrustMessage(m message.TaggedMessage) {
	// Receive the message.Trust type using m.Message()

	mes := m.Message()
	t := mes.(message.Trust)

	mutex := sync.Mutex{}
	mutex.Lock()

	// Gossip IS reputation map ---> one thread will read it, one thread will write it. 
	// Shallow copy introduced in Compile

	for key, value := range t.Gossip {
		rep, exists := a.reputationMap[key]
		if exists {
			diff := rep - value
			norm := diff * (a.reputationMap[m.Sender()] / 100)
			a.reputationMap[key] = rep + norm
		} else {
			a.reputationMap[key] = value
		}

	}
	a.socialCap[m.Sender()] += 1
	mutex.Unlock()
	//fmt.Println("sender is", t.Recipients, m.Sender(), a.socialCap[m.Sender()])
	// This function is type void - you can do whatever you want with it. I would suggest keeping a local dictionary

}
