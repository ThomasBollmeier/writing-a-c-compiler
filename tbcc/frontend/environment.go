package frontend

import (
	"errors"
	"fmt"
)

type Environment struct {
	parent   *Environment
	identMap map[string]*EnvEntry
}

type linkage int

const (
	linkNone linkage = iota
	linkStatic
	linkExternal
)

type identCategory int

const (
	idCatVariable identCategory = iota
	idCatFunction
	idCatParameter
)

type EnvEntry struct {
	uniqueName string
	linkage    linkage
	category   identCategory
	typeInfo   TypeInfo
}

func (ee *EnvEntry) GetTypeInfo() TypeInfo {
	return ee.typeInfo
}

func NewEnvironment(parent *Environment) *Environment {
	return &Environment{
		parent:   parent,
		identMap: make(map[string]*EnvEntry),
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
	linkage linkage,
	category identCategory,
	typeInfo TypeInfo,
) {
	entry := &EnvEntry{
		uniqueName: uniqueName,
		linkage:    linkage,
		category:   category,
		typeInfo:   typeInfo,
	}

	env.identMap[name] = entry

}

func (env *Environment) Get(name string) (*EnvEntry, *Environment) {
	ret, ok := env.identMap[name]
	if ok {
		return ret, env
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

type Environments struct {
	Global map[string]*EnvEntry
	Block  *Environment
}

func NewEnvironments() *Environments {
	return &Environments{
		Global: make(map[string]*EnvEntry),
		Block:  NewEnvironment(nil),
	}
}

func (envs *Environments) beginBlock() {
	envs.Block = NewEnvironment(envs.Block)
}

func (envs *Environments) endBlock() {
	envs.Block = envs.Block.getParent()
}

func (envs *Environments) set(
	name string,
	uniqueName string,
	linkage linkage,
	category identCategory,
	typeInfo TypeInfo,
) {
	envEntry := &EnvEntry{
		uniqueName: uniqueName,
		linkage:    linkage,
		category:   category,
		typeInfo:   typeInfo,
	}

	if linkage != linkNone {
		envs.Global[name] = envEntry
	}

	envs.Block.identMap[name] = envEntry
}

func (envs *Environments) Get(name string) (*EnvEntry, *Environment) {
	return envs.Block.Get(name)
}

func (envs *Environments) Lookup(name string) (string, error) {
	return envs.Block.Lookup(name)
}
