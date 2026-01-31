package generator_test

import (
	"context"
	"strings"
	"testing"

	"github.com/aperturerobotics/protobuf-cpp-gen/compiler"
	"github.com/aperturerobotics/protobuf-cpp-gen/generator"
)

func TestGenerateSimpleMessage(t *testing.T) {
	config := &generator.Config{
		OutputDir:   ".",
		ImportPaths: []string{"../testdata"},
		LiteRuntime: true,
	}

	gen := generator.New(config)

	ctx := context.Background()
	result, err := gen.Generate(ctx, "simple.proto")
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if len(result.Files) == 0 {
		t.Fatal("No files generated")
	}

	// Should have header and source
	if len(result.Files) != 2 {
		t.Fatalf("Expected 2 files, got %d", len(result.Files))
	}

	var headerFile, sourceFile *generator.OutputFile
	for _, f := range result.Files {
		if strings.HasSuffix(f.Name, ".pb.h") {
			headerFile = f
		} else if strings.HasSuffix(f.Name, ".pb.cc") {
			sourceFile = f
		}
	}

	if headerFile == nil {
		t.Fatal("No header file generated")
	}
	if sourceFile == nil {
		t.Fatal("No source file generated")
	}

	// Verify header content
	headerContent := string(headerFile.Content)

	// Check for header guard
	if !strings.Contains(headerContent, "#ifndef") || !strings.Contains(headerContent, "#define") {
		t.Error("Header guard not found")
	}

	// Check for namespace
	if !strings.Contains(headerContent, "namespace example {") {
		t.Error("Namespace declaration not found")
	}

	// Check for SimpleMessage class
	if !strings.Contains(headerContent, "class SimpleMessage") {
		t.Error("SimpleMessage class not found")
	}

	// Check for MessageLite base class
	if !strings.Contains(headerContent, "MessageLite") {
		t.Error("MessageLite base class not found")
	}

	// Check for field accessors
	if !strings.Contains(headerContent, "name()") {
		t.Error("name() accessor not found")
	}
	if !strings.Contains(headerContent, "set_name(") {
		t.Error("set_name() not found")
	}

	// Verify source content
	sourceContent := string(sourceFile.Content)

	// Check for include of header
	if !strings.Contains(sourceContent, `#include "simple.pb.h"`) {
		t.Error("Header include not found in source")
	}

	// Check for constructor
	if !strings.Contains(sourceContent, "SimpleMessage::SimpleMessage()") {
		t.Error("Constructor not found")
	}

	// Check for Clear method
	if !strings.Contains(sourceContent, "void SimpleMessage::Clear()") {
		t.Error("Clear() method not found")
	}

	t.Logf("Header file content:\n%s", headerContent)
	t.Logf("Source file content:\n%s", sourceContent)
}

func TestGenerateEnum(t *testing.T) {
	config := &generator.Config{
		OutputDir:   ".",
		ImportPaths: []string{"../testdata"},
		LiteRuntime: true,
	}

	gen := generator.New(config)

	ctx := context.Background()
	result, err := gen.Generate(ctx, "simple.proto")
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	var headerFile *generator.OutputFile
	for _, f := range result.Files {
		if strings.HasSuffix(f.Name, ".pb.h") {
			headerFile = f
			break
		}
	}

	if headerFile == nil {
		t.Fatal("No header file generated")
	}

	headerContent := string(headerFile.Content)

	// Check for enum definition
	if !strings.Contains(headerContent, "enum Status : int {") {
		t.Error("Status enum not found")
	}

	// Check for enum values
	if !strings.Contains(headerContent, "UNKNOWN = 0") {
		t.Error("UNKNOWN enum value not found")
	}
	if !strings.Contains(headerContent, "ACTIVE = 1") {
		t.Error("ACTIVE enum value not found")
	}

	// Check for enum helper functions
	if !strings.Contains(headerContent, "Status_IsValid") {
		t.Error("Status_IsValid not found")
	}
	if !strings.Contains(headerContent, "Status_Name") {
		t.Error("Status_Name not found")
	}
}

