package state

import "github.com/benbjohnson/immutable"

type AgentState struct {
	Hp           uint
	Stamina      uint
	Attack       uint
	Defense      uint
	BonusAttack  uint
	BonusDefense uint
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
	AgentState    map[string]AgentState
}

type View struct {
	CurrentLevel  uint
	HpPool        uint
	MonsterHealth uint
	MonsterAttack uint
	AgentState    *immutable.Map[string, AgentState]
}

func (s *State) ToView() *View {
	b := immutable.NewMapBuilder[string, AgentState](nil)
	for uuid, state := range s.AgentState {
		b.Set(uuid, state)
	}

	return &View{
		CurrentLevel:  s.CurrentLevel,
		HpPool:        s.HpPool,
		MonsterHealth: s.MonsterHealth,
		MonsterAttack: s.MonsterAttack,
		AgentState:    b.Map(),
	}
}
