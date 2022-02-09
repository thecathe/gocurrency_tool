package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"os"
)

// // //
// // // Scope Manager
// // //
type ScopeManager struct {
	Scopes       ScopeMap     // all scopes
	Decls        DataDeclMap  // all declarations
	WorkingStack ScopeIDStack // working stack
	FileSet      *token.FileSet
	FileSrc      string
}

// // //
// // // Scope
// // //
type ScopeMap map[ScopeID]Scope
type ScopeID string
type Scope struct {
	ID                ScopeID
	Node              ast.Node
	Type              ScopeType
	Position          token.Position
	ParamIDs          []ScopeDataID
	ContainedData     ScopeDataMap
	ContainedDataDecl []DataDeclID
	ChildScopeIDs     []ScopeID
	ScopeIDStack      ScopeIDStack
}

// // //
// // // Scope Stack
// // //
type ScopeIDStackID string

// ID is a trace of the ScopeIDs
// ParsingScope signals if the top of the stack is still being initialised.
type ScopeIDStack struct {
	ID       ScopeIDStackID
	StartPos token.Pos
	EndPos   token.Pos
	Size     int
	ScopeIDs []ScopeID
}

// // //
// // // Scope Data
// // //
type ScopeDataMap map[ScopeDataID]ScopeData
type ScopeDataID string
type ScopeData struct {
	ID       ScopeDataID
	Node     ast.Node
	Name     string
	Position token.Position
	Value    string
	DataType []DataType
	DeclID   DataDeclID
	IsParam  bool
	IsDecl   bool
}

// // //
// // // Scope Type
// // //
type ScopeType string

const (
	SCOPE_NONE               ScopeType = "None"
	SCOPE_FILE               ScopeType = "File"
	SCOPE_GEN_DECL           ScopeType = "Decl"
	SCOPE_FUNC_DECL          ScopeType = "Function"
	SCOPE_IF                 ScopeType = "If"
	SCOPE_SELECT             ScopeType = "Select"
	SCOPE_SWITCH             ScopeType = "Switch"
	SCOPE_TYPE_SWITCH        ScopeType = "Type Switch"
	SCOPE_FOR                ScopeType = "For Loop"
	SCOPE_RANGE              ScopeType = "Ranged For Loop"
	SCOPE_GO_NAMED           ScopeType = "Goroutine (Named)"
	SCOPE_GO_ANONYMOUS       ScopeType = "Goroutine (Anonymous)"
	SCOPE_ANONYMOUS_FUNCTION ScopeType = "Anonymous Function"
)

// // //
// // // Data Decl
// // //
type DataDeclMap map[DataDeclID]DataDecl
type DataDeclID string
type DataDecl struct {
	ID       DataDeclID
	Position token.Position
	DataType DataType
	DeclType DeclType
	Params   []ScopeData
}

type DataType string

const (
	DATA_NONE          DataType = "None"
	DATA_VARIABLE      DataType = "Variable"
	DATA_CHANNEL       DataType = "Channel"
	DATA_ASYNC_CHANNEL DataType = "Asynchronous Channel"
	DATA_SYNC_CHANNEL  DataType = "Synchronous Channel"
	DATA_INT           DataType = "Integer"
	DATA_STRING        DataType = "String"
)

type DeclType string

const (
	DECL_NONE     DeclType = "None"
	DECL_ASSIGNED DeclType = "Assigned"
	DECL_DECLARED DeclType = "Declared"
	DECL_FUNCTION DeclType = "Function"
)

// // //
// // // Parse Type
// // //
type ParseType string

const (
	PARSE_NONE        ParseType = "None"
	PARSE_SCOPE       ParseType = "Scope"
	PARSE_DECL        ParseType = "Declaration"
	PARSE_ASSIGN      ParseType = "Assignment"
	PARSE_DUPL        ParseType = "Duplicate Node, Already Parsed"
	PARSE_SCOPE_END   ParseType = "End of Scope"
	PARSE_STACK_EMPTY ParseType = "Working Stack is Empty"
)

// // // // // //
// // // // // //
// // // // // //
// // // // // // Scope Manager
// // // // // //
// // // // // //
// // // // // //

