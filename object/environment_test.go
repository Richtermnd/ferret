package object_test

import (
	"testing"

	"github.com/Richtermnd/ferret/object"
)

func TestSubEnvironment(t *testing.T) {
	env := object.NewEnv()
	env.Set("a", &object.Integer{1})
	env.Set("b", &object.Integer{2})
	subEnv := env.SubEnv()
	subEnv.Set("c", &object.Integer{3})
	subEnv.Set("b", &object.Integer{4})
	env.Set("a", &object.Integer{5})
	t.Log(subEnv)

	if _, ok := env.Get("c"); ok {
		t.Errorf("env can get subEnv values\n")
	}

	envB, ok := env.Get("b")
	if !ok {
		t.Errorf("env: can't get b")
	}

	subEnvB, ok := subEnv.Get("b")
	if !ok {
		t.Errorf("subenv: can't get b")
	}

	if v := envB.(*object.Integer).Value; v != 2 {
		t.Errorf("env: wrong b expected: %d got: %d\n", 2, v)
	}
	if v := subEnvB.(*object.Integer).Value; v != 4 {
		t.Errorf("subenv: wrong b expected: %d got: %d\n", 4, v)
	}

	if obj, ok := subEnv.Get("a"); !ok {
		t.Errorf("subEnv: cannot get 'a' from outer env\n")
	} else if v := obj.(*object.Integer).Value; v == 1 {
		t.Errorf("subEnv: 'a' not update a value\n")
	} else if v != 5 {
		t.Errorf("subEnv: 'a' wrong value %d\n", v)
	}
}
