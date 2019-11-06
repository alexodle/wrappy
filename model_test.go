package wrappy

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_DeepCopy(t *testing.T) {
	t1 := &TopLevelType{Type: &ModeledType{
		BaseType: BaseType{
			Name:      "string",
			IsPtr:     true,
			IsBuiltin: true,
		},
		LocalNameForPkg: "string",
	}}

	t2 := t1.DeepCopy().(*TopLevelType)
	t2.Type.(*ModeledType).LocalNameForPkg = "int"
	t2.Type.(*ModeledType).Name = "int"
	t2.Type.(*ModeledType).IsPtr = false

	require.Equal(t, &TopLevelType{Type: &ModeledType{
		BaseType: BaseType{
			Name:      "string",
			IsPtr:     true,
			IsBuiltin: true,
		},
		LocalNameForPkg: "string",
	}}, t1)

	require.Equal(t, &TopLevelType{Type: &ModeledType{
		BaseType: BaseType{
			Name:      "int",
			IsPtr:     false,
			IsBuiltin: true,
		},
		LocalNameForPkg: "int",
	}}, t2)
}
