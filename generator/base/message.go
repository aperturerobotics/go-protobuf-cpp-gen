package base

import (
	"github.com/aperturerobotics/go-protobuf-cpp-gen/generator/genfile"
	"github.com/aperturerobotics/go-protobuf-cpp-gen/internal/naming"
	"github.com/aperturerobotics/go-protobuf-cpp-gen/internal/types"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// GenerateMessageHeader generates the class declaration for a message.
func GenerateMessageHeader(g *genfile.GeneratedFile, msg protoreflect.MessageDescriptor, config *genfile.Config) {
	className := naming.ClassName(string(msg.Name()))

	baseClass := config.MessageBaseClass()
	if baseClass != "" {
		g.P("class ", className, " : public ", baseClass, " {")
	} else {
		g.P("class ", className, " {")
	}
	g.P(" public:")
	g.In()

	// Default constructor
	g.P(className, "();")

	// Destructor
	if baseClass != "" {
		g.P("~", className, "() override;")
	} else {
		g.P("~", className, "();")
	}

	// Copy constructor and assignment
	g.P(className, "(const ", className, "& other);")
	g.P(className, "(", className, "&& other) noexcept;")
	g.P(className, "& operator=(const ", className, "& other);")
	g.P(className, "& operator=(", className, "&& other) noexcept;")

	g.P()

	// Clear method
	g.P("void Clear();")

	// Serialization methods
	g.P()
	g.P("// Serialization")
	g.P("size_t ByteSizeLong() const;")
	g.P("bool SerializeToString(::std::string* output) const;")
	g.P("::uint8_t* SerializeToArray(::uint8_t* target, size_t size) const;")
	g.P("bool ParseFromString(const ::std::string& data);")
	g.P("bool ParseFromArray(const ::uint8_t* data, size_t size);")

	// Field accessors
	generateFieldAccessors(g, msg)

	g.Out()
	g.P()
	g.P(" private:")
	g.In()

	// Field storage
	generateFieldStorage(g, msg)

	g.Out()
	g.P("};")
	g.P()
}

func generateFieldAccessors(g *genfile.GeneratedFile, msg protoreflect.MessageDescriptor) {
	fields := msg.Fields()
	if fields.Len() == 0 {
		return
	}

	g.P()
	g.P("// Field accessors")

	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		generateSingleFieldAccessors(g, field)
	}
}

func generateSingleFieldAccessors(g *genfile.GeneratedFile, field protoreflect.FieldDescriptor) {
	fieldName := naming.FieldName(string(field.Name()))
	accessorName := naming.AccessorName(string(field.Name()))

	g.P()
	g.P("// ", field.Name())

	if field.IsMap() {
		generateMapAccessors(g, field, fieldName, accessorName)
	} else if field.IsList() {
		generateRepeatedAccessors(g, field, fieldName, accessorName)
	} else if field.Message() != nil {
		generateMessageAccessors(g, field, fieldName, accessorName)
	} else if field.Enum() != nil {
		generateEnumAccessors(g, field, fieldName, accessorName)
	} else {
		generateScalarAccessors(g, field, fieldName, accessorName)
	}
}

func generateScalarAccessors(g *genfile.GeneratedFile, field protoreflect.FieldDescriptor, fieldName, accessorName string) {
	cppType := types.ScalarTypeToCpp(field.Kind())
	typeName := cppType.Name

	// Getter
	if types.IsStringType(field.Kind()) {
		g.P("inline const ", typeName, "& ", accessorName, "() const { return ", fieldName, "; }")
	} else {
		g.P("inline ", typeName, " ", accessorName, "() const { return ", fieldName, "; }")
	}

	// Setter
	setName := naming.SetAccessorName(string(field.Name()))
	if types.IsStringType(field.Kind()) {
		g.P("inline void ", setName, "(const ", typeName, "& value) { ", fieldName, " = value; }")
		g.P("inline void ", setName, "(", typeName, "&& value) { ", fieldName, " = std::move(value); }")
	} else {
		g.P("inline void ", setName, "(", typeName, " value) { ", fieldName, " = value; }")
	}

	// Clear
	clearName := naming.ClearAccessorName(string(field.Name()))
	defaultVal := types.DefaultValue(field.Kind())
	if types.IsStringType(field.Kind()) {
		g.P("inline void ", clearName, "() { ", fieldName, ".clear(); }")
	} else {
		g.P("inline void ", clearName, "() { ", fieldName, " = ", defaultVal, "; }")
	}

	// Mutable pointer for strings/bytes
	if types.IsStringType(field.Kind()) {
		mutableName := naming.MutableAccessorName(string(field.Name()))
		g.P("inline ", typeName, "* ", mutableName, "() { return &", fieldName, "; }")
	}
}

