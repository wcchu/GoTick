package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
)

// BoardSize is the length/width of the board
const BoardSize = 3

// arrayEqualsInteger checks whether all elements in an array is equal to a certain integer
func arrayEqualsInteger(array []int, integer int) bool {
	for _, element := range array {
		if element != integer {
			return false
		}
	}
	return true
}

// playGame runs an episode and lets players (if robot) remember what they've learnt
func playGame(p1, p2 player, e environment) {
	var l location
	e.initializeEnvironment()
	pid := -1
	for !e.gameOver {
		// current player takes action
		if pid == -1 {
			l = p1.robotActs(e)
		} else {
			l = p2.robotActs(e)
		}

		// update environment by the action
		e.updateGameStatus(l, pid)

		// update state history
		state := e.getState()
		p1.updateHistory(state)
		p2.updateHistory(state)

		// switch player
		pid = -pid
	}

	p1.robotUpdatesValues(e)
	p2.robotUpdatesValues(e)

	return
}

// exportValues writes state values of the player to a csv file
func exportValues(vs stateValues, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for state, value := range vs {
		row := []string{strconv.FormatInt(state, 10), strconv.FormatFloat(value, 'g', 5, 64)}
		err := writer.Write(row)
		if err != nil {
			log.Fatal("Cannot write to file", err)
		}
	}

	return
}

// train two robots to play
func main() {
	p1 := player{}
	p1.initializeRobot(-1, 0.1, 0.5, 0.5, 0.01)
	p2 := player{}
	p2.initializeRobot(1, 0.1, 0.5, 0.5, 0.01)
	e := environment{}

	numEpisodes := 10000
	for episode := 0; episode < numEpisodes; episode++ {
		log.Printf("episode = %v", episode)
		playGame(p1, p2, e)
	}

	exportValues(p1.intel.values, "p1_values.csv")
	exportValues(p2.intel.values, "p2_values.csv")

}
