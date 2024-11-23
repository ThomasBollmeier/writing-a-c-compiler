package frontend

import (
	"errors"
	"fmt"
)

type Environment struct {
	parent   *Environment
	identMap map[string]EnvEntry
}

type identCategory int

const (
	idCatVariable identCategory = iota
	idCatFunction
	idCatParameter
)

type EnvEntry struct {
	uniqueName string
	isExternal bool
	category   identCategory
	typeInfo   TypeInfo
}

func (ee *EnvEntry) GetTypeInfo() TypeInfo {
	return ee.typeInfo
}

func NewEnvironment(parent *Environment) *Environment {
	return &Environment{
		parent:   parent,
		identMap: make(map[string]EnvEntry),
	}
}

func (env *Environment) getParent() *Environment {
	return env.parent
}

func (env *Environment) getGlobal() *Environment {
	ret := env
	for {
		if ret.parent == nil {
			return ret
		}
		ret = ret.parent
	}
}

func (env *Environment) set(
	name string,
	uniqueName string,
	isExternal bool,
	category identCategory,
	typeInfo TypeInfo,
) {
	entry := EnvEntry{
		uniqueName: uniqueName,
		isExternal: isExternal,
		category:   category,
		typeInfo:   typeInfo,
	}

	env.identMap[name] = entry

	// Externally linked names must be added to
	// the global Environment as well
	if isExternal {
		env.getGlobal().identMap[name] = entry
	}

}

func (env *Environment) Get(name string) (*EnvEntry, *Environment) {
	ret, ok := env.identMap[name]
	if ok {
		return &ret, env
	}
	if env.parent != nil {
		return env.parent.Get(name)
	}

	return nil, nil
}

func (env *Environment) Lookup(name string) (string, error) {
	entry, definingEnv := env.Get(name)
	if definingEnv == nil {
		return "", errors.New(fmt.Sprintf("identifier '%s' is not defined", name))
	}
	return entry.uniqueName, nil
}
