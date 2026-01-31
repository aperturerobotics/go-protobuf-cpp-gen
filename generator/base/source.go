package base

import (
	"github.com/aperturerobotics/protobuf-cpp-gen/generator/genfile"
	"github.com/aperturerobotics/protobuf-cpp-gen/internal/naming"
	"github.com/aperturerobotics/protobuf-cpp-gen/internal/types"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// GenerateMessageSource generates the implementation (.pb.cc) for a message.
func GenerateMessageSource(g *genfile.GeneratedFile, msg protoreflect.MessageDescriptor, config *genfile.Config) {
	className := naming.ClassName(string(msg.Name()))

	// Constructor
	generateConstructor(g, msg, className)

	// Destructor
	generateDestructor(g, msg, className)

	// Copy constructor
	generateCopyConstructor(g, msg, className)

	// Move constructor
	generateMoveConstructor(g, msg, className)

	// Assignment operators
	generateAssignmentOperators(g, msg, className)

	// Clear method
	generateClearMethod(g, msg, className)

	// ByteSizeLong
	generateByteSizeLong(g, msg, className)

	// Serialize methods
	generateSerializeMethods(g, msg, className)

	// Parse methods
	generateParseMethods(g, msg, className)
}

func generateConstructor(g *genfile.GeneratedFile, msg protoreflect.MessageDescriptor, className string) {
	g.P(className, "::", className, "()")
	g.In()

	// Initialize fields with default values
	fields := msg.Fields()
	first := true
	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		if field.IsMap() || field.IsList() || field.Message() != nil {
			continue // Complex types don't need explicit initialization
		}

		prefix := ":"
		if !first {
			prefix = ","
		}
		first = false

		fieldName := naming.FieldName(string(field.Name()))

		// Enums need static_cast, scalars use direct value
		if field.Enum() != nil {
			enumType := FullyQualifiedEnumName(field.Enum())
			g.P(prefix, " ", fieldName, "(static_cast<", enumType, ">(0))")
		} else {
			defaultVal := types.DefaultValue(field.Kind())
			g.P(prefix, " ", fieldName, "(", defaultVal, ")")
		}
	}

	g.Out()
	g.P("{}")
	g.P()
}

func generateDestructor(g *genfile.GeneratedFile, msg protoreflect.MessageDescriptor, className string) {
	g.P(className, "::~", className, "() {}")
	g.P()
}

func generateCopyConstructor(g *genfile.GeneratedFile, msg protoreflect.MessageDescriptor, className string) {
	g.P(className, "::", className, "(const ", className, "& other) {")
	g.In()
	g.P("*this = other;")
	g.Out()
	g.P("}")
	g.P()
}

func generateMoveConstructor(g *genfile.GeneratedFile, msg protoreflect.MessageDescriptor, className string) {
	g.P(className, "::", className, "(", className, "&& other) noexcept {")
	g.In()
	g.P("*this = std::move(other);")
	g.Out()
	g.P("}")
	g.P()
}

func generateAssignmentOperators(g *genfile.GeneratedFile, msg protoreflect.MessageDescriptor, className string) {
	// Copy assignment
	g.P(className, "& ", className, "::operator=(const ", className, "& other) {")
	g.In()
	g.P("if (this != &other) {")
	g.In()

	fields := msg.Fields()
	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		fieldName := naming.FieldName(string(field.Name()))
		g.P(fieldName, " = other.", fieldName, ";")
	}

	g.Out()
	g.P("}")
	g.P("return *this;")
	g.Out()
	g.P("}")
	g.P()

	// Move assignment
	g.P(className, "& ", className, "::operator=(", className, "&& other) noexcept {")
	g.In()
	g.P("if (this != &other) {")
	g.In()

	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		fieldName := naming.FieldName(string(field.Name()))
		g.P(fieldName, " = std::move(other.", fieldName, ");")
	}

	g.Out()
	g.P("}")
	g.P("return *this;")
	g.Out()
	g.P("}")
	g.P()
}

func generateClearMethod(g *genfile.GeneratedFile, msg protoreflect.MessageDescriptor, className string) {
	g.P("void ", className, "::Clear() {")
	g.In()

	fields := msg.Fields()
	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		fieldName := naming.FieldName(string(field.Name()))

		if field.IsList() {
			g.P(fieldName, ".clear();") // std::vector::clear()
		} else if field.IsMap() {
			g.P(fieldName, ".clear();") // std::map::clear()
		} else if field.Message() != nil {
			g.P(fieldName, ".Clear();") // Our Clear() method
		} else if field.Kind() == protoreflect.StringKind || field.Kind() == protoreflect.BytesKind {
			g.P(fieldName, ".clear();")
		} else if field.Enum() != nil {
			enumType := FullyQualifiedEnumName(field.Enum())
			g.P(fieldName, " = static_cast<", enumType, ">(0);")
		} else {
			defaultVal := types.DefaultValue(field.Kind())
			g.P(fieldName, " = ", defaultVal, ";")
		}
	}

	g.Out()
	g.P("}")
	g.P()
}

func generateByteSizeLong(g *genfile.GeneratedFile, msg protoreflect.MessageDescriptor, className string) {
	g.P("size_t ", className, "::ByteSizeLong() const {")
	g.In()
	g.P("size_t total_size = 0;")
	g.P()

	fields := msg.Fields()
	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		generateFieldSizeCalculation(g, field)
	}

	g.P()
	g.P("return total_size;")
	g.Out()
	g.P("}")
	g.P()
}

func generateFieldSizeCalculation(g *genfile.GeneratedFile, field protoreflect.FieldDescriptor) {
	fieldName := naming.FieldName(string(field.Name()))
	tagSize := tagSize(int(field.Number()))

	if field.IsList() {
		g.P("// ", field.Name(), " (repeated)")
		g.P("for (size_t i = 0; i < ", fieldName, ".size(); i++) {")
		g.In()
		generateSingleFieldSize(g, field, fieldName+"[i]", tagSize)
		g.Out()
		g.P("}")
	} else if field.IsMap() {
		g.P("// ", field.Name(), " (map)")
		g.P("for (const auto& entry : ", fieldName, ") {")
		g.In()
		g.P("total_size += ", tagSize, ";  // outer tag")
		g.P("size_t entry_size = 0;")
		// Calculate key size (field 1)
		generateMapKeySize(g, field.MapKey(), "entry.first")
		// Calculate value size (field 2)
		generateMapValueSize(g, field.MapValue(), "entry.second")
		g.P("total_size += ::google::protobuf::io::CodedOutputStream::VarintSize32(static_cast<uint32_t>(entry_size));")
		g.P("total_size += entry_size;")
		g.Out()
		g.P("}")
	} else {
		g.P("// ", field.Name())
		generateSingleFieldSize(g, field, fieldName, tagSize)
	}
}

