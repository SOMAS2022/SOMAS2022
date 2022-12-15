package message

import (
	"github.com/benbjohnson/immutable"
)

type ArrayInfo struct {
	num       int
	stringArr immutable.List[string]
}

func NewArrayInfo(num int, strArr []string) ArrayInfo {
	builder := immutable.NewListBuilder[string]()
	for _, str := range strArr {
		builder.Append(str)
	}

	return ArrayInfo{num: num, stringArr: *builder.List()}
}

func (s ArrayInfo) GetStringArr() []string {
	strArr := make([]string, 0, s.stringArr.Len())
	it := s.stringArr.Iterator()
	for !it.Done() {
		_, str := it.Next()
		strArr = append(strArr, str)
	}
	return strArr
}

func (s ArrayInfo) GetNum() int {
	return s.num
}

func (s ArrayInfo) sealedMessage() {
	// TODO implement me
	panic("implement me")
}

func (s ArrayInfo) sealedInform() {
	panic("implement me")
}

func (s ArrayInfo) sealedFightInform() {
	panic("implement me")
}
