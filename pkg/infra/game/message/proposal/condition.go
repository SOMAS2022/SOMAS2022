package proposal

type Attribute uint

const (
	Health Attribute = iota
	Stamina
	TotalAttack
	TotalDefence
)

type Comparator uint

const (
	GreaterThan Comparator = iota
	LessThan
)

type Value = uint

type Condition interface {
	sealedCondition()
}

type ComparativeCondition struct {
	Attribute
	Comparator
	Value
}

func NewComparativeCondition(attribute Attribute, comparator Comparator, value Value) *ComparativeCondition {
	return &ComparativeCondition{Attribute: attribute, Comparator: comparator, Value: value}
}

func (c ComparativeCondition) sealedCondition() {
}

// DefectorCondition todo: add functionality for this later
type DefectorCondition struct {
}

func NewDefectorCondition() *DefectorCondition {
	return &DefectorCondition{}
}

func (d DefectorCondition) sealedCondition() {
}
