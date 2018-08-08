package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
)

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

	var s string // the "print out" of the board
	for state, valueHistory := range p.mind.valhist {

		b := stateToBoard(state, "x") // use "x" as self to express board
		s = s + strconv.FormatInt(state, 10) + "\n" + printBoard(&b, false) + "\n"

		for time, value := range valueHistory {
			if math.Mod(float64(time), float64(nPrintHistory)) == 0.0 {
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
	}

	d := []byte(s)
	ioutil.WriteFile(filename2, d, 0644)

	fmt.Printf("%v's value histories of the oldest %v states saved into %v \n", p.name, len(p.mind.valhist), filename)
	return
}