func generateSingleFieldSize(g *genfile.GeneratedFile, field protoreflect.FieldDescriptor, accessor string, tagSize int) {
	switch field.Kind() {
	case protoreflect.BoolKind:
		g.P("total_size += ", tagSize, " + 1;  // tag + bool")
	case protoreflect.Int32Kind:
		// Negative int32 values are encoded as 10-byte varints (sign-extended to 64-bit)
		g.P("total_size += ", tagSize, " + (", accessor, " < 0 ? 10 : ::google::protobuf::io::CodedOutputStream::VarintSize32(static_cast<uint32_t>(", accessor, ")));")
	case protoreflect.Sint32Kind:
		// ZigZag encoding ensures all values are positive, so VarintSize32 is correct
		g.P("total_size += ", tagSize, " + ::google::protobuf::io::CodedOutputStream::VarintSize32(::google::protobuf::internal::WireFormatLite::ZigZagEncode32(", accessor, "));")
	case protoreflect.Uint32Kind:
		g.P("total_size += ", tagSize, " + ::google::protobuf::io::CodedOutputStream::VarintSize32(", accessor, ");")
	case protoreflect.Int64Kind:
		// All int64 values use VarintSize64
		g.P("total_size += ", tagSize, " + ::google::protobuf::io::CodedOutputStream::VarintSize64(static_cast<uint64_t>(", accessor, "));")
	case protoreflect.Sint64Kind:
		// ZigZag encoding ensures all values are positive
		g.P("total_size += ", tagSize, " + ::google::protobuf::io::CodedOutputStream::VarintSize64(::google::protobuf::internal::WireFormatLite::ZigZagEncode64(", accessor, "));")
	case protoreflect.Uint64Kind:
		g.P("total_size += ", tagSize, " + ::google::protobuf::io::CodedOutputStream::VarintSize64(", accessor, ");")
	case protoreflect.Fixed32Kind, protoreflect.Sfixed32Kind, protoreflect.FloatKind:
		g.P("total_size += ", tagSize, " + 4;  // tag + fixed32")
	case protoreflect.Fixed64Kind, protoreflect.Sfixed64Kind, protoreflect.DoubleKind:
		g.P("total_size += ", tagSize, " + 8;  // tag + fixed64")
	case protoreflect.StringKind, protoreflect.BytesKind:
		g.P("total_size += ", tagSize, " + ::google::protobuf::io::CodedOutputStream::VarintSize32(static_cast<uint32_t>(", accessor, ".size())) + ", accessor, ".size();")
	case protoreflect.MessageKind:
		g.P("{")
		g.In()
		g.P("size_t msg_size = ", accessor, ".ByteSizeLong();")
		g.P("total_size += ", tagSize, " + ::google::protobuf::io::CodedOutputStream::VarintSize32(static_cast<uint32_t>(msg_size)) + msg_size;")
		g.Out()
		g.P("}")
	case protoreflect.EnumKind:
		g.P("total_size += ", tagSize, " + ::google::protobuf::io::CodedOutputStream::VarintSize32(static_cast<uint32_t>(", accessor, "));")
	}
}

func generateMapKeySize(g *genfile.GeneratedFile, field protoreflect.FieldDescriptor, accessor string) {
	// Map keys are field 1 in the entry
	switch field.Kind() {
	case protoreflect.BoolKind:
		g.P("entry_size += 1 + 1;  // key tag + bool")
	case protoreflect.Int32Kind:
		g.P("entry_size += 1 + (", accessor, " < 0 ? 10 : ::google::protobuf::io::CodedOutputStream::VarintSize32(static_cast<uint32_t>(", accessor, ")));")
	case protoreflect.Int64Kind:
		g.P("entry_size += 1 + ::google::protobuf::io::CodedOutputStream::VarintSize64(static_cast<uint64_t>(", accessor, "));")
	case protoreflect.Uint32Kind:
		g.P("entry_size += 1 + ::google::protobuf::io::CodedOutputStream::VarintSize32(", accessor, ");")
	case protoreflect.Uint64Kind:
		g.P("entry_size += 1 + ::google::protobuf::io::CodedOutputStream::VarintSize64(", accessor, ");")
	case protoreflect.Sint32Kind:
		g.P("entry_size += 1 + ::google::protobuf::io::CodedOutputStream::VarintSize32(::google::protobuf::internal::WireFormatLite::ZigZagEncode32(", accessor, "));")
	case protoreflect.Sint64Kind:
		g.P("entry_size += 1 + ::google::protobuf::io::CodedOutputStream::VarintSize64(::google::protobuf::internal::WireFormatLite::ZigZagEncode64(", accessor, "));")
	case protoreflect.Fixed32Kind, protoreflect.Sfixed32Kind:
		g.P("entry_size += 1 + 4;")
	case protoreflect.Fixed64Kind, protoreflect.Sfixed64Kind:
		g.P("entry_size += 1 + 8;")
	case protoreflect.StringKind:
		g.P("entry_size += 1 + ::google::protobuf::io::CodedOutputStream::VarintSize32(static_cast<uint32_t>(", accessor, ".size())) + ", accessor, ".size();")
	}
}

