package parser_test

import (
	"fmt"
	"github.com/denysvitali/latex-parser/pkg/parser"
	"github.com/denysvitali/latex-parser/pkg/parser/symbols"
	"github.com/denysvitali/latex-parser/pkg/tokenizer"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestParser(t *testing.T){
	tkn, err := tokenizer.Open("../../resources/test/1.tex")
	assert.Nil(t, err)

	p := parser.New(tkn)
	ast, err := p.Parse()

	assert.Nil(t, err)
	assert.NotNil(t, ast)

	parseAst(ast, 0)
}

func parseAst(ast []symbols.Symbol, depth int) {
	for _, v := range ast {
		switch v.Type() {
		case symbols.Env:
			envSymbol := v.(symbols.EnvSymbol)
			fmt.Printf("%s%s - %s (env: %s)\n", tabs(depth), v.Type(), v.Name(), envSymbol.Environment)
			parseAst(envSymbol.Statements, depth + 1)
		case symbols.Macro:
			macroSymbol := v.(symbols.MacroSymbol)
			fmt.Printf("%s%s - %s (s: %v, c: %v)\n", tabs(depth), v.Type(), v.Name(),
				macroSymbol.SquareArgs,
				macroSymbol.CurlyArgs)
		case symbols.Text:
			textSymbol := v.(symbols.TextSymbol)
			fmt.Printf("%s%s - %s (%s)\n", tabs(depth), v.Type(), v.Name(), textSymbol.Content)
		}
	}
}

func tabs(depth int) string {
	var result []string
	for i := 0; i< depth; i++ {
		result = append(result, "\t")
	}

	return strings.Join(result, "")
}