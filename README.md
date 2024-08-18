# fenneg
Flat encoding code generation for Go.

[![Go Reference](https://pkg.go.dev/badge/github.com/sirkon/fenneg.svg)](https://pkg.go.dev/github.com/sirkon/fenneg)

* [Installation\.](#installation)
* [What it is about\.](#what-it-is-about)
    * [How the final utility is expected to work\.](#how-the-final-utility-is-expected-to-work)
        * [List of types subported out of the box\.](#list-of-types-subported-out-of-the-box)
    * [Auto\-supported types\.](#auto-supported-types)
    * [Custom types\.](#custom-types)
    * [LogRecorder code generation details\.](#logrecorder-code-generation-details)
* [Not just log encoding\.](#not-just-log-encoding)
* [Usage example\.](#usage-example)
* [IMPORTANT: Binary compatibility details\.](#important-binary-compatibility-details)
* [TODO](#todo)




# Installation.

```shell
go get github.com/sirkon/fenneg
```

# What it is about.

It helps to create compact and performant binary encoding and decoding for:

- Operation logs.
- Structures.

This library can be used to create both standalone utilities with CLI or using runners like it is done
in the [example](./example/main.go). Beware though it is not recommended to use this library as a dependency
of your projects directly. An approach with some inner module of your project that will generate
things is preferable, something placed in `internal/tools/`.:q

# Operation logs.

Imagine we have some KV storage. This means we have a snapshot and an operations log with operations records:

```antlr4
Operation:
    : Create(key string, value []byte)
    | Update(key string, value []byte)
    | Delete(key string)
    ;
```

Go's first choice to abstract these operation will be:

```go
// LogOperationsRecorder to write operations into an operation log.
type LogOperationsRecorder interface{
    Create(key string, value []byte) error
    Update(key string, value []byte) error
    Delete(key string) error
}

// LogOperationDispatcher to dispatch operations retrieved from an operation log
type LogOperationDispatcher interface{
    Create(key string, value []byte) error
    Update(key string, value []byte) error
    Delete(key string) error
}
```

And here we have:

- `LogOperationsRecorder` encodes parameters of methods calls into a binary form and save encoded data into a "physical" log.
- `LogOperationDispatcher` methods are being called by a **dispatcher** that decodes an operation retrieved from a physical
  log.

Where:

- `LogOperationsRecorder` implementation seems to be easy enough target for a code generation.
- Likewise, a dispatcher is a reverse for `LogOperationsRecorder` and is also easy enough for a codegen. It is

    ```go
    func logRecordDispatch(disp LogOperationDispatcher, rec []byte) error {
        ...
    }
    ```

- Although `LogOperationDispatcher` is an exact match for the `LogRecorder` as an interface, but it is an actual
  business logic, which is to be written by a user.

So, this code generator is about rendering a code for a `LogOperationRecorder`'s t7y8uio];
\\;'l implementation and a dispatching.

## How the final utility is expected to work.

1. Write an interface `A`.
2. Write a type `B` having two methods a generated code will rely on to encode events:
    - `allocateBuffer(n int) []byte` method returning an empty slice with capacity ≥ n.
    - May be a `writeBuffer(buf []byte) <returnTuple>` method to write encoded events back.
      This method defines returns of encoding methods, they will have the same return tuple
      as `writeBuffer` or will be just `[]byte` otherwise.
2. Write a dispatcher type `C` what implements `A`. This type will be used to handle decoded events.
2. Run utility pointing `A`, `B`, may be `C` (this is optional) and a dispatching function name.
    - Methods to encode events in `<c>_generated.go` file, here `<c>.go` is the file where `C` is defined.
    - Dispatch function `<name>(h <B>|<A>, data []byte) error` in the same `<c>_generated.go` file.
    - When `C` is set the dispatching function will use `C` directly instead of using generic interface `A`.

Arguments of the `LogRecorder` interface are having their own types. Some are supported out of the box, so as types
satisfying certain predefined interface. And you can define your own codegen steps for certain types too.
This kind of customization is a reason why this thing is a framework rather than a ready to use utility.

### List of types subported out of the box.

| type           |
|----------------|
| `bool`         |
| `int8`         |
| `int16`        |
| `int32`        |
| `int64`        |
| `uint8`        |
| `uint16`       |
| `uint32`       |
| `uint64`       |
| `intypes.VI16` |
| `intypes.VI32` |
| `intypes.VI64` |
| `intypes.VU16` |
| `intypes.VU32` |
| `intypes.VU64` |
| `float32`      |
| `float64`      |
| `[N]byte`      |
| `[]byte`       |
| `[][]byte`     |
| `string`       |
| `[]int16`      |
| `[]int32`      |
| `[]int64`      |
| `[]uint16`     |
| `[]uint32`     |
| `[]uint64`     |

Here:

- `intypes.VIX` and `intypes.VUX` are defined in the [sirkon/intypes](https://github.com/sirkon/intypes) package and
  their sole purpose is to represent `int16..64` and `uint16..64` with uleb128 encoding applied rather than a regular
  little endian encoding. I mean, if your `LogRecorder` interface will have, say, `intypes.VU64` argument type in one
  of its methods, the argument type will be replaced to `uint64` in both recorder and handler implementations.

You can also make `fenneg.Chill()` call and then there will be a support for:

| type         | notice               |
|--------------|----------------------|
| `int`        | Treated as `int64`.  |
| `uint`       | Treated as `uint64`. |
| `intypes.VI` |                      |
| `intypes.VU` |                      |

All these types are called `builtins`.

## Auto-supported types.

```go
type Encoder interface{
    Len() int 
    // Encode must append to the dst slice and returns
    // the resulted slice. 
    Encode([]byte) []byte
}

type Decoder interface{
    // its job.
    // Decode returns the rest of the data after it ends
    Decode([]byte) ([]byte, error)
}
```

any type that:

- Satisfies the first interface.
- The type itself or a pointer of the type satisfies the second interface.

Will be handled automatically. Beware though, values of this type must be usable at their zero state.

## Custom types.

You may define custom encoding and decoding for your own types.

You need to implement [Handler](./handler.go) interface and register a handler factory for them using
either CustomHandler option (`HandleByName` or `NewHandler`).

## LogRecorder code generation details.

The general scheme of data encoding is:

``` mermaid
graph TD
    size[Compute an output size of a data]
    allc[Allocate a buffer to keep the encoding data]
    encd[Encode data into the buffer]
    stor[Save buffer]
    
    size --> allc %% Relies on user defined method.
    allc --> encd %% Relies on generated code. 
    encd --> stor %% Relies on user defined method.
```

# Structures binary flat encoding and decoding.

Log encoding/decoding was an original reason to develop this library for, but it turned out
soon I also have something similar for structs too. So, the library was extended to handle them
as well.

# Usage example.

There's an example of the framework usage in the [example](./example/main.go) folder.

# IMPORTANT: Binary compatibility details.

We start from a recorder interface which has a set of (public) methods M<sub>1</sub>, M<sub>2</sub>,…, M<sub>n</sub>,
each having its own set of arguments:

```go
type XXXRecorder interface{
	M1(…)
	…
	Mn(…)
}
```

Each one is getting an encoder. And we have a dispatching procedure, which gets encoded data, decode it and decides
what method to call then. It requires some kind of method reference to be a part of the encoding. The generator
does this that way:

| Method number (uint32 kind) | Arguments encoding |
|-----------------------------|--------------------|

Where method number is its index in the list of methods.

What does it mean? It means we MUST NOT change a method order in any way: no reordering, no insertion, only appends
are allowed.

Arguments are encoded in the order they come in their methods. This means everything what have been said for methods
is applied to them. A little notice on appended arguments though: their encoders should tolerate empty buffer –
this means you can't add arguments with builtin types, they are not like this.

Conclusion:

* You cannot reorder anything generally.
* You can only append methods to the end and append arguments with custom handler which tolerates empty buffer
  on decoding.
* You can rename arguments and methods whatever you like, because only positional information matters for both
  encoding and decoding.

Imagine we have `Op` method in our recorder and want to replace it with updated version. The best approach will be
to rename `Op` -> `DeprecatedOp` and append a new `Op` method. This will do the trick.

You can even remove `DeprecatedOp` arguments altogether at some point, once you are sure there are no records of
the deprecated `Op` in your logs anymore. But don't remove the method or make it private nevertheless, cause the
order.

And it is even harder for structures: you can only append new fields to the end of their list and these fields must
have your custom type that handles no data left case.

# TODO

- [ ] Support types definitions where an underlying type is one of the builtins.
- [ ] Provide support for pointers over numeric, boolean and string types + nil []byte values.
- [ ] Provide auto-support for struct types where all fields are supported.