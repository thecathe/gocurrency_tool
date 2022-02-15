package scopemanager

import (
	"fmt"
	"go/ast"
	"go/token"
)

// ScopeManager
//
type ScopeManager struct {
	ScopeMap *MapOfScopes
	Stack    *StackOfIDs
	FileSet  token.FileSet
	FileSrc  string
	Cache    struct {
		AwaitedFunction   AwaitedFunction
		ExpectingFunction bool
	}
}

//
func NewScopeManager() *ScopeManager {
	return &ScopeManager{}
}

// AwaitedFunction
//
type AwaitedFunction struct {
	Name string
	Pos  token.Pos
	Args *[]ast.Expr
}

func NewAwaitedFunction(node *ast.Node) AwaitedFunction {
	var awaited_function AwaitedFunction

	awaited_function.Name = (*node).(*ast.CallExpr).Fun.(*ast.Ident).Name
	awaited_function.Pos = (*node).Pos()
	awaited_function.Args = &(*node).(*ast.CallExpr).Args

	return awaited_function
}

// ParseType
//
type ParseType string

const (
	PARSE_NONE         ParseType = "None"
	PARSE_FAIL         ParseType = "Failed"
	PARSE_FAIL_DEFAULT ParseType = "Failed by Default"
	PARSE_PASS         ParseType = "Pass, Nothing of Interest"
	PARSE_SCOPE        ParseType = "Scope"
	PARSE_DECL         ParseType = "Declaration"
	PARSE_ASSIGN       ParseType = "Assignment"
	PARSE_SCOPE_END    ParseType = "End of Scope"
	PARSE_STACK_EMPTY  ParseType = "Working Stack is Empty"
)

// Node Type
//
type NodeType string

const (
	NODE_TYPE_OTHER NodeType = "Other"
	NODE_TYPE_SCOPE NodeType = "Scope"
	NODE_TYPE_VAR   NodeType = "Var"
)

//
func (sm *ScopeManager) ParseNode(node *ast.Node) (*ScopeManager, ParseType) {

	// Check if leaving current scope
	if scope, ok := (*sm).Peek(); ok && (*node).Pos() > (*scope.Node).End() {
		// if current node starts after the current scope ends, left current scope
		if _sm, ok := (*sm).Pop(); ok {
			(*sm) = *_sm
		} else {
			// failed
			return sm, PARSE_FAIL
		}
		return sm, PARSE_SCOPE_END
	}

	switch node_type := (*node).(type) {
	case *ast.GenDecl:
		// variable declaration
		switch node_type.Tok {
		case token.VAR:
			// global variable block
		case token.CONST:
			// constant block
		case token.IMPORT:
			// import block
		default:
			// unaccounted for
			return sm, PARSE_FAIL_DEFAULT
		}
	case *ast.FuncLit:
		// anonymous function
	case *ast.FuncDecl:
		// function declaration

	case *ast.AssignStmt:
		// variable assignment, find decl
		switch node_type.Tok {
		case token.DEFINE:
			// declaration
		case token.ASSIGN:
			// variable assignment
		default:
			// unnaccounted for
			return sm, PARSE_FAIL_DEFAULT
		}
	// Expect scope
	case *ast.CallExpr:
		// function call
		(*sm).Cache.AwaitedFunction = NewAwaitedFunction(node)
		(*sm).Cache.ExpectingFunction = true
	}

	return sm, PARSE_PASS

}

// Creates a new Scope and adds it to ScopeMap
func (sm *ScopeManager) NewScope(node *ast.Node, scope_type ScopeType) *ScopeManager {
	var scope Scope = *NewScope(node, scope_type)

	// add scope to map
	(*(*sm).ScopeMap)[scope.ID] = &scope

	return sm
}

// Returns token.Pos of the Scope at the top of the Stack
func (sm *ScopeManager) Pos() token.Pos {
	if scope, ok := (*sm).Peek(); ok {
		return (*scope).Pos()
	}
	return token.NoPos
}

// Returns token.End of the Scope at the top of the Stack
func (sm *ScopeManager) End() token.Pos {
	if scope, ok := (*sm).Peek(); ok {
		return (*scope).End()
	}
	return token.NoPos
}

// Returns FileSet.Position of the Scope at the top of the Stack
func (sm *ScopeManager) Position() token.Position {
	if scope, ok := (*sm).Peek(); ok {
		return (*sm).FileSet.Position((*scope).Pos())
	}
	return (*sm).FileSet.Position(token.NoPos)
}

