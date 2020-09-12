package tokenizer

import "text/scanner"

type TokenType string

const (
	Identifier TokenType = "Identifier"
	OpenCurly TokenType = "OpenCurly"
	CloseCurly TokenType = "CloseCurly"
	OpenSquare TokenType = "OpenSquare"
	CloseSquare TokenType = "CloseSquare"
	Dollar TokenType = "Dollar"
	EOF TokenType = "EOF"
	NewLine TokenType = "NewLine"
	Percent TokenType = "Percent"
	Text TokenType = "Text"
	Unknown TokenType = "Unknown"
)

type Token struct {
	Type  TokenType
	Name  string
	Value interface{}
	Pos   scanner.Position
}


type TextToken struct {
	Token
	Value string
}