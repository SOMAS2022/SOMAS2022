package main

import (
	"github.com/google/uuid"
	"infra/teams/team1"
	"os"
	"testing"
)

func TestFoo(t *testing.T) {
	t.Parallel()
	main()
}

func TestTrainQ(t *testing.T) {

	// Number of games to run
	numGames := 1

	// Delete old logs
	//os.Remove("teams/team1/post_log.csv")

	// Turn on QLogging
	os.Setenv("QLOGGING", "true")

	// Run numGames full game rounds
	for i := 0; i < numGames; i++ {

		// Set a game id for the game
		id := uuid.New()
		os.Setenv("GAMEID", id.String())

		// Run one game
		main()

		// Postprocess the logs from that specific game
		team1.PostprocessLog()
	}

}

func TestMain(m *testing.M) {
	m.Run()
}
