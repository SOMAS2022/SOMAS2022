package main

import (
	"infra/server/agent"
	"infra/server/commons"
	"infra/server/decision"
	"infra/server/state"
)

/*
Each agent is initialised with these modes of communication
Whereby the peer channels enable bidirectional p2p communication
Enables peers to send and receive messages with broadcasting possible via non-blocking sends
*/

const numLevels = 60

const numAgents = 10

func main() {
	agentMap := make(map[uint]agent.BaseAgent)

	stateChannels := make(map[uint]chan<- state.State)
	decisionChannels := make(map[uint]<-chan decision.Decision)

	for i := uint(0); i < numAgents; i++ {
		//todo: add peer channels

		stateChan := make(chan state.State)
		decisionChan := make(chan decision.Decision)

		stateChannels[i] = stateChan
		decisionChannels[i] = decisionChan

		agentMap[i] = agent.BaseAgent{Communication: commons.Communication{
			Peer:     nil,
			Receiver: stateChan,
			Sender:   decisionChan,
		},
			Id: i}
	}

	for level := 0; level < numLevels; level++ {

	}
}

func handleFight(state *state.State, agents map[uint]agent.BaseAgent, decisionChannels map[uint]<-chan decision.Decision) {
	for _, baseAgent := range agents {
		go baseAgent.Logic.HandleFight(*state, baseAgent)
	}
	decisions := make(map[uint]decision.FightDecision)

	for u, decisionC := range decisionChannels {
		handleFightDecision(decisionC, decisions, u)
	}

	//todo: process updates and modify state here

}

func handleFightDecision(decisionC <-chan decision.Decision, decisions map[uint]decision.FightDecision, u uint) {
	for {
		received := <-decisionC
		switch d := received.(type) {
		case decision.FightDecision:
			decisions[u] = d
			return
		default:
			continue
		}
	}
}
