package compile

import (
	"fmt"
	"sort"

	"github.com/kalo-build/morphe-go/pkg/registry"
	rcfg "github.com/kalo-build/morphe-go/pkg/registry/cfg"
	"github.com/kalo-build/plugin-morphe-morphemap/pkg/mapdef"
)

// ScaffoldConfig holds the configuration for map scaffolding.
type ScaffoldConfig struct {
	RegistryConfig rcfg.MorpheLoadRegistryConfig
	ExternalConfig *rcfg.MorpheLoadRegistryConfig
	OutputPath     string
	Mappings       []MappingRequest
}

// MappingRequest defines a single source→target mapping to scaffold.
type MappingRequest struct {
	Name        string
	SourceAlias string
	SourcePath  string
	TargetAlias string
	TargetPath  string
}

// ScaffoldMorpheMaps generates skeleton .map files for the configured mappings.
func ScaffoldMorpheMaps(config ScaffoldConfig) error {
	// Load the local registry
	localRegistry, err := registry.LoadMorpheRegistry(
		registry.LoadMorpheRegistryHooks{},
		config.RegistryConfig,
	)
	if err != nil {
		return fmt.Errorf("failed to load local Morphe registry: %w", err)
	}

	// Load external registry if configured
	var externalRegistry *registry.Registry
	if config.ExternalConfig != nil {
		extReg, err := registry.LoadMorpheRegistry(
			registry.LoadMorpheRegistryHooks{},
			*config.ExternalConfig,
		)
		if err != nil {
			return fmt.Errorf("failed to load external Morphe registry: %w", err)
		}
		externalRegistry = extReg
	}

	// Generate each mapping
	for _, req := range config.Mappings {
		m, err := scaffoldMap(req, localRegistry, externalRegistry)
		if err != nil {
			return fmt.Errorf("failed to scaffold map %q: %w", req.Name, err)
		}

		if err := mapdef.WriteMapToFile(m, config.OutputPath); err != nil {
			return fmt.Errorf("failed to write map %q: %w", req.Name, err)
		}
	}

	return nil
}

// scaffoldMap generates a skeleton MorpheMap for a single mapping request.
func scaffoldMap(req MappingRequest, localReg *registry.Registry, externalReg *registry.Registry) (*mapdef.MorpheMap, error) {
	m := &mapdef.MorpheMap{
		Name: req.Name,
		Aliases: map[string]string{
			req.SourceAlias: req.SourcePath,
			req.TargetAlias: req.TargetPath,
		},
	}

	// Try to resolve target type fields to pre-populate the scaffold
	targetFields := resolveTypeFields(req.TargetPath, localReg, externalReg)
	sourceFields := resolveTypeFields(req.SourcePath, localReg, externalReg)

	// Check if both are enums
	targetIsEnum := isEnumType(req.TargetPath, localReg, externalReg)
	sourceIsEnum := isEnumType(req.SourcePath, localReg, externalReg)

	if targetIsEnum && sourceIsEnum {
		// Scaffold an enum map
		targetEntries := resolveEnumEntries(req.TargetPath, localReg, externalReg)
		sourceEntries := resolveEnumEntries(req.SourcePath, localReg, externalReg)
		m.Entries = scaffoldEntries(req.TargetAlias, targetEntries, req.SourceAlias, sourceEntries)
	} else {
		// Scaffold a field map
		m.Fields = scaffoldFields(req.TargetAlias, targetFields, req.SourceAlias, sourceFields)
	}

	return m, nil
}

// scaffoldFields generates skeleton field mappings.
// For each target field, it tries to find a matching source field by name.
func scaffoldFields(targetAlias string, targetFields []string, sourceAlias string, sourceFields []string) mapdef.FieldMappings {
	fields := make(mapdef.FieldMappings)
	sourceFieldSet := make(map[string]bool)
	for _, f := range sourceFields {
		sourceFieldSet[f] = true
	}

	sort.Strings(targetFields)

	for _, targetField := range targetFields {
		key := targetAlias + "." + targetField

		// Try exact name match
		if sourceFieldSet[targetField] {
			fields[key] = mapdef.FieldMappingValue{
				IsScalar: true,
				Scalar:   sourceAlias + "." + targetField,
			}
		} else {
			// No match found - leave a TODO placeholder
			fields[key] = mapdef.FieldMappingValue{
				IsScalar: false,
				Object: &mapdef.FieldMapping{
					From: sourceAlias + ".TODO",
				},
			}
		}
	}

	return fields
}

