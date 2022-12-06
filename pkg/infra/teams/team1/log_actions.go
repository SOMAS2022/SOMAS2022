package team1

import (
	"bufio"
	"fmt"
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

func LogAction(view *state.View, agentID string, action decision.FightAction) {

	// Uses environment variables for gameID
	gameID := os.Getenv("GAMEID")
	if gameID == "" {
		panic("GAMEID variable not set")
	}

	agentState, _ := view.AgentState().Get(agentID)

	message := []byte(fmt.Sprintf("%s,%s,%d,%d,%d,%d,%d,%d,%d,%d,agentRound,meanRound\n",
		gameID,
		agentID,
		action,
		view.CurrentLevel(),
		view.MonsterAttack(),
		view.MonsterHealth(),
		agentState.Hp,
		agentState.Stamina,
		agentState.TotalAttack(),
		agentState.TotalDefense(),
	))

	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile("teams/team1/log.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write([]byte(message)); err != nil {
		f.Close() // ignore error; Write error takes precedence
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

type gameAgent struct {
	game  string
	agent string
}

func PostprocessLog() {
	// Write round each agent reached to log after games are done

	roundTracker := make(map[gameAgent]int)

	file, err := os.Open("teams/team1/log.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

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

	file, err = os.Open("teams/team1/log.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner = bufio.NewScanner(file)

	f, err := os.OpenFile("teams/team1/post_log.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

		if _, err := f.Write([]byte(snapshotString + "\n")); err != nil {
			f.Close() // ignore error; Write error takes precedence
			log.Fatal(err)
		}
	}

	f.Close()

}
