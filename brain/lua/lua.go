package lua

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

type Function interface {
	Name() string
	Script() string
	OneToOne(string) (string, error)
	OneToMany(string) ([]string, error)
}

type _func struct {
	name   string
	script string
}

func (f *_func) OneToOne(str string) (string, error) {
	return CallOneToOne(f.script, f.name, str)
}

func (f *_func) OneToMany(str string) ([]string, error) {
	return CallOneToMany(f.script, f.name, str)
}

func (f *_func) Name() string {
	return f.name
}

func (f *_func) Script() string {
	return f.script
}

func LoadFromPath(path string) ([]Function, error) {
	funcs := []Function{}

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("loadLuaFromPath: %s %w", path, err)
	}

	for _, finfo := range files {
		if filepath.Ext(finfo.Name()) == ".lua" {
			fpath := filepath.Join(path, finfo.Name())
			byts, err := os.ReadFile(fpath)
			if err != nil {
				return nil, fmt.Errorf("loadLuaFromPath: %s %w", fpath, err)
			}
			lf := _func{name: strings.TrimSuffix(finfo.Name(), ".lua"), script: string(byts)}
			funcs = append(funcs, &lf)
		}
	}
	return funcs, nil
}

func CallOneToOne(script, fName, p string) (string, error) {
	L := lua.NewState()
	defer L.Close()

	// Load the Lua script
	if err := L.DoString(script); err != nil {
		return "", err
	}

	// Call function with params
	ops := lua.P{
		Fn:      L.GetGlobal(fName),
		NRet:    1,
		Protect: true,
	}
	if err := L.CallByParam(ops, lua.LString(p)); err != nil {
		return "", err
	}

	luaStr := L.Get(-1)
	if luaStr.Type() != lua.LTString {
		return "", fmt.Errorf("return %v but a string was expected", luaStr.Type())
	}
	str := luaStr.(lua.LString)
	return str.String(), nil
}

func CallOneToMany(script, fName, p string) ([]string, error) {
	L := lua.NewState()
	defer L.Close()

	// Load the Lua script
	if err := L.DoString(script); err != nil {
		return nil, err
	}

	// Call function with params
	ops := lua.P{
		Fn:      L.GetGlobal(fName),
		NRet:    1,
		Protect: true,
	}
	if err := L.CallByParam(ops, lua.LString(p)); err != nil {
		return nil, err
	}

	luaTable := L.Get(-1)
	if luaTable.Type() != lua.LTTable {
		return nil, fmt.Errorf("return is not a table")
	}
	tbl := luaTable.(*lua.LTable)

	var strings []string
	tbl.ForEach(func(key lua.LValue, val lua.LValue) {
		strings = append(strings, val.String())
	})

	return strings, nil
}
