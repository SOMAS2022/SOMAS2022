package trade

import (
	"fmt"
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/message"
	"infra/game/state"
	"infra/logging"
	"time"

	"github.com/benbjohnson/immutable"
)

// HandleTrade
// A complete trading stage contains several rounds.
// In each round, the following steps take place in order:
// 1. Each agent can respond to one of the trading negotiations it is involved in OR propose a new trade to another agent.
// 2. Main thread collects trade messages from all agents, and updated the state accordingly.
// 3. Collected message will be forwarded to corresponding target agents in the start of next round.
func HandleTrade(s state.State, agents map[commons.ID]agent.Agent, channelsMap map[commons.ID]chan message.TaggedMessage, round uint, roundLimit uint) {
	// track offers made by each agent, no repeated offers are allowed
	// i.e. only one offer of a specific item from an agent to another agent is allowed to exist simultaneously
	availableWeapons := make(map[commons.ID][]state.Item)
	availableShields := make(map[commons.ID][]state.Item)
	// track all ongoing negotiations
	negotiations := make(map[commons.TradeID]message.TradeNegotiation)

	// extract inventory from agents
	for agentID, agentState := range s.AgentState {
		availableWeapons[agentID] = commons.ImmutableListToSlice(agentState.Weapons)
		availableShields[agentID] = commons.ImmutableListToSlice(agentState.Shields)
	}

	for r := uint(0); r < round; r++ {
		starts := make(map[commons.ID]chan interface{})
		closures := make(map[commons.ID]chan interface{})
		responses := make(map[commons.ID]chan message.TradeMessage)

		for id, a := range agents {
			a := a
			agentState := s.AgentState[a.BaseAgent.ID()]
			availWeapons := commons.SliceToImmutableList(availableWeapons[id])
			availShields := commons.SliceToImmutableList(availableShields[id])
			requests := FindNegotiations(id, negotiations)

			start := make(chan interface{})
			starts[id] = start
			closure := make(chan interface{})
			closures[id] = closure
			response := make(chan message.TradeMessage)
			responses[id] = response

			go (&a).HandleTrade(agentState, availWeapons, availShields, &requests, start, closure, response)
		}
		// start all agents
		for _, startMessage := range starts {
			startMessage <- nil
		}
		// handle responses from agents
		for agentID, response := range responses {
			negotiation := <-response
			negotiations, availableWeapons, availableShields = HandleTradeMessage(agentID, negotiation, negotiations, availableWeapons, availableShields, s.AgentState)
		}
		// timeout for agents to respond
		time.Sleep(100 * time.Millisecond)
		for _, closeMessage := range closures {
			closeMessage <- nil
		}
		// filter out outdated negotiations
		for id, negotiation := range negotiations {
			negotiation.RoundNum++
			if negotiation.RoundNum > roundLimit {
				logging.Log(logging.Trace, nil, fmt.Sprintf("Negotiation %s between %s and %s is outdated", id, negotiation.Agent1, negotiation.Agent2))
				delete(negotiations, id)
			} else {
				negotiations[id] = negotiation
			}
		}
		logging.Log(logging.Info, logging.LogField{
			"round":          r,
			"numNegotiation": len(negotiations),
		}, fmt.Sprintf("Round %d: %d ongoing negotiations", r, len(negotiations)))
	}

	// End of trade stage, update agent inventory
	for agentID := range agents {
		agentState := s.AgentState[agentID]
		agentState.Weapons = *commons.SliceToImmutableList(availableWeapons[agentID])
		agentState.Shields = *commons.SliceToImmutableList(availableShields[agentID])
		s.AgentState[agentID] = agentState
	}
}
func HandleTradeMessage(agentID commons.ID, negotiation message.TradeMessage,
	negotiations map[commons.TradeID]message.TradeNegotiation,
	availableWeapons map[commons.ID][]state.Item,
	availableShields map[commons.ID][]state.Item,
	agentState map[commons.ID]state.AgentState,
) (
	newNegotiations map[commons.TradeID]message.TradeNegotiation,
	newAvailWeapons map[commons.ID][]state.Item,
	newAvailShields map[commons.ID][]state.Item,
) {
	newNegotiations = negotiations
	newAvailWeapons = availableWeapons
	newAvailShields = availableShields

	switch msg := negotiation.(type) {
	case message.TradeAbstain:
	case message.TradeResponse:
		newNegotiations, newAvailWeapons, newAvailShields = HandleTradeResponse(agentID, msg, negotiations, availableWeapons, availableShields, agentState)
	case message.TradeRequest:
		newNegotiations, newAvailWeapons, newAvailShields = HandleTradeRequest(agentID, msg, negotiations, availableWeapons, availableShields)
	}

	return newNegotiations, newAvailWeapons, newAvailShields
}

func HandleTradeRequest(agentID commons.ID, msg message.TradeRequest,
	negotiations map[commons.TradeID]message.TradeNegotiation,
	availableWeapons map[commons.ID][]state.Item,
	availableShields map[commons.ID][]state.Item,
) (
	newNegotiations map[commons.TradeID]message.TradeNegotiation,
	newAvailWeapons map[commons.ID][]state.Item,
	newAvailShields map[commons.ID][]state.Item,
) {
	newNegotiations = negotiations
	newAvailWeapons = availableWeapons
	newAvailShields = availableShields

	// add new negotiation to ongoing negotiations
	negotiation := message.NewTradeNegotiation(agentID, msg.CounterPartyID, msg.Offer, msg.Demand)
	newNegotiations[negotiation.Id] = negotiation
	// remove offered item from available items
	if msg.Offer.ItemType == commons.Weapon {
		newAvailWeapons[agentID] = RemoveItem(newAvailWeapons[agentID], msg.Offer.Item)
	} else {
		newAvailShields[agentID] = RemoveItem(newAvailShields[agentID], msg.Offer.Item)
	}

	return newNegotiations, newAvailWeapons, newAvailShields
}

