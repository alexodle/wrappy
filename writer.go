package wrappy

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

func WriteCode(files []*File) {
	for _, f := range files {
		writeFile(f)
	}
	for _, f := range files {
		fullPath, err := filepath.Abs(f.Path)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Formatting new file: %s\n", fullPath)
		cmd := exec.Command("goimports", "-w", fullPath)
		output, err := cmd.CombinedOutput()
		if err != nil {
			panic(fmt.Errorf("gofmt failure: %s", string(output)))
		}
	}
}

type vvar struct {
	basename string
	i        int
	t        Type
}

func (v *vvar) next(t Type) vvar {
	return vvar{
		basename: v.basename,
		i:        v.i + 1,
		t:        t,
	}
}

func (v *vvar) name() string {
	if v.i == 0 {
		return v.basename
	}
	return fmt.Sprintf("%s_%d", v.basename, v.i)
}

type vvarlist []vvar

func (v vvarlist) names() []string {
	var names []string
	for _, v := range v {
		names = append(names, v.name())
	}
	return names
}

func writeFile(f *File) {
	fullPath, err := filepath.Abs(f.Path)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Writing new file: %s\n", fullPath)
	if err := os.MkdirAll(filepath.Dir(fullPath), os.ModePerm); err != nil {
		panic(err)
	}
	w, err := os.Create(fullPath)
	if err != nil {
		panic(err)
	}
	defer func() {
		err = w.Close()
		if err != nil {
			panic(err)
		}
	}()

	printf(w, "package %s\n\n", f.Package.Name)

	writeImports(w, f.Imports)
	printf(w, "\n")

	writeInterfaces(w, f.Interfaces)
	printf(w, "\n")
}

func writeImports(w io.Writer, imps ImportStore) {
	var impStrs []string
	for _, imp := range imps {
		impStrs = append(impStrs, fmt.Sprintf("import %s \"%s\"\n", imp.ExplicitName, imp.Path))
	}
	sort.Strings(impStrs)
	for _, s := range impStrs {
		printf(w, s)
	}
}

func writeInterfaces(w io.Writer, interfaces InterfaceList) {
	for _, iface := range interfaces {
		printf(w, "type %s interface {\n", iface.Name)
		printf(w, "GetImpl() *%s\n", iface.OriginalStructTypeName)
		for _, meth := range iface.Methods {
			if len(meth.ReturnType) > 0 {
				printf(w, "%s(%s) (%s)\n", meth.Name, formatParams(meth.Params), formatParams(meth.ReturnType))
			} else {
				printf(w, "%s(%s)\n", meth.Name, formatParams(meth.Params))
			}
		}
		printf(w, "}\n\n")
		writeWrapperStruct(w, iface)
	}
}

func writeWrapperStruct(w io.Writer, iface *Interface) {
	printf(w, "func New%s(impl *%s) %s {\n", iface.Name, iface.OriginalStructTypeName, iface.Name)
	printf(w, "return &%s{impl: impl}\n", iface.WrapperStruct.Name)
	printf(w, "}\n\n")

	printf(w, "type %s struct {\n", iface.WrapperStruct.Name)
	printf(w, "impl *%s\n", iface.OriginalStructTypeName)
	printf(w, "}\n\n")

	printf(w, "func (o *%s) GetImpl() *%s {\n", iface.WrapperStruct.Name, iface.OriginalStructTypeName)
	printf(w, "return o.impl\n")
	printf(w, "}\n\n")
	for _, method := range iface.WrapperStruct.PublicMethods {
		if method.ReturnType != nil {
			printf(w, "func (%s *%s) %s(%s) (%s) {\n", method.Receiver.Name, iface.WrapperStruct.Name, method.Name, formatParams(method.Params), formatParams(method.ReturnType))
		} else {
			printf(w, "func (%s *%s) %s(%s) {\n", method.Receiver.Name, iface.WrapperStruct.Name, method.Name, formatParams(method.Params))
		}

		newVars := unwrapParams(w, method.Params)
		returnVars := applyToImpl(w, method, newVars)
		if len(method.ReturnType) > 0 {
			returnVars = wrapParams(w, method, returnVars)
			printf(w, "return %s\n", strings.Join(returnVars.names(), ", "))
		}

		printf(w, "}\n\n")
	}
}

func isStarred(t Type) bool {
	curr := t
	for {
		switch tt := curr.(type) {
		case *TopLevelType:
			curr = tt.Type
		case *BaseType:
			return tt.IsPtr
		case *ModeledType:
			return tt.IsPtr
		case *ArrayType:
			return tt.IsPtr
		case *MapType:
			return tt.IsPtr
		}
	}
}

