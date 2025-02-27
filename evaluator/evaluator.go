package evaluator

import (
	"github.com/Richtermnd/ferret/ast"
	"github.com/Richtermnd/ferret/object"
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)

	case *ast.ExpressionStatement:
		return Eval(node.Expr)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}

	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right)
	}
	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	var res object.Object
	for _, stmt := range stmts {
		res = Eval(stmt)
	}
	return res
}

func evalPrefixExpression(op string, right object.Object) object.Object {
	switch op {
	case "-":
		return evalMinusPrefixOperator(right)
	}
	return nil
}

func evalMinusPrefixOperator(right object.Object) object.Object {
	switch right.Type() {
	case object.INTEGER_OBJ:
		return &object.Integer{Value: -right.(*object.Integer).Value}
	case object.FLOAT_OBJ:
		return &object.Float{Value: -right.(*object.Float).Value}
	}
	return nil
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
	return nil
}

func add(left, right object.Object) object.Object {
	leftAdder, ok := left.(object.Adder)
	if ok {
		res := leftAdder.Add(right)
		if res != nil {
			return res
		}
	}
	rightAdder, ok := right.(object.Adder)
	if !ok {
		return nil
	}
	return rightAdder.Radd(left)
}

func sub(left, right object.Object) object.Object {
	var res object.Object
	leftSuber, ok := left.(object.Suber)
	if ok {
		res = leftSuber.Sub(right)
		if res != nil {
			return res
		}
	}
	rightSuber, ok := right.(object.Suber)
	if !ok {
		return nil
	}
	return rightSuber.Rsub(left)
}

func mul(left, right object.Object) object.Object {
	var res object.Object
	leftMuler, ok := left.(object.Muler)
	if ok {
		res = leftMuler.Mul(right)
		if res != nil {
			return res
		}
	}
	rightMuler, ok := right.(object.Muler)
	if !ok {
		return nil
	}
	return rightMuler.Rmul(left)
}

func div(left, right object.Object) object.Object {
	var res object.Object
	leftDiver, ok := left.(object.Diver)
	if ok {
		res = leftDiver.Div(right)
		if res != nil {
			return res
		}
	}

	rightDiver, ok := right.(object.Diver)
	if !ok {
		return nil
	}
	return rightDiver.Rdiv(left)
}
