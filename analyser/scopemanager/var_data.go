package scopemanager

import (
	"go/ast"
	"go/token"
)

// MapOfDecls
//
type MapOfDecls map[ID]*VarDecl

func NewMapOfDecls() *MapOfDecls {
	return &MapOfDecls{}
}

// VarDecl
//
type VarDecl struct {
	Node   *ast.Node
	Label  string
	Type   VarType
	Values []VarValue
	Token  token.Token
}

//
func (var_decl *VarDecl) Pos() token.Pos {
	return (*var_decl.Node).Pos()
}

//
func (var_decl *VarDecl) End() token.Pos {
	return (*var_decl.Node).End()
}

//
func (var_decl *VarDecl) ID() token.Pos {
	return (*var_decl.Node).Pos()
}

// GeneralVarType
//
type GeneralVarType string

const (
	VAR_DATA_TYPE_OTHER         GeneralVarType = "Other"
	VAR_DATA_TYPE_INT           GeneralVarType = "Int"
	VAR_DATA_TYPE_STRING        GeneralVarType = "String"
	VAR_DATA_TYPE_STRUCT        GeneralVarType = "Struct"
	VAR_DATA_TYPE_ASYNC_CHANNEL GeneralVarType = "Async. Channel"
	VAR_DATA_TYPE_SYNC_CHANNEL  GeneralVarType = "Sync. Channel"
	VAR_DATA_FUNC_RET           GeneralVarType = "Function Return"
)

// VarValue
//
type VarValue struct {
	Value   string
	Pos     token.Pos
	ScopeID ID
}
