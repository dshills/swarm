package brain

import (
	"path/filepath"

	"github.com/dshills/swarm/comm"
	"github.com/dshills/swarm/def"
	"github.com/dshills/swarm/lua"
)

const (
	modelPath  = "models.yaml"
	brainPath  = "brain.yaml"
	layersPath = "layers"
	luaPath    = "lua"
)

func Load(swarmPath string) (Brain, error) {
	// Load the Brain file
	bpath := filepath.Join(swarmPath, brainPath)
	brainDef, err := def.LoadBrainFromPath(bpath)
	if err != nil {
		return nil, err
	}

	for _, l := range brainDef.Layers {
		lpath := filepath.Join(swarmPath, layersPath, l+".yaml")
		ld, err := def.LoadLayerFromPath(lpath)
		if err != nil {
			return nil, err
		}
		brainDef.LayerDefs = append(brainDef.LayerDefs, *ld)
	}

	modPath := filepath.Join(swarmPath, modelPath)
	com := comm.Comms{}
	if err := com.Configure(modPath); err != nil {
		return nil, err
	}

	luaPath := filepath.Join(swarmPath, luaPath)
	fns, err := lua.LoadFromPath(luaPath)
	if err != nil {
		return nil, err
	}

	return newBrain(brainDef, &com, fns)
}
