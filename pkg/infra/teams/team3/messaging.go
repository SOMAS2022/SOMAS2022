package team3

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/message"
	"sync"
)

func agentInList(agentID commons.ID, messageList []commons.ID) bool {
	for _, id := range messageList {
		if agentID == id {
			return true
		}
	}
	return false
}

// This is where you must compile your trust message. My example implementation takes ALL agents from the agent map **
func (a *AgentThree) CompileTrustMessage(agentMap map[commons.ID]agent.Agent) message.Trust {
	// fmt.Println("AGENT 3 COMPOSED: message.Trust")

	// faireness = the function ids --> reputation number: which is the gossip
	num := int(a.samplePercent * float64(len(agentMap)))
	agentsToMessage := make([]commons.ID, num)

	// Check TSN first
	i := 0
	for _, k := range a.TSN {
		if i == num {
			break
		}
		agentsToMessage[i] = k
		i++
	}

	// Then fill remaining spots from the rest
	for k := range agentMap {
		if i == num {
			break
		}

		// try next ID if agent is already being messaged
		if agentInList(k, agentsToMessage) {
			continue
		}

		agentsToMessage[i] = k
		i++
	}

	// declare new trust message
	trustMsg := new(message.Trust)

	// avoid concurrent write to/read from trust.Gossip
	repMapDeepCopy := make(map[commons.ID]float64)

	for k, v := range a.reputationMap {
		repMapDeepCopy[k] = v
	}

	trustMsg.MakeNewTrust(agentsToMessage, repMapDeepCopy)

	// send off
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
