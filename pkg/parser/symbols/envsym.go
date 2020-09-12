package symbols

/*
	A symbol that represents an environment
 */

const Env SymbolType = "EnvSymbol"

type EnvSymbol struct {
	Environment string
	Statements []Symbol
}

func (e EnvSymbol) Type() SymbolType {
	return Env
}

func (e EnvSymbol) Name() string {
	return e.Environment
}

var _ Symbol = EnvSymbol{}
var _ Symbol = &EnvSymbol{}