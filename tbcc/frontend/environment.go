package frontend

import (
	"errors"
	"fmt"
)

type environment struct {
	parent *environment
	varMap map[string]string
}

func newEnvironment(parent *environment) *environment {
	return &environment{
		parent: parent,
		varMap: make(map[string]string),
	}
}

func (env *environment) getParent() *environment {
	return env.parent
}

func (env *environment) set(name string, uniqueName string) {
	env.varMap[name] = uniqueName
}

func (env *environment) isSet(name string) bool {
	_, ok := env.varMap[name]
	return ok
}

func (env *environment) lookup(name string) (string, error) {
	ret, ok := env.varMap[name]
	if ok {
		return ret, nil
	}
	if env.parent != nil {
		return env.parent.lookup(name)
	}

	return "", errors.New(fmt.Sprintf("variable '%s' is not defined", name))
}
