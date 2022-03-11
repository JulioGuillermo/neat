package neat

import (
	"sort"
)

type NEAT struct {
	Recurrent  bool
	InputSize  int
	OutputSize int

	PopulationSize int
	Survivors      int

	MutRate        float64
	MutSize        float64
	ChangeBiasRate float64
	NewNeuronRate  float64

	Activation Activation

	Generation int
	Population []*Individual
}

// Constructor
func MakeNEAT(InputSize, OutputSize, PopulationSize, Survivors int,
	MutRate, MutSize, ChangeBiasRate, NewNeuronRate float64,
	activation Activation, Recurrent bool) *NEAT {
	population := make([]*Individual, PopulationSize)

	for i := 0; i < PopulationSize; i++ {
		population[i] = MakeIndividual(InputSize, OutputSize, Recurrent)
	}

	return &NEAT{
		Recurrent:      Recurrent,
		InputSize:      InputSize,
		OutputSize:     OutputSize,
		PopulationSize: PopulationSize,
		Survivors:      Survivors,

		MutRate:        MutRate,
		MutSize:        MutSize,
		ChangeBiasRate: ChangeBiasRate,
		NewNeuronRate:  NewNeuronRate,

		Activation: activation,

		Generation: 0,
		Population: population,
	}
}

// Evolution
func (neat *NEAT) cross(father, mother *Individual) *Individual {
	baby := MakeIndividual(neat.InputSize, neat.OutputSize, neat.Recurrent)

	// Create baby neurons
	neuronsLen := len(father.Neurons)
	baby.Neurons = make([]TNeuron, neuronsLen)
	neuronsRef := make([]bool, neuronsLen)
	for i := 0; i < neuronsLen; i++ {
		baby.Neurons[i] = MakeNamedNeuron(father.Neurons[i].GetName())
	}

	// Clone parent neurons to baby
	var neuron, fatherNeuron, motherNeuron TNeuron
	var fatherConnections, motherConnections []int
	var connectionsLen int
	for i := 0; i < neuronsLen; i++ {
		neuron = baby.Neurons[i]

		fatherNeuron = father.Neurons[i]
		fatherConnections = fatherNeuron.GetConnectionsIndex(father.Neurons)
		connectionsLen = fatherNeuron.ConnectionsLength()

		if i < len(mother.Neurons) {
			motherNeuron = mother.Neurons[i]
			motherConnections = motherNeuron.GetConnectionsIndex(mother.Neurons)
		} else {
			motherNeuron = fatherNeuron
			motherConnections = motherNeuron.GetConnectionsIndex(father.Neurons)
		}

		// Set bias
		if probability(0.5) {
			neuron.SetBias(fatherNeuron.GetBias())
		} else {
			neuron.SetBias(motherNeuron.GetBias())
		}

		// Set connections
		//neuron.Weights = make([]float64, connectionsLen)
		//neuron.Connections = make([]*Neuron, connectionsLen)
		for j := 0; j < connectionsLen; j++ {
			/*if probability(0.5) && j < len(motherConnections) && motherConnections[j] < neuronsLen {
			      neuron.Weights[j] = motherNeuron.Weights[j]
			      neuron.Connections[j] = baby.Neurons[motherConnections[j]]
			      neuronsRef[motherConnections[j]] = true
			  } else {
			      neuron.Weights[j] = fatherNeuron.Weights[j]
			      neuron.Connections[j] = baby.Neurons[fatherConnections[j]]
			      neuronsRef[fatherConnections[j]] = true
			  }*/
			// Add loop check
			if j < len(motherConnections) && motherConnections[j] < neuronsLen {
				if probability(0.5) && !checkLoops(neuron, baby.Neurons[motherConnections[j]]) {
					//neuron.Weights[j] = motherNeuron.Weights[j]
					//neuron.Connections[j] = baby.Neurons[motherConnections[j]]
					neuron.SetWeights(append(neuron.GetWeights(), motherNeuron.GetWeight(j)))
					neuron.SetConnections(append(neuron.GetConnections(), baby.Neurons[motherConnections[j]]))
					neuronsRef[motherConnections[j]] = true
				} else if !checkLoops(neuron, baby.Neurons[fatherConnections[j]]) {
					//neuron.Weights[j] = fatherNeuron.Weights[j]
					//neuron.Connections[j] = baby.Neurons[fatherConnections[j]]
					neuron.SetWeights(append(neuron.GetWeights(), fatherNeuron.GetWeight(j)))
					neuron.SetConnections(append(neuron.GetConnections(), baby.Neurons[fatherConnections[j]]))
					neuronsRef[fatherConnections[j]] = true
				}
			} else if !checkLoops(neuron, baby.Neurons[fatherConnections[j]]) {
				//neuron.Weights[j] = fatherNeuron.Weights[j]
				//neuron.Connections[j] = baby.Neurons[fatherConnections[j]]
				neuron.SetWeights(append(neuron.GetWeights(), fatherNeuron.GetWeight(j)))
				neuron.SetConnections(append(neuron.GetConnections(), baby.Neurons[fatherConnections[j]]))
				neuronsRef[fatherConnections[j]] = true
			}
		}
	}

	// Removing unrefered neurons
	neurons := baby.Neurons[:baby.InputSize+baby.OutputSize]
	for i := baby.InputSize + baby.OutputSize; i < neuronsLen; i++ {
		if neuronsRef[i] {
			neurons = append(neurons, baby.Neurons[i])
		}
	}
	baby.Neurons = neurons

	if probability(neat.MutRate) {
		baby.Mutate(neat.NewNeuronRate, neat.ChangeBiasRate, neat.MutSize)
	}

	return baby
}

