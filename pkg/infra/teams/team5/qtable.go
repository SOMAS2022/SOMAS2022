package team5

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

func (qt *Qtable) getMaxFR(fSas []SaPair) float32 {
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

func (qt *Qtable) learn(reward float32, maxFR float32) {
	qt.table[qt.saTaken] += qt.alpha * (reward + qt.gamma*maxFR - qt.table[qt.saTaken])
}

func NewQTable(alpha float32, gamma float32) *Qtable {
	qt := new(Qtable)
	qt.table = make(map[SaPair]float32)
	qt.alpha = alpha
	qt.gamma = gamma
	return qt
}
