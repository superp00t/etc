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
	Mbigint
	Mint
	Muint
	Mbool
	Mdate
	McheckLater
	Menum
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
	case "bool":
		return mkspec(Mbool)
	case "string":
		return mkspec(Mstring)
	case "uint8", "byte":
		return mkspec(Muint8)
	case "uint16":
		return mkspec(Muint16)
	case "uint32":
		return mkspec(Muint32)
	case "uint64":
		return mkspec(Muint64)
	case "bytes":
		return mkspec(Mbytes)
	case "float32":
		return mkspec(Mfloat32)
	case "uuid":
		return mkspec(Muuid)
	case "int":
		return mkspec(Mint)
	case "time", "uint":
		return mkspec(Muint)
	case "bigint":
		return mkspec(Mbigint)
	case "float64":
		return mkspec(Mfloat64)
	case "date":
		return mkspec(Mdate)
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

type RPC_Service struct {
	Procs map[string]*RPC_Proc
}

type RPC_Proc struct {
	Name         string
	RequestType  string
	ResponseType string
}

type Syntax struct {
	PackageName string

	Enums map[string]*enumMap

	Pragmas      map[string]bool
	MessageTypes map[string]*SpecType

	RPC map[string]*RPC_Service
}

type enumMap struct {
	keys []string
	vals []uint64
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
	s.RPC = make(map[string]*RPC_Service)
	s.MessageTypes = make(map[string]*SpecType)
	s.Enums = make(map[string]*enumMap)

	for i := 0; i < len(tokens); {
		tk := tokens[i]
		if tk.T == TEnum {
			i++
			tk2 := tokens[i]
			if tk2.T != TName {
				return nil, fmt.Errorf("(%d, %d) unexpected token %s", tk2.Ln, tk2.Col, tk2.S)
			}

			i++
			tk3 := tokens[i]
			if tk3.T != TOpenBracket {
				return nil, fmt.Errorf("(%d, %d) expected open bracket, got \"%s\"", tk3.Ln, tk3.Col, tk3.S)
			}

			en := new(enumMap)
			var o uint64

			i++
			for y := i; y < len(tokens); y++ {
				if tokens[y].T == TCloseBracket {
					i = y + 1
					break
				}

				en.keys = append(en.keys, tokens[y].S)
				en.vals = append(en.vals, o)
				o += 1
			}

			s.Enums[tk2.S] = en
			tk = tokens[i]
		}

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

		if tk.T == TRPC {
			rpc := new(RPC_Service)
			rpc.Procs = make(map[string]*RPC_Proc)
			rnm := ""

			if len(tokens) < i+2 {
				return nil, fmt.Errorf("(%d, %d) Error parsing struct, unexpected EOF", tk.Ln, tk.Col)
			}

			rnm = tokens[i+1].S
			s.RPC[rnm] = rpc

			if tokens[i+2].T != TOpenBracket {
				return nil, fmt.Errorf("(%d, %d) Error parsing rpc, expected open bracket", tk.Ln, tk.Col)
			}

			i += 3

			for {
				tk = tokens[i]
				if tk.T == TCloseBracket {
					i++
					break
				}

				if tk.T == TName {
					if !strings.Contains(tk.S, "(") && !strings.Contains(tk.S, ")") {
						return nil, fmt.Errorf("(%d, %d) Error parsing rpc, needs function", tk.Ln, tk.Col)
					}

					nm := []rune(tk.S)
					rrn := nm[len(nm)-1:][0]
					if rrn != ')' {
						return nil, fmt.Errorf("(%d, %d) Error parsing rpc, needs function parentheses (%c)", tk.Ln, tk.Col, rrn)
					}

					funcNameS := strings.Split(tk.S, "(")
					funcName := funcNameS[0]

					centerTypeS := strings.Split(funcNameS[1], ")")
					centerType := centerTypeS[0]
					if centerType == "" {
						centerType = "void"
					}

					i++
					tk = tokens[i]
					if tk.T != TReturns {
						return nil, fmt.Errorf("(%d, %d) Error parsing rpc, needs function return ->", tk.Ln, tk.Col)
					}
					i++
					tk = tokens[i]
					if tk.T != TName {
						return nil, fmt.Errorf("(%d, %d) Error parsing rpc, needs function return ->", tk.Ln, tk.Col)
					}

					returnType := tk.S

					rpc.Procs[funcName] = &RPC_Proc{
						funcName,
						centerType,
						returnType,
					}
					i++
				}
			}
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
				var typeSpec *SpecType
				if s.Enums[tz.S] != nil {
					typeSpec = mkspec(Menum)
				} else {
					typeSpec = s.CheckIfTypeExists(tz.S)
				}

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
