# Etc (WIP)

Etc is an encoding system. You can write manual encoding functions with it, or you can use my format language called EtcSchema.

EtcSchema draws heavy inspiration from Google's Protocol Buffers.

Etc is definitely unstable and **NOT** production ready, so proceed with caution.

etc.Buffer is similar to `bytes.Buffer`, but with number and string methods so you don't have to manually type in `encoding/binary` functions.

it is intended to be fault-tolerant, and as such does not return errors for most reading and writing functions.

## Encoding with EtcSchema

Create a EtcSchema file like this:

`exampleData.etcschema`
```go
// Use compression (optional)
#pragma zlib-compress

struct exampleData {
  uint64  time_ms

  float32 coordinates[]
}
```

```
$ etc-schema-gen exampleData.etcschema --go_out=./pkg/ --pkg=pkg
```

Will create pkg/pkg.etc.go:

```go
package pkg

type ExampleData struct {
	Time_ms     uint64    `json:"time_ms"`
	Coordinates []float32 `json:"coordinates"`
}

func UnmarshalExampleData(data []byte) (*ExampleData, error) {
	var err error
	input, err := etc.ZlibDecompress(data)
	if err != nil {
		return nil, err
	}
	v := new(ExampleData)
	d := etc.MkBuffer(input)
	v.Time_ms = d.ReadUint64()
	ln_Coordinates := int(d.ReadUint32())
	for _i := 0; _i < ln_Coordinates; _i++ {
		v.Coordinates = append(v.Coordinates, d.ReadFloat32())
	}
	return v, err
}

func (v *ExampleData) Marshal() []byte {
	d := etc.NewBuffer()
	d.WriteUint64(v.Time_ms)
	d.WriteUint32(uint32(len(v.Coordinates)))
	for _i := 0; _i < len(v.Coordinates); _i++ {
		e := v.Coordinates[_i]
		d.WriteFloat32(e)
	}
	return etc.ZlibCompress(d.Bytes())
}

```

## Coming soon™

- RPC interface definition
- imports
- speed increases
- Code generators for JavaScript, Lua, and C
