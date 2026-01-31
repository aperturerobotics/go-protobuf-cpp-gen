package types

import (
	"testing"

	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestScalarTypeToCpp(t *testing.T) {
	tests := []struct {
		kind     protoreflect.Kind
		expected string
	}{
		{protoreflect.BoolKind, "bool"},
		{protoreflect.Int32Kind, "::int32_t"},
		{protoreflect.Int64Kind, "::int64_t"},
		{protoreflect.Uint32Kind, "::uint32_t"},
		{protoreflect.Uint64Kind, "::uint64_t"},
		{protoreflect.Sint32Kind, "::int32_t"},
		{protoreflect.Sint64Kind, "::int64_t"},
		{protoreflect.Fixed32Kind, "::uint32_t"},
		{protoreflect.Fixed64Kind, "::uint64_t"},
		{protoreflect.Sfixed32Kind, "::int32_t"},
		{protoreflect.Sfixed64Kind, "::int64_t"},
		{protoreflect.FloatKind, "float"},
		{protoreflect.DoubleKind, "double"},
		{protoreflect.StringKind, "::std::string"},
		{protoreflect.BytesKind, "::std::string"},
	}

	for _, tt := range tests {
		t.Run(tt.kind.String(), func(t *testing.T) {
			result := ScalarTypeToCpp(tt.kind)
			if result.Name != tt.expected {
				t.Errorf("ScalarTypeToCpp(%v).Name = %q, want %q", tt.kind, result.Name, tt.expected)
			}
		})
	}
}

func TestDefaultValue(t *testing.T) {
	tests := []struct {
		kind     protoreflect.Kind
		expected string
	}{
		{protoreflect.BoolKind, "false"},
		{protoreflect.Int32Kind, "0"},
		{protoreflect.Int64Kind, "0"},
		{protoreflect.Uint32Kind, "0u"},
		{protoreflect.Uint64Kind, "0u"},
		{protoreflect.FloatKind, "0.0f"},
		{protoreflect.DoubleKind, "0.0"},
		{protoreflect.StringKind, `""`},
		{protoreflect.BytesKind, `""`},
		{protoreflect.EnumKind, "0"},
		{protoreflect.MessageKind, "nullptr"},
	}

	for _, tt := range tests {
		t.Run(tt.kind.String(), func(t *testing.T) {
			result := DefaultValue(tt.kind)
			if result != tt.expected {
				t.Errorf("DefaultValue(%v) = %q, want %q", tt.kind, result, tt.expected)
			}
		})
	}
}

func TestIsScalarType(t *testing.T) {
	scalarKinds := []protoreflect.Kind{
		protoreflect.BoolKind,
		protoreflect.Int32Kind, protoreflect.Int64Kind,
		protoreflect.Uint32Kind, protoreflect.Uint64Kind,
		protoreflect.Sint32Kind, protoreflect.Sint64Kind,
		protoreflect.Fixed32Kind, protoreflect.Fixed64Kind,
		protoreflect.Sfixed32Kind, protoreflect.Sfixed64Kind,
		protoreflect.FloatKind, protoreflect.DoubleKind,
		protoreflect.StringKind, protoreflect.BytesKind,
	}

	nonScalarKinds := []protoreflect.Kind{
		protoreflect.EnumKind,
		protoreflect.MessageKind,
		protoreflect.GroupKind,
	}

	for _, kind := range scalarKinds {
		if !IsScalarType(kind) {
			t.Errorf("IsScalarType(%v) = false, want true", kind)
		}
	}

	for _, kind := range nonScalarKinds {
		if IsScalarType(kind) {
			t.Errorf("IsScalarType(%v) = true, want false", kind)
		}
	}
}

func TestIsNumericType(t *testing.T) {
	numericKinds := []protoreflect.Kind{
		protoreflect.Int32Kind, protoreflect.Int64Kind,
		protoreflect.Uint32Kind, protoreflect.Uint64Kind,
		protoreflect.Sint32Kind, protoreflect.Sint64Kind,
		protoreflect.Fixed32Kind, protoreflect.Fixed64Kind,
		protoreflect.Sfixed32Kind, protoreflect.Sfixed64Kind,
		protoreflect.FloatKind, protoreflect.DoubleKind,
	}

	nonNumericKinds := []protoreflect.Kind{
		protoreflect.BoolKind,
		protoreflect.StringKind, protoreflect.BytesKind,
		protoreflect.EnumKind, protoreflect.MessageKind,
	}

	for _, kind := range numericKinds {
		if !IsNumericType(kind) {
			t.Errorf("IsNumericType(%v) = false, want true", kind)
		}
	}

	for _, kind := range nonNumericKinds {
		if IsNumericType(kind) {
			t.Errorf("IsNumericType(%v) = true, want false", kind)
		}
	}
}

func TestIsStringType(t *testing.T) {
	if !IsStringType(protoreflect.StringKind) {
		t.Error("IsStringType(StringKind) = false, want true")
	}
	if !IsStringType(protoreflect.BytesKind) {
		t.Error("IsStringType(BytesKind) = false, want true")
	}
	if IsStringType(protoreflect.Int32Kind) {
		t.Error("IsStringType(Int32Kind) = true, want false")
	}
}

func TestGetWireType(t *testing.T) {
	tests := []struct {
		kind     protoreflect.Kind
		expected WireType
	}{
		{protoreflect.BoolKind, WireVarint},
		{protoreflect.Int32Kind, WireVarint},
		{protoreflect.Int64Kind, WireVarint},
		{protoreflect.Uint32Kind, WireVarint},
		{protoreflect.Uint64Kind, WireVarint},
		{protoreflect.Sint32Kind, WireVarint},
		{protoreflect.Sint64Kind, WireVarint},
		{protoreflect.EnumKind, WireVarint},
		{protoreflect.Fixed64Kind, WireFixed64},
		{protoreflect.Sfixed64Kind, WireFixed64},
		{protoreflect.DoubleKind, WireFixed64},
		{protoreflect.Fixed32Kind, WireFixed32},
		{protoreflect.Sfixed32Kind, WireFixed32},
		{protoreflect.FloatKind, WireFixed32},
		{protoreflect.StringKind, WireBytes},
		{protoreflect.BytesKind, WireBytes},
		{protoreflect.MessageKind, WireBytes},
	}

	for _, tt := range tests {
		t.Run(tt.kind.String(), func(t *testing.T) {
			result := GetWireType(tt.kind)
			if result != tt.expected {
				t.Errorf("GetWireType(%v) = %d, want %d", tt.kind, result, tt.expected)
			}
		})
	}
}
