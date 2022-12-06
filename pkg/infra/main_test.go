package main

import (
	"flag"
	"github.com/google/uuid"
	"infra/logging"
	"infra/teams/team1"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	// define flags
	useJSONFormatter := flag.Bool("j", false, "Whether to output logs in JSON")
	debug := flag.Bool("d", false, "Whether to run in debug mode. If false, only logs with level info or above will be shown")
	flag.Parse()

	logging.InitLogger(*useJSONFormatter, *debug)

	// Delete old logs
	os.Remove("teams/team1/log.csv")
	os.Remove("teams/team1/post_log.csv")

	for i := 0; i < 1; i++ {
		id := uuid.New()
		os.Setenv("GAMEID", id.String())
		agentMap, globalState, gameConfig := initialise()
		gameLoop(globalState, agentMap, gameConfig)
	}

	team1.PostprocessLog()

}
