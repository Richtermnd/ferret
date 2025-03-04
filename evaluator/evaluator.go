package evaluator

import (
	"github.com/Richtermnd/ferret/ast"
	"github.com/Richtermnd/ferret/object"
	"github.com/Richtermnd/ferret/token"
)

var (
	TRUE  = &object.Bool{Value: true}
	FALSE = &object.Bool{Value: false}
)

func Eval(env *object.Environment, node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(env, node.Statements)

	case *ast.Identifier:
		return evalIdentifier(env, node)

	case *ast.BlockStatement:
		return evalStatements(env.SubEnv(), node.Statements)

	case *ast.LetStatement:
		value := Eval(env, node.Value)
		if object.IsError(value) {
			return value
		}
		env.Set(node.Name.Value, value)

	case *ast.ExpressionStatement:
		return Eval(env, node.Expr)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}

	case *ast.BooleanLiteral:
		return boolFromNative(node.Value)

	case *ast.PrefixExpression:
		right := Eval(env, node.Right)
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(env, node.Left)
		right := Eval(env, node.Right)
		return evalInfixExpression(node.Token, left, right)
	}

	return nil
}

func evalStatements(env *object.Environment, stmts []ast.Statement) object.Object {
	var res object.Object
	for _, stmt := range stmts {
		res = Eval(env, stmt)
	}
	return res
}

func evalIdentifier(env *object.Environment, node *ast.Identifier) object.Object {
	obj, ok := env.Get(node.Value)
	if !ok {
		return object.NewError(object.NOT_FOUND_ERR, "%s", node.Value)
	}
	return obj
}

func evalPrefixExpression(op string, right object.Object) object.Object {
	switch op {
	case "!":
		return not(right)
	case "-":
		return evalMinusPrefixOperator(right)
	}
	return object.NewError(object.UNKNOWN_OPERATOR_ERR, "%s%s", op, right.Type())
}

func evalMinusPrefixOperator(right object.Object) object.Object {
	switch right.Type() {
	case object.INTEGER_OBJ:
		return &object.Integer{Value: -right.(*object.Integer).Value}
	case object.FLOAT_OBJ:
		return &object.Float{Value: -right.(*object.Float).Value}
	}
	return object.NewError(object.UNSUPPORTED_ERR, "-%s", right.Type())
}

// TODO: split this func in smaller pieces
func evalInfixExpression(tok token.Token, left, right object.Object) object.Object {
	if tok.Type == token.AND || tok.Type == token.OR {
		return evalLogicExpression(tok, left, right)
	}
	switch tok.Type {
	case token.ADD:
		return add(left, right)
	case token.SUB:
		return sub(left, right)
	case token.MUL:
		return mul(left, right)
	case token.DIV:
		return div(left, right)
	}

	leftCmp, lok := left.(object.Compared)
	rightCmp, rok := right.(object.Compared)
	if lok && rok {
		// comparing
		switch tok.Type {
		case token.EQ:
			return eq(leftCmp, rightCmp)
		case token.NEQ:
			return neq(leftCmp, rightCmp)
		case token.LT:
			return lt(leftCmp, rightCmp)
		case token.LEQ:
			return leq(leftCmp, rightCmp)
		case token.GT:
			return gt(leftCmp, rightCmp)
		case token.GEQ:
			return geq(leftCmp, rightCmp)
		}
	}
	return object.NewError(object.UNKNOWN_OPERATOR_ERR, "%s %s %s", left.Type(), tok.String(), right.Type())
}

func evalLogicExpression(tok token.Token, left, right object.Object) object.Object {
	var a, b bool

	if lbool, ok := left.(*object.Bool); ok {
		a = lbool.Value
	} else if lbooler, ok := left.(object.Booler); ok {
		a = lbooler.AsBool().Value
	} else {
		return object.NewError(object.UNEXPECTED, "expected: bool or booler got: %s", left.Type())
	}

	if rbool, ok := right.(*object.Bool); ok {
		b = rbool.Value
	} else if rbooler, ok := right.(object.Booler); ok {
		b = rbooler.AsBool().Value
	} else {
		return object.NewError(object.UNEXPECTED, "expected: bool or booler got: %s", right.Type())
	}

	switch tok.Type {
	case token.AND:
		return boolFromNative(a && b)
	case token.OR:
		return boolFromNative(a || b)
	}
	return object.NewError(object.UNEXPECTED, "unexpected operator (probably a bug): %s", tok.Literal)
}

