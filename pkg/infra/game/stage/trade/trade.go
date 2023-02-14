package trade

import (
	"fmt"
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/message"
	"infra/game/stage/trade/internal"
	"infra/game/state"
	"infra/logging"
	"time"
)

// HandleTrade
// A complete trading stage contains several rounds.
// In each round, the following steps take place in order:
// 1. Each agent can respond to one of the trading negotiations it is involved in OR propose a new trade to another agent.
// 2. Main thread collects trade messages from all agents, and updated the state accordingly.
// 3. Collected message will be forwarded to corresponding target agents in the start of next round.
func HandleTrade(s state.State, agents map[commons.ID]agent.Agent, round uint, roundLimit uint) {
	// track offers made by each agent, no repeated offers are allowed
	// i.e. only one offer of a specific item from an agent to another agent is allowed to exist simultaneously
	availableWeapons := make(map[commons.ID][]state.Item)
	availableShields := make(map[commons.ID][]state.Item)
	// track all ongoing negotiations
	negotiations := make(map[commons.TradeID]message.TradeNegotiation)
	info := internal.NewInfo(negotiations, *internal.NewInventory(availableWeapons, availableShields))
	// extract inventory from agents
	for agentID, agentState := range s.AgentState {
		info.Inventory.Weapons()[agentID] = commons.ImmutableListToSlice(agentState.Weapons)
		info.Inventory.Shields()[agentID] = commons.ImmutableListToSlice(agentState.Shields)
	}

	for r := uint(0); r < round; r++ {
		starts := make(map[commons.ID]chan interface{})
		closures := make(map[commons.ID]chan interface{})
		responses := make(map[commons.ID]chan message.TradeMessage)

		for id, a := range agents {
			a := a
			start := make(chan interface{})
			starts[id] = start
			closure := make(chan interface{})
			closures[id] = closure
			response := make(chan message.TradeMessage)
			responses[id] = response

			go (&a).HandleTrade(s.AgentState[a.BaseAgent.ID()], NewTradeInfo(id, info), start, closure, response)
		}
		// start all agents
		for _, startMessage := range starts {
			startMessage <- nil
		}
		// handle responses from agents
		for agentID, response := range responses {
			negotiation := <-response
			HandleTradeMessage(agentID, negotiation, info, s.AgentState)
		}
		// timeout for agents to respond
		time.Sleep(25 * time.Millisecond)
		for id, closure := range closures {
			closure <- nil
			close(closure)
			close(starts[id])
			close(responses[id])
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
		// logging.Log(logging.Info, logging.LogField{
		// 	"round":          r,
		// 	"numNegotiation": len(negotiations),
		// }, fmt.Sprintf("Round %d: %d ongoing negotiations", r, len(negotiations)))
	}

	// End of trade stage, update agent inventory
	for agentID := range agents {
		agentState := s.AgentState[agentID]
		agentState.Weapons = *commons.SliceToImmutableList(availableWeapons[agentID])
		agentState.Shields = *commons.SliceToImmutableList(availableShields[agentID])
		s.AgentState[agentID] = agentState
	}
}

func NewTradeInfo(agentID commons.ID, info *internal.Info) message.TradeInfo {
	return message.TradeInfo{
		Negotiations: FindNegotiations(agentID, info.Negotiations()),
		Weapons:      commons.ListToImmutableList(info.Inventory.Weapons()[agentID]),
		Shields:      commons.ListToImmutableList(info.Inventory.Shields()[agentID]),
	}
}

func HandleTradeMessage(agentID commons.ID, negotiation message.TradeMessage,
	info *internal.Info,
	agentState map[commons.ID]state.AgentState,
) {
	switch msg := negotiation.(type) {
	case message.TradeAbstain:
	case message.TradeResponse:
		HandleTradeResponse(agentID, msg, info, agentState)
	case message.TradeRequest:
		HandleTradeRequest(agentID, msg, info)
	}
}

func HandleTradeRequest(agentID commons.ID, msg message.TradeRequest,
	info *internal.Info,
) {
	// add new negotiation to ongoing negotiations
	negotiation := message.NewTradeNegotiation(agentID, msg.CounterPartyID, msg.Offer, msg.Demand)
	info.Negotiations()[negotiation.Id] = negotiation
	// remove offered item from available items
	if msg.Offer.ItemType == commons.Weapon {
		info.Weapons()[agentID] = RemoveItem(info.Weapons()[agentID], msg.Offer.Item)
	} else {
		info.Shields()[agentID] = RemoveItem(info.Shields()[agentID], msg.Offer.Item)
	}
}

func HandleTradeResponse(agentID commons.ID, msg message.TradeResponse,
	info *internal.Info,
	agentState map[commons.ID]state.AgentState,
) {
	switch resp := msg.(type) {
	case message.TradeAccept:
		negotiation := info.Negotiations()[resp.TradeID]
		if negotiation.Notarize(agentState) {
			ExecuteTrade(&info.Inventory, negotiation)
		}
		RemoveFromNegotiation(resp.TradeID, agentID, info.Negotiations())
	case message.TradeReject:
		negotiation := info.Negotiations()[resp.TradeID]
		RemoveFromNegotiation(resp.TradeID, agentID, info.Negotiations())
		PutBackItems(&info.Inventory, negotiation)
	case message.TradeBargain:
		negotiation := info.Negotiations()[resp.TradeID]
		if !negotiation.IsInvolved(agentID) {
			return
		}
		// check if offered item is still available
		ItemIsAvailable(info.Inventory, agentID, resp.Offer)
		// update ongoing negotiations
		negotiation.UpdateDemand(agentID, resp.Demand)
		oldOffer, replaceOffer := negotiation.UpdateOffer(agentID, resp.Offer)
		info.Negotiations()[resp.TradeID] = negotiation
		// update available items
		if replaceOffer {
			if oldOffer.ItemType == commons.Weapon {
				info.Weapons()[agentID] = append(info.Weapons()[agentID], oldOffer.Item)
				info.Weapons()[agentID] = RemoveItem(info.Weapons()[agentID], resp.Offer.Item)
			} else {
				info.Shields()[agentID] = append(info.Shields()[agentID], oldOffer.Item)
				info.Shields()[agentID] = RemoveItem(info.Shields()[agentID], resp.Offer.Item)
			}
		}
	}
}

func RemoveFromNegotiation(tradeID commons.TradeID, agentID commons.ID, negotiations map[commons.TradeID]message.TradeNegotiation) {
	negotiation := negotiations[tradeID]
	if negotiation.IsInvolved(agentID) {
		delete(negotiations, tradeID)
	}
}

func ItemIsAvailable(inventory internal.Inventory, agentID commons.ID, offer message.TradeOffer) bool {
	switch offer.ItemType {
	case commons.Weapon:
		return ContainsItem(inventory.Shields()[agentID], agentID, offer.Item)
	case commons.Shield:
		return ContainsItem(inventory.Shields()[agentID], agentID, offer.Item)
	}
	return false
}

func ContainsItem(inventory []state.Item, agentID commons.ID, item state.Item) bool {
	for _, v := range inventory {
		if v.Id() == item.Id() {
			return true
		}
	}
	return false
}

func AddItem(available map[commons.ID][]state.Item, agentID commons.ID, item state.Item) {
	availableList := available[agentID]
	availableList = append(availableList, item)

	available[agentID] = availableList
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

func PutBackItems(inventory *internal.Inventory, negotiation message.TradeNegotiation) {
	if negotiation.Condition1.Offer.IsValid {
		switch negotiation.Condition1.Offer.ItemType {
		case commons.Weapon:
			AddItem(inventory.Weapons(), negotiation.Agent1, negotiation.Condition1.Offer.Item)
		case commons.Shield:
			AddItem(inventory.Shields(), negotiation.Agent1, negotiation.Condition1.Offer.Item)
		}
	}
	if negotiation.Condition2.Offer.IsValid {
		switch negotiation.Condition2.Offer.ItemType {
		case commons.Weapon:
			AddItem(inventory.Weapons(), negotiation.Agent2, negotiation.Condition2.Offer.Item)
		case commons.Shield:
			AddItem(inventory.Shields(), negotiation.Agent2, negotiation.Condition2.Offer.Item)
		}
	}
}

// ExecuteTrade
// switch the offered items between the two agents
// empty offer is allowed
func ExecuteTrade(inventory *internal.Inventory, negotiation message.TradeNegotiation) {
	condition1 := negotiation.Condition1
	if condition1.Offer.IsValid {
		switch condition1.Offer.ItemType {
		case commons.Weapon:
			AddItem(inventory.Weapons(), negotiation.Agent2, condition1.Offer.Item)
		case commons.Shield:
			AddItem(inventory.Shields(), negotiation.Agent2, condition1.Offer.Item)
		}
	}

	condition2 := negotiation.Condition2
	if condition2.Offer.IsValid {
		switch condition2.Offer.ItemType {
		case commons.Weapon:
			AddItem(inventory.Weapons(), negotiation.Agent1, condition2.Offer.Item)
		case commons.Shield:
			AddItem(inventory.Shields(), negotiation.Agent1, condition2.Offer.Item)
		}
	}
}

// FindNegotiations
// Find all negotiations that the given agent is involved in
func FindNegotiations(agentID commons.ID, negotiations map[commons.TradeID]message.TradeNegotiation) map[commons.TradeID]message.TradeNegotiation {
	result := make(map[commons.TradeID]message.TradeNegotiation)
	for tradeID, negotiation := range negotiations {
		if negotiation.IsInvolved(agentID) {
			result[tradeID] = negotiation
		}
	}
	return result
}