// Returns a new ScopeManager.
// The node provided should be of type *ast.File, or it will fail.
// Contains the file scope and working stack, and all fields are initialised.
//
// (name of file, fileset of file, ast of file) -> (scope manager, ok)
func NewScopeManager(filename string, fileset *token.FileSet, node ast.Node) (*ScopeManager, bool) {
	var scope_manager ScopeManager = ScopeManager{}

	// first scope should be file
	switch node.(type) {
	case *ast.File:
		// do nothing
	default:
		// failure
		fmt.Printf("Scope, NewScopeManager(): Initial Scope should be a file.\n")
		return &scope_manager, false
	}

	// must be able to read file.
	if _file_src, err := os.ReadFile(filename); err == nil {
		// got source file
		scope_manager.FileSrc = string(_file_src)
		scope_manager.FileSet = fileset

		scope_manager.Scopes = *NewScopeMap()
		scope_manager.Decls = *NewDataDeclMap()

		// create scope
		if _scope_manager, _initial_scope_id, ok := scope_manager.NewScope(node); ok && _initial_scope_id == "initial_scope" {
			scope_manager = *_scope_manager
			// update position
			var _initial_scope Scope = _scope_manager.Scopes[_initial_scope_id]
			_initial_scope = *_initial_scope.UpdateID()
			// update scope to manager
			_scope_manager.Scopes[_initial_scope.ID] = _initial_scope

			// add scope id to stack
			scope_manager.WorkingStack = NewWorkingScopeStack(_initial_scope.ID, _initial_scope.Node.Pos(), _initial_scope.Node.End())
			scope_manager.WorkingStack = *scope_manager.WorkingStack.UpdateID()

			//
			//
			//
			//
			scope_manager.PrintScopes()
			scope_manager.PrintWorkingStack()
			//
			//
			//
			//

			fmt.Printf("Scope, NewScopeManager(): Creation Successful.\n")
			return &scope_manager, true
		} else {
			fmt.Printf("Scope, NewScopeManager(): Failed to create inital Scope.\n")
			scope_manager = *_scope_manager
		}
	}
	// something went wrong
	fmt.Printf("Scope, NewScopeManager(): Creation: Something went wrong.\n")
	return &scope_manager, false
}

//
func (scope_manager *ScopeManager) NewDecl(node ast.Node) (*ScopeManager, DataDeclID, bool) {
	if _decl, ok := NewDataDecl(node); ok {
		_data_decl_map := (*scope_manager).Decls.Add(*_decl)
		(*scope_manager).Decls = *_data_decl_map
		return scope_manager, _decl.ID, true
	} else {
		return scope_manager, _decl.ID, false
	}
}

// Given the node of a 'scope', creates a new scope and adds it to Scopes and pdates the workign stack.
// Once the scope has been successfully created, assign its position and update its ID.
//
// In the case of the initial file scope, skips the position and ID and returns prematurely to let the Scope Manager handle the rest.
func (scope_manager *ScopeManager) NewScope(node ast.Node) (*ScopeManager, ScopeID, bool) {
	// check scopes isnt nil
	if (*scope_manager).Scopes == nil {
		fmt.Print("Scope.go: ScopeManager.NewScope(), ScopeManager.Scopes is nil.\n")
	}

	if scope, ok := NewScope(node); ok {
		// check not already found
		if _, ok := (*scope_manager).Scopes[scope.ID]; ok {
			var _temp_scope_map ScopeMap = (*scope_manager).Scopes
			var _temp_scope Scope = _temp_scope_map[scope.ID]
			var _existing_scope *Scope = &_temp_scope
			fmt.Printf("Scope.go: ScopeManager.NewScope(), Scope already found.\nShared Scope ID: \"%s\"\nNewScope: {\n%s}\nOldScope: {\n%s}\n", scope.ID, scope.ToString(), _existing_scope.ToString())
			// scope already found, check stack is correct
			if _scope_manager, ok := (*scope_manager).VerifyWorkingStack(); ok {
				scope_manager = _scope_manager
				// stack correct, set to top of stack
				if _scope_manager, ok := (*scope_manager).ReduceStackTo(scope.ID); ok {
					scope_manager = _scope_manager
				} else {
					fmt.Print("Scope.go: ScopeManager.NewScope(), ScopeManager.ReduceStackTo() failed.\n")
					scope_manager = _scope_manager
				}
			} else {
				fmt.Print("Scope.go: ScopeManager.NewScope(), ScopeManager.VerifyWorkingStack() failed.\n")
				scope_manager = _scope_manager
				return scope_manager, scope.ID, false
			}
		} else {
			// new scope, continue setting up

			// set position
			scope.Position = (*scope_manager).FileSet.Position(node.Pos())

			// check scope type
			if _scope_manager, _scope_type := (*scope_manager).GetScopeType(scope.ID); _scope_type == SCOPE_NONE {
				// not even a scope, quit
				return _scope_manager, scope.ID, false
			} else if _scope_type == SCOPE_FILE {
				// do nothing if file scope, let manager fix this
				return _scope_manager, "initial_scope", true
			} else {
				// update scope manager
				scope_manager = _scope_manager
			}

			// update id (requires position and type)
			scope = scope.UpdateID()

			// get any parameters/arguments
			scope_manager = (*scope_manager).GetScopeParams(scope.ID)

			// add scope to map
			scope_manager = (*scope_manager).AddScope(*scope)

			//
			//
			//
			//
			(*scope_manager).PrintScopes()
			(*scope_manager).PrintWorkingStack()
			//
			//
			//
			//
		}
		scope_manager = (*scope_manager).PushStack(scope.ID)
		(*scope_manager).WorkingStack = *(*scope_manager).WorkingStack.UpdateID()
		return scope_manager, scope.ID, true
	} else {
		fmt.Print("Scope.go: ScopeManager.NewScope(), NewScope() failed.\n")
		return scope_manager, scope.ID, false
	}
}

