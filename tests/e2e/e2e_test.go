// Package e2e provides end-to-end tests for the protobuf-cpp-gen code generator.
package e2e

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/aperturerobotics/go-protobuf-cpp-gen/generator"
)

// testDir returns the path to the e2e test directory.
func testDir(t *testing.T) string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to get test directory")
	}
	return filepath.Dir(file)
}

// TestGenerateProtos tests that all proto files can be parsed and generated.
func TestGenerateProtos(t *testing.T) {
	dir := testDir(t)
	protoDir := filepath.Join(dir, "protos")
	outDir := filepath.Join(dir, "generated")

	// Clean up generated directory
	os.RemoveAll(outDir)
	if err := os.MkdirAll(outDir, 0755); err != nil {
		t.Fatalf("failed to create output directory: %v", err)
	}

	// Proto files to generate
	protoFiles := []string{
		"scalars.proto",
		"enums.proto",
		"repeated.proto",
		"nested.proto",
		"maps.proto",
		"common.proto",
		"imports.proto",
	}

	config := &generator.Config{
		OutputDir:   outDir,
		ImportPaths: []string{protoDir},
		LiteRuntime: true,
	}

	gen := generator.New(config)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, proto := range protoFiles {
		t.Run(proto, func(t *testing.T) {
			result, err := gen.Generate(ctx, proto)
			if err != nil {
				t.Fatalf("failed to generate %s: %v", proto, err)
			}

			if len(result.Files) != 2 {
				t.Fatalf("expected 2 files (header + source), got %d", len(result.Files))
			}

			// Write the generated files
			for _, file := range result.Files {
				outPath := filepath.Join(outDir, filepath.Base(file.Name))
				if err := os.WriteFile(outPath, file.Content, 0644); err != nil {
					t.Fatalf("failed to write %s: %v", outPath, err)
				}
				t.Logf("Generated: %s (%d bytes)", outPath, len(file.Content))
			}
		})
	}
}

// TestGeneratedCodeCompiles tests that the generated C++ code compiles with CMake.
func TestGeneratedCodeCompiles(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping compilation test in short mode")
	}

	// Check for cmake
	if _, err := exec.LookPath("cmake"); err != nil {
		t.Skip("cmake not found, skipping compilation test")
	}

	// Check for protobuf
	if !hasProtobuf() {
		t.Skip("protobuf not found, skipping compilation test")
	}

	dir := testDir(t)
	buildDir := filepath.Join(dir, "build")

	// Clean up build directory
	os.RemoveAll(buildDir)
	if err := os.MkdirAll(buildDir, 0755); err != nil {
		t.Fatalf("failed to create build directory: %v", err)
	}

	// First, generate the proto files
	t.Run("Generate", func(t *testing.T) {
		TestGenerateProtos(t)
	})

	// Run cmake configure
	t.Run("CMakeConfigure", func(t *testing.T) {
		cmd := exec.Command("cmake", "..")
		cmd.Dir = buildDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			t.Fatalf("cmake configure failed: %v", err)
		}
	})

	// Run cmake build
	t.Run("CMakeBuild", func(t *testing.T) {
		cmd := exec.Command("cmake", "--build", ".", "--parallel")
		cmd.Dir = buildDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			t.Fatalf("cmake build failed: %v", err)
		}
	})
}

// TestGeneratedCodeRuns tests that the generated C++ tests pass.
func TestGeneratedCodeRuns(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping run tests in short mode")
	}

	dir := testDir(t)
	buildDir := filepath.Join(dir, "build")

	// Check if build directory exists and has executables
	if _, err := os.Stat(buildDir); os.IsNotExist(err) {
		t.Skip("build directory not found, run TestGeneratedCodeCompiles first")
	}

	tests := []string{
		"compile_test",
		"scalar_test",
		"enum_test",
		"repeated_test",
		"nested_test",
		"map_test",
		"serialize_test",
	}

	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			execPath := filepath.Join(buildDir, test)
			if runtime.GOOS == "windows" {
				execPath += ".exe"
			}

			if _, err := os.Stat(execPath); os.IsNotExist(err) {
				t.Skipf("test executable %s not found", test)
			}

			cmd := exec.Command(execPath)
			cmd.Dir = buildDir

			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			if err := cmd.Run(); err != nil {
				t.Fatalf("test %s failed:\nstdout: %s\nstderr: %s\nerror: %v",
					test, stdout.String(), stderr.String(), err)
			}

			t.Logf("Output:\n%s", stdout.String())
		})
	}
}

// TestCTestIntegration runs the tests using ctest.
func TestCTestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping ctest in short mode")
	}

	dir := testDir(t)
	buildDir := filepath.Join(dir, "build")

	// Check if build directory exists
	if _, err := os.Stat(buildDir); os.IsNotExist(err) {
		t.Skip("build directory not found, run TestGeneratedCodeCompiles first")
	}

	// Check for ctest
	if _, err := exec.LookPath("ctest"); err != nil {
		t.Skip("ctest not found")
	}

	cmd := exec.Command("ctest", "--output-on-failure")
	cmd.Dir = buildDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("ctest failed: %v", err)
	}
}

