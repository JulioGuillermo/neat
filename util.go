package neat

import (
	"fmt"
	"math/rand"
)

// Check loops
func checkLoops(neuron, from TNeuron) bool {
	if neuron == from {
		return true
	}
	for i := 0; i < from.ConnectionsLength(); i++ {
		if checkLoops(neuron, from.GetConnection(i)) {
			return true
		}
	}
	return false
}

// Random
//func random() float64 {
//	return rand.Float64()
//}

func randWeight() float64 {
	return rand.Float64()*2 - 1
}

func randScaleWeight(scale float64) float64 {
	return randWeight() * scale
}

func probability(pro float64) bool {
	return rand.Float64() < pro
}

func randIndex(from, to int) int {
	return rand.Intn(to-from) + from
}

func randInt(max int) int {
	return rand.Intn(max)
}

// Neurons name
func inputName(i int) string {
	return "I-" + fmt.Sprint(i)
}

func outputName(i int) string {
	return "O-" + fmt.Sprint(i)
}

// Array
func getIndex(neuron TNeuron, neurons []TNeuron) int {
	for i := 0; i < len(neurons); i++ {
		if neuron == neurons[i] {
			return i
		}
	}
	return -1
}
