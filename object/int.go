package object

import "fmt"

type Integer struct {
	Value int64
}

func (o Integer) Type() ObjectType { return INTEGER_OBJ }
func (o Integer) Inspect() string  { return fmt.Sprintf("%d", o.Value) }

func (o *Integer) Add(right Object) Object {
	switch right := right.(type) {
	case *Integer:
		return &Integer{Value: o.Value + right.Value}
	case *Float:
		return &Float{Value: float64(o.Value) + right.Value}
	default:
		return NewError(UNSUPPORTED_ERR, "%s %s %s", o.Type(), "+", right.Type())
	}
}

func (o *Integer) Radd(left Object) Object {
	return o.Add(left)
}

func (o *Integer) Sub(right Object) Object {
	switch right := right.(type) {
	case *Integer:
		return &Integer{Value: o.Value - right.Value}
	case *Float:
		return &Float{Value: float64(o.Value) - right.Value}
	default:
		return NewError(UNSUPPORTED_ERR, "%s %s %s", o.Type(), "-", right.Type())
	}
}

func (o *Integer) Rsub(left Object) Object {
	switch left := left.(type) {
	case *Integer:
		return &Integer{Value: left.Value - o.Value}
	case *Float:
		return &Float{Value: left.Value - float64(o.Value)}
	default:
		return NewError(UNSUPPORTED_ERR, "%s %s %s", left.Type(), "-", o.Type())
	}
}

func (o *Integer) Mul(right Object) Object {
	switch right := right.(type) {
	case *Integer:
		return &Integer{Value: o.Value * right.Value}
	case *Float:
		return &Float{Value: float64(o.Value) * right.Value}
	default:
		return NewError(UNSUPPORTED_ERR, "%s %s %s", o.Type(), "*", right.Type())
	}
}

func (o *Integer) Rmul(left Object) Object {
	return o.Mul(left)
}

func (o *Integer) Div(right Object) Object {
	switch right := right.(type) {
	case *Integer:
		return &Integer{Value: o.Value / right.Value}
	case *Float:
		return &Float{Value: float64(o.Value) / right.Value}
	default:
		return NewError(UNSUPPORTED_ERR, "%s %s %s", o.Type(), "/", right.Type())
	}
}

func (o *Integer) Rdiv(left Object) Object {
	switch left := left.(type) {
	case *Integer:
		return &Integer{Value: left.Value / o.Value}
	case *Float:
		return &Float{Value: left.Value / float64(o.Value)}
	default:
		return NewError(UNSUPPORTED_ERR, "%s %s %s", left.Type(), "/", o.Type())
	}
}
