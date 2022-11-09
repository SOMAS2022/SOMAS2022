package message

type Payload interface {
	isPayload()
}

type Type int64

const (
	Proposal Type = iota
	Something
	SomethingElse
)

type Message struct {
	mType   Type
	payload Payload
}
