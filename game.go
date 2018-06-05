package main

import (
	"fmt"
	"math"
	"math/rand"
)

func createSessions(players []player) error {
	for {
		var newSess bool
		fmt.Printf("Create a session? (t/f): ")
		_, errS := fmt.Scanf("%t", &newSess)
		if errS != nil {
			return errS
		}
		if !newSess {
			break
		}
		// start a new session
		fmt.Print("players are: \n")
		for i, p := range players {
			fmt.Printf("#%v %v \n", i, p.name)
		}
		var i1, i2, n int
		fmt.Printf("pick two players (# #): ")
		_, errP := fmt.Scanf("%d%d", &i1, &i2)
		if errP != nil {
			return errP
		}
		fmt.Printf("how many episodes: ")
		_, errE := fmt.Scanf("%d", &n)
		if errE != nil {
			return errE
		}
		runSession(&players[i1], &players[i2], n)
	}
	return nil
}

// run an episode and let players (if robot) remember what they've learnt
func runEpisode(p1, p2 *player) {
	var loc location
	var env environment
	var withHuman bool // if episode is played by at least a human, make reports more often
	env.initializeEnvironment()
	if p1.being == "human" || p2.being == "human" {
		withHuman = true
	}

	// p1 always starts first and uses "x"
	p1.symbol = "x"
	p2.symbol = "o"
	if withHuman {
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

	if withHuman {
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
