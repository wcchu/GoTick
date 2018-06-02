package main

import (
	"log"
	"math"
)

type location [2]int

type board [][]int

type environment struct {
	board    board
	winner   int
	gameOver bool
}

// initializeEnvironment initializes environment
func (env *environment) initializeEnvironment() {
	board := make(board, BoardSize)
	for irow := range board {
		board[irow] = make([]int, BoardSize)
	}
	env.board = board
	env.winner = 0
	env.gameOver = false
	return
}

// getState hashes the game board including player locations into an integer (state)
func (env *environment) getState() int64 { // "state" is an encoded description of the whole board
	var k, h, v int64
	for _, row := range env.board {
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

// checkBoardForWinner checks the board and finds the winner of the game
// (if game is tie or not over yet, winner is 0)
func checkBoardForWinner(b board) int {
	ps := [2]int{-1, 1}

	// check rows
	for _, row := range b {
		for _, p := range ps {
			if arrayEqualsInteger(row, p) {
				return p
			}
		}
	}

	// check columns
	for icol := range b[0] {
		// collection is an array composed by elements of this column
		collection := []int{}
		for irow := range b {
			collection = append(collection, b[irow][icol])
		}
		for _, p := range ps {
			if arrayEqualsInteger(collection, p) {
				return p
			}
		}
	}

	// check diagonal top-left to bottom-right
	var targetArray []int
	for i := range b {
		targetArray = append(targetArray, b[i][i])
	}
	for _, p := range ps {
		if arrayEqualsInteger(targetArray, p) {
			return p
		}
	}

	// check diagonal top-right to bottom-left
	targetArray = []int{}
	for i := range b {
		targetArray = append(targetArray, b[i][BoardSize-1-i])
	}
	for _, p := range ps {
		if arrayEqualsInteger(targetArray, p) {
			return p
		}
	}

	// no player wins
	return 0
}

// checkBoardForOccupancy checks the board and finds how many locations are still unoccupied
func checkBoardForOccupancy(b board) int {
	availables := 0
	for _, row := range b {
		for _, element := range row {
			if element == 0 {
				availables++
			}
		}
	}
	return availables
}

// updateGameStatus looks at the board following a move and updates the winner and the game-over
func (env *environment) updateGameStatus(loc location, pid int) {
	// add the new move on the board
	env.board[loc[0]][loc[1]] = pid

	// update winner
	env.winner = checkBoardForWinner(env.board)

	// update gameOver
	if env.winner != 0 || checkBoardForOccupancy(env.board) == 0 {
		env.gameOver = true
		return
	}
	env.gameOver = false
	return
}

// printBoard prints the board with players on it
func (env *environment) printBoard() {
	// draw board
	for _, row := range env.board {
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
func (env *environment) reward(pid int) float64 {
	if env.gameOver {
		if env.winner == pid { // this player wins
			return 1.0
		} else if env.winner == 0 { // tie
			return 0.5
		}
	}
	// player loses or game not over yet
	return 0.0
}
