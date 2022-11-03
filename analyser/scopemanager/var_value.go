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

	var _var_type VarType = NewVarType()
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
		// check bool
		if value_expr.Name == "true" {
			_var_type.Type = VAR_DATA_TYPE_BOOL
		} else if value_expr.Name == "false" {
			_var_type.Type = VAR_DATA_TYPE_BOOL
		} else {
			// find corresponding decl
			if _decl_id, _scope_id, _elevated, _elevated_id := (*sm).FindDeclID(value_expr.Name); _scope_id != "" {
				_var_type.Type = (*(*sm).Decls)[_decl_id].Type.Type

				if _elevated {
					_scope_id = _elevated_id
				}

				if _index, _var_value := (*(*sm).Decls)[_decl_id].FindValue(_scope_id); _index >= 0 {
					value.Value = fmt.Sprintf("%v", _var_value.Value)
				} else {
					value.Value = fmt.Sprintf("ident: %v", value_expr.Name)
					log.FailureLog("NewVarValue; *ast.Ident: label: \"%s\" did not yield a decl_id", value_expr.Name)
				}
			} else {
				value.Value = fmt.Sprintf("ident: %v", value_expr.Name)
				log.FailureLog("NewVarValue; *ast.Ident: label: \"%s\" did not yield a decl_id", value_expr.Name)
			}
		}

	// add as is
	default:

		switch inner_expr := value_expr.(type) {

		// Type from Function
		case *ast.CallExpr:
			_var_type.Info["ValueTrace"] = fmt.Sprintf("%s > %s", _var_type.Info["ValueTrace"], "*ast.CallExpr")
			_var_type = *(*sm).CallExprVarType(&_var_type, inner_expr)
			// patch var value
			if _var_type.Type == VAR_DATA_TYPE_ASYNC_CHAN || _var_type.Type == VAR_DATA_TYPE_SYNC_CHAN || _var_type.Type == VAR_DATA_TYPE_CHAN {
				value.Value = string(_var_type.Data.Get(len(_var_type.Data) - 1))
			}

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

	log.DebugLog("NewVarValue: %s", value.Value)
	return value, _var_type
}
