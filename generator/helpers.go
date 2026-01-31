package generator

import (
	"google.golang.org/protobuf/reflect/protoreflect"
)

// KeySize returns the size of a protobuf tag for a field.
func KeySize(fieldNumber protoreflect.FieldNumber, wireType int) int {
	x := uint32(fieldNumber)<<3 | uint32(wireType)
	size := 0
	for size = 0; x > 127; size++ {
		x >>= 7
	}
	size++
	return size
}

// Wire types as defined in the protobuf wire format.
const (
	WireVarint  = 0
	WireFixed64 = 1
	WireBytes   = 2
	WireFixed32 = 5
)

var wireTypes = map[protoreflect.Kind]int{
	protoreflect.BoolKind:     WireVarint,
	protoreflect.EnumKind:     WireVarint,
	protoreflect.Int32Kind:    WireVarint,
	protoreflect.Sint32Kind:   WireVarint,
	protoreflect.Uint32Kind:   WireVarint,
	protoreflect.Int64Kind:    WireVarint,
	protoreflect.Sint64Kind:   WireVarint,
	protoreflect.Uint64Kind:   WireVarint,
	protoreflect.Sfixed32Kind: WireFixed32,
	protoreflect.Fixed32Kind:  WireFixed32,
	protoreflect.FloatKind:    WireFixed32,
	protoreflect.Sfixed64Kind: WireFixed64,
	protoreflect.Fixed64Kind:  WireFixed64,
	protoreflect.DoubleKind:   WireFixed64,
	protoreflect.StringKind:   WireBytes,
	protoreflect.BytesKind:    WireBytes,
	protoreflect.MessageKind:  WireBytes,
	protoreflect.GroupKind:    WireBytes,
}

// ProtoWireType returns the wire type for a protobuf kind.
func ProtoWireType(k protoreflect.Kind) int {
	return wireTypes[k]
}

// IsScalarKind returns true if the kind is a scalar type.
func IsScalarKind(k protoreflect.Kind) bool {
	switch k {
	case protoreflect.BoolKind,
		protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
		protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind,
		protoreflect.Uint32Kind, protoreflect.Fixed32Kind,
		protoreflect.Uint64Kind, protoreflect.Fixed64Kind,
		protoreflect.FloatKind, protoreflect.DoubleKind,
		protoreflect.StringKind, protoreflect.BytesKind:
		return true
	}
	return false
}

// IsVarintKind returns true if the kind is encoded as a varint.
func IsVarintKind(k protoreflect.Kind) bool {
	switch k {
	case protoreflect.BoolKind,
		protoreflect.Int32Kind, protoreflect.Sint32Kind,
		protoreflect.Int64Kind, protoreflect.Sint64Kind,
		protoreflect.Uint32Kind, protoreflect.Uint64Kind,
		protoreflect.EnumKind:
		return true
	}
	return false
}

// IsFixed32Kind returns true if the kind is encoded as fixed32.
func IsFixed32Kind(k protoreflect.Kind) bool {
	switch k {
	case protoreflect.Sfixed32Kind, protoreflect.Fixed32Kind, protoreflect.FloatKind:
		return true
	}
	return false
}

// IsFixed64Kind returns true if the kind is encoded as fixed64.
func IsFixed64Kind(k protoreflect.Kind) bool {
	switch k {
	case protoreflect.Sfixed64Kind, protoreflect.Fixed64Kind, protoreflect.DoubleKind:
		return true
	}
	return false
}

// IsBytesKind returns true if the kind is encoded as bytes.
func IsBytesKind(k protoreflect.Kind) bool {
	switch k {
	case protoreflect.StringKind, protoreflect.BytesKind, protoreflect.MessageKind:
		return true
	}
	return false
}
