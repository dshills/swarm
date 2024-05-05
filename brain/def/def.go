package def

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Definition struct {
	LayerDefs []LayerDef
	Brain     string   `yaml:"Brain"`
	Layers    []string `yaml:"Layers"`
}

func LoadBrainFromPath(path string) (*Definition, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("loadBrainFromPath: %s %w", path, err)
	}
	defer file.Close()

	def := Definition{}
	if err := yaml.NewDecoder(file).Decode(&def); err != nil {
		return nil, fmt.Errorf("loadBrainFromPath: %s %w", path, err)
	}

	return &def, nil
}

type LayerDef struct {
	Persona          string   `yaml:"Persona"`
	Prompt           string   `yaml:"Prompt"`
	NeuronModels     []string `yaml:"NeuronModels"`
	IgnoreContext    bool     `yaml:"IgnoreContext"`
	ChangeTaskFns    []string `yaml:"ChangeTaskFns"`
	ChangeResultFns  []string `yaml:"ChangeResultFns"`
	ResultToTasksFns []string `yaml:"ResultToTasksFns"`
}

func LoadLayerFromPath(path string) (*LayerDef, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("loadLayerFromPath: %s %w", path, err)
	}
	ld := LayerDef{}
	if err := yaml.NewDecoder(file).Decode(&ld); err != nil {
		return nil, fmt.Errorf("loadLayerFromPath: %s %w", path, err)
	}
	return &ld, nil
}