func getOriginalType(t *TopLevelType) Type {
	if t.OriginalType != nil {
		return t.OriginalType
	}
	return t.Type
}

func withPtr(t Type) Type {
	t = t.DeepCopy()
	switch tt := t.(type) {
	case *ModeledType:
		tt.IsPtr = true
	case *ArrayType:
		tt.IsPtr = true
	case *MapType:
		tt.IsPtr = true
	default:
		panic(fmt.Sprintf("Unsupported type: %T", t))
	}
	return t
}

func fieldsNeedsRefing(p *Param) bool {
	origType := getOriginalType(p.Type)
	switch tt := origType.(type) {
	case *ModeledType:
		return !tt.IsPtr && !tt.IsBuiltin
	case *ArrayType, *MapType:
		return false
	default:
		panic(fmt.Sprintf("unsupported type: %T", origType))
	}
}

func applyToImpl(w io.Writer, method *Method, vars vvarlist) vvarlist {
	if method.IsFieldSetter {
		printf(w, "%s.impl.%s = %s\n", method.Receiver.Name, method.Field.Name, vars[0].name())
		return nil
	} else if method.IsFieldGetter {
		rt := getOriginalType(method.ReturnType[0].Type)
		// Field getters must be called by reference in order to actually preserve updates
		if fieldsNeedsRefing(method.ReturnType[0]) {
			printf(w, "retval := &%s.impl.%s\n", method.Receiver.Name, method.Field.Name)
			return vvarlist{{basename: "retval", t: withPtr(rt)}}
		} else {
			printf(w, "retval := %s.impl.%s\n", method.Receiver.Name, method.Field.Name)
			return vvarlist{{basename: "retval", t: rt}}
		}
	} else if method.ReturnType != nil {
		var newVars vvarlist
		for i, p := range method.ReturnType {
			v := vvar{basename: fmt.Sprintf("retval%d", i), t: getOriginalType(p.Type)}
			newVars = append(newVars, v)
		}
		printf(w, "%s := %s.impl.%s(%s)\n", strings.Join(newVars.names(), ", "), method.Receiver.Name, method.Name, strings.Join(vars.names(), ", "))
		return newVars
	}
	printf(w, "%s.impl.%s(%s)\n", method.Receiver.Name, method.Name, strings.Join(vars.names(), ", "))
	return nil
}

func convertArrayType(w io.Writer, oldT, newT *ArrayType, oldArrayVar vvar) vvar {
	newArrayVar := oldArrayVar.next(newT)
	derefArray := ""

	unwrap := func() {
		if newT.IsPtr {
			printf(w, "%s = &%s{}\n", newArrayVar.name(), formatTypeWithoutLeadingPtr(newT))
		}

		innerVar := vvar{basename: "it", t: oldT.Type}
		printf(w, "for _, %s := range %s%s {\n", innerVar.name(), derefArray, oldArrayVar.name())
		innerVar = convertType(w, oldT.Type, newT.Type, innerVar)

		printf(w, "%s%s = append(%s%s, %s)\n", derefArray, newArrayVar.name(), derefArray, newArrayVar.name(), innerVar.name())
		printf(w, "}\n")
	}

	printf(w, "var %s %s\n", newArrayVar.name(), formatType(newT))
	if oldArrayVar.t.(*ArrayType).IsPtr {
		derefArray = "*"
		printf(w, "if %s != nil {\n", oldArrayVar.name())
		unwrap()
		printf(w, "}\n")
	} else {
		unwrap()
	}

	return newArrayVar
}

func convertMapType(w io.Writer, oldT, newT *MapType, oldMapVar vvar) vvar {
	newMapVar := oldMapVar.next(newT)
	derefMap := ""
	refMap := ""

	unwrap := func() {
		printf(w, "%s = %s%s{}\n", newMapVar.name(), refMap, formatTypeWithoutLeadingPtr(newT))

		innerVar := vvar{basename: "it", t: oldT.ValueType}

		printf(w, "for k, %s := range %s%s {\n", innerVar.name(), derefMap, oldMapVar.name())
		innerVar = convertType(w, oldT.ValueType, newT.ValueType, innerVar)

		if derefMap != "" {
			printf(w, "(%s%s)[k] = %s\n", derefMap, newMapVar.name(), innerVar.name())
		} else {
			printf(w, "%s[k] = %s\n", newMapVar.name(), innerVar.name())
		}
		printf(w, "}\n")
	}

	printf(w, "var %s %s\n", newMapVar.name(), formatType(newT))

	if newT.IsPtr {
		derefMap = "*"
		refMap = "&"
		printf(w, "if %s != nil {\n", oldMapVar.name())
		unwrap()
		printf(w, "}\n")
	} else {
		unwrap()
	}

	return newMapVar
}

