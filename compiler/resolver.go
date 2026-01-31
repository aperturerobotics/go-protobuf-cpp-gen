package compiler

import (
	"context"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/bufbuild/protocompile"
	"github.com/bufbuild/protocompile/linker"
)

// VendoredResolver resolves imports with support for vendored paths.
// It handles mappings like "google/protobuf/..." to vendored locations.
type VendoredResolver struct {
	// ImportPaths are directories to search for imported .proto files.
	ImportPaths []string

	// VendorMappings maps import prefixes to vendored paths.
	// For example: "google/protobuf" -> "vendor/github.com/aperturerobotics/protobuf/src/google/protobuf"
	VendorMappings map[string]string
}

// FindFileByPath implements protocompile.Resolver.
func (r *VendoredResolver) FindFileByPath(path string) (protocompile.SearchResult, error) {
	// Try vendor mappings first
	for prefix, vendorPath := range r.VendorMappings {
		if strings.HasPrefix(path, prefix) {
			suffix := strings.TrimPrefix(path, prefix)
			mappedPath := vendorPath + suffix
			for _, importPath := range r.ImportPaths {
				fullPath := filepath.Join(importPath, mappedPath)
				if info, err := os.Stat(fullPath); err == nil && !info.IsDir() {
					return protocompile.SearchResult{Source: &fileOpener{path: fullPath}}, nil
				}
			}
		}
	}

	// Fall back to normal import path resolution
	for _, importPath := range r.ImportPaths {
		fullPath := filepath.Join(importPath, path)
		if info, err := os.Stat(fullPath); err == nil && !info.IsDir() {
			return protocompile.SearchResult{Source: &fileOpener{path: fullPath}}, nil
		}
	}

	return protocompile.SearchResult{}, fs.ErrNotExist
}

// fileOpener implements io.Reader and returns file contents.
type fileOpener struct {
	path string
	file *os.File
}

func (f *fileOpener) Read(p []byte) (n int, err error) {
	if f.file == nil {
		f.file, err = os.Open(f.path)
		if err != nil {
			return 0, err
		}
	}
	return f.file.Read(p)
}

func (f *fileOpener) Close() error {
	if f.file != nil {
		return f.file.Close()
	}
	return nil
}

// MultiResolver combines multiple resolvers, trying each in order.
type MultiResolver struct {
	Resolvers []protocompile.Resolver
}

// FindFileByPath implements protocompile.Resolver.
func (r *MultiResolver) FindFileByPath(path string) (protocompile.SearchResult, error) {
	for _, resolver := range r.Resolvers {
		result, err := resolver.FindFileByPath(path)
		if err == nil {
			return result, nil
		}
	}
	return protocompile.SearchResult{}, fs.ErrNotExist
}

// ReaderAt implements the interface needed by protocompile for source files.
type ReaderAt interface {
	io.Reader
	io.Closer
}

// NewSourceResolver creates a protocompile.SourceResolver with the given import paths.
func NewSourceResolver(importPaths []string) *protocompile.SourceResolver {
	return &protocompile.SourceResolver{
		ImportPaths: importPaths,
	}
}

// CompileWithResolver compiles .proto files using a custom resolver.
func CompileWithResolver(ctx context.Context, resolver protocompile.Resolver, files ...string) (linker.Files, error) {
	compiler := protocompile.Compiler{
		Resolver:       protocompile.WithStandardImports(resolver),
		SourceInfoMode: protocompile.SourceInfoStandard,
	}
	return compiler.Compile(ctx, files...)
}
