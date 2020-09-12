package symbols

const NewLine SymbolType = "NewLine"

type NewLineSymbol struct {
	NewLineName string
	CurlyArgs []Symbol
	SquareArgs []Symbol
}

func (m NewLineSymbol) Type() SymbolType {
	return NewLine
}

func (m NewLineSymbol) Name() string {
	return m.NewLineName
}

var _ Symbol = NewLineSymbol{}
var _ Symbol = &NewLineSymbol{}