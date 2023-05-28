# NEAT

This is a neuro-evolution of augmenting topologies library.
It uses a genetic algorithm to evolve neural networks.

This is useful when you don't have a dataset to train your neural network,
for example when you need an agent to interact with an environment or
to learn to play some games.

It's similar to [StaticNeuroGenetic](https://github.com/JulioGuillermo/static_neuro_genetic) library,
but this can also evolve the topology of the neural network.

## How to use

Create an evaluation function

```go
// Eval an individual
func eval(agents *staticneurogenetic.SNG, individual int) {
    inputs := [][]float64 {
        []float64 {0, 0},
        []float64 {0, 1},
        []float64 {1, 0},
        []float64 {1, 1},
    }
    targets := []float64 {
        1,
        0,
        0,
        1
    }

    agents.SetFitness(individual, 0)
    for i, input := range inputs {
        // Get individual output ([]float64)
        output := agents.Output(individual, input)
        // Calculate how wrong is the output
        dif := math.abs(targets[i] - output[0])
        // Added to the fitness
        agents.AddFitness(individual, 1 - dif)
    }
}

// Eval each individual
func evalAll(agents *staticneurogenetic.SNG) {
    for i := range agents.Population {
        eval(agents, i)
    }
}
```

Create a new set of agents

```go
agents := staticneurogenetic.NewSNG(
    2,                              //Input size
    1,                              //Output size
    300,                            //PopulationSize (number of individual to work with)
    10,                             //Survivors (number of individual that will not change in next generation and to use as parents)
    0.1,                            //Probability to mutate a new individual
    0.1,                            //Maximun size of mutations
    0.2,                            //Probability to change bias
    0.2,                            //Probability to add new neurons
    neat.GetActivation("sigmoid"),  //Activation function for the neural network
    false,                          //Not recurrent
)
```

To train the agents we just need to get the next generation

```go
for i := 0; i < 300; i++ {
    evalAll(agents)
    agents.NextGeneration() //Evolve each neural networks
}
```