func generateMapValueSize(g *genfile.GeneratedFile, field protoreflect.FieldDescriptor, accessor string) {
	// Map values are field 2 in the entry
	switch field.Kind() {
	case protoreflect.BoolKind:
		g.P("entry_size += 1 + 1;  // value tag + bool")
	case protoreflect.Int32Kind:
		g.P("entry_size += 1 + (", accessor, " < 0 ? 10 : ::google::protobuf::io::CodedOutputStream::VarintSize32(static_cast<uint32_t>(", accessor, ")));")
	case protoreflect.Int64Kind:
		g.P("entry_size += 1 + ::google::protobuf::io::CodedOutputStream::VarintSize64(static_cast<uint64_t>(", accessor, "));")
	case protoreflect.Uint32Kind:
		g.P("entry_size += 1 + ::google::protobuf::io::CodedOutputStream::VarintSize32(", accessor, ");")
	case protoreflect.Uint64Kind:
		g.P("entry_size += 1 + ::google::protobuf::io::CodedOutputStream::VarintSize64(", accessor, ");")
	case protoreflect.Sint32Kind:
		g.P("entry_size += 1 + ::google::protobuf::io::CodedOutputStream::VarintSize32(::google::protobuf::internal::WireFormatLite::ZigZagEncode32(", accessor, "));")
	case protoreflect.Sint64Kind:
		g.P("entry_size += 1 + ::google::protobuf::io::CodedOutputStream::VarintSize64(::google::protobuf::internal::WireFormatLite::ZigZagEncode64(", accessor, "));")
	case protoreflect.Fixed32Kind, protoreflect.Sfixed32Kind, protoreflect.FloatKind:
		g.P("entry_size += 1 + 4;")
	case protoreflect.Fixed64Kind, protoreflect.Sfixed64Kind, protoreflect.DoubleKind:
		g.P("entry_size += 1 + 8;")
	case protoreflect.StringKind, protoreflect.BytesKind:
		g.P("entry_size += 1 + ::google::protobuf::io::CodedOutputStream::VarintSize32(static_cast<uint32_t>(", accessor, ".size())) + ", accessor, ".size();")
	case protoreflect.MessageKind:
		g.P("{")
		g.In()
		g.P("size_t msg_size = ", accessor, ".ByteSizeLong();")
		g.P("entry_size += 1 + ::google::protobuf::io::CodedOutputStream::VarintSize32(static_cast<uint32_t>(msg_size)) + msg_size;")
		g.Out()
		g.P("}")
	case protoreflect.EnumKind:
		g.P("entry_size += 1 + ::google::protobuf::io::CodedOutputStream::VarintSize32(static_cast<uint32_t>(", accessor, "));")
	}
}

func generateMapKeyWrite(g *genfile.GeneratedFile, field protoreflect.FieldDescriptor, accessor string) {
	// Map keys are field 1 in the entry
	switch field.Kind() {
	case protoreflect.BoolKind:
		g.Pf("target = ::google::protobuf::internal::WireFormatLite::WriteBoolToArray(1, %s, target);", accessor)
	case protoreflect.Int32Kind:
		g.Pf("target = ::google::protobuf::internal::WireFormatLite::WriteInt32ToArray(1, %s, target);", accessor)
	case protoreflect.Int64Kind:
		g.Pf("target = ::google::protobuf::internal::WireFormatLite::WriteInt64ToArray(1, %s, target);", accessor)
	case protoreflect.Uint32Kind:
		g.Pf("target = ::google::protobuf::internal::WireFormatLite::WriteUInt32ToArray(1, %s, target);", accessor)
	case protoreflect.Uint64Kind:
		g.Pf("target = ::google::protobuf::internal::WireFormatLite::WriteUInt64ToArray(1, %s, target);", accessor)
	case protoreflect.Sint32Kind:
		g.Pf("target = ::google::protobuf::internal::WireFormatLite::WriteSInt32ToArray(1, %s, target);", accessor)
	case protoreflect.Sint64Kind:
		g.Pf("target = ::google::protobuf::internal::WireFormatLite::WriteSInt64ToArray(1, %s, target);", accessor)
	case protoreflect.Fixed32Kind:
		g.Pf("target = ::google::protobuf::internal::WireFormatLite::WriteFixed32ToArray(1, %s, target);", accessor)
	case protoreflect.Fixed64Kind:
		g.Pf("target = ::google::protobuf::internal::WireFormatLite::WriteFixed64ToArray(1, %s, target);", accessor)
	case protoreflect.Sfixed32Kind:
		g.Pf("target = ::google::protobuf::internal::WireFormatLite::WriteSFixed32ToArray(1, %s, target);", accessor)
	case protoreflect.Sfixed64Kind:
		g.Pf("target = ::google::protobuf::internal::WireFormatLite::WriteSFixed64ToArray(1, %s, target);", accessor)
	case protoreflect.StringKind:
		g.Pf("target = ::google::protobuf::internal::WireFormatLite::WriteStringToArray(1, %s, target);", accessor)
	}
}

