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

var StartingHealthPoints = 1000
var StartingStamina = 2000

// var StartingHealthPoints, _ = strconv.ParseInt(os.Getenv("STARTING_HP"), 10, 0)
// var StartingStamina, _ = strconv.ParseInt(os.Getenv("BASE_STAMINA"), 10, 0)

var granularity = 10

// var granularity, _ = strconv.ParseInt(os.Getenv("GRAN"), 10, 0)

var HealthQuant = int(StartingHealthPoints) / granularity
var StaminaQuant = int(StartingStamina) / granularity

type HiddenAgentState struct {
	Hp           HealthRange
	Stamina      StaminaRange
	Attack       uint
	Defense      uint
	BonusAttack  uint
	BonusDefense uint
	Defector     Defector
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
		healthRange := (state.Hp / uint(HealthQuant)) * uint(HealthQuant)

		staminaRange := (state.Stamina / uint(StaminaQuant)) * uint(StaminaQuant)

		b.Set(uuid, HiddenAgentState{
			Hp:           HealthRange(healthRange),
			Stamina:      StaminaRange(staminaRange),
			Attack:       state.Attack,
			Defense:      state.Defense,
			BonusAttack:  state.BonusAttack(),
			BonusDefense: state.BonusDefense(),
			Defector:     state.Defector,
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
