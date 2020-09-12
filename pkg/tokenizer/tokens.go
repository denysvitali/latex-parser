package tokenizer

type TokenType string

const (
	Identifier TokenType = "Identifier"
	OpenCurly TokenType = "OpenCurly"
	CloseCurly TokenType = "CloseCurly"
	EOF TokenType = "EOF"
	NewLine TokenType = "NewLine"
	Percent TokenType = "Percent"
	Text TokenType = "Text"
	Unknown TokenType = "Unknown"
)
