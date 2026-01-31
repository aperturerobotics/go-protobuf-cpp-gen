package generator

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aperturerobotics/go-protobuf-cpp-gen/compiler"
	"github.com/aperturerobotics/go-protobuf-cpp-gen/generator/base"
	"github.com/aperturerobotics/go-protobuf-cpp-gen/generator/genfile"
	"github.com/bufbuild/protocompile/linker"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Generator is the main C++ code generator.
type Generator struct {
	cfg      *Config
	compiler *compiler.Compiler
	features *FeatureRegistry
}

// New creates a new Generator with the given configuration.
func New(config *Config) *Generator {
	if config == nil {
		config = DefaultConfig()
	}
	return &Generator{
		cfg: config,
		compiler: &compiler.Compiler{
			ImportPaths: config.ImportPaths,
		},
		features: NewFeatureRegistry(),
	}
}

// Config returns the generator configuration.
func (g *Generator) Config() *Config {
	return g.cfg
}

// Features returns the feature registry.
func (g *Generator) Features() *FeatureRegistry {
	return g.features
}

// GenerateResult holds the result of code generation.
type GenerateResult struct {
	Files []*OutputFile
}

// OutputFile represents a generated output file.
type OutputFile struct {
	Name    string
	Content []byte
}

// Generate compiles the given .proto files and generates C++ code.
func (g *Generator) Generate(ctx context.Context, files ...string) (*GenerateResult, error) {
	linkedFiles, err := g.compiler.Compile(ctx, files...)
	if err != nil {
		return nil, fmt.Errorf("failed to compile proto files: %w", err)
	}
	return g.GenerateFromDescriptors(linkedFiles)
}

// GenerateFromDescriptors generates C++ code from already-compiled file descriptors.
func (g *Generator) GenerateFromDescriptors(files linker.Files) (*GenerateResult, error) {
	result := &GenerateResult{}
	for _, file := range files {
		headerFile, sourceFile, err := g.generateFile(file)
		if err != nil {
			return nil, fmt.Errorf("failed to generate %s: %w", file.Path(), err)
		}
		result.Files = append(result.Files, headerFile, sourceFile)
	}
	return result, nil
}

func (g *Generator) generateFile(file protoreflect.FileDescriptor) (*OutputFile, *OutputFile, error) {
	baseName := strings.TrimSuffix(file.Path(), ".proto")
	headerName := baseName + ".pb.h"
	sourceName := baseName + ".pb.cc"

	header := genfile.NewGeneratedFile(headerName, g.cfg)
	source := genfile.NewGeneratedFile(sourceName, g.cfg)

	if err := g.generateHeader(header, file); err != nil {
		return nil, nil, fmt.Errorf("header generation failed: %w", err)
	}
	if err := g.generateSource(source, file); err != nil {
		return nil, nil, fmt.Errorf("source generation failed: %w", err)
	}

	return &OutputFile{Name: headerName, Content: header.Content()},
		&OutputFile{Name: sourceName, Content: source.Content()}, nil
}

func (g *Generator) generateHeader(gf *genfile.GeneratedFile, file protoreflect.FileDescriptor) error {
	gf.IncludeSystem("cstdint")
	gf.IncludeSystem("string")
	gf.IncludeSystem("vector")
	gf.IncludeSystem("map")
	if header := g.cfg.MessageBaseHeader(); header != "" {
		gf.IncludeRuntime(header)
	}

	// Include headers for imported proto files
	imports := file.Imports()
	for i := 0; i < imports.Len(); i++ {
		importFile := imports.Get(i).FileDescriptor
		if importFile != nil {
			importBaseName := strings.TrimSuffix(importFile.Path(), ".proto")
			gf.IncludeUser(importBaseName + ".pb.h")
		}
	}

	base.GenerateHeaderPreamble(gf, file)
	base.OpenNamespace(gf, file)
	base.GenerateForwardDeclarations(gf, file)

	enums := file.Enums()
	for i := 0; i < enums.Len(); i++ {
		base.GenerateEnum(gf, enums.Get(i))
	}

	messages := file.Messages()
	for i := 0; i < messages.Len(); i++ {
		if err := g.generateMessageHeader(gf, messages.Get(i)); err != nil {
			return err
		}
	}

	base.CloseNamespace(gf, file)
	base.GenerateHeaderPostamble(gf, file)
	return nil
}

func (g *Generator) generateSource(gf *genfile.GeneratedFile, file protoreflect.FileDescriptor) error {
	baseName := strings.TrimSuffix(file.Path(), ".proto")
	gf.IncludeUser(baseName + ".pb.h")
	gf.IncludeSystem("cstring")
	gf.IncludeRuntime("google/protobuf/io/coded_stream.h")
	gf.IncludeRuntime("google/protobuf/wire_format_lite.h")

	base.GenerateSourcePreamble(gf, file)
	base.OpenNamespace(gf, file)

	messages := file.Messages()
	for i := 0; i < messages.Len(); i++ {
		if err := g.generateMessageSource(gf, messages.Get(i)); err != nil {
			return err
		}
	}

	base.CloseNamespace(gf, file)
	return nil
}

func (g *Generator) generateMessageHeader(gf *genfile.GeneratedFile, msg protoreflect.MessageDescriptor) error {
	if msg.IsMapEntry() {
		return nil
	}

	nestedEnums := msg.Enums()
	for i := 0; i < nestedEnums.Len(); i++ {
		base.GenerateEnum(gf, nestedEnums.Get(i))
	}

	nestedMessages := msg.Messages()
	for i := 0; i < nestedMessages.Len(); i++ {
		if err := g.generateMessageHeader(gf, nestedMessages.Get(i)); err != nil {
			return err
		}
	}

	base.GenerateMessageHeader(gf, msg, g.cfg)
	g.features.GenerateHeader(msg)
	return nil
}

func (g *Generator) generateMessageSource(gf *genfile.GeneratedFile, msg protoreflect.MessageDescriptor) error {
	if msg.IsMapEntry() {
		return nil
	}

	nestedMessages := msg.Messages()
	for i := 0; i < nestedMessages.Len(); i++ {
		if err := g.generateMessageSource(gf, nestedMessages.Get(i)); err != nil {
			return err
		}
	}

	base.GenerateMessageSource(gf, msg, g.cfg)
	g.features.GenerateSource(msg)
	return nil
}

// WriteResult writes the generated files to the output directory.
func (g *Generator) WriteResult(result *GenerateResult) error {
	for _, file := range result.Files {
		outPath := filepath.Join(g.cfg.OutputDir, file.Name)
		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", outPath, err)
		}
		if err := os.WriteFile(outPath, file.Content, 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", outPath, err)
		}
	}
	return nil
}
