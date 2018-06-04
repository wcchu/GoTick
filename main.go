package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
)

// Global constants
const boardSize = 3               // length/width of the board
var symbols = [2]string{"x", "o"} // player symbols on the board

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
		if s == "o" {
			s = "x"
			loc = p1.playerActs(env)
		} else {
			s = "o"
			loc = p2.playerActs(env)
		}

		// update environment by the action
		env.updateGameStatus(loc, s)

		// update state history
		state := env.getState(s)
		p1.updateHistory(state)
		p2.updateHistory(state)
	}

	if p1.being == "human" || p2.being == "human" {
		// someone is a human being, let's announce the winner
		log.Print("game over")
		printBoard(&env.board)
		if env.winner != "" { // there's a winner
			if env.winner == p1.symbol {
				log.Printf("%v is the winner", p1.name)
			} else {
				log.Printf("%v is the winner", p2.name)
			}
		} else {
			log.Print("draw")
		}
	}

	// grow some intelligence
	p1.updateValues(env)
	p2.updateValues(env)

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
	r2d2.initializeRobot(
		"R2-D2",
		robotSpecs{eps: 0.1, alp: 0.5, mean: 0.0, fluc: 0.01, draw: 1.0},
		false)
	termino.initializeRobot(
		"Terminator",
		robotSpecs{eps: 0.1, alp: 0.5, mean: 0.0, fluc: 0.01, draw: -1.0},
		false)
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
	log.Printf("r2d2 won %v times", r2d2.wins)
	exportValues(r2d2.mind.values, "robot1_values.csv")
	log.Printf("termino won %v times", termino.wins)
	exportValues(termino.mind.values, "robot2_values.csv")
	fmt.Print("training session ends \n")

	// human plays with r2d2
	var aHuman player
	aHuman.initializeHuman("A Human")
	playGame(&r2d2, &aHuman, env)

}
