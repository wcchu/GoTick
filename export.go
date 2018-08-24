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
func exportValues(name string, values stateValues) {
	filename := name + ".values.csv"
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for state, value := range values {
		row := []string{strconv.FormatInt(state, 10), strconv.FormatFloat(value, 'g', 5, 64)}
		err := writer.Write(row)
		if err != nil {
			log.Fatal("Cannot write to file", err)
		}
	}
	fmt.Printf("%v has %v state-values, saved into %v \n", name, len(values), filename)
	return
}

// write state values of the player to a csv file
func exportValueHistory(name string, vhist stateValueHistory) {
	filename := name + ".oldest_states_hist.csv"
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	filename2 := name + ".oldest_states.txt"

	var s string // the "print out" of the board
	for state, valueHistory := range vhist {

		b, sym := stateToBoard(state)
		s = s + strconv.FormatInt(state, 10) + "\n" + "player plays " + sym + "\n" + printBoard(&b, false) + "\n"

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

	fmt.Printf("%v's value histories of the oldest %v states saved into %v \n", name, len(vhist), filename)
	return
}
