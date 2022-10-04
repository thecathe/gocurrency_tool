package scopemanager

import (
	"fmt"
	"go/ast"
	"go/token"
	"os"

	"github.com/thecathe/gocurrency_tool/analyser/log"
)

// ScopeManager
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

func (sm *ScopeManager) CheckAwaitedFunction(node *ast.Node) (*ScopeManager, bool) {

	// x:=(*node).(*ast.Ident).

	return sm, false
}

func (sm *ScopeManager) ParseNode(node ast.Node) (*ScopeManager, ParseType) {

	// if not first scope
	if _size := (*sm).StackSize(); _size > 0 {
		// Check if leaving current scope
		if outer_scope, ok := (*sm).Peek(); ok {
			// check not leaving file
			if outer_scope.Type != SCOPE_TYPE_FILE {
				if (node).Pos() > (*(outer_scope).Node).End() {
					// if current node starts after the current scope ends, left current scope
					// log.DebugLog("Analyser; ParseNode, Exiting Scope: %d > %d\n", (node).Pos(), (*(outer_scope).Node).End())
					if _sm, ok := (*sm).Pop(); ok {
						sm = _sm
						// log.DebugLog("Analyser; ParseNode, Exiting Scope: %d -> %d\n", _size, (*sm).StackSize())
						log.GeneralLog("Analyser; ParseNode, Exiting Scope: %s\n", outer_scope.ID)
						// do not return, continue add scope
						// return sm, PARSE_SCOPE_EXIT
					} else {
						// failed
						log.FailureLog("Analyser; ParseNode, StackPop.\n")
						return sm, PARSE_FAIL_STACK_POP
					}
				} // continue
			}
		} else {
			log.FailureLog("Analyser; ParseNode, StackPeek: Size %d\n", _size)
			return sm, PARSE_FAIL_STACK_PEEK
		}
	} else {
		log.DebugLog("Analyser; ParseNode: First Scope\n\n")
		if scope, ok := (*sm).Peek(); ok {
			log.WarningLog("Analyser; ParseNode, StackPeek Successful : %d | Scope: %s\n", _size, scope.ID)
		}
	}

	// Check for each ScopeType
	switch node_type := (node).(type) {

	// Debug: import
	case *ast.ImportSpec:
		// check outerscope
		// scope or vardecl, depends on outerscope
		if outer_scope, ok := (*sm).Peek(); ok {
			if outer_scope.Type == SCOPE_TYPE_FILE_IMPORT {
				log.DebugLog("Analyser; ParseNode, Package: ImportSpec\n")
			} else {
				log.WarningLog("Analyser; ParseNode, Package: Unknown ImportSpec\n")
			}
		}
		return sm, PARSE_NONE

	// Scope: Package
	case *ast.Package:
		log.DebugLog("Analyser; ParseNode, Package\n")
		sm = (*sm).NewScope(node, SCOPE_TYPE_PACKAGE)
		return sm, PARSE_PACKAGE

	// Scope: File
	case *ast.File:
		log.DebugLog("Analyser; ParseNode, File\n")
		sm = (*sm).NewScope(node, SCOPE_TYPE_FILE)
		return sm, PARSE_FILE

	// Scope: GenDecl
	case *ast.GenDecl:
		// scope or vardecl, depends on outerscope
		if outer_scope, ok := (*sm).Peek(); ok {
			log.DebugLog("Analyser; ParseNode, %s contains GenDecl\n", outer_scope.Type)

			switch outer_scope.Type {

			// if file, this is new scope of global decl
			case SCOPE_TYPE_FILE:
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
					log.FailureLog("Analyser; ParseNode, DeclTok\n")
					return sm, PARSE_FAIL_DECL_TOKEN
				}

			// if imports, skip
			case SCOPE_TYPE_FILE_IMPORT:
				log.DebugLog("Analyser; ParseNode, GenDecl: Import.\n")
				return sm, PARSE_FILE_IMPORT

			default:
				log.DebugLog("Analyser; ParseNode, GenDecl Unknown outerscope: %s\n", outer_scope.Type)
				return sm, PARSE_NONE
			}
		} else {
			log.FailureLog("Analyser; ParseNode, StackPeek: Size %d\n", (*sm).StackSize())
			return sm, PARSE_FAIL_STACK_PEEK
		}

	// Scope: Goroutine
	case *ast.GoStmt:
		log.DebugLog("Analyser; ParseNode, GoStmt\n")
		sm = (*sm).NewScope(node, SCOPE_TYPE_GOROUTINE)
		return sm, PARSE_GO_STMT

	// Scope: Anon Function
	case *ast.FuncLit:
		log.DebugLog("Analyser; ParseNode, FuncLit\n")
		sm = (*sm).NewScope(node, SCOPE_TYPE_FUNC_DECL)
		return sm, PARSE_FUNC_LIT

	// Scope: Function
	case *ast.FuncDecl: // line 1914
		log.DebugLog("Analyser; ParseNode, FuncDecl\n")
		sm = (*sm).NewScope(node, SCOPE_TYPE_FUNC_DECL)
		return sm, PARSE_FUNC_DECL

	// Scope: FuncCall
	case *ast.CallExpr:
		log.DebugLog("Analyser; ParseNode, CallExpr\n")
		// check function call
		if outer_scope, ok := (*sm).Peek(); ok {
			if outer_scope.Type == SCOPE_TYPE_GOROUTINE {
				sm = (*sm).NewScope(node, SCOPE_TYPE_FUNC_CALL)
				return sm, PARSE_FUNC_DECL
			} else {
				log.DebugLog("Analyser; ParseNode, CallExpr: not in goroutine.\n")
				return sm, PARSE_NONE
			}
		} else {
			log.FailureLog("Analyser; ParseNode, StackPeek: Size %d\n", (*sm).StackSize())
			return sm, PARSE_FAIL_STACK_PEEK
		}

	// Scope: If Statement
	case *ast.IfStmt:
		log.DebugLog("Analyser; ParseNode, IfStmt\n")
		sm = (*sm).NewScope(node, SCOPE_TYPE_IF)
		return sm, PARSE_IF_STMT

	// Scope: Select Statement
	case *ast.SelectStmt:
		log.DebugLog("Analyser; ParseNode, SelectStmt\n")
		sm = (*sm).NewScope(node, SCOPE_TYPE_SELECT)
		return sm, PARSE_SELECT_STMT

	// Scope: Switch Statement
	case *ast.SwitchStmt:
		log.DebugLog("Analyser; ParseNode, SwitchStmt\n")
		sm = (*sm).NewScope(node, SCOPE_TYPE_SWITCH)
		return sm, PARSE_SWTICH_STMT

	// Scope: Switch Type Statement
	case *ast.TypeSwitchStmt:
		log.DebugLog("Analyser; ParseNode, TypeSwitchStmt\n")
		sm = (*sm).NewScope(node, SCOPE_TYPE_TYPE_SWITCH)
		return sm, PARSE_TYPE_SWITCH_STMT

	// Scope: For Loop Statement
	case *ast.ForStmt:
		log.DebugLog("Analyser; ParseNode, ForStmt\n")
		sm = (*sm).NewScope(node, SCOPE_TYPE_FOR)
		return sm, PARSE_FOR_STMT

	// Scope: Ranged For Loop Statement
	case *ast.RangeStmt:
		log.DebugLog("Analyser; ParseNode, RangeStmt\n")
		sm = (*sm).NewScope(node, SCOPE_TYPE_RANGE)
		return sm, PARSE_RANGE_STMT

	// Var: Params
	case *ast.FieldList: // line 1924
		log.DebugLog("Analyser; ParseNode, FieldList\n")
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
					log.FailureLog("Analyser; ParseNode, FieldList\n")
					return sm, PARSE_FAIL_FIELD_LIST
				}
			} else {
				// log.DebugLog("Analyser; ParseNode, FieldList: not function params.\n")
				return sm, PARSE_NONE
			}
		} else {
			log.FailureLog("Analyser; ParseNode, StackPeek: Size %d\n", (*sm).StackSize())
			return sm, PARSE_FAIL_STACK_PEEK
		}

	// Var: Declaration
	case *ast.DeclStmt:
		log.DebugLog("Analyser; ParseNode, DeclStmt\n")

		// TODO
		switch _decl := node_type.Decl.(type) {
		case *ast.GenDecl:
			// should happen once?
			for _, _spec := range _decl.Specs {
				switch spec := _spec.(type) {
				case *ast.ValueSpec:
					(*sm).NewVarDecl(spec, _decl.Tok)
				}
			}
		}

		return sm, PARSE_DECL

	// Var: VarDecl, global or scoped
	case *ast.ValueSpec:
		log.DebugLog("Analyser; ParseNode, ValueSpec\n")
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
					log.DebugLog("Analyser; ParseNode, ValueSpec: package var\n")
					if _sm, ok := (*sm).NewVarDecl(node, token.VAR); ok {
						sm = _sm
					}
					return sm, PARSE_PACKAGE_VAR

				case SCOPE_TYPE_PACKAGE_CONST:
					log.DebugLog("Analyser; ParseNode, ValueSpec: package var\n")
					if _sm, ok := (*sm).NewVarDecl(node, token.CONST); ok {
						sm = _sm
					}
					return sm, PARSE_PACKAGE_CONST

				case SCOPE_TYPE_FILE_IMPORT:
					log.DebugLog("Analyser; ParseNode, ValueSpec: package var\n")
					if _sm, ok := (*sm).NewVarDecl(node, token.IMPORT); ok {
						sm = _sm
					}
					return sm, PARSE_FILE_IMPORT
				// not accounted for
				default:
					log.FailureLog("Analyser; ParseNode, ValueSpec\n")
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
			log.FailureLog("Analyser; ParseNode, StackPeek: Size %d\n", (*sm).StackSize())
			return sm, PARSE_FAIL_STACK_PEEK
		}

	// Var: Assign or Decl
	case *ast.AssignStmt:
		log.DebugLog("Analyser; ParseNode, AssignStmt\n")
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
		log.VerboseLog("Analyser; ParseNode, Default: Nonthing of interest\n")
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