func (scope_manager *ScopeManager) AddScope(scope Scope) *ScopeManager {
	(*scope_manager).Scopes = *(*scope_manager).Scopes.Add(scope)
	return scope_manager
}

func (scope_manager *ScopeManager) PushStack(scope_id ScopeID) *ScopeManager {
	(*scope_manager).WorkingStack = *(*scope_manager).WorkingStack.Push(scope_id)
	return scope_manager
}

func (scope_manager *ScopeManager) PopStack() (*ScopeManager, ScopeID, bool) {
	_working_stack, _scope_id, _ok := (*scope_manager).WorkingStack.Pop()
	(*scope_manager).WorkingStack = *_working_stack
	return scope_manager, _scope_id, _ok
}

func (scope_manager *ScopeManager) PeekStack() (*ScopeManager, ScopeID, bool) {
	_working_stack, _scope_id, _ok := (*scope_manager).WorkingStack.Peek()
	(*scope_manager).WorkingStack = *_working_stack
	return scope_manager, _scope_id, _ok
}

func (scope_manager *ScopeManager) PrintScopes() {
	fmt.Printf("\n\nPrintScopes(): %d\n", len((*scope_manager).Scopes))
	var scope_index int = 0
	for scope_id, _ := range (*scope_manager).Scopes {
		fmt.Printf("\t%d: \"%s\"\n", scope_index, scope_id)
		scope_index++
	}
	fmt.Printf("End Scopes.\n\n")
}

func (scope_manager *ScopeManager) PrintWorkingStack() {
	fmt.Printf("\n\nPrintWorkingStack(): %d\n", len((*scope_manager).WorkingStack.ScopeIDs))
	for stack_index, scope_id := range (*scope_manager).WorkingStack.ScopeIDs {
		fmt.Printf("%d: \"%s\"\n", stack_index, scope_id)
	}
	fmt.Printf("End Stacks.\n\n")
}

