package main

import (
	"flag"
	"infra/logging"
	"testing"
)

func TestMain(m *testing.M) {

	// define flags
	useJSONFormatter := flag.Bool("j", false, "whether to use JSONFormatter for logging")
	flag.Parse()

	logging.InitLogger(*useJSONFormatter)

	for i := 0; i < 1; i++ {
		agentMap, globalState, gameConfig := initialise()
		gameLoop(globalState, agentMap, gameConfig)
	}

}
