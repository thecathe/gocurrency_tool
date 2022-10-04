package scopemanager

import (
	"go/ast"
	"go/token"
)

// AwaitedFunction
//
type AwaitedFunction struct {
	Name     string
	Pos      token.Pos
	Args     *[]ast.Expr
	ParentID ID
}

func NewAwaitedFunction(node *ast.Node, id ID) AwaitedFunction {
	var awaited_function AwaitedFunction

	awaited_function.Name = (*node).(*ast.CallExpr).Fun.(*ast.Ident).Name
	awaited_function.Pos = (*node).Pos()
	awaited_function.Args = &(*node).(*ast.CallExpr).Args
	awaited_function.ParentID = id

	return awaited_function
}