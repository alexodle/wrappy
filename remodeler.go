package wrappy

import (
	"fmt"
	"sort"
	"strings"
)

func Remodel(structs StructStore, inputDir, outputDir string) []*File {
	modeler := &remodeler{structStore: structs, wrapperStore: InterfaceStore{}, inputDir: inputDir, outputDir: outputDir}
	modeler.buildWrappers()
	modeler.fillWrappers()

	seenFiles := map[string]struct{}{}
	var filesList []*File
	for _, iface := range modeler.wrapperStore {
		if _, ok := seenFiles[iface.File.Path]; ok {
			continue
		}
		filesList = append(filesList, iface.File)

		sort.Sort(iface.File.Interfaces)
		for _, iface := range iface.File.Interfaces {
			sort.Sort(iface.Methods)
		}

		seenFiles[iface.File.Path] = struct{}{}
	}
	return filesList
}

type remodeler struct {
	structStore  StructStore
	wrapperStore InterfaceStore
	inputDir     string
	outputDir    string
}

func shouldWrap(st *Struct) bool {
	return len(st.PublicMethods) > 0 || len(st.PublicFields) > 0
}

func (m *remodeler) buildWrappers() {
	newFiles := map[string]*File{}
	for key, st := range m.structStore {
		if !shouldWrap(st) {
			continue
		}

		newPath := strings.Replace(st.File.Path, m.inputDir, m.outputDir, 1)
		newFile, ok := newFiles[newPath]
		if !ok {
			newFile = &File{
				Path: newPath,
				Imports: ImportStore{
					"orig_" + st.File.Package.Name: &Import{ExplicitName: "orig_" + st.File.Package.Name, Path: st.File.Package.Path},
				},
				Package: &Package{
					Name: st.File.Package.Name,
					Path: strings.Replace(st.File.Package.Path, packagePath(m.inputDir), packagePath(m.outputDir), 1),
				},
			}
			newFiles[newPath] = newFile
		}

		newIFace := &Interface{
			File:                   newFile,
			Name:                   st.Name,
			Methods:                MethodList{},
			OriginalStruct:         st,
			OriginalStructTypeName: "orig_" + st.File.Package.Name + "." + st.Name,
		}
		m.wrapperStore[key] = newIFace
		newFile.Interfaces = append(newFile.Interfaces, newIFace)
	}
}

func (m *remodeler) fillWrappers() {
	for _, iface := range m.wrapperStore {
		newStructName := strings.ToLower(iface.OriginalStruct.Name[0:1]) + iface.OriginalStruct.Name[1:] + "Wrapper"
		recvParam := &Param{Name: "wrapperRcvr", Type: &TopLevelType{Type: &ModeledType{
			BaseType:        BaseType{Name: newStructName, IsPtr: true},
			LocalNameForPkg: newStructName,
		}}}

		iface.WrapperStruct = &Struct{
			File:          iface.File,
			Name:          newStructName,
			PublicMethods: MethodList{},
		}

		fields, fieldImps := m.convertTypesForFile(iface.File, iface.OriginalStruct.PublicFields)
		iface.File.Imports.AddAll(fieldImps)

		for _, f := range fields {
			if isUnsupportedType(f.Type) {
				fmt.Printf("WARN: skipping getter/setter methods for field:%s.%s, at least one param or return type is not currently supported\n", iface.Name, f.Name)
				continue
			}
			setParams := ParamsList{&Param{Name: "v", Type: f.Type}}
			getReturnType := ParamsList{&Param{Type: f.Type}}
			iface.Methods = append(iface.Methods,
				&Method{
					Name:       "Get" + f.Name,
					ReturnType: getReturnType,
				},
				&Method{
					Name:   "Set" + f.Name,
					Params: setParams,
				},
			)
			iface.WrapperStruct.PublicMethods = append(iface.WrapperStruct.PublicMethods,
				&Method{
					Name:          "Get" + f.Name,
					Receiver:      recvParam,
					ReturnType:    getReturnType,
					IsFieldGetter: true,
					Field:         f,
				},
				&Method{
					Name:          "Set" + f.Name,
					Receiver:      recvParam,
					Params:        setParams,
					IsFieldSetter: true,
					Field:         f,
				},
			)
		}

		for _, method := range iface.OriginalStruct.PublicMethods {
			if isUnsupportedMethod(method) {
				fmt.Printf("WARN: skipping method:%s.%s, at least one param or return type is not currently supported\n", iface.Name, method.Name)
				continue
			}

			params, imps1 := m.convertTypesForFile(iface.File, method.Params)
			iface.File.Imports.AddAll(imps1)

			returnType, imps2 := m.convertTypesForFile(iface.File, method.ReturnType)
			iface.File.Imports.AddAll(imps2)

			iface.Methods = append(iface.Methods, &Method{
				Name:       method.Name,
				Params:     params,
				ReturnType: returnType,
			})
			iface.WrapperStruct.PublicMethods = append(iface.WrapperStruct.PublicMethods, &Method{
				Name:       method.Name,
				Params:     params,
				ReturnType: returnType,
				Receiver:   recvParam,
			})
		}
	}
}

