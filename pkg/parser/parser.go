package parser

import (
	"errors"
	"fmt"
	"github.com/denysvitali/latex-parser/pkg/parser/symbols"
	"github.com/denysvitali/latex-parser/pkg/tokenizer"
	"path/filepath"
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

func (p *Parser) Parse() ([]symbols.Symbol, error){
	return p.parse(tokenizer.Unknown, "")
}

func (p *Parser) parse(tokenType tokenizer.TokenType, envName string) ([]symbols.Symbol, error) {
	var statements []symbols.Symbol

	for ;; {
		currentToken := p.tokenizer.Peek()
		if currentToken.Type == tokenizer.EOF {
			p.tokenizer.Next()
			break
		}

		if currentToken.Type == tokenizer.Identifier {
			if currentToken.Name == "end" {
				token := p.tokenizer.Next()
				// Check if this envName is the environment we're going to end
				env, err := p.curlySingleArgument()
				if err != nil {
					return statements, fmt.Errorf("cannot parse begin curly argument: %v", err)
				}
				symbol := symbols.EndEnvSymbol{
					Environment: env,
					Position: token.Pos,
				}
				statements = append(statements, symbol)
				if env == envName {
					return statements, nil
				}
				continue
			}

			if currentToken.Name == "newenvironment" {
				firstToken := p.tokenizer.Next()

				// Let's read first the environment name
				env, err := p.curlySingleArgument()
				if err != nil {
					return statements, fmt.Errorf("cannot parse begin curly argument: %v", err)
				}

				// Optional square arg:
				squareArgs, err := p.squareArguments()
				if err == nil {

				}

				_, _ = p.comment()
				p.skiptext()

				// We will interpret this as "MacroText"
				macroBegin, err := p.macroText(tokenizer.OpenCurly, tokenizer.CloseCurly)
				if err != nil {
					return statements, err
				}
				_, _ = p.comment()
				p.skiptext()

				macroEnd, err := p.macroText(tokenizer.OpenCurly, tokenizer.CloseCurly)
				if err != nil {
					return statements, err
				}

				newEnvironment := symbols.NewEnvironmentSymbol{
					Environment: env,
					Statements:  []symbols.Symbol{macroBegin, macroEnd},
					SquareArgs: squareArgs,
					Position:    firstToken.Pos,
				}

				statements = append(statements, newEnvironment)
				continue
			}

			_,_ = p.comment()
			envMacro, err := p.envMacro()
			if err != nil {
				return statements, err
			}
			statements = append(statements, envMacro)
			continue
		}

		if currentToken.Type == tokenizer.Dollar {

			if tokenType == tokenizer.Dollar {
				return statements, nil
			}

			// Inline Math Mode
			symbol, err := p.inlineMath()
			if err != nil {
				return statements, err
			}
			statements = append(statements, symbol)
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

		if currentToken.Type == tokenizer.NewLine {
			p.tokenizer.Next()
			newLineSymbol := symbols.NewLineSymbol{}
			statements = append(statements, newLineSymbol)
			continue
		}

		if currentToken.Type == tokenizer.OpenCurly {
			p.tokenizer.Next()
			parsedStatements, err := p.parse(tokenizer.CloseCurly, "")
			if err != nil {
				return statements, err
			}

			next := p.tokenizer.Next()
			if next.Type != tokenizer.CloseCurly {
				return statements, fmt.Errorf("invalid token: } expected but %v found", next)
			}

			statements = append(statements, symbols.CurlyEnvSymbol {
				Statements: parsedStatements,
			})
			continue
		}

		if currentToken.Type == tokenizer.CloseCurly {
			if tokenType == tokenizer.CloseCurly {
				return statements, nil
			}
		}

		if currentToken.Type == tokenizer.CloseSquare {
			if tokenType == tokenizer.CloseSquare {
				return statements, nil
			}
		}

		panic(fmt.Sprintf("unimplemented token %s", currentToken.Type))
	}

	return statements, nil
}

func (p *Parser) comment() (symbols.Symbol, error){
	var symbol symbols.Symbol
	if p.tokenizer.Peek().Type != tokenizer.Percent {
		return symbol, errors.New("comment doesn't start with a %")
	}
	p.tokenizer.Next()
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

func (p *Parser) macro() (symbols.Symbol, error) {
	token := p.tokenizer.Next()

	var statement symbols.Symbol
	var squareArgs []symbols.Symbol

	// Loop because we can have multiple (2?) arguments to a macro
	// e.g: \dfrac{1}{x}
	var curlyArguments [][]symbols.Symbol
	for i:=0; i<5; i++ {
		if p.tokenizer.Peek().Type == tokenizer.OpenSquare {
			// The envMacro has some optional args
			var err error
			squareArgs, err = p.squareArguments()
			if err != nil {
				return statement, err
			}
			continue
		}

		if p.tokenizer.Peek().Type == tokenizer.OpenCurly {
			// The envMacro has some arguments, let's parse them
			var err error
			cArg, err := p.curlyArguments()
			if err != nil {
				return statement, err
			}

			curlyArguments = append(curlyArguments, cArg)
			continue
		}
	}

	macro := symbols.MacroSymbol {
		MacroName: token.Name,
		CurlyArgs: curlyArguments,
		SquareArgs: squareArgs,
	}

	return macro, nil
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

	token := p.tokenizer.Peek()

	if token.Type != tokenizer.Identifier {
		return statement, errors.New("token is not an envMacro")
	}

	if token.Name == "begin" {
		// Start of an environment (can contain multiple statements)
		token = p.tokenizer.Next()
		envName, err := p.curlySingleArgument()
		if err != nil {
			return statement, fmt.Errorf("cannot parse begin curly argument: %v", err)
		}

		var squareArgs []symbols.Symbol
		if p.tokenizer.Peek().Type == tokenizer.OpenSquare {
			// The envMacro has some optional args
			var err error
			squareArgs, err = p.squareArguments()
			if err != nil {
				return statement, err
			}
		}

		var s []symbols.Symbol

		switch envName{
		case "equation", "align", "equation*", "align*":
			var equationSymbol symbols.Symbol
			equationSymbol, err = p.equationEnv(envName)
			s = []symbols.Symbol{equationSymbol}
		default:
			s, err = p.parse(tokenizer.Identifier, envName)
		}


		if err != nil {
			return statement, fmt.Errorf("unable to parse inside of \\begin{%s}: %v", envName, err)
		}

		statement = symbols.EnvSymbol{
			Environment: envName,
			Statements: s,
			SquareArgs: squareArgs,
			Position: token.Pos,
		}
		return statement, nil
	}
	if token.Name == "include" || token.Name == "input"{
		// Include files. I guess this should be part of the interpreter instead, but since I need a complete
		// AST for my use case, I included it here 🤷
		token = p.tokenizer.Next()

		basePath := filepath.Dir(p.tokenizer.OriginalFilePath)
		fileName, err := p.curlySingleArgument()
		if err != nil {
			return statement, fmt.Errorf("unable to find included file %s", fileName)
		}

		if token.Name == "include" {
			fileName = fileName + ".tex"
		}
		finalPath := filepath.Join(basePath, fileName)

		tkz2, err := tokenizer.Open(finalPath)
		if err != nil {
			return statement, fmt.Errorf("unable to initialize nested tokenizer: %v", err)
		}

		p2 := New(tkz2)
		symb, err := p2.Parse()

		return symbols.IncludeSymbol{
			Path: fileName,
			Statements: symb,
			Position: token.Pos,
		}, nil
	}

	return p.macro()
}

func (p *Parser) squareArguments() ([]symbols.Symbol, error) {
	var args []symbols.Symbol
	var err error
	if p.tokenizer.Next().Type != tokenizer.OpenSquare {
		return args, errors.New("expected [")
	}

	args, err = p.argBody(tokenizer.CloseSquare)
	if err != nil {
		return args, err
	}

	if p.tokenizer.Next().Type != tokenizer.CloseSquare {
		return args, errors.New("expected ]")
	}

	return args, nil
}


func (p *Parser) curlyArguments() ([]symbols.Symbol, error) {
	if p.tokenizer.Next().Type != tokenizer.OpenCurly {
		return []symbols.Symbol{}, errors.New("expected {")
	}

	args, err := p.argBody(tokenizer.CloseCurly)
	if err != nil {
		return args, err
	}

	if p.tokenizer.Next().Type != tokenizer.CloseCurly {
		return args, errors.New("expected }")
	}

	return args, nil
}

func(p *Parser) argBody(stopWith tokenizer.TokenType) ([]symbols.Symbol, error){
	return p.parse(stopWith, "")
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

	curlyArg = textToken.Value.(string)

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

func (p *Parser) equationEnv(envName string) (symbols.Symbol, error) {
	var symbol symbols.Symbol
	var statements []symbols.Symbol

	peek := p.tokenizer.Peek()
	if peek.Type != tokenizer.Identifier {
		return symbol, fmt.Errorf("invalid token %v found, Identifier expected", peek.Type)
	}

	// Body
	shouldStop := false
	for ;; {
		peek := p.tokenizer.Peek()

		switch peek.Type {
		case tokenizer.Dollar:
			p.tokenizer.Next()
			shouldStop = true

		case tokenizer.Identifier:
			if peek.Name == "end" {
				// Handle end of this env
				token := p.tokenizer.Next()
				t := p.tokenizer.Next()
				if t.Type != tokenizer.OpenCurly {
					return symbol, fmt.Errorf("expected token {, fourd %v at %v", t.Type, t.Pos)
				}
				t = p.tokenizer.Next()
				if t.Type != tokenizer.Text {
					return symbol, fmt.Errorf("expected text, fourd %v at %v", t.Type, t.Pos)
				}
				currentEnv := t.Value.(string)
				if currentEnv != envName {
					return symbol, fmt.Errorf("wrong environment being ended: expected %v but fourd %v at %v",
						envName, currentEnv, t.Pos)
				}
				t = p.tokenizer.Next()
				if t.Type != tokenizer.CloseCurly {
					return symbol, fmt.Errorf("expected token }, fourd %v at %v", t.Type, t.Pos)
				}

				endEnv := symbols.EndEnvSymbol{
					Environment: currentEnv,
					Position: token.Pos,
				}
				statements = append(statements, endEnv)

			}
			statement, err := p.macro()
			if err != nil {
				return symbol, err
			}

			statements = append(statements, statement)
			continue
		case tokenizer.Text:
			token := p.tokenizer.Next()
			txtSymb := symbols.TextSymbol{
				Content: token.Value.(string),
				Position: token.Pos,
			}
			statements = append(statements, txtSymb)
			continue
		case tokenizer.OpenSquare:
			token := p.tokenizer.Next()
			txtSymb := symbols.TextSymbol{
				Content: "[",
				Position: token.Pos,
			}
			statements = append(statements, txtSymb)
			continue
		case tokenizer.CloseSquare:
			token := p.tokenizer.Next()
			txtSymb := symbols.TextSymbol{
				Content: "]",
				Position: token.Pos,
			}
			statements = append(statements, txtSymb)
			continue
		default:
			panic(fmt.Sprintf("%v in InlineMath is unimplemented", peek))
		}

		if shouldStop {
			break
		}
	}

	return symbols.InlineMathSymbol{Statements: statements}, nil
}

/*
	Inline Math

	Anything delimited by $...$ is considered inline math. It behaves like a \begin{env}...\end{env}
	but its contents are not parsed by a general math parser. The content of the Inline Math mode
	is different from the classical LaTeX environment. This environment is used to write equations, therefore
	the result of this mode in our parser would be an EquationSymbol.
 */
func (p *Parser) inlineMath() (symbols.Symbol, error) {
	var symbol symbols.Symbol
	var statements []symbols.Symbol

	firstToken := p.tokenizer.Next()
	if firstToken.Type != tokenizer.Dollar {
		return symbol, errors.New("unexpected firstToken %s, $ expected")
	}

	// Body
	shouldStop := false
	for ;; {
		peek := p.tokenizer.Peek()

		switch peek.Type {
		case tokenizer.Dollar:
			p.tokenizer.Next()
			shouldStop = true

		case tokenizer.Identifier:
			statement, err := p.macro()
			if err != nil {
				return symbol, err
			}

			statements = append(statements, statement)
			continue
		case tokenizer.Text:
			token := p.tokenizer.Next()
			txtSymb := symbols.TextSymbol{
				Content: token.Value.(string),
				Position: token.Pos,
			}
			statements = append(statements, txtSymb)
			continue
		case tokenizer.OpenSquare:
			token := p.tokenizer.Next()
			txtSymb := symbols.TextSymbol{
				Content: "[",
				Position: token.Pos,
			}
			statements = append(statements, txtSymb)
			continue
		case tokenizer.CloseSquare:
			token := p.tokenizer.Next()
			txtSymb := symbols.TextSymbol{
				Content: "]",
				Position: token.Pos,
			}
			statements = append(statements, txtSymb)
			continue
		default:
			panic(fmt.Sprintf("%v in InlineMath is unimplemented", peek))
		}

		if shouldStop {
			break
		}
	}

	return symbols.InlineMathSymbol{
		Statements: statements,
		Position: firstToken.Pos,
	}, nil
}

func (p *Parser) macroText(startsWith tokenizer.TokenType, endsWith tokenizer.TokenType) (symbols.Symbol, error){
	// A text that will later on be parsed by the Macro Executor.
	// A Macro Text can contain anything.

	var symbol symbols.Symbol
	firstToken := p.tokenizer.Peek()
	if firstToken.Type != startsWith {
		return symbol, fmt.Errorf("unexpected token %v, %v was expected at %v",
			firstToken.Type,
			startsWith,
			firstToken.Pos,
		)
	}
	p.tokenizer.Next()

	var tokens []tokenizer.Token

	// Macro Text
	shouldStop := false
	seenStartTokens := 0
	for ;; {
		t := p.tokenizer.Peek()

		switch t.Type {
		case endsWith:
			tokens = append(tokens, t)
			if seenStartTokens == 0 {
				p.tokenizer.Next()
				shouldStop = true
			} else {
				seenStartTokens--
				p.tokenizer.Next()
				tokens = append(tokens, t)
			}
		case tokenizer.Percent:
			// Skip comments
			for ;; {
				t := p.tokenizer.Next()
				if t.Type == tokenizer.Text {
					if strings.Contains(t.Value.(string), "\n") {
						break
					}
				}
			}
		case startsWith:
			p.tokenizer.Next()
			seenStartTokens++
		default:
			p.tokenizer.Next()
			tokens = append(tokens, t)
		}

		if shouldStop {
			break
		}
	}

	return symbols.MacroTextSymbol{
		Tokens: tokens,
		Position: firstToken.Pos,
	}, nil
}

func (p *Parser) skiptext() {
	for ;; {
		peek := p.tokenizer.Peek()
		if peek.Type == tokenizer.Text {
			p.tokenizer.Next()
			continue
		}
		break
	}
}
