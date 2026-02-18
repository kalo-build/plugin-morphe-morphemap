package mapdef_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"

	"github.com/kalo-build/plugin-morphe-morphemap/pkg/mapdef"
)

type UnmarshalTestSuite struct {
	suite.Suite
}

func TestUnmarshalTestSuite(t *testing.T) {
	suite.Run(t, new(UnmarshalTestSuite))
}

func (suite *UnmarshalTestSuite) TestUnmarshalFieldMappingValue_ScalarString() {
	input := `Src.Name`

	var node yaml.Node
	suite.NoError(yaml.Unmarshal([]byte(input), &node))

	var val mapdef.FieldMappingValue
	unmarshalErr := val.UnmarshalYAML(node.Content[0])

	suite.NoError(unmarshalErr)
	suite.True(val.IsScalar)
	suite.Equal("Src.Name", val.Scalar)
	suite.Nil(val.Object)
}

func (suite *UnmarshalTestSuite) TestUnmarshalFieldMappingValue_ScalarConstant() {
	input := `"some constant value"`

	var node yaml.Node
	suite.NoError(yaml.Unmarshal([]byte(input), &node))

	var val mapdef.FieldMappingValue
	unmarshalErr := val.UnmarshalYAML(node.Content[0])

	suite.NoError(unmarshalErr)
	suite.True(val.IsScalar)
	suite.Equal("some constant value", val.Scalar)
}

func (suite *UnmarshalTestSuite) TestUnmarshalFieldMappingValue_ObjectMapping() {
	input := `
from: Src.FieldX
cast: int
required: true
errorCode: FIELD_REQUIRED
`

	var node yaml.Node
	suite.NoError(yaml.Unmarshal([]byte(input), &node))

	var val mapdef.FieldMappingValue
	unmarshalErr := val.UnmarshalYAML(node.Content[0])

	suite.NoError(unmarshalErr)
	suite.False(val.IsScalar)
	suite.NotNil(val.Object)
	suite.Equal("Src.FieldX", val.Object.From)
	suite.Equal("int", val.Object.Cast)
	suite.True(val.Object.Required)
	suite.Equal("FIELD_REQUIRED", val.Object.ErrorCode)
}

func (suite *UnmarshalTestSuite) TestUnmarshalFieldMappingValue_ObjectWithValueMap() {
	input := `
from: Src.Status
valueMap:
  active: Active
  inactive: Inactive
`

	var node yaml.Node
	suite.NoError(yaml.Unmarshal([]byte(input), &node))

	var val mapdef.FieldMappingValue
	unmarshalErr := val.UnmarshalYAML(node.Content[0])

	suite.NoError(unmarshalErr)
	suite.False(val.IsScalar)
	suite.NotNil(val.Object)
	suite.Equal("Src.Status", val.Object.From)
	suite.Len(val.Object.ValueMap, 2)
	suite.Equal("Active", val.Object.ValueMap["active"])
	suite.Equal("Inactive", val.Object.ValueMap["inactive"])
}

func (suite *UnmarshalTestSuite) TestUnmarshalFieldMappings_NormalMap() {
	input := `
Tgt.Name: Src.Name
Tgt.Code: Src.Code
Tgt.Status:
  from: Src.Status
  cast: string
`

	var fm mapdef.FieldMappings
	unmarshalErr := yaml.Unmarshal([]byte(input), &fm)

	suite.NoError(unmarshalErr)
	suite.Len(fm, 3)

	nameVal := fm["Tgt.Name"]
	suite.True(nameVal.IsScalar)
	suite.Equal("Src.Name", nameVal.Scalar)

	statusVal := fm["Tgt.Status"]
	suite.False(statusVal.IsScalar)
	suite.Equal("Src.Status", statusVal.Object.From)
	suite.Equal("string", statusVal.Object.Cast)
}

func (suite *UnmarshalTestSuite) TestUnmarshalFieldMappings_ScalarRejected() {
	input := `auto`

	var fm mapdef.FieldMappings
	unmarshalErr := yaml.Unmarshal([]byte(input), &fm)

	suite.Error(unmarshalErr)
}

func (suite *UnmarshalTestSuite) TestMarshalFieldMappingValue_Scalar() {
	val := mapdef.FieldMappingValue{
		IsScalar: true,
		Scalar:   "Src.Name",
	}

	result, marshalErr := val.MarshalYAML()

	suite.NoError(marshalErr)
	suite.Equal("Src.Name", result)
}

func (suite *UnmarshalTestSuite) TestMarshalFieldMappingValue_Object() {
	val := mapdef.FieldMappingValue{
		IsScalar: false,
		Object: &mapdef.FieldMapping{
			From: "Src.FieldX",
			Cast: "int",
		},
	}

	result, marshalErr := val.MarshalYAML()

	suite.NoError(marshalErr)
	suite.IsType(&mapdef.FieldMapping{}, result)
	fm := result.(*mapdef.FieldMapping)
	suite.Equal("Src.FieldX", fm.From)
	suite.Equal("int", fm.Cast)
}
