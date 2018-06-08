package main

import (
	"fmt"
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

func (env *environment) reportEpisode(p1, p2 *player) {
	printBoard(&env.board)
	if env.gameOver {
		fmt.Print("Game Over - ")
	}
	if env.winner != "" { // there's a winner
		if env.winner == p1.symbol {
			fmt.Printf("%v is the winner \n\n", p1.name)
		} else {
			fmt.Printf("%v is the winner \n\n", p2.name)
		}
	} else {
		fmt.Print("draw \n\n")
	}
	return
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

// hash the game board in the player's perspective into an integer
// NOTE: For each player, each location's status is viewed only as occupied either by him/herself or
//       by the opponent, regardless of the actual symbol ("x" or "o") there.
func boardToState(b *board, symbol string) int64 {
	var k, h, v int64
	for _, row := range *b {
		for _, element := range row {
			if element == symbol { // occupied by current player
				v = 0
			} else if element == "" { // empty
				v = 1
			} else { // occupied by opponent
				v = 2
			}
			h += int64(math.Pow(3, float64(k))) * v
			k++
		}
	}
	return h
}

//
func stateToBoard(h int64, symbol string) board {
	// assign opponent's symbol
	var otherSymbol string
	if symbol == "x" {
		otherSymbol = "o"
	} else {
		otherSymbol = "x"
	}

	b := make(board, boardSize)
	k := boardSize*boardSize - 1
	for irow := boardSize - 1; irow >= 0; irow-- {
		r := make([]string, boardSize)
		for ielement := boardSize - 1; ielement >= 0; ielement-- {
			base := int64(math.Pow(3, float64(k)))
			v := h / base
			if v == 0 {
				r[ielement] = symbol
			} else if v == 1 {
				r[ielement] = ""
			} else {
				r[ielement] = otherSymbol
			}
			h -= v * base
			k--
		}
		b[irow] = r
	}
	return b
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

// rowFilled checks whether all elements in a string array are equal to a certain string
func rowFilled(array []string, s string) bool {
	for _, element := range array {
		if element != s {
			return false
		}
	}
	return true
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

// get reward for a certain player by knowing the winner
func getReward(w, s string, d float64) float64 {
	if w == s { // this player wins
		return 1.0
	} else if w == "" { // draw
		return d
	}
	// this player loses
	return 0.0
}
