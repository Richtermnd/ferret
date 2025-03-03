package object

import "fmt"

type Float struct {
	Value float64
}

func (o Float) Type() ObjectType { return FLOAT_OBJ }
func (o Float) Inspect() string  { return fmt.Sprintf("%f", o.Value) }

func (o *Float) Add(right Object) Object {
	switch right := right.(type) {
	case *Integer:
		return &Float{Value: o.Value + float64(right.Value)}
	case *Float:
		return &Float{Value: o.Value + right.Value}
	default:
		return NewError(UNSUPPORTED_ERR, "%s %s %s", o.Type(), "+", right.Type())
	}
}

func (o *Float) Radd(left Object) Object {
	return o.Add(left)
}

func (o *Float) Sub(right Object) Object {
	switch right := right.(type) {
	case *Integer:
		return &Float{Value: o.Value - float64(right.Value)}
	case *Float:
		return &Float{Value: o.Value - right.Value}
	default:
		return NewError(UNSUPPORTED_ERR, "%s %s %s", o.Type(), "-", right.Type())
	}
}

func (o *Float) Rsub(left Object) Object {
	switch left := left.(type) {
	case *Integer:
		return &Float{Value: float64(left.Value) - o.Value}
	case *Float:
		return &Float{Value: left.Value - o.Value}
	default:
		return NewError(UNSUPPORTED_ERR, "%s %s %s", left.Type(), "-", o.Type())
	}
}

func (o *Float) Mul(right Object) Object {
	switch right := right.(type) {
	case *Integer:
		return &Float{Value: o.Value * float64(right.Value)}
	case *Float:
		return &Float{Value: o.Value * right.Value}
	default:
		return NewError(UNSUPPORTED_ERR, "%s %s %s", o.Type(), "*", right.Type())
	}
}

func (o *Float) Rmul(left Object) Object {
	return o.Mul(left)
}

func (o *Float) Div(right Object) Object {
	switch right := right.(type) {
	case *Integer:
		return &Float{Value: o.Value / float64(right.Value)}
	case *Float:
		return &Float{Value: o.Value / right.Value}
	default:
		return NewError(UNSUPPORTED_ERR, "%s %s %s", o.Type(), "/", right.Type())
	}
}

func (o *Float) Rdiv(left Object) Object {
	switch left := left.(type) {
	case *Integer:
		return &Float{Value: float64(left.Value) / o.Value}
	case *Float:
		return &Float{Value: left.Value / o.Value}
	default:
		return NewError(UNSUPPORTED_ERR, "%s %s %s", left.Type(), "/", o.Type())
	}
}

func (o *Float) LesserThan(right Object) Object {
	switch r := right.(type) {
	case *Integer:
		return &Bool{Value: o.Value < float64(r.Value)}
	case *Float:
		return &Bool{Value: float64(o.Value) < r.Value}
	default:
		return NewError(UNSUPPORTED_ERR, "%s and %s not comparable", o.Type(), r.Type())
	}
}

func (o *Float) Equal(right Object) Object {
	switch r := right.(type) {
	case *Integer:
		return &Bool{Value: o.Value == float64(r.Value)}
	case *Float:
		return &Bool{Value: float64(o.Value) == r.Value}
	default:
		return NewError(UNSUPPORTED_ERR, "%s and %s not comparable", o.Type(), r.Type())
	}
}
