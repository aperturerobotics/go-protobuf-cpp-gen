// Package genfile provides the generated file abstraction for C++ code generation.
package genfile

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

// Config holds configuration options for the C++ code generator.
type Config struct {
	// RuntimeIncludeBase is the base path for runtime includes.
	// e.g., "vendor/github.com/aperturerobotics/protobuf/src"
	// If empty, uses the standard "google/protobuf" path.
	RuntimeIncludeBase string

	// LiteRuntime uses MessageLite instead of Message (smaller footprint).
	LiteRuntime bool

	// OutputDir is the directory where generated files are written.
	OutputDir string

	// ImportPaths are directories to search for imported .proto files.
	ImportPaths []string

	// SourceRelativeOutput generates files relative to the input .proto file location.
	SourceRelativeOutput bool

	// IncludeSourceInfo includes source code comments in generated files.
	IncludeSourceInfo bool
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		RuntimeIncludeBase: "",
		LiteRuntime:        true,
		OutputDir:          ".",
		ImportPaths:        []string{"."},
	}
}

// RuntimeInclude returns the include path for a runtime header.
func (c *Config) RuntimeInclude(header string) string {
	if c.RuntimeIncludeBase != "" {
		return c.RuntimeIncludeBase + "/" + header
	}
	return header
}

// MessageBaseClass returns the base class for generated messages.
// For lite runtime, we use no base class (standalone structs) to avoid
// the complexity of protobuf's internal class hierarchy.
func (c *Config) MessageBaseClass() string {
	if c.LiteRuntime {
		return "" // No base class for lite runtime
	}
	return "::google::protobuf::Message"
}

// MessageBaseHeader returns the header file for the message base class.
// Returns empty string for lite runtime (no base class).
func (c *Config) MessageBaseHeader() string {
	if c.LiteRuntime {
		return "" // No base class header for lite runtime
	}
	return "google/protobuf/message.h"
}

// CppInclude represents a C++ include directive.
type CppInclude struct {
	Path   string
	System bool // true for <...>, false for "..."
}

// CppIdent represents a C++ identifier that may need namespace qualification.
type CppIdent struct {
	Namespace string // e.g., "foo::bar"
	Name      string // e.g., "MyClass"
}

// String returns the fully qualified name.
func (i CppIdent) String() string {
	if i.Namespace == "" {
		return i.Name
	}
	return "::" + i.Namespace + "::" + i.Name
}

// GeneratedFile represents a file being generated.
type GeneratedFile struct {
	filename string
	buf      bytes.Buffer
	includes map[CppInclude]bool
	Config   *Config
	indent   int
}

// NewGeneratedFile creates a new generated file.
func NewGeneratedFile(filename string, config *Config) *GeneratedFile {
	return &GeneratedFile{
		filename: filename,
		includes: make(map[CppInclude]bool),
		Config:   config,
	}
}

// Filename returns the output filename.
func (g *GeneratedFile) Filename() string {
	return g.filename
}

// P prints a line to the generated file with the current indentation.
// Arguments are printed sequentially. CppIdent values are formatted as fully qualified names.
func (g *GeneratedFile) P(v ...interface{}) {
	for i := 0; i < g.indent; i++ {
		g.buf.WriteString("  ")
	}
	for _, x := range v {
		switch x := x.(type) {
		case CppIdent:
			fmt.Fprint(&g.buf, x.String())
		case *CppIdent:
			fmt.Fprint(&g.buf, x.String())
		default:
			fmt.Fprint(&g.buf, x)
		}
	}
	g.buf.WriteByte('\n')
}

// Pf prints a formatted line to the generated file with the current indentation.
func (g *GeneratedFile) Pf(format string, args ...interface{}) {
	for i := 0; i < g.indent; i++ {
		g.buf.WriteString("  ")
	}
	fmt.Fprintf(&g.buf, format, args...)
	g.buf.WriteByte('\n')
}

// In increases the indentation level.
func (g *GeneratedFile) In() {
	g.indent++
}

// Out decreases the indentation level.
func (g *GeneratedFile) Out() {
	if g.indent > 0 {
		g.indent--
	}
}

// Include adds an include directive to the file.
func (g *GeneratedFile) Include(path string, system bool) {
	g.includes[CppInclude{Path: path, System: system}] = true
}

// IncludeSystem adds a system include (<...>).
func (g *GeneratedFile) IncludeSystem(path string) {
	g.Include(path, true)
}

// IncludeUser adds a user include ("...").
func (g *GeneratedFile) IncludeUser(path string) {
	g.Include(path, false)
}

// IncludeRuntime adds a runtime include using the configured base path.
func (g *GeneratedFile) IncludeRuntime(header string) {
	g.IncludeUser(g.Config.RuntimeInclude(header))
}

// QualifiedName returns the fully qualified C++ name for an identifier.
func (g *GeneratedFile) QualifiedName(ident CppIdent) string {
	return ident.String()
}

// Content returns the complete file content including includes.
func (g *GeneratedFile) Content() []byte {
	var result bytes.Buffer

	// Sort includes: system includes first, then user includes
	var systemIncludes, userIncludes []string
	for inc := range g.includes {
		if inc.System {
			systemIncludes = append(systemIncludes, inc.Path)
		} else {
			userIncludes = append(userIncludes, inc.Path)
		}
	}
	sort.Strings(systemIncludes)
	sort.Strings(userIncludes)

	for _, path := range systemIncludes {
		fmt.Fprintf(&result, "#include <%s>\n", path)
	}
	for _, path := range userIncludes {
		fmt.Fprintf(&result, "#include \"%s\"\n", path)
	}
	if len(systemIncludes) > 0 || len(userIncludes) > 0 {
		result.WriteByte('\n')
	}
	result.Write(g.buf.Bytes())

	return result.Bytes()
}

// String returns the file content as a string.
func (g *GeneratedFile) String() string {
	return string(g.Content())
}

// WriteString writes a raw string to the file without indentation.
func (g *GeneratedFile) WriteString(s string) {
	g.buf.WriteString(s)
}

// Indent returns a string with the current indentation.
func (g *GeneratedFile) Indent() string {
	return strings.Repeat("  ", g.indent)
}
