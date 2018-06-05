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

	//var r2d2, termino, person player
	//r2d2.initializeRobot("R2-D2", robotSpecs{eps: 0.1, alp: 0.8, mean: 0.5, fluc: 0.2, draw: 1.0}, false)
	//termino.initializeRobot("Terminator", robotSpecs{eps: 0.1, alp: 0.5, mean: 0.5, fluc: 0.2, draw: 0.0}, false)
	//person.initializeHuman("Somebody")

	// train the two robots
	//runSession(&r2d2, &termino, 50000)

	// Somebody vs. R2-D2
	//r2d2.mind.verb = true
	//runSession(&person, &r2d2, 3)

}
