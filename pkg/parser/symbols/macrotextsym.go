package symbols

import (
	"github.com/denysvitali/latex-parser/pkg/tokenizer"
	"text/scanner"
)

/*
	A symbol that represents the content of a macro argument
	E.g: \newenvironment{\begin{env}}{\end{env}}

	This content (in this case, \begin{env} and \end{env}) won't be parsed - we'll keep a list
	of tokens in the struct so that at a later stage, we can parse the snippet when it is used - if we fancy.
*/

const MacroText SymbolType = "MacroTextSymbol"

type MacroTextSymbol struct {
	Tokens []tokenizer.Token
	Position scanner.Position
}

func (m MacroTextSymbol) Pos() scanner.Position {
	return m.Position
}

func (m MacroTextSymbol) Type() SymbolType {
	return MacroText
}

func (m MacroTextSymbol) Name() string {
	return "Macro Text"
}

var _ Symbol = MacroTextSymbol{}
var _ Symbol = &MacroTextSymbol{}