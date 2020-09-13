package symbols

import "text/scanner"

/*
	A symbol that represents an inline math environment
	E.g: $\mathrm{x} = \sqrt{2}$
*/

const InlineMath SymbolType = "InlineMathSymbol"

type InlineMathSymbol struct {
	Statements []Symbol
	Position scanner.Position
}

func (i InlineMathSymbol) Pos() scanner.Position {
	return i.Position
}

func (i InlineMathSymbol) Type() SymbolType {
	return InlineMath
}

func (i InlineMathSymbol) Name() string {
	return "Inline Math"
}

var _ Symbol = InlineMathSymbol{}
var _ Symbol = &InlineMathSymbol{}