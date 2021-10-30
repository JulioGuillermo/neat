package neat

import (
	"sort"
)

type NEAT struct {
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
	activation Activation) *NEAT {
	population := make([]*Individual, PopulationSize)

	for i := 0; i < PopulationSize; i++ {
		population[i] = MakeIndividual(InputSize, OutputSize)
	}

	return &NEAT{
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
func (self *NEAT) cross(father, mother *Individual) *Individual {
	baby := MakeIndividual(self.InputSize, self.OutputSize)

	// Create baby neurons
	neuronsLen := len(father.Neurons)
	baby.Neurons = make([]*Neuron, neuronsLen)
	neuronsRef := make([]bool, neuronsLen)
	for i := 0; i < neuronsLen; i++ {
		baby.Neurons[i] = MakeNamedNeuron(father.Neurons[i].Name)
	}

	// Clone parent neurons to baby
	var neuron, fatherNeuron, motherNeuron *Neuron
	var fatherConnections, motherConnections []int
	var connectionsLen int
	for i := 0; i < neuronsLen; i++ {
		neuron = baby.Neurons[i]

		fatherNeuron = father.Neurons[i]
		fatherConnections = fatherNeuron.GetConnectionsIndex(father.Neurons)
		connectionsLen = len(fatherNeuron.Connections)

		if i < len(mother.Neurons) {
			motherNeuron = mother.Neurons[i]
			motherConnections = motherNeuron.GetConnectionsIndex(mother.Neurons)
		} else {
			motherNeuron = fatherNeuron
			motherConnections = motherNeuron.GetConnectionsIndex(father.Neurons)
		}

		// Set bias
		if probability(0.5) {
			neuron.Bias = fatherNeuron.Bias
		} else {
			neuron.Bias = motherNeuron.Bias
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
					neuron.Weights = append(neuron.Weights, motherNeuron.Weights[j])
					neuron.Connections = append(neuron.Connections, baby.Neurons[motherConnections[j]])
					neuronsRef[motherConnections[j]] = true
				} else if !checkLoops(neuron, baby.Neurons[fatherConnections[j]]) {
					//neuron.Weights[j] = fatherNeuron.Weights[j]
					//neuron.Connections[j] = baby.Neurons[fatherConnections[j]]
					neuron.Weights = append(neuron.Weights, fatherNeuron.Weights[j])
					neuron.Connections = append(neuron.Connections, baby.Neurons[fatherConnections[j]])
					neuronsRef[fatherConnections[j]] = true
				}
			} else if !checkLoops(neuron, baby.Neurons[fatherConnections[j]]) {
				//neuron.Weights[j] = fatherNeuron.Weights[j]
				//neuron.Connections[j] = baby.Neurons[fatherConnections[j]]
				neuron.Weights = append(neuron.Weights, fatherNeuron.Weights[j])
				neuron.Connections = append(neuron.Connections, baby.Neurons[fatherConnections[j]])
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

	if probability(self.MutRate) {
		baby.Mutate(self.NewNeuronRate, self.ChangeBiasRate, self.MutSize)
	}

	return baby
}

func (self *NEAT) NextGeneration() {
	self.Generation++

	// Sort population by fitness
	sort.SliceStable(self.Population, func(x, y int) bool {
		return self.Population[x].GetFitness() > self.Population[y].GetFitness()
	})

	// Reset fitness
	for i := 0; i < self.PopulationSize; i++ {
		self.Population[i].SetFitness(0)
	}

	// Create the new generation
	var father, mother int
	for i := self.Survivors; i < self.PopulationSize; i++ {
		// Select father and mother from survivors
		father = randInt(self.Survivors)
		mother = randInt(self.Survivors)
		for father == mother {
			mother = randInt(self.Survivors)
		}
		if father > mother {
			father, mother = mother, father
		}

		// Create the new individual
		self.Population[i] = self.cross(self.Population[father], self.Population[mother])
	}
}

// Output
func (self *NEAT) Output(index int, input []float64) []float64 {
	return self.Population[index].Output(self.Activation, input)
}

// Fitness
func (self *NEAT) GetFitness(index int) float64 {
	return self.Population[index].GetFitness()
}

func (self *NEAT) SetFitness(index int, fitness float64) {
	self.Population[index].SetFitness(fitness)
}

func (self *NEAT) AddFitness(index int, fitness float64) {
	self.Population[index].AddFitness(fitness)
}

// Generation
func (self *NEAT) GetGeneration() int {
	return self.Generation
}

// Serialization
type JsonNEAT struct {
	InputSize  int `json:"input_size"`
	OutputSize int `json:"output_size"`

	PopulationSize int `json:"population_size"`
	Survivors      int `json:"survivors"`

	MutRate        float64 `json:"mutation_rate"`
	MutSize        float64 `json:"mutation_size"`
	ChangeBiasRate float64 `json:"change_bias_rate"`
	NewNeuronRate  float64 `json:"new_neuron_rate"`

	Activation string `json:"activation"`

	Generation int              `json:"generation"`
	Population []JsonIndividual `json:"population"`
}

func (self *NEAT) GetJsonNEAT() JsonNEAT {
	population := make([]JsonIndividual, self.PopulationSize)
	for i := 0; i < self.PopulationSize; i++ {
		population[i] = self.Population[i].GetJsonIndividual()
	}
	return JsonNEAT{
		InputSize:  self.InputSize,
		OutputSize: self.OutputSize,

		PopulationSize: self.PopulationSize,
		Survivors:      self.Survivors,

		MutRate:        self.MutRate,
		MutSize:        self.MutSize,
		ChangeBiasRate: self.ChangeBiasRate,
		NewNeuronRate:  self.NewNeuronRate,

		Activation: self.Activation.GetString(),

		Generation: self.Generation,
		Population: population,
	}
}

func MakeNEATFromJsonNEAT(jsonNeat JsonNEAT) *NEAT {
	population := make([]*Individual, jsonNeat.PopulationSize)
	for i := 0; i < jsonNeat.PopulationSize; i++ {
		population[i] = MakeIndividualFromJson(jsonNeat.Population[i], jsonNeat.InputSize, jsonNeat.OutputSize)
	}

	return &NEAT{
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

func (self *NEAT) Save(path string) error {
	return Save(self, path)
}
