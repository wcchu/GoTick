package main

import (
	"log"
	"math"
)

type location [2]int

type board [][]string

type environment struct {
	board    board
	winner   string
	gameOver bool
}

// initializeEnvironment initializes environment
func (env *environment) initializeEnvironment() {
	board := make(board, boardSize)
	for irow := range board {
		board[irow] = make([]string, boardSize)
	}
	env.board = board
	env.winner = ""
	env.gameOver = false
	return
}

// getState hashes the game board in the player's perspective into an integer
// NOTE: For each player, each location's status is viewed only as occupied either by him/herself or
//       by the opponent, regardless of the actual symbol ("x" or "o") there.
func (env *environment) getState(symbol string) int64 {
	var k, h, v int64
	for _, row := range env.board {
		for _, element := range row {
			if element == symbol { // occupied by current player
				v = 0
			} else if element == "" { // empty
				v = 1
			} else { // occupied by opponent
				v = 2
			}
			h = h + int64(math.Pow(3, float64(k)))*v
			k = k + 1
		}
	}
	return h
}

// updateGameStatus looks at the board following a move and updates the winner and the game-over
func (env *environment) updateGameStatus(loc location, symbol string) {
	// add new move on the board
	env.board[loc[0]][loc[1]] = symbol
	// update status
	env.winner = getWinner(env.board)
	if env.winner != "" || getEmpties(env.board) == 0 {
		env.gameOver = true
	} else {
		env.gameOver = false
	}
	return
}

func padSymbol(s string) string {
	if len(s) == 0 {
		s = "     "
	} else if len(s) == 1 {
		s = "  " + s + "  "
	} else if len(s) == 4 {
		s = " " + s
	}
	s += "|"
	return s
}

// printBoard prints the board with players on it
func printBoard(b *board) {
	// draw board
	for _, row := range *b {
		log.Print("-------------------")
		rowPrint := "|"
		for _, element := range row {
			rowPrint += padSymbol(element)
		}
		log.Print(rowPrint)
	}
	log.Print("-------------------")
	return
}

// check the current board and find the winner
func getWinner(b board) string {
	// rows
	for _, row := range b {
		for _, p := range symbols {
			if rowFilled(row, p) {
				return p
			}
		}
	}
	// columns
	for icol := range b[0] {
		// collection is an array composed by elements of this column
		collection := []string{}
		for irow := range b {
			collection = append(collection, b[irow][icol])
		}
		for _, p := range symbols {
			if rowFilled(collection, p) {
				return p
			}
		}
	}
	// top-left to bottom-right
	var targetArray []string
	for i := range b {
		targetArray = append(targetArray, b[i][i])
	}
	for _, p := range symbols {
		if rowFilled(targetArray, p) {
			return p
		}
	}
	// top-right to bottom-left
	targetArray = []string{}
	for i := range b {
		targetArray = append(targetArray, b[i][boardSize-1-i])
	}
	for _, p := range symbols {
		if rowFilled(targetArray, p) {
			return p
		}
	}
	// no winner found
	return ""
}

// check number of empty spots
func getEmpties(b board) int {
	n := 0
	for _, row := range b {
		for _, element := range row {
			if element == "" {
				n++
			}
		}
	}
	return n
}

// get reward of the current game state for a certain player: non-zero only if a winner is found
func getReward(w string, s string) float64 {
	if w == s { // this player wins
		return 1.0
	} else if w == "" { // draw or game not over yet
		return 0.0
	}
	// this player loses
	return -1.0
}
