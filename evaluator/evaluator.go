package evaluator

import (
	"github.com/Richtermnd/ferret/ast"
	"github.com/Richtermnd/ferret/object"
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
