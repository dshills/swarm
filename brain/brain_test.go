package brain_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/dshills/swarm/brain"
)

func TestThink(t *testing.T) {
	fpath := filepath.Join("examples", "swarmtest")
	br, err := brain.Load(fpath)
	if err != nil {
		t.Fatal(err)
	}

	task := "What is Formula 1"
	out := br.Think(task)

	res := <-out
	for _, r := range res.Results {
		fmt.Println(r)
	}
	fmt.Println(res.String())
}