func generateMapValueWrite(g *genfile.GeneratedFile, field protoreflect.FieldDescriptor, accessor string) {
	// Map values are field 2 in the entry
	switch field.Kind() {
	case protoreflect.BoolKind:
		g.Pf("target = ::google::protobuf::internal::WireFormatLite::WriteBoolToArray(2, %s, target);", accessor)
	case protoreflect.Int32Kind:
		g.Pf("target = ::google::protobuf::internal::WireFormatLite::WriteInt32ToArray(2, %s, target);", accessor)
	case protoreflect.Int64Kind:
		g.Pf("target = ::google::protobuf::internal::WireFormatLite::WriteInt64ToArray(2, %s, target);", accessor)
	case protoreflect.Uint32Kind:
		g.Pf("target = ::google::protobuf::internal::WireFormatLite::WriteUInt32ToArray(2, %s, target);", accessor)
	case protoreflect.Uint64Kind:
		g.Pf("target = ::google::protobuf::internal::WireFormatLite::WriteUInt64ToArray(2, %s, target);", accessor)
	case protoreflect.Sint32Kind:
		g.Pf("target = ::google::protobuf::internal::WireFormatLite::WriteSInt32ToArray(2, %s, target);", accessor)
	case protoreflect.Sint64Kind:
		g.Pf("target = ::google::protobuf::internal::WireFormatLite::WriteSInt64ToArray(2, %s, target);", accessor)
	case protoreflect.Fixed32Kind:
		g.Pf("target = ::google::protobuf::internal::WireFormatLite::WriteFixed32ToArray(2, %s, target);", accessor)
	case protoreflect.Fixed64Kind:
		g.Pf("target = ::google::protobuf::internal::WireFormatLite::WriteFixed64ToArray(2, %s, target);", accessor)
	case protoreflect.Sfixed32Kind:
		g.Pf("target = ::google::protobuf::internal::WireFormatLite::WriteSFixed32ToArray(2, %s, target);", accessor)
	case protoreflect.Sfixed64Kind:
		g.Pf("target = ::google::protobuf::internal::WireFormatLite::WriteSFixed64ToArray(2, %s, target);", accessor)
	case protoreflect.FloatKind:
		g.Pf("target = ::google::protobuf::internal::WireFormatLite::WriteFloatToArray(2, %s, target);", accessor)
	case protoreflect.DoubleKind:
		g.Pf("target = ::google::protobuf::internal::WireFormatLite::WriteDoubleToArray(2, %s, target);", accessor)
	case protoreflect.StringKind:
		g.Pf("target = ::google::protobuf::internal::WireFormatLite::WriteStringToArray(2, %s, target);", accessor)
	case protoreflect.BytesKind:
		g.Pf("target = ::google::protobuf::internal::WireFormatLite::WriteBytesToArray(2, %s, target);", accessor)
	case protoreflect.MessageKind:
		g.P("{")
		g.In()
		g.P("target = ::google::protobuf::internal::WireFormatLite::WriteTagToArray(2, ::google::protobuf::internal::WireFormatLite::WIRETYPE_LENGTH_DELIMITED, target);")
		g.P("size_t msg_size = ", accessor, ".ByteSizeLong();")
		g.P("target = ::google::protobuf::io::CodedOutputStream::WriteVarint32ToArray(static_cast<uint32_t>(msg_size), target);")
		g.P("target = ", accessor, ".SerializeToArray(target, msg_size);")
		g.Out()
		g.P("}")
	case protoreflect.EnumKind:
		g.Pf("target = ::google::protobuf::internal::WireFormatLite::WriteEnumToArray(2, static_cast<int>(%s), target);", accessor)
	}
}

func generateSerializeMethods(g *genfile.GeneratedFile, msg protoreflect.MessageDescriptor, className string) {
	// SerializeToString
	g.P("bool ", className, "::SerializeToString(::std::string* output) const {")
	g.In()
	g.P("size_t size = ByteSizeLong();")
	g.P("output->resize(size);")
	g.P("if (size == 0) return true;")
	g.P("::uint8_t* buffer = reinterpret_cast<::uint8_t*>(&(*output)[0]);")
	g.P("::uint8_t* end = SerializeToArray(buffer, size);")
	g.P("return end != nullptr;")
	g.Out()
	g.P("}")
	g.P()

	// SerializeToArray - internal helper
	g.P("::uint8_t* ", className, "::SerializeToArray(::uint8_t* target, size_t size) const {")
	g.In()

	fields := msg.Fields()
	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		generateFieldSerialization(g, field)
	}

	g.P("return target;")
	g.Out()
	g.P("}")
	g.P()
}

func generateFieldSerialization(g *genfile.GeneratedFile, field protoreflect.FieldDescriptor) {
	fieldName := naming.FieldName(string(field.Name()))
	fieldNum := int(field.Number())

	if field.IsList() {
		g.P("// ", field.Name(), " (repeated)")
		g.P("for (size_t i = 0; i < ", fieldName, ".size(); i++) {")
		g.In()
		generateSingleFieldWrite(g, field, fieldName+"[i]", fieldNum)
		g.Out()
		g.P("}")
	} else if field.IsMap() {
		g.P("// ", field.Name(), " (map)")
		g.P("for (const auto& entry : ", fieldName, ") {")
		g.In()
		// Write outer tag (field number, wire type LENGTH_DELIMITED)
		g.Pf("target = ::google::protobuf::internal::WireFormatLite::WriteTagToArray(%d, ::google::protobuf::internal::WireFormatLite::WIRETYPE_LENGTH_DELIMITED, target);", fieldNum)
		// Calculate entry size
		g.P("size_t entry_size = 0;")
		generateMapKeySize(g, field.MapKey(), "entry.first")
		generateMapValueSize(g, field.MapValue(), "entry.second")
		// Write entry length
		g.P("target = ::google::protobuf::io::CodedOutputStream::WriteVarint32ToArray(static_cast<uint32_t>(entry_size), target);")
		// Write key (field 1)
		generateMapKeyWrite(g, field.MapKey(), "entry.first")
		// Write value (field 2)
		generateMapValueWrite(g, field.MapValue(), "entry.second")
		g.Out()
		g.P("}")
	} else {
		g.P("// ", field.Name())
		generateSingleFieldWrite(g, field, fieldName, fieldNum)
	}
}

