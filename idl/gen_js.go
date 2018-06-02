package idl

import (
	"fmt"
	"sort"
)

func defaultJS(i *SpecType) string {
	switch i.Type {
	case Mbytes:
		return "[]"
	case Mbigint:
		return "new etc.Bn(0)"
	case Mbool:
		return "false"
	case Mfloat32, Mfloat64:
		return "0.0"
	case Muint, Mint, Muint64, Muint32, Muint16:
		return "new etc.Bn(0)"
	case Mstring:
		return ""
	case Mstruct:
		return "{}"
	default:
		return "null"
	}
}

func (s *Syntax) SortedMsgTypes() []string {
	var st []string
	for k := range s.MessageTypes {
		st = append(st, k)
	}

	sort.Strings(st)
	return st
}

func (s *Syntax) GenerateJS() (string, string) {
	src := `const etc = require("etc-js");`
	src += "\n\n"

	for _, t := range s.SortedMsgTypes() {
		src += "\n"
		src += fmt.Sprintf("function %s() {\n", t)
		mt := s.MessageTypes[t]
		for _, v := range mt.StructTypes {
			src += fmt.Sprintf("\tthis.%s = %s;\n", v.FieldName, defaultJS(v))
		}

		src += "}\n"
	}

	return src, ""
}
