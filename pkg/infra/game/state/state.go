package state

type AgentState struct {
	Hp            uint
	Attack        uint
	Defense       uint
	AbilityPoints uint
	BonusAttack   uint
	BonusDefense  uint
}

func (a AgentState) TotalAttack() uint {
	return a.Attack + a.BonusAttack
}

func (a AgentState) TotalDefense() uint {
	return a.Defense + a.BonusDefense
}

type State struct {
	CurrentLevel  uint
	HpPool        uint
	MonsterHealth uint
	MonsterAttack uint
	AgentState    map[uint]AgentState
}
