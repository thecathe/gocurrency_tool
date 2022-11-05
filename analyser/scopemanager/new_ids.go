package scopemanager

import (
	"fmt"
	"go/ast"
	"math"
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


// merge sort
// returns given scope ids ordered by their scopes beginning position
func (sm *ScopeManager) SortIDs(_id_type string,_ids *[]ID) *[]ID {

	// return self
	var len_ids int = len(*_ids)

	if len_ids <= 1 {
		return _ids
	}

	var _new_ids []ID = make([]ID, 0)

	// split
	var half_len int = int(math.Floor(float64(len_ids) / 2.0))
	var side_a IDs = make(IDs, 0)
	var side_b IDs = make(IDs, 0)

	for i := 0; i < len_ids; i++ {
		if i < half_len {
			side_a = append(side_a, (*_ids)[i])
		} else {
			side_b = append(side_b, (*_ids)[i])
		}
	}

	// log.DebugLog("Init: [%s]\n\t%02d | A: [%s]\n\t%02d | B: [%s]\n\n", strings.Join([]string(*MakeIDs(scope_ids).MakeString()), ", "), len([]string(*side_a.MakeString())), strings.Join([]string(*side_a.MakeString()), ", "), len([]string(*side_b.MakeString())), strings.Join([]string(*side_b.MakeString()), ", "))
	// recursive call
	side_a = *(*sm).SortIDs(_id_type, side_a.MakeIDs())
	side_b = *(*sm).SortIDs(_id_type, side_b.MakeIDs())

	// log.DebugLog("Post: [%s]\n\t%02d | A: [%s]\n\t%02d | B: [%s]\n\n", strings.Join([]string(*MakeIDs(scope_ids).MakeString()), ", "), len([]string(*side_a.MakeString())), strings.Join([]string(*side_a.MakeString()), ", "), len([]string(*side_b.MakeString())), strings.Join([]string(*side_b.MakeString()), ", "))

	// merge
	for i := 0; i < len_ids; i++ {
		// log.DebugLog("Merge: [%s]\n\t%02d | A: [%s]\n\t%02d | B: [%s]\n\t%02d | C: [%s]\n\n", strings.Join([]string(*MakeIDs(scope_ids).MakeString()), ", "), len([]string(*side_a.MakeString())), strings.Join([]string(*side_a.MakeString()), ", "), len([]string(*side_b.MakeString())), strings.Join([]string(*side_b.MakeString()), ", "), len([]string(*MakeIDs(&_new_ids).MakeString())), strings.Join([]string(*MakeIDs(&_new_ids).MakeString()), ", "))

		// check if either are empty
		if len(side_a) == 0 {
			_new_ids = append(_new_ids, side_b...)

		} else if len(side_b) == 0 {
			_new_ids = append(_new_ids, side_a...)

		} else {

			switch _id_type {

			case "scope":
				// keep adding one with lowest pos
				if (*(*sm).ScopeMap)[(side_a)[0]].Pos() < (*(*sm).ScopeMap)[(side_b)[0]].Pos() {
					_new_ids = append(_new_ids, (side_a)[0])
					// remove from a
					side_a = side_a[1:]
				} else {
					_new_ids = append(_new_ids, (side_b)[0])
					// remove from b
					side_b = side_b[1:]
				}

			case "decl":
				// keep adding one with lowest pos
				if (*(*sm).Decls)[(side_a)[0]].Pos() < (*(*sm).Decls)[(side_b)[0]].Pos() {
					_new_ids = append(_new_ids, (side_a)[0])
					// remove from a
					side_a = side_a[1:]
				} else {
					_new_ids = append(_new_ids, (side_b)[0])
					// remove from b
					side_b = side_b[1:]
				}

			}
		}
	}

	// log.DebugLog("Merged: %02d | C: [%s]\n\n", len([]string(*MakeIDs(&_new_ids).MakeString())), strings.Join([]string(*MakeIDs(&_new_ids).MakeString()), ", "))
	return &_new_ids
}