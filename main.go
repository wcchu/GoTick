package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
)

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
func playGame(p1, p2 agent, e environment) {
	var l location
	e.initializeEnvironment()
	pid := -1
	for !e.gameOver {
		// current player takes action
		if pid == -1 {
			l = p1.actAgent(e)
		} else {
			l = p2.actAgent(e)
		}

		// update environment by the action
		e.updateGameStatus(l, pid)

		// update state history
		state := e.getState()
		p1.updateStateHistory(state)
		p2.updateStateHistory(state)

		// switch player
		pid = -pid
	}

	p1.updateValues(e)
	p2.updateValues(e)

	return
}

// exportValues writes state values of the agent to a csv file
func exportValues(vs values, filename string) {
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
	p1 := agent{}
	p1.initializeAgent(-1)
	p2 := agent{}
	p2.initializeAgent(1)
	e := environment{}

	numEpisodes := 10000
	for episode := 0; episode < numEpisodes; episode++ {
		log.Printf("episode = %v", episode)
		playGame(p1, p2, e)
	}

	exportValues(p1.values, "p1_values.csv")

	// let's see what p1 have learnt
	//valueArray := rankStateValues(p1.values)
	//for i := 0; i < 5; i++ {
	//	log.Printf("best state %v, state = %v, value = %v", i, valueArray[i].state, valueArray[i].value)
	//}

}
