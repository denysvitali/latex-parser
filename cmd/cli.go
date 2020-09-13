package main

import (
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/denysvitali/latex-parser/pkg/parser"
	"github.com/denysvitali/latex-parser/pkg/parser/symbols"
	"github.com/denysvitali/latex-parser/pkg/tokenizer"
	"github.com/sirupsen/logrus"
	"regexp"
	"strings"
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
			fmt.Print(strings.Replace(text.Content, "\n", " ", -1))
		case symbols.Env:
			env := v.(symbols.EnvSymbol)
			fmt.Printf("Env: %v\n", env.Environment)
			printText(env.Statements)
		case symbols.NewLine:
			fmt.Print("\n")
		case symbols.Macro:
			macro := v.(symbols.MacroSymbol)
			if isTextModifier(macro.MacroName) {
				printText(macro.CurlyArgs[0])
				continue
			}
			if isHeading(macro.MacroName){
				if len(macro.CurlyArgs) != 1 {
					continue
				}
				if len(macro.CurlyArgs[0]) != 1 {
					continue
				}
				if macro.CurlyArgs[0][0].Type() != symbols.Text {
					continue
				}
				content := macro.CurlyArgs[0][0].(symbols.TextSymbol)
				fmt.Printf("\n\n%s\n\n", content.Content)
				continue
			}

			fmt.Printf("macro not handled: %v\n", macro.MacroName)

		case symbols.Include:
			include := v.(symbols.IncludeSymbol)
			printText(include.Statements)
		case symbols.CurlyEnv:
			curlyEnv := v.(symbols.CurlyEnvSymbol)
			printText(curlyEnv.Statements)
		case symbols.InlineMath:
			// Discard
			fmt.Printf("Inline Math: %v", v)
		default:
			panic(fmt.Errorf("unhandled %v", v.Type()))
		}
	}
}

func isTextModifier(name string) bool {
	switch name {
	case "textit", "textbf", "emph", "underline":
		return true
	}

	return false
}

func isHeading(name string) bool {
	heading := regexp.MustCompile("^(?:(?:sub){,2}section|chapter|paragraph)(?:\\*|)$")
	return heading.MatchString(name)
}
