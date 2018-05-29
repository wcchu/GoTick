// Use the update rule: V(s) = V(s) + alpha*(V(s') - V(s))
package main

import (
	"log"
	"math"
	"math/rand"
)

// BoardSize is the length/width of the board
const BoardSize = 3 // TODO: utilize this const

// structs: agent, human, environment

// values is an agent's memory of each state's value
type values map[int64]float64

type agent struct {
	epsilon      float64
	alpha        float64
	identity     int
	stateHistory []int64
	values       values
}

type human struct {
	identity     int
	stateHistory []int64
}

type location [2]int

type environment struct {
	board    [][]int
	winner   int
	gameOver bool
}

// methods for agent

func (a *agent) initializeAgent(pid int) {
	a.identity = pid
	a.epsilon = 0.1
	a.alpha = 0.5
	a.stateHistory = []int64{}
	a.values = values{}
	return
}

func (a *agent) resetAgentHistory() {
	a.stateHistory = []int64{}
	return
}

func (a *agent) setAgentIdentity(id int) {
	a.identity = id
	return
}

func (a *agent) actAgent(env environment) (actionLocation location) {
	if rand.Float64() < a.epsilon {
		// take a random action
		possibleLocations := []location{}
		for irow, row := range env.board {
			for ielement, element := range row {
				if element == 0 {
					possibleLocations = append(possibleLocations, location{irow, ielement})
				}
			}
		}
		pickedIndex := rand.Intn(len(possibleLocations))
		actionLocation = possibleLocations[pickedIndex]
		log.Printf("player %v acts randomly, picks location %v", a.identity, actionLocation)
	} else {
		// choose the best action based on current values of states
		bestValue := -1.0
		for irow, row := range env.board {
			for ielement, element := range row {
				if element == 0 { // location is empty; look up value if move here
					env.board[irow][ielement] = a.identity // assume if agent move here
					state := env.getState()                // state with this move
					env.board[irow][ielement] = 0          // delete this action
					// look up value for the hypothetical state
					stateValue, ok := a.values[state]
					if !ok {
						// agent has no record of this state, use a default value
						// TODO: generate default value in different ways
						stateValue = 0.5
					}
					if stateValue > bestValue { // update move and best value
						bestValue = stateValue
						actionLocation = location{irow, ielement}
					}
				}
			}
		}
		log.Printf("player %v acts based on best value, picks location %v", a.identity, actionLocation)
	}
	return actionLocation
}

func (a *agent) updateStateHistory(state int64) {
	a.stateHistory = append(a.stateHistory, state)
	return
}

// updateValues should only be run at the end of an episode
func (a *agent) updateValues(env environment) {
	reward := env.reward(a.identity)
	target := reward
	// loop backward from the last state to the first along stateHistory
	// i is the index of a.stateHistory array
	for i := len(a.stateHistory) - 1; i >= 0; i-- {
		state := a.stateHistory[i]
		existingValue, ok := a.values[state]
		if !ok {
			// agent has no record of this state, use a default value
			// TODO: generate default value in different ways
			existingValue = 0.5
		}
		updatedValue := existingValue + a.alpha*(target-existingValue)
		a.values[state] = updatedValue
		target = updatedValue
	}
	a.resetAgentHistory() // state history is reset but memory of state values is kept
	log.Printf("agent %v's memory size is %v, content is %+v", a.identity, len(a.values), a.values)
	return
}

// methods for human

// methods for environment

func (e *environment) initializeEnvironment() {
	indices := []int{0, 1, 2}
	board := make([][]int, len(indices))
	for i := range indices {
		row := make([]int, len(indices))
		for j := range indices {
			row[j] = 0
		}
		board[i] = row
	}
	e.board = board
	e.winner = 0
	e.gameOver = false
	return
}

func (e *environment) getState() int64 { // "state" is an encoded description of the whole board
	var k, h, v int64
	for _, row := range e.board {
		for _, element := range row {
			if element == 0 {
				v = 0
			} else if element == -1 {
				v = 1
			} else if element == 1 {
				v = 2
			}
			h = h + int64(math.Pow(3, float64(k)))*v
			k = k + 1
		}
	}
	return h
}

// update gameOver and winner
func (e *environment) updateGameStatus(l location, p int) {
	// add a player on the board
	e.board[l[0]][l[1]] = p

	players := [2]int{-1, 1}

	// check rows
	for _, row := range e.board {
		for _, player := range players {
			if arrayEqualsInteger(row, player) {
				e.winner = player
				e.gameOver = true
				return
			}
		}
	}

	// check columns
	for icol := range e.board[0] {
		// collection is an array composed by elements of this column
		collection := []int{}
		for irow := range e.board {
			collection = append(collection, e.board[irow][icol])
		}
		for _, player := range players {
			if arrayEqualsInteger(collection, player) {
				e.winner = player
				e.gameOver = true
				return
			}
		}
	}

	// check diagonal top-left to bottom-right
	var targetArray []int
	for i := range e.board {
		targetArray = append(targetArray, e.board[i][i])
	}
	for _, player := range players {
		if arrayEqualsInteger(targetArray, player) {
			e.winner = player
			e.gameOver = true
			return
		}
	}

	// check diagonal top-right to bottom-left
	targetArray = []int{}
	for i := range e.board {
		targetArray = append(targetArray, e.board[i][len(e.board)-1-i])
	}
	for _, player := range players {
		if arrayEqualsInteger(targetArray, player) {
			e.winner = player
			e.gameOver = true
			return
		}
	}

	// no winner found
	e.winner = 0

	// check draw
	e.gameOver = true // temporarily assum game is over. But is it true?
	for _, row := range e.board {
		for _, element := range row {
			if element == 0 { // there are still unoccupied spots, game it not over.
				e.gameOver = false
			}
		}
	}

	return
}

func (e *environment) printBoard() {
	// draw board
	for _, row := range e.board {
		log.Print("------------------")
		rowPrint := ""
		for _, element := range row {
			if element == -1 {
				rowPrint += "x  |"
			} else if element == 1 {
				rowPrint += "o  |"
			} else {
				rowPrint += "   |"
			}
		}
		log.Print(rowPrint)
	}
	log.Print("------------------")
	return
}

func (e *environment) reward(id int) float64 {
	if e.gameOver && e.winner == id {
		return 1.0
	}
	return 0.0
}

// general functions

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

// run game

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
