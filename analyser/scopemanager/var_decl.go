package scopemanager

import (
	"go/ast"
	"go/token"
)

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

//
func (decl *VarDecl) AddValue(value VarValue) *VarDecl {
	decl.Values = append(decl.Values, value)
	return decl
}

// MapOfDecls
//
type MapOfDecls map[ID]*VarDecl

func NewMapOfDecls() *MapOfDecls {
	return &MapOfDecls{}
}