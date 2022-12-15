package message

import (
	"fmt"
	"infra/game/commons"
	"infra/game/state"

	"github.com/google/uuid"
)

type TradeNegotiation struct {
	Id         commons.TradeID
	RoundNum   uint // increment in each round, terminate negotiation when this number reaches limit
	Agent1     commons.ID
	Agent2     commons.ID
	Condition1 TradeCondition
	Condition2 TradeCondition
}

func NewTradeNegotiation(agentID commons.ID, counterPartyID commons.ID, offer TradeOffer, demand TradeDemand) TradeNegotiation {
	condition := TradeCondition{
		Offer:  offer,
		Demand: demand,
	}
	return TradeNegotiation{
		Id:         uuid.New().String(),
		Agent1:     agentID,
		Agent2:     counterPartyID,
		RoundNum:   0,
		Condition1: condition,
	}
}

func (negotiation *TradeNegotiation) sealedMessage() {}

func (negotiation *TradeNegotiation) IsInvolved(agentID commons.ID) bool {
	return negotiation.Agent1 == agentID || negotiation.Agent2 == agentID
}

func (negotiation *TradeNegotiation) GetOffer(agentID commons.ID) (TradeOffer, bool) {
	switch agentID {
	case negotiation.Agent1:
		return negotiation.Condition1.Offer, true
	case negotiation.Agent2:
		return negotiation.Condition2.Offer, true
	default:
		return TradeOffer{}, false
	}
}

func (negotiation *TradeNegotiation) GetDemand(agentID commons.ID) (demand TradeDemand, ok bool) {
	switch agentID {
	case negotiation.Agent1:
		return negotiation.Condition1.Demand, true
	case negotiation.Agent2:
		return negotiation.Condition2.Demand, true
	default:
		return TradeDemand{}, false
	}
}

func (negotiation *TradeNegotiation) GetCounterParty(agentID commons.ID) (commons.ID, bool) {
	switch agentID {
	case negotiation.Agent1:
		return negotiation.Agent2, true
	case negotiation.Agent2:
		return negotiation.Agent1, true
	default:
		panic(fmt.Sprintf("agent %s is not involved in this negotiation", agentID))
	}
}

// Notarize
// A trade is valid iff:
// 1. exactly 2 agents are involved
// 2. both agents are valid
// 3. both agents had the chance to make offer and demand
// 4. both agents' offer are valid if they exist
// if a trade is valid and one of the agent has offered nothing, the trade is considered as a donation
func (negotiation *TradeNegotiation) Notarize(agents map[commons.ID]state.AgentState) (success bool) {
	numberOfValidAgents := 0

	if agent1, ok := agents[negotiation.Agent1]; ok {
		numberOfValidAgents++
		if negotiation.Condition1.Offer.IsValid {
			if !agent1.HasItem(negotiation.Condition1.Offer.ItemType, negotiation.Condition1.Offer.Item.Id()) {
				return false
			}
		}
	}

	if agent2, ok := agents[negotiation.Agent2]; ok {
		numberOfValidAgents++
		if negotiation.Condition2.Offer.IsValid {
			if !agent2.HasItem(negotiation.Condition2.Offer.ItemType, negotiation.Condition2.Offer.Item.Id()) {
				return false
			}
		}
	}

	return numberOfValidAgents == 2 && negotiation.RoundNum >= 1
}

func (negotiation *TradeNegotiation) UpdateOffer(agentID commons.ID, offer TradeOffer) (oldOffer TradeOffer, ok bool) {
	switch agentID {
	case negotiation.Agent1:
		oldOffer = negotiation.Condition1.Offer
		ok = true
		negotiation.Condition1.Offer = offer
	case negotiation.Agent2:
		oldOffer = negotiation.Condition2.Offer
		ok = true
		negotiation.Condition2.Offer = offer
	default:
		return TradeOffer{}, false
	}

	return
}

func (negotiation *TradeNegotiation) UpdateDemand(agentID commons.ID, demand TradeDemand) {
	switch agentID {
	case negotiation.Agent1:
		negotiation.Condition1.Demand = demand
	case negotiation.Agent2:
		negotiation.Condition2.Demand = demand
	default:
	}
}
