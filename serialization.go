package neat

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"io/ioutil"
	"os"
)

// Serializers
type SerializedNeuron struct {
	Name        string    `json:"name"`
	Bias        float64   `json:"bias"`
	RWeight     float64   `json:"r_weight"`
	Weights     []float64 `json:"weights"`
	Connections []int     `json:"connections"`
}

type SerializedIndividual struct {
	Recurrent bool               `json:"recurrent"`
	Fitness   float64            `json:"fitness"`
	Neurons   []SerializedNeuron `json:"neurons"`
}

type SerializedNEAT struct {
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

	Generation int                    `json:"generation"`
	Population []SerializedIndividual `json:"population"`
}

func SaveAsJson(neat *NEAT, path string) error {
	bytes, err := json.Marshal(neat.GetSerializedNEAT())
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func LoadFromJson(path string) (*NEAT, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(file)
	if err != nil {
		return nil, err
	}

	var serializedNeat SerializedNEAT
	err = decoder.Decode(&serializedNeat)
	if err != nil {
		return nil, err
	}

	return MakeNEATFromSerializedNEAT(serializedNeat), nil
}

func SaveAsBin(neat *NEAT, path string) error {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(neat.GetSerializedNEAT())
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, buff.Bytes(), 0777)
	if err != nil {
		return err
	}

	return nil
}

func LoadFromBin(path string) (*NEAT, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	decoder := gob.NewDecoder(file)
	if err != nil {
		return nil, err
	}

	var serializedNeat SerializedNEAT
	err = decoder.Decode(&serializedNeat)
	if err != nil {
		return nil, err
	}

	return MakeNEATFromSerializedNEAT(serializedNeat), nil
}
