package neat

type Neuron struct {
	Name        string
	Bias        float64
	Weights     []float64
	Connections []TNeuron

	output  float64
	control bool
}

func NewNeuron() *Neuron {
	return &Neuron{
		Bias:        randWeight(),
		Weights:     []float64{},
		Connections: []TNeuron{},
	}
}

func NewNamedNeuron(name string) *Neuron {
	return &Neuron{
		Name:        name,
		Bias:        randWeight(),
		Weights:     []float64{},
		Connections: []TNeuron{},
	}
}

func (neuron *Neuron) SetName(name string) {
	neuron.Name = name
}

func (neuron *Neuron) GetName() string {
	return neuron.Name
}

func (neuron *Neuron) SetBias(bias float64) {
	neuron.Bias = bias
}

func (neuron *Neuron) GetBias() float64 {
	return neuron.Bias
}

func (neuron *Neuron) SetWeights(weights []float64) {
	neuron.Weights = weights
}

func (neuron *Neuron) GetWeights() []float64 {
	return neuron.Weights
}

func (neuron *Neuron) GetWeight(i int) float64 {
	return neuron.Weights[i]
}

func (neuron *Neuron) GetConnections() []TNeuron {
	return neuron.Connections
}

func (neuron *Neuron) SetConnections(c []TNeuron) {
	neuron.Connections = c
}

// Output methods
func (neuron *Neuron) Reset() {
	neuron.control = false
}

func (neuron *Neuron) SetOutput(output float64) {
	// For the input neurons
	neuron.control = true
	neuron.output = output
}

func (neuron *Neuron) CalOutput(activation Activation) {
	neuron.output = neuron.Bias
	for i := 0; i < len(neuron.Connections); i++ {
		neuron.output += neuron.Connections[i].Output(activation) * neuron.Weights[i]
	}
	neuron.output = activation.Activate(neuron.output)
	neuron.control = true
}

func (neuron *Neuron) Output(activation Activation) float64 {
	if !neuron.control {
		neuron.CalOutput(activation)
	}
	return neuron.output
}

// Connections
func (neuron *Neuron) GetConnectionsIndex(neurons []TNeuron) []int {
	index := make([]int, len(neuron.Connections))
	for i := 0; i < len(neuron.Connections); i++ {
		index[i] = getIndex(neuron.Connections[i], neurons)
	}
	return index
}

func (neuron *Neuron) ConnectionsLength() int {
	return len(neuron.Connections)
}

func (neuron *Neuron) GetConnection(i int) TNeuron {
	return neuron.Connections[i]
}

// Mutations methods
func (neuron *Neuron) HasConnections() bool {
	return len(neuron.Connections) > 0
}

func (neuron *Neuron) find(n TNeuron) int {
	for i := 0; i < len(neuron.Connections); i++ {
		if n == neuron.Connections[i] {
			return i
		}
	}
	return -1
}

func (neuron *Neuron) Mutate(n TNeuron, newNeuronRate, mutSize float64) TNeuron {
	neuron_index := neuron.find(n)
	if neuron_index == -1 {
		// there is not connection to the given neuron
		neuron.Connections = append(neuron.Connections, n)
		neuron.Weights = append(neuron.Weights, randWeight())
	} else {
		// there is a connection to the given neuron
		if probability(newNeuronRate) {
			// add a new neuron in that connection
			newNeuron := NewNeuron()
			newNeuron.Bias = 0.0
			newNeuron.Weights = []float64{1.0}
			newNeuron.Connections = []TNeuron{neuron.Connections[neuron_index]}
			newNeuron.Reset()
			neuron.Connections[neuron_index] = newNeuron
			return newNeuron
		}
		// change the connection weight
		neuron.Weights[neuron_index] += randScaleWeight(mutSize)
	}
	return nil
}

func (neuron *Neuron) MutateBias(mutSize float64) {
	neuron.Bias += randScaleWeight(mutSize)
}

// Serialize
func MakeNeuronFromSerialized(serializedNeuron SerializedNeuron) *Neuron {
	return &Neuron{
		Name:    serializedNeuron.Name,
		Bias:    serializedNeuron.Bias,
		Weights: serializedNeuron.Weights,
	}
}

func (neuron *Neuron) SetConnectionsFromIndex(neurons []TNeuron, index []int) {
	connectionsLen := len(index)
	neuron.Connections = make([]TNeuron, connectionsLen)
	for i := 0; i < connectionsLen; i++ {
		neuron.Connections[i] = neurons[index[i]]
	}
}

func (neuron *Neuron) GetSerializedNeuron(neurons []TNeuron) SerializedNeuron {
	return SerializedNeuron{
		Name:        neuron.Name,
		Bias:        neuron.Bias,
		Weights:     neuron.Weights,
		Connections: neuron.GetConnectionsIndex(neurons),
	}
}
