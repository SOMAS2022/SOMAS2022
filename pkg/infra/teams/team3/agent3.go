package team3

import (
	"infra/config"
	"infra/game/agent"
	"infra/game/commons"
	"infra/game/decision"
	"infra/logging"
	"math"
	"math/rand"

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
	personality           int
	sanctioned            int
	statsQueue            StatsQueue
	change_init           float64
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

		stat := Stats{1000, 0, 0, 0}
		stat2 := Stats{1000, 0, 0, 0}
		stat3 := Stats{1000, 0, 0, 0}
		a.statsQueue.addStat(stat)
		a.statsQueue.addStat(stat2)
		a.statsQueue.addStat(stat3)

		a.change_init = 0
	}

	a.AT = int(AS.Attack + AS.BonusAttack())
	a.SH = int(AS.Defense + AS.BonusDefense())

	// a.fightDecisionsHistory = *history

	// a.sendGossipMessage(baseAgent)

	// if preLog != postLog {
	// 	fmt.Println("MSG RECEIVED")
	// }

	a.UpdateTotalUtility(baseAgent)
	a.ResetContacts()
	a.UpdateTSN(baseAgent)

	//stat := Stats{1.0, 2, 3, 6}
	//stat := Stats{1.0, 2, 3, 6}
	//a.statsQueue.addStat(stat)
	//// fmt.Println("AVG: ", a.statsQueue)
	enable := config.EnvToBool("UPDATE_PERSONALITY", true)
	if enable {
		a.UpdatePersonality(baseAgent)
	}
}

func (a *AgentThree) Sanctioning() int {
	return 50
}

func (a *AgentThree) PruneAgentList(agentMap map[commons.ID]agent.Agent) map[commons.ID]agent.Agent {
	// fmt.Println("Agent 3")
	prunned := make(map[commons.ID]agent.Agent)
	for id, agent := range agentMap {
		// Compare to 50 in order to sanction
		toSanctionOrNot := rand.Intn(100)
		if toSanctionOrNot > a.Sanctioning() {
			prunned[id] = agent
		}
	}
	// fmt.Println(len(agentMap))
	// fmt.Println(len(prunned))
	return prunned
}

func (a *AgentThree) UpdatePersonality(baseAgent agent.BaseAgent) {
	avg_ini := a.statsQueue.averageStats()

	AS := baseAgent.AgentState()

	AT := AS.Attack + AS.BonusAttack()
	SH := AS.Defense + AS.BonusDefense()
	HP := AS.Hp
	ST := AS.Stamina

	stat := Stats{HP, ST, AT, SH}
	a.statsQueue.addStat(stat)
	avg_now := a.statsQueue.averageStats()

	change_now := avg_now.Health - avg_ini.Health
	P_C := ((change_now - a.change_init) / math.Abs(change_now))
	increment := (P_C * 5)

	if math.IsNaN(increment) {
		increment = 0
	}
	if increment >= 5 {
		increment = 5
	}
	if increment <= -5 {
		increment = -5
	}

	a.personality = a.personality + int(math.Ceil(increment))

	if a.personality >= 100 {
		a.personality = 100
	}

	if a.personality <= 0 {
		a.personality = 0
	}

	a.change_init = (avg_now.Health - avg_ini.Health)
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
		sanctioned:        0,
		statsQueue:        *makeStatsQueue(3),
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
		sanctioned:        0,
		statsQueue:        *makeStatsQueue(3),
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
		sanctioned:        0,
		statsQueue:        *makeStatsQueue(3),
	}
}