func generateSingleFieldWrite(g *genfile.GeneratedFile, field protoreflect.FieldDescriptor, accessor string, fieldNum int) {
	switch field.Kind() {
	case protoreflect.BoolKind:
		g.Pf("target = ::google::protobuf::internal::WireFormatLite::WriteBoolToArray(%d, %s, target);", fieldNum, accessor)
	case protoreflect.Int32Kind:
		g.Pf("target = ::google::protobuf::internal::WireFormatLite::WriteInt32ToArray(%d, %s, target);", fieldNum, accessor)
	case protoreflect.Int64Kind:
		g.Pf("target = ::google::protobuf::internal::WireFormatLite::WriteInt64ToArray(%d, %s, target);", fieldNum, accessor)
	case protoreflect.Uint32Kind:
		g.Pf("target = ::google::protobuf::internal::WireFormatLite::WriteUInt32ToArray(%d, %s, target);", fieldNum, accessor)
	case protoreflect.Uint64Kind:
		g.Pf("target = ::google::protobuf::internal::WireFormatLite::WriteUInt64ToArray(%d, %s, target);", fieldNum, accessor)
	case protoreflect.Sint32Kind:
		g.Pf("target = ::google::protobuf::internal::WireFormatLite::WriteSInt32ToArray(%d, %s, target);", fieldNum, accessor)
	case protoreflect.Sint64Kind:
		g.Pf("target = ::google::protobuf::internal::WireFormatLite::WriteSInt64ToArray(%d, %s, target);", fieldNum, accessor)
	case protoreflect.Fixed32Kind:
		g.Pf("target = ::google::protobuf::internal::WireFormatLite::WriteFixed32ToArray(%d, %s, target);", fieldNum, accessor)
	case protoreflect.Fixed64Kind:
		g.Pf("target = ::google::protobuf::internal::WireFormatLite::WriteFixed64ToArray(%d, %s, target);", fieldNum, accessor)
	case protoreflect.Sfixed32Kind:
		g.Pf("target = ::google::protobuf::internal::WireFormatLite::WriteSFixed32ToArray(%d, %s, target);", fieldNum, accessor)
	case protoreflect.Sfixed64Kind:
		g.Pf("target = ::google::protobuf::internal::WireFormatLite::WriteSFixed64ToArray(%d, %s, target);", fieldNum, accessor)
	case protoreflect.FloatKind:
		g.Pf("target = ::google::protobuf::internal::WireFormatLite::WriteFloatToArray(%d, %s, target);", fieldNum, accessor)
	case protoreflect.DoubleKind:
		g.Pf("target = ::google::protobuf::internal::WireFormatLite::WriteDoubleToArray(%d, %s, target);", fieldNum, accessor)
	case protoreflect.StringKind:
		g.Pf("target = ::google::protobuf::internal::WireFormatLite::WriteStringToArray(%d, %s, target);", fieldNum, accessor)
	case protoreflect.BytesKind:
		g.Pf("target = ::google::protobuf::internal::WireFormatLite::WriteBytesToArray(%d, %s, target);", fieldNum, accessor)
	case protoreflect.MessageKind:
		// For nested messages: write tag, length, then serialize the message
		g.P("{")
		g.In()
		g.Pf("target = ::google::protobuf::internal::WireFormatLite::WriteTagToArray(%d, ::google::protobuf::internal::WireFormatLite::WIRETYPE_LENGTH_DELIMITED, target);", fieldNum)
		g.P("size_t msg_size = ", accessor, ".ByteSizeLong();")
		g.P("target = ::google::protobuf::io::CodedOutputStream::WriteVarint32ToArray(static_cast<uint32_t>(msg_size), target);")
		g.P("target = ", accessor, ".SerializeToArray(target, msg_size);")
		g.Out()
		g.P("}")
	case protoreflect.EnumKind:
		g.Pf("target = ::google::protobuf::internal::WireFormatLite::WriteEnumToArray(%d, static_cast<int>(%s), target);", fieldNum, accessor)
	}
}

func generateFieldSerializationCoded(g *genfile.GeneratedFile, field protoreflect.FieldDescriptor) {
	fieldName := naming.FieldName(string(field.Name()))
	fieldNum := int(field.Number())

	if field.IsList() {
		g.P("// ", field.Name(), " (repeated)")
		g.P("for (size_t i = 0; i < ", fieldName, ".size(); i++) {")
		g.In()
		generateSingleFieldWriteCoded(g, field, fieldName+"[i]", fieldNum)
		g.Out()
		g.P("}")
	} else if field.IsMap() {
		g.P("// ", field.Name(), " (map) - TODO")
	} else {
		g.P("// ", field.Name())
		generateSingleFieldWriteCoded(g, field, fieldName, fieldNum)
	}
}

func generateSingleFieldWriteCoded(g *genfile.GeneratedFile, field protoreflect.FieldDescriptor, accessor string, fieldNum int) {
	switch field.Kind() {
	case protoreflect.BoolKind:
		g.Pf("::google::protobuf::internal::WireFormatLite::WriteBool(%d, %s, output);", fieldNum, accessor)
	case protoreflect.Int32Kind:
		g.Pf("::google::protobuf::internal::WireFormatLite::WriteInt32(%d, %s, output);", fieldNum, accessor)
	case protoreflect.Int64Kind:
		g.Pf("::google::protobuf::internal::WireFormatLite::WriteInt64(%d, %s, output);", fieldNum, accessor)
	case protoreflect.Uint32Kind:
		g.Pf("::google::protobuf::internal::WireFormatLite::WriteUInt32(%d, %s, output);", fieldNum, accessor)
	case protoreflect.Uint64Kind:
		g.Pf("::google::protobuf::internal::WireFormatLite::WriteUInt64(%d, %s, output);", fieldNum, accessor)
	case protoreflect.Sint32Kind:
		g.Pf("::google::protobuf::internal::WireFormatLite::WriteSInt32(%d, %s, output);", fieldNum, accessor)
	case protoreflect.Sint64Kind:
		g.Pf("::google::protobuf::internal::WireFormatLite::WriteSInt64(%d, %s, output);", fieldNum, accessor)
	case protoreflect.Fixed32Kind:
		g.Pf("::google::protobuf::internal::WireFormatLite::WriteFixed32(%d, %s, output);", fieldNum, accessor)
	case protoreflect.Fixed64Kind:
		g.Pf("::google::protobuf::internal::WireFormatLite::WriteFixed64(%d, %s, output);", fieldNum, accessor)
	case protoreflect.Sfixed32Kind:
		g.Pf("::google::protobuf::internal::WireFormatLite::WriteSFixed32(%d, %s, output);", fieldNum, accessor)
	case protoreflect.Sfixed64Kind:
		g.Pf("::google::protobuf::internal::WireFormatLite::WriteSFixed64(%d, %s, output);", fieldNum, accessor)
	case protoreflect.FloatKind:
		g.Pf("::google::protobuf::internal::WireFormatLite::WriteFloat(%d, %s, output);", fieldNum, accessor)
	case protoreflect.DoubleKind:
		g.Pf("::google::protobuf::internal::WireFormatLite::WriteDouble(%d, %s, output);", fieldNum, accessor)
	case protoreflect.StringKind:
		g.Pf("::google::protobuf::internal::WireFormatLite::WriteString(%d, %s, output);", fieldNum, accessor)
	case protoreflect.BytesKind:
		g.Pf("::google::protobuf::internal::WireFormatLite::WriteBytes(%d, %s, output);", fieldNum, accessor)
	case protoreflect.MessageKind:
		// For nested messages: serialize to string and write as bytes
		g.P("{")
		g.In()
		g.P("::std::string msg_data;")
		g.P(accessor, ".SerializeToString(&msg_data);")
		g.Pf("::google::protobuf::internal::WireFormatLite::WriteBytes(%d, msg_data, output);", fieldNum)
		g.Out()
		g.P("}")
	case protoreflect.EnumKind:
		g.Pf("::google::protobuf::internal::WireFormatLite::WriteEnum(%d, static_cast<int>(%s), output);", fieldNum, accessor)
	}
}

