package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
)

type stateCounts map[int64]uint            // each state maps to how many times it's encountered
type stateValues map[int64]float64         // each state maps to a value
type stateValueHistory map[int64][]float64 // each state maps to an array of values

type robotSpecs struct {
	alp float64 // learning rate; if zero, use weighted average to update the value
	eps float64 // epsilon-greedy search
	gam float64 // discount factor
}

type mind struct {
	specs   robotSpecs
	counts  stateCounts       // count number of times each state has appeared
	valhist stateValueHistory // historic values of N oldest states in the robot's record
	values  stateValues       // most updated values of the robot's known states
	verb    bool              // verbose
}

type player struct {
	name    string  // name of the player
	symbol  string  // "x" plays first, "o" plays second. Each episode assigns symbols randomly.
	being   string  // human or robot
	history []int64 // history of states played in the episode
	wins    int     // number of wins
	mind    mind    // empty if human
}

type playerPair [2]player

func createPlayers() []player {
	// number of players
	var N uint
	for {
		fmt.Print("Enter number of players: ")
		_, err := fmt.Scanf("%d", &N)
		if err == nil {
			break
		}
	}

	// define each player
	players := make([]player, N)
	for i := range players {
		var name string
		var isRobot bool
		// name
		for {
			fmt.Printf("Enter name of player #%v: ", i)
			_, err := fmt.Scanf("%s", &name)
			if err == nil {
				break
			}
		}
		// being
		for {
			fmt.Printf("robot? (t/f): ")
			_, err := fmt.Scanf("%t", &isRobot)
			if err == nil {
				break
			}
		}
		if isRobot {
			// specs
			var a, e, g float64
			fmt.Printf("specs (alp eps gam) / click enter to use default values: ")
			_, err := fmt.Scanf("%f%f%f", &a, &e, &g)
			if err != nil {
				a, e, g = alpha, epsilon, gamma
				fmt.Printf("use default specs %v %v %v \n", a, e, g)
			}
			players[i].initializeRobot(name, robotSpecs{alp: a, eps: e, gam: g}, false)
		} else {
			players[i].initializeHuman(name)
		}
	}
	fmt.Print("*** Done creating players *** \n\n")
	return players
}

func (p *player) initializeRobot(name string, rs robotSpecs, verb bool) {
	p.name = name
	p.symbol = ""
	p.being = "robot"
	p.history = []int64{}
	p.wins = 0
	p.mind.specs = rs
	p.mind.counts = stateCounts{}
	p.mind.valhist = stateValueHistory{}
	p.mind.values = stateValues{}
	p.mind.verb = verb
	return
}

func (p *player) initializeHuman(name string) {
	p.name = name
	p.symbol = ""
	p.being = "human"
	p.history = []int64{}
	p.wins = 0
	p.mind = mind{}
	return
}

// resetHistory resets the state history of a player
func (p *player) resetHistory() {
	p.history = []int64{}
	return
}

// append the new state to the player's state history within the episode
func (p *player) updateStateSequence(state int64) {
	p.history = append(p.history, state)
	return
}

func (p *player) getOldestNStates(state int64) {
	if p.being == "robot" && len(p.mind.valhist) < nOldest { // record up to N states in valhist
		p.mind.valhist[state] = []float64{}
	}
	return
}

func (p *player) playerActs(env environment) (actionLocation location) {
	if p.being == "robot" {
		return p.robotActs(env)
	} else if p.being == "human" {
		return p.humanActs(env)
	}
	fmt.Printf("player %v is an unknown creature; the game board explodes \n", p.name)
	os.Exit(1)
	return
}

func (p *player) humanActs(env environment) (actionLocation location) {
	printBoard(&env.board, true)
	for {
		var x, y int
		fmt.Print("Enter location (x y): ")
		_, err := fmt.Scanf("%d%d", &x, &y)
		if err == nil {
			l := location{x, y}
			if env.board[l[0]][l[1]] == "" {
				fmt.Printf("you are making a move to %v \n", l)
				return l
			}
		}
		// invalid move, re-enter location
		fmt.Print("invalid move \n")
	}
}

