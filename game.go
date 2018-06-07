package main

import (
	"fmt"
	"math"
	"math/rand"
)

func createSessions(players []player) {
	for {

		// user input
		var newSess bool
		for {
			fmt.Printf("Create a session? (t/f): ")
			_, err := fmt.Scanf("%t", &newSess)
			if err == nil {
				break
			}
		}
		if !newSess {
			break
		}

		// start a new session
		fmt.Print("available players are: \n")
		for i, p := range players {
			fmt.Printf("#%v %v \n", i, p.name)
		}
		var i1, i2, n int
		if len(players) > 2 { // more than 2 available players, choose 2
			for {
				fmt.Printf("pick two players (# #): ")
				_, err := fmt.Scanf("%d%d", &i1, &i2)
				if err == nil {
					break
				}
			}
		} else {
			i1, i2 = 0, 1
		}

		for {
			fmt.Printf("how many episodes: ")
			_, err := fmt.Scanf("%d", &n)
			if err == nil {
				break
			}
		}

		// run session
		fmt.Printf("*** Session starts: %v and %v play %v episodes *** \n", players[i1].name, players[i2].name, n)
		runSession(&players[i1], &players[i2], n)
	}

	return
}

func runSession(p1, p2 *player, nEpisodes int) {
	// set up reporting parameters
	r := false                // report more frequently
	v := false                // robot is verbose
	if p1.being != p2.being { // human vs robot
		r = true // report more frequently
		for {
			fmt.Printf("set robot to verbose? (t/f): ")
			_, err := fmt.Scanf("%t", &v)
			if err == nil {
				break
			}
		}
	}
	if p1.being == "robot" {
		p1.mind.verb = v
	}
	if p2.being == "robot" {
		p2.mind.verb = v
	}

	// run episodes
	for episode := 0; episode < nEpisodes; episode++ {
		if math.Mod(float64(episode+1), 1000) == 0 && p1.being == p2.being {
			fmt.Printf("episode #%v \n", episode)
		}
		// for each episode, randomly pick the first player
		if rand.Float64() < 0.5 {
			runEpisode(p1, p2, r)
		} else {
			runEpisode(p2, p1, r)
		}
	}
	if p1.being == "robot" {
		p1.exportValues()
		p1.exportValueHistory()
	}
	if p2.being == "robot" {
		p2.exportValues()
		p2.exportValueHistory()
	}
	fmt.Printf("*** Session ends - %v won %v times / %v won %v times *** \n\n", p1.name, p1.wins, p2.name, p2.wins)

	return
}

// run an episode and let players (if robot) remember what they've learnt
func runEpisode(p1, p2 *player, report bool) {
	var loc location
	var env environment
	env.initializeEnvironment()

	// p1 always starts first and uses "x"
	p1.symbol = "x"
	p2.symbol = "o"
	if report {
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
		p1.updateStateSequence(state)
		p2.updateStateSequence(state)

		// remember 5 oldest states
		p1.getFiveOldestStates(state)
		p2.getFiveOldestStates(state)
	}

	if report {
		env.reportEpisode(p1, p2)
	}

	// grow some brain
	p1.updatePlayerRecord(env)
	p2.updatePlayerRecord(env)

	return
}
