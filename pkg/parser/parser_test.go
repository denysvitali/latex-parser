package parser_test

import (
	"github.com/denysvitali/latex-parser/pkg/parser"
	"github.com/denysvitali/latex-parser/pkg/tokenizer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParser(t *testing.T){
	tkn, err := tokenizer.Open("../../resources/test/1.tex")
	assert.Nil(t, err)

	p := parser.New(tkn)
	ast, err := p.Parse()

	assert.NotNil(t, ast)
}
