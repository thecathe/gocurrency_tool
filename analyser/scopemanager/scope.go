package scopemanager

import (
	"fmt"
	"go/ast"
	"go/token"
	"math"

	"github.com/thecathe/gocurrency_tool/analyser/log"
)

type ScopeType string

const (
	SCOPE_TYPE_NONE          ScopeType = "No Scope Found"
	SCOPE_TYPE_PACKAGE       ScopeType = "Package"              // *ast.Package
	SCOPE_TYPE_FILE          ScopeType = "File"                 // *ast.File
	SCOPE_TYPE_FILE_IMPORT   ScopeType = "File Import"          // Peek().Type == File && *ast.GenDecl.tok == import
	SCOPE_TYPE_PACKAGE_VAR   ScopeType = "Package Var"          // Peek().Type == File && *ast.GenDecl.tok == var
	SCOPE_TYPE_PACKAGE_CONST ScopeType = "Package Const"        // Peek().Type == File && *ast.GenDecl.tok == const
	SCOPE_TYPE_FUNC_CALL     ScopeType = "Function Call"        // *ast.CallExpr
	SCOPE_TYPE_FUNC_DECL     ScopeType = "Function Declaration" // *ast.FuncDecl
	SCOPE_TYPE_IF            ScopeType = "If"                   // *ast.
	SCOPE_TYPE_SELECT        ScopeType = "Select"               // *ast.
	SCOPE_TYPE_SWITCH        ScopeType = "Switch"               // *ast.
	SCOPE_TYPE_TYPE_SWITCH   ScopeType = "Type Switch"          // *ast.
	SCOPE_TYPE_FOR           ScopeType = "For Loop"             // *ast.
	SCOPE_TYPE_RANGE         ScopeType = "Ranged For Loop"      // *ast.
	SCOPE_TYPE_DECL          ScopeType = "Declaration"          // *ast.DeclStmt
	SCOPE_TYPE_GOROUTINE     ScopeType = "Goroutine"            // *ast.GoStmt
	// SCOPE_TYPE_GO_NAMED      ScopeType = "Goroutine (Named)"     // *ast.
	// SCOPE_TYPE_GO_ANONYMOUS  ScopeType = "Goroutine (Anonymous)" // *ast.
)

// Scope
type Scope struct {
	ID          ID
	Node        *ast.Node
	Decls       map[ID][]string // not just declarations but assignments too
	Type        ScopeType
	ElevatedIDs *IDs // array of scope ids that elevates all their own decls
	Elevate     bool // signifies if scope should have its decls/assignments elevated to outerscope
}

// Creates a new Scope and adds it to ScopeMap
func (sm *ScopeManager) NewScope(node ast.Node, scope_type ScopeType) *ScopeManager {
	var scope Scope = *NewScope(node, scope_type)

	// add id to stack
	sm = (*sm).Push(scope.ID)

	// add scope to map
	(*(*sm).ScopeMap)[scope.ID] = &scope

	log.GeneralLog("Analyser; NewScope %d: %s\n\n", (*sm).Stack.Size(), scope.ID)

	return sm
}

// Returns a pointer to a Scope.
func NewScope(node ast.Node, scope_type ScopeType) *Scope {
	var scope Scope

	scope.Node = &node
	scope.Type = scope_type
	scope.Decls = make(map[ID][]string, 0)
	scope.ElevatedIDs = NewIDs()

	// should be elevated?
	switch scope_type {
	case SCOPE_TYPE_PACKAGE:
		scope.Elevate = true

	case SCOPE_TYPE_PACKAGE_CONST:
		scope.Elevate = true

	case SCOPE_TYPE_PACKAGE_VAR:
		scope.Elevate = true

	default:
		scope.Elevate = false
	}

	// set ID
	scope.ID = NewScopeID(*scope.Node, scope.Type)

	return &scope
}

// adds decl to map of decl ids and corresponding labels
func (scope *Scope) AddDecl(decl_id ID, decl_label string) *Scope {
	if _, ok := (*scope).Decls[decl_id]; ok {
		// add to existing
		(*scope).Decls[decl_id] = append((*scope).Decls[decl_id], decl_label)
	} else {
		// make new
		(*scope).Decls[decl_id] = append(make([]string, 1), decl_label)
	}
	// (*scope).Decls = (*scope).ElevatedIDs.Append(decl_id)

	return scope
}

// adds scope if to array of elevated ids
func (scope *Scope) ElevateID(scope_id ID) *Scope {
	(*scope).ElevatedIDs = (*scope).ElevatedIDs.Append(scope_id)

	return scope
}

// Returns the Pos of Scope.Node
func (scope *Scope) Pos() token.Pos {
	return (*scope.Node).Pos()
}

// Returns the End of Scope.Node
func (scope *Scope) End() token.Pos {
	return (*scope.Node).End()
}

// MapOfScopes
type MapOfScopes map[ID]*Scope

func NewMapOfScopes() *MapOfScopes {
	return &MapOfScopes{}
}

func (ms *MapOfScopes) ToString() string {
	var _string = ""
	for _, scope := range *ms {
		var _temp = fmt.Sprintf("\nScope: %s\n\tType: %s", scope.ID, scope.Type)
		_string = _string + _temp
	}
	return _string
}

func (ms *MapOfScopes) Size() int {
	return len(*ms)
}

// merge sort
// returns given scope ids ordered by their scopes beginning position
func (sm *ScopeManager) SortScopes(scope_ids *[]ID) *[]ID {

	// return self
	var len_ids int = len(*scope_ids)

	if len_ids <= 1 {
		return scope_ids
	}

	var _new_ids []ID = make([]ID, 0)

	// split
	var half_len int = int(math.Floor(float64(len_ids) / 2.0))
	var side_a IDs = make(IDs, 0)
	var side_b IDs = make(IDs, 0)

	for i := 0; i < len_ids; i++ {
		if i < half_len {
			side_a = append(side_a, (*scope_ids)[i])
		} else {
			side_b = append(side_b, (*scope_ids)[i])
		}
	}

	// log.DebugLog("Init: [%s]\n\t%02d | A: [%s]\n\t%02d | B: [%s]\n\n", strings.Join([]string(*MakeIDs(scope_ids).MakeString()), ", "), len([]string(*side_a.MakeString())), strings.Join([]string(*side_a.MakeString()), ", "), len([]string(*side_b.MakeString())), strings.Join([]string(*side_b.MakeString()), ", "))
	// recursive call
	side_a = *(*sm).SortScopes(side_a.MakeIDs())
	side_b = *(*sm).SortScopes(side_b.MakeIDs())

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
		}
	}

	// log.DebugLog("Merged: %02d | C: [%s]\n\n", len([]string(*MakeIDs(&_new_ids).MakeString())), strings.Join([]string(*MakeIDs(&_new_ids).MakeString()), ", "))
	return &_new_ids
}
