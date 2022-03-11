package neat

type Individual struct {
	Recurrent  bool
	InputSize  int
	OutputSize int

	Fitness float64

	Neurons []TNeuron
}

func MakeIndividual(InputSize, OutputSize int, Recurrent bool) *Individual {
	neurons := []TNeuron{}

	if Recurrent {
		// Make input neurons
		// They are the firsts neurons in the neurons array
		for i := 0; i < InputSize; i++ {
			neurons = append(neurons, MakeNamedRecurrentNeuron(inputName(i)))
		}
		// Make output neurons
		// They are the firsts neurons after the inputs neurons in the neurons array
		for i := 0; i < OutputSize; i++ {
			neurons = append(neurons, MakeNamedRecurrentNeuron(outputName(i)))
		}
	} else {
		// Make input neurons
		// They are the firsts neurons in the neurons array
		for i := 0; i < InputSize; i++ {
			neurons = append(neurons, MakeNamedNeuron(inputName(i)))
		}
		// Make output neurons
		// They are the firsts neurons after the inputs neurons in the neurons array
		for i := 0; i < OutputSize; i++ {
			neurons = append(neurons, MakeNamedNeuron(outputName(i)))
		}
	}

	return &Individual{
		Recurrent:  Recurrent,
		InputSize:  InputSize,
		OutputSize: OutputSize,

		Fitness: 0,

		Neurons: neurons,
	}
}

func (individual *Individual) Output(activation Activation, input []float64) []float64 {
	// Put the inputs into the input neurons
	for i := 0; i < individual.InputSize; i++ {
		individual.Neurons[i].SetOutput(input[i])
	}
	// Reset the others neurons
	for i := individual.InputSize; i < len(individual.Neurons); i++ {
		individual.Neurons[i].Reset()
	}

	// Get the output
	output := make([]float64, individual.OutputSize)
	for i := individual.InputSize; i < individual.InputSize+individual.OutputSize; i++ {
		output[i-individual.InputSize] = individual.Neurons[i].Output(activation)
	}
	return output
}

func (individual *Individual) Mutate(newNeuronRate, changeBiasRate, mutSize float64) {
	// Select any neuron to mutate except the inputs neurons
	neuron_to_mutate := individual.Neurons[randIndex(individual.InputSize, len(individual.Neurons))]
	if probability(changeBiasRate) && neuron_to_mutate.HasConnections() {
		neuron_to_mutate.MutateBias(mutSize)
	} else {
		// Select any neuron to connect except the outputs neurons
		neuron := individual.Neurons[(randIndex(individual.OutputSize, len(individual.Neurons))+individual.InputSize)%len(individual.Neurons)]
		// "neuron_to_mutate" will has a connection to "neuron"
		// to prevent a connections loop the selected "neuron" could not has a connection to "neuron_to_mutate"
		for checkLoops(neuron_to_mutate, neuron) {
			neuron = individual.Neurons[(randIndex(individual.OutputSize, len(individual.Neurons))+individual.InputSize)%len(individual.Neurons)]
		}

		// Mutate the selected neuron to the other one
		newNeuron := neuron_to_mutate.Mutate(neuron, newNeuronRate, mutSize)
		// if the newNeuron is not nil the mutation create a new neuron between the other two
		if newNeuron != nil {
			individual.Neurons = append(individual.Neurons, newNeuron)
		}
	}
}

// Fitness
func (individual *Individual) GetFitness() float64 {
	return individual.Fitness
}

func (individual *Individual) SetFitness(fitness float64) {
	individual.Fitness = fitness
}

func (individual *Individual) AddFitness(fitness float64) {
	individual.Fitness += fitness
}

// Serialization
func MakeIndividualFromJson(jsonIndividual JsonIndividual, inputSize, outputSize int) *Individual {
	neuronsLen := len(jsonIndividual.Neurons)
	neurons := make([]TNeuron, neuronsLen)

	for i := 0; i < neuronsLen; i++ {
		neurons[i] = MakeNeuronFromJson(jsonIndividual.Neurons[i])
	}
	for i := 0; i < neuronsLen; i++ {
		neurons[i].SetConnectionsFromIndex(neurons, jsonIndividual.Neurons[i].Connections)
	}

	return &Individual{
		Recurrent:  jsonIndividual.Recurrent,
		InputSize:  inputSize,
		OutputSize: outputSize,
		Fitness:    jsonIndividual.Fitness,
		Neurons:    neurons,
	}
}

func (individual *Individual) GetJsonIndividual() JsonIndividual {
	neuronsLen := len(individual.Neurons)
	neurons := make([]JsonNeuron, neuronsLen)

	for i := 0; i < neuronsLen; i++ {
		neurons[i] = individual.Neurons[i].GetJsonNeuron(individual.Neurons)
	}

	return JsonIndividual{
		Recurrent: individual.Recurrent,
		Fitness:   individual.Fitness,
		Neurons:   neurons,
	}
}
