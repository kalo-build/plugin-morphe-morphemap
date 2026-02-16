package compile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"

	rcfg "github.com/kalo-build/morphe-go/pkg/registry/cfg"
	"github.com/kalo-build/plugin-morphe-morphemap/internal/testutils"
	"github.com/kalo-build/plugin-morphe-morphemap/pkg/compile"
	"github.com/kalo-build/plugin-morphe-morphemap/pkg/mapdef"
)

type CompileTestSuite struct {
	suite.Suite

	TestDirPath string

	ModelsDirPath string
	EnumsDirPath  string
}

func TestCompileTestSuite(t *testing.T) {
	suite.Run(t, new(CompileTestSuite))
}

func (suite *CompileTestSuite) SetupTest() {
	suite.TestDirPath = testutils.GetTestDirPath()

	suite.ModelsDirPath = filepath.Join(suite.TestDirPath, "registry", "minimal", "models")
	suite.EnumsDirPath = filepath.Join(suite.TestDirPath, "registry", "minimal", "enums")
}

func (suite *CompileTestSuite) TearDownTest() {
	suite.TestDirPath = ""
}

func (suite *CompileTestSuite) TestScaffoldMorpheMaps_FieldMap() {
	workingDirPath := filepath.Join(suite.TestDirPath, "working-field")
	suite.Nil(os.Mkdir(workingDirPath, 0755))
	defer os.RemoveAll(workingDirPath)

	config := compile.ScaffoldConfig{
		RegistryConfig: rcfg.MorpheLoadRegistryConfig{
			RegistryModelsDirPath:     suite.ModelsDirPath,
			RegistryEnumsDirPath:      suite.EnumsDirPath,
			RegistryStructuresDirPath: filepath.Join(suite.TestDirPath, "registry", "minimal", "structures"),
			RegistryEntitiesDirPath:   filepath.Join(suite.TestDirPath, "registry", "minimal", "entities"),
		},
		OutputPath: workingDirPath,
		Mappings: []compile.MappingRequest{
			{
				Name:        "OrgToProject",
				SourceAlias: "Org",
				SourcePath:  "Organization",
				TargetAlias: "Proj",
				TargetPath:  "Project",
			},
		},
	}

	compileErr := compile.ScaffoldMorpheMaps(config)

	suite.NoError(compileErr)

	// Verify map file was created
	mapPath := filepath.Join(workingDirPath, "org_to_project.map")
	suite.FileExists(mapPath)

	// Load and verify contents
	m, loadErr := mapdef.LoadMapFromFile(mapPath)
	suite.NoError(loadErr)
	suite.Equal("OrgToProject", m.Name)
	suite.Equal("Organization", m.Aliases["Org"])
	suite.Equal("Project", m.Aliases["Proj"])
	suite.NotEmpty(m.Fields)
}

func (suite *CompileTestSuite) TestScaffoldMorpheMaps_EnumMap() {
	workingDirPath := filepath.Join(suite.TestDirPath, "working-enum")
	suite.Nil(os.Mkdir(workingDirPath, 0755))
	defer os.RemoveAll(workingDirPath)

	config := compile.ScaffoldConfig{
		RegistryConfig: rcfg.MorpheLoadRegistryConfig{
			RegistryModelsDirPath:     suite.ModelsDirPath,
			RegistryEnumsDirPath:      suite.EnumsDirPath,
			RegistryStructuresDirPath: filepath.Join(suite.TestDirPath, "registry", "minimal", "structures"),
			RegistryEntitiesDirPath:   filepath.Join(suite.TestDirPath, "registry", "minimal", "entities"),
		},
		OutputPath: workingDirPath,
		Mappings: []compile.MappingRequest{
			{
				Name:        "StatusToPriority",
				SourceAlias: "Src",
				SourcePath:  "Status",
				TargetAlias: "Tgt",
				TargetPath:  "Priority",
			},
		},
	}

	compileErr := compile.ScaffoldMorpheMaps(config)

	suite.NoError(compileErr)

	// Verify map file was created
	mapPath := filepath.Join(workingDirPath, "status_to_priority.map")
	suite.FileExists(mapPath)

	// Load and verify it's an enum map
	m, loadErr := mapdef.LoadMapFromFile(mapPath)
	suite.NoError(loadErr)
	suite.Equal("StatusToPriority", m.Name)
	suite.NotEmpty(m.Entries)
	suite.Equal(mapdef.MapTypeEnum, m.InferMapType())
}

func (suite *CompileTestSuite) TestScaffoldMorpheMaps_MultipleRequests() {
	workingDirPath := filepath.Join(suite.TestDirPath, "working-multi")
	suite.Nil(os.Mkdir(workingDirPath, 0755))
	defer os.RemoveAll(workingDirPath)

	config := compile.ScaffoldConfig{
		RegistryConfig: rcfg.MorpheLoadRegistryConfig{
			RegistryModelsDirPath:     suite.ModelsDirPath,
			RegistryEnumsDirPath:      suite.EnumsDirPath,
			RegistryStructuresDirPath: filepath.Join(suite.TestDirPath, "registry", "minimal", "structures"),
			RegistryEntitiesDirPath:   filepath.Join(suite.TestDirPath, "registry", "minimal", "entities"),
		},
		OutputPath: workingDirPath,
		Mappings: []compile.MappingRequest{
			{
				Name:        "OrgToProject",
				SourceAlias: "Org",
				SourcePath:  "Organization",
				TargetAlias: "Proj",
				TargetPath:  "Project",
			},
			{
				Name:        "StatusToPriority",
				SourceAlias: "Src",
				SourcePath:  "Status",
				TargetAlias: "Tgt",
				TargetPath:  "Priority",
			},
		},
	}

	compileErr := compile.ScaffoldMorpheMaps(config)

	suite.NoError(compileErr)

	suite.FileExists(filepath.Join(workingDirPath, "org_to_project.map"))
	suite.FileExists(filepath.Join(workingDirPath, "status_to_priority.map"))
}

func (suite *CompileTestSuite) TestScaffoldMorpheMaps_FieldMapMatchesSourceFields() {
	workingDirPath := filepath.Join(suite.TestDirPath, "working-match")
	suite.Nil(os.Mkdir(workingDirPath, 0755))
	defer os.RemoveAll(workingDirPath)

	config := compile.ScaffoldConfig{
		RegistryConfig: rcfg.MorpheLoadRegistryConfig{
			RegistryModelsDirPath:     suite.ModelsDirPath,
			RegistryEnumsDirPath:      suite.EnumsDirPath,
			RegistryStructuresDirPath: filepath.Join(suite.TestDirPath, "registry", "minimal", "structures"),
			RegistryEntitiesDirPath:   filepath.Join(suite.TestDirPath, "registry", "minimal", "entities"),
		},
		OutputPath: workingDirPath,
		Mappings: []compile.MappingRequest{
			{
				Name:        "OrgSelfMap",
				SourceAlias: "Src",
				SourcePath:  "Organization",
				TargetAlias: "Tgt",
				TargetPath:  "Organization",
			},
		},
	}

	compileErr := compile.ScaffoldMorpheMaps(config)

	suite.NoError(compileErr)

	mapPath := filepath.Join(workingDirPath, "org_self_map.map")
	m, loadErr := mapdef.LoadMapFromFile(mapPath)
	suite.NoError(loadErr)

	// Organization has 3 fields: ID, Code, Name
	// All should have matching source fields since same type
	suite.Len(m.Fields, 3)

	// All mappings should be scalar (name match)
	for _, val := range m.Fields {
		suite.True(val.IsScalar, "all fields should be scalar for same-type mapping")
	}
}