func generateParseMethods(g *genfile.GeneratedFile, msg protoreflect.MessageDescriptor, className string) {
	// ParseFromString
	g.P("bool ", className, "::ParseFromString(const ::std::string& data) {")
	g.In()
	g.P("return ParseFromArray(reinterpret_cast<const ::uint8_t*>(data.data()), data.size());")
	g.Out()
	g.P("}")
	g.P()

	// ParseFromArray
	g.P("bool ", className, "::ParseFromArray(const ::uint8_t* data, size_t size) {")
	g.In()
	g.P("::google::protobuf::io::CodedInputStream stream(data, static_cast<int>(size));")
	g.P("::google::protobuf::io::CodedInputStream* input = &stream;")
	g.P("::uint32_t tag;")
	g.P("while ((tag = input->ReadTag()) != 0) {")
	g.In()
	g.P("switch (::google::protobuf::internal::WireFormatLite::GetTagFieldNumber(tag)) {")

	fields := msg.Fields()
	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		generateFieldParsing(g, field)
	}

	g.P("default:")
	g.In()
	g.P("if (!::google::protobuf::internal::WireFormatLite::SkipField(input, tag)) return false;")
	g.P("break;")
	g.Out()
	g.P("}")
	g.Out()
	g.P("}")
	g.P("return true;")
	g.Out()
	g.P("}")
	g.P()
}

func generateFieldParsing(g *genfile.GeneratedFile, field protoreflect.FieldDescriptor) {
	fieldName := naming.FieldName(string(field.Name()))
	fieldNum := int(field.Number())

	g.Pf("case %d: {", fieldNum)
	g.In()

	if field.IsList() {
		generateRepeatedFieldParsing(g, field, fieldName)
	} else if field.IsMap() {
		generateMapFieldParsing(g, field, fieldName)
	} else {
		generateSingleFieldParsing(g, field, fieldName)
	}

	g.P("break;")
	g.Out()
	g.P("}")
}

func generateSingleFieldParsing(g *genfile.GeneratedFile, field protoreflect.FieldDescriptor, fieldName string) {
	switch field.Kind() {
	case protoreflect.BoolKind:
		g.P("::uint64_t value;")
		g.P("if (!input->ReadVarint64(&value)) return false;")
		g.P(fieldName, " = value != 0;")
	case protoreflect.Int32Kind:
		// Use ReadVarint64 because negative int32 values are encoded as 10-byte varints
		g.P("::uint64_t value;")
		g.P("if (!input->ReadVarint64(&value)) return false;")
		g.P(fieldName, " = static_cast<::int32_t>(value);")
	case protoreflect.Int64Kind:
		g.P("::uint64_t value;")
		g.P("if (!input->ReadVarint64(&value)) return false;")
		g.P(fieldName, " = static_cast<::int64_t>(value);")
	case protoreflect.Uint32Kind:
		// Use ReadVarint64 for safety - 5-byte varints can occur
		g.P("::uint64_t value;")
		g.P("if (!input->ReadVarint64(&value)) return false;")
		g.P(fieldName, " = static_cast<::uint32_t>(value);")
	case protoreflect.Uint64Kind:
		g.P("::uint64_t value;")
		g.P("if (!input->ReadVarint64(&value)) return false;")
		g.P(fieldName, " = value;")
	case protoreflect.Sint32Kind:
		g.P("::uint64_t value;")
		g.P("if (!input->ReadVarint64(&value)) return false;")
		g.P(fieldName, " = ::google::protobuf::internal::WireFormatLite::ZigZagDecode32(static_cast<::uint32_t>(value));")
	case protoreflect.Sint64Kind:
		g.P("::uint64_t value;")
		g.P("if (!input->ReadVarint64(&value)) return false;")
		g.P(fieldName, " = ::google::protobuf::internal::WireFormatLite::ZigZagDecode64(value);")
	case protoreflect.Fixed32Kind:
		g.P("::uint32_t value;")
		g.P("if (!input->ReadLittleEndian32(&value)) return false;")
		g.P(fieldName, " = value;")
	case protoreflect.Fixed64Kind:
		g.P("::uint64_t value;")
		g.P("if (!input->ReadLittleEndian64(&value)) return false;")
		g.P(fieldName, " = value;")
	case protoreflect.Sfixed32Kind:
		g.P("::uint32_t value;")
		g.P("if (!input->ReadLittleEndian32(&value)) return false;")
		g.P(fieldName, " = static_cast<::int32_t>(value);")
	case protoreflect.Sfixed64Kind:
		g.P("::uint64_t value;")
		g.P("if (!input->ReadLittleEndian64(&value)) return false;")
		g.P(fieldName, " = static_cast<::int64_t>(value);")
	case protoreflect.FloatKind:
		g.P("::uint32_t value;")
		g.P("if (!input->ReadLittleEndian32(&value)) return false;")
		g.P("::memcpy(&", fieldName, ", &value, sizeof(float));")
	case protoreflect.DoubleKind:
		g.P("::uint64_t value;")
		g.P("if (!input->ReadLittleEndian64(&value)) return false;")
		g.P("::memcpy(&", fieldName, ", &value, sizeof(double));")
	case protoreflect.StringKind, protoreflect.BytesKind:
		g.P("::uint32_t length;")
		g.P("if (!input->ReadVarint32(&length)) return false;")
		g.P("if (!input->ReadString(&", fieldName, ", static_cast<int>(length))) return false;")
	case protoreflect.MessageKind:
		g.P("::uint32_t length;")
		g.P("if (!input->ReadVarint32(&length)) return false;")
		g.P("auto limit = input->PushLimit(static_cast<int>(length));")
		g.P("// Read the sub-message bytes and parse")
		g.P("const void* data_ptr;")
		g.P("int data_size;")
		g.P("if (!input->GetDirectBufferPointer(&data_ptr, &data_size)) return false;")
		g.P("if (static_cast<size_t>(data_size) < length) return false;")
		g.P("if (!", fieldName, ".ParseFromArray(reinterpret_cast<const ::uint8_t*>(data_ptr), length)) return false;")
		g.P("input->Skip(static_cast<int>(length));")
		g.P("input->PopLimit(limit);")
	case protoreflect.EnumKind:
		g.P("::uint32_t value;")
		g.P("if (!input->ReadVarint32(&value)) return false;")
		enumType := cppEnumType(field)
		g.P(fieldName, " = static_cast<", enumType, ">(value);")
	}
}

