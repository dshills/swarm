package lua

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

func TestCallOneToMany(t *testing.T) {
	script := `
	function oneToMany(str)
			local result = {}
			for word in str:gmatch("%w+") do
				table.insert(result, word)
			end
			return result
		end
	`
	fname := "oneToMany"
	param := "This is a test of the emergency broadcast system."

	strs, err := CallOneToMany(script, fname, param)
	if err != nil {
		t.Error(err)
	}
	exp := 9
	if len(strs) != exp {
		t.Errorf("Expected %v strings got %v", exp, len(strs))
	}
}

func TestLoadFromPath(t *testing.T) {
	path := filepath.Join("..", "examples", "swarmtest", "lua")
	apath, err := filepath.Abs(path)
	if err != nil {
		t.Fatal(err)
	}
	fns, err := LoadFromPath(apath)
	if err != nil {
		t.Fatal(err)
	}
	if len(fns) == 0 {
		t.Errorf("Expected loaded functions got none")
	}
}

func TestCallOneToOne(t *testing.T) {
	script := `
function oneToOne(str)
  str = string.gsub(str, "walk", "jump")
	return str
end
	`
	fname := "oneToOne"
	param := "Immediatly walk to your nearest shelter without delay"

	str, err := CallOneToOne(script, fname, param)
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(str, "jump") {
		fmt.Printf("Expected jump got %v", str)
	}
}
