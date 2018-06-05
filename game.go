package main

import (
	"fmt"
	"math"
	"math/rand"
)

// run an episode and let players (if robot) remember what they've learnt
func runEpisode(p1, p2 *player) {
	var loc location
	var env environment
	var crossSpecies bool // if episode is cross-species, make reports more often
	env.initializeEnvironment()
	if p1.being != p2.being {
		crossSpecies = true
	}

	// p1 always starts first and uses "x"
	p1.symbol = "x"
	p2.symbol = "o"
	if crossSpecies {
		fmt.Printf("\n %v(%v) starts first \n", p1.name, p1.symbol)
	}
	s := "o" // current player
	for !env.gameOver {
		// switch player and take action
		if s == "o" {
			s = "x"
			loc = p1.playerActs(env)
		} else {
			s = "o"
			loc = p2.playerActs(env)
		}

		// update environment by the action
		env.updateGameStatus(loc, s)

		// update state history
		state := env.getState(s)
		p1.updateHistory(state)
		p2.updateHistory(state)
	}

	if crossSpecies {
		env.reportEpisode(p1, p2)
	}

	// grow some brain
	p1.updateValues(env)
	p2.updateValues(env)

	return
}

func runSession(p1, p2 *player, nEpisodes int) {
	for episode := 0; episode < nEpisodes; episode++ {
		if math.Mod(float64(episode+1), 1000) == 0 {
			fmt.Printf("episode #%v \n", episode)
		}
		// for each episode, randomly pick the first player
		if rand.Float64() < 0.5 {
			runEpisode(p1, p2)
		} else {
			runEpisode(p2, p1)
		}
	}
	fmt.Printf("Session ends - %v won %v times; %v won %v times \n", p1.name, p1.wins, p2.name, p2.wins)
	if p1.being == "robot" {
		p1.exportValues()
	}
	if p2.being == "robot" {
		p2.exportValues()
	}
}
