package idl

import (
	"fmt"
	"regexp"
	"strings"
)

type MType int

const (
	Mstring MType = iota
	Muint8
	Muint16
	Muint32
	Muint64
	Mbytes
	Mfloat32
	Mfloat64
	Mstruct
	Muuid
	Mint
	Muint
	McheckLater
)

func mkspec(m MType) *SpecType {
	return &SpecType{
		Type: m,
	}
}

func (s *Syntax) CheckIfTypeExists(str string) *SpecType {
	if mt := s.MessageTypes[str]; mt != nil {
		return mt
	}

	switch str {
	case "string":
		return mkspec(Mstring)
	case "uint8", "byte":
		return mkspec(Muint8)
	case "uint16":
		return mkspec(Muint16)
	case "uint32":
		return mkspec(Muint32)
	case "time", "uint64":
		return mkspec(Muint64)
	case "bytes":
		return mkspec(Mbytes)
	case "float32":
		return mkspec(Mfloat32)
	case "uuid":
		return mkspec(Muuid)
	case "int":
		return mkspec(Mint)
	case "uint":
		return mkspec(Muint)
	case "float64":
		return mkspec(Mfloat64)
	default:
		c := mkspec(McheckLater)
		c.StructName = str
		return c
	}
}

type SpecType struct {
	Type        MType
	StructTypes []*SpecType
	StructName  string
	FieldName   string
	ArrayType   bool

	Ln, Col int
}

type Syntax struct {
	PackageName string

	Pragmas      map[string]bool
	MessageTypes map[string]*SpecType
}

func (s *Syntax) FixUnknownTypes() error {
	for _, v := range s.MessageTypes {
		for _, st := range v.StructTypes {
			if st.Type == McheckLater {
				if n := s.MessageTypes[st.StructName]; n != nil {
					st.Type = Mstruct
				} else {
					return fmt.Errorf("(%d, %d) Undefined struct type %s", st.Ln, st.Col, st.StructName)
				}
			}
		}
	}

	return nil
}

func ParseTokens(tokens []*TokenPos) (*Syntax, error) {
	s := &Syntax{}
	s.Pragmas = make(map[string]bool)
	s.MessageTypes = make(map[string]*SpecType)

	for i := 0; i < len(tokens); {
		tk := tokens[i]
		if tk.T == TPragma {
			if len(tokens) < i+2 {
				return nil, fmt.Errorf("(%d, %d) Error parsing Pragma: unexpected EOF", tk.Ln, tk.Col)
			}
			ta := tokens[i+1]
			if ta.T != TName {
				return nil, fmt.Errorf("(%d, %d) Error parsing Pragma: pragma is reserved keyword", tk.Ln, tk.Col)
			}

			s.Pragmas[ta.S] = true
			i += 2
			continue
		}

		if tk.T == TStruct {
			if len(tokens) < i+2 {
				return nil, fmt.Errorf("(%d, %d) Error parsing struct, unexpected EOF", tk.Ln, tk.Col)
			}
			ta := tokens[i+1]

			name := ta.S
			if !isValidName(name) {
				return nil, fmt.Errorf("(%d, %d) Invalid name \"%s\"", ta.Ln, ta.Col, name)
			}

			i += 3

			spcc := new(SpecType)
			spcc.Type = Mstruct
			spcc.StructName = name

			for {
				tz := tokens[i]
				if tz.T == TCloseBracket {
					i++
					break
				}
				typeSpec := s.CheckIfTypeExists(tz.S)
				typeSpec.FieldName = tokens[i+1].S
				if strings.HasSuffix(typeSpec.FieldName, "[]") {
					typeSpec.ArrayType = true
					typeSpec.FieldName = strings.Replace(typeSpec.FieldName, "[]", "", -1)
				}
				typeSpec.Ln = tz.Ln
				typeSpec.Col = tz.Col
				spcc.StructTypes = append(spcc.StructTypes, typeSpec)
				i += 2
				if i > (len(tokens) - 1) {
					break
				}
			}

			s.MessageTypes[name] = spcc

			continue
		}

		panic(fmt.Errorf("Unknwon type %s %s", tk.T, tk.S))
	}

	return s, s.FixUnknownTypes()
}

func isValidName(n string) bool {
	b, err := regexp.MatchString("^[a-zA-Z0-9_]*$", n)
	if err != nil {
		panic(err)
	}

	return b
}

func Parse(src string) (*Syntax, error) {
	t, err := Lex(src)
	if err != nil {
		return nil, err
	}

	s, err := ParseTokens(t)
	if err != nil {
		return nil, err
	}

	return s, nil
}
