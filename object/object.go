package object

type ObjectType string

const (
	ERROR_OBJ = "ERROR"
	FLOAT_OBJ = "FLOAT"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type ArithmeticObject interface {
	Adder
	Suber
	Muler
	Diver
}

type Adder interface {
	Object
	Add(right Object) Object
	Radd(left Object) Object
}

type Suber interface {
	Object
	Sub(right Object) Object
	Rsub(left Object) Object
}

type Muler interface {
	Object
	Mul(right Object) Object
	Rmul(left Object) Object
}

type Diver interface {
	Object
	Div(right Object) Object
	Rdiv(left Object) Object
}

type Booler interface {
	AsBool() Bool
}

type Compared interface {
	Object
	LesserThan(right Object) Object
	Equal(right Object) Object
}
