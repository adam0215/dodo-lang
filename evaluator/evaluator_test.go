package evaluator

import (
	"bytes"
	"dodo-lang/lexer"
	"dodo-lang/object"
	"dodo-lang/parser"
	"os"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"!5", false},
		{"!10", false},
		{"!!5", true},
		{"!-5", false},
		{"!!-5", true},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestEvalStringExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"foobar"`, "foobar"},
		{`"hello world!"`, "hello world!"},
		{`"5"`, "5"},
		{`"-5"`, "-5"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testStringObject(t, evaluated, tt.expected)
	}
}

func TestStringConcatenation(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"foo" + "bar"`, "foobar"},
		{`"foo" + "bar" + "baz"`, "foobarbaz"},
		{`"foo " + "bar!"`, "foo bar!"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testStringObject(t, evaluated, tt.expected)
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}

	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d",
			len(result.Elements))
	}

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!!true", true},
		{"!false", true},
		{"!!false", false},
		{"!5", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { return 10 }", 10},
		{"if (false) { return 10 }", nil},
		{"if (1) { return 10 }", 10},
		{"if (1 < 2) { return 10 }", 10},
		{"if (1 > 2) { return 10 }", nil},
		{"if (1 > 2) { return 10 } else { return 20 }", 20},
		{"if (1 < 2) { return 10 } else { return 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestForExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`let mut count = 0;

		  for (count < 10) {
		    count = count + 1;
		  }

		  count;
			`, 10},
		{`let mut count = 0;

		  for (count < 10) {
		    count = count + 1;

			if (count == 5) {
			  return;
			}
		  }

		  count;
			`, 5},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)

		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"let i = 0; [1][i];",
			1,
		},
		{
			"[1, 2, 3][1 + 1];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]",
			2,
		},
		{
			"[1, 2, 3][3]",
			nil,
		},
		{
			// -1 returns last element in array
			"[1, 2, 3][-1]",
			3,
		},
		{`"hello world"[2]`, "l"},
		{
			`let myStr = "foobar"; let i = 1; myStr[i]`,
			"o",
		},
		{
			"[1, 2, 3].0",
			1,
		},
		{
			`"hello world".2`,
			"l",
		},
		{
			`let str = "hello world"; (str.2) + (str.4) + (str.9);`,
			"lol",
		},
		{
			"let myArray = [1, 2, 3]; myArray.2;",
			3,
		},
		{
			"let myArray = [1, 2, 3]; (myArray.0) + (myArray.1) + (myArray.2);",
			6,
		},
		{
			"let myArray = [1, 2, 3]; let i = myArray.0; myArray.i",
			2,
		},
		{
			"[1, 2, 3].3",
			nil,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			testStringObject(t, evaluated, expected)
		default:
			testNullObject(t, evaluated)
		}
	}
}

func TestDotExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"hello world".len()`, 11},
		{`"hello world".first()`, "h"},
		{`"hello world".last()`, "d"},
		{`"hello world".rest()`, "ello world"},
		{`[1, 2, 3].push(4)`, "[1, 2, 3, 4]"},
		{`[1, 2, 3].len()`, 3},
		{`1.len`, "cannot index *object.Integer"},
		{`"hello world".doesnotexist()`, "identifier not found: doesnotexist"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			switch result := evaluated.(type) {
			case *object.Error:

				if result.Message != expected {
					t.Errorf("wrong error message. expected=%q, got=%q", expected, result.Message)
				}
			case *object.String:
				testStringObject(t, evaluated, expected)
			case *object.Array:
				testArrayObject(t, evaluated, expected)
			}
		default:
			testNullObject(t, evaluated)
		}
	}
}

func TestPipeExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`let add = fn(x, y) {x + y};
		  let result = 5 |> add(10, $);

		  result;`, 15},
		{`let add = fn(x, y) {x + y};
		  let sub = fn(x, y) {x - y};
		  let result = sub(10, 3) |> add($, 10);

		  result;`, 17},
		{`"hello".len() |> push([1, 2, 3, 4], $)`, "[1, 2, 3, 4, 5]"},
		// TODO: Add support for using placeholder in function invokations using dot expressions
		// {`4 |> [1, 2, 3].push($)`, "[1, 2, 3, 4]"},
		// {`1.len`, "argument to `len` not supported, got INTEGER"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			switch result := evaluated.(type) {
			case *object.Error:

				if result.Message != expected {
					t.Errorf("wrong error message. expected=%q, got=%q", expected, result.Message)
				}
			case *object.String:
				testStringObject(t, evaluated, expected)
			case *object.Array:
				testArrayObject(t, evaluated, expected)
			}
		default:
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{"if (10 > 1) { return 10; }", 10},
		{
			`
if (10 > 1) {
  if (10 > 1) {
    return 10;
  }

  return 1;
}
`,
			10,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"true + false + true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`"foo" - "bar"`,
			"unknown operator: STRING - STRING",
		},
		{
			`
if (10 > 1) {
  if (10 > 1) {
    return true + false;
  }

  return 1;
}
`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{"foobar", "identifier not found: foobar"},
		{`{"name": "Dodo"}[fn(x) { x }];`, "type of FUNCTION cannot be used as hash key"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		testError(t, evaluated, tt.expectedMessage)
	}
}

