package main

import (
	"log"
)

// BoardSize is the length/width of the board
const BoardSize = 3 // TODO: utilize this const

func arrayEqualsInteger(array []int, integer int) bool {
	for _, element := range array {
		if element != integer {
			return false
		}
	}
	return true
}

// run an episode
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

func main() {
	p1 := agent{}
	p1.initializeAgent(-1)
	p2 := agent{}
	p2.initializeAgent(1)
	e := environment{}

	numEpisodes := 5
	for episode := 0; episode < numEpisodes; episode++ {
		log.Printf("episode = %v", episode)
		playGame(p1, p2, e)
	}
}
