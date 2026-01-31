// Package types provides proto to C++ type mapping.
package types

import (
	"google.golang.org/protobuf/reflect/protoreflect"
)

// CppType represents a C++ type with its include requirements.
type CppType struct {
	// Name is the C++ type name (e.g., "::std::string", "::int32_t").
	Name string

	// Include is the header file to include for this type (e.g., "<string>", "<cstdint>").
	Include string

	// IsSystem indicates if Include is a system include (<...>) vs user include ("...").
	IsSystem bool
}

// ScalarTypeToCpp maps protobuf scalar types to C++ types.
func ScalarTypeToCpp(kind protoreflect.Kind) CppType {
	switch kind {
	case protoreflect.BoolKind:
		return CppType{Name: "bool", Include: "", IsSystem: true}
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return CppType{Name: "::int32_t", Include: "<cstdint>", IsSystem: true}
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return CppType{Name: "::int64_t", Include: "<cstdint>", IsSystem: true}
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return CppType{Name: "::uint32_t", Include: "<cstdint>", IsSystem: true}
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return CppType{Name: "::uint64_t", Include: "<cstdint>", IsSystem: true}
	case protoreflect.FloatKind:
		return CppType{Name: "float", Include: "", IsSystem: true}
	case protoreflect.DoubleKind:
		return CppType{Name: "double", Include: "", IsSystem: true}
	case protoreflect.StringKind:
		return CppType{Name: "::std::string", Include: "<string>", IsSystem: true}
	case protoreflect.BytesKind:
		return CppType{Name: "::std::string", Include: "<string>", IsSystem: true}
	case protoreflect.EnumKind:
		// Enum type name is determined by the enum descriptor
		return CppType{Name: "", Include: "", IsSystem: false}
	case protoreflect.MessageKind, protoreflect.GroupKind:
		// Message type name is determined by the message descriptor
		return CppType{Name: "", Include: "", IsSystem: false}
	default:
		return CppType{Name: "void", Include: "", IsSystem: true}
	}
}

// DefaultValue returns the C++ default value for a protobuf scalar type.
func DefaultValue(kind protoreflect.Kind) string {
	switch kind {
	case protoreflect.BoolKind:
		return "false"
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return "0"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return "0"
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return "0u"
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return "0u"
	case protoreflect.FloatKind:
		return "0.0f"
	case protoreflect.DoubleKind:
		return "0.0"
	case protoreflect.StringKind, protoreflect.BytesKind:
		return `""`
	case protoreflect.EnumKind:
		return "0"
	case protoreflect.MessageKind, protoreflect.GroupKind:
		return "nullptr"
	default:
		return ""
	}
}

// IsScalarType returns true if the kind is a scalar type (not message/enum).
func IsScalarType(kind protoreflect.Kind) bool {
	switch kind {
	case protoreflect.BoolKind,
		protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
		protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind,
		protoreflect.Uint32Kind, protoreflect.Fixed32Kind,
		protoreflect.Uint64Kind, protoreflect.Fixed64Kind,
		protoreflect.FloatKind, protoreflect.DoubleKind,
		protoreflect.StringKind, protoreflect.BytesKind:
		return true
	default:
		return false
	}
}

// IsNumericType returns true if the kind is a numeric type.
func IsNumericType(kind protoreflect.Kind) bool {
	switch kind {
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
		protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind,
		protoreflect.Uint32Kind, protoreflect.Fixed32Kind,
		protoreflect.Uint64Kind, protoreflect.Fixed64Kind,
		protoreflect.FloatKind, protoreflect.DoubleKind:
		return true
	default:
		return false
	}
}

// IsStringType returns true if the kind is a string or bytes type.
func IsStringType(kind protoreflect.Kind) bool {
	return kind == protoreflect.StringKind || kind == protoreflect.BytesKind
}

// WireType returns the wire type for a protobuf kind.
type WireType int

const (
	WireVarint  WireType = 0
	WireFixed64 WireType = 1
	WireBytes   WireType = 2
	WireFixed32 WireType = 5
)

// GetWireType returns the wire type for a protobuf kind.
func GetWireType(kind protoreflect.Kind) WireType {
	switch kind {
	case protoreflect.BoolKind,
		protoreflect.Int32Kind, protoreflect.Int64Kind,
		protoreflect.Uint32Kind, protoreflect.Uint64Kind,
		protoreflect.Sint32Kind, protoreflect.Sint64Kind,
		protoreflect.EnumKind:
		return WireVarint
	case protoreflect.Fixed64Kind, protoreflect.Sfixed64Kind, protoreflect.DoubleKind:
		return WireFixed64
	case protoreflect.Fixed32Kind, protoreflect.Sfixed32Kind, protoreflect.FloatKind:
		return WireFixed32
	case protoreflect.StringKind, protoreflect.BytesKind, protoreflect.MessageKind, protoreflect.GroupKind:
		return WireBytes
	default:
		return WireVarint
	}
}

// RepeatedFieldType returns the C++ container type for a repeated field.
func RepeatedFieldType(kind protoreflect.Kind, runtimePrefix string) string {
	switch kind {
	case protoreflect.MessageKind, protoreflect.GroupKind:
		return runtimePrefix + "::google::protobuf::RepeatedPtrField"
	default:
		return runtimePrefix + "::google::protobuf::RepeatedField"
	}
}

// MapFieldType returns the C++ type for a map field.
func MapFieldType(keyType, valueType, runtimePrefix string) string {
	return runtimePrefix + "::google::protobuf::Map<" + keyType + ", " + valueType + ">"
}