func generateRepeatedFieldParsing(g *genfile.GeneratedFile, field protoreflect.FieldDescriptor, fieldName string) {
	switch field.Kind() {
	case protoreflect.BoolKind:
		g.P("::uint64_t value;")
		g.P("if (!input->ReadVarint64(&value)) return false;")
		g.P(fieldName, ".push_back(value != 0);")
	case protoreflect.Int32Kind:
		g.P("::uint64_t value;")
		g.P("if (!input->ReadVarint64(&value)) return false;")
		g.P(fieldName, ".push_back(static_cast<::int32_t>(value));")
	case protoreflect.Int64Kind:
		g.P("::uint64_t value;")
		g.P("if (!input->ReadVarint64(&value)) return false;")
		g.P(fieldName, ".push_back(static_cast<::int64_t>(value));")
	case protoreflect.Uint32Kind:
		g.P("::uint64_t value;")
		g.P("if (!input->ReadVarint64(&value)) return false;")
		g.P(fieldName, ".push_back(static_cast<::uint32_t>(value));")
	case protoreflect.Uint64Kind:
		g.P("::uint64_t value;")
		g.P("if (!input->ReadVarint64(&value)) return false;")
		g.P(fieldName, ".push_back(value);")
	case protoreflect.Sint32Kind:
		g.P("::uint64_t value;")
		g.P("if (!input->ReadVarint64(&value)) return false;")
		g.P(fieldName, ".push_back(::google::protobuf::internal::WireFormatLite::ZigZagDecode32(static_cast<::uint32_t>(value)));")
	case protoreflect.Sint64Kind:
		g.P("::uint64_t value;")
		g.P("if (!input->ReadVarint64(&value)) return false;")
		g.P(fieldName, ".push_back(::google::protobuf::internal::WireFormatLite::ZigZagDecode64(value));")
	case protoreflect.Fixed32Kind:
		g.P("::uint32_t value;")
		g.P("if (!input->ReadLittleEndian32(&value)) return false;")
		g.P(fieldName, ".push_back(value);")
	case protoreflect.Fixed64Kind:
		g.P("::uint64_t value;")
		g.P("if (!input->ReadLittleEndian64(&value)) return false;")
		g.P(fieldName, ".push_back(value);")
	case protoreflect.Sfixed32Kind:
		g.P("::uint32_t value;")
		g.P("if (!input->ReadLittleEndian32(&value)) return false;")
		g.P(fieldName, ".push_back(static_cast<::int32_t>(value));")
	case protoreflect.Sfixed64Kind:
		g.P("::uint64_t value;")
		g.P("if (!input->ReadLittleEndian64(&value)) return false;")
		g.P(fieldName, ".push_back(static_cast<::int64_t>(value));")
	case protoreflect.FloatKind:
		g.P("::uint32_t raw;")
		g.P("if (!input->ReadLittleEndian32(&raw)) return false;")
		g.P("float value;")
		g.P("::memcpy(&value, &raw, sizeof(float));")
		g.P(fieldName, ".push_back(value);")
	case protoreflect.DoubleKind:
		g.P("::uint64_t raw;")
		g.P("if (!input->ReadLittleEndian64(&raw)) return false;")
		g.P("double value;")
		g.P("::memcpy(&value, &raw, sizeof(double));")
		g.P(fieldName, ".push_back(value);")
	case protoreflect.StringKind, protoreflect.BytesKind:
		// For strings, read length first, then read into new element
		g.P("::uint32_t length;")
		g.P("if (!input->ReadVarint32(&length)) return false;")
		g.P(fieldName, ".emplace_back();")
		g.P("if (!input->ReadString(&", fieldName, ".back(), static_cast<int>(length))) return false;")
	case protoreflect.MessageKind:
		// For messages, use emplace_back and parse into the last element
		g.P(fieldName, ".emplace_back();")
		g.P("::uint32_t length;")
		g.P("if (!input->ReadVarint32(&length)) return false;")
		g.P("auto limit = input->PushLimit(static_cast<int>(length));")
		g.P("const void* data_ptr;")
		g.P("int data_size;")
		g.P("if (!input->GetDirectBufferPointer(&data_ptr, &data_size)) return false;")
		g.P("if (static_cast<size_t>(data_size) < length) return false;")
		g.P("if (!", fieldName, ".back().ParseFromArray(reinterpret_cast<const ::uint8_t*>(data_ptr), length)) return false;")
		g.P("input->Skip(static_cast<int>(length));")
		g.P("input->PopLimit(limit);")
	case protoreflect.EnumKind:
		g.P("::uint32_t value;")
		g.P("if (!input->ReadVarint32(&value)) return false;")
		// Enums use int storage
		g.P(fieldName, ".push_back(static_cast<int>(value));")
	}
}

func generateMapFieldParsing(g *genfile.GeneratedFile, field protoreflect.FieldDescriptor, fieldName string) {
	keyField := field.MapKey()
	valueField := field.MapValue()

	// Declare key and value variables with their C++ types
	keyType := mapKeyCppType(keyField)
	valueType := mapValueCppType(valueField)

	g.P("// Map entry")
	g.P("::uint32_t entry_length;")
	g.P("if (!input->ReadVarint32(&entry_length)) return false;")
	g.P("auto entry_limit = input->PushLimit(static_cast<int>(entry_length));")
	g.P(keyType, " map_key{};")
	g.P(valueType, " map_value{};")
	g.P("while (true) {")
	g.In()
	g.P("::uint32_t entry_tag = input->ReadTag();")
	g.P("if (entry_tag == 0) break;")
	g.P("switch (::google::protobuf::internal::WireFormatLite::GetTagFieldNumber(entry_tag)) {")
	g.P("case 1: {  // key")
	g.In()
	generateMapEntryParsing(g, keyField, "map_key")
	g.P("break;")
	g.Out()
	g.P("}")
	g.P("case 2: {  // value")
	g.In()
	generateMapEntryParsing(g, valueField, "map_value")
	g.P("break;")
	g.Out()
	g.P("}")
	g.P("default:")
	g.In()
	g.P("if (!::google::protobuf::internal::WireFormatLite::SkipField(input, entry_tag)) return false;")
	g.P("break;")
	g.Out()
	g.P("}")
	g.Out()
	g.P("}")
	g.P("input->PopLimit(entry_limit);")
	g.P(fieldName, "[map_key] = std::move(map_value);")
}

