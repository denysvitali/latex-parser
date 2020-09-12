package main

import (
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/denysvitali/latex-parser/pkg/parser"
	"github.com/denysvitali/latex-parser/pkg/parser/symbols"
	"github.com/denysvitali/latex-parser/pkg/tokenizer"
	"github.com/sirupsen/logrus"
)

type GetTextCmd struct {
	Input string `arg:"positional,required"`
}

var args struct {
	GetTextCmd *GetTextCmd `arg:"subcommand:text"`
}

func main(){
	arg.MustParse(&args)

	if args.GetTextCmd != nil {
		tkz, err := tokenizer.Open(args.GetTextCmd.Input)
		if err != nil {
			logrus.Fatal(err)
		}

		p := parser.New(tkz)
		symb, err := p.Parse()

		if err != nil {
			logrus.Fatal(err)
		}

		printText(symb)

	}
}

func printText(symb []symbols.Symbol) {
	for _, v := range symb {
		switch v.Type() {
		case symbols.Text:
			text := v.(symbols.TextSymbol)
			fmt.Print(text.Content)
		case symbols.Env:
			env := v.(symbols.EnvSymbol)
			printText(env.Statements)
		case symbols.NewLine:
			fmt.Print("\n")
		}
	}
}
