package main

import "sort"

type stateValue struct {
	state int64
	value float64
}

type stateValues []stateValue

// rankStateValues converts the agent's memory (state-value map) to an array sorted by value
func rankStateValues(vs values) stateValues {
	stateValueArray := make(stateValues, len(vs))
	i := 0
	for s, v := range vs {
		stateValueArray[i] = stateValue{state: s, value: v}
		i++
	}
	sort.Sort(sort.Reverse(stateValueArray))
	return stateValueArray
}

func (vs stateValues) Len() int {
	return len(vs)
}

func (vs stateValues) Less(i, j int) bool {
	return vs[i].value < vs[j].value
}

func (vs stateValues) Swap(i, j int) {
	vs[i], vs[j] = vs[j], vs[i]
}
