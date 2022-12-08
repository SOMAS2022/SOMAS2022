package team1

import (
	"bufio"
	"fmt"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/state"
	"log"
	"os"
	"strconv"
	"strings"
)

func meanAbove(levelSlice []int, level int) float64 {

	var levelsAbove []int

	for _, value := range levelSlice {
		if value >= level {
			levelsAbove = append(levelsAbove, value)
		}
	}

	sum := 0

	for _, value := range levelsAbove {
		sum += value
	}

	return float64(sum) / float64(len(levelsAbove))

}

func LogDecisions(decisions map[commons.ID]decision.FightAction, gs state.State) {

	//fmt.Println(os.Getenv("QLOGGING"), "hello")

	// Check if logging is turned on
	if os.Getenv("QLOGGING") != "true" {
		// Logging is not turned on so return
		return
	}

	// Get gameID from environment variable
	gameID := os.Getenv("GAMEID")
	if gameID == "" {
		panic("GAMEID variable not set")
	}

	// If logging file doesn't exist, create it, or append to the file
	f, err := os.OpenFile("pkg/infra/teams/team1/log.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// Loop through the list of decisions
	for agentID, decision := range decisions {

		// Get state of specific agent
		agentState := gs.AgentState[agentID]

		// Calculate the message to log
		message := []byte(fmt.Sprintf("%s,%s,%d,%d,%d,%d,%d,%d,%d,%d,%d,agentRound,meanRound\n",
			gameID,
			agentID,
			decision,
			gs.CurrentLevel,
			gs.MonsterAttack,
			gs.MonsterHealth,
			agentState.Hp,
			agentState.Stamina,
			agentState.TotalAttack(gs),
			agentState.TotalDefense(gs),
			agentState.Weapons.Len(),
		))

		if _, err := f.Write([]byte(message)); err != nil {
			log.Fatal(err)
		}
	}

	f.Close()
}

type gameAgent struct {
	game  string
	agent string
}

func PostprocessLog() {
	// Fill in missing values in log.csv logFile and write to post_log.csv

	// Check if logging is turned on
	if os.Getenv("QLOGGING") != "true" {
		// Logging is not turned on so return
		return
	}

	// Map to keep track of how far agents made it
	roundTracker := make(map[gameAgent]int)

	// Open log log file
	logFile, err := os.Open("pkg/infra/teams/team1/log.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	// Create a new scanner on logFile
	scanner := bufio.NewScanner(logFile)

	// Scan through log to find out when each agent died
	for scanner.Scan() {
		snapshotString := scanner.Text()
		snapshot := strings.Split(snapshotString, ",")
		key := gameAgent{snapshot[0], snapshot[1]}
		if val, ok := roundTracker[key]; ok {
			level, _ := strconv.Atoi(snapshot[3])

			if level > val {
				roundTracker[key] = level
			}
		} else {
			roundTracker[key], _ = strconv.Atoi(snapshot[3])
		}
	}

	// Calculate the mean number of rounds for remaining agents at each level
	meanRoundTracker := make(map[string][]int)
	for key, value := range roundTracker {
		meanRoundTracker[key.game] = append(meanRoundTracker[key.game], value)
	}

	logFile.Close()
	// Open log log file
	logFile, err = os.Open("pkg/infra/teams/team1/log.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	scanner = bufio.NewScanner(logFile)

	// Open a new file to store post processed logs
	postLogFile, err := os.OpenFile("pkg/infra/teams/team1/post_log.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	for scanner.Scan() {
		snapshotString := scanner.Text()
		snapshot := strings.Split(snapshotString, ",")
		key := gameAgent{snapshot[0], snapshot[1]}

		currentLevel, _ := strconv.Atoi(snapshot[3])
		meanAboveString := fmt.Sprintf("%f", meanAbove(meanRoundTracker[key.game], currentLevel))
		snapshotString = strings.Replace(snapshotString, "meanRound", meanAboveString, 1)

		snapshotString = strings.Replace(snapshotString, "agentRound", strconv.Itoa(roundTracker[key]), 1)

		if _, err := postLogFile.Write([]byte(snapshotString + "\n")); err != nil {
			postLogFile.Close() // ignore error; Write error takes precedence
			log.Fatal(err)
		}
	}

	// Close log files
	logFile.Close()
	postLogFile.Close()

	// Delete unprocessed log-logFile
	os.Remove("pkg/infra/teams/team1/log.csv")

}
