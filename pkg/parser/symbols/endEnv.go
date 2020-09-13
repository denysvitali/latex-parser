package symbols

import "text/scanner"

/*
	A symbol that represents a Curly Bracket Environment
	E.g:
	{
		Hello there
	}
*/

const EndEnv SymbolType = "EndEnv"

type EndEnvSymbol struct {
	Environment string
	Position scanner.Position
}

func (e EndEnvSymbol) Pos() scanner.Position {
	return e.Position
}

func (e EndEnvSymbol) Type() SymbolType {
	return EndEnv
}

func (e EndEnvSymbol) Name() string {
	return "End Environment"
}

var _ Symbol = EndEnvSymbol{}
var _ Symbol = &EndEnvSymbol{}