package main

import (
	"flag"
	"infra/logging"
	"testing"
)

func TestMain(m *testing.M) {

	// define flags
	useJSONFormatter := flag.Bool("j", false, "Whether to output logs in JSON")
	debug := flag.Bool("d", false, "Whether to run in debug mode. If false, only logs with level info or above will be shown")
	flag.Parse()

	logging.InitLogger(*useJSONFormatter, *debug)

	for i := 0; i < 1; i++ {
		agentMap, globalState, gameConfig := initialise()
		gameLoop(globalState, agentMap, gameConfig)
	}

}
