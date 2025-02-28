package object

import (
	"fmt"
	"strings"
)

type Environment struct {
	outer *Environment
	env   map[string]Object
}

func NewEnv() *Environment {
	return &Environment{
		env: make(map[string]Object),
	}
}

func (e *Environment) Get(s string) (Object, bool) {
	obj, ok := e.env[s]
	if !ok && e.outer != nil {
		return e.outer.Get(s)
	}
	return obj, ok
}

func (e *Environment) Set(s string, obj Object) Object {
	e.env[s] = obj
	return obj
}

func (e *Environment) SubEnv() *Environment {
	ne := NewEnv()
	ne.outer = e
	return ne
}

func (e *Environment) String() string {
	sb := strings.Builder{}
	e.buildString(&sb, 0)
	return sb.String()
}

func (e *Environment) buildString(sb *strings.Builder, level int) int {
	if e.outer != nil {
		level = e.outer.buildString(sb, level)
	}
	indent := strings.Repeat("  ", level)
	fmt.Fprintf(sb, "%s%s\n", indent, "{")
	for k, v := range e.env {
		fmt.Fprintf(sb, "%s  %s: %T %+v\n", indent, k, v, v)
	}
	fmt.Fprintf(sb, "%s%s\n", indent, "}")
	return level + 1
}
