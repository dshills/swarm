package brain_test

import (
	"path/filepath"
	"testing"

	"github.com/dshills/swarm/brain"
)

func TestLoad(t *testing.T) {
	fp := filepath.Join("examples", "swarmtest")
	_, err := brain.Load(fp)
	if err != nil {
		t.Fatal(err)
	}
}
