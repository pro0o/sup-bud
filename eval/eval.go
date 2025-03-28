package eval

import (
	"context"
	"fmt"
	"time"

	"github.com/pro0o/sup-bud/ast"
	"github.com/pro0o/sup-bud/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

type EvalOptions struct {
	MaxDepth int
	Timeout  time.Duration
}

func EvalWithOptions(node ast.Node, env *object.Environment, opts EvalOptions) object.Object {
	ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
	defer cancel()

	resultChan := make(chan object.Object, 1)
	errChan := make(chan error, 1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				errChan <- fmt.Errorf("panic during evaluation: %v", r)
			}
		}()

		result := evalWithDepthTracking(node, env, opts.MaxDepth)
		resultChan <- result
	}()

	select {
	case result := <-resultChan:
		return result
	case err := <-errChan:
		return newError("Evaluation error: %v", err)
	case <-ctx.Done():
		return newError("Evaluation timed out after %v", opts.Timeout)
	}
}

func evalWithDepthTracking(node ast.Node, env *object.Environment, maxDepth int) object.Object {
	if maxDepth <= 0 {
		return newError("Max recursion depth reached, Slow down brotherrrr—")
	}

	switch node := node.(type) {
	case *ast.Program:
		return evalProgramWithDepthTracking(node, env, maxDepth)

	case *ast.ExpressionStatement:
		return evalWithDepthTracking(node.Expression, env, maxDepth-1)

	case *ast.ReturnStatement:
		val := evalWithDepthTracking(node.ReturnValue, env, maxDepth-1)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.LetStatement:
		val := evalWithDepthTracking(node.Value, env, maxDepth-1)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
		return nil

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		right := evalWithDepthTracking(node.Right, env, maxDepth-1)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := evalWithDepthTracking(node.Left, env, maxDepth-1)
		if isError(left) {
			return left
		}
		right := evalWithDepthTracking(node.Right, env, maxDepth-1)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)

	case *ast.BlockStatement:
		return evalBlockStatementWithDepthTracking(node, env, maxDepth)

	case *ast.IfExpression:
		return evalIfExpressionWithDepthTracking(node, env, maxDepth)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env, Body: body}

	case *ast.CallExpression:
		function := evalWithDepthTracking(node.Function, env, maxDepth-1)
		if isError(function) {
			return function
		}
		args := evalExpressionsWithDepthTracking(node.Arguments, env, maxDepth-1)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunctionWithDepthTracking(function, args, maxDepth-1)

	default:
		return nil
	}
}

func evalProgramWithDepthTracking(program *ast.Program, env *object.Environment, maxDepth int) object.Object {
	var result object.Object
	for _, statement := range program.Statements {
		result = evalWithDepthTracking(statement, env, maxDepth-1)
		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}

func evalBlockStatementWithDepthTracking(block *ast.BlockStatement, env *object.Environment, maxDepth int) object.Object {
	var result object.Object
	for _, statement := range block.Statements {
		result = evalWithDepthTracking(statement, env, maxDepth-1)
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}
	return result
}

func evalIfExpressionWithDepthTracking(ie *ast.IfExpression, env *object.Environment, maxDepth int) object.Object {
	condition := evalWithDepthTracking(ie.Condition, env, maxDepth-1)
	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return evalWithDepthTracking(ie.Consequence, env, maxDepth-1)
	} else if ie.Alternative != nil {
		return evalWithDepthTracking(ie.Alternative, env, maxDepth-1)
	} else {
		return NULL
	}
}

func applyFunctionWithDepthTracking(bud object.Object, args []object.Object, maxDepth int) object.Object {
	function, ok := bud.(*object.Function)
	if !ok {
		return newError("not a function: %s", bud.Type())
	}
	extendedEnv := extendFunctionEnv(function, args)
	evaluated := evalWithDepthTracking(function.Body, extendedEnv, maxDepth)
	return unwrapReturnValue(evaluated)
}

func evalExpressionsWithDepthTracking(
	exps []ast.Expression,
	env *object.Environment,
	maxDepth int,
) []object.Object {
	var result []object.Object
	for _, e := range exps {
		evaluated := evalWithDepthTracking(e, env, maxDepth-1)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

// extension of func env
func extendFunctionEnv(
	bud *object.Function,
	args []object.Object,
) *object.Environment {
	env := object.NewEnclosedEnvironment(bud.Env)
	for paramIdx, param := range bud.Parameters {
		env.Set(param.Value, args[paramIdx])
	}
	return env
}

// func evalStatements(stmts []ast.Statement) object.Object {
// 	var result object.Object
// 	for _, statement := range stmts {
// 		result = Eval(statement)
// 		if returnValue, ok := result.(*object.ReturnValue); ok {
// 			return returnValue.Value
// 		}
// 	}
// 	return result
// }

// booll
func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

// prefix
func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

// prefix <- bang
func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

// prefix <- minus
func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}
	// if obj is int, we allcoate a new obbj to wrap a
	// negated version to this val
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

// infix
func evalInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {

	switch {
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s",
			left.Type(), operator, right.Type())
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

// infix <- int
func evalIntegerInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	// lavde lang doesnt allow pointer comparison for int objs.
	// *obj.int alwats allocates new instacnes of obj.integer
	// thus slower than bool exp.
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case TRUE:
		return true
	case FALSE:
		return false
	case NULL:
		return false
	default:
		return true
	}
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func evalIdentifier(
	node *ast.Identifier,
	env *object.Environment) object.Object {
	val, ok := env.Get(node.Value)
	if !ok {
		return newError("identifier not found: %s", node.Value)
	}
	return val
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}