func testError(t *testing.T, obj object.Object, expectedError string) bool {
	errObj, ok := obj.(*object.Error)

	if !ok {
		t.Errorf("no error object returned. got=%T(%+v)",
			obj, obj)
		return false
	}

	if errObj.Message != expectedError {
		t.Errorf("wrong error message. expected=%q, got=%q",
			expectedError, errObj.Message)
		return false
	}

	return true
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
		{"let mut a = 5; a;", 5},
		{"let mut a = 5 * 5; a;", 25},
		{"let mut a = 5; let b = a; b;", 5},
		{"let mut a = 5; let b = a; let c = a + b + 5; c;", 15},
		{"let mut a = 5; let a = 5; a;", "identifier 'a' already exists"},
		{"let a = 5; let a = 5; a;", "identifier 'a' already exists"},
	}

	for _, tt := range tests {
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, testEval(tt.input), int64(expected))
		case string:
			testError(t, testEval(tt.input), expected)
		}
	}
}

func TestReassignmentStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"let mut a = 3; a = 5; a;", 5},
		{"let mut a = 3; let b = a; a = a + b + 5; a;", 11},
		{"let a = 5; a = 3; a;", "identifier 'a' is not mutable"},
	}

	for _, tt := range tests {
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, testEval(tt.input), int64(expected))
		case string:
			testError(t, testEval(tt.input), string(expected))
		}
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v",
			fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"

	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
	let newAdder = fn(x) { fn(y) { x + y }; };
	let addTwo = newAdder(2);

	addTwo(2);`

	testIntegerObject(t, testEval(input), 4)
}

func TestBuiltinFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len("")`, 0},
		{`len([1, 2, 3])`, 3},
		{`len([])`, 0},
		{`len(1)`, "argument to `len` not supported, got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. got=2, expected=1"},
		{`rest([1, 2, 3])`, "[2, 3]"},
		{`rest("hello world")`, "ello world"},
		{`rest([])`, NULL},
		{`rest("")`, NULL},
		{`first([1, 2, 3])`, 1},
		{`first([])`, NULL},
		{`first("")`, NULL},
		{`last([1, 2, 3])`, 3},
		{`last([])`, NULL},
		{`last("")`, NULL},
		{`push([1, 2, 3], 4);`, "[1, 2, 3, 4]"},
		{`let myArray = [1, 2, 3];
		  let newArray = push(myArray, 4);
		  newArray;`, "[1, 2, 3, 4]"},
		{`typeof("hello world")`, "STRING"},
		{`typeof(9)`, "INTEGER"},
		{`typeof(fn (x) { 420; })`, "FUNCTION"},
		{`typeof(len)`, "BUILTIN"},
		{`typeof("one", "two")`, "wrong number of arguments. got=2, expected=1"},
		{`debug("hello world")`, "hello world\n"},
		{`debug(9)`, "9\n"},
	}

	for _, tt := range tests {
		pipeReader, pipeWriter, _ := os.Pipe()
		os.Stdout = pipeWriter

		evaluated := testEval(tt.input)

		pipeWriter.Close()

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			switch result := evaluated.(type) {
			case *object.Error:
				if result.Message != expected {
					t.Errorf("wrong error message. expected=%q, got=%q", expected, result.Message)
				}
			case *object.String:
				testStringObject(t, result, string(expected))
			case *object.Array:
				testArrayObject(t, result, expected)
			case *object.Null:
				var buf bytes.Buffer
				buf.ReadFrom(pipeReader)
				capturedStd := buf.String()

				if capturedStd != expected {
					t.Errorf("wrong std message. expected=%q, got=%q", expected, capturedStd)
				}
			}
		}
	}
}

func TestHashLiterals(t *testing.T) {
	input := `let two = "two";
	{
		"one": 10 - 9,
		two: 1 + 1,
		"thr" + "ee": 6 / 2,
		4: 4,
		true: 5,
		false: 6
	}`

	evaluated := testEval(input)
	result, ok := evaluated.(*object.HashMap)
	if !ok {
		t.Fatalf("Eval didn't return Hash. got=%T (%+v)", evaluated, evaluated)
	}

	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		TRUE.HashKey():                             5,
		FALSE.HashKey():                            6,
	}

	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong num of pairs. got=%d", len(result.Pairs))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}

		testIntegerObject(t, pair.Value, expectedValue)
	}
}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`{"foo": 5}["foo"]`,
			5,
		},
		{
			`{"foo": 5}["bar"]`,
			nil,
		},
		{
			`let key = "foo"; {"foo": 5}[key]`,
			5,
		},
		{
			`{}["foo"]`,
			nil,
		},
		{
			`{5: 5}[5]`,
			5,
		},
		{
			`{true: 5}[true]`,
			5,
		},
		{
			`{false: 5}[false]`,
			5,
		},
		//
		{
			`let key = "foo"; {"foo": 5}.key`,
			5,
		},
		{
			`let map = {"foo": 5, true: 3}; (map."foo") + (map.true)`,
			8,
		},
		{
			`{}."foo"`,
			nil,
		},
		{
			`{5: 5}.5`,
			5,
		},
		{
			`{false: 5}.false`,
			5,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)

	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, expected=%d", result.Value, expected)
		return false
	}

	return true
}

func testStringObject(t *testing.T, obj object.Object, expected string) bool {
	result, ok := obj.(*object.String)

	if !ok {
		t.Errorf("object is not String. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%s, expected=%s", result.Value, expected)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)

	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, expected=%t", result.Value, expected)
		return false
	}

	return true
}

func testArrayObject(t *testing.T, obj object.Object, expected string) bool {
	result, ok := obj.(*object.Array)

	if !ok {
		t.Errorf("object is not Array. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Inspect() != expected {
		t.Errorf("object has wrong form. got=%s, expected=%s", result.Inspect(), expected)
		return false
	}

	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}
