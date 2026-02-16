package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	rcfg "github.com/kalo-build/morphe-go/pkg/registry/cfg"
	"github.com/kalo-build/plugin-morphe-morphemap/pkg/compile"
)

// PluginConfig represents the configuration passed to the plugin by Kalo CLI
type PluginConfig struct {
	// Store-based paths (mounted by CLI)
	Stores map[string]StoreConfig `json:"stores,omitempty"`

	// Legacy direct paths (for backward compatibility)
	InputPath  string `json:"inputPath,omitempty"`
	OutputPath string `json:"outputPath,omitempty"`

	// Plugin-specific config
	Config  ConfigEntries `json:"config,omitempty"`
	Verbose bool          `json:"verbose,omitempty"`
}

// StoreConfig represents a store configuration from Kalo CLI
type StoreConfig struct {
	ID        uint32 `json:"id"`
	Type      string `json:"type"`
	MountPath string `json:"mountPath,omitempty"`
}

// ConfigEntries holds plugin-specific configuration
type ConfigEntries struct {
	Mappings []MappingConfig `json:"mappings,omitempty"`
}

// MappingConfig defines a single source→target mapping to scaffold
type MappingConfig struct {
	Name        string `json:"name"`
	SourceAlias string `json:"sourceAlias"`
	SourcePath  string `json:"sourcePath"`
	TargetAlias string `json:"targetAlias"`
	TargetPath  string `json:"targetPath"`
}

// Exit codes
const (
	ExitSuccess         = 0
	ExitCompileFailed   = 1
	ExitMissingConfig   = 3
	ExitInvalidConfig   = 4
	ExitInputPathError  = 12
	ExitOutputPathError = 13
)

// logInfo prints info messages only when verbose mode is enabled
func logInfo(verbose bool, format string, args ...interface{}) {
	if verbose {
		fmt.Fprintf(os.Stdout, format+"\n", args...)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: plugin-morphe-morphemap <config>")
		fmt.Fprintln(os.Stderr, "  config: JSON string with store configurations")
		os.Exit(ExitMissingConfig)
	}

	// Parse configuration
	rawConfig := os.Args[1]
	var config PluginConfig
	if err := json.Unmarshal([]byte(rawConfig), &config); err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing config JSON:", err)
		os.Exit(ExitInvalidConfig)
	}

	// Determine paths - prefer store mounts, fall back to legacy paths
	var registryPath, externalPath, outputPath string

	if config.Stores != nil {
		for _, store := range config.Stores {
			switch store.MountPath {
			case "/registry":
				registryPath = "/registry"
			case "/external":
				externalPath = "/external"
			case "/input":
				if registryPath == "" {
					registryPath = "/input"
				}
			case "/output":
				outputPath = "/output"
			}
		}
	}

	// Fall back to legacy paths
	if registryPath == "" && config.InputPath != "" {
		registryPath = config.InputPath
	}
	if outputPath == "" && config.OutputPath != "" {
		outputPath = config.OutputPath
	}

	// Validate required paths
	if registryPath == "" {
		fmt.Fprintln(os.Stderr, "Error: registry input path is required (mount /registry store or provide inputPath)")
		os.Exit(ExitInputPathError)
	}
	if outputPath == "" {
		fmt.Fprintln(os.Stderr, "Error: output path is required (mount /output store or provide outputPath)")
		os.Exit(ExitOutputPathError)
	}

	// Validate mappings config
	if len(config.Config.Mappings) == 0 {
		fmt.Fprintln(os.Stderr, "Error: at least one mapping must be configured")
		os.Exit(ExitInvalidConfig)
	}

	for i, m := range config.Config.Mappings {
		if m.Name == "" || m.SourceAlias == "" || m.SourcePath == "" || m.TargetAlias == "" || m.TargetPath == "" {
			fmt.Fprintf(os.Stderr, "Error: mapping[%d] is missing required fields (name, sourceAlias, sourcePath, targetAlias, targetPath)\n", i)
			os.Exit(ExitInvalidConfig)
		}
	}

	verbose := config.Verbose
	logInfo(verbose, "Registry path: %s", registryPath)
	if externalPath != "" {
		logInfo(verbose, "External path: %s", externalPath)
	}
	logInfo(verbose, "Output path: %s", outputPath)
	logInfo(verbose, "Scaffolding %d mapping(s)", len(config.Config.Mappings))

	// Build registry config
	registryConfig := rcfg.MorpheLoadRegistryConfig{
		RegistryModelsDirPath:     filepath.Join(registryPath, "models"),
		RegistryEntitiesDirPath:   filepath.Join(registryPath, "entities"),
		RegistryEnumsDirPath:      filepath.Join(registryPath, "enums"),
		RegistryStructuresDirPath: filepath.Join(registryPath, "structures"),
	}

	// Build external registry config (optional)
	var externalConfig *rcfg.MorpheLoadRegistryConfig
	if externalPath != "" {
		externalConfig = &rcfg.MorpheLoadRegistryConfig{
			RegistryModelsDirPath:     filepath.Join(externalPath, "models"),
			RegistryEntitiesDirPath:   filepath.Join(externalPath, "entities"),
			RegistryEnumsDirPath:      filepath.Join(externalPath, "enums"),
			RegistryStructuresDirPath: filepath.Join(externalPath, "structures"),
		}
	}

	// Convert mapping configs
	mappings := make([]compile.MappingRequest, len(config.Config.Mappings))
	for i, m := range config.Config.Mappings {
		mappings[i] = compile.MappingRequest{
			Name:        m.Name,
			SourceAlias: m.SourceAlias,
			SourcePath:  m.SourcePath,
			TargetAlias: m.TargetAlias,
			TargetPath:  m.TargetPath,
		}
	}

	// Build scaffold config
	scaffoldConfig := compile.ScaffoldConfig{
		RegistryConfig: registryConfig,
		ExternalConfig: externalConfig,
		OutputPath:     outputPath,
		Mappings:       mappings,
	}

	// Run scaffolding
	logInfo(verbose, "Starting map scaffolding...")
	if err := compile.ScaffoldMorpheMaps(scaffoldConfig); err != nil {
		fmt.Fprintln(os.Stderr, "Map scaffolding failed:", err)
		os.Exit(ExitCompileFailed)
	}

	logInfo(verbose, "Map scaffolding completed successfully")
	os.Exit(ExitSuccess)
}
