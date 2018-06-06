package main

// Global constants
const boardSize = 3               // length/width of the board
var symbols = [2]string{"x", "o"} // player symbols on the board

// main
func main() {

	// create players
	players, _ := createPlayers() // TODO: report error

	// create sessions
	_ = createSessions(players) // TODO: report error

}
