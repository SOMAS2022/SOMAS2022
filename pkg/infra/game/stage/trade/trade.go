package trade

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/message"
	"infra/game/state"
	"time"

	"github.com/benbjohnson/immutable"
)

// A complete trading stage contains several rounds.
// In each round, the following steps take place in order:
// 1. Each agent can response to one of the trading negotiations it is involved in OR propose a new trade to another agent.
// 2. Main thread collects all trade messages from all agents, and updated the state accordingly.
// 3. Collected message will be forwarded to corresponding target agents in the start of next round.
func HandleTrade(state state.State, agents map[commons.ID]agent.Agent, channelsMap map[commons.ID]chan message.TaggedMessage,
	roundLimit uint, perRoundLimit uint,
) {
	// track offers made by each agent, no repeated offers are allowed
	// i.e. only one offer of a specific item from an agent to another agent is allowed to exist simultaneously
	offers := make(map[commons.ID][]commons.ItemID)
	// track all ongoing negotiations
	negotiations := make(map[commons.TradeID]message.TradeNegotiation)

	for round := uint(0); round < roundLimit; round++ {
		starts := make(map[commons.ID]chan interface{})
		closures := make(map[commons.ID]chan interface{})
		responses := make(map[commons.ID]chan message.TradeMessage)

		for id, a := range agents {
			a := a
			agentState := state.AgentState[a.BaseAgent.ID()]
			offeredItems := commons.ListToImmutable(offers[id])
			requests := FindNegotiations(id, negotiations)

			start := make(chan interface{})
			starts[id] = start
			closure := make(chan interface{})
			closures[id] = closure
			response := make(chan message.TradeMessage)
			responses[id] = response

			go (&a).HandleTrade(agentState, roundLimit, &offeredItems, &requests, start, closure, response)
		}
		// start all agents
		for _, startMessage := range starts {
			startMessage <- nil
		}
		// handle responses from agents
		for _, response := range responses {
			negotiation := <-response
			switch msg := negotiation.(type) {
			case message.TradeResponse:
				switch msg.(type) {
				case message.TradeAccept:
					// TODO: remove this negotiation from ongoing negotiations
					// TODO: update both agents' inventory
					// TODO: remove the offered item from offers
				case message.TradeReject:
					// TODO: remove this negotiation from ongoing negotiations
					// TODO: remove the offered item from offers
				case message.TradeBargain:
					// TODO: update this negotiation in ongoing negotiations
					// TODO: update offers
				}
			case message.TradeRequest:
				// TODO: create a new negotiation and add it to ongoing negotiations
				// TODO: update offers
			}
		}
		// timeout for agents to respond
		time.Sleep(100 * time.Millisecond)
		for _, closeMessage := range closures {
			closeMessage <- nil
		}
		// filterout outdated negotiations
		for id, negotiation := range negotiations {
			if negotiation.UpdateRoundNum() >= roundLimit {
				delete(negotiations, id)
			}
		}
	}
}

// Find all negotiations that the given agent is involved in
func FindNegotiations(agentID commons.ID, negotiations map[commons.TradeID]message.TradeNegotiation) immutable.Map[commons.TradeID, message.TradeNegotiation] {
	b := immutable.NewMapBuilder[commons.TradeID, message.TradeNegotiation](nil)
	for tradeID, negotiation := range negotiations {
		if negotiation.IsInvolved(agentID) {
			b.Set(tradeID, negotiation)
		}
	}
	return *b.Map()
}
