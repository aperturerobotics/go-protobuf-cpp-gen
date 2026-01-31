# go-protobuf-cpp-gen

[![GoDoc Widget]][GoDoc] [![Go Report Card Widget]][Go Report Card]

> A pure-Go C++ protobuf code generator using protocompile for parsing.

[GoDoc]: https://godoc.org/github.com/aperturerobotics/go-protobuf-cpp-gen
[GoDoc Widget]: https://godoc.org/github.com/aperturerobotics/go-protobuf-cpp-gen?status.svg
[Go Report Card Widget]: https://goreportcard.com/badge/github.com/aperturerobotics/go-protobuf-cpp-gen
[Go Report Card]: https://goreportcard.com/report/github.com/aperturerobotics/go-protobuf-cpp-gen

## Related Projects

- [protobuf-go-lite](https://github.com/aperturerobotics/protobuf-go-lite) - Lightweight Go protobuf without reflection
- [protocompile](https://github.com/bufbuild/protocompile) - Pure-Go protobuf parser
- [protobuf](https://github.com/protocolbuffers/protobuf) - Google's Protocol Buffers

## Purpose

This project generates standalone C++ classes from `.proto` files without requiring the `protoc` compiler. It uses [protocompile] for pure-Go protobuf parsing, making it easy to integrate into Go-based build systems.

[protocompile]: https://github.com/bufbuild/protocompile

### Why Standalone Classes?

Unlike the official `protoc` C++ output which requires `MessageLite` inheritance and the full protobuf runtime, this generator produces:

- **Standalone C++ classes** with no base class inheritance
- **STL containers** (`std::vector`, `std::map`) instead of protobuf containers
- **Minimal dependencies** - only requires protobuf headers for wire format utilities
- **Header-only friendly** - all accessors are inline

This approach is useful for:

- Embedded systems with limited runtime support
- Projects that want protobuf serialization without the full runtime
- Integration with codebases that prefer STL containers

## Features

- Pure-Go implementation using protocompile for parsing
- Generates `.pb.h` headers and `.pb.cc` source files
- Supports all protobuf field types:
  - Scalar types (bool, int32, int64, uint32, uint64, float, double, string, bytes)
  - Signed integers with ZigZag encoding (sint32, sint64)
  - Fixed-width integers (fixed32, fixed64, sfixed32, sfixed64)
  - Enums with `_IsValid`, `_Name`, `_MIN`, `_MAX` helpers
  - Nested messages (arbitrary depth)
  - Repeated fields (scalars and messages)
  - Map fields (all key/value type combinations)
- Full serialization/parsing support via protobuf wire format
- Proper handling of negative int32 values (10-byte varint encoding)

## Installation

```bash
go install github.com/aperturerobotics/go-protobuf-cpp-gen/cmd/protoc-gen-cpp-lite@latest
```

## Usage

### As a Library

```go
package main

import (
    "context"
    "fmt"
    "os"

    "github.com/aperturerobotics/go-protobuf-cpp-gen/generator"
)

func main() {
    config := &generator.Config{
        OutputDir:   "./generated",
        ImportPaths: []string{"./protos"},
        LiteRuntime: true,
    }

    gen := generator.New(config)
    ctx := context.Background()

    result, err := gen.Generate(ctx, "example.proto")
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    for _, file := range result.Files {
        fmt.Printf("Generated: %s (%d bytes)\n", file.Name, len(file.Content))
    }
}
```

### Generated Code Example

Given this proto file:

```protobuf
syntax = "proto3";
package example;

message Person {
    string name = 1;
    int32 age = 2;
    repeated string emails = 3;
}
```

The generator produces:

```cpp
namespace example {

class Person {
 public:
  Person();
  ~Person();
  Person(const Person& other);
  Person(Person&& other) noexcept;
  Person& operator=(const Person& other);
  Person& operator=(Person&& other) noexcept;

  void Clear();

  // Serialization
  size_t ByteSizeLong() const;
  bool SerializeToString(::std::string* output) const;
  bool ParseFromString(const ::std::string& data);

  // name
  inline const ::std::string& name() const { return name_; }
  inline void set_name(const ::std::string& value) { name_ = value; }
  inline void clear_name() { name_.clear(); }
  inline ::std::string* mutable_name() { return &name_; }

  // age
  inline ::int32_t age() const { return age_; }
  inline void set_age(::int32_t value) { age_ = value; }
  inline void clear_age() { age_ = 0; }

  // emails (repeated)
  inline int emails_size() const { return emails_.size(); }
  inline const ::std::string& emails(int index) const { return emails_[index]; }
  inline void add_emails(const ::std::string& value) { emails_.push_back(value); }
  inline void clear_emails() { emails_.clear(); }

 private:
  ::std::string name_;
  ::int32_t age_;
  ::std::vector<::std::string> emails_;
};

}  // namespace example
```

## Building with CMake

The generated code requires protobuf headers for wire format utilities:

```cmake
find_package(Protobuf REQUIRED)

add_library(my_protos
    generated/example.pb.cc
)
target_link_libraries(my_protos ${Protobuf_LIBRARIES})
target_include_directories(my_protos PUBLIC ${Protobuf_INCLUDE_DIRS})
```

## Testing

Run the Go tests:

```bash
go test ./...
```

Run the end-to-end C++ tests (requires CMake and protobuf installed):

```bash
go test -v ./tests/e2e/...
```

This will:
1. Generate C++ code from test proto files
2. Build the generated code with CMake
3. Run C++ test executables that verify accessors, serialization, and parsing

## License

MIT
