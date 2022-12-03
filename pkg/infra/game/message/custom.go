package message

type CustomPayload interface {
	IsCustomPayload()
}

type CustomInfo struct {
	payload CustomPayload
}

func (c CustomInfo) sealedMessage() {
}

func (c CustomInfo) sealedInform() {
}

func (c CustomInfo) sealedCustomInform() {
}

func (c CustomInfo) Payload() CustomPayload {
	return c.payload
}

func NewCustomMessage(p CustomPayload) CustomInfo {
	return CustomInfo{payload: p}
}
