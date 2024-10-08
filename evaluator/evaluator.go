package evaluator

import (
	"dodo-lang/ast"
	"dodo-lang/object"
	"fmt"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, expected=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"rest": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, expected=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				if len(arg.Elements) > 0 {

					arrLen := len(arg.Elements)
					newArr := make([]object.Object, arrLen-1, arrLen-1)
					copy(newArr, arg.Elements[1:])

					return &object.Array{Elements: newArr}
				}

				return NULL
			case *object.String:
				if len(arg.Value) > 0 {
					return &object.String{Value: arg.Value[1:]}
				}

				return NULL
			default:
				return newError("argument to `rest` not supported, got %s", args[0].Type())
			}
		},
	},
	"first": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, expected=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				if len(arg.Elements) > 0 {
					return arg.Elements[0]
				}

				return NULL
			case *object.String:
				if len(arg.Value) > 0 {
					return &object.String{Value: string(arg.Value[0])}
				}

				return NULL
			default:
				return newError("argument to `first` not supported, got %s", args[0].Type())
			}
		},
	},
	"last": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, expected=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				max := len(arg.Elements) - 1
				if len(arg.Elements) > 0 {
					return arg.Elements[max]
				}

				return NULL
			case *object.String:
				max := len(arg.Value) - 1
				if len(arg.Value) > 0 {
					return &object.String{Value: string(arg.Value[max])}
				}

				return NULL
			default:
				return newError("argument to `last` not supported, got %s", args[0].Type())
			}
		},
	},
	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, expected=2", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				length := len(arg.Elements)
				newArr := make([]object.Object, length+1, length+1)
				copy(newArr, arg.Elements)
				newArr[length] = args[1]

				return &object.Array{Elements: newArr}
			default:
				return newError("argument to `last` not supported, got %s", args[0].Type())
			}
		},
	},
	"typeof": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, expected=1", len(args))
			}

			arg := args[0].(object.Object)

			return &object.String{Value: string(arg.Type())}
		},
	},
	"debug": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, expected=1", len(args))
			}

			arg := args[0].(object.Object)

			fmt.Printf("%s\n", arg.Inspect())

			return NULL
		},
	},
	"println": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}

			return NULL
		},
	},
	"printf": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) < 2 {
				return newError("wrong number of arguments. got=%d, expected at least 2", len(args))
			}

			formatStr, ok := args[0].(*object.String)

			if !ok {
				return newError("first argument has to be a string. got=%s", args[0].Type())
			}

			var templateArgs []any

			for _, arg := range args[1:] {
				switch a := arg.(type) {
				case *object.String:
					templateArgs = append(templateArgs, a.Value)
				case *object.Integer:
					templateArgs = append(templateArgs, a.Value)
				default:
					return newError("only strings and integers can be used with 'printf'. got=%s", a.Type())
				}
			}

			fmt.Printf(formatStr.Value, templateArgs...)
			fmt.Print("\n")

			return NULL
		},
	},
}

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.LetStatement:
		val := Eval(node.Value, env)

		if isError(val) {
			return val
		}

		if _, ok := env.Get(node.Name.Value); ok {
			return newError("identifier '%s' already exists", node.Name.Value)
		}

		env.Set(node.Name.Value, node.Mutable, val)
	case *ast.ReassignmentStatement:
		val := Eval(node.Value, env)

		if isError(val) {
			return val
		}

		if !env.IsMutable(node.Ident.Value) {
			return newError("identifier '%s' is not mutable", node.Ident.Value)
		}

		env.Set(node.Ident.Value, true, val)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)

		if isError(val) {
			return val
		}

		return &object.ReturnValue{Value: val}

	// Expressions
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBooleanToBooleanObject(node.Value)
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)

		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}

		return &object.Array{Elements: elements}
	case *ast.HashLiteral:
		return evalHashLiteral(node, env)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)

		if isError(right) {
			return right
		}

		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)

		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)

		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.ForExpression:
		return evalForExpression(node, env)
	case *ast.IndexExpression:
		left := Eval(node.Left, env)

		if isError(left) {
			return left
		}

		index := Eval(node.Index, env)

		if isError(index) {
			return index
		}

		return evalIndexExpression(left, index)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Body: body, Env: env}
	case *ast.CallExpression:
		function := Eval(node.Function, env)

		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env)

		// Return instantly if an error is encountered when evaluating the arguments
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args)
	}

	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()

			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

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

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBooleanToBooleanObject(left == right)
	case operator == "!=":
		return nativeBooleanToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

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

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
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
		return nativeBooleanToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBooleanToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBooleanToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBooleanToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)

	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func evalForExpression(ie *ast.ForExpression, env *object.Environment) object.Object {
	// TODO: Add variable reassignment and enclose loops in their own environments

	condition := Eval(ie.Condition, env)
	var result object.Object = NULL

	if isError(condition) {
		return condition
	}

	for {
		condition = Eval(ie.Condition, env)

		if isError(condition) {
			return condition
		}

		if !isTruthy(condition) {
			break
		}

		result = Eval(ie.Body, env)

		if isError(result) {
			return result
		}

		if rv, isReturnValue := result.(*object.ReturnValue); isReturnValue {
			return rv.Value
		}
	}

	return NULL
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: " + node.Value)
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)

		if isError(evaluated) {
			return []object.Object{evaluated}
		}

		result = append(result, evaluated)
	}

	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {

	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)

		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %s", fn.Type())
	}
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for pI, param := range fn.Parameters {
		env.Set(param.Value, fn.Env.IsMutable(param.Value), args[pI])
	}

	return env
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASHMAP_OBJ:
		return evalHashMapIndexExpression(left, index)
	case left.Type() == object.STRING_OBJ:
		return evalStringIndexExpression(left, index)
	}

	return newError("cannot index %T", left)
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arr := array.(*object.Array)
	idx, ok := index.(*object.Integer)

	if !ok {
		return newError("type of %s cannot be used to index %s", index.Type(), arr.Type())
	}

	max := int64(len(arr.Elements) - 1)

	if idx.Value == -1 {
		return arr.Elements[max]
	} else if idx.Value < 0 || idx.Value > max {
		return NULL
	}

	return arr.Elements[idx.Value]
}

func evalHashMapIndexExpression(hashMap, index object.Object) object.Object {
	hashObject := hashMap.(*object.HashMap)
	key, ok := index.(object.Hashable)

	if !ok {
		return newError("type of %s cannot be used as hash key", index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]

	if !ok {
		return NULL
	}

	return pair.Value
}

func evalStringIndexExpression(str, index object.Object) object.Object {
	strObject := str.(*object.String)
	idx, ok := index.(*object.Integer)

	if !ok {
		return newError("type of %s cannot be used to index %s", index.Type(), strObject.Type())
	}

	max := int64(len(strObject.Value) - 1)

	if idx.Value == -1 {
		return &object.String{Value: string(strObject.Value[max])}
	} else if idx.Value < 0 || idx.Value > max {
		return NULL
	}

	return &object.String{Value: string(strObject.Value[idx.Value])}
}

func evalDotExpression(left, fn object.Object, args []object.Object) object.Object {
	switch result := fn.(type) {
	case *object.Builtin:
		allArgs := append([]object.Object{left}, args...)
		return result.Fn(allArgs...)
	}

	return newError("%s does not exist on type %s", fn.Inspect(), left.Type())
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)

		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)

		if !ok {
			return newError("type of '%s' cannot be used as hash key", key.Type())
		}

		value := Eval(valueNode, env)

		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.HashMap{Pairs: pairs}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func nativeBooleanToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func newError(format string, a ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}

	return false
}
