package message

import (
	"fmt"
	"infra/game/commons"
	"infra/game/state"
	"infra/logging"

	"github.com/benbjohnson/immutable"
	"github.com/google/uuid"
)

type TradeNegotiation struct {
	// terminate negotiation when this number decrement to 0
	id           commons.TradeID
	initiator    commons.ID
	counterParty commons.ID
	roundNum     uint
	conditions   immutable.Map[commons.ID, TradeCondition]
}

func NewTradeOffer(itemType commons.ItemType, idx uint, isEmpty bool, weapon immutable.List[state.Item], shield immutable.List[state.Item]) (offer TradeOffer, ok bool) {
	if isEmpty {
		return TradeOffer{IsEmpty: true}, true
	}
	var inventory immutable.List[state.Item]
	if itemType == commons.Weapon {
		inventory = weapon
	} else if itemType == commons.Shield {
		inventory = shield
	}
	if idx > uint(inventory.Len()) {
		return TradeOffer{}, false
	}
	item := inventory.Get(int(idx))
	return TradeOffer{ItemType: itemType, Item: item, IsEmpty: false}, true
}

func NewTradeDemand(itemType commons.ItemType, minValue uint) TradeDemand {
	return TradeDemand{ItemType: itemType, MinValue: minValue}
}

func NewTradeNegotiation(agentID commons.ID, counterPartyID commons.ID, offer TradeOffer, demand TradeDemand) TradeNegotiation {
	conditions := immutable.NewMap[commons.ID, TradeCondition](nil)
	conditions = conditions.Set(agentID, TradeCondition{
		Offer:  offer,
		Demand: demand,
	})
	return TradeNegotiation{
		id:           uuid.New().String(),
		initiator:    agentID,
		counterParty: counterPartyID,
		roundNum:     0,
		conditions:   *conditions,
	}
}

func (negotiation TradeNegotiation) sealedMessage() {}

func (negotiation TradeNegotiation) IsInvolved(agentID commons.ID) bool {
	_, ok := negotiation.conditions.Get(agentID)
	return ok
}

func (negotiation TradeNegotiation) GetID() commons.TradeID {
	return negotiation.id
}

func (negotiation TradeNegotiation) GetRoundNum() uint {
	return negotiation.roundNum
}

func (negotiation TradeNegotiation) GetOffer(agentID commons.ID) (TradeOffer, bool) {
	condition, ok := negotiation.conditions.Get(agentID)
	if ok {
		return condition.Offer, true
	}
	return TradeOffer{}, false
}

func (negotiation TradeNegotiation) GetDemand(agentID commons.ID) (demand TradeDemand, ok bool) {
	condition, ok := negotiation.conditions.Get(agentID)
	if ok {
		return condition.Demand, true
	}
	return TradeDemand{}, false
}

func (negotiation TradeNegotiation) GetConditions() immutable.Map[commons.ID, TradeCondition] {
	return negotiation.conditions
}

func (negotiation TradeNegotiation) GetCounterParty(agentID commons.ID) (commons.ID, bool) {
	switch negotiation.conditions.Len() {
	case 0, 1:
		return "", false
	case 2:
		itr := negotiation.conditions.Iterator()
		for !itr.Done() {
			id, _, _ := itr.Next()
			logging.Log(logging.Debug, nil, fmt.Sprintf("id: %s, agentID: %s", id, agentID))
			if id != agentID {
				return id, true
			}
		}
		return "", false
	default:
		panic(fmt.Sprintf("negotiation should have at most 2 agents, %d found", negotiation.conditions.Len()))
	}
}

// A trade is valid iff:
// 1. exactly 2 agents are involved
// 2. both agents are valid
// 3. both agents had the chance to make offer and demand
// 4. both agents' offer are valid if they exist
// if a trade is valid and one of the agent has offered nothing, the trade is considered as a donation
func (negotiation *TradeNegotiation) Notarize(agents map[commons.ID]state.AgentState) (success bool) {
	numberOfValidAgents := 0
	itr := negotiation.conditions.Iterator()
	for !itr.Done() {
		id, condition, _ := itr.Next()
		agent, ok := agents[id]
		ok = ok && (id == negotiation.initiator || id == negotiation.counterParty)
		if ok {
			numberOfValidAgents++
			if !agent.HasItem(condition.Offer.ItemType, condition.Offer.Item.Id()) {
				return false
			}
		}
	}
	return numberOfValidAgents == 2 && negotiation.roundNum > 1
}

func (negotiation *TradeNegotiation) UpdateRoundNum() uint {
	negotiation.roundNum = negotiation.roundNum + 1
	return negotiation.roundNum
}

func (negotiation *TradeNegotiation) UpdateOffer(agentID commons.ID, offer TradeOffer) (oldOffer TradeOffer, ok bool) {
	if agentID != negotiation.initiator || agentID != negotiation.counterParty {
		return TradeOffer{}, false
	}
	condition, ok := negotiation.conditions.Get(agentID)
	if ok {
		negotiation.conditions.Set(agentID, TradeCondition{
			Offer:  offer,
			Demand: condition.Demand,
		})
		return condition.Offer, true
	}
	return TradeOffer{}, false
}

func (negotiation *TradeNegotiation) UpdateDemand(agentID commons.ID, demand TradeDemand) (oldDemand TradeDemand, ok bool) {
	if agentID != negotiation.initiator || agentID != negotiation.counterParty {
		return TradeDemand{}, false
	}
	condition, ok := negotiation.conditions.Get(agentID)
	if ok {
		negotiation.conditions.Set(agentID, TradeCondition{
			Offer:  condition.Offer,
			Demand: demand,
		})
		return condition.Demand, true
	}
	return TradeDemand{}, false
}

func (condition TradeCondition) GetItemType() commons.ItemType {
	return condition.Offer.ItemType
}

func (condition TradeCondition) GetItem() state.Item {
	return condition.Offer.Item
}
