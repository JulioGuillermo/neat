package neat

type RecurrentNeuron struct {
	Name            string
	Bias            float64
	RecurrentWeight float64
	Weights         []float64
	Connections     []TNeuron

	old_output float64
	output     float64
	control    bool
}

func NewRecurrentNeuron() *RecurrentNeuron {
	return &RecurrentNeuron{
		Bias:        randWeight(),
		Weights:     []float64{},
		Connections: []TNeuron{},

		old_output: 0,
	}
}

func NewNamedRecurrentNeuron(name string) *RecurrentNeuron {
	return &RecurrentNeuron{
		Name:        name,
		Bias:        randWeight(),
		Weights:     []float64{},
		Connections: []TNeuron{},

		old_output: 0,
	}
}

func (neuron *RecurrentNeuron) SetName(name string) {
	neuron.Name = name
}

func (neuron *RecurrentNeuron) GetName() string {
	return neuron.Name
}

func (neuron *RecurrentNeuron) SetBias(bias float64) {
	neuron.Bias = bias
}

func (neuron *RecurrentNeuron) GetBias() float64 {
	return neuron.Bias
}

func (neuron *RecurrentNeuron) SetWeights(weights []float64) {
	neuron.Weights = weights
}

func (neuron *RecurrentNeuron) GetWeights() []float64 {
	return neuron.Weights
}

func (neuron *RecurrentNeuron) GetWeight(i int) float64 {
	return neuron.Weights[i]
}

func (neuron *RecurrentNeuron) GetConnections() []TNeuron {
	return neuron.Connections
}

func (neuron *RecurrentNeuron) SetConnections(c []TNeuron) {
	neuron.Connections = c
}

// Output methods
func (neuron *RecurrentNeuron) Reset() {
	neuron.control = false
}

func (neuron *RecurrentNeuron) SetOutput(output float64) {
	// For the input neurons
	neuron.control = true
	neuron.output = output
}

func (neuron *RecurrentNeuron) CalOutput(activation Activation) {
	neuron.output = neuron.Bias
	neuron.output += neuron.RecurrentWeight * neuron.old_output
	for i := 0; i < len(neuron.Connections); i++ {
		neuron.output += neuron.Connections[i].Output(activation) * neuron.Weights[i]
	}
	neuron.output = activation.Activate(neuron.output)
	neuron.control = true
	neuron.old_output = neuron.output
}

func (neuron *RecurrentNeuron) Output(activation Activation) float64 {
	if !neuron.control {
		neuron.CalOutput(activation)
	}
	return neuron.output
}

// Connections
func (neuron *RecurrentNeuron) GetConnectionsIndex(neurons []TNeuron) []int {
	index := make([]int, len(neuron.Connections))
	for i := 0; i < len(neuron.Connections); i++ {
		index[i] = getIndex(neuron.Connections[i], neurons)
	}
	return index
}

func (neuron *RecurrentNeuron) ConnectionsLength() int {
	return len(neuron.Connections)
}

func (neuron *RecurrentNeuron) GetConnection(i int) TNeuron {
	return neuron.Connections[i]
}

// Mutations methods
func (neuron *RecurrentNeuron) HasConnections() bool {
	return len(neuron.Connections) > 0
}

func (neuron *RecurrentNeuron) find(n TNeuron) int {
	for i := 0; i < len(neuron.Connections); i++ {
		if n == neuron.Connections[i] {
			return i
		}
	}
	return -1
}

func (neuron *RecurrentNeuron) Mutate(n TNeuron, newNeuronRate, mutSize float64) TNeuron {
	neuron_index := neuron.find(n)
	if neuron_index == -1 {
		// there is not connection to the given neuron
		neuron.Connections = append(neuron.Connections, n)
		neuron.Weights = append(neuron.Weights, randWeight())
	} else {
		// there is a connection to the given neuron
		if probability(newNeuronRate) {
			// add a new neuron in that connection
			newNeuron := NewRecurrentNeuron()
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

func (neuron *RecurrentNeuron) MutateBias(mutSize float64) {
	if probability(0.5) {
		neuron.RecurrentWeight += randScaleWeight(mutSize)
	} else {
		neuron.Bias += randScaleWeight(mutSize)
	}
}

// Serialize
func MakeRecurrentNeuronFromSerialized(serializedNeuron SerializedNeuron) *RecurrentNeuron {
	return &RecurrentNeuron{
		Name:            serializedNeuron.Name,
		Bias:            serializedNeuron.Bias,
		RecurrentWeight: serializedNeuron.RWeight,
		Weights:         serializedNeuron.Weights,
	}
}

func (neuron *RecurrentNeuron) SetConnectionsFromIndex(neurons []TNeuron, index []int) {
	connectionsLen := len(index)
	neuron.Connections = make([]TNeuron, connectionsLen)
	for i := 0; i < connectionsLen; i++ {
		neuron.Connections[i] = neurons[index[i]]
	}
}

func (neuron *RecurrentNeuron) GetSerializedNeuron(neurons []TNeuron) SerializedNeuron {
	return SerializedNeuron{
		Name:        neuron.Name,
		Bias:        neuron.Bias,
		RWeight:     neuron.RecurrentWeight,
		Weights:     neuron.Weights,
		Connections: neuron.GetConnectionsIndex(neurons),
	}
}