func convertType(w io.Writer, oldT, newT Type, oldVar vvar) vvar {
	if oldT.Equal(newT) {
		return oldVar
	}

	if _, ok := newT.(*ModeledType); ok {
		newVar := oldVar.next(newT)
		if oldT.(*ModeledType).Interface != nil {
			if newT.(*ModeledType).IsPtr {
				printf(w, "%s := %s.GetImpl()\n", newVar.name(), oldVar.name())
			} else {
				printf(w, "%s := *%s.GetImpl()\n", newVar.name(), oldVar.name())
			}
			return newVar

		} else if newT.(*ModeledType).Interface != nil {
			if isStarred(oldVar.t) {
				printf(w, "%s := %s(%s)\n", newVar.name(), newT.(*ModeledType).NewFuncNameForPkg, oldVar.name())
			} else {
				printf(w, "%s := %s(&%s)\n", newVar.name(), newT.(*ModeledType).NewFuncNameForPkg, oldVar.name())
			}
			return newVar

		} else if isStarred(oldVar.t) != newT.(*ModeledType).IsPtr {
			if newT.(*ModeledType).IsPtr {
				printf(w, "%s := &%s\n", newVar.name(), oldVar.name())
			} else {
				printf(w, "%s := *%s\n", newVar.name(), oldVar.name())
			}
			return newVar
		}

		return oldVar
	}

	switch tt := oldT.(type) {
	case *ArrayType:
		return convertArrayType(w, tt, newT.(*ArrayType), oldVar)
	case *MapType:
		return convertMapType(w, tt, newT.(*MapType), oldVar)
	default:
		panic(fmt.Errorf("unsupported type: %T->%T", newT, oldT))
	}
	return oldVar
}

func unwrapParam(w io.Writer, currVar vvar, p *Param) vvar {
	return convertType(w, currVar.t, getOriginalType(p.Type), currVar)
}

func unwrapParams(w io.Writer, params ParamsList) vvarlist {
	var unwrappedVars vvarlist
	for _, p := range params {
		v := vvar{basename: p.Name, t: p.Type.Type}
		unwrappedVars = append(unwrappedVars, unwrapParam(w, v, p))
	}
	return unwrappedVars
}

func wrapParam(w io.Writer, p *Param, v vvar) vvar {
	return convertType(w, v.t, p.Type.Type, v)
}

func wrapParams(w io.Writer, method *Method, vars vvarlist) vvarlist {
	for i, p := range method.ReturnType {
		vars[i] = wrapParam(w, p, vars[i])
	}
	return vars
}

func formatParams(params ParamsList) string {
	var strs []string
	for _, p := range params {
		typeStr := formatType(p.Type)
		if p.Name != "" {
			strs = append(strs, fmt.Sprintf("%s %s", p.Name, typeStr))
		} else {
			strs = append(strs, typeStr)
		}
	}
	return strings.Join(strs, ", ")
}

func formatTypeWithoutLeadingPtr(t Type) string {
	return formatTypeWithOptions(t, true)
}

func formatType(t Type) string {
	return formatTypeWithOptions(t, false)
}

func formatTypeWithOptions(t Type, ignoreLeadingPtr bool) string {
	switch tt := t.(type) {
	case *TopLevelType:
		return formatTypeWithOptions(tt.Type, ignoreLeadingPtr)
	case *ModeledType:
		if tt.IsPtr && !ignoreLeadingPtr {
			return "*" + tt.LocalNameForPkg
		} else {
			return tt.LocalNameForPkg
		}
	case *ArrayType:
		if tt.IsPtr && !ignoreLeadingPtr {
			return "*[]" + formatTypeWithOptions(tt.Type, false)
		} else {
			return "[]" + formatTypeWithOptions(tt.Type, false)
		}
	case *MapType:
		if tt.IsPtr && !ignoreLeadingPtr {
			return fmt.Sprintf("*map[%s]%s",
				formatTypeWithOptions(tt.KeyType, false),
				formatTypeWithOptions(tt.ValueType, false))
		} else {
			return fmt.Sprintf("map[%s]%s",
				formatTypeWithOptions(tt.KeyType, false),
				formatTypeWithOptions(tt.ValueType, false))
		}
	}
	panic(fmt.Errorf("unsupported type: %T", t))
}

func printf(w io.Writer, s string, args ...interface{}) {
	_, err := fmt.Fprintf(w, s, args...)
	if err != nil {
		panic(err)
	}
}
