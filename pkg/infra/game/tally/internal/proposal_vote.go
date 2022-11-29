package internal

import "infra/game/commons"

type VoteCount struct {
	commons.ProposalID
	uint
}
