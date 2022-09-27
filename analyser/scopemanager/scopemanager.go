package scopemanager

import (
	"fmt"
	"go/ast"
	"go/token"
	"os"
)

// ScopeManager
//
type ScopeManager struct {
	ScopeMap *MapOfScopes
	Stack    *StackOfIDs
	FileSet  *token.FileSet
	FileSrc  string
	Decls    *MapOfDecls
	// Cache    struct {
	// 	AwaitedFunction   AwaitedFunction
	// 	ExpectingFunction bool
	// }
}

//
func NewScopeManager(filename string, fileset *token.FileSet) (*ScopeManager, error) {
	var sm ScopeManager

	sm.ScopeMap = NewMapOfScopes()
	sm.Stack = NewStackOfIDs()
	sm.FileSet = fileset

	if file_src, err := os.ReadFile(filename); err == nil {
		sm.FileSrc = string(file_src)
	} else {
		// failed
		return &sm, err
	}

	sm.Decls = NewMapOfDecls()

	return &sm, nil
}

// AwaitedFunction
//
type AwaitedFunction struct {
	Name     string
	Pos      token.Pos
	Args     *[]ast.Expr
	ParentID ID
}

func NewAwaitedFunction(node *ast.Node, id ID) AwaitedFunction {
	var awaited_function AwaitedFunction

	awaited_function.Name = (*node).(*ast.CallExpr).Fun.(*ast.Ident).Name
	awaited_function.Pos = (*node).Pos()
	awaited_function.Args = &(*node).(*ast.CallExpr).Args
	awaited_function.ParentID = id

	return awaited_function
}

func (sm *ScopeManager) CheckAwaitedFunction(node *ast.Node) (*ScopeManager, bool) {

	// x:=(*node).(*ast.Ident).

	return sm, false
}

// ParseType
//
type ParseType string