func (neat *NEAT) NextGeneration() {
	neat.Generation++

	// Sort population by fitness
	sort.SliceStable(neat.Population, func(x, y int) bool {
		return neat.Population[x].GetFitness() > neat.Population[y].GetFitness()
	})

	// Reset fitness
	for i := 0; i < neat.PopulationSize; i++ {
		neat.Population[i].SetFitness(0)
	}

	// Create the new generation
	var father, mother int
	for i := neat.Survivors; i < neat.PopulationSize; i++ {
		// Select father and mother from survivors
		father = randInt(neat.Survivors)
		mother = randInt(neat.Survivors)
		for father == mother {
			mother = randInt(neat.Survivors)
		}
		if father > mother {
			father, mother = mother, father
		}

		// Create the new individual
		neat.Population[i] = neat.cross(neat.Population[father], neat.Population[mother])
	}
}

// Output
func (neat *NEAT) Output(index int, input []float64) []float64 {
	return neat.Population[index].Output(neat.Activation, input)
}

// Fitness
func (neat *NEAT) GetFitness(index int) float64 {
	return neat.Population[index].GetFitness()
}

func (neat *NEAT) SetFitness(index int, fitness float64) {
	neat.Population[index].SetFitness(fitness)
}

func (neat *NEAT) AddFitness(index int, fitness float64) {
	neat.Population[index].AddFitness(fitness)
}

// Generation
func (neat *NEAT) GetGeneration() int {
	return neat.Generation
}

// Serialization
func (neat *NEAT) GetJsonNEAT() JsonNEAT {
	population := make([]JsonIndividual, neat.PopulationSize)
	for i := 0; i < neat.PopulationSize; i++ {
		population[i] = neat.Population[i].GetJsonIndividual()
	}
	return JsonNEAT{
		Recurrent:  neat.Recurrent,
		InputSize:  neat.InputSize,
		OutputSize: neat.OutputSize,

		PopulationSize: neat.PopulationSize,
		Survivors:      neat.Survivors,

		MutRate:        neat.MutRate,
		MutSize:        neat.MutSize,
		ChangeBiasRate: neat.ChangeBiasRate,
		NewNeuronRate:  neat.NewNeuronRate,

		Activation: neat.Activation.GetString(),

		Generation: neat.Generation,
		Population: population,
	}
}

func MakeNEATFromJsonNEAT(jsonNeat JsonNEAT) *NEAT {
	population := make([]*Individual, jsonNeat.PopulationSize)
	for i := 0; i < jsonNeat.PopulationSize; i++ {
		population[i] = MakeIndividualFromJson(jsonNeat.Population[i], jsonNeat.InputSize, jsonNeat.OutputSize)
	}

	return &NEAT{
		Recurrent:  jsonNeat.Recurrent,
		InputSize:  jsonNeat.InputSize,
		OutputSize: jsonNeat.OutputSize,

		PopulationSize: jsonNeat.PopulationSize,
		Survivors:      jsonNeat.Survivors,

		MutRate:        jsonNeat.MutRate,
		MutSize:        jsonNeat.MutSize,
		ChangeBiasRate: jsonNeat.ChangeBiasRate,
		NewNeuronRate:  jsonNeat.NewNeuronRate,

		Activation: GetActivation(jsonNeat.Activation),

		Generation: jsonNeat.Generation,
		Population: population,
	}
}

func (neat *NEAT) Save(path string) error {
	return Save(neat, path)
}
