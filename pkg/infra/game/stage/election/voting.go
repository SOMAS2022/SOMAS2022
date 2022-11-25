package election

import (
	"infra/game/commons"
	"infra/game/decision"
	"infra/logging"
	"math/rand"
)

/*
	Can implement multiple voting strategies in this file.
	Problems arise when agents are given more than 2
	choices.

	C <= 2
	0. Confidence Vote (Yes/No)
	1. Single/Double Choice plurality
	C > 2
	2. Plurality (1st rank choices)
	3. Runoff (ask agents to vote between best 2 options)
	4. Borda Count
	5. Instant Runoff
	6. Approval
	7. Copeland Scoring
*/

func singleChoicePlurality(ballots map[commons.ID]decision.Ballot) (commons.ID, uint) {
	// Count number of votes collected for each candidate
	votes := make(map[commons.ID]uint)
	for _, ballot := range ballots {
		votes[ballot[0]]++
	}

	// Find the candidate(s) with max number of votes
	var maxNumVotes uint
	var winners []commons.ID
	for agentId, numVotes := range votes {
		if numVotes > maxNumVotes {
			maxNumVotes = numVotes
			winners = []commons.ID{agentId}
		} else if numVotes == maxNumVotes {
			winners = append(winners, agentId)
		}
	}

	// Randomly pick one if no valid votes
	if maxNumVotes == 0 {
		candidateIds := make([]commons.ID, len(ballots))
		i := uint(0)
		for id := range ballots {
			candidateIds[i] = id
			i++
		}
		logging.Log(logging.Warn, nil, "No Valid Votes")
		return candidateIds[len(ballots)], 0
	}

	// Randomly choose one if there are more than one winner
	var winner commons.ID
	if len(winners) > 1 {
		logging.Log(
			logging.Info,
			logging.LogField{"winners": winners},
			"Multiple candidates with a winning number of votes",
		)
		randIdx := rand.Intn(len(winners))
		winner = winners[uint(randIdx)]
	} else {
		winner = winners[0]
	}

	return winner, 100 * maxNumVotes / uint(len(ballots))
}
