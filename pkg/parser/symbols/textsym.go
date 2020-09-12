package symbols

/*
	A symbol that represents some text
*/

const Text SymbolType = "Text"

type TextSymbol struct {
	Content string
}

func (m TextSymbol) Type() SymbolType {
	return Text
}

func (m TextSymbol) Name() string {
	return "Text"
}

var _ Symbol = TextSymbol{}
var _ Symbol = &TextSymbol{}
