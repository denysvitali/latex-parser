package parser

import (
	"errors"
	"fmt"
	"github.com/denysvitali/latex-parser/pkg/parser/symbols"
	"github.com/denysvitali/latex-parser/pkg/tokenizer"
	"strings"
)

type Parser struct {
	tokenizer *tokenizer.Tokenizer
}

func New(tokenizer *tokenizer.Tokenizer) *Parser {
	return &Parser {
		tokenizer,
	}
}

func (p *Parser) Parse() ([]symbols.Symbol, error) {
	var statements []symbols.Symbol

	for ;; {
		currentToken := p.tokenizer.Peek()
		if currentToken.Type == tokenizer.EOF {
			p.tokenizer.Next()
			break
		}

		if currentToken.Type == tokenizer.Identifier {
			envMacro, err := p.envMacro()
			if err != nil {
				return statements, err
			}
			statements = append(statements, envMacro)
			continue
		}

		if currentToken.Type == tokenizer.Percent {
			// Comment, not added to the statements
			_, err := p.comment()
			if err != nil {
				return statements, err
			}
			continue
		}

		if currentToken.Type == tokenizer.Text {
			p.tokenizer.Next()
			textSymbol := symbols.TextSymbol{Content: currentToken.Value.(string)}
			statements = append(statements, textSymbol)
			continue
		}

		panic(fmt.Sprintf("unimplemented token %s", currentToken.Type))
	}

	return statements, nil
}

func (p *Parser) comment() (symbols.Symbol, error){
	var symbol symbols.Symbol
	if p.tokenizer.Next().Type != tokenizer.Percent {
		return symbol, errors.New("comment doesn't start with a %")
	}

	// Ignore everything until newline

	for ;; {
		// Find newline
		token := p.tokenizer.Next()

		if token.Type != tokenizer.Text {
			continue
		}

		// Is text
		if strings.Contains(token.Value.(string), "\n") {
			break
		}
	}

	return nil, nil
}

func (p *Parser) envMacro() (symbols.Symbol, error) {
	var statement symbols.Symbol
	/*
		An envMacro is usually a macro that can get executed and defined by the user.
		Identifiers are basically functions that return values and can therefore have arguments.
		These arguments are passed via curly brackets { } and square brackets [ ], for example:
		\newcommand{name}[num][default]{definition}

		Additionally, the special envMacro \begin{xxx} and \end{xxx} will define an "environment".
		These environments are usually used to run certain macros before and after the content
		defined between those two entries.

		Environments can additionally be nested. A simple example is:
			\begin{figure}[h!]
				\begin{center}
					\includegraphics{image.png}
					\caption{An image}
				\end{center}
			\end{figure}
	*/

	token := p.tokenizer.Next()

	if token.Type != tokenizer.Identifier {
		return statement, errors.New("token is not an envMacro")
	}

	if token.Name == "begin" {
		// Start of an environment (can contain multiple statements)
		envName, err := p.curlySingleArgument()
		if err != nil {
			return statement, fmt.Errorf("cannot parse begin curly argument: %v", err)
		}
		statement = symbols.EnvSymbol{
			Environment: envName,
		}

		err = p.endEnvironment(envName)
		if err != nil {
			return statement, err
		}

		return statement, nil
	}

	var curlyArguments []string

	if p.tokenizer.Peek().Type == tokenizer.OpenCurly {
		// The envMacro has some arguments, let's parse them
		var err error
		curlyArguments, err = p.curlyArguments()
		if err != nil {
			return statement, err
		}
	}

	var squareArgs []string

	if p.tokenizer.Peek().Type == tokenizer.OpenSquare {
		// The envMacro has some optional args
		var err error
		squareArgs, err = p.squareArguments()
		if err != nil {
			return statement, err
		}
	}

	macro := symbols.MacroSymbol {
		MacroName: token.Name,
		CurlyArgs: curlyArguments,
		SquareArgs: squareArgs,
	}

	return macro, nil
}

func (p *Parser) squareArguments() ([]string, error) {
	var squareArgs []string
	if p.tokenizer.Next().Type != tokenizer.OpenSquare {
		return squareArgs, errors.New("expected [")
	}

	textToken := p.tokenizer.Next()
	if textToken.Type != tokenizer.Text {
		return squareArgs, errors.New("expected text")
	}
	squareArgs = strings.Split(textToken.Value.(string), ",")

	if p.tokenizer.Next().Type != tokenizer.CloseSquare {
		return squareArgs, errors.New("expected ]")
	}

	return squareArgs, nil
}


func (p *Parser) curlyArguments() ([]symbols.Symbol, error) {
	var curlyArgs []symbols.Symbol
	if p.tokenizer.Next().Type != tokenizer.OpenCurly {
		return curlyArgs, errors.New("expected {")
	}

	for ;; {
		// Scan everything until } is found
		if p.tokenizer.Peek().Type == tokenizer.CloseCurly {
			p.tokenizer.Next()
			return curlyArgs, nil
		}

		if p.tokenizer.Peek().Type == tokenizer.Percent {
			_ , err := p.comment()
			if err != nil {
				return curlyArgs, err
			}
			continue
		}

		if p.tokenizer.Peek().Type == tokenizer.Identifier {
			// TODO: Execute Macro
			macro, _ := p.envMacro()
		}
	}

	if textToken.Type != tokenizer.Text {
		return curlyArgs, errors.New("expected text")
	}
	curlyArgs = strings.Split(textToken.Value.(string), ",")

	if p.tokenizer.Next().Type != tokenizer.CloseCurly {
		return curlyArgs, errors.New("expected }")
	}

	return curlyArgs, nil
}

func (p *Parser) curlySingleArgument() (string, error) {
	curlyArg := ""
	if p.tokenizer.Next().Type != tokenizer.OpenCurly {
		return curlyArg, errors.New("expected {")
	}

	textToken := p.tokenizer.Next()
	if textToken.Type != tokenizer.Text {
		return curlyArg, errors.New("expected text")
	}

	if p.tokenizer.Next().Type != tokenizer.CloseCurly {
		return curlyArg, errors.New("expected }")
	}

	return curlyArg, nil
}

func (p *Parser) endEnvironment(name string) error {
	endId := p.tokenizer.Next()

	if endId.Type != tokenizer.Identifier {
		return errors.New(fmt.Sprintf("expected envMacro, found %s", endId.Type))
	}

	if endId.Name != "end" {
		return errors.New(fmt.Sprintf("expected \\end, found \\%s", endId.Name))
	}

	envName, err := p.curlySingleArgument()
	if err != nil {
		return err
	}

	if envName != name {
		return errors.New(fmt.Sprintf("expected \\end{%s}, found \\end{%s}", name, envName))
	}

	return nil
}
