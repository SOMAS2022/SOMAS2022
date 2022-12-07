package trade

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/message"
	"infra/game/state"
	"time"
)

// A complete trading(exchange) stage contains several rounds.
// In each round, the following steps take place in order:
// 1. Each agent can propose an exchange of item(weapon or shield) to another agent.
// 2. Main thread collects exchange proposals from all agents
// 3. Collected proposals are forwarded to corresponding target agents
// 4. Target agents response if they want to accept the proposal and what they are willing to offer in return.
// 5. Agents initiated the exchange decide whether to accept the exchange.

func HandleTrade(state state.State, agents map[commons.ID]agent.Agent, channelsMap map[commons.ID]chan message.TaggedMessage, roundLimit uint) {
	starts := make(map[commons.ID]chan<- interface{})
	closures := make(map[commons.ID]chan<- interface{})

	for id, a := range agents {
		a := a
		agentState := state.AgentState[a.BaseAgent.ID()]

		start := make(chan interface{})
		starts[id] = start
		closure := make(chan interface{})
		closures[id] = closure

		go (&a).HandleTrade(agentState, roundLimit, start, closure)
	}
	for _, startMessage := range starts {
		startMessage <- nil
	}
	time.Sleep(100 * time.Millisecond)
	for _, closeMessage := range closures {
		closeMessage <- nil
	}
}
