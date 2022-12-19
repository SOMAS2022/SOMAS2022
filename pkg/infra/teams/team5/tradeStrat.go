package team5

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
	"math/rand"

	"github.com/benbjohnson/immutable"
)

func (t5 *Agent5) HandleTradeNegotiation(a agent.BaseAgent, T_TradeInfo message.TradeInfo) message.TradeMessage {
	t5.round += 1
	if t5.round > 4 {
		t5.round -= 5
	}
	leaderfightpredictionmap := FindBestStrategy(a.View())
	//T_fightDecision, _ := leaderfightpredictionmap.Get(a.ID())
	T_fightDecision := decision.FightAction(rand.Intn(3))
	if len(T_TradeInfo.Negotiations) == 0 && t5.round <= 2 && T_fightDecision != decision.Cower { //no offers then request
		return T_request(a, T_fightDecision, leaderfightpredictionmap)
	} else if len(T_TradeInfo.Negotiations) != 0 && t5.round <= 2 { //process requests
		return firstResponse(a, T_TradeInfo, T_fightDecision, leaderfightpredictionmap)
	} else if len(T_TradeInfo.Negotiations) != 0 && T_fightDecision == decision.Cower { //process requests
		return secondResponse(a, T_TradeInfo, leaderfightpredictionmap)
	} else { //no offers then abstain
		return message.TradeAbstain{}
	}
}

func T_request(a agent.BaseAgent, T_fightDecision decision.FightAction, T_map *immutable.Map[commons.ID, decision.FightAction]) message.TradeRequest {
	var fiveRequest message.TradeRequest
	if T_fightDecision == decision.Attack {
		shieldList := a.AgentState().Shields
		itr := shieldList.Iterator()
		s_max := uint(0)
		maxIndex := uint(0)
		for !itr.Done() {
			index, s_item := itr.Next()
			if s_item.Value() > s_max {
				maxIndex = uint(index)
				s_max = s_item.Value()
			}
		}
		T_id_attack := getTradePartnerAttack(a, T_map)
		agent := a.AgentState()
		makeOffer, _ := message.NewTradeOffer(commons.Shield, maxIndex, a.AgentState().Weapons, a.AgentState().Shields)
		makeDemand := message.NewTradeDemand(commons.Weapon, agent.BonusAttack())
		fiveRequest = message.TradeRequest{
			CounterPartyID: T_id_attack,
			Offer:          makeOffer,
			Demand:         makeDemand,
		}
	} else {
		weaponList := a.AgentState().Weapons
		itr := weaponList.Iterator()
		w_max := uint(0)
		maxIndex := uint(0)
		for !itr.Done() {
			index, w_item := itr.Next()
			if w_item.Value() > w_max {
				maxIndex = uint(index)
				w_max = w_item.Value()
			}
		}
		T_id_defend := getTradePartnerDefend(a, T_map)
		agent := a.AgentState()
		makeOffer, _ := message.NewTradeOffer(commons.Weapon, maxIndex, a.AgentState().Weapons, a.AgentState().Shields)
		makeDemand := message.NewTradeDemand(commons.Shield, agent.BonusDefense())
		fiveRequest = message.TradeRequest{
			CounterPartyID: T_id_defend,
			Offer:          makeOffer,
			Demand:         makeDemand,
		}
	}
	return fiveRequest
}

func getTradePartnerAttack(a agent.BaseAgent, T_map *immutable.Map[commons.ID, decision.FightAction]) commons.ID {
	agent_view := a.View()
	agent_state := agent_view.AgentState()
	itr := agent_state.Iterator()
	numOfAgents := agent_state.Len()
	ranNum := rand.Intn(numOfAgents - 1)
	for ranNum >= 0 {
		itr.Next()
		ranNum -= 1
	}
	for !itr.Done() {
		id, T_HiddenAgentState, _ := itr.Next()
		agent := a.AgentState()
		parterDecision, _ := T_map.Get(id)
		if parterDecision == decision.Cower && T_HiddenAgentState.BonusAttack > agent.BonusAttack() {
			return id
		}
	}
	return ""
}
func getTradePartnerDefend(a agent.BaseAgent, T_map *immutable.Map[commons.ID, decision.FightAction]) commons.ID {
	agent_view := a.View()
	agent_state := agent_view.AgentState()
	itr := agent_state.Iterator()
	numOfAgents := agent_state.Len()
	ranNum := rand.Intn(numOfAgents - 1)
	for ranNum >= 0 {
		itr.Next()
		ranNum -= 1
	}
	for !itr.Done() {
		id, T_HiddenAgentState, _ := itr.Next()
		agent := a.AgentState()
		parterDecision, _ := T_map.Get(id)
		if parterDecision == decision.Cower && T_HiddenAgentState.BonusDefense > agent.BonusDefense() {
			return id
		}
	}
	return ""
}

