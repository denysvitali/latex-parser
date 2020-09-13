package symbols

import "text/scanner"

/*
	A symbol that represents some text
*/

const Text SymbolType = "Text"

type TextSymbol struct {
	Content string
	Position scanner.Position
}

func (t TextSymbol) Pos() scanner.Position {
	return t.Position
}

func (t TextSymbol) Type() SymbolType {
	return Text
}

func (t TextSymbol) Name() string {
	return "Text"
}

var _ Symbol = TextSymbol{}
var _ Symbol = &TextSymbol{}
