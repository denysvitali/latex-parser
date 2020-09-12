package tokenizer

import (
	"os"
	"text/scanner"
)

type Tokenizer struct {
	scanner          *scanner.Scanner
	OriginalFilePath string

	tokenStack []Token
}

func getFirstRune(s string) rune {
	runes := []rune(s)
	return runes[0]
}

func (t *Tokenizer) Peek() Token {
	if len(t.tokenStack) == 0 {
		t.tokenStack = append(t.tokenStack, t.Next())
		return t.tokenStack[0]
	}

	return t.tokenStack[0]
}

func (t *Tokenizer) Next() Token {
	if len(t.tokenStack) != 0 {
		lastToken := t.tokenStack[0]
		t.tokenStack = t.tokenStack[1:]
		return lastToken
	}

	initialPos := t.scanner.Pos()
	r := t.scanner.Peek()
	rs := string(r)

	if rs == "\\" {
		t.scanner.Next()
		peekStr := string(t.scanner.Peek())

		if peekStr == "\\" {
			// New line
			t.scanner.Next()
			return Token {
				Type: NewLine,
				Pos: initialPos,
			}
		}

		if peekStr == "$" {
			t.scanner.Next()
			return Token{
				Type: Text,
				Value: "$",
				Pos: initialPos,
			}
		}

		if peekStr == "%" {
			t.scanner.Next()
			return Token {
				Type: Text,
				Value: "%",
				Pos: initialPos,
			}
		}

		return t.identifier(initialPos)
	}

	if rs == "{" {
		t.scanner.Next()
		return Token{
			Type: OpenCurly,
			Pos: initialPos,
		}
	}

	if rs == "}" {
		t.scanner.Next()
		return Token {
			Type: CloseCurly,
			Pos: initialPos,
		}
	}

	if rs == "[" {
		t.scanner.Next()
		return Token {
			Type: OpenSquare,
			Pos: initialPos,
		}
	}

	if rs == "]" {
		t.scanner.Next()
		return Token {
			Type: CloseSquare,
			Pos: initialPos,
		}
	}

	if rs == "%" {
		t.scanner.Next()
		return Token {
			Type: Percent,
			Pos: initialPos,
		}
	}

	if rs == "$" {
		t.scanner.Next()
		return Token {
			Type: Dollar,
			Pos: initialPos,
		}
	}

	if r == scanner.EOF {
		t.scanner.Next()
		return Token {
			Type: EOF,
			Pos: initialPos,
		}
	}

	return t.text(initialPos)
}

func isAlpha(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func isNum(r rune) bool {
	return r >= '0' && r <= 9
}

func isAlphaNum(r rune) bool {
	return isAlpha(r) || isNum(r)
}

func isLatexChar(r rune) bool {
	switch r {
	case '{', '}', '\\', '%', '[', ']', '$':
		return true
	}

	return false
}

func isValidTextChar(r rune) bool {
	return isAlphaNum(r) || !isLatexChar(r)
}

func (t Tokenizer) text(pos scanner.Position) Token {
	var tokenContent []rune
	for ;; {
		currentRune := t.scanner.Peek()
		if !isValidTextChar(currentRune) {
			break
		}

		if currentRune == -1 {
			break
		}
		tokenContent = append(tokenContent, t.scanner.Next())
	}

	if len(tokenContent) == 0 {
		return Token{
			Type: Unknown,
			Pos: pos,
		}
	}

	return Token{
		Type: Text,
		Pos: pos,
		Value: string(tokenContent),
	}
}

func (t Tokenizer) identifier(pos scanner.Position) Token {
	var id []rune
	var currentRune rune

	for ;; {
		currentRune = t.scanner.Peek()
		if !isValidIdentifierChar(currentRune) {
			break
		}

		currentRune = t.scanner.Next()
		id = append(id, currentRune)
	}

	return Token{
		Name: string(id),
		Type: Identifier,
		Pos: pos,
	}
}

func isValidIdentifierChar(r rune) bool {
	return isAlphaNum(r) || r == '*'
}

func Open(path string) (*Tokenizer, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	var s scanner.Scanner
	s.Init(file)

	return &Tokenizer{scanner: &s, OriginalFilePath: path}, nil
}
