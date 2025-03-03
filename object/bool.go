package object

import "fmt"

const BOOL_OBJ ObjectType = "BOOL"

type Bool struct {
	Value bool
}

func (o Bool) Type() ObjectType { return BOOL_OBJ }
func (o Bool) Inspect() string  { return fmt.Sprintf("%t", o.Value) }

func (o Bool) AsInt() *Integer {
	if o.Value {
		return &Integer{Value: 1}
	}
	return &Integer{0}
}

func (o Bool) AsFloat() *Float {
	if o.Value {
		return &Float{Value: 1}
	}
	return &Float{Value: 0}
}

func (o *Bool) Add(right Object) Object {
	switch right := right.(type) {
	case *Bool:
		return o.AsInt().Add(right.AsInt())
	case *Integer:
		return o.AsInt().Add(right)
	case *Float:
		return o.AsFloat().Add(right)
	default:
		return NewError(UNSUPPORTED_ERR, "%s %s %s", o.Type(), "+", right.Type())
	}
}

func (o *Bool) Radd(left Object) Object {
	return o.Add(left)
}

func (o *Bool) Sub(right Object) Object {
	switch right := right.(type) {
	case *Bool:
		return o.AsInt().Sub(right.AsInt())
	case *Integer:
		return o.AsInt().Sub(right)
	case *Float:
		return o.AsFloat().Sub(right)
	default:
		return NewError(UNSUPPORTED_ERR, "%s %s %s", o.Type(), "-", right.Type())
	}
}

func (o *Bool) Rsub(left Object) Object {
	switch left := left.(type) {
	case *Bool:
		return o.AsInt().Rsub(left.AsInt())
	case *Integer:
		return o.AsInt().Rsub(left)
	case *Float:
		return o.AsFloat().Rsub(left)
	default:
		return NewError(UNSUPPORTED_ERR, "%s %s %s", left.Type(), "-", o.Type())
	}
}

func (o *Bool) Mul(right Object) Object {
	switch right := right.(type) {
	case *Bool:
		return o.AsInt().Mul(right.AsInt())
	case *Integer:
		return o.AsInt().Mul(right)
	case *Float:
		return o.AsFloat().Mul(right)
	default:
		return NewError(UNSUPPORTED_ERR, "%s %s %s", o.Type(), "*", right.Type())
	}
}

func (o *Bool) Rmul(left Object) Object {
	return o.Mul(left)
}

func (o *Bool) LesserThan(right Object) Object {
	switch r := right.(type) {
	case *Bool:
		return &Bool{Value: o.AsInt().Value < r.AsInt().Value}
	case *Integer:
		return &Bool{Value: o.AsInt().Value < r.Value}
	case *Float:
		return &Bool{Value: o.AsFloat().Value < r.Value}
	default:
		return NewError(UNSUPPORTED_ERR, "%s and %s not comparable", o.Type(), r.Type())
	}
}

func (o *Bool) Equal(right Object) Object {
	switch r := right.(type) {
	case *Bool:
		return &Bool{Value: o.Value == r.Value}
	case *Integer:
		return &Bool{Value: (r.Value == 0) != o.Value}
	case *Float:
		return &Bool{Value: (r.Value == 0) != o.Value}
	default:
		return NewError(UNSUPPORTED_ERR, "%s and %s not comparable", o.Type(), r.Type())
	}
}
