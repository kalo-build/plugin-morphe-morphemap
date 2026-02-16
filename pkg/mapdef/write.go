package mapdef

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// WriteMapToFile writes a MorpheMap to a YAML file.
func WriteMapToFile(m *MorpheMap, dirPath string) error {
	filename := toMapFileName(m.Name)
	filePath := filepath.Join(dirPath, filename)

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("failed to create directory for %s: %w", filePath, err)
	}

	data, err := yaml.Marshal(m)
	if err != nil {
		return fmt.Errorf("failed to marshal map %q to YAML: %w", m.Name, err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write map file %s: %w", filePath, err)
	}

	return nil
}

// toMapFileName converts a map name to a file name.
// Example: "IS24HouseBuyToRealEstateListing" → "is24_house_buy_to_real_estate_listing.map"
func toMapFileName(name string) string {
	var result strings.Builder
	for i, r := range name {
		if i > 0 && r >= 'A' && r <= 'Z' {
			prev := rune(name[i-1])
			if prev >= 'a' && prev <= 'z' {
				result.WriteRune('_')
			} else if i+1 < len(name) {
				next := rune(name[i+1])
				if next >= 'a' && next <= 'z' {
					result.WriteRune('_')
				}
			}
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String()) + ".map"
}
