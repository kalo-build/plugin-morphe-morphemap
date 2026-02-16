package mapdef_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kalo-build/plugin-morphe-morphemap/pkg/mapdef"
)

type TypesTestSuite struct {
	suite.Suite
}

func TestTypesTestSuite(t *testing.T) {
	suite.Run(t, new(TypesTestSuite))
}

func (suite *TypesTestSuite) TestInferMapType_FieldMap() {
	m := mapdef.MorpheMap{
		Name:    "TestFieldMap",
		Aliases: map[string]string{"Src": "Source", "Tgt": "Target"},
		Fields: mapdef.FieldMappings{
			"Tgt.Name": mapdef.FieldMappingValue{IsScalar: true, Scalar: "Src.Name"},
		},
	}

	result := m.InferMapType()

	suite.Equal(mapdef.MapTypeField, result)
}

func (suite *TypesTestSuite) TestInferMapType_EnumMap() {
	m := mapdef.MorpheMap{
		Name:    "TestEnumMap",
		Aliases: map[string]string{"Src": "StatusA", "Tgt": "StatusB"},
		Entries: map[string]string{
			"Tgt.Active": "Src.Active",
		},
	}

	result := m.InferMapType()

	suite.Equal(mapdef.MapTypeEnum, result)
}

func (suite *TypesTestSuite) TestInferMapType_EmptyDefaultsToField() {
	m := mapdef.MorpheMap{
		Name:    "EmptyMap",
		Aliases: map[string]string{"Src": "Source", "Tgt": "Target"},
	}

	result := m.InferMapType()

	suite.Equal(mapdef.MapTypeField, result)
}