func HandleTradeResponse(agentID commons.ID, msg message.TradeResponse,
	negotiations map[commons.TradeID]message.TradeNegotiation,
	availableWeapons map[commons.ID][]state.Item,
	availableShields map[commons.ID][]state.Item,
	agentState map[commons.ID]state.AgentState,
) (
	newNegotiations map[commons.TradeID]message.TradeNegotiation,
	newAvailWeapons map[commons.ID][]state.Item,
	newAvailShields map[commons.ID][]state.Item,
) {
	newNegotiations = negotiations
	newAvailWeapons = availableWeapons
	newAvailShields = availableShields

	switch resp := msg.(type) {
	case message.TradeAccept:
		negotiation := negotiations[resp.TradeID]
		if negotiation.Notarize(agentState) {
			newAvailWeapons, newAvailShields = ExecuteTrade(newAvailWeapons, newAvailShields, negotiation)
		}
		RemoveFromNegotiation(resp.TradeID, agentID, newNegotiations)
	case message.TradeReject:
		negotiation := negotiations[resp.TradeID]
		RemoveFromNegotiation(resp.TradeID, agentID, newNegotiations)
		newAvailWeapons, newAvailShields = PutBackItems(newAvailWeapons, newAvailShields, negotiation)
	case message.TradeBargain:
		negotiation := negotiations[resp.TradeID]
		if !negotiation.IsInvolved(agentID) {
			break
		}
		// update ongoing negotiations
		negotiation.UpdateDemand(agentID, resp.Demand)
		oldOffer, replaceOffer := negotiation.UpdateOffer(agentID, resp.Offer)
		newNegotiations[resp.TradeID] = negotiation
		// update available items
		if replaceOffer {
			if oldOffer.ItemType == commons.Weapon {
				newAvailWeapons[agentID] = append(newAvailWeapons[agentID], oldOffer.Item)
				newAvailWeapons[agentID] = RemoveItem(newAvailWeapons[agentID], resp.Offer.Item)
			} else {
				newAvailShields[agentID] = append(newAvailShields[agentID], oldOffer.Item)
				newAvailShields[agentID] = RemoveItem(newAvailShields[agentID], resp.Offer.Item)
			}
		}
	}
	return newNegotiations, newAvailWeapons, newAvailShields
}

func RemoveFromNegotiation(tradeID commons.TradeID, agentID commons.ID, negotiations map[commons.TradeID]message.TradeNegotiation) {
	negotiation := negotiations[tradeID]
	if negotiation.IsInvolved(agentID) {
		delete(negotiations, tradeID)
	}
}

func AddItem(available map[commons.ID][]state.Item, agentID commons.ID, item state.Item) map[commons.ID][]state.Item {
	availableList := available[agentID]
	availableList = append(availableList, item)

	available[agentID] = availableList
	return available
}

func RemoveItem(available []state.Item, item state.Item) []state.Item {
	for idx, availableItem := range available {
		if availableItem.Id() == item.Id() {
			available = append(available[:idx], available[idx+1:]...)
			return available
		}
	}
	return available
}

func PutBackItems(weapons map[commons.ID][]state.Item, shields map[commons.ID][]state.Item, negotiation message.TradeNegotiation) (newWeapons map[commons.ID][]state.Item, newShields map[commons.ID][]state.Item) {
	if negotiation.Condition1.Offer.IsValid {
		switch negotiation.Condition1.Offer.ItemType {
		case commons.Weapon:
			newWeapons = AddItem(weapons, negotiation.Agent1, negotiation.Condition1.Offer.Item)
		case commons.Shield:
			newShields = AddItem(shields, negotiation.Agent1, negotiation.Condition1.Offer.Item)
		}
	}
	if negotiation.Condition2.Offer.IsValid {
		switch negotiation.Condition2.Offer.ItemType {
		case commons.Weapon:
			newWeapons = AddItem(weapons, negotiation.Agent2, negotiation.Condition2.Offer.Item)
		case commons.Shield:
			newShields = AddItem(shields, negotiation.Agent2, negotiation.Condition2.Offer.Item)
		}
	}

	return
}

// switch the offered items between the two agents
// empty offer is allowed
func ExecuteTrade(weapons map[commons.ID][]state.Item, shields map[commons.ID][]state.Item, negotiation message.TradeNegotiation) (newWeapons map[commons.ID][]state.Item, newShields map[commons.ID][]state.Item) {
	condition1 := negotiation.Condition1
	if condition1.Offer.IsValid {
		switch condition1.Offer.ItemType {
		case commons.Weapon:
			newWeapons = AddItem(weapons, negotiation.Agent2, condition1.Offer.Item)
		case commons.Shield:
			newShields = AddItem(shields, negotiation.Agent2, condition1.Offer.Item)
		}
	}

	condition2 := negotiation.Condition2
	if condition2.Offer.IsValid {
		switch condition2.Offer.ItemType {
		case commons.Weapon:
			newWeapons = AddItem(weapons, negotiation.Agent1, condition2.Offer.Item)
		case commons.Shield:
			newShields = AddItem(shields, negotiation.Agent1, condition2.Offer.Item)
		}
	}

	return
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
