package neat

type TNeuron interface {
	GetName() string
	SetName(string)

	GetBias() float64
	SetBias(float64)

	GetWeight(int) float64
	GetWeights() []float64
	SetWeights([]float64)

	ConnectionsLength() int
	GetConnection(int) TNeuron
	GetConnections() []TNeuron
	SetConnections([]TNeuron)

	Reset()
	SetOutput(float64)
	CalOutput(Activation)
	Output(Activation) float64

	GetConnectionsIndex([]TNeuron) []int
	HasConnections() bool

	Mutate(TNeuron, float64, float64) TNeuron
	MutateBias(float64)

	SetConnectionsFromIndex([]TNeuron, []int)
	GetSerializedNeuron([]TNeuron) SerializedNeuron
}
