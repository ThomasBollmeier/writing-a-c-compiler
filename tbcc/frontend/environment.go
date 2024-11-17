package frontend

import (
	"errors"
	"fmt"
)

type environment struct {
	parent   *environment
	identMap map[string]envEntry
}

type identCategory int

const (
	idCatVariable identCategory = iota
	idCatFunction
	idCatParameter
)

type envEntry struct {
	uniqueName string
	isExternal bool
	category   identCategory
	typeInfo   typeInfo
}

func newEnvironment(parent *environment) *environment {
	return &environment{
		parent:   parent,
		identMap: make(map[string]envEntry),
	}
}

func (env *environment) getParent() *environment {
	return env.parent
}

func (env *environment) getGlobal() *environment {
	ret := env
	for {
		if ret.parent == nil {
			return ret
		}
		ret = ret.parent
	}
}

func (env *environment) set(
	name string,
	uniqueName string,
	isExternal bool,
	category identCategory,
	typeInfo typeInfo,
) {
	entry := envEntry{
		uniqueName: uniqueName,
		isExternal: isExternal,
		category:   category,
		typeInfo:   typeInfo,
	}

	env.identMap[name] = entry

	// Externally linked names must be added to
	// the global environment as well
	if isExternal {
		env.getGlobal().identMap[name] = entry
	}

}

func (env *environment) get(name string) (*envEntry, *environment) {
	ret, ok := env.identMap[name]
	if ok {
		return &ret, env
	}
	if env.parent != nil {
		return env.parent.get(name)
	}

	return nil, nil
}

func (env *environment) lookup(name string) (string, error) {
	entry, definingEnv := env.get(name)
	if definingEnv == nil {
		return "", errors.New(fmt.Sprintf("identifier '%s' is not defined", name))
	}
	return entry.uniqueName, nil
}
