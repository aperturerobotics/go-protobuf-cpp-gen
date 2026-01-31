// Package generator provides C++ code generation from protobuf descriptors.
package generator

import "github.com/aperturerobotics/go-protobuf-cpp-gen/generator/genfile"

// Config is an alias to genfile.Config for convenience.
type Config = genfile.Config

// DefaultConfig returns a Config with sensible defaults.
var DefaultConfig = genfile.DefaultConfig
