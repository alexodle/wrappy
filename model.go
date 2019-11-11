package wrappy

import "sort"

type StructStore map[string]*Struct
type InterfaceStore map[string]*Interface
type ImportStore map[string]*Import

func (i *ImportStore) AddAll(other ImportStore) {
	for k, imp := range other {
		(*i)[k] = imp
	}
}

func (i *ImportStore) ToSortedList() ImportList {
	var l ImportList
	for _, v := range *i {
		l = append(l, v)
	}
	sort.Sort(l)
	return l
}

type ParamsList []*Param
type MethodList []*Method
type InterfaceList []*Interface
type ImportList []*Import

type Import struct {
	ImplicitName string
	ExplicitName string
	Path         string
}

type Package struct {
	Name string
	Path string
}

func (p *Package) DeepCopy() *Package {
	if p == nil {
		return nil
	}
	return &Package{
		Name: p.Name,
		Path: p.Path,
	}
}

type Interface struct {
	File                   *File
	Name                   string
	Methods                MethodList
	OriginalStruct         *Struct
	OriginalStructTypeName string
	WrapperStruct          *Struct
}

type File struct {
	Path       string
	Imports    ImportStore
	Package    *Package
	Interfaces InterfaceList
}

type Struct struct {
	Name          string
	File          *File
	PublicMethods MethodList
	PublicFields  ParamsList
}

func (s *Struct) FullName() string {
	return s.File.Package.Path + "." + s.Name
}

type Method struct {
	Name          string
	Receiver      *Param
	Params        ParamsList
	ReturnType    ParamsList
	IsFieldSetter bool
	IsFieldGetter bool
	Field         *Param
}

type Param struct {
	Name string
	Type *TopLevelType
}

// Sorting

func (l ImportList) Len() int {
	return len(l)
}
func (l ImportList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
func (l ImportList) Less(i, j int) bool {
	return l[i].ExplicitName < l[j].ExplicitName
}

func (l InterfaceList) Len() int {
	return len(l)
}
func (l InterfaceList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
func (l InterfaceList) Less(i, j int) bool {
	return l[i].Name < l[j].Name
}

func (l MethodList) Len() int {
	return len(l)
}
func (l MethodList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
func (l MethodList) Less(i, j int) bool {
	return l[i].Name < l[j].Name
}

// Types

type Type interface {
	DeepCopy() Type
	Equal(other Type) bool
}

type TopLevelType struct {
	OriginalType Type
	Type         Type
}

func (t *TopLevelType) DeepCopy() Type {
	tt := &TopLevelType{}
	if t.Type != nil {
		tt.Type = t.Type.DeepCopy()
	}
	if t.OriginalType != nil {
		tt.OriginalType = t.OriginalType.DeepCopy()
	}
	return tt
}

func (t *TopLevelType) Equal(other Type) bool {
	tt, ok := other.(*TopLevelType)
	if !ok {
		return false
	}
	if (t.OriginalType == nil) != (tt.OriginalType == nil) {
		return false
	}
	if (t.Type == nil) != (tt.Type == nil) {
		return false
	}
	if t.Type != nil && !t.Type.Equal(tt.Type) {
		return false
	}
	if t.OriginalType != nil && !t.OriginalType.Equal(tt.OriginalType) {
		return false
	}
	return true
}

type BaseType struct {
	Name           string
	IsBuiltin      bool
	Package        *Package
	IsPtr          bool
	UnderlyingType string
}

func (t *BaseType) FullName() string {
	if t.IsBuiltin {
		return t.Name
	}
	return t.Package.Path + "." + t.Name
}

func (t *BaseType) DeepCopy() Type {
	t2 := *t
	t2.Package = t2.Package.DeepCopy()
	return &t2
}

func (t *BaseType) Equal(other Type) bool {
	tt, ok := other.(*BaseType)
	if !ok {
		return false
	}
	if (t.Package == nil) != (tt.Package == nil) {
		return false
	}
	if t.Package != nil && *t.Package != *tt.Package {
		return false
	}
	return t.Name == tt.Name && t.IsPtr == tt.IsPtr && t.IsBuiltin == tt.IsBuiltin && t.UnderlyingType == t.UnderlyingType
}

type ModeledType struct {
	BaseType
	LocalNameForPkg   string
	NewFuncNameForPkg string
	Interface         *Interface
}

func (t *ModeledType) DeepCopy() Type {
	t2 := *t
	return &t2
}

func (t *ModeledType) Equal(other Type) bool {
	tt, ok := other.(*ModeledType)
	if !ok {
		return false
	}
	if !t.BaseType.Equal(&tt.BaseType) {
		return false
	}
	if (t.Interface == nil) != (tt.Interface == nil) {
		return false
	}
	if t.Interface != nil && t.Interface.Name != tt.Interface.Name {
		return false
	}
	return t.LocalNameForPkg == tt.LocalNameForPkg && t.NewFuncNameForPkg == tt.NewFuncNameForPkg
}

type ArrayType struct {
	Type  Type
	IsPtr bool
}

func (t *ArrayType) DeepCopy() Type {
	t2 := *t
	t2.Type = t2.Type.DeepCopy()
	return &t2
}

func (t *ArrayType) Equal(other Type) bool {
	tt, ok := other.(*ArrayType)
	if !ok {
		return false
	}
	if (t.Type == nil) != (tt.Type == nil) {
		return false
	}
	if t.Type != nil && !t.Type.Equal(tt.Type) {
		return false
	}
	return t.IsPtr == t.IsPtr
}

type MapType struct {
	KeyType   Type
	ValueType Type
	IsPtr     bool
}

func (t *MapType) DeepCopy() Type {
	t2 := *t
	t2.ValueType = t2.ValueType.DeepCopy()
	t2.KeyType = t2.KeyType.DeepCopy()
	return &t2
}

func (t *MapType) Equal(other Type) bool {
	tt, ok := other.(*MapType)
	if !ok {
		return false
	}
	if (t.ValueType == nil) != (tt.ValueType == nil) {
		return false
	}
	if (t.KeyType == nil) != (tt.KeyType == nil) {
		return false
	}
	if t.ValueType != nil && !t.ValueType.Equal(tt.ValueType) {
		return false
	}
	if t.KeyType != nil && !t.KeyType.Equal(tt.KeyType) {
		return false
	}
	return t.IsPtr == t.IsPtr
}

type UnsupportedType struct {
	AstType string
}

func (t *UnsupportedType) DeepCopy() Type {
	t2 := *t
	return &t2
}

func (t *UnsupportedType) Equal(other Type) bool {
	tt, ok := other.(*UnsupportedType)
	if !ok {
		return false
	}
	return t.AstType == tt.AstType
}
