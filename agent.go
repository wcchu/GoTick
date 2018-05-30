package main

import (
	"log"
	"math/rand"
)

// values is an agent's memory of each state's value
type values map[int64]float64

type agent struct {
	epsilon      float64
	alpha        float64
	identity     int
	stateHistory []int64
	values       values
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
		log.Printf("player %v, random action, %v", a.identity, actionLocation)
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
		log.Printf("player %v, best value %v, %v", a.identity, bestValue, actionLocation)
	}
	return actionLocation
}

func (a *agent) updateStateHistory(state int64) {
	a.stateHistory = append(a.stateHistory, state)
	return
}

// updateValues should only be run at the end of an episode
// Use the update rule: V(s) = V(s) + alpha*(V(s') - V(s))
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
	log.Printf("agent %v's memory size is %v", a.identity, len(a.values))
	return
}
