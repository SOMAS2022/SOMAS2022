package agent

type HPPool interface {
	DonateToHpPool(baseAgent BaseAgent) uint
}
