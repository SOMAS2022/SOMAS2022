package state

import (
	"infra/game/commons"
	"infra/game/decision"

	"github.com/benbjohnson/immutable"
)

type View struct {
	currentLevel    uint
	hpPool          uint
	monsterHealth   uint
	monsterAttack   uint
	agentState      *immutable.Map[commons.ID, HiddenAgentState]
	currentLeader   commons.ID
	leaderManifesto decision.Manifesto
}

type (
	HealthRange  uint
	StaminaRange uint
)

const (
	LowHealth  uint = 250 // 25% starting HP
	MidHealth  uint = 500
	HighHealth uint = 750 // 75% starting HP
)

const (
	LowStamina  uint = 500 // 25% starting HP
	MidStamina  uint = 1000
	HighStamina uint = 1500 // 75% starting HP
)

type HiddenAgentState struct {
	Hp           HealthRange
	Stamina      StaminaRange
	Attack       uint
	Defense      uint
	BonusAttack  uint
	BonusDefense uint
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

func (v *View) AgentState() immutable.Map[commons.ID, HiddenAgentState] {
	return *v.agentState
}

func (v *View) CurrentLeader() commons.ID {
	return v.currentLeader
}

func (v *View) LeaderManifesto() decision.Manifesto {
	return v.leaderManifesto
}

func (s *State) ToView() View {
	b := immutable.NewMapBuilder[commons.ID, HiddenAgentState](nil)

	for uuid, state := range s.AgentState {
		healthRange := MidHealth

		if state.Hp < LowHealth {
			healthRange = LowHealth
		} else if state.Hp > HighHealth {
			healthRange = HighHealth
		}

		staminaRange := MidStamina

		if state.Stamina < LowStamina {
			staminaRange = LowHealth
		} else if state.Stamina > HighStamina {
			staminaRange = HighHealth
		}

		b.Set(uuid, HiddenAgentState{
			Hp:           HealthRange(healthRange),
			Stamina:      StaminaRange(staminaRange),
			Attack:       state.Attack,
			Defense:      state.Defense,
			BonusAttack:  state.BonusAttack(*s),
			BonusDefense: state.BonusDefense(*s),
		})
	}

	return View{
		currentLevel:    s.CurrentLevel,
		hpPool:          s.HpPool,
		monsterHealth:   s.MonsterHealth,
		monsterAttack:   s.MonsterAttack,
		agentState:      b.Map(),
		currentLeader:   s.CurrentLeader,
		leaderManifesto: s.LeaderManifesto,
	}
}