func (scope_manager *ScopeManager) GetScopeType(scope_id ScopeID) (*ScopeManager, ScopeType) {
	var scope_type ScopeType
	var _scope Scope = (*scope_manager).Scopes[scope_id]
	// scopes
	switch node_type := _scope.Node.(type) {
	case *ast.File:
		scope_type = SCOPE_FILE
	case *ast.GenDecl:
		scope_type = SCOPE_GEN_DECL
	case *ast.FuncDecl:
		scope_type = SCOPE_FUNC_DECL
	case *ast.IfStmt:
		scope_type = SCOPE_IF
	case *ast.SelectStmt:
		scope_type = SCOPE_SELECT
	case *ast.SwitchStmt:
		scope_type = SCOPE_SWITCH
	case *ast.TypeSwitchStmt:
		scope_type = SCOPE_TYPE_SWITCH
	case *ast.ForStmt:
		scope_type = SCOPE_FOR
	case *ast.RangeStmt:
		scope_type = SCOPE_RANGE
	case *ast.GoStmt:
		// check if to func or anonymous
		switch _call := node_type.Call.Fun.(type) {
		case *ast.Ident:
			// named function
			switch _call.Obj.Decl.(type) {
			case *ast.FuncDecl:
				scope_type = SCOPE_GO_NAMED
			}
		case *ast.FuncLit:
			// anonymous function
			scope_type = SCOPE_GO_ANONYMOUS
		}
	case *ast.FuncLit:
		scope_type = SCOPE_ANONYMOUS_FUNCTION
	default:
		// not supported
		scope_type = SCOPE_NONE
	}
	return scope_manager, scope_type
}

// node *ast.FieldList
func (scope_manager *ScopeManager) GetScopeParams(scope_id ScopeID) *ScopeManager {
	var scope Scope = (*scope_manager).Scopes[scope_id]
	var node ast.Node = (*scope_manager).Scopes[scope_id].Node

	var param_ids []ScopeDataID = make([]ScopeDataID, 0)
	// for function parameters
	if scope.Type == SCOPE_FUNC_DECL || scope.Type == SCOPE_ANONYMOUS_FUNCTION || scope.Type == SCOPE_GO_NAMED || scope.Type == SCOPE_GO_ANONYMOUS {
		// adjust node for goroutine
		var node_params []*ast.Field
		if scope.Type == SCOPE_GO_ANONYMOUS {
			// goroutine anon func params
			node_params = node.(*ast.GoStmt).Call.Fun.(*ast.FuncLit).Type.Params.List
		} else if scope.Type == SCOPE_GO_NAMED {
			// goroutine named func params
			node_params = node.(*ast.GoStmt).Call.Fun.(*ast.Ident).Obj.Decl.(*ast.FuncDecl).Type.Params.List
		} else {
			// func params
			node_params = node.(*ast.FuncDecl).Type.Params.List
		}

		// each param
		for _, _field_node := range node_params {
			var param ScopeData
			// update id
			// _param =
			param.Node = _field_node
			param.Name = _field_node.Names[0].Obj.Name
			param.Position = (*scope_manager).FileSet.Position(_field_node.Pos())
			// because this is a param
			param.Value = ""
			param.IsParam = true
			param.IsDecl = false

			param.DataType = []DataType{}
			param.DeclID = ""

			// add to param_ids
			param_ids = append(param_ids, param.ID)
			// add to scopes contained data
			(*scope_manager).Scopes[scope_id].ContainedData[param.ID] = param
		}
	} else if true {
		// find some in select, switch, if, for etc...
		//
		//
		//
		//
		//
		//
		//
		//
		//
		//
		//
		//
		//
		//
		//
		//
		//
	}
	// update
	scope.ParamIDs = param_ids
	(*scope_manager).Scopes[scope_id] = scope
	return scope_manager
}

func (scope_manager *ScopeManager) GetContainedDataTypeParams(scope_id ScopeID, data_id ScopeDataID) (*ScopeManager, []DataType) {
	var data_type []DataType = make([]DataType, 0)
	var scope_data ScopeData = (*scope_manager).Scopes[scope_id].ContainedData[data_id]

	switch node_type := scope_data.Node.(type) {
	case *ast.ChanType:
		// found channel declaration
		// check or sync or async
		// channel, add surface level
		data_type = append(data_type, DATA_CHANNEL)
		// check for any types of types...// check for any types of types
		var keep_looking bool = true
		var current_node ast.Node = node_type
		for keep_looking {
			_scope_manager, _child, more_children := (*scope_manager).SeekDataType(current_node)
			scope_manager = _scope_manager

			data_type = append(data_type)

			if more_children {
				current_node = _child
			} else {
				keep_looking = false
			}
		}
	case *ast.Ident:
		// simple, no more reccursion
		switch node_type.Name {
		case "int":
			data_type = append(data_type, DATA_INT)
		case "string":
			data_type = append(data_type, DATA_STRING)
		default:
			// something else
			data_type = append(data_type, DATA_VARIABLE)
		}
	}

	return scope_manager, data_type
}

