package election

import (
	"fmt"
	"math/rand"

	"infra/game/commons"
	"infra/game/decision"
	"infra/logging"
)

/*
	Can implement multiple voting strategies in this file.
	Problems arise when agents are given more than 2
	choices.

	C <= 2
	1. Single/Double Choice plurality
	C > 2
	2. Plurality (1st rank choices)
	3. Runoff (ask agents to vote between best 2 options)
	4. Borda Count
	5. Instant Runoff
	6. Approval
	7. Copeland Scoring
*/

func singleChoicePlurality(ballots []decision.Ballot) (commons.ID, uint) {
	// Count number of votes collected for each candidate
	votes := make(map[commons.ID]uint)

	for _, ballot := range ballots {
		if len(ballot) > 0 {
			votes[ballot[0]]++
		}
	}

	// Find the candidate(s) with max number of votes
	var maxNumVotes uint
	winners := make([]commons.ID, 0)

	for agentID, numVotes := range votes {
		if numVotes > maxNumVotes {
			maxNumVotes = numVotes
			winners = []commons.ID{agentID}
		} else if numVotes == maxNumVotes {
			winners = append(winners, agentID)
		}
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

// Borda count
// 1. ignore empty ballots
// 2. assume points shared if not shown in non-empty ballots
// 3. randomly select one if multiple agents get the max score
func BordaCount(ballots []decision.Ballot, aliveAgentIDs []commons.ID) (winner commons.ID, maxScore float64) {
	N := len(aliveAgentIDs)
	updated := make(map[commons.ID]bool)
	scores := make(map[commons.ID]float64)

	// Initalise updated to false for all agents,
	for _, id := range aliveAgentIDs {
		updated[id] = false
	}

	// Fill scores
	for _, ballot := range ballots {
		// Ignore empty ballot
		if len(ballot) < 1 {
			continue
		}

		// Reset updated to false for all agents
		for _, id := range aliveAgentIDs {
			updated[id] = false
		}

		// Assign score for agents shown in ballot
		K := 0
		for idx, candidateID := range ballot {
			K = idx
			scores[candidateID] += float64(N) - float64(K) + 1
			updated[candidateID] = true
		}

		// Assign score for agents not shown in ballot
		remaining := 0
		for _, isUpdated := range updated {
			if isUpdated {
				remaining++
			}
		}
		logging.Log(
			logging.Debug,
			logging.LogField{"remainig": remaining},
			fmt.Sprintf("%d agents not shown in ballot", remaining),
		)
		sharedScore := (float64(N) - float64(K) + 1) / float64(remaining)
		for candidateID, _ := range updated {
			scores[candidateID] += sharedScore
		}
	}

	// Find max score
	for _, score := range scores {
		if score > maxScore {
			maxScore = score
		}
	}

	// Find the candidate(s) with max score
	winners := make([]commons.ID, 0)
	for agentID, score := range scores {
		if score == maxScore {
			winners = append(winners, agentID)
		}
	}

	// Randomly choose one if there are more than one winner
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

	return winner, maxScore
}
