package generator

import (
	"fmt"
	"sort"

	"github.com/aperturerobotics/go-protobuf-cpp-gen/generator/genfile"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Feature is a function that creates a FeatureGenerator for a generated file.
type Feature func(gen *GeneratedFile) FeatureGenerator

// FeatureGenerator generates code for a specific feature.
type FeatureGenerator interface {
	// Name returns the feature name.
	Name() string
	// GenerateHeader generates header code for a message.
	GenerateHeader(msg protoreflect.MessageDescriptor) bool
	// GenerateSource generates source code for a message.
	GenerateSource(msg protoreflect.MessageDescriptor) bool
}

var defaultFeatures = make(map[string]Feature)

// RegisterFeature registers a feature by name.
func RegisterFeature(name string, feat Feature) {
	defaultFeatures[name] = feat
}

// findFeatures returns the features for the given names.
func findFeatures(featureNames []string) ([]Feature, error) {
	required := make(map[string]Feature)
	for _, name := range featureNames {
		if name == "all" {
			required = defaultFeatures
			break
		}
		feat, ok := defaultFeatures[name]
		if !ok {
			return nil, fmt.Errorf("unknown feature: %q", name)
		}
		required[name] = feat
	}

	type namefeat struct {
		name string
		feat Feature
	}
	var sorted []namefeat
	for name, feat := range required {
		sorted = append(sorted, namefeat{name, feat})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].name < sorted[j].name
	})

	var features []Feature
	for _, sp := range sorted {
		features = append(features, sp.feat)
	}
	return features, nil
}

// GeneratedFile extends genfile.GeneratedFile with additional context.
type GeneratedFile struct {
	*genfile.GeneratedFile
}

// NewGeneratedFile creates a new GeneratedFile wrapping a genfile.GeneratedFile.
func NewGeneratedFile(gf *genfile.GeneratedFile) *GeneratedFile {
	return &GeneratedFile{GeneratedFile: gf}
}

// FeatureRegistry holds registered features for a generator instance.
type FeatureRegistry struct {
	features []FeatureGenerator
}

// NewFeatureRegistry creates a new feature registry.
func NewFeatureRegistry() *FeatureRegistry {
	return &FeatureRegistry{}
}

// Register adds a feature generator to the registry.
func (r *FeatureRegistry) Register(f FeatureGenerator) {
	r.features = append(r.features, f)
}

// GenerateHeader runs all features' header generation for a message.
func (r *FeatureRegistry) GenerateHeader(msg protoreflect.MessageDescriptor) {
	for _, f := range r.features {
		f.GenerateHeader(msg)
	}
}

// GenerateSource runs all features' source generation for a message.
func (r *FeatureRegistry) GenerateSource(msg protoreflect.MessageDescriptor) {
	for _, f := range r.features {
		f.GenerateSource(msg)
	}
}