// ast.field, ast.assignstmt, ast.valuespec
func (scope_manager *ScopeManager) GetContainedDataType(scope_id ScopeID, data_id ScopeDataID) (*ScopeManager, []DataType) {
	var data_type []DataType = make([]DataType, 0)
	var scope_data ScopeData = (*scope_manager).Scopes[scope_id].ContainedData[data_id]

	if scope_data.IsDecl {

		switch node_type := scope_data.Node.(type) {
		case *ast.CallExpr:
			// declaration params
			if len(node_type.Args) == 1 {
				data_type = append(data_type, DATA_SYNC_CHANNEL)
			} else if len(node_type.Args) == 2 {
				data_type = append(data_type, DATA_ASYNC_CHANNEL)
			}
			_scope_manager, _data_type := scope_manager.GetContainedDataTypeParams(scope_id, data_id)
			scope_manager = _scope_manager
			data_type = append(data_type, _data_type...)
		}

	} else if scope_data.IsParam {

		_scope_manager, _data_type := scope_manager.GetContainedDataTypeParams(scope_id, data_id)
		scope_manager = _scope_manager
		data_type = append(data_type, _data_type...)

	} else {

	}

	// update manager
	scope_data.DataType = data_type
	(*scope_manager).Scopes[scope_id].ContainedData[data_id] = scope_data
	return scope_manager, data_type
}

func (scope_manager *ScopeManager) SeekDataType(node ast.Node) (*ScopeManager, ast.Node, bool) {

	return scope_manager, node, false
}

func (scope_manager *ScopeManager) GetDataType() {

}

// given a node, adds it to the relevant fields in the manager
func (scope_manager *ScopeManager) ParseNode(node ast.Node) (*ScopeManager, ParseType, string, bool) {

	// check not parsing node twice
	if node.Pos() == (*scope_manager).WorkingStack.StartPos && node.End() == (*scope_manager).WorkingStack.EndPos {
		return scope_manager, PARSE_DUPL, "", false
	}

	// check if leaving current scope
	if node.End() == (*scope_manager).WorkingStack.EndPos {
		if _scope_manager, scope_id, ok := (*scope_manager).PopStack(); ok {
			scope_manager = _scope_manager
			return scope_manager, PARSE_SCOPE_END, string(scope_id), true
		} else {
			scope_manager = _scope_manager
			return scope_manager, PARSE_SCOPE_END, string(scope_id), false
		}
	}

	// must be something to parse
	switch decl := node.(type) {

	case *ast.GenDecl:
		// variable declaration
		switch decl.Tok {
		case token.VAR:
			// global variable block
		case token.CONST:
			// constant block
		case token.IMPORT:
			// import block
		default:
			// unaccounted for
			return scope_manager, PARSE_DECL, "", false
		}

	case *ast.FuncLit:
		// goroutine function
		if _scope_manager, parent_scope_id, ok := (*scope_manager).PeekStack(); ok {
			scope_manager = _scope_manager
			// check current scope is go type
			var parent_scope_type ScopeType = (*scope_manager).Scopes[parent_scope_id].Type
			if parent_scope_type == SCOPE_GO_ANONYMOUS || parent_scope_type == SCOPE_GO_NAMED {
				// func lit is inside go stmt
				// make new scope
				if _scope_manager, scope_id, ok := (*scope_manager).NewScope(node); ok {
					// scope creation success
					return _scope_manager, PARSE_SCOPE, string(scope_id), true
				} else {
					// new scope failed
					return _scope_manager, PARSE_SCOPE, string(scope_id), false
				}
			} else {
				// cannot accept ast.funclit not inside a go statment
				return scope_manager, PARSE_SCOPE, string(parent_scope_id), false
			}
		} else {
			return _scope_manager, PARSE_STACK_EMPTY, string(parent_scope_id), false
		}

	case *ast.FuncDecl:
		// function declaration
		// decl.Body.l

	case *ast.AssignStmt:
		// assignment, check decl
		switch decl.Tok {
		case token.DEFINE:
			// variable declaration
		case token.ASSIGN:
			// updating existing variable
		default:
			// unaccounted for
			return scope_manager, PARSE_ASSIGN, "", false
		}
	}
	// not something interesting
	return scope_manager, PARSE_NONE, "", false
}