// TestGeneratedHeaderStructure verifies the structure of generated headers.
func TestGeneratedHeaderStructure(t *testing.T) {
	dir := testDir(t)
	outDir := filepath.Join(dir, "generated")

	// Generate if not already done
	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		t.Run("Generate", func(t *testing.T) {
			TestGenerateProtos(t)
		})
	}

	tests := []struct {
		file     string
		contains []string
	}{
		{
			file: "scalars.pb.h",
			contains: []string{
				"#ifndef",
				"#define",
				"namespace test {",
				"namespace scalars {",
				"class AllScalars",
				"class DefaultValues",
				"bool_val()",
				"set_bool_val(",
				"int32_val()",
				"set_int32_val(",
				"string_val()",
				"set_string_val(",
				"mutable_string_val()",
				"clear_bool_val()",
				"ByteSizeLong()",
				"SerializeToString(",
				"ParseFromString(",
				"#endif",
			},
		},
		{
			file: "enums.pb.h",
			contains: []string{
				"enum Status : int",
				"STATUS_UNKNOWN = 0",
				"STATUS_ACTIVE = 1",
				"Status_IsValid(",
				"Status_Name(",
				"Status_MIN",
				"Status_MAX",
				"enum Priority : int",
				"PRIORITY_LOW = 10",
				"class EnumMessage",
				"class MessageWithNestedEnum",
				"enum NestedStatus",
			},
		},
		{
			file: "repeated.pb.h",
			contains: []string{
				"class RepeatedScalars",
				"std::vector<",
				"bool_vals_size()",
				"add_bool_vals(",
				"clear_bool_vals()",
				"class Item",
				"class RepeatedMessages",
				"add_items()",
				"items_size()",
				"mutable_items(",
			},
		},
		{
			file: "nested.pb.h",
			contains: []string{
				"class Outer",
				"class Middle",
				"class Inner",
				"mutable_middle()",
				"mutable_inner()",
				"class Document",
				"class Header",
				"class Body",
				"class Section",
				"class Footer",
				"class Level1",
				"class Level2",
				"class Level3",
				"class Level4",
			},
		},
		{
			file: "maps.pb.h",
			contains: []string{
				"class MapScalars",
				"std::map<",
				"string_to_string()",
				"mutable_string_to_string()",
				"string_to_string_size()",
				"clear_string_to_string()",
				"class MapMessages",
				"class MapValue",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.file, func(t *testing.T) {
			content, err := os.ReadFile(filepath.Join(outDir, tc.file))
			if err != nil {
				t.Fatalf("failed to read %s: %v", tc.file, err)
			}

			contentStr := string(content)
			for _, expected := range tc.contains {
				if !strings.Contains(contentStr, expected) {
					t.Errorf("expected %s to contain %q", tc.file, expected)
				}
			}
		})
	}
}

// TestGeneratedSourceStructure verifies the structure of generated source files.
func TestGeneratedSourceStructure(t *testing.T) {
	dir := testDir(t)
	outDir := filepath.Join(dir, "generated")

	// Generate if not already done
	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		t.Run("Generate", func(t *testing.T) {
			TestGenerateProtos(t)
		})
	}

	tests := []struct {
		file     string
		contains []string
	}{
		{
			file: "scalars.pb.cc",
			contains: []string{
				`#include "scalars.pb.h"`,
				"AllScalars::AllScalars()",
				"AllScalars::~AllScalars()",
				"AllScalars::Clear()",
				"AllScalars::ByteSizeLong()",
				"AllScalars::SerializeToString(",
				"AllScalars::SerializeToArray(",
				"AllScalars::ParseFromString(",
				"AllScalars::ParseFromArray(",
				"operator=(",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.file, func(t *testing.T) {
			content, err := os.ReadFile(filepath.Join(outDir, tc.file))
			if err != nil {
				t.Fatalf("failed to read %s: %v", tc.file, err)
			}

			contentStr := string(content)
			for _, expected := range tc.contains {
				if !strings.Contains(contentStr, expected) {
					t.Errorf("expected %s to contain %q", tc.file, expected)
				}
			}
		})
	}
}

// hasProtobuf checks if protobuf is available for compilation.
func hasProtobuf() bool {
	// Try pkg-config first
	cmd := exec.Command("pkg-config", "--exists", "protobuf")
	if err := cmd.Run(); err == nil {
		return true
	}

	// Try cmake find
	tmpDir, err := os.MkdirTemp("", "protobuf-check")
	if err != nil {
		return false
	}
	defer os.RemoveAll(tmpDir)

	cmakeContent := `
cmake_minimum_required(VERSION 3.14)
project(protobuf_check)
find_package(Protobuf REQUIRED)
`
	cmakePath := filepath.Join(tmpDir, "CMakeLists.txt")
	if err := os.WriteFile(cmakePath, []byte(cmakeContent), 0644); err != nil {
		return false
	}

	buildDir := filepath.Join(tmpDir, "build")
	os.MkdirAll(buildDir, 0755)

	cmd = exec.Command("cmake", "..")
	cmd.Dir = buildDir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("protobuf check failed: %s\n", stderr.String())
		return false
	}

	return true
}

// BenchmarkGeneration benchmarks the code generation.
func BenchmarkGeneration(b *testing.B) {
	dir := testDir(&testing.T{})
	protoDir := filepath.Join(dir, "protos")

	config := &generator.Config{
		OutputDir:   b.TempDir(),
		ImportPaths: []string{protoDir},
		LiteRuntime: true,
	}

	gen := generator.New(config)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := gen.Generate(ctx, "scalars.proto")
		if err != nil {
			b.Fatalf("generation failed: %v", err)
		}
	}
}
