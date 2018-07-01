package main

// Global constants
const boardSize = 3      // length/width of the board
const nOldest = 9        // number of oldest states to remember
const printSteps = false // print board and plan at each step
const epsilon = 0.1      // default epsilon (probability to take random action)
const gamma = 0.9        // default alpha (learning rate)
const initialValue = 0.0 // a (non-ending) state's initial value before iteration
const fluctuation = 1e-5 // the amplitude of fluctuation for initialValue
const winReward = 1.0    // reward for winning the game
const drawReward = 0.5   // reward for draw game
const loseReward = 0.0   // reward for losing the game

// main
func main() {

	// create players
	players := createPlayers()

	// create sessions
	createSessions(players)

}
