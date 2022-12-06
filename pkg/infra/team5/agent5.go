package team5

import (
	"infra/game/state"
	"infra/team5/Pkg5/allocation"
	"infra/team5/Pkg5/commons5"
)

type returnType struct {
	allocationMessage
}

func team5(s *state.View) {
	var clashSet []allocation.ClashItem
	var winnerSet []allocation.ClashItem
	agentstates := commons5.FetchAllAgents(s)
	winnerSet = allocation.AllocMessageHandler(agentstates, clashSet)
}
