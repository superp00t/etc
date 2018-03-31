package idl

type Token int

//go:generate stringer -type=Token
const (
	TName Token = iota
	TNewline
	TPragma
	TStruct
	TOpenBracket
	TCloseBracket
	TServer
)

var KwMap = map[string]Token{
	"struct":  TStruct,
	"#pragma": TPragma,
	"server":  TServer,
	"{":       TOpenBracket,
	"}":       TCloseBracket,
}

type TokenPos struct {
	T       Token
	Ln, Col int

	S string
}

func Lex(input string) ([]*TokenPos, error) {
	var t []*TokenPos
	var cur *TokenPos
	ln, col := 0, 0
	curString := ""

scan:
	for i := 0; i < len(input); i++ {
		v := input[i]
		// ignore comments
		if v == '/' && input[i+1] == '/' {
			i += 1
			for x := i; ; x++ {
				col++
				if input[x] == '\n' {
					col = 0
					i = x
					ln++
					continue scan
				}
			}
		}

		if v == ' ' || v == '\t' {
			col++
			if curString != "" {
				cur.T = KwMap[curString]
				cur.S = curString
				curString = ""
				cur.Ln = ln
				cur.Col = col
				t = append(t, cur)
				cur = new(TokenPos)
			}
			continue
		}

		if v == '\n' {
			if curString != "" {
				cur.T = KwMap[curString]
				cur.S = curString
				curString = ""
				cur.Ln = ln
				cur.Col = col
				t = append(t, cur)
				cur = new(TokenPos)
			}
			ln++
			col = 0
			continue
		}

		curString += string(v)

		if cur == nil {
			cur = new(TokenPos)
			cur.Col = col
			cur.Ln = ln
		}

		col++
	}

	return t, nil
}
