package team5

import (
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/message"
)

func HandleTradeNegotiation(a agent.BaseAgent, T_TradeInfo message.TradeInfo, round uint) message.TradeMessage {
	leaderfightpredictionmap := FindBestStrategy(a.View())
	T_fightDecision, _ := leaderfightpredictionmap.Get(a.ID())
	if len(T_TradeInfo.Negotiations) == 0 && round < 2 && T_fightDecision != decision.Cower { //no offers then request
		return T_request(a, T_fightDecision)
	} else if round < 2 { //process requests
		return firstResponse(a, T_TradeInfo, T_fightDecision)
	} else if len(T_TradeInfo.Negotiations) != 0 { //process requests
		return secondResponse(a, T_TradeInfo, T_fightDecision)
	} else { //no offers then abstain
		return message.TradeAbstain{}
	}
}

func T_request(a agent.BaseAgent, T_fightDecision decision.FightAction) message.TradeRequest {
	T_id_attack := getTradePartnerAttack(a, T_fightDecision)
	T_id_defend := getTradePartnerDefend(a, T_fightDecision)
	var fiveRequest message.TradeRequest
	if T_fightDecision == decision.Attack {
		makeOffer, _ := message.NewTradeOffer(commons.Shield, 0, a.AgentState().Weapons, a.AgentState().Shields)
		makeDemand := message.NewTradeDemand(commons.Weapon, a.AgentState().BonusAttack())
		fiveRequest = message.TradeRequest{
			CounterPartyID: T_id_attack,
			Offer:          makeOffer,
			Demand:         makeDemand,
		}
	} else {
		makeOffer, _ := message.NewTradeOffer(commons.Weapon, 0, a.AgentState().Weapons, a.AgentState().Shields)
		makeDemand := message.NewTradeDemand(commons.Shield, a.AgentState().BonusDefense())
		fiveRequest = message.TradeRequest{
			CounterPartyID: T_id_defend,
			Offer:          makeOffer,
			Demand:         makeDemand,
		}
	}
	return fiveRequest
}

func getTradePartnerAttack(a agent.BaseAgent, T_fightDecision decision.FightAction) commons.ID {
	for id, T_HiddenAgentState := range *a.View().AgentState() {
		if T_fightDecision == decision.Cower && T_HiddenAgentState.BonusAttack > a.AgentState().Weapons.Get(0).value {
			return id
		}
	}
	return ""
}
func getTradePartnerDefend(a agent.BaseAgent, T_fightDecision decision.FightAction) commons.ID {
	for id, T_HiddenAgentState := range a.View().AgentState() {
		if T_fightDecision == decision.Cower && T_HiddenAgentState.BonusDefense > a.AgentState().Shields.Get(0).value {
			return id
		}
	}
	return ""
}

// if fighting ask for stuff, if cowering process a proposal to accept
func firstResponse(a agent.BaseAgent, T_Negotiations message.TradeInfo, T_fightDecision decision.FightAction) message.TradeMessage {
	if T_fightDecision != decision.Cower {
		return T_request(a, T_fightDecision)
	}
	//if i'm cowering accept first offer from attack/defend guy
	for _, T_TradeNegotiation := range T_Negotiations.Negotiations { // go through request buffer
		if T_TradeNegotiation.RoundNum == 1 && T_fightDecision != decision.Cower {
			return message.TradeAccept{
				TradeID: T_TradeNegotiation.Id}
		}
	}
	return message.TradeAbstain{}
}

func secondResponse(a agent.BaseAgent, T_Negotiations message.TradeInfo, T_fightDecision decision.FightAction) message.TradeMessage {
	//if i'm cowering accept first offer from attack/defend guy
	for _, T_TradeNegotiation := range T_Negotiations.Negotiations { // go through request buffer
		if T_TradeNegotiation.RoundNum == 1 && T_fightDecision != decision.Cower { //need change
			return message.TradeAccept{
				TradeID: T_TradeNegotiation.Id}
		}
	}
	//if everone requesting is cowering then abstain
	return message.TradeAbstain{}
}
