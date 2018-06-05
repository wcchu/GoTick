package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
)

type stateValues map[int64]float64

type robotSpecs struct {
	eps  float64 // epsilon-greedy search
	alp  float64 // learning rate
	mean float64 // default value for an unseen state
	fluc float64 // random flucuation for the above default value
	draw float64 // reward for draw game (between winning 1 and losing -1)
}

type mind struct {
	specs  robotSpecs
	values stateValues // state values that the robot has learnt
	verb   bool        // verbose
}

type player struct {
	name    string  // name of the player
	symbol  string  // "x" plays first, "o" plays second. Each episode assigns symbols randomly.
	being   string  // human or robot
	history []int64 // history of states played in the episode
	wins    int     // number of wins
	mind    mind    // empty if human
}

func createPlayers() ([]player, error) {
	var N uint
	fmt.Print("Enter number of players: ")
	_, errN := fmt.Scanf("%d", &N)
	if errN == nil {
		players := make([]player, N)
		for i := range players {
			var name string
			var isRobot bool
			// name
			fmt.Printf("Enter name of player #%v: ", i)
			_, errName := fmt.Scanf("%s", &name)
			if errName != nil {
				return []player{}, errName
			}
			// being
			fmt.Printf("Robot? (t/f): ")
			_, errIsRobot := fmt.Scanf("%t", &isRobot)
			if errIsRobot != nil {
				return []player{}, errIsRobot
			}
			if isRobot {
				// specs
				var e, a, m, f, d float64
				fmt.Printf("Specs (eps alp mean fluc draw): ")
				_, errSpecs := fmt.Scanf("%f%f%f%f%f", &e, &a, &m, &f, &d)
				if errSpecs != nil {
					return []player{}, errSpecs
				}
				players[i].initializeRobot(name, robotSpecs{eps: e, alp: a, mean: m, fluc: f, draw: d}, false)
			} else {
				players[i].initializeHuman(name)
			}
		}
		return players, nil
	}
	return []player{}, errN
}

func (p *player) initializeRobot(name string, rs robotSpecs, verb bool) {
	p.name = name
	p.symbol = ""
	p.being = "robot"
	p.history = []int64{}
	p.wins = 0
	p.mind.specs = rs
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

// updateHistory append the new state to the player's state history within the episode
func (p *player) updateHistory(state int64) {
	p.history = append(p.history, state)
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
	fmt.Printf("%v's %v state values saved into %v \n", p.name, len(p.mind.values), filename)
	return
}

func (p *player) playerActs(env environment) (actionLocation location) {
	if p.being == "robot" {
		return p.robotActs(env)
	} else if p.being == "human" {
		return p.humanActs(env)
	}
	fmt.Printf("player %v is a non-being; the game board explodes \n", p.name)
	os.Exit(1)
	return
}

func (p *player) humanActs(env environment) (actionLocation location) {
	printBoard(&env.board)
	for {
		var x, y int
		fmt.Print("Enter location (x y): ")
		_, err := fmt.Scanf("%d%d", &x, &y)
		if err == nil {
			l := location{x, y}
			if env.board[l[0]][l[1]] == "" {
				log.Printf("You are making a move to %v", l)
				return l
			}
		}
		// invalid move, re-enter location
		fmt.Print("invalid move \n")
	}
}

// robotActs determines what location the robot moves to
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
					env.board[irow][ielement] = p.symbol // board after this move
					testState := env.getState(p.symbol)  // state after this move
					testWinner := getWinner(env.board)   // winner after this move
					testEmpties := getEmpties(env.board) // empty spots after this move
					env.board[irow][ielement] = ""       // revert this action
					// get value for the test state
					testValue, ok := p.mind.values[testState]
					if !ok { // there's no record of this state
						if testWinner != "" || testEmpties == 0 { // test state is final state, use reward as value
							testValue = getReward(testWinner, p.symbol, p.mind.specs.draw)
						} else { // test state is not final state, use default value
							testValue = defaultValue(p.mind.specs.mean, p.mind.specs.fluc)
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
		if p.mind.verb {
			log.Printf("player %v(%v)'s plan board:", p.name, p.symbol)
			printBoard(&plan)
			log.Printf("player %v(%v) takes action at %v \n", p.name, p.symbol, actionLocation)
		}
	}
	return actionLocation
}

func (p *player) updateValues(env environment) {
	if p.being == "robot" {
		p.robotUpdatesValues(env)
	}
	return
}

// robotUpdatesvalues should only be run at the end of an episode
// Use the update rule: V(s) = V(s) + alpha*(V(s') - V(s))
func (p *player) robotUpdatesValues(env environment) {
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
				existingValue = defaultValue(p.mind.specs.mean, p.mind.specs.fluc)
			}
			updatedValue = existingValue + p.mind.specs.alp*(target-existingValue)
		}
		p.mind.values[state] = updatedValue
		target = updatedValue
	}
	if env.winner == p.symbol {
		p.wins++
	}
	p.resetHistory() // state history is reset but values of state values is kept
	return
}

// defaultValue generates a value of certain mean and certain randomness
func defaultValue(defaultMean, fluctuation float64) float64 {
	return defaultMean + fluctuation*(rand.Float64()-0.5)
}