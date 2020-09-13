package symbols

import "text/scanner"

/*
	A symbol that represents a Curly Bracket Environment
	E.g:
	{
		Hello there
	}
*/

const CurlyEnv SymbolType = "CurlyEnv"

type CurlyEnvSymbol struct {
	Statements []Symbol
	Position scanner.Position
}

func (c CurlyEnvSymbol) Pos() scanner.Position {
	return c.Position
}

func (c CurlyEnvSymbol) Type() SymbolType {
	return CurlyEnv
}

func (c CurlyEnvSymbol) Name() string {
	return "Curly Environment"
}

var _ Symbol = CurlyEnvSymbol{}
var _ Symbol = &CurlyEnvSymbol{}