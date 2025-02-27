package object

import (
	"fmt"
)

type ObjectType string

const (
	INTEGER_OBJ = "INTEGER"
	FLOAT_OBJ   = "FLOAT"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (o Integer) Type() ObjectType { return INTEGER_OBJ }
func (o Integer) Inspect() string  { return fmt.Sprintf("%d", o.Value) }

type Float struct {
	Value float64
}

func (o Float) Type() ObjectType { return FLOAT_OBJ }
func (o Float) Inspect() string  { return fmt.Sprintf("%f", o.Value) }

type Adder interface {
	Add(right Object) Object
	Radd(left Object) Object
}

func (o *Integer) Add(right Object) Object {
	switch right := right.(type) {
	case *Integer:
		return &Integer{Value: o.Value + right.Value}
	case *Float:
		return &Float{Value: float64(o.Value) + right.Value}
	default:
		return nil
	}
}

func (o *Integer) Radd(left Object) Object {
	return o.Add(left)
}

func (o *Float) Add(right Object) Object {
	switch right := right.(type) {
	case *Integer:
		return &Float{Value: o.Value + float64(right.Value)}
	case *Float:
		return &Float{Value: o.Value + right.Value}
	default:
		return nil
	}
}

func (o *Float) Radd(left Object) Object {
	return o.Add(left)
}

type Suber interface {
	Sub(right Object) Object
	Rsub(left Object) Object
}

func (o *Integer) Sub(right Object) Object {
	switch right := right.(type) {
	case *Integer:
		return &Integer{Value: o.Value - right.Value}
	case *Float:
		return &Float{Value: float64(o.Value) - right.Value}
	default:
		return nil
	}
}

func (o *Integer) Rsub(left Object) Object {
	switch left := left.(type) {
	case *Integer:
		return &Integer{Value: left.Value - o.Value}
	case *Float:
		return &Float{Value: left.Value - float64(o.Value)}
	default:
		return nil
	}
}

func (o *Float) Sub(right Object) Object {
	switch right := right.(type) {
	case *Integer:
		return &Float{Value: o.Value - float64(right.Value)}
	case *Float:
		return &Float{Value: o.Value - right.Value}
	default:
		return nil
	}
}

func (o *Float) Rsub(left Object) Object {
	switch left := left.(type) {
	case *Integer:
		return &Float{Value: float64(left.Value) - o.Value}
	case *Float:
		return &Float{Value: left.Value - o.Value}
	default:
		return nil
	}
}

type Muler interface {
	Mul(right Object) Object
	Rmul(left Object) Object
}

func (o *Integer) Mul(right Object) Object {
	switch right := right.(type) {
	case *Integer:
		return &Integer{Value: o.Value * right.Value}
	case *Float:
		return &Float{Value: float64(o.Value) * right.Value}
	default:
		return nil
	}
}

func (o *Integer) Rmul(left Object) Object {
	return o.Mul(left)
}

func (o *Float) Mul(right Object) Object {
	switch right := right.(type) {
	case *Integer:
		return &Float{Value: o.Value * float64(right.Value)}
	case *Float:
		return &Float{Value: o.Value * right.Value}
	default:
		return nil
	}
}

func (o *Float) Rmul(left Object) Object {
	return o.Mul(left)
}

type Diver interface {
	Div(right Object) Object
	Rdiv(left Object) Object
}

func (o *Integer) Div(right Object) Object {
	switch right := right.(type) {
	case *Integer:
		return &Integer{Value: o.Value / right.Value}
	case *Float:
		return &Float{Value: float64(o.Value) / right.Value}
	default:
		return nil
	}
}

func (o *Integer) Rdiv(left Object) Object {
	switch left := left.(type) {
	case *Integer:
		return &Integer{Value: left.Value / o.Value}
	case *Float:
		return &Float{Value: left.Value / float64(o.Value)}
	default:
		return nil
	}
}

func (o *Float) Div(right Object) Object {
	switch right := right.(type) {
	case *Integer:
		return &Float{Value: o.Value / float64(right.Value)}
	case *Float:
		return &Float{Value: o.Value / right.Value}
	default:
		return nil
	}
}

func (o *Float) Rdiv(left Object) Object {
	switch left := left.(type) {
	case *Integer:
		return &Float{Value: float64(left.Value) / o.Value}
	case *Float:
		return &Float{Value: left.Value / o.Value}
	default:
		return nil
	}
}