func generateEnumAccessors(g *genfile.GeneratedFile, field protoreflect.FieldDescriptor, fieldName, accessorName string) {
	enumType := cppEnumType(field)

	// Getter
	g.P("inline ", enumType, " ", accessorName, "() const { return ", fieldName, "; }")

	// Setter
	setName := naming.SetAccessorName(string(field.Name()))
	g.P("inline void ", setName, "(", enumType, " value) { ", fieldName, " = value; }")

	// Clear - enums default to 0
	clearName := naming.ClearAccessorName(string(field.Name()))
	g.P("inline void ", clearName, "() { ", fieldName, " = static_cast<", enumType, ">(0); }")
}

func generateMessageAccessors(g *genfile.GeneratedFile, field protoreflect.FieldDescriptor, fieldName, accessorName string) {
	msgType := FullyQualifiedClassName(field.Message())

	// has_* for optional message fields
	if field.HasPresence() {
		hasName := naming.HasAccessorName(string(field.Name()))
		g.P("inline bool ", hasName, "() const { return true; }  // TODO: Track presence")
	}

	// Getter (const reference)
	g.P("inline const ", msgType, "& ", accessorName, "() const { return ", fieldName, "; }")

	// Mutable getter
	mutableName := naming.MutableAccessorName(string(field.Name()))
	g.P("inline ", msgType, "* ", mutableName, "() { return &", fieldName, "; }")

	// Clear
	clearName := naming.ClearAccessorName(string(field.Name()))
	g.P("inline void ", clearName, "() { ", fieldName, ".Clear(); }")
}

func generateRepeatedAccessors(g *genfile.GeneratedFile, field protoreflect.FieldDescriptor, fieldName, accessorName string) {
	// Size
	sizeName := naming.SizeAccessorName(string(field.Name()))
	g.P("inline int ", sizeName, "() const { return ", fieldName, ".size(); }")

	// Determine element type
	var elemType string
	var storageType string // For container template parameter
	isMessage := field.Message() != nil
	isEnum := field.Enum() != nil
	isString := types.IsStringType(field.Kind())

	if isMessage {
		elemType = FullyQualifiedClassName(field.Message())
		storageType = elemType
	} else if isEnum {
		elemType = cppEnumType(field)
		storageType = "int" // Enums use int storage
	} else {
		elemType = types.ScalarTypeToCpp(field.Kind()).Name
		storageType = elemType
	}

	// Element getter by index - use vector-style access for messages, RepeatedField for scalars
	if isMessage {
		// Use std::vector for messages (no MessageLite requirement)
		g.P("inline const ", elemType, "& ", accessorName, "(int index) const { return ", fieldName, "[index]; }")
		g.P("inline ", elemType, "* ", naming.MutableAccessorName(string(field.Name())), "(int index) { return &", fieldName, "[index]; }")
	} else if isString {
		g.P("inline const ", elemType, "& ", accessorName, "(int index) const { return ", fieldName, "[index]; }")
		g.P("inline ", elemType, "* ", naming.MutableAccessorName(string(field.Name())), "(int index) { return &", fieldName, "[index]; }")
	} else if isEnum {
		// Cast from int storage to enum type
		g.P("inline ", elemType, " ", accessorName, "(int index) const { return static_cast<", elemType, ">(", fieldName, "[index]); }")
	} else {
		g.P("inline ", elemType, " ", accessorName, "(int index) const { return ", fieldName, "[index]; }")
	}

	// Add
	addName := naming.AddAccessorName(string(field.Name()))
	if isMessage {
		// For messages, push_back and return pointer to last element
		g.P("inline ", elemType, "* ", addName, "() { ", fieldName, ".emplace_back(); return &", fieldName, ".back(); }")
	} else if isString {
		g.P("inline void ", addName, "(const ", elemType, "& value) { ", fieldName, ".push_back(value); }")
		g.P("inline ", elemType, "* ", addName, "() { ", fieldName, ".emplace_back(); return &", fieldName, ".back(); }")
	} else if isEnum {
		// Cast enum to int for storage
		g.P("inline void ", addName, "(", elemType, " value) { ", fieldName, ".push_back(static_cast<int>(value)); }")
	} else {
		g.P("inline void ", addName, "(", elemType, " value) { ", fieldName, ".push_back(value); }")
	}

	// Clear
	clearName := naming.ClearAccessorName(string(field.Name()))
	g.P("inline void ", clearName, "() { ", fieldName, ".clear(); }")

	// Direct access to repeated field
	if isMessage {
		g.P("inline const ::std::vector<", elemType, ">& ", accessorName, "() const { return ", fieldName, "; }")
		g.P("inline ::std::vector<", elemType, ">* ", naming.MutableAccessorName(string(field.Name())), "() { return &", fieldName, "; }")
	} else if isString {
		g.P("inline const ::std::vector<", elemType, ">& ", accessorName, "() const { return ", fieldName, "; }")
	} else {
		// Scalar or enum - use std::vector with storage type
		g.P("inline const ::std::vector<", storageType, ">& ", accessorName, "() const { return ", fieldName, "; }")
	}
}

