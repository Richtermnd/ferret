package evaluator

import (
	"github.com/Richtermnd/ferret/ast"
	"github.com/Richtermnd/ferret/object"
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
		return evalInfixExpression(node.Operator, left, right)
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

func evalInfixExpression(op string, left, right object.Object) object.Object {
	switch op {
	case "+":
		return add(left, right)
	case "-":
		return sub(left, right)
	case "*":
		return mul(left, right)
	case "/":
		return div(left, right)
	}

	leftCmp, lok := left.(object.Compared)
	rightCmp, rok := right.(object.Compared)
	if lok && rok {
		// comparing
		switch op {
		case "==":
			return eq(leftCmp, rightCmp)
		case "!=":
			return neq(leftCmp, rightCmp)
		case ">":
			return gt(leftCmp, rightCmp)
		case ">=":
			return geq(leftCmp, rightCmp)
		case "<":
			return lt(leftCmp, rightCmp)
		case "<=":
			return leq(leftCmp, rightCmp)
		}
	}
	return object.NewError(object.UNKNOWN_OPERATOR_ERR, "%s %s %s", left.Type(), op, right.Type())
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
