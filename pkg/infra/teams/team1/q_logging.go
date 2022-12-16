package team1

import (
	"bufio"
	"fmt"
	"infra/game/commons"
	"infra/game/decision"
	"infra/game/state"
	"infra/teams/team1/internal"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

var logPath string = "teams/team1/"

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
	f, err := os.OpenFile(logPath+"log.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
			agentState.TotalAttack(),
			agentState.TotalDefense(),
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
	logFile, err := os.Open(logPath + "log.csv")
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
	logFile, err = os.Open(logPath + "log.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	scanner = bufio.NewScanner(logFile)

	// Open a new file to store post processed logs
	postLogFile, err := os.OpenFile(logPath+"post_log.csv", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	for scanner.Scan() {

		if rand.Float64() < 0.9 {
			continue
		}

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
	os.Remove(logPath + "log.csv")

}

func LearnStrategies() {
	// Read from post_log.csv to get data

	// Check if logging is turned on
	if os.Getenv("QLOGGING") != "true" {
		// Logging is not turned on so return
		return
	}

	// Arrays to keep track of data for each action
	// For each of 3 actions, associated table with rows of experience
	var action_data [3][][]float64

	// Array of observations
	var coop_obs [3][]float64
	var self_obs [3][]float64

	// Open log file
	logFile, err := os.Open(logPath + "post_log.csv")
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

		decision, _ := strconv.Atoi(snapshot[2])
		if len(snapshot) == 13 {
			action_data[decision] = append(action_data[decision], getStateFromSnapshot(snapshot))

			// TODO Should really write these arrays directly instead of writing to csv first
			//Cooperative: Mean round
			coop_obs[decision] = append(coop_obs[decision], internal.StringToFloat(snapshot[12]))
			//Selfish: Agent round
			self_obs[decision] = append(self_obs[decision], internal.StringToFloat(snapshot[11]))
		}
	}

	logFile.Close()

	// Writing to weights.csv

	// If logging file doesn't exist, create it, or append to the file
	weightFile, err := os.OpenFile(logPath+"weights.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// Reduce size of tables

	// Train q function and write weights to csv
	// Write cooperative weights (Defend, Cower, Attack)
	for action := 0; action < 3; action++ {
		coop_w := internal.FitLinReg(action_data[action], coop_obs[action])

		for _, weight := range coop_w {
			weightFile.Write([]byte(fmt.Sprintf("%f,", weight)))
		}

		weightFile.Write([]byte("\n"))
	}

	// Write selfish weights (Defend, Cower, Attack)
	for action := 0; action < 3; action++ {
		self_w := internal.FitLinReg(action_data[action], self_obs[action])

		for _, weight := range self_w {
			weightFile.Write([]byte(fmt.Sprintf("%f,", weight)))
		}

		weightFile.Write([]byte("\n"))
	}

	weightFile.Close()

}

func getStateFromSnapshot(snapshot []string) []float64 {
	// Collect csv string experience values to a float array
	return []float64{
		internal.StringToFloat(snapshot[6]), // Hp
		internal.StringToFloat(snapshot[7]), // Stamina
		internal.StringToFloat(snapshot[8]), // TotalAttack
		internal.StringToFloat(snapshot[9]), // TotalDefense
		rand.Float64(),                      //internal.StringToFloat(snapshot[3]), // CurrLevel
		internal.StringToFloat(snapshot[5]), // MonsterHealth
		internal.StringToFloat(snapshot[4]), // MonsterAttack
	}
}
