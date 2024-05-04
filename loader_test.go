package brain_test

import (
	"path/filepath"
	"testing"

	brain "github.com/dshills/swarm"
)

func TestLoad(t *testing.T) {
	fp := filepath.Join("examples", "swarmtest")
	_, err := brain.Load(fp)
	if err != nil {
		t.Fatal(err)
	}
}
