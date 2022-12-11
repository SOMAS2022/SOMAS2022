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

type AndCondition struct {
	condA Condition
	condB Condition
}

func (a AndCondition) CondA() Condition {
	return a.condA
}

func (a AndCondition) CondB() Condition {
	return a.condB
}

func NewAndCondition(condA Condition, condB Condition) *AndCondition {
	return &AndCondition{condA: condA, condB: condB}
}

func (a AndCondition) sealedCondition() {
}

type OrCondition struct {
	condA Condition
	condB Condition
}

func (a OrCondition) CondA() Condition {
	return a.condA
}

func (a OrCondition) CondB() Condition {
	return a.condB
}

func NewOrCondition(condA Condition, condB Condition) *OrCondition {
	return &OrCondition{condA: condA, condB: condB}
}

func (a OrCondition) sealedCondition() {
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
