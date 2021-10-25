package neat

type Individual struct {
    InputSize       int
    OutputSize      int

    Fitness         float64

    Neurons         []*Neuron
}

func MakeIndividual(InputSize, OutputSize int) *Individual {
    neurons := []*Neuron{}

    // Make input neurons
    // They are the firsts neurons in the neurons array
    for i := 0; i < InputSize; i ++ {
        neurons = append(neurons, MakeNamedNeuron(inputName(i)))
    }
    // Make output neurons
    // They are the firsts neurons after the inputs neurons in the neurons array
    for i := 0; i < OutputSize; i ++ {
        neurons = append(neurons, MakeNamedNeuron(outputName(i)))
    }

    return &Individual {
        InputSize:      InputSize,
        OutputSize:     OutputSize,

        Fitness:        0,

        Neurons:        neurons,
    }
}

func (self *Individual) Output(activation Activation, input []float64) []float64 {
    // Put the inputs into the input neurons
    for i := 0; i < self.InputSize; i ++ {
        self.Neurons[i].SetOutput(input[i])
    }
    // Reset the others neurons
    for i := self.InputSize; i < len(self.Neurons); i ++ {
        self.Neurons[i].Reset()
    }

    // Get the output
    output := make([]float64, self.OutputSize)
    for i := self.InputSize; i < self.InputSize + self.OutputSize; i ++ {
        output[i - self.InputSize] = self.Neurons[i].Output(activation)
    }
    return output
}

func (self *Individual) Mutate(newNeuronRate, changeBiasRate, mutSize float64) {
    // Select any neuron to mutate except the inputs neurons
    neuron_to_mutate := self.Neurons[randIndex(self.InputSize, len(self.Neurons))]
    if probability(changeBiasRate) && neuron_to_mutate.HasConnections() {
        neuron_to_mutate.MutateBias(mutSize)
    } else {
        // Select any neuron to connect except the outputs neurons
        neuron := self.Neurons[(randIndex(self.OutputSize, len(self.Neurons)) + self.InputSize) % len(self.Neurons)]
        // "neuron_to_mutate" will has a connection to "neuron"
        // to prevent a connections loop the selected "neuron" could not has a connection to "neuron_to_mutate"
        for checkLoops(neuron_to_mutate, neuron) {
            neuron = self.Neurons[(randIndex(self.OutputSize, len(self.Neurons)) + self.InputSize) % len(self.Neurons)]
        }

        // Mutate the selected neuron to the other one
        newNeuron := neuron_to_mutate.Mutate(neuron, newNeuronRate, mutSize)
        // if the newNeuron is not nil the mutation create a new neuron between the other two
        if newNeuron != nil {
            self.Neurons = append(self.Neurons, newNeuron)
        }
    }
}

// Fitness
func (self *Individual) GetFitness() float64 {
    return self.Fitness
}

func (self *Individual) SetFitness(fitness float64) {
    self.Fitness = fitness
}

func (self *Individual) AddFitness(fitness float64) {
    self.Fitness += fitness
}

// Serialization
type JsonIndividual struct {
    Neurons     []JsonNeuron    `json:"neurons"`
}

func MakeIndividualFromJson(jsonIndividual JsonIndividual, inputSize, outputSize int) *Individual {
    neuronsLen := len(jsonIndividual.Neurons)
    neurons := make([]*Neuron, neuronsLen)

    for i := 0; i < neuronsLen; i ++ {
        neurons[i] = MakeNeuronFromJson(jsonIndividual.Neurons[i])
    }
    for i := 0; i < neuronsLen; i ++ {
        neurons[i].SetConnectionsFromIndex(neurons, jsonIndividual.Neurons[i].Connections)
    }

    return &Individual {
        InputSize:      inputSize,
        OutputSize:     outputSize,
        Fitness:        0,
        Neurons:        neurons,
    }
}

func (self *Individual) GetJsonIndividual() JsonIndividual {
    neuronsLen := len(self.Neurons)
    neurons := make([]JsonNeuron, neuronsLen)

    for i := 0; i < neuronsLen; i ++ {
        neurons[i] = self.Neurons[i].GetJsonNeuron(self.Neurons)
    }

    return JsonIndividual {
        Neurons:    neurons,
    }
}
