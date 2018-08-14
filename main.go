package main

import (
	"math/rand"
	"time"
)

// Global constants
const boardSize = 3       // length/width of the board
const nOldest = 9         // number of oldest states to remember
const printSteps = false  // print board and plan at each step
const alpha = 0.5         // default alpha (learning rate)
const epsilon = 0.1       // default epsilon (probability to take random action)
const gamma = 0.9         // default gamma (decay rate)
const initialValue = 0.0  // a (non-ending) state's initial value before iteration
const fluctuation = 0.01  // the amplitude of fluctuation for initialValue
const winReward = 1.0     // reward for winning the game
const drawReward = 0.5    // reward for draw game
const loseReward = 0.0    // reward for losing the game
const nPrintHistory = 100 // print value history every N points

// main
func main() {
	// set random seed to time
	rand.Seed(time.Now().UTC().UnixNano())

	// create players
	players := createPlayers()

	// create sessions
	createSessions(players)

}
