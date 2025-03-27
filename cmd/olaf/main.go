package main

import (
	"olaf/eval"
	"olaf/lexer"
	"olaf/object"
	"olaf/parser"
	"strings"
	"syscall/js"
	"time"
)

func main() {
	js.Global().Set("evaluateOlaf", js.FuncOf(evaluateOlaf))
	select {}
}

func evaluateOlaf(this js.Value, args []js.Value) interface{} {
	evalOptions := eval.EvalOptions{
		MaxDepth: 200,
		Timeout:  5 * time.Second,
	}

	if len(args) < 1 {
		return map[string]interface{}{
			"error": "No code provided",
		}
	}

	code := args[0].String()
	l := lexer.New(code)

	var errors []string

	if l.HasErrors() {
		errors = append(errors, l.Errors()...)
	}

	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		errors = append(errors, p.Errors()...)
	}

	if len(errors) > 0 {
		var errorBuilder strings.Builder
		lexerErrors := l.FormatErrors()
		parserErrors := p.FormatErrors()

		if lexerErrors != "" {
			errorBuilder.WriteString(lexerErrors)
		}

		if parserErrors != "" {
			if lexerErrors != "" {
				errorBuilder.WriteString("\n")
			}
			errorBuilder.WriteString(parserErrors)
		}

		return map[string]interface{}{
			"error": errorBuilder.String(),
		}
	}

	env := object.NewEnvironment()

	evaluated := eval.EvalWithOptions(program, env, evalOptions)

	if evaluated == nil {
		return map[string]interface{}{
			"result": "null",
		}
	}

	if errorObj, ok := evaluated.(*object.Error); ok {
		return map[string]interface{}{
			"error": errorObj.Message,
		}
	}

	return map[string]interface{}{
		"result": evaluated.Inspect(),
	}
}
