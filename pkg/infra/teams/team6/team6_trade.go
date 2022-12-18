package team6

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/message"
	"infra/game/state"
	"math/rand"

	"github.com/benbjohnson/immutable"
)

func (a *Team6Agent) HandleTradeNegotiation(agent agent.BaseAgent, trades message.TradeInfo) message.TradeMessage {
	if len(trades.Negotiations) == 0 {
		return newTrade(agent)
	} else {
		return respondToTrades(agent, trades)
	}
}

func respondToTrades(agent agent.BaseAgent, trades message.TradeInfo) message.TradeMessage {
	agentState := agent.AgentState()
	attack := agentState.BonusAttack()
	defense := agentState.BonusDefense()

	for id, negotiation := range trades.Negotiations {
		counterparty, _ := negotiation.GetCounterParty(agent.ID())
		offer, _ := negotiation.GetOffer(counterparty)
		demand, _ := negotiation.GetDemand(counterparty)

		reject := false

		if negotiation.Agent1 == agent.ID() {
			//Offer sent by us, check if Agent 2 has responded
			if offer.Item.Id() == "" {
				continue
			} else {
				ourDemand, _ := negotiation.GetDemand(agent.ID())
				if offer.ItemType == ourDemand.ItemType && offer.Item.Value() >= ourDemand.MinValue {
					return message.TradeAccept{TradeID: id}
				} else {
					reject = true
				}

			}
		} else {
			//Offer sent by another agent
			if (demand.ItemType == commons.Shield && demand.MinValue > defense) || (demand.ItemType == commons.Weapon && demand.MinValue > attack) {
				// Can't meet demands
				reject = true
			} else {
				// Can fulfill request if we wish
				if offer.ItemType == demand.ItemType {
					if offer.Item.Value() <= demand.MinValue {
						reject = true
					} else {
						idx := worstItemToTrade(agent, demand)
						if idx == -1 {
							reject = true
						} else {
							ourOffer, _ := message.NewTradeOffer(demand.ItemType, uint(idx), agentState.Weapons, agentState.Shields)
							ourDemand := message.NewTradeDemand(offer.ItemType, offer.Item.Value())
							return message.TradeBargain{TradeID: id, Offer: ourOffer, Demand: ourDemand}
						}
					}
				} else {
					// Offering different item type to demand
					weapons := agentState.Weapons
					shields := agentState.Shields

					var offerGain uint
					var demandLoss uint

					if offer.ItemType == commons.Weapon {
						offerGain = offer.Item.Value() - agentState.BonusAttack()
					} else {
						offerGain = offer.Item.Value() - agentState.BonusDefense()
					}

					itemToGiveIdx := worstItemToTrade(agent, demand)
					demandLoss = diffBestSecondItems(agent, demand.ItemType)

					if offerGain > demandLoss {
						// We will participate
						ourOffer, _ := message.NewTradeOffer(demand.ItemType, uint(itemToGiveIdx), weapons, shields)
						ourDemand := message.NewTradeDemand(offer.ItemType, offer.Item.Value())
						return message.TradeBargain{TradeID: id, Offer: ourOffer, Demand: ourDemand}
					} else {
						reject = true
					}
				}
			}
		}
		// Would rather accept offer B than reject A in a given round
		if reject && len(trades.Negotiations) == 1 {
			return message.TradeReject{TradeID: id}
		}

	}
	return message.TradeAbstain{}
}

func diffBestSecondItems(agent agent.BaseAgent, itemType commons.ItemType) uint {
	agentState := agent.AgentState()
	if itemType == commons.Weapon {
		best := agentState.BonusAttack()
		if agentState.Weapons.Len() == 1 {
			return best
		} else {
			_, itm := secondBestItems(agent, commons.Weapon)
			return (best - itm.Value())
		}
	} else {
		best := agentState.BonusDefense()
		if agentState.Shields.Len() == 1 {
			return best
		} else {
			_, itm := secondBestItems(agent, commons.Weapon)
			return (best - itm.Value())
		}
	}
}

func worstItemToTrade(agent agent.BaseAgent, demand message.TradeDemand) int {
	var items immutable.List[state.Item]
	if demand.ItemType == commons.Shield {
		items = agent.AgentState().Shields
	} else {
		items = agent.AgentState().Weapons
	}

	lowestVal := uint(10000) // Just a big number
	lowestValIdx := -1
	it := items.Iterator()
	for !it.Done() {
		idx, itm := it.Next()
		if itm.Value() < uint(lowestVal) && itm.Value() >= demand.MinValue {
			lowestVal = itm.Value()
			lowestValIdx = idx
		}
	}

	return lowestValIdx
}