// scaffoldEntries generates skeleton enum entry mappings.
func scaffoldEntries(targetAlias string, targetEntries []string, sourceAlias string, sourceEntries []string) map[string]string {
	entries := make(map[string]string)
	sourceEntrySet := make(map[string]bool)
	for _, e := range sourceEntries {
		sourceEntrySet[e] = true
	}

	sort.Strings(targetEntries)

	for _, targetEntry := range targetEntries {
		key := targetAlias + "." + targetEntry

		// Try exact name match
		if sourceEntrySet[targetEntry] {
			entries[key] = sourceAlias + "." + targetEntry
		} else {
			entries[key] = sourceAlias + ".TODO"
		}
	}

	return entries
}

// resolveTypeFields returns the field names for a Morphe type path.
func resolveTypeFields(typePath string, localReg *registry.Registry, externalReg *registry.Registry) []string {
	// Try local registry first
	if fields := resolveFieldsFromRegistry(typePath, localReg); fields != nil {
		return fields
	}
	// Try external registry
	if externalReg != nil {
		if fields := resolveFieldsFromRegistry(typePath, externalReg); fields != nil {
			return fields
		}
	}
	return nil
}

// resolveFieldsFromRegistry tries to resolve fields from a single registry.
func resolveFieldsFromRegistry(typePath string, reg *registry.Registry) []string {
	if reg == nil {
		return nil
	}

	// Try as model
	for name, model := range reg.GetAllModels() {
		if name == typePath || matchesPath(name, typePath) {
			var fields []string
			for fieldName := range model.Fields {
				fields = append(fields, fieldName)
			}
			return fields
		}
	}

	// Try as structure
	for name, structure := range reg.GetAllStructures() {
		if name == typePath || matchesPath(name, typePath) {
			var fields []string
			for fieldName := range structure.Fields {
				fields = append(fields, fieldName)
			}
			return fields
		}
	}

	// Try as entity
	for name, entity := range reg.GetAllEntities() {
		if name == typePath || matchesPath(name, typePath) {
			var fields []string
			for fieldName := range entity.Fields {
				fields = append(fields, fieldName)
			}
			return fields
		}
	}

	return nil
}

// isEnumType checks if a type path resolves to an enum.
func isEnumType(typePath string, localReg *registry.Registry, externalReg *registry.Registry) bool {
	if localReg != nil {
		for name := range localReg.GetAllEnums() {
			if name == typePath || matchesPath(name, typePath) {
				return true
			}
		}
	}
	if externalReg != nil {
		for name := range externalReg.GetAllEnums() {
			if name == typePath || matchesPath(name, typePath) {
				return true
			}
		}
	}
	return false
}

// resolveEnumEntries returns entry keys for a Morphe enum type path.
func resolveEnumEntries(typePath string, localReg *registry.Registry, externalReg *registry.Registry) []string {
	resolveFromReg := func(reg *registry.Registry) []string {
		if reg == nil {
			return nil
		}
		for name, enum := range reg.GetAllEnums() {
			if name == typePath || matchesPath(name, typePath) {
				var entries []string
				for entryKey := range enum.Entries {
					entries = append(entries, entryKey)
				}
				return entries
			}
		}
		return nil
	}

	if entries := resolveFromReg(localReg); entries != nil {
		return entries
	}
	return resolveFromReg(externalReg)
}

// matchesPath checks if a registry type name matches a path.
// Handles cases where the registry uses simple names but the path includes directory prefixes.
func matchesPath(registryName string, path string) bool {
	// Simple case: direct match
	if registryName == path {
		return true
	}
	// Path may have directory prefix (e.g., "immoscout24/RealEstateHouseBuy")
	// Registry name may be just "RealEstateHouseBuy"
	// For now, check if the registry name matches the last segment
	lastSlash := len(path)
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			lastSlash = i
			break
		}
	}
	if lastSlash < len(path) {
		return registryName == path[lastSlash+1:]
	}
	return false
}
