package team5

import (
	"infra/game/commons"
	"infra/game/message"
	"infra/game/state"
	"infra/game/stage/trade"
	"infra/game/stage/trade/internal"
)

func HandleTradeNegotiation(a agent.BaseAgent, T_TradeInfo message.TradeInfo, round uint) message.TradeMessage {
	leaderfightpredictionmap := FindBestStrategy(a.view)
	T_fightDecision := leaderfightpredictionmap[a.id]
	if len(T_TradeInfo.Negotiations) == 0 && round < 2 && T_fightDecision!= Cower{ //no offers then request
		return T_request(a, T_fightDecision)
	}else if round < 2{ //process requests 
		return firstResponse(a, T_TradeInfo.Negotiations, T_fightDecision)
	}else if len(T_TradeInfo.Negotiations) != 0{ //process requests
		return secondResponse(a, T_TradeInfo.Negotiations, T_fightDecision)
	}else{ //no offers then abstain
		return message.TradeAbstain{}
	}
}

func T_request(a agent.BaseAgent, T_fightDecision decision.FightAction) message.TradeRequest{
	T_id_attack := getTradePartnerAttack(a, T_fightDecision)
	T_id_defend := getTradePartnerDefend(a, T_fightDecision)
	if(T_fightDecision == decision.Fight){
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

func getTradePartnerAttack(a agent.BaseAgent, T_fightDecision decision.FightAction) commons.ID{
	for id, T_HiddenAgentState := range a.view.agentState{
		if T_fightDecision == decision.Cower && T_HiddenAgentState.BonusAttack > a.latestState.Weapons[0]{
			return id
		}
	}
	return ""
}
func getTradePartnerDefend(a agent.BaseAgent, T_fightDecision decision.FightAction) commons.ID{
	for id, T_HiddenAgentState := range a.view.agentState{
		if T_fightDecision == decision.Cower && T_HiddenAgentState.BonusDefense > a.latestState.Shields[0]{
			return id
		}
	}
	return ""
}
//if fighting ask for stuff, if cowering process a proposal to accept
func firstResponse(a agent.BaseAgent, T_Negotiations TradeInfo.Negotiations, T_fightDecision decision.FightAction) message.TradeMessage{
	if T_fightDecision != decision.Cower{ 
		return T_request(a, T_fightDecision)
	}
	//if i'm cowering accept first offer from attack/defend guy
	for id,T_TradeNegotiation := range T_Negotiations{ // go through request buffer
		if T_TradeNegotiation.RoundNum == 1 &&  T_fightDecision != decision.Cower{
			return message.TradeAccept{
				TradeID:	T_TradeNegotiation.Id}
		}
	}
	return message.TradeAbstain{}
}

func secondResponse(a agent.BaseAgent, T_Negotiations TradeInfo.Negotiations, T_fightDecision decision.FightAction) message.TradeMessage{
	T_fightmode := decision.FightAction //need change
	//if i'm cowering accept first offer from attack/defend guy
	for id,T_TradeNegotiation := range T_Negotiations{ // go through request buffer
		if T_TradeNegotiation.RoundNum == 1 &&  T_fightDecision != decision.Cower{//need change
			return message.TradeAccept{
				TradeID:	T_TradeNegotiation.Id}
		}
	}
	//if everone requesting is cowering then abstain
	return message.TradeAbstain{}
}