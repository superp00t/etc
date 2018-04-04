# Etc: Efficient Transfer Coding

Etc is an encoding system. You can write manual encoding functions with it, or you can use my format language called EtcSchema.

EtcSchema draws heavy inspiration from Google's Protocol Buffers & gRPC. It uses an offset-based system, in lieu of syntax: you just read one type, and move on to the next. The tokenization of data relies on types remaining a fixed size, or having some kind of termination header: for example, with C strings, you write a string and append a NULL character to the end. With an LEB128 integer, it is terminated by the eigth bit.

Etc is definitely unstable and **NOT** production ready, so proceed with caution.

`etc.Buffer` is similar to `bytes.Buffer`, but with number and string methods so you don't have to manually type in `encoding/binary` functions.

it is intended to be fault-tolerant, and as such does not return errors for most reading and writing functions.

## Encoding with EtcSchema

You can define data structures with EtcSchema and make RPC request functions.

`exampleData.etcschema`
```ruby
# Use compression (optional)
# use zlib-compress

struct exampleData {
  uint64      time_ms
  uuid        id
  float32     coordinates[]
}

rpc exampleRPC {
	requestData(void) -> exampleData
	postData(exampleData) -> void # use void to declare empty requests and responses
}
```

## Things to be aware of:

- In EtcSchema, the types `int` and `uint` are 64-bit integer types, but using variable length encoding (LEB128). For fixed bit length encoding, use types with the length in the name, e.g. `uint16`. However, this is not recommended as signed ints are not supported this way. For true, variable-length integers, use the `bigint` type.

## Coming soon™

- imports
- speed increases
- Code generators for JavaScript, Lua, and C
