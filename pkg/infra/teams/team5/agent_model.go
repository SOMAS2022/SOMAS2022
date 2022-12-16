package team5

import (
	"infra/game/commons"
	"infra/game/decision"

	"github.com/benbjohnson/immutable"
)

type Agent struct {
	ID      string
	Hp      uint
	Attack  uint
	Defense uint
	Action  uint
}

type Result struct {
	Death  int      // 死亡数
	Damage uint     // 伤害数
	Agents []*Agent // 组合
}

// 转换返回结果
func ConvertToImmutable(agents []*Agent, agentsAll []*Agent) *immutable.Map[commons.ID, decision.FightAction] {
	b := immutable.NewMapBuilder[commons.ID, decision.FightAction](nil)
	// 已经做出攻击/防御决定的人
	agentMap := make(map[string]struct{})
	for _, item := range agents {
		agentMap[item.ID] = struct{}{}
		b.Set(item.ID, decision.FightAction(item.Action))
	}
	// 遍历所有人，找出没做决定的人并逃跑
	for _, item := range agentsAll {
		if _, ok := agentMap[item.ID]; !ok {
			b.Set(item.ID, decision.FightAction(decision.Cower))
		}
	}
	return b.Map()
}
