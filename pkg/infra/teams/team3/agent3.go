package team3

import (
	"infra/config"
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/logging"
	"math"

	"github.com/benbjohnson/immutable"
)

type AgentThree struct {
	AT                    int
	SH                    int
	uR                    map[commons.ID]int
	uP                    map[commons.ID]int
	uC                    map[commons.ID]int
	utilityScore          map[commons.ID]int
	TSN                   []commons.ID
	contactsLastRound     map[commons.ID]bool
	chairTolerance        int
	proposalTolerance     map[commons.ID]int
	fightDecisionsHistory commons.ImmutableList[decision.ImmutableFightResult]
	reputationMap         map[commons.ID]float64
	w1Map                 map[commons.ID]float64
	w2Map                 map[commons.ID]float64
	pastHPMap             map[commons.ID]int
	pastStaminaMap        map[commons.ID]int
	personality           int
	sanctioned            int
	statsQueue            StatsQueue
	change_init           float64
	alpha                 float64
	Soc_cap               int
	sample_percent        float64
	numAgents             int
}

// Update internal parameters at the end of each stage
func (a *AgentThree) UpdateInternalState(baseAgent agent.BaseAgent, _ *commons.ImmutableList[decision.ImmutableFightResult], votes *immutable.Map[decision.Intent, uint], log chan<- logging.AgentLog) {
	AS := baseAgent.AgentState()
	view := baseAgent.View()
	// Initialise utils
	if view.CurrentLevel() == 1 {
		a.utilityScore = a.InitUtility(baseAgent)
		a.uR = a.InitUtility(baseAgent)
		a.uP = a.InitUtility(baseAgent)
		a.uC = a.InitUtility(baseAgent)
		// initialise stats and add to the queue.
		stat := Stats{1000, 0, 0, 0}
		stat2 := Stats{1000, 0, 0, 0}
		stat3 := Stats{1000, 0, 0, 0}
		a.statsQueue.addStat(stat)
		a.statsQueue.addStat(stat2)
		a.statsQueue.addStat(stat3)
		viewAS := view.AgentState()
		a.numAgents = viewAS.Len()
	}
	// fetch total attack and defence
	a.AT = int(AS.Attack + AS.BonusAttack())
	a.SH = int(AS.Defense + AS.BonusDefense())

	// a.fightDecisionsHistory = *history
	// a.sendGossipMessage(baseAgent)
	// if preLog != postLog {
	// 	fmt.Println("MSG RECEIVED")
	// }

	// update parameters
	a.Reputation(baseAgent)
	a.UpdateTotalUtility(baseAgent)
	a.ResetContacts()
	a.UpdateTSN(baseAgent)

	// if personality enabled, update it
	enablePersonalityUpdate := config.EnvToBool("UPDATE_PERSONALITY", true)
	if enablePersonalityUpdate {
		// update internal personality
		a.UpdatePersonality(baseAgent)
	}

	// a.CalcReputation(baseAgent)

	//fmt.Println(a.SocialCapital(baseAgent))
	//a.SocialCapital(baseAgent)
}

func (a *AgentThree) GetStats() (int, int) {
	return a.personality, a.sanctioned
}

func (a *AgentThree) UpdatePersonality(baseAgent agent.BaseAgent) {
	// calculate the agent average stats before round
	avg_ini := a.statsQueue.averageStats()

	// get agents stats after round
	AS := baseAgent.AgentState()
	AT := AS.Attack + AS.BonusAttack()
	SH := AS.Defense + AS.BonusDefense()
	HP := AS.Hp
	ST := AS.Stamina
	// add to queue and calculate the new average
	stat := Stats{HP, ST, AT, SH}
	a.statsQueue.addStat(stat)
	avg_now := a.statsQueue.averageStats()
	// calculate difference
	change_now := avg_now.Health - avg_ini.Health
	// calculate % change
	P_C := ((change_now - a.change_init) / math.Abs(change_now))
	// scale
	increment := (P_C * a.alpha)

	// keep with max perosnality swing
	if math.IsNaN(increment) {
		increment = 0
	}
	if increment >= 5 {
		increment = 5
	}
	if increment <= -5 {
		increment = -5
	}
	// update personality
	a.personality = a.personality + int(math.Ceil(increment))
	// keep within maxMin personality
	if a.personality >= 100 {
		a.personality = 100
	}
	if a.personality <= 0 {
		a.personality = 0
	}
	// reset initial change to new value.
	a.change_init = change_now
}

func NewAgentThreeNeutral() agent.Strategy {
	dis := config.EnvToUint("COLLECTIVE_PER", 50)
	// dis, _ := strconv.ParseInt(os.Getenv("PASSIVE_PER"), 10, 0)
	return &AgentThree{
		utilityScore:      CreateUtility(),
		uR:                CreateUtility(),
		uP:                CreateUtility(),
		uC:                CreateUtility(),
		chairTolerance:    0,
		proposalTolerance: make(map[commons.ID]int, 0),
		personality:       int(dis),
		reputationMap:     make(map[commons.ID]float64, 0),
		w1Map:             make(map[commons.ID]float64, 0),
		w2Map:             make(map[commons.ID]float64, 0),
		pastHPMap:         make(map[commons.ID]int, 0),
		pastStaminaMap:    make(map[commons.ID]int, 0),
		sanctioned:        0,
		statsQueue:        *makeStatsQueue(3),
		change_init:       0,
		alpha:             5,
		sample_percent:    0.25,
	}
}

func NewAgentThreePassive() agent.Strategy {
	dis := config.EnvToUint("SELFLESS_PER", 75)
	// dis, _ := strconv.ParseInt(os.Getenv("COLLECTIVE_PER"), 10, 0)
	return &AgentThree{
		utilityScore:      CreateUtility(),
		uR:                CreateUtility(),
		uP:                CreateUtility(),
		uC:                CreateUtility(),
		chairTolerance:    0,
		proposalTolerance: make(map[commons.ID]int, 0),
		personality:       int(dis),
		reputationMap:     make(map[commons.ID]float64, 0),
		w1Map:             make(map[commons.ID]float64, 0),
		w2Map:             make(map[commons.ID]float64, 0),
		pastHPMap:         make(map[commons.ID]int, 0),
		pastStaminaMap:    make(map[commons.ID]int, 0),
		sanctioned:        0,
		statsQueue:        *makeStatsQueue(3),
		change_init:       0,
		alpha:             5,
		sample_percent:    0.25,
	}
}
func NewAgentThreeAggressive() agent.Strategy {
	dis := config.EnvToUint("SELFISH_PER", 25)
	// dis, _ := strconv.ParseInt(os.Getenv("SELFISH_PER"), 10, 0)
	return &AgentThree{
		utilityScore:      CreateUtility(),
		uR:                CreateUtility(),
		uP:                CreateUtility(),
		uC:                CreateUtility(),
		chairTolerance:    0,
		proposalTolerance: make(map[commons.ID]int, 0),
		personality:       int(dis),
		reputationMap:     make(map[commons.ID]float64, 0),
		w1Map:             make(map[commons.ID]float64, 0),
		w2Map:             make(map[commons.ID]float64, 0),
		pastHPMap:         make(map[commons.ID]int, 0),
		pastStaminaMap:    make(map[commons.ID]int, 0),
		sanctioned:        0,
		statsQueue:        *makeStatsQueue(3),
		change_init:       0,
		alpha:             5,
		sample_percent:    0.25,
	}
}