// func GetSpecDataTypes(node []ast.Spec) {

// }

func GetSpecDataType(node ast.ValueSpec) *[]DataType {
	var data_type []DataType

	// for _, name := range node.Names {
	// 	// name.Name
	// }

	return &data_type
}

func (scope_manager *ScopeManager) ReduceStackTo(scope_id ScopeID) (*ScopeManager, bool) {
	if _scope_manager, ok := (*scope_manager).WorkingStackContains(scope_id); ok {
		scope_manager = _scope_manager
		for _scope_manager, _scope_id, ok := (*scope_manager).PeekStack(); ok; {
			scope_manager = _scope_manager
			// top of stack is desired
			if _scope_id == scope_id {
				return scope_manager, true
			}
			// continue
			(*scope_manager).PopStack()
		}
	} else {
		scope_manager = _scope_manager
	}
	return scope_manager, false
}

func (scope_manager *ScopeManager) WorkingStackContains(scope_id ScopeID) (*ScopeManager, bool) {
	_stack, ok := (*scope_manager).WorkingStack.ScopeIDStackContains(scope_id)
	scope_manager.WorkingStack = *_stack
	return scope_manager, ok
}

func (scope_manager *ScopeManager) VerifyWorkingStack() (*ScopeManager, bool) {
	// check that each scope in stack has next in nested id
	for i, scope_id := range (*scope_manager).WorkingStack.ScopeIDs {
		// scope first one
		if i > 0 {
			var prev_scope_id ScopeID = (*scope_manager).WorkingStack.ScopeIDs[i-1]
			var prev_scope Scope = (*scope_manager).Scopes[prev_scope_id]
			if _prev_scope, ok := prev_scope.ChildScopeIDsContains(scope_id); ok {
				(*scope_manager).Scopes[prev_scope_id] = *_prev_scope
				continue
			} else {
				(*scope_manager).Scopes[prev_scope_id] = *_prev_scope
				return scope_manager, false
			}
		}
	}
	return scope_manager, true
}

func (scope *Scope) ChildScopeIDsContains(scope_id ScopeID) (*Scope, bool) {
	for _, _id := range (*scope).ChildScopeIDs {
		if _id == scope_id {
			return scope, true
		}
	}
	return scope, false
}

func (scope_id_stack *ScopeIDStack) ScopeIDStackContains(scope_id ScopeID) (*ScopeIDStack, bool) {
	for _, _id := range (*scope_id_stack).ScopeIDs {
		if _id == scope_id {
			return scope_id_stack, true
		}
	}
	return scope_id_stack, false
}

// // //
// // // Scope Stack
// // //
func NewWorkingScopeStack(_initial_scope_id ScopeID, start_pos token.Pos, end_pos token.Pos) ScopeIDStack {
	var scope_id_stack ScopeIDStack
	scope_id_stack.Size = 1
	// add scope id slice to stack
	scope_id_stack.ScopeIDs = []ScopeID{_initial_scope_id}
	// generate id string
	scope_id_stack = *scope_id_stack.UpdateID()

	scope_id_stack.StartPos = start_pos
	scope_id_stack.EndPos = end_pos

	return scope_id_stack
}

// a trace of the stack using the ids of the scopes contained.
func (scope_id_stack *ScopeIDStack) UpdateID() *ScopeIDStack {
	var _scope_stack_id ScopeIDStackID
	for _, _scope_id := range (*scope_id_stack).ScopeIDs {
		_scope_stack_id = ScopeIDStackID(fmt.Sprintf("%sL>- ScopeID: %s\n", _scope_stack_id, _scope_id))
	}
	(*scope_id_stack).ID = _scope_stack_id
	return scope_id_stack
}

func (scope_id_stack *ScopeIDStack) Push(scope_id ScopeID) *ScopeIDStack {
	(*scope_id_stack).Size++
	(*scope_id_stack).ScopeIDs = append((*scope_id_stack).ScopeIDs, scope_id)
	return scope_id_stack
}

