package symbols

/*
	A symbol that represents a Curly Bracket Environment
	E.g:
	{
		Hello there
	}
*/

const CurlyEnv SymbolType = "CurlyEnv"

type CurlyEnvSymbol struct {
	Statements []Symbol
}

func (e CurlyEnvSymbol) Type() SymbolType {
	return CurlyEnv
}

func (e CurlyEnvSymbol) Name() string {
	return "Curly Environment"
}

var _ Symbol = CurlyEnvSymbol{}
var _ Symbol = &CurlyEnvSymbol{}