func TestGenerateRepeatedFields(t *testing.T) {
	config := &generator.Config{
		OutputDir:   ".",
		ImportPaths: []string{"../testdata"},
		LiteRuntime: true,
	}

	gen := generator.New(config)

	ctx := context.Background()
	result, err := gen.Generate(ctx, "simple.proto")
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	var headerFile *generator.OutputFile
	for _, f := range result.Files {
		if strings.HasSuffix(f.Name, ".pb.h") {
			headerFile = f
			break
		}
	}

	if headerFile == nil {
		t.Fatal("No header file generated")
	}

	headerContent := string(headerFile.Content)

	// Check for RepeatedFields class
	if !strings.Contains(headerContent, "class RepeatedFields") {
		t.Error("RepeatedFields class not found")
	}

	// Check for repeated field accessors
	if !strings.Contains(headerContent, "names_size()") {
		t.Error("names_size() not found")
	}
	if !strings.Contains(headerContent, "add_names(") {
		t.Error("add_names() not found")
	}

	// Check for RepeatedField/RepeatedPtrField storage
	if !strings.Contains(headerContent, "RepeatedPtrField") {
		t.Error("RepeatedPtrField not found")
	}
	if !strings.Contains(headerContent, "RepeatedField") {
		t.Error("RepeatedField not found")
	}
}

func TestGenerateNestedMessage(t *testing.T) {
	config := &generator.Config{
		OutputDir:   ".",
		ImportPaths: []string{"../testdata"},
		LiteRuntime: true,
	}

	gen := generator.New(config)

	ctx := context.Background()
	result, err := gen.Generate(ctx, "simple.proto")
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	var headerFile *generator.OutputFile
	for _, f := range result.Files {
		if strings.HasSuffix(f.Name, ".pb.h") {
			headerFile = f
			break
		}
	}

	if headerFile == nil {
		t.Fatal("No header file generated")
	}

	headerContent := string(headerFile.Content)

	// Check for Outer class
	if !strings.Contains(headerContent, "class Outer") {
		t.Error("Outer class not found")
	}

	// Check for Inner class (nested)
	if !strings.Contains(headerContent, "class Inner") {
		t.Error("Inner class not found")
	}

	// Check for mutable_inner accessor
	if !strings.Contains(headerContent, "mutable_inner()") {
		t.Error("mutable_inner() not found")
	}
}

func TestGenerateMapField(t *testing.T) {
	config := &generator.Config{
		OutputDir:   ".",
		ImportPaths: []string{"../testdata"},
		LiteRuntime: true,
	}

	gen := generator.New(config)

	ctx := context.Background()
	result, err := gen.Generate(ctx, "simple.proto")
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	var headerFile *generator.OutputFile
	for _, f := range result.Files {
		if strings.HasSuffix(f.Name, ".pb.h") {
			headerFile = f
			break
		}
	}

	if headerFile == nil {
		t.Fatal("No header file generated")
	}

	headerContent := string(headerFile.Content)

	// Check for MapMessage class
	if !strings.Contains(headerContent, "class MapMessage") {
		t.Error("MapMessage class not found")
	}

	// Check for Map type
	if !strings.Contains(headerContent, "::google::protobuf::Map<") {
		t.Error("Map type not found")
	}
}

func TestCompilerIntegration(t *testing.T) {
	comp := &compiler.Compiler{
		ImportPaths: []string{"../testdata"},
	}

	ctx := context.Background()
	files, err := comp.Compile(ctx, "simple.proto")
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	if len(files) == 0 {
		t.Fatal("No files compiled")
	}

	// Verify the file was parsed correctly
	file := files[0]
	if file.Path() != "simple.proto" {
		t.Errorf("Expected file path simple.proto, got %s", file.Path())
	}

	// Check package
	if string(file.Package()) != "example" {
		t.Errorf("Expected package example, got %s", file.Package())
	}

	// Check messages
	messages := file.Messages()
	if messages.Len() == 0 {
		t.Error("No messages found in compiled file")
	}

	foundSimple := false
	for i := 0; i < messages.Len(); i++ {
		if string(messages.Get(i).Name()) == "SimpleMessage" {
			foundSimple = true
			break
		}
	}
	if !foundSimple {
		t.Error("SimpleMessage not found in compiled file")
	}
}
