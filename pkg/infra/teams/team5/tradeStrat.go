package team5

import (
	"infra/game/commons"
	"infra/game/message"
	"infra/game/state"
	"infra/game/stage/trade"
	"infra/game/stage/trade/internal"
)

func HandleTradeNegotiation(a agent.BaseAgent, T_TradeInfo message.TradeInfo, round uint) message.TradeMessage {

	if len(T_TradeInfo.Negotiations) == 0 && round < 2 && leaderfightpredictionmap[a.id]!= Cower{ //need change, no offers then request
		return T_request()
	}else if round < 2{ //process requests 
		return firstResponse(a, T_TradeInfo.Negotiations)
	}else if len(T_TradeInfo.Negotiations) != 0{ //process requests
		return secondResponse(a, T_TradeInfo.Negotiations)
	}else{ //no offers then abstain
		return message.TradeAbstain{}
	}
}

func T_request(a agent.BaseAgent) message.TradeRequest{
	T_id_attack := getTradePartnerAttack(a)
	T_id_defend := getTradePartnerDefend(a)
	if(leaderfightpredictionmap[id] == decision.Fight){
		makeOffer := message.NewTradeOffer(Weapon, a.latestState.Weapons[0], weapon immutable.List[state.Item], shield immutable.List[state.Item])
		makeDemand := message.NewTradeDemand(Weapon, minValue uint)
	}else {
		makeOffer := message.NewTradeOffer(Shield, idx uint, weapon immutable.List[state.Item], shield immutable.List[state.Item])
		makeDemand := message.NewTradeDemand(Shield, minValue uint)
	}
		fiveRequest:= TradeRequest{
			CounterPartyID: T_id,
			Offer:		    makeOffer,
			Demand:			makeDemand,	
		}
		return fiveRequest
}

func getTradePartnerAttack(a agent.BaseAgent) commons.ID{
	for id, T_HiddenAgentState := range a.view.agentState{
		if leaderfightpredictionmap[id] == decision.Cower && T_HiddenAgentState.BonusAttack > a.latestState.Weapons[0]{
			return id
		}
	}
	return ""
}
func getTradePartnerDefend(a agent.BaseAgent) commons.ID{
	for id, T_HiddenAgentState := range a.view.agentState{
		if leaderfightpredictionmap[id] == decision.Cower && T_HiddenAgentState.BonusDefense > a.latestState.Shields[0]{
			return id
		}
	}
	return ""
}
//if fighting ask for stuff, if cowering process a proposal to accept
func firstResponse(a agent.BaseAgent, T_Negotiations TradeInfo.Negotiations) message.TradeMessage{
	T_fightmode := leaderfightpredictionmap[a.Id] //need change
	if T_fightmode != decision.Cower{ //need change
		return T_request(a)
	}
	//if i'm cowering accept first offer from attack/defend guy
	for id,T_TradeNegotiation := range T_Negotiations{ // go through request buffer
		if T_TradeNegotiation.RoundNum == 1 &&  leaderfightpredictionmap[id] != decision.Cower{//need change
			return message.TradeAccept{
				TradeID:	T_TradeNegotiation.Id}
		}
	}
	return message.TradeAbstain{}
}

func secondResponse(a agent.BaseAgent, T_Negotiations TradeInfo.Negotiations) message.TradeMessage{
	T_fightmode := decision.FightAction //need change
	//if i'm cowering accept first offer from attack/defend guy
	for id,T_TradeNegotiation := range T_Negotiations{ // go through request buffer
		if T_TradeNegotiation.RoundNum == 1 &&  leaderfightpredictionmap[id] != decision.Cower{//need change
			return message.TradeAccept{
				TradeID:	T_TradeNegotiation.Id}
		}
	}
	//if everone requesting is cowering then abstain
	return message.TradeAbstain{}
}