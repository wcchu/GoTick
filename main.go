package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
)

// Global constants
const boardSize = 3               // length/width of the board
var symbols = [2]string{"x", "o"} // player symbols on the board

// write state values of the player to a csv file
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

// main
func main() {

	// create players
	var r2d2, termino, person player
	r2d2.initializeRobot(
		"R2-D2",
		robotSpecs{eps: 0.1, alp: 0.5, mean: 0.0, fluc: 0.01, draw: 1.0},
		false)
	termino.initializeRobot(
		"Terminator",
		robotSpecs{eps: 0.1, alp: 0.5, mean: 0.0, fluc: 0.01, draw: -1.0},
		false)
	person.initializeHuman("Somebody")

	// train the two robots
	runSession(&r2d2, &termino, 10000)
	exportValues(r2d2.mind.values, "robot1_values.csv")
	exportValues(termino.mind.values, "robot2_values.csv")

	// Somebody vs. R2-D2
	runSession(&person, &r2d2, 3)

}