const (
	// general
	PARSE_NONE       ParseType = "None"
	PARSE_SCOPE_EXIT ParseType = "None"

	// scopes
	PARSE_PACKAGE          ParseType = "None"
	PARSE_PACKAGE_CONST    ParseType = "None"
	PARSE_PACKAGE_VAR      ParseType = "None"
	PARSE_FILE             ParseType = "None"
	PARSE_FILE_IMPORT      ParseType = "None"
	PARSE_FUNC_DECL        ParseType = "None"
	PARSE_FUNC_LIT         ParseType = "None"
	PARSE_GO_STMT          ParseType = "None"
	PARSE_FOR_STMT         ParseType = "None"
	PARSE_RANGE_STMT       ParseType = "None"
	PARSE_IF_STMT          ParseType = "None"
	PARSE_SELECT_STMT      ParseType = "None"
	PARSE_SWTICH_STMT      ParseType = "None"
	PARSE_TYPE_SWITCH_STMT ParseType = "None"

	// vars
	PARSE_ASSIGN ParseType = "None"
	PARSE_DECL   ParseType = "None"
	PARSE_DEFINE ParseType = "None"
	// PARSE_FUNC_CALL_PARAMS ParseType = "None"
	PARSE_FUNC_DECL_PARAMS ParseType = "None"

	// fails
	PARSE_FAIL              ParseType = "None"
	PARSE_FAIL_ASSIGN_TOKEN ParseType = "None"
	PARSE_FAIL_DECL_TOKEN   ParseType = "None"
	PARSE_FAIL_DEFAULT      ParseType = "None"
	PARSE_FAIL_FIELD_LIST   ParseType = "None"
	PARSE_FAIL_GEN_DECL     ParseType = "None"
	PARSE_FAIL_STACK_PEEK   ParseType = "None"
	PARSE_FAIL_STACK_POP    ParseType = "None"
	PARSE_FAIL_VALUE_SPEC   ParseType = "None"
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
	if scope, ok := (*sm).Peek(); ok {
		if (*node).Pos() > (*scope.Node).End() {
			// if current node starts after the current scope ends, left current scope
			if _sm, ok := (*sm).Pop(); ok {
				sm = _sm
			} else {
				// failed
				return sm, PARSE_FAIL_STACK_POP
			}
			return sm, PARSE_SCOPE_EXIT
		} // continue
	} else {
		return sm, PARSE_FAIL_STACK_PEEK
	}

	// Check for each ScopeType
	switch node_type := (*node).(type) {

	// Scope: Package
	case *ast.Package:
		sm = (*sm).NewScope(node, SCOPE_TYPE_PACKAGE)
		return sm, PARSE_PACKAGE

	// Scope: File
	case *ast.File:
		sm = (*sm).NewScope(node, SCOPE_TYPE_FILE)
		return sm, PARSE_FILE

	// Scope: GenDecl
	case *ast.GenDecl:
		// scope or vardecl, depends on outerscope
		if outer_scope, ok := (*sm).Peek(); ok {
			// if file, this is new scope of global decl
			if outer_scope.Type == SCOPE_TYPE_FILE {
				// global decls: const import var
				switch node_type.Tok {
				case token.VAR:
					sm = (*sm).NewScope(node, SCOPE_TYPE_PACKAGE_VAR)
					return sm, PARSE_PACKAGE_VAR

				case token.CONST:
					sm = (*sm).NewScope(node, SCOPE_TYPE_PACKAGE_CONST)
					return sm, PARSE_PACKAGE_CONST

				case token.IMPORT:
					sm = (*sm).NewScope(node, SCOPE_TYPE_FILE_IMPORT)
					return sm, PARSE_FILE_IMPORT

				default:
					return sm, PARSE_FAIL_DECL_TOKEN
				}
			} else {
				return sm, PARSE_FAIL_GEN_DECL
			}
		} else {
			return sm, PARSE_FAIL_STACK_PEEK
		}

	// Scope: Goroutine
	case *ast.GoStmt:
		sm = (*sm).NewScope(node, SCOPE_TYPE_GOROUTINE)
		return sm, PARSE_GO_STMT

	// Scope: Anon Function
	case *ast.FuncLit:
		sm = (*sm).NewScope(node, SCOPE_TYPE_FUNC_DECL)
		return sm, PARSE_FUNC_LIT

	// Scope: Function
	case *ast.FuncDecl: // line 1914
		sm = (*sm).NewScope(node, SCOPE_TYPE_FUNC_DECL)
		return sm, PARSE_FUNC_DECL

	// Scope: FuncCall
	case *ast.CallExpr:
		// check function call
		if outer_scope, ok := (*sm).Peek(); ok {
			if outer_scope.Type == SCOPE_TYPE_GOROUTINE {
				sm = (*sm).NewScope(node, SCOPE_TYPE_FUNC_CALL)
				return sm, PARSE_FUNC_DECL
			} else {
				return sm, PARSE_NONE
			}
		} else {
			return sm, PARSE_FAIL_STACK_PEEK
		}

	// Scope: If Statement
	case *ast.IfStmt:
		sm = (*sm).NewScope(node, SCOPE_TYPE_IF)
		return sm, PARSE_IF_STMT

	// Scope: If Statement
	case *ast.SelectStmt:
		sm = (*sm).NewScope(node, SCOPE_TYPE_SELECT)
		return sm, PARSE_SELECT_STMT

	// Scope: Switch Statement
	case *ast.SwitchStmt:
		sm = (*sm).NewScope(node, SCOPE_TYPE_SWITCH)
		return sm, PARSE_SWTICH_STMT

	// Scope: Switch Type Statement
	case *ast.TypeSwitchStmt:
		sm = (*sm).NewScope(node, SCOPE_TYPE_TYPE_SWITCH)
		return sm, PARSE_TYPE_SWITCH_STMT

	// Scope: For Loop Statement
	case *ast.ForStmt:
		sm = (*sm).NewScope(node, SCOPE_TYPE_FOR)
		return sm, PARSE_FOR_STMT

	// Scope: Ranged For Loop Statement
	case *ast.RangeStmt:
		sm = (*sm).NewScope(node, SCOPE_TYPE_RANGE)
		return sm, PARSE_RANGE_STMT

	// Var: Params
	case *ast.FieldList: // line 1924
		// scope or vardecl, depends on outerscope
		if outer_scopes, ok := (*sm).PeekX(2); ok {

			var func_call Scope = *outer_scopes[0]
			var func_decl Scope = *outer_scopes[1]

			// if these are function parameters
			if func_call.Type == SCOPE_TYPE_FUNC_CALL && node_type.List != nil {
				// if file, this is new scope of global decl
				switch func_decl.Type {

				// params line 1385
				// args line 1843

				// Func Decl Parameters
				case SCOPE_TYPE_FUNC_DECL:

					// extract each param as decl, take values from passed args
					// TODO 
					// for _index, _param := range node_type.List {
					// 	// new decl

					// 	// find arg from parent
					// 	// node_type
					// }

					return sm, PARSE_FUNC_DECL_PARAMS

				// not accounted for
				default:
					return sm, PARSE_FAIL_FIELD_LIST
				}
			} else {
				return sm, PARSE_NONE
			}
		} else {
			return sm, PARSE_FAIL_STACK_PEEK
		}

	// Var: VarDecl, global or scoped
	case *ast.ValueSpec:
		// check outerscopes context
		if outer_scopes, ok := (*sm).PeekX(2); ok {

			// TODO 
			// var file_scope Scope = *outer_scopes[0]
			// var Scope = *outer_scopes[1]

			// if file > gendecl > node
			if outer_scopes[0].Type == SCOPE_TYPE_FILE {

				// package decl or file import
				switch outer_scopes[1].Type {

				case SCOPE_TYPE_PACKAGE_VAR:
					if _sm, ok := (*sm).NewVarDecl(node, token.VAR); ok {
						sm = _sm
					}
					return sm, PARSE_PACKAGE_VAR

				case SCOPE_TYPE_PACKAGE_CONST:
					if _sm, ok := (*sm).NewVarDecl(node, token.CONST); ok {
						sm = _sm
					}
					return sm, PARSE_PACKAGE_CONST

				case SCOPE_TYPE_FILE_IMPORT:
					if _sm, ok := (*sm).NewVarDecl(node, token.IMPORT); ok {
						sm = _sm
					}
					return sm, PARSE_FILE_IMPORT
				// not accounted for
				default:
					return sm, PARSE_FAIL_VALUE_SPEC
				}
			} else {
				// scoped decl
				if _sm, ok := (*sm).NewVarDecl(node, token.VAR); ok {
					sm = _sm
				}
				return sm, PARSE_DECL
			}
		} else {
			return sm, PARSE_FAIL_STACK_PEEK
		}

	// Var: Assign or Decl
	case *ast.AssignStmt:
		switch node_type.Tok {

		// VarDecl: Define (:=)
		case token.DEFINE:
			if _sm, ok := (*sm).NewVarDecl(node, token.DEFINE); ok {
				sm = _sm
			}
			return sm, PARSE_DEFINE

		// VarData: Assignment
		case token.ASSIGN:
			// find label decl
			switch var_label_type := node_type.Lhs[0].(type) {
			case *ast.Ident:
				var decl_id ID = (*sm).FindDeclID(var_label_type.Name)
				// update values
				sm = (*sm).VarDeclAddValue(decl_id, var_label_type, node_type.Pos())

			}
			return sm, PARSE_ASSIGN

		// unnaccounted for
		default:
			return sm, PARSE_FAIL_ASSIGN_TOKEN
		}

	// nothing of interest
	default:
		return sm, PARSE_NONE
	}
}

