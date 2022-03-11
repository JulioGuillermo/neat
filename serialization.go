package neat

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Serializers
type JsonNeuron struct {
	Name        string    `json:"name"`
	Bias        float64   `json:"bias"`
	RWeight     float64   `json:"r_weight"`
	Weights     []float64 `json:"weights"`
	Connections []int     `json:"connections"`
}

type JsonIndividual struct {
	Recurrent bool         `json:"recurrent"`
	Fitness   float64      `json:"fitness"`
	Neurons   []JsonNeuron `json:"neurons"`
}

type JsonNEAT struct {
	Recurrent  bool `json:"recurrent"`
	InputSize  int  `json:"input_size"`
	OutputSize int  `json:"output_size"`

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

func Save(neat *NEAT, path string) error {
	bytes, err := json.Marshal(neat.GetJsonNEAT())
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func Load(path string) (*NEAT, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(file)
	if err != nil {
		return nil, err
	}

	var jsonNeat JsonNEAT
	err = decoder.Decode(&jsonNeat)
	if err != nil {
		return nil, err
	}

	return MakeNEATFromJsonNEAT(jsonNeat), nil
}
