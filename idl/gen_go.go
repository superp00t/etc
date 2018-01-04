package idl

import (
	"fmt"
	"strings"
)

func (s *Syntax) GenerateGo() string {
	src := fmt.Sprintf("package %s\n\n", s.PackageName)
	encS := ""
	encodeS := ""

	for k, v := range s.MessageTypes {
		sname := goNormalizeExportField(k)
		// generate decoder
		encodeS += "func (v *" + sname + ") Marshal() []byte {\n\td := etc.NewBuffer()\n"
		encS += "func Unmarshal" + sname + "(data []byte) (*" + sname + ", error) {\n"
		if s.Pragmas["zlib-compress"] == true {
			encS += "\tvar err error\n\tinput, err := etc.ZlibDecompress(data)\n\tif err != nil {\n\t\treturn nil, err\n\t}\n\t"
		} else {
			encS += "\tinput := data\n"
		}

		encS += "\tv := new(" + sname + ")\n\td := etc.MkBuffer(input)\n"

		src += fmt.Sprintf("type %s struct {\n", sname)
		for _, field := range v.StructTypes {
			fname := goNormalizeExportField(field.FieldName)
			ap := ""
			if field.ArrayType {
				ap = "[]"
			}
			src += fmt.Sprintf("\t%s\t%s%s", fname, ap, goType(field))
			src += "\t" + "`" + "json:\"" + field.FieldName + "\"`\n"

			encName := "v." + fname
			if field.ArrayType {
				encodeS += "\td.WriteUint32(uint32(len(v." + fname + ")))\n"
				encodeS += "\tfor _i := 0; _i < len(v." + fname + "); _i++ {\n\t\te := v." + fname + "[_i]\n\t"
				encName = "e"

				l := "ln_" + fname
				encS += "\t" + l + " := int(d.ReadUint32())\n\tfor _i := 0; _i < " + l + "; _i++ {\n\t"
			}

			if field.Type == Mstruct {
				encodeS += "\td.WriteLimitedBytes(" + encName + ".Marshal())\n"
				encS += "\tdcc, err := etc.ZlibDecompress(d.ReadLimitedBytes())\n\tif err != nil {\n\t\treturn nil, err\n\t}\n\t"
				encS += "\tchnk, err := Unmarshal" + field.StructName + "(dcc)\n\t\tif err != nil {\n\t\t\treturn nil, err\n\t\t}\n"
				if field.ArrayType {
					encS += "\t\tv." + fname + " = append(v." + fname + ", chnk)\n"
				} else {
					encS += "\tv." + fname + " = chnk\n"
				}
			} else {
				argName := "v." + fname
				if field.ArrayType {
					argName = "e"
					encS += "\tv." + fname + " = append(v." + fname + ", " + goReadFunc(field) + ")\n"
				} else {
					encS += "\tv." + fname + " = " + goReadFunc(field) + "\n"
				}
				encodeS += "\t" + goWriteFunc(field, argName) + "\n"
			}

			if field.ArrayType {
				encodeS += "\t}\n"
				encS += "\t}\n"
			}
		}

		returnArg := "return d.Bytes()"
		if s.Pragmas["zlib-compress"] == true {
			returnArg = "return etc.ZlibCompress(d.Bytes())"
		}

		encodeS += "\t" + returnArg + "\n}\n\n"
		encS += "\treturn v, err\n}\n\n"
		src += fmt.Sprintf("}\n\n")

	}

	return src + encS + encodeS
}

func goType(m *SpecType) string {
	switch m.Type {
	case Mbytes:
		return "[]byte"
	case Mstring:
		return "string"
	case Muint16:
		return "uint16"
	case Muint32:
		return "uint32"
	case Muint64:
		return "uint64"
	case Muint8:
		return "uint8"
	case Mfloat32:
		return "float32"
	case Mfloat64:
		return "float64"
	case Mstruct:
		return "*" + m.StructName
	default:
		return "interface{}"
	}
}

func goWriteFunc(m *SpecType, fname string) string {
	switch m.Type {
	case Mbytes:
		return fmt.Sprintf("d.WriteLimitedBytes(%s)", fname)
	case Mstring:
		return fmt.Sprintf("d.WriteCString(%s)", fname)
	case Muint16:
		return fmt.Sprintf("d.WriteUint16(%s)", fname)
	case Muint32:
		return fmt.Sprintf("d.WriteUint32(%s)", fname)
	case Muint64:
		return fmt.Sprintf("d.WriteUint64(%s)", fname)
	case Muint8:
		return fmt.Sprintf("d.WriteByte(%s)", fname)
	case Mfloat32:
		return fmt.Sprintf("d.WriteFloat32(%s)", fname)
	case Mfloat64:
		return fmt.Sprintf("d.WriteFloat64(%s)", fname)
	case Mstruct:
		return "/* NOT YET IMPLEMENTED */"
	default:
		return "interface{}"
	}
}

func goReadFunc(m *SpecType) string {
	switch m.Type {
	case Mbytes:
		return "d.ReadLimitedBytes()"
	case Mstring:
		return "d.ReadCString()"
	case Muint16:
		return "d.ReadUint16()"
	case Muint32:
		return "d.ReadUint32()"
	case Muint64:
		return "d.ReadUint64()"
	case Muint8:
		return "d.ReadByte()"
	case Mfloat32:
		return "d.ReadFloat32()"
	case Mfloat64:
		return "d.ReadFloat64()"
	case Mstruct:
		return "/* NOT YET IMPLEMENTED */"
	default:
		return "interface{}"
	}
}

func goNormalizeExportField(input string) string {
	in := []rune(input)
	tail := input[1:]
	firstChar := string(in[0])
	firstChar = strings.ToUpper(firstChar)

	return firstChar + string(tail)
}
