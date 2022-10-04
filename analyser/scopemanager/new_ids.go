package scopemanager

import (
	"fmt"
	"go/ast"
)

// ID
type ID string

// Returns IDTrace of Scope IDs and their index in the slice.
// IDTrace concatenation in the form: index, ID: ...
func NewScopeID(node ast.Node, scope_type ScopeType) ID {
	return ID(fmt.Sprintf("{SCOPE, %s: %v - %v}", scope_type, (node).Pos(), (node).End()))
}

// may be obselete
func NewVarID(node ast.Node, var_context VarContext) ID {
	return ID(fmt.Sprintf("{VAR, %s: %v - %v}", var_context, (node).Pos(), (node).End()))
}

func NewVarDeclID(label string, scope_id ID) ID {
	return ID(fmt.Sprintf("{DECL, %s: %s", label, scope_id))
}
