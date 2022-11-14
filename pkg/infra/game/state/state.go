package state

import (
	"github.com/benbjohnson/immutable"
	"infra/game/commons"
)

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
	AgentState    map[commons.ID]AgentState
}

type View struct {
	currentLevel  uint
	hpPool        uint
	monsterHealth uint
	monsterAttack uint
	agentState    *immutable.Map[commons.ID, AgentState]
}

func (v *View) CurrentLevel() uint {
	return v.currentLevel
}

func (v *View) HpPool() uint {
	return v.hpPool
}

func (v *View) MonsterHealth() uint {
	return v.monsterHealth
}

func (v *View) MonsterAttack() uint {
	return v.monsterAttack
}

func (v *View) AgentState() *immutable.Map[commons.ID, AgentState] {
	return v.agentState
}

func (s *State) ToView() *View {
	b := immutable.NewMapBuilder[commons.ID, AgentState](nil)
	for uuid, state := range s.AgentState {
		b.Set(uuid, state)
	}

	return &View{
		currentLevel:  s.CurrentLevel,
		hpPool:        s.HpPool,
		monsterHealth: s.MonsterHealth,
		monsterAttack: s.MonsterAttack,
		agentState:    b.Map(),
	}
}
