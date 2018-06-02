package idl

import (
	"fmt"
	"github.com/superp00t/etc"
	"io"
	"log"
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

var keyword = map[string]Token{
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
	t := new(TokenPos)
	t.T = keyword[l.currentString]
	t.S = l.currentString
	l.currentString = ""
	t.Ln = l.Ln
	t.Col = l.Col
	return t, nil
}

func (l *Lexer) ReadToken() (*TokenPos, error) {
	curString := []rune{}

	if l.eof {
		return nil, io.EOF
	}

	for {
		r, _, err := l.ReadRune()
		fmt.Println(err, "?", string(r))
		if err != nil {
			fmt.Println(err)
			log.Fatal(l.String())
		}
	}
	log.Fatal(l.String())

	for {
		rn, _, err := l.ReadRune()
		if err != nil {
			// Terminate keyword upon EOF.
			rn = ' '
			if err == io.EOF {
				l.eof = true
			}
		}

		if rn == '#' {
			l.IgnoreComments()
			continue
		}

		if rn == '\n' {
			l.Ln += 1
			l.Col = 0
			if string(curString) != "" {
				return l.terminate()
			}
			continue
		}

		if rn == '\r' {
			l.Ln += 1
			l.Col = 0
			continue
		}

		l.Col += 1

		if (rn == ' ' || rn == '\t') && string(curString) == "" && (l.Ln == 0 && l.Col == 0) {
			return l.terminate()
		}

		if rn == ' ' || rn == '\t' {
			return l.terminate()
		}

		fmt.Println(err, ">", string(rn))

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

		t = append(t, tk)
	}

	return t, nil
}