// Returns the first Scope found that contains a VarDecl of the Label provided.
// Goes through the Stack from top to bottom.
func (sm *ScopeManager) FindDeclID(label string) ID {

	for i := 0; i < (*sm).Stack.Size(); i++ {
		// from top of stack
		if scope_id, ok := (*sm).Stack.Get((*sm).Stack.Size() - i); ok {
			// get possible decl id
			var decl_id ID = NewVarDeclID(label, scope_id)
			// check if in this scope
			if _, ok := (*(*sm.ScopeMap)[scope_id].Decls)[decl_id]; ok {
				return decl_id
			}
		}
	}
	// not found
	return ID("")
}

// Adds a given value to a VarDecls slice of Values
func (sm *ScopeManager) VarDeclAddValue(decl_id ID, expr ast.Expr, pos token.Pos) *ScopeManager {
	var value VarValue = (*sm).NewVarValue(expr, pos)
	(*sm.Decls)[decl_id] = (*sm.Decls)[decl_id].AddValue(value)
	return sm
}

//
func (decl *VarDecl) AddValue(value VarValue) *VarDecl {
	decl.Values = append(decl.Values, value)
	return decl
}

// Creates a new VarDecl and adds it the the MapOfVarDecl
// Node should be of type *ast.ValueSpec or *ast.AssignStmt
func (sm *ScopeManager) NewVarDecl(node *ast.Node, tok token.Token) (*ScopeManager, bool) {

	switch node_type := (*node).(type) {
	case *ast.ValueSpec:
		// decl found
		var var_decl *VarDecl

		var_decl.Label = node_type.Names[0].Name
		var_decl.Node = node
		var_decl.Token = tok

		var_decl.Type = (*sm).NewVarType(node)

		// check for value
		if node_type.Values != nil {
			for _, value_expr := range node_type.Values {
				// add to values
				var_decl = var_decl.AddValue((*sm).NewVarValue(value_expr, node_type.Pos()))
			}

		}

		// add to ScopeManager
		(*sm.Decls)[(*sm).NewVarDeclID(var_decl)] = var_decl

	case *ast.AssignStmt:
		// variable assignment, find decl
		switch node_type.Tok {
		case token.DEFINE:
			// check decl
			// for each decl
			for index, expr := range node_type.Lhs {
				// ensure ident
				switch expr_ident := expr.(type) {
				case *ast.Ident:
					// declaration
					var var_decl *VarDecl

					var_decl.Label = expr_ident.Name
					var_decl.Node = node
					var_decl.Type = (*sm).NewVarType(node)
					var_decl.Token = token.DEFINE

					// add to values
					var_decl = var_decl.AddValue((*sm).NewVarValue(node_type.Rhs[index], node_type.Pos()))

					// add to ScopeManager
					(*sm.Decls)[(*sm).NewVarDeclID(var_decl)] = var_decl
				}
			}
		default:
			// unnaccounted for
			return sm, false
		}
	}

	return sm, true
}

