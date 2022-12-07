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
func HandleTrade(s state.State, agents map[commons.ID]agent.Agent, channelsMap map[commons.ID]chan message.TaggedMessage,
	roundLimit uint, perRoundLimit uint,
) {
	// track offers made by each agent, no repeated offers are allowed
	// i.e. only one offer of a specific item from an agent to another agent is allowed to exist simultaneously
	availableWeapons := make(map[commons.ID][]state.Item)
	availableShields := make(map[commons.ID][]state.Item)
	// track all ongoing negotiations
	negotiations := make(map[commons.TradeID]message.TradeNegotiation)

	// extract inventory from agents
	for agentID, agentState := range s.AgentState {
		availableWeapons[agentID] = ImmutableListToSlice(agentState.Weapons)
		availableShields[agentID] = ImmutableListToSlice(agentState.Shields)
	}

	for round := uint(0); round < roundLimit; round++ {
		starts := make(map[commons.ID]chan interface{})
		closures := make(map[commons.ID]chan interface{})
		responses := make(map[commons.ID]chan message.TradeMessage)

		for id, a := range agents {
			a := a
			agentState := s.AgentState[a.BaseAgent.ID()]
			availableWeapons := SliceToImmutableList(availableWeapons[id])
			availableShields := SliceToImmutableList(availableShields[id])
			requests := FindNegotiations(id, negotiations)

			start := make(chan interface{})
			starts[id] = start
			closure := make(chan interface{})
			closures[id] = closure
			response := make(chan message.TradeMessage)
			responses[id] = response

			go (&a).HandleTrade(agentState, roundLimit, availableWeapons, availableShields, &requests, start, closure, response)
		}
		// start all agents
		for _, startMessage := range starts {
			startMessage <- nil
		}
		// handle responses from agents
		for agentID, response := range responses {
			negotiation := <-response
			switch msg := negotiation.(type) {
			case message.TradeResponse:
				switch resp := msg.(type) {
				case message.TradeAccept:
					negotiation := negotiations[resp.TradeID]
					RemoveFromeNegotiation(resp.TradeID, agentID, negotiations)
					availableWeapons, availableShields = ExecuteTrade(availableWeapons, availableShields, negotiation)
				case message.TradeReject:
					negotiation := negotiations[resp.TradeID]
					RemoveFromeNegotiation(resp.TradeID, agentID, negotiations)
					availableWeapons, availableShields = PutBackFromNegotiation(availableWeapons, availableShields, negotiation)
				case message.TradeBargain:
					negotiation := negotiations[resp.TradeID]
					if !negotiation.IsInvolved(agentID) {
						break
					}
					// update ongoing negotiations
					oldOffer, replaceOffer := negotiation.UpdateOffer(agentID, resp.Offer)
					negotiations[resp.TradeID] = negotiation
					// update available items
					if replaceOffer {
						if oldOffer.ItemType == commons.Weapon {
							availableWeapons[agentID] = append(availableWeapons[agentID], oldOffer.Item)
							availableWeapons[agentID] = RemoveItem(availableWeapons[agentID], resp.Offer.Item)
						} else {
							availableShields[agentID] = append(availableShields[agentID], oldOffer.Item)
							availableShields[agentID] = RemoveItem(availableShields[agentID], resp.Offer.Item)
						}
					}
				}
			case message.TradeRequest:
				counterPartyID, offer, demand := msg.CounterPartyID, msg.Offer, msg.Demand
				// add new negotiation to ongoing negotiations
				negotiation := message.NewTradeNegotiation(agentID, counterPartyID, offer, demand)
				negotiations[negotiation.GetID()] = negotiation
				// remove offered item from available items
				if offer.ItemType == commons.Weapon {
					availableWeapons[agentID] = RemoveItem(availableWeapons[agentID], offer.Item)
				} else {
					availableShields[agentID] = RemoveItem(availableWeapons[agentID], offer.Item)
				}
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

	for agentID := range agents {
		agentState := s.AgentState[agentID]
		agentState.Weapons = *SliceToImmutableList(availableWeapons[agentID])
		agentState.Shields = *SliceToImmutableList(availableShields[agentID])
		s.AgentState[agentID] = agentState
	}
}

func RemoveFromeNegotiation(tradeID commons.TradeID, agentID commons.ID, negotiations map[commons.TradeID]message.TradeNegotiation) {
	negotiation := negotiations[tradeID]
	if negotiation.IsInvolved(agentID) {
		delete(negotiations, tradeID)
	}
}

func ImmutableListToSlice(list immutable.List[state.Item]) []state.Item {
	slice := make([]state.Item, list.Len())
	for i := 0; i < list.Len(); i++ {
		slice[i] = list.Get(i)
	}
	return slice
}

func SliceToImmutableList(slice []state.Item) *immutable.List[state.Item] {
	list := immutable.NewListBuilder[state.Item]()
	for _, item := range slice {
		list.Append(item)
	}
	return list.List()
}

func AddItem(available map[commons.ID][]state.Item, agentID commons.ID, item state.Item) map[commons.ID][]state.Item {
	availableList := available[agentID]
	availableList = append(availableList, item)

	available[agentID] = availableList
	return available
}

func RemoveItem(available []state.Item, item state.Item) []state.Item {
	for idx, availableItem := range available {
		if availableItem == item {
			available = append(available[:idx], available[idx+1:]...)
			return available
		}
	}
	return available
}

func PutBackFromNegotiation(weapons map[commons.ID][]state.Item, shields map[commons.ID][]state.Item, negotiation message.TradeNegotiation) (newWeapons map[commons.ID][]state.Item, newShields map[commons.ID][]state.Item) {
	conditions := negotiation.GetConditions()
	iter := conditions.Iterator()
	for !iter.Done() {
		agentID, condition, _ := iter.Next()
		switch condition.GetItemType() {
		case commons.Weapon:
			newWeapons = AddItem(weapons, agentID, condition.GetItem())
		case commons.Shield:
			newShields = AddItem(shields, agentID, condition.GetItem())
		}
	}

	return newWeapons, newShields
}

func ExecuteTrade(weapons map[commons.ID][]state.Item, shields map[commons.ID][]state.Item, negotiation message.TradeNegotiation) (newWeapons map[commons.ID][]state.Item, newShields map[commons.ID][]state.Item) {
	conditions := negotiation.GetConditions()
	iter := conditions.Iterator()
	agentIDs := []commons.ID{}

	for !iter.Done() {
		agentID, _, _ := iter.Next()
		agentIDs = append(agentIDs, agentID)
	}

	agentID1 := agentIDs[0]
	agentID2 := agentIDs[1]

	condition1, _ := conditions.Get(agentID1)
	switch condition1.GetItemType() {
	case commons.Weapon:
		newWeapons = AddItem(weapons, agentID2, condition1.GetItem())
	case commons.Shield:
		newShields = AddItem(shields, agentID2, condition1.GetItem())
	}

	condition2, _ := conditions.Get(agentID2)
	switch condition2.GetItemType() {
	case commons.Weapon:
		newWeapons = AddItem(weapons, agentID2, condition2.GetItem())
	case commons.Shield:
		newShields = AddItem(shields, agentID2, condition2.GetItem())
	}

	return newWeapons, newShields
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
