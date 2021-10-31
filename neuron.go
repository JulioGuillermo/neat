package neat

type Neuron struct {
    Name        string
    Bias        float64
    Weights     []float64
    Connections []*Neuron

    output      float64
    control     bool
}

func MakeNeuron() *Neuron {
    return &Neuron {
        Bias:           randWeight(),
        Weights:        []float64{},
        Connections:    []*Neuron{},
    }
}

func MakeNamedNeuron(name string) *Neuron {
    return &Neuron {
        Name:           name,
        Bias:           randWeight(),
        Weights:        []float64{},
        Connections:    []*Neuron{},
    }
}

// Output methods
func (self *Neuron) Reset() {
    self.control = false
}

func (self *Neuron) SetOutput(output float64) {
    // For the input neurons
    self.control = true
    self.output = output
}

func (self *Neuron) CalOutput(activation Activation) {
    self.output = self.Bias
    for i := 0; i < len(self.Connections); i ++ {
        self.output += self.Connections[i].Output(activation) * self.Weights[i]
    }
    self.output = activation.Activate(self.output)
    self.control = true
}

func (self *Neuron) Output(activation Activation) float64 {
    if !self.control {
        self.CalOutput(activation)
    }
    return self.output
}

// Connections
func (self *Neuron) GetConnectionsIndex(neurons []*Neuron) []int {
    index := make([]int, len(self.Connections))
    for i := 0; i < len(self.Connections); i ++ {
        index[i] = getIndex(self.Connections[i], neurons)
    }
    return index
}

// Mutations methods
func (self *Neuron) HasConnections() bool {
    return len(self.Connections) > 0
}

func (self *Neuron) find(neuron *Neuron) int {
    for i := 0; i < len(self.Connections); i ++ {
        if neuron == self.Connections[i] {
            return i
        }
    }
    return -1
}

func (self *Neuron) Mutate(neuron *Neuron, newNeuronRate, mutSize float64) *Neuron {
    neuron_index := self.find(neuron)
    if neuron_index == -1 {
        // there is not connection to the given neuron
        self.Connections = append(self.Connections, neuron)
        self.Weights = append(self.Weights, randWeight())
    } else {
        // there is a connection to the given neuron
        if probability(newNeuronRate) {
            // add a new neuron in that connection
            newNeuron := MakeNeuron()
            newNeuron.Bias = 0.0
            newNeuron.Weights = []float64{1.0}
            newNeuron.Connections = []*Neuron{self.Connections[neuron_index]}
            newNeuron.Reset()
            self.Connections[neuron_index] = newNeuron
            return newNeuron
        }
        // change the connection weight
        self.Weights[neuron_index] += randScaleWeight(mutSize)
    }
    return nil
}

func (self *Neuron) MutateBias(mutSize float64) {
    self.Bias += randScaleWeight(mutSize)
}

// Serialize
type JsonNeuron struct {
    Name        string      `json:"name"`
    Bias        float64     `json:"bias"`
    Weights     []float64   `json:"weights"`
    Connections []int       `json:"connections"`
}

func MakeNeuronFromJson(jsonNeuron JsonNeuron) *Neuron {
    return &Neuron {
        Name:           jsonNeuron.Name,
        Bias:           jsonNeuron.Bias,
        Weights:        jsonNeuron.Weights,
    }
}

func (self *Neuron) SetConnectionsFromIndex(neurons []*Neuron, index []int) {
    connectionsLen := len(index)
    self.Connections = make([]*Neuron, connectionsLen)
    for i := 0; i < connectionsLen; i ++ {
        self.Connections[i] = neurons[index[i]]
    }
}

func (self *Neuron) GetJsonNeuron(neurons []*Neuron) JsonNeuron {
    return JsonNeuron {
        Name:           self.Name,
        Bias:           self.Bias,
        Weights:        self.Weights,
        Connections:    self.GetConnectionsIndex(neurons),
    }
}
