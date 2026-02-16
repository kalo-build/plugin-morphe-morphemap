package mapdef_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kalo-build/plugin-morphe-morphemap/internal/testutils"
	"github.com/kalo-build/plugin-morphe-morphemap/pkg/mapdef"
)

type LoadTestSuite struct {
	suite.Suite

	TestDirPath string
}

func TestLoadTestSuite(t *testing.T) {
	suite.Run(t, new(LoadTestSuite))
}

func (suite *LoadTestSuite) SetupTest() {
	suite.TestDirPath = testutils.GetTestDirPath()
}

func (suite *LoadTestSuite) TearDownTest() {
	suite.TestDirPath = ""
}

func (suite *LoadTestSuite) TestLoadMapFromFile_ValidFieldMap() {
	filePath := filepath.Join(suite.TestDirPath, "maps", "valid_field.map")

	m, loadErr := mapdef.LoadMapFromFile(filePath)

	suite.NoError(loadErr)
	suite.NotNil(m)
	suite.Equal("OrgToProject", m.Name)
	suite.Len(m.Aliases, 2)
	suite.Equal("Organization", m.Aliases["Org"])
	suite.Equal("Project", m.Aliases["Proj"])
	suite.Len(m.Fields, 3)
}

func (suite *LoadTestSuite) TestLoadMapFromFile_ValidEnumMap() {
	filePath := filepath.Join(suite.TestDirPath, "maps", "valid_enum.map")

	m, loadErr := mapdef.LoadMapFromFile(filePath)

	suite.NoError(loadErr)
	suite.NotNil(m)
	suite.Equal("StatusToPriority", m.Name)
	suite.Len(m.Aliases, 2)
	suite.Len(m.Entries, 3)
	suite.Equal("Src.Inactive", m.Entries["Tgt.Low"])
}

func (suite *LoadTestSuite) TestLoadMapFromFile_MissingName() {
	filePath := filepath.Join(suite.TestDirPath, "maps", "missing_name.map")

	m, loadErr := mapdef.LoadMapFromFile(filePath)

	suite.Error(loadErr)
	suite.Nil(m)
	suite.Contains(loadErr.Error(), "missing required 'name' field")
}

func (suite *LoadTestSuite) TestLoadMapFromFile_MissingAliases() {
	filePath := filepath.Join(suite.TestDirPath, "maps", "missing_aliases.map")

	m, loadErr := mapdef.LoadMapFromFile(filePath)

	suite.Error(loadErr)
	suite.Nil(m)
	suite.Contains(loadErr.Error(), "missing required 'aliases' field")
}

func (suite *LoadTestSuite) TestLoadMapFromFile_NonexistentFile() {
	filePath := filepath.Join(suite.TestDirPath, "maps", "does_not_exist.map")

	m, loadErr := mapdef.LoadMapFromFile(filePath)

	suite.Error(loadErr)
	suite.Nil(m)
}

func (suite *LoadTestSuite) TestLoadMapsFromDirectory_ValidDir() {
	dirPath := filepath.Join(suite.TestDirPath, "maps")

	maps, loadErr := mapdef.LoadMapsFromDirectory(dirPath)

	// Should load valid files but fail on invalid ones
	// The directory contains files with missing name/aliases which will error
	suite.Error(loadErr)
	_ = maps
}

func (suite *LoadTestSuite) TestLoadMapsFromDirectory_NonexistentDir() {
	dirPath := filepath.Join(suite.TestDirPath, "nonexistent")

	maps, loadErr := mapdef.LoadMapsFromDirectory(dirPath)

	suite.Error(loadErr)
	suite.Nil(maps)
}
