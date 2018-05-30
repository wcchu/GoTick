package main

import (
	"log"
	"math"
)

type location [2]int

type environment struct {
	board    [][]int
	winner   int
	gameOver bool
}

// initializeEnvironment initializes environment
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

// getState hashes the game board including player locations into an integer (state)
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

// updateGameStatus looks at the current board and updates the winner and the game-over
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

// printBoard prints the board with players on it
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

// reward tells the reward of the current game state for a certain player
func (e *environment) reward(player int) float64 {
	if e.gameOver && e.winner == player {
		return 1.0
	}
	return 0.0
}
