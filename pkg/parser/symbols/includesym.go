package symbols

const Include SymbolType = "Include"

type IncludeSymbol struct {
	Path string
	Statements []Symbol
}

func (i IncludeSymbol) Type() SymbolType {
	return Include
}

func (i IncludeSymbol) Name() string {
	return "Include(" + i.Path + ")"
}

var _ Symbol = IncludeSymbol{}
var _ Symbol = &IncludeSymbol{}