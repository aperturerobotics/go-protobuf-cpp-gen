// protoc-gen-cpp-lite is a pure-Go C++ code generator for Protocol Buffers.
//
// It can be used in two modes:
//   1. Standalone mode: uses protocompile to parse .proto files directly
//   2. Plugin mode: receives CodeGeneratorRequest from protoc
//
// Standalone usage:
//
//	protoc-gen-cpp-lite --out=gen --import_path=. file.proto
//
// Plugin usage (with protoc):
//
//	protoc --cpp-lite_out=gen --plugin=protoc-gen-cpp-lite file.proto
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aperturerobotics/protobuf-cpp-gen/compiler"
	"github.com/aperturerobotics/protobuf-cpp-gen/generator"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	// Check if we're being called as a protoc plugin (stdin has data)
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// Plugin mode - stdin has data
		if err := runPlugin(); err != nil {
			fmt.Fprintf(os.Stderr, "protoc-gen-cpp-lite: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Standalone mode
	if err := runStandalone(); err != nil {
		fmt.Fprintf(os.Stderr, "protoc-gen-cpp-lite: %v\n", err)
		os.Exit(1)
	}
}

func runStandalone() error {
	var (
		outDir             string
		importPaths        stringSlice
		runtimeIncludeBase string
		liteRuntime        bool
	)

	flag.StringVar(&outDir, "out", ".", "Output directory for generated files")
	flag.Var(&importPaths, "import_path", "Import path for .proto files (can be specified multiple times)")
	flag.StringVar(&runtimeIncludeBase, "runtime_include", "", "Base path for runtime includes")
	flag.BoolVar(&liteRuntime, "lite", true, "Use lite runtime (MessageLite)")
	flag.Parse()

	files := flag.Args()
	if len(files) == 0 {
		return fmt.Errorf("no input files specified")
	}

	// Add current directory to import paths if empty
	if len(importPaths) == 0 {
		importPaths = []string{"."}
	}

	config := &generator.Config{
		OutputDir:          outDir,
		ImportPaths:        importPaths,
		RuntimeIncludeBase: runtimeIncludeBase,
		LiteRuntime:        liteRuntime,
	}

	gen := generator.New(config)

	ctx := context.Background()
	result, err := gen.Generate(ctx, files...)
	if err != nil {
		return err
	}

	return gen.WriteResult(result)
}

func runPlugin() error {
	// Read the CodeGeneratorRequest from stdin
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	var req pluginpb.CodeGeneratorRequest
	if err := proto.Unmarshal(data, &req); err != nil {
		return fmt.Errorf("failed to parse CodeGeneratorRequest: %w", err)
	}

	// Parse options from the parameter field
	config := parseParameter(req.GetParameter())

	// Compile the proto files using protocompile
	comp := &compiler.Compiler{
		ImportPaths: config.ImportPaths,
	}

	ctx := context.Background()
	linkedFiles, err := comp.Compile(ctx, req.GetFileToGenerate()...)
	if err != nil {
		return writeError(fmt.Sprintf("compilation failed: %v", err))
	}

	// Generate code
	gen := generator.New(config)
	result, err := gen.GenerateFromDescriptors(linkedFiles)
	if err != nil {
		return writeError(fmt.Sprintf("generation failed: %v", err))
	}

	// Build the response
	resp := &pluginpb.CodeGeneratorResponse{
		SupportedFeatures: proto.Uint64(uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)),
	}

	for _, file := range result.Files {
		resp.File = append(resp.File, &pluginpb.CodeGeneratorResponse_File{
			Name:    proto.String(file.Name),
			Content: proto.String(string(file.Content)),
		})
	}

	// Write the response to stdout
	out, err := proto.Marshal(resp)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	if _, err := os.Stdout.Write(out); err != nil {
		return fmt.Errorf("failed to write response: %w", err)
	}

	return nil
}

func writeError(msg string) error {
	resp := &pluginpb.CodeGeneratorResponse{
		Error: proto.String(msg),
	}
	out, err := proto.Marshal(resp)
	if err != nil {
		return err
	}
	_, err = os.Stdout.Write(out)
	return err
}

func parseParameter(param string) *generator.Config {
	config := generator.DefaultConfig()

	if param == "" {
		return config
	}

	for _, opt := range strings.Split(param, ",") {
		parts := strings.SplitN(opt, "=", 2)
		key := parts[0]
		value := ""
		if len(parts) > 1 {
			value = parts[1]
		}

		switch key {
		case "import_path":
			config.ImportPaths = append(config.ImportPaths, value)
		case "runtime_include":
			config.RuntimeIncludeBase = value
		case "lite":
			config.LiteRuntime = value != "false"
		}
	}

	return config
}

// stringSlice is a flag.Value that accumulates multiple string values.
type stringSlice []string

func (s *stringSlice) String() string {
	return strings.Join(*s, ",")
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}
