package tokenizer_test

import (
	"fmt"
	"github.com/denysvitali/latex-parser/pkg/tokenizer"
	"github.com/stretchr/testify/assert"
	"testing"
	"text/scanner"
)

func TestTokenizer1(t *testing.T){
	tk, err := tokenizer.Open("../../resources/test/tokenizer/1.tex")
	assert.Nil(t, err)
	assert.NotNil(t, tk)

	assert.Equal(t, tokenizer.Token{
		Type: tokenizer.Identifier,
		Name: "documentclass",
		Pos: scanner.Position{
			Filename: "",
			Offset:   0,
			Line:     1,
			Column:   1,
		},
	}, tk.Next())

	assert.Equal(t, tokenizer.Token{
		Type: tokenizer.OpenCurly,
		Pos: scanner.Position {
			Filename: "",
			Offset:   14,
			Line:     1,
			Column:   15,
		},
	}, tk.Next())

	assert.Equal(t, tokenizer.Token{
		Type: tokenizer.Text,
		Value: "article",
		Pos: scanner.Position{
			Filename: "",
			Offset:   15,
			Line:     1,
			Column:   16,
		},
	}, tk.Next())

	assert.Equal(t, tokenizer.Token{
		Type: tokenizer.CloseCurly,
		Pos: scanner.Position{
			Filename: "",
			Offset:   22,
			Line:     1,
			Column:   23,
		},
	}, tk.Next())

	assert.Equal(t, tokenizer.Token{
		Type: tokenizer.Text,
		Pos: scanner.Position{
			Filename: "",
			Offset:   23,
			Line:     1,
			Column:   24,
		},
		Value: " ",
	}, tk.Next())

	assert.Equal(t, tokenizer.Token{
		Type: tokenizer.Percent,
		Pos: scanner.Position{
			Filename: "",
			Offset:   24,
			Line:     1,
			Column:   25,
		},
	}, tk.Next())

	assert.Equal(t, tokenizer.Token{
		Type: tokenizer.Text,
		Pos: scanner.Position{
			Filename: "",
			Offset:   25,
			Line:     1,
			Column:   26,
		},
		Value: " Starts an article\n",
	}, tk.Next())

	assert.Equal(t, tokenizer.Token{
		Type: tokenizer.Identifier,
		Pos: scanner.Position{
			Filename: "",
			Offset:   44,
			Line:     2,
			Column:   1,
		},
		Name: "usepackage",
	}, tk.Next())

	assert.Equal(t, tokenizer.Token{
		Type: tokenizer.OpenCurly,
		Pos: scanner.Position{
			Filename: "",
			Offset:   55,
			Line:     2,
			Column:   12,
		},
	}, tk.Next())

	assert.Equal(t, tokenizer.Token{
		Type: tokenizer.Text,
		Pos: scanner.Position{
			Filename: "",
			Offset:   56,
			Line:     2,
			Column:   13,
		},
		Value: "amsmath",
	}, tk.Next())

	assert.Equal(t, tokenizer.Token{
		Type: tokenizer.CloseCurly,
		Pos: scanner.Position{
			Filename: "",
			Offset:   63,
			Line:     2,
			Column:   20,
		},
	}, tk.Next())

	assert.Equal(t, tokenizer.Token{
		Type: tokenizer.Text,
		Pos: scanner.Position{
			Filename: "",
			Offset:   64,
			Line:     2,
			Column:   21,
		},
		Value: " ",
	}, tk.Next())

	assert.Equal(t, tokenizer.Token{
		Type: tokenizer.Percent,
		Pos: scanner.Position{
			Filename: "",
			Offset:   65,
			Line:     2,
			Column:   22,
		},
	}, tk.Next())

	assert.Equal(t, tokenizer.Token{
		Type: tokenizer.Text,
		Pos: scanner.Position{
			Filename: "",
			Offset:   66,
			Line:     2,
			Column:   23,
		},
		Value: " Imports amsmath\n",
	}, tk.Next())
}

func TestTokenizer2(t *testing.T){
	tk, err := tokenizer.Open("../../resources/test/tokenizer/1.tex")
	assert.Nil(t, err)

	for ;; {
		next := tk.Next()
		if next.Type == tokenizer.EOF {
			return
		}

		fmt.Printf("%+v\n", next)
	}
}