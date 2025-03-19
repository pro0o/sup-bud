package main

import (
	"olaf/eval"
	"olaf/lexer"
	"olaf/object"
	"olaf/parser"
	"syscall/js"
)

func main() {
	// Define the eval function that will be exposed to JavaScript
	js.Global().Set("evaluateOlaf", js.FuncOf(evaluateOlaf))

	// Keep the program running
	select {}
}

// evaluateOlaf - JS-callable function to evaluate Olaf code
func evaluateOlaf(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return map[string]interface{}{
			"error": "No code provided",
		}
	}

	// Get code from JS
	code := args[0].String()

	// Create environment
	env := object.NewEnvironment()

	// Create lexer, parser
	l := lexer.New(code)
	p := parser.New(l)

	// Parse program
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		var errorStr string
		for _, msg := range p.Errors() {
			errorStr += msg + "\n"
		}
		return map[string]interface{}{
			"error": errorStr,
		}
	}

	// Evaluate program
	evaluated := eval.Eval(program, env)

	if evaluated == nil {
		return map[string]interface{}{
			"result": "null",
		}
	}

	return map[string]interface{}{
		"result": evaluated.Inspect(),
	}
}
