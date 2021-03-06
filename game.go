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
		runSession(&playerPair{players[i1], players[i2]}, n)
	}

	return
}

func runSession(ps *playerPair, nEpisodes int) {
	// set up reporting parameters
	r := false                      // report more frequently
	v := false                      // robot is verbose
	if ps[0].being != ps[1].being { // human vs robot
		r = true // report more frequently
		for {
			fmt.Printf("set robot to verbose? (t/f): ")
			_, err := fmt.Scanf("%t", &v)
			if err == nil {
				break
			}
		}
	}
	for i := range ps {
		if ps[i].being == "robot" {
			ps[i].mind.verb = v
		}
	}

	// run episodes
	for episode := 0; episode < nEpisodes; episode++ {
		epiNum := episode + 1 // epiNum starts from 1 which is more human readable
		if math.Mod(float64(epiNum), nPrintEpisode) == 0 && ps[0].being == ps[1].being {
			fmt.Printf("episode #%v \n", epiNum)
		}
		runEpisode(ps, r, episode == 0)
	}

	// robot export values
	for i := range ps {
		if ps[i].being == "robot" {
			exportValues(ps[i].name, ps[i].mind.values)
			exportValueHistory(ps[i].name, ps[i].mind.demohist)
		}
	}
	fmt.Printf("*** Session ends - %v won %v times / %v won %v times *** \n\n", ps[0].name, ps[0].wins, ps[1].name, ps[1].wins)

	return
}

// run an episode and let players (if robot) remember what they've learnt
func runEpisode(ps *playerPair, report, firstEpisode bool) {
	var loc location
	var env environment
	if printSteps { // global const to force reporting
		report = true
	}
	env.initializeEnvironment()

	// randomly assign 0 or 1 as the first player ("x")
	first := rand.Perm(2)[0]
	second := 1 - first

	// first player uses "x"
	ps[first].symbol = "x"
	ps[second].symbol = "o"
	if report {
		fmt.Printf("\n %v(%v) starts first \n", ps[first].name, ps[first].symbol)
	}
	s := "o" // current player
	for !env.gameOver {
		// switch player and take action
		if s == "o" {
			s = "x"
			loc = ps[first].playerActs(env)
		} else {
			s = "o"
			loc = ps[second].playerActs(env)
		}

		// update environment by the action
		env.updateGameStatus(loc, s)

		// update state history and remember the demo states
		for i := range ps {
			// The same board is encoded differently by the two players;
			// each location is viewed not as "x" or "o", but instead as Me or You.
			state := boardToState(&env.board, ps[i].symbol)
			ps[i].updateStateSequence(state)
		}
	}

	if firstEpisode {
		for i := range ps {
			ps[i].getDemoStates()
		}
	}

	if report {
		env.summarizeEpisode(&ps[first], &ps[second])
	}

	// grow some brain
	ps[first].updatePlayerRecord(env)
	ps[second].updatePlayerRecord(env)

	return
}
