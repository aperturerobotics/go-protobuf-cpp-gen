// Package compiler provides a wrapper around protocompile for parsing .proto files.
package compiler

import (
	"context"

	"github.com/bufbuild/protocompile"
	"github.com/bufbuild/protocompile/linker"
)

// Compiler wraps protocompile to provide .proto file parsing and linking.
type Compiler struct {
	// ImportPaths are directories to search for imported .proto files.
	ImportPaths []string

	// Resolver is a custom resolver for locating .proto files.
	// If nil, a SourceResolver using ImportPaths is used.
	Resolver protocompile.Resolver
}

// Compile parses and links the given .proto files, returning linked file descriptors.
func (c *Compiler) Compile(ctx context.Context, files ...string) (linker.Files, error) {
	resolver := c.Resolver
	if resolver == nil {
		resolver = &protocompile.SourceResolver{
			ImportPaths: c.ImportPaths,
		}
	}

	compiler := protocompile.Compiler{
		Resolver:       protocompile.WithStandardImports(resolver),
		SourceInfoMode: protocompile.SourceInfoStandard,
	}
	return compiler.Compile(ctx, files...)
}
