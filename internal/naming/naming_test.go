package naming

import "testing"

func TestCppName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"simple", "simple"},
		{"snake_case", "snake_case"},
		{"class", "class_"}, // C++ keyword
		{"void", "void_"},   // C++ keyword
		{"name-with-dash", "name_with_dash"},
		{"name.with.dot", "name_with_dot"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := CppName(tt.input)
			if result != tt.expected {
				t.Errorf("CppName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestClassName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"MyMessage", "MyMessage"},
		{"my_message", "my_message"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ClassName(tt.input)
			if result != tt.expected {
				t.Errorf("ClassName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestNamespace(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"foo", "foo"},
		{"foo.bar", "foo::bar"},
		{"foo.bar.baz", "foo::bar::baz"},
		{"google.protobuf", "google::protobuf"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := Namespace(tt.input)
			if result != tt.expected {
				t.Errorf("Namespace(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestNamespaceParts(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"", nil},
		{"foo", []string{"foo"}},
		{"foo.bar", []string{"foo", "bar"}},
		{"foo.bar.baz", []string{"foo", "bar", "baz"}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := NamespaceParts(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("NamespaceParts(%q) = %v, want %v", tt.input, result, tt.expected)
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("NamespaceParts(%q)[%d] = %q, want %q", tt.input, i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestFieldName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"field", "field_"},
		{"my_field", "my_field_"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := FieldName(tt.input)
			if result != tt.expected {
				t.Errorf("FieldName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestAccessorName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"field", "field"},
		{"my_field", "my_field"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := AccessorName(tt.input)
			if result != tt.expected {
				t.Errorf("AccessorName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestHeaderGuard(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"simple.proto", "SIMPLE_PROTO_"},
		{"path/to/file.proto", "PATH_TO_FILE_PROTO_"},
		{"foo-bar.proto", "FOO_BAR_PROTO_"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := HeaderGuard(tt.input)
			if result != tt.expected {
				t.Errorf("HeaderGuard(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToPascalCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"snake_case", "SnakeCase"},
		{"already", "Already"},
		{"multiple_words_here", "MultipleWordsHere"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ToPascalCase(tt.input)
			if result != tt.expected {
				t.Errorf("ToPascalCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"PascalCase", "pascal_case"},
		{"already", "already"},
		{"MyMessage", "my_message"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ToSnakeCase(tt.input)
			if result != tt.expected {
				t.Errorf("ToSnakeCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
