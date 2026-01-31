package base

import (
	"github.com/aperturerobotics/protobuf-cpp-gen/generator/genfile"
	"github.com/aperturerobotics/protobuf-cpp-gen/internal/naming"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// GenerateEnum generates a C++ enum definition from a protobuf enum descriptor.
func GenerateEnum(g *genfile.GeneratedFile, enum protoreflect.EnumDescriptor) {
	enumName := naming.EnumName(string(enum.Name()))

	g.P("enum ", enumName, " : int {")
	g.In()

	values := enum.Values()
	for i := 0; i < values.Len(); i++ {
		value := values.Get(i)
		valueName := naming.EnumValueName(string(value.Name()))
		valueNum := value.Number()

		comma := ","
		if i == values.Len()-1 {
			comma = ""
		}

		g.Pf("%s = %d%s", valueName, valueNum, comma)
	}

	g.Out()
	g.P("};")
	g.P()

	// Generate helper functions for the enum
	generateEnumHelpers(g, enum, enumName)
}

func generateEnumHelpers(g *genfile.GeneratedFile, enum protoreflect.EnumDescriptor, enumName string) {
	// IsValid function
	g.P("inline bool ", enumName, "_IsValid(int value) {")
	g.In()
	g.P("switch (value) {")
	g.In()

	values := enum.Values()
	for i := 0; i < values.Len(); i++ {
		value := values.Get(i)
		g.Pf("case %d:", value.Number())
	}
	g.In()
	g.P("return true;")
	g.Out()

	g.P("default:")
	g.In()
	g.P("return false;")
	g.Out()
	g.Out()
	g.P("}")
	g.Out()
	g.P("}")
	g.P()

	// Name function (returns string name of enum value)
	g.P("inline const char* ", enumName, "_Name(", enumName, " value) {")
	g.In()
	g.P("switch (value) {")
	g.In()

	for i := 0; i < values.Len(); i++ {
		value := values.Get(i)
		valueName := naming.EnumValueName(string(value.Name()))
		g.Pf("case %s: return \"%s\";", valueName, value.Name())
	}

	g.P("default: return \"\";")
	g.Out()
	g.P("}")
	g.Out()
	g.P("}")
	g.P()

	// Min and max value constants
	if values.Len() > 0 {
		minVal := values.Get(0).Number()
		maxVal := values.Get(0).Number()
		for i := 1; i < values.Len(); i++ {
			v := values.Get(i).Number()
			if v < minVal {
				minVal = v
			}
			if v > maxVal {
				maxVal = v
			}
		}

		g.Pf("constexpr %s %s_MIN = static_cast<%s>(%d);", enumName, enumName, enumName, minVal)
		g.Pf("constexpr %s %s_MAX = static_cast<%s>(%d);", enumName, enumName, enumName, maxVal)
		g.P()
	}
}