// Returns VarValue using ast.ValueSpec .Values[]ast.Expr and .Pos
func (sm *ScopeManager) NewVarValue(expr ast.Expr, pos token.Pos) VarValue {
	var value VarValue

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
		var expr_str string
		switch value_expr.(type) {
		case *ast.CallExpr:
			expr_str = "CallExpr"
		case *ast.BinaryExpr:
			expr_str = "BinaryExpr"
		case *ast.UnaryExpr:
			expr_str = "UnaryExpr"
		default:
			expr_str = "Other"
		}
		value.Value = fmt.Sprintf("%v", expr_str)

	}

	return value
}

// VarType
// DataType is a list of Types
// Argument contains specific arguments like "BufferSize" for channels.
// Use ParseInt etc for extracting values from resulting string.
type VarType struct {
	Type GeneralVarType
	Data []string
	Info map[string]string
}

// Retruns VarType when node is:
// - *ast.AssignStmt
// - *ast.ValueSpec
// - *ast.FieldList
func (sm *ScopeManager) NewVarType(node *ast.Node) VarType {

	var var_type VarType
	var_type.Data = make([]string, 0)

	switch node_type := (*node).(type) {

	// Define (:=) assignment
	case *ast.AssignStmt:
		switch node_type.Tok {
		case token.DEFINE:
			// take type from rhs
			switch rhs_expr := node_type.Rhs[0].(type) {

			// int or string
			case *ast.BasicLit:
				switch rhs_expr.Kind {
				// int
				case token.INT:
					var_type.Type = VAR_DATA_TYPE_INT
				// string
				case token.STRING:
					var_type.Type = VAR_DATA_TYPE_STRING
				// struct
				case token.STRUCT:
					var_type.Type = VAR_DATA_TYPE_STRUCT
				// unnaccounted for
				default:
					var_type.Type = VAR_DATA_TYPE_OTHER
				}

			// Data received from channel
			case *ast.UnaryExpr:
				// received from channel
				if rhs_expr.Op == token.ARROW {
					var channel_name string = rhs_expr.X.(*ast.Ident).Name
					// search outwardly for first decl of this label
					if channel_decl_id := (*sm).FindDeclID(channel_name); channel_decl_id != "" {
						// copy
						copy(var_type.Data, (*sm.Decls)[channel_decl_id].Type.Data[1:])
					}
				}

			// Type from Function
			case *ast.CallExpr:
				switch rhs_expr.Fun.(*ast.Ident).Name {
				// Channel or Slice
				case "make":
					switch rhs_expr.Args[0].(type) {

					// Channel
					case *ast.ChanType:
						// If Async. Channel
						if len(rhs_expr.Args) > 1 {
							var_type.Type = VAR_DATA_TYPE_ASYNC_CHANNEL
							// get buffer size
							switch rhs_expr.Args[1].(type) {
							// Buffer inline
							case *ast.BasicLit:
								var_type.Info["BufferSize"] = fmt.Sprintf("%v", rhs_expr.Args[1].(*ast.BasicLit).Value)

							// Buffer from var
							case *ast.Ident:
								// search outwardly for first decl of this label
								if var_decl_id := (*sm).FindDeclID(rhs_expr.Args[1].(*ast.Ident).Name); var_decl_id != "" {
									// get data type
									var_type.Info["BufferSize"] = fmt.Sprintf("%v", rhs_expr.Args[1].(*ast.BasicLit).Value)
									// var_type.Data = (*sm.Decls)[channel_decl_id].Data
								}

							default:
								var_type.Type = VAR_DATA_TYPE_OTHER
							}
						} else {
							// sync channel
							var_type.Type = VAR_DATA_TYPE_SYNC_CHANNEL
						}
						// get channel type
						var iterate_type ast.Node = rhs_expr.Args[0]
						var loop_type bool = true
						for loop_type {
							// current data is chan
							switch iter_type := iterate_type.(type) {
							case *ast.ChanType:
								var_type.Data = append(var_type.Data, "chan")
								// get next type
								iterate_type = iter_type.Value
							case *ast.Ident:
								var_type.Data = append(var_type.Data, iter_type.Name)
								// exit loop
								loop_type = false
							default:
								// exit loop if not accounted for
								loop_type = false
							}
						}
					default:
						var_type.Type = VAR_DATA_TYPE_OTHER
					}
				}

			// some other func
			case *ast.CompositeLit:
				switch rhs_expr.Type.(type) {
				// Get them all
				case *ast.SelectorExpr:

					var sel_expr []string = ExtractExpr(rhs_expr.Type)
					// add to type
					var_type.Data = sel_expr
					var_type.Type = VAR_DATA_FUNC_RET

					// context
					var_type.Info["Function"] = sel_expr[len(sel_expr)-1]

				default:
					var_type.Type = VAR_DATA_TYPE_OTHER
				}
			default:
				var_type.Type = VAR_DATA_TYPE_OTHER
			}
		}

	// Params
	case *ast.Field:

	// Declaration
	case *ast.ValueSpec:
		// look in type field
		if node_type.Type != nil {
			switch value_type := node_type.Type.(type) {

			// Pointer
			case *ast.Ident:
				switch value_type.Name {
				// int
				case "int":
					var_type.Type = VAR_DATA_TYPE_INT
				// string
				case "string":
					var_type.Type = VAR_DATA_TYPE_STRING
				// unnaccounted for
				default:
					var_type.Type = VAR_DATA_TYPE_OTHER
				}
			// unnaccounted for
			default:
				var_type.Type = VAR_DATA_TYPE_OTHER
			}
		}

		// get type from value if possible and couldnt get it from type field
		if node_type.Values != nil || var_type.Type == VAR_DATA_TYPE_OTHER {
			switch value_value := node_type.Type.(type) {

			// Pointer
			case *ast.StarExpr:

			// int or string
			case *ast.BasicLit:
				switch value_value.Kind {
				// int
				case token.INT:
					var_type.Type = VAR_DATA_TYPE_INT
				// string
				case token.STRING:
					var_type.Type = VAR_DATA_TYPE_STRING
				// struct
				case token.STRUCT:
					var_type.Type = VAR_DATA_TYPE_STRUCT
				// unnaccounted for
				default:
					var_type.Type = VAR_DATA_TYPE_OTHER
				}
			}
		}
	// unnaccounted for
	default:
		var_type.Type = VAR_DATA_TYPE_OTHER
	}

	return var_type
}

