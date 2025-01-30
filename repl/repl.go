package repl

import (
	"fmt"
	"io"
	"olaf/eval"
	"olaf/lexer"
	"olaf/object"
	"olaf/parser"
	"strings"

	"github.com/chzyer/readline"
)

const (
	PROMPT       = ">> "
	CONTINUATION = "... "
)

type REPL struct {
	env             *object.Environment
	history         []string
	multilineBuffer []string
	rl              *readline.Instance
}

func NewREPL() (*REPL, error) {
	rl, err := readline.New(PROMPT)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize readline: %v", err)
	}

	return &REPL{
		env:             object.NewEnvironment(),
		history:         make([]string, 0),
		multilineBuffer: make([]string, 0),
		rl:              rl,
	}, nil
}

func (r *REPL) handleSpecialCommands(line string) bool {
	switch strings.TrimSpace(line) {
	case ":h":
		fmt.Println("Available commands:")
		fmt.Println(":h    - Show this help message")
		fmt.Println(":z - Show command history")
		fmt.Println(":c   - Clear the screen")
		fmt.Println(":q    - Exit the REPL")
		fmt.Println("Use \\ at the end of a line for multi-line input")
		return true
	case ":z":
		for i, cmd := range r.history {
			fmt.Printf("%d: %s\n", i+1, cmd)
		}
		return true
	case ":c":
		readline.ClearScreen(r.rl)
		return true
	case ":q":
		r.rl.Close()
		return true
	}
	return false
}

func (r *REPL) Start(out io.Writer) {
	defer r.rl.Close()

	for {
		line, err := r.rl.Readline()
		if err != nil {
			if err == io.EOF {
				return
			}
			fmt.Fprintf(out, "Error reading input: %v\n", err)
			continue
		}

		// Handle special commands
		if strings.HasPrefix(line, ":") {
			if r.handleSpecialCommands(line) {
				continue
			}
		}

		// Handle multi-line input
		if strings.HasSuffix(strings.TrimSpace(line), "\\") {
			r.multilineBuffer = append(r.multilineBuffer, strings.TrimSuffix(strings.TrimSpace(line), "\\"))
			r.rl.SetPrompt(CONTINUATION)
			continue
		}

		// Process the complete input
		if len(r.multilineBuffer) > 0 {
			r.multilineBuffer = append(r.multilineBuffer, line)
			line = strings.Join(r.multilineBuffer, "\n")
			r.multilineBuffer = make([]string, 0)
			r.rl.SetPrompt(PROMPT)
		}

		if strings.TrimSpace(line) == "" {
			continue
		}

		// Add to history
		r.history = append(r.history, line)

		// Parse and evaluate
		l := lexer.New(line)
		p := parser.New(l)
		p.SetDebugMode(true)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := eval.Eval(program, r.env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

// Start creates and starts a new REPL instance
func Start(in io.Reader, out io.Writer) {
	repl, err := NewREPL()
	if err != nil {
		fmt.Fprintf(out, "Error initializing REPL: %v\n", err)
		return
	}
	repl.Start(out)
}
