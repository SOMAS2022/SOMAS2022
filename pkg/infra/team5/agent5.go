package team5

import (
	"infra/game/message"
	"infra/game/state"
	"infra/team5/Pkg5/allocation"
	"infra/team5/Pkg5/leaderFight"
)

type ReturnType struct {
	allocationMessage  *message.TaggedMessage
	leaderFightMessage *message.TaggedMessage
}

func team5(allocM *message.TaggedMessage, s *state.View) ReturnType {
	var returnType ReturnType

	returnType.allocationMessage = allocation.AllocMessageHandler(allocM, s)
	returnType.leaderFightMessage = leaderFight.LeaderFightMessageHandler(s)

	return returnType
}