// Returns []string containing selectorexpor x, sel of compositelit in each element
// ast.Expr should be of type *ast.SelectorExpr
func ExtractExpr(current_sel_expr ast.Expr) []string {

	var sel_expr []string = make([]string, 0)
	var loop bool = true
	for loop {
		// For recursion on X,
		switch outer_sel_type := current_sel_expr.(type) {
		// only loops whilst SelectorExpr
		case *ast.SelectorExpr:
			// still more to go, add sel to beginning of slice
			sel_expr = append([]string{outer_sel_type.Sel.Name}, sel_expr...)

			// extracting from x
			switch inner_sel_type := outer_sel_type.X.(type) {

			// Selector
			case *ast.SelectorExpr:
				// make x selector
				current_sel_expr = inner_sel_type

			// Ident
			case *ast.Ident:
				// add x to beginning
				sel_expr = append([]string{inner_sel_type.Name}, sel_expr...)
				// end loop
				loop = false
			}

		default:
			loop = false
		}
	}

	return sel_expr
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

// Returns X amount of Scopes at the top of the Stack, and bool if successful
func (sm *ScopeManager) PeekX(x int) ([]*Scope, bool) {
	if scope_ids, ok := (*sm).Stack.PeekX(x); ok {
		var scopes []*Scope
		for _, scope_id := range scope_ids {
			scopes = append(scopes, (*sm.ScopeMap)[scope_id])
		}
		return scopes, true
	}
	return []*Scope{}, false
}

// Returns the Scope at the top of the Stack, and bool if successful
func (sm *ScopeManager) Peek() (*Scope, bool) {
	if scope_id, ok := (*sm).PeekID(); ok {
		return (*(*sm).ScopeMap)[scope_id], true
	}
	return &Scope{}, false
}

// Returns the ID of the Scope at the top of the Stack, and bool if successful
func (sm *ScopeManager) PeekID() (ID, bool) {
	if scope_id, ok := (*sm).Stack.Peek(); ok {
		return scope_id, true
	}
	return ID(""), false
}

// Returns the Scope ID at the given Index, from 0.
func (sm *ScopeManager) Get(index int) (ID, bool) {
	return (*sm).Stack.Get(index)
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

// MapOfScopes
//
type MapOfScopes map[ID]*Scope

func NewMapOfScopes() *MapOfScopes {
	return &MapOfScopes{}
}

func (ms *MapOfScopes) ToString() string {
	var _string = ""
	for _, scope := range (*ms) {
		var _temp = fmt.Sprintf("\nScope: %d\n\tType: %d", scope.ID, scope.Type)
		_string = _string + _temp
	}
	return _string
}

// Scope
//
type Scope struct {
	ID    ID
	Node  *ast.Node
	Decls *ScopeDeclMap
	Type  ScopeType
}

// Returns a pointer to a Scope.
func NewScope(node *ast.Node, scope_type ScopeType) *Scope {
	var scope Scope

	scope.Node = node
	scope.Type = scope_type
	scope.Decls = NewScopeDeclMap()

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

// ScopeDeclMap
// Label => NewVarDeclID().ID
type ScopeDeclMap map[ID]ID

func NewScopeDeclMap() *ScopeDeclMap {
	return &ScopeDeclMap{}
}

// ID
//
type ID string

// Returns IDTrace of Scope IDs and their index in the slice.
// IDTrace concatenation in the form: index, ID: ...
func NewScopeID(node *ast.Node, scope_type ScopeType) ID {
	return ID(fmt.Sprintf("{SCOPE, %s: %v - %v}", scope_type, (*node).Pos(), (*node).End()))
}

// may be obselete
func NewVarID(node *ast.Node, var_context VarContext) ID {
	return ID(fmt.Sprintf("{VAR, %s: %v - %v}", var_context, (*node).Pos(), (*node).End()))
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
