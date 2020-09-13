package symbols

import (
	"text/scanner"
)

type Symbol interface {
	Type() SymbolType
	Name() string
	Pos() scanner.Position
}

type SymbolType string
