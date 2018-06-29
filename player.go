package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
)

type stateCounts map[int64]uint            // each state maps to how many times it's encountered
type stateValues map[int64]float64         // each state maps to a value
type stateValueHistory map[int64][]float64 // each state maps to an array of values

type robotSpecs struct {
	eps  float64 // epsilon-greedy search
	alp  float64 // learning rate
	draw float64 // reward for draw game (between winning 1 and losing -1)
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
			var e, a, d float64
			fmt.Printf("specs (eps alp draw) / click enter to use default values: ")
			_, err := fmt.Scanf("%f%f%f", &e, &a, &d)
			if err != nil {
				e, a, d = epsilon, alpha, drawReward
				fmt.Printf("use default specs %v %v %v \n", e, a, d)
			}
			players[i].initializeRobot(name, robotSpecs{eps: e, alp: a, draw: d}, false)
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

// write state values of the player to a csv file
func (p *player) exportValues() {
	filename := p.name + "_values.csv"
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for state, value := range p.mind.values {
		row := []string{strconv.FormatInt(state, 10), strconv.FormatFloat(value, 'g', 5, 64)}
		err := writer.Write(row)
		if err != nil {
			log.Fatal("Cannot write to file", err)
		}
	}
	fmt.Printf("%v has %v state-values, saved into %v \n", p.name, len(p.mind.values), filename)
	return
}

// write state values of the player to a csv file
func (p *player) exportValueHistory() {
	filename := p.name + "_oldest_states_hist.csv"
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	filename2 := p.name + "_oldest_states.txt"

	var s string
	for state, valueHistory := range p.mind.valhist {

		b := stateToBoard(state, p.symbol)
		s = s + strconv.FormatInt(state, 10) + "\n" + printBoard(&b, false) + "\n"

		for time, value := range valueHistory {
			row := []string{
				strconv.FormatInt(state, 10),
				strconv.Itoa(time),
				strconv.FormatFloat(value, 'g', 5, 64)}
			err := writer.Write(row)
			if err != nil {
				log.Fatal("Cannot write to file", err)
			}
		}
	}

	d := []byte(s)
	ioutil.WriteFile(filename2, d, 0644)

	fmt.Printf("%v's value histories of the oldest %v states saved into %v \n", p.name, len(p.mind.valhist), filename)
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
		bestValue := -1.0
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
					// get value for the test state
					testValue, ok := p.mind.values[testState]
					if !ok { // there's no record of this state
						if testWinner != "" || testEmpties == 0 { // test state is final state, use reward as value
							testValue = getReward(testWinner, p.symbol, p.mind.specs.draw)
						} else { // test state is not final state, use default value
							testValue = defaultValue()
						}
					}
					plan[irow][ielement] = strconv.FormatFloat(testValue, 'f', 2, 64)
					// update move and best value
					if testValue > bestValue {
						bestValue = testValue
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
	return
}

// should only be run at the end of an episode
// update rule: V(s) = V(s) + alpha*(V(s') - V(s))
func (p *player) updateStateValues(env environment) {
	reward := getReward(env.winner, p.symbol, p.mind.specs.draw)
	target := reward
	// loop backward from the last state to the first along history
	// i is the index of a.history array
	for i := len(p.history) - 1; i >= 0; i-- {
		state := p.history[i]
		var updatedValue float64
		if i == len(p.history)-1 {
			// If the state is the final state, the value is the reward. The robot should
			// just remember this state-value pair immediately.
			updatedValue = target
		} else {
			// If the state is not the final state, update its value in the regular way
			existingValue, ok := p.mind.values[state]
			if !ok {
				existingValue = defaultValue()
			}
			updatedValue = existingValue + p.mind.specs.alp*(target-existingValue)
		}
		p.mind.values[state] = updatedValue
		target = updatedValue
	}
	p.resetHistory() // state history is reset but values of state values is kept
	return
}

// generate a value of certain mean and certain randomness
func defaultValue() float64 {
	v := initialValue + fluctuation*(rand.Float64()-0.5)
	v = math.Min(math.Max(v, 0.0), 1.0) // the value is bound by [0, 1]
	return v
}

// should be run right after updateStateValues()
func (p *player) updateStateValueHistory(env environment) {
	for state := range p.mind.valhist {
		p.mind.valhist[state] = append(p.mind.valhist[state], p.mind.values[state])
	}
	return
}

//
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