func add(left, right object.Object) object.Object {
	leftAdder, ok := left.(object.Adder)
	if ok {
		res := leftAdder.Add(right)
		if !object.IsError(res) {
			return res
		}
	}
	rightAdder, ok := right.(object.Adder)
	if !ok {
		return object.NewError(object.NOT_IMPLEMENTED_ERR, "%s + %s", left.Type(), right.Type())
	}
	return rightAdder.Radd(left)
}

func sub(left, right object.Object) object.Object {
	var res object.Object
	leftSuber, ok := left.(object.Suber)
	if ok {
		res = leftSuber.Sub(right)
		if !object.IsError(res) {
			return res
		}
	}
	rightSuber, ok := right.(object.Suber)
	if !ok {
		return object.NewError(object.NOT_IMPLEMENTED_ERR, "%s - %s", left.Type(), right.Type())
	}
	return rightSuber.Rsub(left)
}

func mul(left, right object.Object) object.Object {
	var res object.Object
	leftMuler, ok := left.(object.Muler)
	if ok {
		res = leftMuler.Mul(right)
		if !object.IsError(res) {
			return res
		}
	}
	rightMuler, ok := right.(object.Muler)
	if !ok {
		return object.NewError(object.NOT_IMPLEMENTED_ERR, "%s * %s", left.Type(), right.Type())
	}
	return rightMuler.Rmul(left)
}

func div(left, right object.Object) object.Object {
	var res object.Object
	leftDiver, ok := left.(object.Diver)
	if ok {
		res = leftDiver.Div(right)
		if !object.IsError(res) {
			return res
		}
	}

	rightDiver, ok := right.(object.Diver)
	if !ok {
		return object.NewError(object.NOT_IMPLEMENTED_ERR, "%s / %s", left.Type(), right.Type())
	}
	return rightDiver.Rdiv(left)
}

func not(v object.Object) object.Object {
	if v.Type() == object.BOOL_OBJ {
		return &object.Bool{Value: !v.(*object.Bool).Value}
	} else if v, ok := v.(object.Booler); ok {
		return &object.Bool{Value: !v.AsBool().Value}
	}
	return object.NewError(object.UNSUPPORTED_ERR, "cannot represent %s as bool", v.Type())
}

func and(left, right object.Object) object.Object {
	var l, r bool
	if left.Type() == object.BOOL_OBJ && right.Type() == object.BOOL_OBJ {
		l = left.(*object.Bool).Value
		r = right.(*object.Bool).Value
		return boolFromNative(l && r)
	}
	lBooler, ok := left.(object.Booler)
	if !ok {
		return object.NewError(object.UNSUPPORTED_ERR, "cannot represent %s as bool", left.Type())
	}
	rBooler, ok := right.(object.Booler)
	if !ok {
		return object.NewError(object.UNSUPPORTED_ERR, "cannot represent %s as bool", right.Type())
	}
	return boolFromNative(lBooler.AsBool().Value && rBooler.AsBool().Value)
}

func or(left, right object.Object) object.Object {
	var l, r bool
	if left.Type() == object.BOOL_OBJ && right.Type() == object.BOOL_OBJ {
		l = left.(*object.Bool).Value
		r = right.(*object.Bool).Value
		return boolFromNative(l || r)
	}
	lBooler, ok := left.(object.Booler)
	if !ok {
		return object.NewError(object.UNSUPPORTED_ERR, "cannot represent %s as bool", left.Type())
	}
	rBooler, ok := right.(object.Booler)
	if !ok {
		return object.NewError(object.UNSUPPORTED_ERR, "cannot represent %s as bool", right.Type())
	}
	return boolFromNative(lBooler.AsBool().Value || rBooler.AsBool().Value)
}

//

func eq(left, right object.Compared) object.Object {
	// a == b <=> b == a
	if res := left.Equal(right); res.Type() != object.ERROR_OBJ {
		return res
	}
	return right.Equal(left)
}

func neq(left, right object.Compared) object.Object {
	// a != b <=> !(a == b)
	return not(eq(left, right))
}

// There is definitely a recursion somewhere here, but tests passed, so...

func lt(left, right object.Compared) object.Object {
	// a < b <=> b >= a
	ltLeft := left.LesserThan(right)
	if ltLeft.Type() != object.ERROR_OBJ {
		return ltLeft
	}
	return geq(right, left)
}

func gt(left, right object.Compared) object.Object {
	// a > b <=> b <= a
	return lt(right, left)
}

func leq(left, right object.Compared) object.Object {
	// a <= b <=> !(a > b)
	return not(gt(left, right))
}

func geq(left, right object.Compared) object.Object {
	// a >= b <=> !(a < b)
	return not(lt(left, right))
}

func boolFromNative(input bool) *object.Bool {
	if input {
		return TRUE
	}
	return FALSE
}