// removes top of stack, returns new stack, element removed and if it was successful
func (scope_id_stack *ScopeIDStack) Pop() (*ScopeIDStack, ScopeID, bool) {
	// sanity check
	if (*scope_id_stack).Size != len((*scope_id_stack).ScopeIDs) {
		fmt.Printf("Scope.go: ScopeIDStack.Pop(), Size and Len don't add up...\n\tSize: %d\n\tLen: %d\n", (*scope_id_stack).Size, len((*scope_id_stack).ScopeIDs))
		return scope_id_stack, ScopeID(""), false
	}

	if (*scope_id_stack).Size > 0 {
		// get scope id to return
		var _popped_scope_id ScopeID = (*scope_id_stack).ScopeIDs[(*scope_id_stack).Size-1]
		// return reduced scope stack
		(*scope_id_stack).ScopeIDs = (*scope_id_stack).ScopeIDs[:(*scope_id_stack).Size-1]
		return scope_id_stack, _popped_scope_id, true
	}
	// empty
	fmt.Printf("Scope.go: ScopeIDStack.Pop(), ScopeIDStack.Scopes is empty.\n")
	return scope_id_stack, ScopeID(""), false
}

func (scope_id_stack *ScopeIDStack) Peek() (*ScopeIDStack, ScopeID, bool) {
	if (*scope_id_stack).Size > 0 {
		return scope_id_stack, scope_id_stack.ScopeIDs[(*scope_id_stack).Size-1], true
	}
	fmt.Printf("Scope.go: ScopeIDStack.Peek(), ScopeIDStack.Scopes is empty.\n")
	return scope_id_stack, ScopeID(""), false
}

func (scope_id_stack *ScopeIDStack) ToString() string {
	var print_string string
	for i, v := range (*scope_id_stack).ScopeIDs {
		print_string = fmt.Sprintf("%s\t%d: %s\n", print_string, i, v)
	}
	return print_string
}

// // //
// // // Scope Map
// // //
func NewScopeMap() *ScopeMap {
	var scope_map *ScopeMap = &ScopeMap{}
	return scope_map
}

func (scope_map *ScopeMap) Add(_scope Scope) *ScopeMap {
	(*scope_map)[_scope.ID] = _scope
	return scope_map
}

func (scope_map *ScopeMap) ToString() string {
	var print_string string
	var scope_index int = 0
	for i, v := range *scope_map {
		print_string = fmt.Sprintf("%sScope %d, \"%s\"\n%v\n", print_string, scope_index, i, v.ToString())
		scope_index++
	}
	return print_string
}

// // //
// // // Scope
// // //
// given a node, tries to make it a scope
func NewScope(node ast.Node) (*Scope, bool) {
	var scope Scope = Scope{}
	scope.Node = node
	scope = *scope.NewScopeDataMap()
	scope.ContainedDataDecl = make([]DataDeclID, 0)
	scope.ChildScopeIDs = make([]ScopeID, 0)

	return &scope, true
}

func (scope *Scope) ToString() string {
	var print_string string
	// id, type, posiiton
	print_string = fmt.Sprintf("\tID: \"%s\"\n\tType: %s\n\tPosition:\n%s\tParams: %d\n", (*scope).ID, (*scope).Type, (*scope).Position.String(), len((*scope).ParamIDs))
	// params
	for i, param_id := range (*scope).ParamIDs {
		print_string = fmt.Sprintf("%s\t\tParam %d ID: %s\n", print_string, i, param_id)
	}
	// contained data
	print_string = fmt.Sprintf("%s\tContainedData: %d\n%s", print_string, len((*scope).ContainedData), (*scope).ContainedData.ToString())
	// contained data decl
	print_string = fmt.Sprintf("%s\tContainedDataDecl: %d\n", print_string, len((*scope).ContainedData))
	for i, v := range (*scope).ContainedDataDecl {
		print_string = fmt.Sprintf("%s\t%d: %s\n", print_string, i, v)
	}
	// child scope ids
	print_string = fmt.Sprintf("%s\tChildScopeIds: %d\n", print_string, len((*scope).ChildScopeIDs))
	for i, v := range (*scope).ChildScopeIDs {
		print_string = fmt.Sprintf("%s\t%d: %s\n", print_string, i, v)
	}
	// scope id stack
	print_string = fmt.Sprintf("%s\tScopeIDStack: %d\n%s", print_string, (*scope).ScopeIDStack.Size, (*scope).ScopeIDStack.ToString())
	// print
	return print_string
}