func generateMapAccessors(g *genfile.GeneratedFile, field protoreflect.FieldDescriptor, fieldName, accessorName string) {
	keyField := field.MapKey()
	valueField := field.MapValue()

	keyType := types.ScalarTypeToCpp(keyField.Kind()).Name
	var valueType string
	if valueField.Message() != nil {
		valueType = FullyQualifiedClassName(valueField.Message())
	} else if valueField.Enum() != nil {
		valueType = FullyQualifiedEnumName(valueField.Enum())
	} else {
		valueType = types.ScalarTypeToCpp(valueField.Kind()).Name
	}

	// Use std::map for standalone generation
	mapType := "::std::map<" + keyType + ", " + valueType + ">"

	// Size
	sizeName := naming.SizeAccessorName(string(field.Name()))
	g.P("inline int ", sizeName, "() const { return ", fieldName, ".size(); }")

	// Const accessor
	g.P("inline const ", mapType, "& ", accessorName, "() const { return ", fieldName, "; }")

	// Mutable accessor
	mutableName := naming.MutableAccessorName(string(field.Name()))
	g.P("inline ", mapType, "* ", mutableName, "() { return &", fieldName, "; }")

	// Clear
	clearName := naming.ClearAccessorName(string(field.Name()))
	g.P("inline void ", clearName, "() { ", fieldName, ".clear(); }")
}

func generateFieldStorage(g *genfile.GeneratedFile, msg protoreflect.MessageDescriptor) {
	fields := msg.Fields()
	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		generateSingleFieldStorage(g, field)
	}
}

func generateSingleFieldStorage(g *genfile.GeneratedFile, field protoreflect.FieldDescriptor) {
	fieldName := naming.FieldName(string(field.Name()))

	if field.IsMap() {
		keyField := field.MapKey()
		valueField := field.MapValue()

		keyType := types.ScalarTypeToCpp(keyField.Kind()).Name
		var valueType string
		if valueField.Message() != nil {
			valueType = FullyQualifiedClassName(valueField.Message())
		} else if valueField.Enum() != nil {
			valueType = FullyQualifiedEnumName(valueField.Enum())
		} else {
			valueType = types.ScalarTypeToCpp(valueField.Kind()).Name
		}

		// Use std::map for maps (standalone, no protobuf runtime dependency)
		g.P("::std::map<", keyType, ", ", valueType, "> ", fieldName, ";")
	} else if field.IsList() {
		// Use std::vector for all repeated fields (no protobuf runtime dependency)
		if field.Message() != nil {
			msgType := FullyQualifiedClassName(field.Message())
			g.P("::std::vector<", msgType, "> ", fieldName, ";")
		} else if field.Enum() != nil {
			// Use int storage for enums
			g.P("::std::vector<int> ", fieldName, ";")
		} else {
			cppType := types.ScalarTypeToCpp(field.Kind())
			g.P("::std::vector<", cppType.Name, "> ", fieldName, ";")
		}
	} else if field.Message() != nil {
		msgType := FullyQualifiedClassName(field.Message())
		g.P(msgType, " ", fieldName, ";")
	} else if field.Enum() != nil {
		enumType := cppEnumType(field)
		g.P(enumType, " ", fieldName, ";")
	} else {
		cppType := types.ScalarTypeToCpp(field.Kind())
		g.P(cppType.Name, " ", fieldName, ";")
	}
}

// cppEnumType returns the C++ type name for an enum field.
func cppEnumType(field protoreflect.FieldDescriptor) string {
	enum := field.Enum()
	if enum == nil {
		return "int"
	}
	return FullyQualifiedEnumName(enum)
}

// FullyQualifiedEnumName returns the fully qualified C++ enum name.
// Since we generate all enums at namespace level (not inside classes),
// we use just the namespace and enum name without parent class names.
func FullyQualifiedEnumName(enum protoreflect.EnumDescriptor) string {
	enumName := naming.EnumName(string(enum.Name()))

	// Get the file for namespace
	file := enum.ParentFile()
	pkg := string(file.Package())

	if pkg == "" {
		return "::" + enumName
	}
	return "::" + naming.Namespace(pkg) + "::" + enumName
}

func joinParts(parts []string) string {
	result := ""
	for i, p := range parts {
		if i > 0 {
			result += "::"
		}
		result += p
	}
	return result
}
