package symbols

import "text/scanner"

/*
	A symbol that represents a Macro
 */

const Macro SymbolType = "Macro"

type MacroSymbol struct {
	MacroName string
	CurlyArgs [][]Symbol
	SquareArgs []Symbol
	Position scanner.Position
}

func (m MacroSymbol) Pos() scanner.Position {
	return m.Position
}

func (m MacroSymbol) Type() SymbolType {
	return Macro
}

func (m MacroSymbol) Name() string {
	return m.MacroName
}

var _ Symbol = MacroSymbol{}
var _ Symbol = &MacroSymbol{}