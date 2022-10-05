package scopemanager

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/thecathe/gocurrency_tool/analyser/log"
)

// VarValue
type VarValue struct {
	Value   string
	Pos     token.Pos
	ScopeID ID
}

// Returns VarValue using ast.ValueSpec .Values[]ast.Expr and .Pos
func (sm *ScopeManager) NewVarValue(expr ast.Expr, pos token.Pos) VarValue {
	var value VarValue

	value.Pos = pos
	if scope_id, ok := (*sm).PeekID(); ok {
		value.ScopeID = scope_id
	}

	switch value_expr := expr.(type) {
	// simple value
	case *ast.BasicLit:
		value.Value = fmt.Sprintf("%v", value_expr.Value)
	case *ast.Ident:
		value.Value = fmt.Sprintf("%v", value_expr.Name)

	// add as is
	default:
		var expr_str string
		switch value_expr.(type) {
		case *ast.CallExpr:
			expr_str = "CallExpr"
		case *ast.BinaryExpr:
			expr_str = "BinaryExpr"
		case *ast.UnaryExpr:
			expr_str = "UnaryExpr"
		default:
			expr_str = "Other"
		}
		value.Value = fmt.Sprintf("%v", expr_str)

	}

	log.DebugLog("Analyser, NewVarValue: %s", value.Value)
	return value
}
