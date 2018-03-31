# Etc: Efficient Transfer Coding

Etc is an encoding system. You can write manual encoding functions with it, or you can use my format language called EtcSchema.

EtcSchema draws heavy inspiration from Google's Protocol Buffers. It uses an offset-based system, in lieu of syntax: you just read one type, and move on to the next. The tokenization of data relies on types remaining a fixed size, or having some kind of termination header: for example, with C strings, you write a string and append a NULL character to the end. With an LEB128 integer, it is terminated by the eigth bit.

Etc is definitely unstable and **NOT** production ready, so proceed with caution.

etc.Buffer is similar to `bytes.Buffer`, but with number and string methods so you don't have to manually type in `encoding/binary` functions.

it is intended to be fault-tolerant, and as such does not return errors for most reading and writing functions.

## Encoding with EtcSchema

Create a EtcSchema file like this:

`exampleData.etcschema`
```c
// Use compression (optional)
// #pragma zlib-compress

struct exampleData {
  uint64      time_ms
  uuid        id
  float32     coordinates[]
}
```

```
$ etc-schema-gen exampleData.etcschema --go_out=./pkg/ --pkg=pkg
```

Will create pkg/pkg.etc.go:

```go
package pkg

import (
	"encoding/json"

	"github.com/superp00t/etc"
)

type ExampleData struct {
	Time_ms     uint64    `json:"time_ms"`
	Id          etc.UUID  `json:"id"`
	Coordinates []float32 `json:"coordinates"`
}

func UnmarshalExampleData(data []byte) (*ExampleData, error) {
	var err error
	input := data
	v := new(ExampleData)
	d := etc.MkBuffer(input)
	v.Time_ms = d.ReadUint64()
	v.Id = d.ReadUUID()
	ln_Coordinates := int(d.Read_LEB128_Uint())
	for _i := 0; _i < ln_Coordinates; _i++ {
		v.Coordinates = append(v.Coordinates, d.ReadFloat32())
	}
	return v, err
}

func (v *ExampleData) Marshal() []byte {
	d := etc.NewBuffer()
	d.WriteUint64(v.Time_ms)
	d.WriteUUID(v.Id)
	d.Write_LEB128_Uint(uint64(len(v.Coordinates)))
	for _i := 0; _i < len(v.Coordinates); _i++ {
		e := v.Coordinates[_i]
		d.WriteFloat32(e)
	}
	return d.Bytes()
}

func (v *ExampleData) Stringify() string {
	b, _ := json.Marshal(v)
	return string(b)
}

```

## Things to be aware of:

- The types `int` and `uint` are aliases for 64-bit integer types, but using variable length encoding (LEB128). For fixed bit length encoding, use types with the length in the name, e.g. `uint16`.

## Coming soon™

- RPC interface definition
- imports
- speed increases
- Code generators for JavaScript, Lua, and C
