// Package naming provides C++ naming conventions for protobuf code generation.
package naming

import (
	"regexp"
	"strings"
	"unicode"
)

var (
	// cppKeywords are reserved C++ keywords that need escaping.
	cppKeywords = map[string]bool{
		"alignas": true, "alignof": true, "and": true, "and_eq": true,
		"asm": true, "auto": true, "bitand": true, "bitor": true,
		"bool": true, "break": true, "case": true, "catch": true,
		"char": true, "char8_t": true, "char16_t": true, "char32_t": true,
		"class": true, "compl": true, "concept": true, "const": true,
		"consteval": true, "constexpr": true, "constinit": true, "const_cast": true,
		"continue": true, "co_await": true, "co_return": true, "co_yield": true,
		"decltype": true, "default": true, "delete": true, "do": true,
		"double": true, "dynamic_cast": true, "else": true, "enum": true,
		"explicit": true, "export": true, "extern": true, "false": true,
		"float": true, "for": true, "friend": true, "goto": true,
		"if": true, "inline": true, "int": true, "long": true,
		"mutable": true, "namespace": true, "new": true, "noexcept": true,
		"not": true, "not_eq": true, "nullptr": true, "operator": true,
		"or": true, "or_eq": true, "private": true, "protected": true,
		"public": true, "register": true, "reinterpret_cast": true, "requires": true,
		"return": true, "short": true, "signed": true, "sizeof": true,
		"static": true, "static_assert": true, "static_cast": true, "struct": true,
		"switch": true, "template": true, "this": true, "thread_local": true,
		"throw": true, "true": true, "try": true, "typedef": true,
		"typeid": true, "typename": true, "union": true, "unsigned": true,
		"using": true, "virtual": true, "void": true, "volatile": true,
		"wchar_t": true, "while": true, "xor": true, "xor_eq": true,
	}

	// Regex to match non-alphanumeric characters (except underscore)
	nonAlphanumeric = regexp.MustCompile(`[^a-zA-Z0-9_]`)
)

// CppName converts a proto identifier to a valid C++ identifier.
func CppName(name string) string {
	// Replace invalid characters with underscores
	result := nonAlphanumeric.ReplaceAllString(name, "_")

	// Escape C++ keywords
	if cppKeywords[strings.ToLower(result)] {
		result = result + "_"
	}

	return result
}

// ClassName converts a proto message name to a C++ class name.
// Proto uses PascalCase, which is also the C++ convention for class names.
func ClassName(name string) string {
	return CppName(name)
}

// EnumName converts a proto enum name to a C++ enum name.
func EnumName(name string) string {
	return CppName(name)
}

// EnumValueName converts a proto enum value to a C++ enum value.
// Proto enum values are typically UPPER_SNAKE_CASE.
func EnumValueName(name string) string {
	return CppName(name)
}

// FieldName converts a proto field name to a C++ field name.
// Proto uses snake_case, C++ member variables typically use snake_case with trailing underscore.
func FieldName(name string) string {
	return CppName(name) + "_"
}

// AccessorName converts a proto field name to a C++ getter/setter name.
// The accessor name is the same as the proto field name (snake_case).
func AccessorName(name string) string {
	return CppName(name)
}

// MutableAccessorName returns the mutable accessor name for a field.
func MutableAccessorName(name string) string {
	return "mutable_" + CppName(name)
}

// HasAccessorName returns the has_* accessor name for optional fields.
func HasAccessorName(name string) string {
	return "has_" + CppName(name)
}

// ClearAccessorName returns the clear_* method name for a field.
func ClearAccessorName(name string) string {
	return "clear_" + CppName(name)
}

// SetAccessorName returns the set_* method name for a field.
func SetAccessorName(name string) string {
	return "set_" + CppName(name)
}

// AddAccessorName returns the add_* method name for repeated fields.
func AddAccessorName(name string) string {
	return "add_" + CppName(name)
}

// SizeAccessorName returns the *_size method name for repeated fields.
func SizeAccessorName(name string) string {
	return CppName(name) + "_size"
}

// Namespace converts a proto package name to a C++ namespace.
// Proto packages use dots (foo.bar.baz), C++ uses :: (foo::bar::baz).
func Namespace(pkg string) string {
	if pkg == "" {
		return ""
	}
	parts := strings.Split(pkg, ".")
	for i, part := range parts {
		parts[i] = CppName(part)
	}
	return strings.Join(parts, "::")
}

// NamespaceParts returns the parts of a proto package as C++ namespace components.
func NamespaceParts(pkg string) []string {
	if pkg == "" {
		return nil
	}
	parts := strings.Split(pkg, ".")
	for i, part := range parts {
		parts[i] = CppName(part)
	}
	return parts
}

// HeaderGuard generates a header guard macro from a file path.
func HeaderGuard(filePath string) string {
	// Convert to uppercase, replace non-alphanumeric with underscore
	guard := strings.ToUpper(filePath)
	guard = nonAlphanumeric.ReplaceAllString(guard, "_")
	return guard + "_"
}

// FullyQualifiedCppName returns the fully qualified C++ name for a type.
func FullyQualifiedCppName(pkg, name string) string {
	if pkg == "" {
		return "::" + ClassName(name)
	}
	return "::" + Namespace(pkg) + "::" + ClassName(name)
}

// OneofCaseName returns the enum case name for a oneof field.
func OneofCaseName(oneofName, fieldName string) string {
	return "k" + ToPascalCase(fieldName)
}

// OneofEnumName returns the enum name for a oneof.
func OneofEnumName(oneofName string) string {
	return ToPascalCase(oneofName) + "Case"
}

// ToPascalCase converts a snake_case string to PascalCase.
func ToPascalCase(s string) string {
	parts := strings.Split(s, "_")
	for i, part := range parts {
		if len(part) > 0 {
			runes := []rune(part)
			runes[0] = unicode.ToUpper(runes[0])
			parts[i] = string(runes)
		}
	}
	return strings.Join(parts, "")
}

// ToSnakeCase converts a PascalCase or camelCase string to snake_case.
func ToSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}
