package decision

type LootAction int64

const (
	Shield LootAction = iota
	Weapon
	HealthPotion
	StaminaPotion
)
