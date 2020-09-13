package symbols

import "text/scanner"

const NewLine SymbolType = "NewLine"

type NewLineSymbol struct {
	Position scanner.Position
}

func (m NewLineSymbol) Pos() scanner.Position {
	return m.Position
}

func (m NewLineSymbol) Type() SymbolType {
	return NewLine
}

func (m NewLineSymbol) Name() string {
	return "New Line"
}

var _ Symbol = NewLineSymbol{}
var _ Symbol = &NewLineSymbol{}