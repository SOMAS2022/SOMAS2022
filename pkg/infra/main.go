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
	agentMap := make(map[uint]agent.Agent)

	stateChannels := make(map[uint]chan<- state.State)
	decisionChannels := make(map[uint]<-chan decision.Decision)

	for i := uint(0); i < numAgents; i++ {
		//todo: add peer channels

		stateChan := make(chan state.State)
		decisionChan := make(chan decision.Decision)

		stateChannels[i] = stateChan
		decisionChannels[i] = decisionChan

		agentMap[i] = agent.Agent{
			BaseAgent: agent.BaseAgent{
				Communication: commons.Communication{
					Peer:     nil,
					Receiver: stateChan,
					Sender:   decisionChan,
				},
				Id: i,
			},
			Strategy: nil,
		}
	}

	for level := 0; level < numLevels; level++ {

	}
}

func handleFight(state *state.State, agents map[uint]agent.Agent, decisionChannels map[uint]<-chan decision.Decision) {
	for _, a := range agents {
		go a.Strategy.HandleFight(*state, a.BaseAgent)
	}
	decisions := make(map[uint]decision.FightChoice)

	for u, decisionC := range decisionChannels {
		handleFightDecision(decisionC, decisions, u)
	}

	//todo: process updates and modify state here

}

func handleFightDecision(decisionC <-chan decision.Decision, decisions map[uint]decision.FightChoice, u uint) {
	for {
		received := <-decisionC
		switch d := received.(type) {
		case decision.FightDecision:
			decisions[u] = d.Choice
			return
		default:
			continue
		}
	}
}
