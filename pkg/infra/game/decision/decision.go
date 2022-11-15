package decision

type LootDecision struct{}

type HPPoolDecision struct{}

type FightAction int64

const (
	Attack FightAction = iota
	Defend
	Cower
)

func CowerPtr() (p *FightAction) {
	p = new(FightAction)
	*p = Cower
	return
}

func AttackPtr() (p *FightAction) {
	p = new(FightAction)
	*p = Attack
	return
}

func DefendPtr() (p *FightAction) {
	p = new(FightAction)
	*p = Defend
	return
}
