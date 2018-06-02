# Etc: Efficient Transfer Coding

Etc is a binary encoding system. You can use its encoding functions manually, or you can use the format description language I call EtcSchema.

[Etc in Rust](https://github.com/superp00t/etc-rs)

[Etc in JavaScript](https://github.com/superp00t/etc-js)

## Example API usage

```go
// Allocate empty buffer
e := etc.NewBuffer()

// Allocate buffer from string
e := etc.FromString("test")

// Allocate buffer from Base64
e := etc.FromBase64("dGVzdA==")

// Allocate buffer from bytes
e := etc.FromBytes([]byte{'t', 'e', 's', 't'})

// Create buffer as an alias to a file, preserving RAM but costing speed at runtime with disk IO
// Quite useful for parsing large files
e, err := etc.FileController("/tmp/newFile")
if err != nil {
  panic(err)
} else {
  e.Write([]byte("test"))
}

// Load a string from a defined 0-4 sector in buffer, ignoring possible zero padding bytes. This is not recommended and only included to support certain protocols and formats.
e.ReadFixedString(4) // "test"

// Write 64-bit integer to Buffer (using LEB128 integer compression)
e.WriteInt(12345678)

e.ReadInt() // 12345678

e.Base64() // dGVzdM7C8QU= (url encoding)
           //   t     e     s     t     [    12345678 (LEB128) ]
e.Bytes()  // [ 0x74, 0x65, 0x73, 0x74, 0xce, 0xc2, 0xf1, 0x05 ]
```

## EtcSchema (work in progress)
EtcSchema draws heavy inspiration from Google's Protocol Buffers & gRPC. It uses a sequential encoding and decoding system, in lieu of textual syntax: you just read/write one type, and then move on to the next. 

*Etc is unstable and **not** production ready, so proceed with caution.*

`etc.Buffer` is similar to `bytes.Buffer`, but with number and string methods so you don't have to manually type in `encoding/binary` functions.

it is intended to be fault-tolerant, and as such does not return errors for most reading and writing functions.

## Encoding with EtcSchema

You can define data structures with EtcSchema and make RPC request functions.

`exampleData.etcschema`
```ruby
# Use compression (optional)
# use zlib-compress

struct exampleData {
 uint64   time_ms
 uuid     id
 float32  coordinates[] # dynamic array type
}

rpc exampleRPC {
 requestData(void)     -> exampleData
 postData(exampleData) -> void # use void to declare empty requests and responses
}
```

## Things to be aware of:

- In EtcSchema, the types `int` and `uint` are 64-bit integer types, but using variable length encoding (LEB128). For fixed bit length encoding, use types with the length in the name, e.g. `uint16`. However, this is not recommended as only unsigned ints are supported. For unlimited-length integers, use the `bigint` type.

## Coming soon™

- enums
- imports
- speed increases
- Code generators for JavaScript, Lua, and C
