package symbols

import (
	"text/scanner"
)

/*
	A new environment symbol
	\newenvironment{my-env}{\begin{env}}{\end{env}}
*/

const NewEnvironment SymbolType = "NewEnvironmentSymbol"

type NewEnvironmentSymbol struct {
	Environment string
	Statements  []Symbol
	Position    scanner.Position
	SquareArgs  []Symbol
}

func (m NewEnvironmentSymbol) Pos() scanner.Position {
	return m.Position
}

func (m NewEnvironmentSymbol) Type() SymbolType {
	return MacroText
}

func (m NewEnvironmentSymbol) Name() string {
	return "New Environment"
}

var _ Symbol = NewEnvironmentSymbol{}
var _ Symbol = &NewEnvironmentSymbol{}