func mapKeyCppType(field protoreflect.FieldDescriptor) string {
	switch field.Kind() {
	case protoreflect.BoolKind:
		return "bool"
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return "::int32_t"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return "::int64_t"
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return "::uint32_t"
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return "::uint64_t"
	case protoreflect.StringKind:
		return "::std::string"
	default:
		return "::std::string" // fallback
	}
}

func mapValueCppType(field protoreflect.FieldDescriptor) string {
	switch field.Kind() {
	case protoreflect.BoolKind:
		return "bool"
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return "::int32_t"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return "::int64_t"
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return "::uint32_t"
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return "::uint64_t"
	case protoreflect.FloatKind:
		return "float"
	case protoreflect.DoubleKind:
		return "double"
	case protoreflect.StringKind, protoreflect.BytesKind:
		return "::std::string"
	case protoreflect.MessageKind:
		return FullyQualifiedClassName(field.Message())
	case protoreflect.EnumKind:
		return FullyQualifiedEnumName(field.Enum())
	default:
		return "::std::string" // fallback
	}
}

func generateMapEntryParsing(g *genfile.GeneratedFile, field protoreflect.FieldDescriptor, varName string) {
	switch field.Kind() {
	case protoreflect.BoolKind:
		g.P("::uint64_t value;")
		g.P("if (!input->ReadVarint64(&value)) return false;")
		g.P(varName, " = value != 0;")
	case protoreflect.Int32Kind:
		g.P("::uint64_t value;")
		g.P("if (!input->ReadVarint64(&value)) return false;")
		g.P(varName, " = static_cast<::int32_t>(value);")
	case protoreflect.Int64Kind:
		g.P("::uint64_t value;")
		g.P("if (!input->ReadVarint64(&value)) return false;")
		g.P(varName, " = static_cast<::int64_t>(value);")
	case protoreflect.Uint32Kind:
		g.P("::uint64_t value;")
		g.P("if (!input->ReadVarint64(&value)) return false;")
		g.P(varName, " = static_cast<::uint32_t>(value);")
	case protoreflect.Uint64Kind:
		g.P("::uint64_t value;")
		g.P("if (!input->ReadVarint64(&value)) return false;")
		g.P(varName, " = value;")
	case protoreflect.Sint32Kind:
		g.P("::uint64_t value;")
		g.P("if (!input->ReadVarint64(&value)) return false;")
		g.P(varName, " = ::google::protobuf::internal::WireFormatLite::ZigZagDecode32(static_cast<::uint32_t>(value));")
	case protoreflect.Sint64Kind:
		g.P("::uint64_t value;")
		g.P("if (!input->ReadVarint64(&value)) return false;")
		g.P(varName, " = ::google::protobuf::internal::WireFormatLite::ZigZagDecode64(value);")
	case protoreflect.Fixed32Kind:
		g.P("::uint32_t value;")
		g.P("if (!input->ReadLittleEndian32(&value)) return false;")
		g.P(varName, " = value;")
	case protoreflect.Fixed64Kind:
		g.P("::uint64_t value;")
		g.P("if (!input->ReadLittleEndian64(&value)) return false;")
		g.P(varName, " = value;")
	case protoreflect.Sfixed32Kind:
		g.P("::uint32_t value;")
		g.P("if (!input->ReadLittleEndian32(&value)) return false;")
		g.P(varName, " = static_cast<::int32_t>(value);")
	case protoreflect.Sfixed64Kind:
		g.P("::uint64_t value;")
		g.P("if (!input->ReadLittleEndian64(&value)) return false;")
		g.P(varName, " = static_cast<::int64_t>(value);")
	case protoreflect.FloatKind:
		g.P("::uint32_t raw;")
		g.P("if (!input->ReadLittleEndian32(&raw)) return false;")
		g.P("::memcpy(&", varName, ", &raw, sizeof(float));")
	case protoreflect.DoubleKind:
		g.P("::uint64_t raw;")
		g.P("if (!input->ReadLittleEndian64(&raw)) return false;")
		g.P("::memcpy(&", varName, ", &raw, sizeof(double));")
	case protoreflect.StringKind, protoreflect.BytesKind:
		g.P("::uint32_t str_length;")
		g.P("if (!input->ReadVarint32(&str_length)) return false;")
		g.P("if (!input->ReadString(&", varName, ", static_cast<int>(str_length))) return false;")
	case protoreflect.MessageKind:
		g.P("::uint32_t msg_length;")
		g.P("if (!input->ReadVarint32(&msg_length)) return false;")
		g.P("auto msg_limit = input->PushLimit(static_cast<int>(msg_length));")
		g.P("const void* msg_data_ptr;")
		g.P("int msg_data_size;")
		g.P("if (!input->GetDirectBufferPointer(&msg_data_ptr, &msg_data_size)) return false;")
		g.P("if (static_cast<size_t>(msg_data_size) < msg_length) return false;")
		g.P("if (!", varName, ".ParseFromArray(reinterpret_cast<const ::uint8_t*>(msg_data_ptr), msg_length)) return false;")
		g.P("input->Skip(static_cast<int>(msg_length));")
		g.P("input->PopLimit(msg_limit);")
	case protoreflect.EnumKind:
		g.P("::uint32_t value;")
		g.P("if (!input->ReadVarint32(&value)) return false;")
		enumType := FullyQualifiedEnumName(field.Enum())
		g.P(varName, " = static_cast<", enumType, ">(value);")
	}
}

func tagSize(fieldNumber int) int {
	tag := fieldNumber << 3
	if tag < 128 {
		return 1
	}
	if tag < 16384 {
		return 2
	}
	if tag < 2097152 {
		return 3
	}
	if tag < 268435456 {
		return 4
	}
	return 5
}
