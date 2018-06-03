package main

import (
	"math/rand"
	"strconv"
)

type stateValues map[int64]float64

type intel struct {
	eps    float64     // epsilon-greedy search
	alp    float64     // learning rate
	mean   float64     // default value for an unseen state
	fluc   float64     // random flucuation for the above default value
	values stateValues // state values that the robot has learnt
	verb   bool        // verbose
}

type player struct {
	name    string  // name of the player
	symbol  string  // "x" plays first, "o" plays second. Each episode assigns symbols randomly.
	being   string  // human or robot
	history []int64 // history of states played in the episode
	intel   intel   // empty if human
}

func (p *player) initializeRobot(name string, eps, alp, mean, fluc float64, verb bool) {
	p.name = name
	p.symbol = ""
	p.being = "robot"
	p.history = []int64{}
	p.intel.eps = eps
	p.intel.alp = alp
	p.intel.mean = mean
	p.intel.fluc = fluc
	p.intel.values = stateValues{}
	p.intel.verb = verb
	return
}

func (p *player) initializeHuman(name string) {
	p.name = name
	p.symbol = ""
	p.being = "human"
	p.history = []int64{}
	p.intel = intel{}
	return
}

// resetHistory resets the state history of a player
func (p *player) resetHistory() {
	p.history = []int64{}
	return
}

// updateHistory append the new state to the player's state history within the episode
func (p *player) updateHistory(state int64) {
	p.history = append(p.history, state)
	return
}

// robotActs determines what location the robot moves to
func (p *player) robotActs(env environment) (actionLocation location) {
	if rand.Float64() < p.intel.eps {
		// take a random action
		possibleLocations := []location{}
		for irow, row := range env.board {
			for ielement, element := range row {
				if element == "" {
					possibleLocations = append(possibleLocations, location{irow, ielement})
				}
			}
		}
		pickedIndex := rand.Intn(len(possibleLocations))
		actionLocation = possibleLocations[pickedIndex]
	} else {
		plan := make(board, boardSize) // only useful for printing out the plan
		// choose the best action based on current values of states
		bestValue := -1.0
		for irow, row := range env.board {
			plan[irow] = make([]string, boardSize)
			for ielement, element := range row {
				plan[irow][ielement] = "  " + element + " "
				if element == "" { // location is empty; look up value if move here
					env.board[irow][ielement] = p.symbol // assume if player moves here
					testState := env.getState(p.symbol)  // state after this move
					testWinner := getWinner(env.board)   // winner after this move
					testEmpties := getEmpties(env.board) // empty spots after this move
					env.board[irow][ielement] = ""       // revert this action
					// get value for the test state
					testValue, ok := p.intel.values[testState]
					if !ok { // agent has no record of this state
						if testWinner != "" || testEmpties == 0 { // test state is final state, use reward as value
							testValue = getReward(testWinner, p.symbol)
						} else { // test state is not final state, use default value
							testValue = defaultValue(p.intel.mean, p.intel.fluc)
						}
					}
					plan[irow][ielement] = " " + strconv.FormatFloat(testValue, 'f', 2, 64)
					// update move and best value
					if testValue > bestValue {
						bestValue = testValue
						actionLocation = location{irow, ielement}
					}
				}
			}
		}

	}
	return actionLocation
}

// robotUpdatesvalues should only be run at the end of an episode
// Use the update rule: V(s) = V(s) + alpha*(V(s') - V(s))
func (p *player) robotUpdatesValues(env environment) {
	reward := getReward(env.winner, p.symbol)
	target := reward
	// loop backward from the last state to the first along history
	// i is the index of a.history array
	for i := len(p.history) - 1; i >= 0; i-- {
		state := p.history[i]
		var updatedValue float64
		if i == len(p.history)-1 {
			// If the state is the final state, the value is the reward. The agent should
			// just remember this state-value pair immediately.
			updatedValue = target
		} else {
			// If the state is not the final state, update its value in the regular way
			existingValue, ok := p.intel.values[state]
			if !ok {
				// agent has no values of this state, set to defaultValue
				existingValue = defaultValue(p.intel.mean, p.intel.fluc)
			}
			updatedValue = existingValue + p.intel.alp*(target-existingValue)
		}
		p.intel.values[state] = updatedValue
		target = updatedValue
	}
	p.resetHistory() // state history is reset but values of state values is kept
	return
}

// defaultValue generates a value of certain mean and certain randomness
func defaultValue(defaultMean, fluctuation float64) float64 {
	return defaultMean + fluctuation*(rand.Float64()-0.5)
}
