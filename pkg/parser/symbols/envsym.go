package symbols

import "text/scanner"

/*
	A symbol that represents an environment
 */

const Env SymbolType = "EnvSymbol"

type EnvSymbol struct {
	Environment string
	Statements []Symbol
	SquareArgs []Symbol
	Position scanner.Position
}

func (e EnvSymbol) Pos() scanner.Position {
	return e.Position
}

func (e EnvSymbol) Type() SymbolType {
	return Env
}

func (e EnvSymbol) Name() string {
	return e.Environment
}

var _ Symbol = EnvSymbol{}
var _ Symbol = &EnvSymbol{}