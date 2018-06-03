package main

import (
	"encoding/csv"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
)

// Global constants
const boardSize = 3               // length/width of the board
var symbols = [2]string{"x", "o"} // player symbols on the board

// rowFilled checks whether all elements in a string array are equal to a certain string
func rowFilled(array []string, s string) bool {
	for _, element := range array {
		if element != s {
			return false
		}
	}
	return true
}

// playGame runs an episode and lets players (if robot) remember what they've learnt
func playGame(p1, p2 *player, env environment) {
	var loc location
	env.initializeEnvironment()
	// p1 always starts first and uses "x"
	p1.symbol = "x"
	p2.symbol = "o"
	s := "o" // current player
	for !env.gameOver {
		// switch player and take action
		// TODO: consider human player as well
		if s == "o" {
			s = "x"
			loc = p1.robotActs(env)
		} else {
			s = "o"
			loc = p2.robotActs(env)
		}

		// update environment by the action
		env.updateGameStatus(loc, s)

		// update state history
		state := env.getState(s)
		p1.updateHistory(state)
		p2.updateHistory(state)
	}

	//log.Print("game over")
	//printBoard(&env.board)

	p1.robotUpdatesValues(env)
	p2.robotUpdatesValues(env)

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
	var r2d2, termino player
	r2d2.initializeRobot("R2-D2", 0.1, 0.5, 0.0, 0.01, false)
	termino.initializeRobot("Terminator", 0.1, 0.5, 0.0, 0.01, false)
	numEpisodes := 10000
	for episode := 0; episode < numEpisodes; episode++ {
		if math.Mod(float64(episode+1), 1000) == 0 {
			log.Printf("episode = %v", episode)
		}
		// for each episode, randomly pick the first player
		if rand.Float64() < 0.5 {
			playGame(&r2d2, &termino, env)
		} else {
			playGame(&termino, &r2d2, env)
		}
	}

	//log.Printf("r2d2 won %v times", r2d2.wins)
	exportValues(r2d2.intel.values, "robot1_values.csv")
	//log.Printf("termino won %v times", termino.wins)
	exportValues(termino.intel.values, "robot2_values.csv")

}