// determine what location the robot moves to
func (p *player) robotActs(env environment) (actionLocation location) {
	if rand.Float64() < p.mind.specs.eps {
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
		bestGain := -1.0
		for irow, row := range env.board {
			plan[irow] = make([]string, boardSize)
			for ielement, element := range row {
				plan[irow][ielement] = element
				if element == "" { // location is empty; find value if player moves here
					env.board[irow][ielement] = p.symbol            // board after this move
					testState := boardToState(&env.board, p.symbol) // state after this move
					testWinner := getWinner(env.board)              // winner after this move
					testEmpties := getEmpties(env.board)            // empty spots after this move
					env.board[irow][ielement] = ""                  // revert this action
					// get gain of the test state
					var testGain float64
					if testWinner != "" || testEmpties == 0 {
						// test state is final state, reward is non-zero, value is zero
						testGain = getReward(testWinner, p.symbol)
					} else {
						testValue, ok := p.mind.values[testState]
						if !ok { // there's no record of this state, use default value
							testValue = defaultValue()
						}
						testGain = p.mind.specs.gam * testValue
					}
					plan[irow][ielement] = strconv.FormatFloat(testGain, 'f', 2, 64)
					if testGain > bestGain {
						bestGain = testGain
						actionLocation = location{irow, ielement}
					}
				}
			}
		}
		if p.mind.verb || printSteps {
			fmt.Printf("player %v(%v)'s plan board: \n", p.name, p.symbol)
			printBoard(&plan, true)
			fmt.Printf("player %v(%v) takes action at %v \n", p.name, p.symbol, actionLocation)
		}
	}
	return actionLocation
}

// append the state-values learnt in each episode to the player's memory
func (p *player) updatePlayerRecord(env environment) {
	if p.symbol == env.winner {
		p.wins++
	}
	if p.being == "robot" {
		p.updateStateValues(env)
		p.updateStateValueHistory(env)
		p.updateStateCounts()
	}
	p.resetHistory()
	return
}

// should only be run at the end of an episode
// update rule: V(s) = V(s) + alpha*(V(s') - V(s))
func (p *player) updateStateValues(env environment) {
	gains := make(map[int64]float64, len(p.history)) // values learned through this episode
	finalReward := getReward(env.winner, p.symbol)
	// loop backward from the last state to the first along history of this episode
	// i is the index of history array
	gain := 0.0
	for i := len(p.history) - 1; i >= 0; i-- {
		state := p.history[i]
		gains[state] = gain
		var reward float64
		if i == len(p.history)-1 {
			reward = finalReward
		}
		gain = reward + p.mind.specs.gam*gain
	}
	// update the state values
	for state, gain := range gains {
		if p.mind.specs.alp == 0.0 {
			// update V by weighted average between new and existing values
			count, ok := p.mind.counts[state]
			if !ok {
				count = 0
			}
			p.mind.values[state] = (float64(count)*p.mind.values[state] + gain) / float64(count+1)
		} else {
			// update V by correction to the new value with learning rate
			oldValue, ok := p.mind.values[state]
			if !ok {
				oldValue = defaultValue()
			}
			p.mind.values[state] = oldValue + p.mind.specs.alp*(gain-oldValue)
		}
	}
	return
}

// generate a value of certain mean and certain randomness
func defaultValue() float64 {
	return initialValue + fluctuation*(rand.Float64()-0.5)
}

// should be run right after updateStateValues()
func (p *player) updateStateValueHistory(env environment) {
	for state := range p.mind.valhist {
		p.mind.valhist[state] = append(p.mind.valhist[state], p.mind.values[state])
	}
	return
}

// update the record of how many times each state has appeared
func (p *player) updateStateCounts() {
	for _, state := range p.history {
		count, ok := p.mind.counts[state]
		if !ok { // this state appears the first time
			count = 0
		}
		p.mind.counts[state] = count + 1
	}
	return
}
