package symbols

/*
	A symbol that represents an inline math environment
	E.g: $\mathrm{x} = \sqrt{2}$
*/

const InlineMath SymbolType = "InlineMathSymbol"

type InlineMathSymbol struct {
	Statements []Symbol
}

func (e InlineMathSymbol) Type() SymbolType {
	return InlineMath
}

func (e InlineMathSymbol) Name() string {
	return "Inline Math"
}

var _ Symbol = InlineMathSymbol{}
var _ Symbol = &InlineMathSymbol{}