func (m *remodeler) convertTypesForFile(f *File, params ParamsList) (ParamsList, ImportStore) {
	var newParams ParamsList
	importStore := ImportStore{}
	for _, p := range params {
		t, imps := m.convertTypeForFile(f, p.Type)
		newParams = append(newParams, &Param{
			Name: p.Name,
			Type: t,
		})
		importStore.AddAll(imps)
	}
	return newParams, importStore
}

func (m *remodeler) convertTypeForFile(f *File, t *TopLevelType) (*TopLevelType, ImportStore) {
	newType, imports, hasWrapper := m.convertTypeForFileRecursive(f, t.Type, false)
	tt := &TopLevelType{Type: newType}
	if hasWrapper {
		origT, origImports, _ := m.convertTypeForFileRecursive(f, t.Type, true)
		imports.AddAll(origImports)
		tt.OriginalType = origT
	}
	return tt, imports
}

func (m *remodeler) convertMapTypeForFile(f *File, t *MapType, ignoreWrappers bool) (Type, ImportStore, bool) {
	// shallow copy
	var newType = t.DeepCopy().(*MapType)

	keyType, imports1, hasWrapper1 := m.convertTypeForFileRecursive(f, t.KeyType, ignoreWrappers)
	if hasWrapper1 {
		panic(fmt.Errorf("map structs not supported"))
	}
	newType.KeyType = keyType

	valType, imports2, hasWrapper2 := m.convertTypeForFileRecursive(f, t.ValueType, ignoreWrappers)
	newType.ValueType = valType

	imports1.AddAll(imports2)
	return newType, imports1, hasWrapper2
}

func (m *remodeler) convertTypeForFileRecursive(f *File, t Type, ignoreWrappers bool) (Type, ImportStore, bool) {
	switch tt := t.(type) {
	case *MapType:
		return m.convertMapTypeForFile(f, tt, ignoreWrappers)
	case *ArrayType:
		newType := tt.DeepCopy().(*ArrayType)
		ct, imps, hasWrapper := m.convertTypeForFileRecursive(f, tt.Type, ignoreWrappers)
		newType.Type = ct
		return newType, imps, hasWrapper
	case *UnsupportedType:
		return t, ImportStore{}, false
	}

	tt := t.(*BaseType)
	if tt.IsBuiltin {
		return &ModeledType{
			BaseType:        *tt,
			LocalNameForPkg: tt.Name,
		}, ImportStore{}, false
	}

	imports := ImportStore{}
	newType := &ModeledType{BaseType: *tt.DeepCopy().(*BaseType)}
	hasWrapper := false
	prefix := "orig_"

	if wrapper, ok := m.wrapperStore[tt.FullName()]; !ignoreWrappers && ok {
		prefix = ""
		hasWrapper = true
		newType.IsPtr = false // interfaces are always inherently pointers
		newType.Interface = wrapper
		newType.Package = wrapper.File.Package

		if wrapper.File.Package.Path == f.Package.Path {
			newType.LocalNameForPkg = wrapper.Name
			newType.NewFuncNameForPkg = "New" + wrapper.Name
			return newType, imports, hasWrapper
		}
	}

	localizeType(newType, imports, prefix)
	return newType, imports, hasWrapper
}

func localizeType(t *ModeledType, imports ImportStore, namePrefix string) {
	importName := namePrefix + t.Package.Name

	requiredImport := &Import{ExplicitName: importName, Path: t.Package.Path}
	imports[importName] = requiredImport

	t.LocalNameForPkg = combine(".", importName, t.Name)
	if t.Interface != nil {
		t.NewFuncNameForPkg = combine(".", importName, "New"+t.Name)
	}
}

func isUnsupportedType(t Type) bool {
	switch tt := t.(type) {
	case *UnsupportedType:
		return true

	case *TopLevelType:
		return isUnsupportedType(tt.Type)
	case *ModeledType, *BaseType:
		return false
	case *MapType:
		return isUnsupportedType(tt.ValueType) || isUnsupportedType(tt.KeyType)
	case *ArrayType:
		return isUnsupportedType(tt.Type)
	}
	panic(fmt.Errorf("unsupported type: %T", t))
}

func hasUnsupportedType(params ParamsList) bool {
	for _, p := range params {
		if isUnsupportedType(p.Type) {
			return true
		}
	}
	return false
}

func isUnsupportedMethod(method *Method) bool {
	return hasUnsupportedType(method.Params) || hasUnsupportedType(method.ReturnType)
}
