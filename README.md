# Etc: Efficient Transfer Coding

Etc is an encoding format. You can use its encoding functions manually, or you can use the format description language I call EtcSchema.

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

- imports
- speed increases
- Code generators for JavaScript, Lua, and C
