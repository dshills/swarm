package brain_test

import (
	"fmt"
	"path/filepath"
	"testing"

	brain "github.com/dshills/swarm"
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
