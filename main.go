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
func playGame(player1, player2 player, env environment) {
	var loc location
	env.initializeEnvironment()
	pid := -1
	for !env.gameOver {
		// current player takes action
		// TODO: consider human player as well
		if pid == -1 {
			loc = player1.robotActs(env)
		} else {
			loc = player2.robotActs(env)
		}

		// update environment by the action
		env.updateGameStatus(loc, pid)

		// update state history
		state := env.getState()
		player1.updateHistory(state)
		player2.updateHistory(state)

		// switch player
		pid = -pid
	}

	player1.robotUpdatesValues(env)
	player2.robotUpdatesValues(env)

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

// train robots and let human play with a robot
func main() {
	var env environment

	// train two robots
	var robot1, robot2 player
	robot1.initializeRobot(-1, 0.1, 0.5, 0.5, 0.01)
	robot2.initializeRobot(1, 0.1, 0.5, 0.5, 0.01)
	numEpisodes := 10000
	for episode := 0; episode < numEpisodes; episode++ {
		log.Printf("episode = %v", episode)
		playGame(robot1, robot2, env)
	}
	exportValues(robot1.intel.values, "robot1_values.csv")
	exportValues(robot2.intel.values, "robot2_values.csv")

}
