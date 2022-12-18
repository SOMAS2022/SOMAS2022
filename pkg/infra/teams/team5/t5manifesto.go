package team5

import "infra/game/agent"

type T5Manifesto struct {
	LeaderCombo       int
	CurrentTermLength int
	TotalTermLength   int
}

func (m *T5Manifesto) updateLeaderCombo(ba agent.BaseAgent) {
	id := ba.ID()
	view := ba.View()
	leaderID := view.CurrentLeader()
	if id == leaderID {
		m.CurrentTermLength += 1
		if m.CurrentTermLength == m.TotalTermLength {
			m.LeaderCombo += 1
		}
	} else if m.CurrentTermLength < m.TotalTermLength {
		m.LeaderCombo = 1
	}
}
