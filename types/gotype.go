package types

//golang types
//all name is include type. format like name:type

// pkg=type
type Packages map[string]*Package
type Package struct {
	Variables map[string]Variable
	Const     map[string]Const
	Func      map[string]Func
	Struct    map[string]Struct
}
type Type struct {
	Kind Kind
}

type Kind uint8

const (
	Invalid Kind = iota
	StructType
	FuncType
	VariableType
	ConstType
)

type Struct struct {
	Kind   Kind
	Fields []Field
	Name   string
	Ident  string //default is null
}
type Field struct {
	Name string
	Tag  string
}
type Func struct {
	Kind Kind
	Self string
	Args []string //type like name:type
	Resp []string
}

type Variable struct {
	Kind  Kind
	Name  string
	Value string
}

type Const struct {
	Kind  Kind
	Name  string
	Value string
}