// Returns the Scope at the top of the Stack, and bool if successful
func (sm *ScopeManager) Peek() (*Scope, bool) {
	if scope_id, ok := (*sm).Stack.Peek(); ok {
		return (*(*sm).ScopeMap)[scope_id], true
	}
	return &Scope{}, false
}

// Adds Scope ID to the top of the Stack
func (sm *ScopeManager) Push(scope_id ID) *ScopeManager {
	(*sm).Stack.Push(scope_id)
	return sm
}

// Removes the Scope ID at the top of the Stack, and a bool if successful
// NOTE: Does not return the ID, use Peek
func (sm *ScopeManager) Pop() (*ScopeManager, bool) {
	stack, ok := (*sm).Stack.Pop()
	if ok {
		(*sm).Stack = stack
	}
	return sm, ok
}

type MapOfScopes map[ID]*Scope

// func (stack *StackOfTraceIDs)

type Scope struct {
	ID   ID
	Node *ast.Node
	Type ScopeType
	Vars *ScopeVarContextMap
	// Trace ScopeTrace // not needed as working stack should keep track
}

// Returns a pointer to a Scope.
// Scope
func NewScope(node *ast.Node, scope_type ScopeType) *Scope {
	var scope Scope

	scope.Node = node
	scope.Type = scope_type
	scope.Vars = NewScopeVarContextMap()

	// set ScopeType
	switch node_type := (*scope.Node).(type) {
	case *ast.File:
		scope.Type = SCOPE_TYPE_FILE
	case *ast.GenDecl:
		scope.Type = SCOPE_TYPE_GEN_DECL
	case *ast.FuncDecl:
		scope.Type = SCOPE_TYPE_FUNC_DECL
	case *ast.IfStmt:
		scope.Type = SCOPE_TYPE_IF
	case *ast.SelectStmt:
		scope.Type = SCOPE_TYPE_SELECT
	case *ast.SwitchStmt:
		scope.Type = SCOPE_TYPE_SWITCH
	case *ast.TypeSwitchStmt:
		scope.Type = SCOPE_TYPE_TYPE_SWITCH
	case *ast.ForStmt:
		scope.Type = SCOPE_TYPE_FOR
	case *ast.RangeStmt:
		scope.Type = SCOPE_TYPE_RANGE
	case *ast.GoStmt:
		// check if to func or anonymous
		switch node_type.Call.Fun.(type) {
		case *ast.Ident:
			// named function
			scope.Type = SCOPE_TYPE_GO_NAMED
		case *ast.FuncLit:
			// anonymous function
			scope.Type = SCOPE_TYPE_GO_ANONYMOUS
		}
	default:
		// not supported
		scope.Type = SCOPE_TYPE_NONE
	}

	// set ID
	scope.ID = NewScopeID(scope.Node, scope.Type)

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

type ScopeVarContextMap map[VarContext]*ScopeVarMap

func NewScopeVarContextMap() *ScopeVarContextMap {
	var _map ScopeVarContextMap

	_map[VAR_CONTEXT_DECLARATION] = NewScopeVarMap()
	_map[VAR_CONTEXT_ASSIGNMENT] = NewScopeVarMap()
	_map[VAR_CONTEXT_PARAMETER] = NewScopeVarMap()

	return &_map
}

type ScopeVarMap map[token.Pos]ScopeVar

func NewScopeVarMap() *ScopeVarMap {
	return &ScopeVarMap{}
}

// ScopeVar
//
type ScopeVar struct {
	Node *ast.Node
}

func NewScopeVar(node ast.Node) {

}

func (scope_var *ScopeVar) ID() token.Pos {
	return (*scope_var.Node).Pos()
}

// VarContext
//
type VarContext string

const (
	VAR_CONTEXT_NONE        VarContext = "None"
	VAR_CONTEXT_DECLARATION VarContext = "Declaration"
	VAR_CONTEXT_ASSIGNMENT  VarContext = "Assignment"
	VAR_CONTEXT_EXPRESSION  VarContext = "Expression"
	VAR_CONTEXT_PARAMETER   VarContext = "Parameter"
)

// TraceID
//
type ID string

// Returns IDTrace of Scope IDs and their index in the slice.
// IDTrace concatenation in the form: index, ID: ...
func NewScopeID(node *ast.Node, scope_type ScopeType) ID {
	return ID(fmt.Sprint("{SCOPE, %s: %v - %v}", scope_type, (*node).Pos(), (*node).End()))
}
func NewVarID(node *ast.Node, var_context VarContext) ID {
	return ID(fmt.Sprint("{VAR, %s: %v - %v}", var_context, (*node).Pos(), (*node).End()))
}
