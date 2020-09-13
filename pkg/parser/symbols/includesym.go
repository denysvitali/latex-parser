package symbols

import "text/scanner"

const Include SymbolType = "Include"

type IncludeSymbol struct {
	Path string
	Statements []Symbol
	Position scanner.Position
}

func (i IncludeSymbol) Pos() scanner.Position {
	return i.Position
}

func (i IncludeSymbol) Type() SymbolType {
	return Include
}

func (i IncludeSymbol) Name() string {
	return "Include(" + i.Path + ")"
}

var _ Symbol = IncludeSymbol{}
var _ Symbol = &IncludeSymbol{}