package mapdef_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kalo-build/plugin-morphe-morphemap/pkg/mapdef"
)

type WriteTestSuite struct {
	suite.Suite
}

func TestWriteTestSuite(t *testing.T) {
	suite.Run(t, new(WriteTestSuite))
}

func (suite *WriteTestSuite) TestWriteMapToFile_BasicFieldMap() {
	tmpDir, err := os.MkdirTemp("", "mapdef-write-test")
	suite.NoError(err)
	defer os.RemoveAll(tmpDir)

	m := &mapdef.MorpheMap{
		Name: "OrgToProject",
		Aliases: map[string]string{
			"Org":  "Organization",
			"Proj": "Project",
		},
		Fields: mapdef.FieldMappings{
			"Proj.Name": mapdef.FieldMappingValue{IsScalar: true, Scalar: "Org.Name"},
			"Proj.Code": mapdef.FieldMappingValue{IsScalar: true, Scalar: "Org.Code"},
		},
	}

	writeErr := mapdef.WriteMapToFile(m, tmpDir)

	suite.NoError(writeErr)

	expectedPath := filepath.Join(tmpDir, "org_to_project.map")
	suite.FileExists(expectedPath)

	content, readErr := os.ReadFile(expectedPath)
	suite.NoError(readErr)
	suite.Contains(string(content), "name: OrgToProject")
	suite.Contains(string(content), "Org: Organization")
	suite.Contains(string(content), "Proj: Project")
}

func (suite *WriteTestSuite) TestWriteMapToFile_EnumMap() {
	tmpDir, err := os.MkdirTemp("", "mapdef-write-test")
	suite.NoError(err)
	defer os.RemoveAll(tmpDir)

	m := &mapdef.MorpheMap{
		Name: "StatusToPriority",
		Aliases: map[string]string{
			"Src": "Status",
			"Tgt": "Priority",
		},
		Entries: map[string]string{
			"Tgt.Low":  "Src.Inactive",
			"Tgt.High": "Src.Active",
		},
	}

	writeErr := mapdef.WriteMapToFile(m, tmpDir)

	suite.NoError(writeErr)

	expectedPath := filepath.Join(tmpDir, "status_to_priority.map")
	suite.FileExists(expectedPath)

	content, readErr := os.ReadFile(expectedPath)
	suite.NoError(readErr)
	suite.Contains(string(content), "name: StatusToPriority")
	suite.Contains(string(content), "entries:")
}

func (suite *WriteTestSuite) TestWriteMapToFile_RoundTrip() {
	tmpDir, err := os.MkdirTemp("", "mapdef-write-test")
	suite.NoError(err)
	defer os.RemoveAll(tmpDir)

	original := &mapdef.MorpheMap{
		Name: "TestRoundTrip",
		Aliases: map[string]string{
			"Src": "SourceType",
			"Tgt": "TargetType",
		},
		Fields: mapdef.FieldMappings{
			"Tgt.FieldA": mapdef.FieldMappingValue{IsScalar: true, Scalar: "Src.FieldA"},
		},
	}

	writeErr := mapdef.WriteMapToFile(original, tmpDir)
	suite.NoError(writeErr)

	filePath := filepath.Join(tmpDir, "test_round_trip.map")
	loaded, loadErr := mapdef.LoadMapFromFile(filePath)

	suite.NoError(loadErr)
	suite.Equal(original.Name, loaded.Name)
	suite.Equal(original.Aliases["Src"], loaded.Aliases["Src"])
	suite.Equal(original.Aliases["Tgt"], loaded.Aliases["Tgt"])
	suite.Len(loaded.Fields, 1)
}
