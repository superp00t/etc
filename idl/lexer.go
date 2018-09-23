package idl

import (
	"io"

	"github.com/superp00t/etc"
)

type Token int

//go:generate stringer -type=Token
const (
	TName Token = iota
	TNewline
	TPragma
	TStruct
	TOpenBracket
	TCloseBracket
	TRPC
	TReturns
	TEnum
)

var keyword = map[string]Token{
	"->":     TReturns,
	"struct": TStruct,
	"use":    TPragma,
	"rpc":    TRPC,
	"{":      TOpenBracket,
	"}":      TCloseBracket,
	"enum":   TEnum,
}

type TokenPos struct {
	T       Token
	Ln, Col int

	S string
}

type Lexer struct {
	*etc.Buffer
	Ln, Col int

	currentString string

	eof bool
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		etc.FromString(input),
		0,
		0,
		"",
		false,
	}
}

func (l *Lexer) IgnoreComments() {
	for {
		rn, _, err := l.ReadRune()
		if err != nil {
			return
		}

		if rn == '\n' {
			l.Col = 0
			l.Ln++
			return
		}

		l.Col++
	}
}

func (l *Lexer) terminate() (*TokenPos, error) {
	if l.currentString == "" {
		return nil, nil
	}
	t := new(TokenPos)
	t.T = keyword[l.currentString]
	t.S = l.currentString
	l.currentString = ""
	t.Ln = l.Ln
	t.Col = l.Col
	return t, nil
}

func (l *Lexer) ReadToken() (*TokenPos, error) {
	for {
		rn, _, err := l.ReadRune()
		if err != nil {
			// Terminate keyword upon EOF.
			if err == io.EOF {
				return l.terminate()
			}
		}

		if rn == '#' {
			l.IgnoreComments()
			continue
		}

		if rn == ' ' || rn == '\t' {
			l.Col += 1
			return l.terminate()
		}

		if rn == '\n' {
			l.Col = 0
			l.Ln += 1
			return l.terminate()
		}

		l.currentString += string(rn)
	}
}

func Lex(input string) ([]*TokenPos, error) {
	var t []*TokenPos

	b := NewLexer(input)
	for b.Available() > 0 {
		tk, err := b.ReadToken()
		if err != nil {
			return nil, err
		}

		if tk != nil {
			t = append(t, tk)
		}
	}

	return t, nil
}