// called immediately after NewScope by manager
func (scope *Scope) SetStack(_scope_id_stack ScopeIDStack) *Scope {
	(*scope).ScopeIDStack = _scope_id_stack
	return scope
}

func (scope *Scope) UpdateID() *Scope {
	// TODO: use Position and ScopeType
	(*scope).ID = ScopeID(fmt.Sprintf("%s:%d:%d:%s", (*scope).Position.Filename, (*scope).Position.Line, (*scope).Position.Offset, (*scope).Type))
	return scope
}

func (scope *Scope) AddChildID(_child_id ScopeID) *Scope {
	(*scope).ChildScopeIDs = append((*scope).ChildScopeIDs, _child_id)
	return scope
}

// // //
// // // Scope Data Map
// // //
func (scope *Scope) NewScopeDataMap() *Scope {
	var scope_data_map ScopeDataMap
	(*scope).ContainedData = scope_data_map
	return scope
}

func (scope_data_map *ScopeDataMap) ToString() string {
	var print_string string = ""
	for i, v := range *scope_data_map {
		print_string = fmt.Sprintf("%s\t\t%s: %v\n", print_string, i, v)
	}
	return print_string
}

// // //
// // // Scope Data
// // // ast.field, ast.assignstmt, ast.valuespec
func NewScopeData(node ast.Node, position token.Position) ScopeData {
	var scope_data ScopeData = ScopeData{}
	scope_data.Node = node

	switch _node_type := node.(type) {
	case *ast.Field:
		// from params
		scope_data.IsParam = true
		scope_data.IsDecl = false
		// get name
		var names string
		for _, name := range _node_type.Names {
			names = fmt.Sprintf("%s,%s", names, name)
		}
		scope_data.Name = names
		// get datatype
		// scope_data.DataType = GetDataType(node)

	case *ast.AssignStmt:
		// declaration
		// _scope_data.Decl = NewDataDecl()
		scope_data.IsParam = false
		scope_data.IsDecl = true
		// get name
		for _, _node_lhs := range _node_type.Lhs {
			switch _lhs := _node_lhs.(type) {
			case *ast.Ident:
				scope_data.Name = _lhs.Name
			}
		}
		// get datatype
		// for _, _node_rhs := range _node_type.Rhs {
		// 	switch _rhs := _node_rhs.(type) {
		// 	case *ast.CallExpr:
		// 		// scope_data.DataType = GetDataType(_rhs)
		// 	}
		// }
		// scope_data.DataType = GetDataType(node)

	case *ast.ValueSpec:
		// generic declaration
		// _scope_data.Decl = NewDataDecl()
		scope_data.IsParam = false
		scope_data.IsDecl = true
		// get name
		var names string
		for _, name := range _node_type.Names {
			names = fmt.Sprintf("%s,%s", names, name)
		}
		scope_data.Name = names
		// get datatype
		// scope_data.DataType = GetDataType(node)
	}
	// add others
	scope_data.Position = position
	return scope_data
}

func (scope_data *ScopeData) ToString() string {
	var print_string string
	// each scope data
	print_string = fmt.Sprintf("\t\tID: %s\n\t\tName: %s\n\t\tPosition: %s\n\t\tValue: %s\n\t\tIsParam: %t\n\t\tIsDecl: %t\n\t\tDeclID: %s\nDataType: %d\n", (*scope_data).ID, (*scope_data).Name, (*scope_data).Position.String(), (*scope_data).Value, (*scope_data).IsParam, (*scope_data).IsDecl, (*scope_data).DeclID, len((*scope_data).DataType))
	for d_i, d_v := range (*scope_data).DataType {
		print_string = fmt.Sprintf("\t\t%d: %s\n", d_i, d_v)
	}
	return print_string
}

// // //
// // // Data Decl Map
// // //
func NewDataDeclMap() *DataDeclMap {
	var data_decl_map *DataDeclMap = &DataDeclMap{}
	return data_decl_map
}

func (data_decl_map *DataDeclMap) Add(_data_decl DataDecl) *DataDeclMap {
	(*data_decl_map)[_data_decl.ID] = _data_decl
	return data_decl_map
}

// // //
// // // Decl Data
// // //
func NewDataDecl(node ast.Node) (*DataDecl, bool) {
	var data_decl *DataDecl = &DataDecl{}

	return data_decl, true
}
