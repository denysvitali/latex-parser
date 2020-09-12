package symbols

/*
	A symbol that represents a Macro
 */

const Macro SymbolType = "Macro"

type MacroSymbol struct {
	MacroName string
	CurlyArgs []string
	SquareArgs []string
}

func (m MacroSymbol) Type() SymbolType {
	return Macro
}

func (m MacroSymbol) Name() string {
	return m.MacroName
}

var _ Symbol = MacroSymbol{}
var _ Symbol = &MacroSymbol{}