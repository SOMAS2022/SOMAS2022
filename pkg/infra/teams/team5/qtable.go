package team5

import (
	"fmt"
	"strconv"
)

type SaPair struct {
	// State is always "health-stamina-attackLevel-defenseLevel",
	// where attackLevel, defenseLevel are relative to other observable agents.
	state string
	// Action is either "attack", "defense", "cower".
	action string
}

type Qtable struct {
	table   map[SaPair]float32
	alpha   float32
	gamma   float32
	saTaken SaPair
}

func (qt *Qtable) GetMaxFR(fSas []SaPair) float32 {
	// find maximum future reward
	var firstMark bool = true
	var maxFR float32 = 0
	for _, fSa := range fSas {
		if firstMark || maxFR < qt.table[fSa] {
			firstMark = false
			maxFR = qt.table[fSa]
		}
	}
	return maxFR
}

func (qt *Qtable) GetMaxQAction(state string) string {
	fstCheck := true
	var maxQ float32 = 0.0
	var maxQPair SaPair
	for sa, q := range qt.table {
		if sa.state == state {
			if fstCheck || q > maxQ {
				fstCheck = false
				maxQ = q
				maxQPair = sa
			}
		}
	}
	if fstCheck {
		return "NoSaPairAvailable"
	}
	return maxQPair.action
}

func (qt *Qtable) Learn(reward float32, maxFR float32) {
	qt.table[qt.saTaken] += qt.alpha * (reward + qt.gamma*maxFR - qt.table[qt.saTaken])
}

func (qt *Qtable) Print() {
	strBuild := "\n ___________________________________________\n"
	for qstate, qvalue := range qt.table {
		strBuild += "|--  " + qstate.state + ":" + qstate.action + "  " + strconv.FormatFloat(float64(qvalue), 'f', 4, 32) + "  --|\n"
	}
	fmt.Print(strBuild)
}

func NewQTable(alpha float32, gamma float32) *Qtable {
	qt := new(Qtable)
	qt.table = make(map[SaPair]float32)
	qt.alpha = alpha
	qt.gamma = gamma
	return qt
}
