package neat

import (
	"math"
	"strings"
)

type Activation interface {
	Activate(float64) float64
	GetString() string
}

// Linear
type Linear struct{}

func (act *Linear) Activate(n float64) float64 {
	return n
}
func (act *Linear) GetString() string {
	return "linear"
}

// Relu
type Relu struct{}

func (act *Relu) Activate(n float64) float64 {
	if n > 0 {
		return n
	}
	return 0
}
func (act *Relu) GetString() string {
	return "relu"
}

// Sigmoid
type Sigmoid struct{}

func (act *Sigmoid) Activate(n float64) float64 {
	return 1.0 / (1.0 + math.Exp(-n))
}
func (act *Sigmoid) GetString() string {
	return "sigmoid"
}

// Tanh
type Tanh struct{}

func (act *Tanh) Activate(n float64) float64 {
	return math.Tanh(n)
}
func (act *Tanh) GetString() string {
	return "tanh"
}

// Sin
type Sin struct{}

func (act *Sin) Activate(n float64) float64 {
	return math.Sin(n)
}
func (act *Sin) GetString() string {
	return "sin"
}

type Sig struct{}

func (act *Sig) Activate(n float64) float64 {
	if n > 0 {
		return 1
	}
	if n < 0 {
		return -1
	}
	return 0
}
func (act *Sig) GetString() string {
	return "sig"
}

func GetActivation(name string) Activation {
	switch strings.ToLower(name) {
	case "linear":
		return &Linear{}
	case "sigmoid":
		return &Sigmoid{}
	case "tanh":
		return &Tanh{}
	case "sin":
		return &Sin{}
	case "sig":
		return &Sig{}
	}
	return &Relu{}
}
