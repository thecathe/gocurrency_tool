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
func (sm *ScopeManager) NewVarValue(expr ast.Expr, pos token.Pos) (VarValue, VarType) {
	var value VarValue
	value.Value = "unknown"

	var _var_type VarType
	_var_type.Type = VAR_DATA_TYPE_NONE

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

		switch inner_expr := value_expr.(type) {

		// Type from Function
		case *ast.CallExpr:
			_var_type = *(*sm).CallExprVarType(inner_expr)

		//
		case *ast.BinaryExpr:
			value.Value = "BinaryExpr"

		//
		case *ast.UnaryExpr:
			value.Value = "UnaryExpr"

		//
		default:
			value.Value = "Other"
		}

	}

	// remove surrounding quotes
	if len(value.Value) >= 2 {
		if value.Value[:1] == "\"" {
			value.Value = value.Value[1:]
		}

		if value.Value[len(value.Value)-1:] == "\"" {
			value.Value = value.Value[:len(value.Value)-1]
		}
	}

	log.DebugLog("Analyser, NewVarValue: %s", value.Value)
	return value, _var_type
}
