package main

// Global constants
const boardSize = 3               // length/width of the board
const nOldest = 9                 // number of oldest states to remember
const printSteps = false          // print board and plan at each step
const fluctuation = 1e-5          // random fluctuation of default state value
var symbols = [2]string{"x", "o"} // player symbols on the board

// main
func main() {

	// create players
	players := createPlayers()

	// create sessions
	createSessions(players)

}
