package scopemanager

import (
	"fmt"
	"go/ast"
	"go/token"

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
	ID    ID
	Node  *ast.Node
	Decls *ScopeDeclMap // not just declarations but assignments too
	Type  ScopeType
}

// Creates a new Scope and adds it to ScopeMap
func (sm *ScopeManager) NewScope(node ast.Node, scope_type ScopeType) *ScopeManager {
	var scope Scope = *NewScope(node, scope_type)

	// add id to stack
	sm = (*sm).Push(scope.ID)

	// add scope to map
	(*(*sm).ScopeMap)[scope.ID] = &scope

	log.GeneralLog("Analyser; NewScope %d: %s\n\n", (*sm).StackSize(), scope.ID)
	return sm
}

// Returns a pointer to a Scope.
func NewScope(node ast.Node, scope_type ScopeType) *Scope {
	var scope Scope

	scope.Node = &node
	scope.Type = scope_type
	scope.Decls = NewScopeDeclMap()

	// set ID
	scope.ID = NewScopeID(*scope.Node, scope.Type)

	return &scope
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

// ScopeDeclMap
// Label => NewVarDeclID().ID
type ScopeDeclMap map[ID]ID

func NewScopeDeclMap() *ScopeDeclMap {
	return &ScopeDeclMap{}
}
