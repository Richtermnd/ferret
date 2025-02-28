package object

type ObjectType string

const (
	ERROR_OBJ   = "ERROR"
	INTEGER_OBJ = "INTEGER"
	FLOAT_OBJ   = "FLOAT"
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
	Add(right Object) Object
	Radd(left Object) Object
}

type Suber interface {
	Sub(right Object) Object
	Rsub(left Object) Object
}

type Muler interface {
	Mul(right Object) Object
	Rmul(left Object) Object
}

type Diver interface {
	Div(right Object) Object
	Rdiv(left Object) Object
}
