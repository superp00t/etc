package idl

import (
	"fmt"
	"sort"

	"github.com/ditashi/jsbeautifier-go/jsbeautifier"
)

func defaultJS(i *SpecType) string {
	if i.ArrayType {
		return "[]"
	}

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
		return `""`
	case Mdate:
		return "new Date(0)"
	case Mstruct:
		return fmt.Sprintf("new %s()", i.StructName)
	case Menum:
		return "0"
	default:
		return "null"
	}
}

func readJS(i *SpecType) string {
	if !i.ArrayType {
		return "this." + i.FieldName + " = " + readJSunderlying(i)
	}

	rjs := readJSunderlying(i)
	src := fmt.Sprintf("var %s_len = input.readUnsignedVarint().toNumber();\n", i.FieldName)
	src += fmt.Sprintf("  this.%s = new Array(%s_len);\n\n", i.FieldName, i.FieldName)
	src += fmt.Sprintf("  for (var _i = 0; _i < %s_len; _i++) {\n", i.FieldName)
	src += fmt.Sprintf("    this.%s[_i] = %s;\n", i.FieldName, rjs)
	src += "  }\n\n"

	return src
}

func readJSunderlying(i *SpecType) string {
	switch i.Type {
	case Mbytes:
		return "input.readLimitedBytes()"
	case Mbigint:
		return "input.readSignedVarint()"
	case Mbool:
		return "input.readBoolean()"
	case Mfloat32:
		return "input.readFloat32()"
	case Mfloat64:
		return "input.readFloat64()"
	case Mint:
		return "input.readSignedVarint()"
	case Muint:
		return "input.readUnsignedVarint()"
	case Muint64:
		return "input.readUint64()"
	case Muint32:
		return "input.readUint32()"
	case Muint16:
		return "input.readUint16()"
	case Muint8:
		return "input.readByte()"
	case Mstring:
		return "input.readString()"
	case Mdate:
		return "input.readDate()"
	case Menum:
		return "input.readUint()"
	case Mstruct:
		return fmt.Sprintf("new %s().decodeBuf(input)", i.StructName)
	default:
		return "null"
	}
}

func writeJS(i *SpecType) string {
	if !i.ArrayType {
		return "_out." + fmt.Sprintf(writeJSunderlying(i), "this."+i.FieldName)
	}

	wt := writeJSunderlying(i)
	src := "\n_out.writeUnsignedVarint(new etc.Bn(this." + i.FieldName + ".length));\n"
	src += "  for (var _i = 0; _i < this." + i.FieldName + ".length; _i++) {\n"
	src += "    _out." + fmt.Sprintf(wt, "this."+i.FieldName+"[_i]") + "\n"
	src += "  }\n"

	return src
}

func writeJSunderlying(i *SpecType) string {
	switch i.Type {
	case Mbytes:
		return "writeLimitedBytes(%s)"
	case Mbigint:
		return "writeSignedVarint(%s)"
	case Mbool:
		return "writeBoolean(%s)"
	case Mfloat32:
		return "writeFloat32(%s)"
	case Mfloat64:
		return "writeFloat64(%s)"
	case Muint:
		return "writeUnsignedVarint(%s)"
	case Muint64:
		return "writeUint64(%s)"
	case Muint32:
		return "writeUint32(%s)"
	case Muint16:
		return "writeUint16(%s)"
	case Muint8:
		return "writeByte(%s)"
	case Mstring:
		return "writeString(%s)"
	case Mstruct:
		return "writeBytes(%s.encode())"
	case Mdate:
		return "writeDate(%s)"
	case Menum:
		return "writeUint(%s)"
	default:
		return "writeByte(0)"
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

func (s *Syntax) SortedEnumTypes() []string {
	var st []string
	for k := range s.Enums {
		st = append(st, k)
	}

	sort.Strings(st)
	return st
}

func spaceString(in, sub *SpecType) string {
	max := 0

	for _, v := range in.StructTypes {
		n := len([]rune(v.FieldName))

		if n > max {
			max = n
		}
	}

	max -= len([]rune(sub.FieldName))

	out := " "
	for i := 0; i < max; i++ {
		out += " "
	}

	return out
}

func (s *Syntax) GenerateJS() (string, string) {
	src := `const etc = require("etc-js");`
	src += "\n\n"

	for _, ens := range s.SortedEnumTypes() {
		en := s.Enums[ens]
		src += "var Enum_" + ens + " = {\n"
		for idx, v := range en.keys {
			src += "  " + v + ": " + fmt.Sprintf("%d,\n", en.vals[idx])
		}
		src += "}\n"
	}

	for _, t := range s.SortedMsgTypes() {
		src += "\n"
		src += `/**`
		src += "\n"
		src += `* @class`
		src += "\n"
		src += `*/`
		src += "\n"
		src += fmt.Sprintf("function %s() {\n", t)
		mt := s.MessageTypes[t]
		for _, v := range mt.StructTypes {
			sp := spaceString(mt, v)
			src += fmt.Sprintf("  this.%s%s= %s;\n", v.FieldName, sp, defaultJS(v))
		}

		src += "}\n\n"

		src += "/**\n* @return {etc.Buffer}\n*/"
		src += fmt.Sprintf("\n%s.prototype.encodeBuf = function(_out) {", mt.StructName)
		for _, v := range mt.StructTypes {
			src += "\n  " + writeJS(v)
		}
		src += "\n  return _out;\n}\n\n"
		src += fmt.Sprintf("/**\n* @param  {etc.Buffer} in\n* @return {%s}\n**/", mt.StructName)
		src += fmt.Sprintf("\n%s.prototype.decodeBuf = function(input) {", t)
		for _, v := range mt.StructTypes {
			src += fmt.Sprintf("\n  %s", readJS(v))
		}
		src += "\n  return this;\n}\n"

		src += fmt.Sprintf("\n/**\n* Serializes %s to a Uint8Array.\n* @return {Uint8Array}\n **/\n%s.prototype.encode = function() {\n", mt.StructName, mt.StructName)
		src += "  var _out = new etc.Buffer();\n  this.encodeBuf(_out);\n  return _out.finish();\n}\n"

		src += fmt.Sprintf("\n/**\n* Deserializes %s from a Uint8Array.\n* @param {Uint8Array} in\n**/\n%s.prototype.decode = function(input) {\n", mt.StructName, mt.StructName)
		src += "  return this.decodeBuf(new etc.Buffer(input));\n}\n"
	}

	src += "\nmodule.exports = {"
	for _, v := range s.SortedMsgTypes() {
		spc := spaceExport(s, v)
		src += "\n  " + v + ": " + spc + v + ","
	}

	for _, ens := range s.SortedEnumTypes() {
		src += "\n  Enum_" + ens + ": " + "Enum_" + ens + ","
	}

	src += "\n};\n"

	return beautifyJS(src), ""
}

func spaceExport(s *Syntax, id string) string {
	max := 0
	for k := range s.MessageTypes {
		n := len([]rune(k))
		if n > max {
			max += n
		}
	}

	max -= len([]rune(id))

	var out string
	for i := 0; i < max; i++ {
		out += " "
	}

	return out
}

func beautifyJS(input string) string {
	options := jsbeautifier.DefaultOptions()
	options["indent_size"] = 2
	code, err := jsbeautifier.Beautify(&input, options)
	if err != nil {
		panic(err)
	}
	return code
}
