package main

import (
	"log"
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

	log.Printf("*** game over ***; winner is %v", e.winner)
	e.printBoard()

	p1.updateValues(e)
	p2.updateValues(e)

	return
}

// train two robots to play
func main() {
	p1 := agent{}
	p1.initializeAgent(-1)
	p2 := agent{}
	p2.initializeAgent(1)
	e := environment{}

	numEpisodes := 1000
	for episode := 0; episode < numEpisodes; episode++ {
		log.Printf("episode = %v", episode)
		playGame(p1, p2, e)
	}

	// let's see what p1 have learnt
	valueArray := rankStateValues(p1.values)
	for i := 0; i < 5; i++ {
		log.Printf("best state %v, state = %v, value = %v", i, valueArray[i].state, valueArray[i].value)
	}

}
