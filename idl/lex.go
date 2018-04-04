package idl

import (
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
)

var KwMap = map[string]Token{
	"->":     TReturns,
	"struct": TStruct,
	"use":    TPragma,
	"rpc":    TRPC,
	"{":      TOpenBracket,
	"}":      TCloseBracket,
}

type TokenPos struct {
	T       Token
	Ln, Col int

	S string
}

type Lexer struct {
	*etc.Buffer

	Ln, Col int
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		etc.FromString(input),
		0,
		0,
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

func (l *Lexer) ReadToken() (*TokenPos, error) {
	curString := []rune{}

	for {
		rn, _, err := l.ReadRune()
		if err != nil {
			// Terminate keyword upon EOF.
			rn = ' '
		}

		if rn == '#' {
			l.IgnoreComments()
			continue
		}

		if rn == '\n' {
			l.Ln += 1
			l.Col = 0
			if string(curString) != "" {
				t := new(TokenPos)
				s := string(curString)
				// This is okay, because the map's zero value is TName, or referring to as of yet undefined keywords.
				t.T = KwMap[s]
				t.S = s
				t.Ln = l.Ln
				t.Col = l.Col
				return t, nil
			}
			continue
		}

		if rn == '\r' {
			l.Ln += 1
			l.Col = 0
			continue
		}

		l.Col += 1

		if (rn == ' ' || rn == '\t') && string(curString) == "" {
			continue
		}

		if rn == ' ' || rn == '\t' {
			t := new(TokenPos)
			s := string(curString)
			t.T = KwMap[s]
			t.S = s
			t.Ln = l.Ln
			t.Col = l.Col
			return t, nil
		}

		curString = append(curString, rn)
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

		t = append(t, tk)
	}

	return t, nil
}
