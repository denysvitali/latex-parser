package symbols

type Symbol interface {
	Type() SymbolType
	Name() string
}

type SymbolType string