// if fighting ask for stuff, if cowering process a proposal to accept
func firstResponse(a agent.BaseAgent, T_Negotiations message.TradeInfo, T_fightDecision decision.FightAction, T_map *immutable.Map[commons.ID, decision.FightAction]) message.TradeMessage {
	if T_fightDecision != decision.Cower { //if fighting request
		return T_request(a, T_fightDecision, T_map)
	}
	//process response/bargain first
	for agent_id, T_TradeNegotiation := range T_Negotiations.Negotiations {
		agentDecision, _ := T_map.Get(agent_id)
		T_condition := T_TradeNegotiation.Condition2.Offer.IsValid
		if T_TradeNegotiation.RoundNum <= 3 && agentDecision != decision.Cower && T_condition {
			return message.TradeAccept{
				TradeID: T_TradeNegotiation.Id}
		}
	}

	//process request
	for agent_id, T_TradeNegotiation := range T_Negotiations.Negotiations {
		agentDecision, _ := T_map.Get(agent_id)
		if T_TradeNegotiation.RoundNum <= 2 && agentDecision != decision.Cower {
			if T_TradeNegotiation.Condition1.Demand.ItemType == commons.Weapon {
				weaponList := a.AgentState().Weapons
				itr := weaponList.Iterator()
				w_max := uint(0)
				maxIndex := uint(0)
				for !itr.Done() {
					index, w_item := itr.Next()
					if w_item.Value() > w_max {
						maxIndex = uint(index)
						w_max = w_item.Value()
					}
				}
				makeOffer, _ := message.NewTradeOffer(commons.Weapon, maxIndex, a.AgentState().Weapons, a.AgentState().Shields)
				return message.TradeBargain{
					TradeID: T_TradeNegotiation.Id,
					Offer:   makeOffer,
					Demand:  message.TradeDemand{}}
			} else {
				shieldList := a.AgentState().Shields
				itr := shieldList.Iterator()
				s_max := uint(0)
				maxIndex := uint(0)
				for !itr.Done() {
					index, s_item := itr.Next()
					if s_item.Value() > s_max {
						maxIndex = uint(index)
						s_max = s_item.Value()
					}
				}
				makeOffer, _ := message.NewTradeOffer(commons.Shield, maxIndex, a.AgentState().Weapons, a.AgentState().Shields)
				return message.TradeBargain{
					TradeID: T_TradeNegotiation.Id,
					Offer:   makeOffer,
					Demand:  message.TradeDemand{}}
			}
		}
	}
	return message.TradeAbstain{}
}

func secondResponse(a agent.BaseAgent, T_Negotiations message.TradeInfo, T_map *immutable.Map[commons.ID, decision.FightAction]) message.TradeMessage {
	//process request
	for agent_id, T_TradeNegotiation := range T_Negotiations.Negotiations {
		agentDecision, _ := T_map.Get(agent_id)
		if T_TradeNegotiation.RoundNum <= 2 && agentDecision != decision.Cower {
			if T_TradeNegotiation.Condition1.Demand.ItemType == commons.Weapon {
				weaponList := a.AgentState().Weapons
				itr := weaponList.Iterator()
				w_max := uint(0)
				maxIndex := uint(0)
				for !itr.Done() {
					index, w_item := itr.Next()
					if w_item.Value() > w_max {
						maxIndex = uint(index)
						w_max = w_item.Value()
					}
				}
				makeOffer, _ := message.NewTradeOffer(commons.Weapon, maxIndex, a.AgentState().Weapons, a.AgentState().Shields)
				return message.TradeBargain{
					TradeID: T_TradeNegotiation.Id,
					Offer:   makeOffer,
					Demand:  message.TradeDemand{}}
			} else {
				shieldList := a.AgentState().Shields
				itr := shieldList.Iterator()
				s_max := uint(0)
				maxIndex := uint(0)
				for !itr.Done() {
					index, s_item := itr.Next()
					if s_item.Value() > s_max {
						maxIndex = uint(index)
						s_max = s_item.Value()
					}
				}
				makeOffer, _ := message.NewTradeOffer(commons.Shield, maxIndex, a.AgentState().Weapons, a.AgentState().Shields)
				return message.TradeBargain{
					TradeID: T_TradeNegotiation.Id,
					Offer:   makeOffer,
					Demand:  message.TradeDemand{}}
			}
		}
	}
	return message.TradeAbstain{}
}