// Creates a new VarDecl and adds it the the MapOfVarDecl
// Node should be of type *ast.ValueSpec or *ast.AssignStmt
func (sm *ScopeManager) NewVarDecl(node ast.Node, tok token.Token) (*ScopeManager, bool) {

	// declaration
	var var_decl VarDecl
	var_decl.Node = &node

	switch node_type := (node).(type) {
	case *ast.ValueSpec:

		var_decl.Label = node_type.Names[0].Name
		var_decl.Type = (*sm).NewVarType(node)
		var_decl.Token = tok

		// check for value
		if node_type.Values != nil {
			for _, value_expr := range node_type.Values {
				// add to values
				var_decl = *(var_decl).AddValue((*sm).NewVarValue(value_expr, node_type.Pos()))
			}

		}

		// add to ScopeManager
		(*sm.Decls)[(*sm).NewVarDeclID(&var_decl)] = &var_decl

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

					var_decl.Label = expr_ident.Name
					var_decl.Type = (*sm).NewVarType(node)
					var_decl.Token = token.DEFINE

					// add to values
					var_decl = *(var_decl).AddValue((*sm).NewVarValue(node_type.Rhs[index], node_type.Pos()))

					// add to ScopeManager
					(*sm.Decls)[(*sm).NewVarDeclID(&var_decl)] = &var_decl
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

// Returns ID consisting of
func (sm *ScopeManager) NewVarDeclID(decl *VarDecl) ID {
	if scope_id, ok := (*sm).PeekID(); ok {
		return NewVarDeclID(decl.Label, scope_id)
	}
	// fail
	return ID("Fail: VarDeclID")
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
func (sm *ScopeManager) NewScope(node ast.Node, scope_type ScopeType) *ScopeManager {
	var scope Scope = *NewScope(node, scope_type)

	// add id to stack
	sm = (*sm).Push(scope.ID)

	// add scope to map
	(*(*sm).ScopeMap)[scope.ID] = &scope

	log.GeneralLog("Analyser; NewScope %d: %s\n\n", (*sm).StackSize(), scope.ID)
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
	// log.DebugLog("Entering peekX\n")
	if scope_ids, ok := (*sm).Stack.PeekX(x); ok {
		// log.DebugLog("Done Scope ID's peekX\n")
		var scopes []*Scope
		for _, scope_id := range *scope_ids {
			scopes = append(scopes, (*sm.ScopeMap)[scope_id])
		}
		// log.DebugLog("Successful peekX\n")
		return scopes, true
	}
	return []*Scope{}, false
}

// Returns the Scope at the top of the Stack, and bool if successful
func (sm *ScopeManager) Peek() (*Scope, bool) {
	if scope_id, ok := (*sm).PeekID(); ok {
		return (*(*sm).ScopeMap)[scope_id], true
	} else {
		return &Scope{}, false
	}
}

// Returns the ID of the Scope at the top of the Stack, and bool if successful
func (sm *ScopeManager) PeekID() (ID, bool) {
	if scope_id, ok := (*sm).Stack.Peek(); ok {
		return scope_id, true
	}
	return ID(""), false
}

// Returns size of index
func (sm *ScopeManager) StackSize() int {
	return (*sm).Stack.Size()
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
	if stack, ok := ((*sm).Stack).Pop(); ok {
		(*sm).Stack = stack
		return sm, true
	}
	return sm, false
}
