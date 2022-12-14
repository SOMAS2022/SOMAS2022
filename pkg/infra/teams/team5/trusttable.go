package team5

import (
	"infra/game/commons"
)

type TrustTable struct {
	indivlTable   map[commons.ID]float32
	leaderTable   map[commons.ID]float32
	baseDecayRate float32
	leadWeight    float32
	indivlWeight  float32
}

func (tt *TrustTable) Decay() {
	for k, v := range tt.indivlTable {
		tt.indivlTable[k] = v * tt.baseDecayRate
	}
	for k, v := range tt.leaderTable {
		tt.leaderTable[k] = v * (tt.baseDecayRate + 0.1)
	}
}

func (tt *TrustTable) PositiveIndivlEvent(agID commons.ID, trustGain float32) {
	_, exist := tt.indivlTable[agID]
	if exist {
		tt.indivlTable[agID] += trustGain
	} else {
		tt.indivlTable[agID] = trustGain
	}

}

func (tt *TrustTable) NegativeIndivlEvent(agID commons.ID, trustLoss float32) {
	_, exist := tt.indivlTable[agID]
	if exist {
		tt.indivlTable[agID] -= trustLoss
	} else {
		tt.indivlTable[agID] = -trustLoss
	}

}

func (tt *TrustTable) PositiveLeaderEvent(agID commons.ID, trustGain float32) {
	_, exist := tt.leaderTable[agID]
	if exist {
		tt.leaderTable[agID] += trustGain
	} else {
		tt.leaderTable[agID] = trustGain
	}
}

func (tt *TrustTable) NegativeLeaderEvent(agID commons.ID, trustLoss float32) {
	_, exist := tt.leaderTable[agID]
	if exist {
		tt.leaderTable[agID] -= trustLoss
	} else {
		tt.leaderTable[agID] = -trustLoss
	}
}

func (tt *TrustTable) EstimateLeadTrust(agID commons.ID) {
	tt.leaderTable[agID] = tt.leadWeight*tt.leaderTable[agID] + tt.indivlWeight*tt.indivlTable[agID]
}

func NewTrustTable() *TrustTable {
	tt := new(TrustTable)
	tt.indivlTable = make(map[commons.ID]float32)
	tt.leaderTable = make(map[commons.ID]float32)
	tt.baseDecayRate = 0.85
	tt.leadWeight = 0.75
	tt.indivlWeight = 0.25
	return tt
}
