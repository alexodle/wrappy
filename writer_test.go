package wrappy

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_convertType(t *testing.T) {
	for _, test := range []struct {
		name           string
		currVar        vvar
		newType        Type
		expectedOutput string
		expectedNewVar *vvar // if nil, expect currVar
	}{
		{
			name: "ensure convert to ptr",
			currVar: vvar{basename: "v", t: &ModeledType{
				LocalNameForPkg: "orig_c.BasicStruct",
				BaseType: BaseType{
					Name:    "BasicStruct",
					Package: &Package{Name: "c", Path: "a/b/c"},
				},
			}},
			newType: &ModeledType{
				LocalNameForPkg: "orig_c.BasicStruct",
				BaseType: BaseType{
					IsPtr:   true,
					Name:    "BasicStruct",
					Package: &Package{Name: "c", Path: "a/b/c"},
				},
			},
			expectedNewVar: &vvar{
				basename: "v",
				i:        1,
				t: &ModeledType{
					LocalNameForPkg: "orig_c.BasicStruct",
					BaseType: BaseType{
						IsPtr:   true,
						Name:    "BasicStruct",
						Package: &Package{Name: "c", Path: "a/b/c"},
					},
				},
			},
			expectedOutput: "v_1 := &v\n",
		},
		{
			name: "ensure convert from ptr",
			currVar: vvar{basename: "v", t: &ModeledType{
				LocalNameForPkg: "orig_c.BasicStruct",
				BaseType: BaseType{
					IsPtr:   true,
					Name:    "BasicStruct",
					Package: &Package{Name: "c", Path: "a/b/c"},
				},
			}},
			newType: &ModeledType{
				LocalNameForPkg: "orig_c.BasicStruct",
				BaseType: BaseType{
					Name:    "BasicStruct",
					Package: &Package{Name: "c", Path: "a/b/c"},
				},
			},
			expectedNewVar: &vvar{
				basename: "v",
				i:        1,
				t: &ModeledType{
					LocalNameForPkg: "orig_c.BasicStruct",
					BaseType: BaseType{
						Name:    "BasicStruct",
						Package: &Package{Name: "c", Path: "a/b/c"},
					},
				},
			},
			expectedOutput: "v_1 := *v\n",
		},
		{
			name: "ensure noop convert is a noop",
			currVar: vvar{basename: "v", t: &ArrayType{
				Type: &ModeledType{
					LocalNameForPkg: "orig_c.BasicStruct",
					BaseType: BaseType{
						Name:    "BasicStruct",
						Package: &Package{Name: "c", Path: "a/b/c"},
					},
				},
			}},
			newType: &ArrayType{
				Type: &ModeledType{
					LocalNameForPkg: "orig_c.BasicStruct",
					BaseType: BaseType{
						Name:    "BasicStruct",
						Package: &Package{Name: "c", Path: "a/b/c"},
					},
				},
			},
			expectedOutput: ``,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			buf := bytes.NewBufferString("")
			newVar := convertType(buf, test.currVar.t, test.newType, test.currVar)
			if test.expectedNewVar != nil {
				require.Equal(t, *test.expectedNewVar, newVar)
			} else {
				require.Equal(t, test.currVar, newVar)
			}
			require.Equal(t, test.expectedOutput, buf.String())
		})
	}
}