func newTrade(agent agent.BaseAgent) message.TradeMessage {
	agentState := agent.AgentState()
	attack := agentState.BonusAttack()
	defense := agentState.BonusDefense()

	if attack == 0 && defense == 0 {
		return message.TradeAbstain{}
	}

	wantWeapon := attack <= defense

	view := agent.View()
	otherAgents := view.AgentState()
	it := otherAgents.Iterator()

	var canTradeIDs []commons.ID

	for !it.Done() {
		agentID, hiddenState, _ := it.Next()
		var val uint
		if wantWeapon {
			val = hiddenState.BonusAttack
		} else {
			val = hiddenState.BonusDefense
		}
		if val-4 > Min(attack, defense) {
			canTradeIDs = append(canTradeIDs, agentID)
		}
	}

	if len(canTradeIDs) == 0 {
		return message.TradeAbstain{}
	} else {
		idx := rand.Intn(len(canTradeIDs))

		var demand message.TradeDemand
		if wantWeapon {
			value := attack + 5
			if attack == 0 {
				value = 0
			}
			demand = message.NewTradeDemand(commons.Weapon, value)
		} else {
			value := defense + 5
			if defense == 0 {
				value = 0
			}
			demand = message.NewTradeDemand(commons.Shield, value)
		}
		return message.TradeRequest{CounterPartyID: canTradeIDs[idx], Offer: tradeOffer(agent, wantWeapon), Demand: demand}
	}
}

func tradeOffer(agent agent.BaseAgent, wantWeapon bool) message.TradeOffer {
	weapons := agent.AgentState().Weapons
	shields := agent.AgentState().Shields

	var offer message.TradeOffer
	itemType, idx := secondBestItem(agent, wantWeapon)
	offer, _ = message.NewTradeOffer(itemType, idx, weapons, shields)
	return offer
}

func secondBestItem(agent agent.BaseAgent, wantWeapon bool) (commons.ItemType, uint) {
	agentState := agent.AgentState()
	weapons := agentState.Weapons
	shields := agentState.Shields

	bestWeapon := bestItemIndex(agent, commons.Weapon)
	bestShield := bestItemIndex(agent, commons.Shield)

	if weapons.Len() == 0 {
		return commons.Shield, bestShield
	} else if shields.Len() == 0 {
		return commons.Weapon, bestWeapon
	}

	weaponIdx, weapon := secondBestItems(agent, commons.Weapon)
	shieldIdx, shield := secondBestItems(agent, commons.Shield)

	if wantWeapon {
		if shield.Value() > agentState.TotalAttack() {
			return commons.Shield, uint(shieldIdx)
		} else {
			return commons.Weapon, bestWeapon
		}
	} else {
		if weapon.Value() > agentState.TotalDefense() {
			return commons.Weapon, uint(weaponIdx)
		} else {
			return commons.Shield, bestShield
		}
	}
}

func bestItemIndex(agent agent.BaseAgent, itemType commons.ItemType) uint {
	agentState := agent.AgentState()
	var items immutable.List[state.Item]
	var id commons.ItemID

	if itemType == commons.Shield {
		items = agentState.Shields
		id = agentState.ShieldInUse
	} else {
		items = agentState.Weapons
		id = agentState.WeaponInUse
	}

	it := items.Iterator()

	for !it.Done() {
		idx, itm := it.Next()
		if itm.Id() == id {
			return uint(idx)
		}
	}
	return 0
}

func secondBestItems(agent agent.BaseAgent, itemType commons.ItemType) (idx int, item state.Item) {
	var items immutable.List[state.Item]
	var bestID commons.ItemID
	if itemType == commons.Shield {
		items = agent.AgentState().Shields
		bestID = agent.AgentState().ShieldInUse
	} else {
		items = agent.AgentState().Weapons
		bestID = agent.AgentState().WeaponInUse
	}

	var highestItem state.Item
	highestItemIdx := -1

	it := items.Iterator()

	for !it.Done() {
		idx, itm := it.Next()
		if itm.Value() > highestItem.Value() && itm.Id() != bestID {
			highestItem = itm
			highestItemIdx = idx
		}
	}

	return highestItemIdx, highestItem
}
