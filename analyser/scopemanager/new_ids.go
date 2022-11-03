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

type IDs []ID

func NewIDs() *IDs {
	var ids IDs = make([]ID, 0)
	return &ids
}

func MakeIDs(_string *[]ID) *IDs {
	var ids IDs = make(IDs, 0)

	for _, _s := range *_string {
		ids = append(ids, _s)
	}

	return &ids
}

func (ids *IDs) MakeString() *[]string {
	var idstrings []string = make([]string, 0)

	for _, _s := range *ids {
		idstrings = append(idstrings, string(_s))
	}

	return &idstrings
}

func (ids *IDs) MakeIDs() *[]ID {
	var idstrings []ID = make([]ID, 0)

	for _, _s := range *ids {
		idstrings = append(idstrings, ID(_s))
	}

	return &idstrings
}

func (_ids *IDs) Append(new_id ID) *IDs {
	var new_ids IDs = make(IDs, len(*_ids)+1)
	copy(new_ids, *_ids)
	new_ids[len(*_ids)] = new_id
	return &new_ids
}
