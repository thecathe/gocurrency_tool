package scopemanager

import (
	"fmt"
	"go/ast"
)

// ID
//
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

// Returns ID consisting of
func (sm *ScopeManager) NewVarDeclID(decl *VarDecl) ID {
	if scope_id, ok := (*sm).PeekID(); ok {
		return NewVarDeclID(decl.Label, scope_id)
	}
	// fail
	return ID("Fail: VarDeclID")
}

//
func NewVarDeclID(label string, scope_id ID) ID {
	return ID(fmt.Sprintf("{DECL, %s: %s", label, scope_id))
}
