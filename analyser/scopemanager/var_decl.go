package scopemanager

import (
	"go/ast"
	"go/token"

	"github.com/thecathe/gocurrency_tool/analyser/log"
)

// VarDecl
type VarDecl struct {
	Node   *ast.Node
	Label  string
	Type   VarType
	Values []VarValue
	Token  token.Token
}

// Creates a new VarDecl and adds it the the MapOfVarDecl
// Node should be of type *ast.ValueSpec or *ast.AssignStmt
func (sm *ScopeManager) NewVarDecl(node ast.Node, tok token.Token) (*ScopeManager, bool) {

	switch node_type := (node).(type) {
	// var decl
	case *ast.ValueSpec:

		// for each var being decl
		for _, name := range node_type.Names {
			// declaration
			var var_decl VarDecl
			var_decl.Node = &node

			var_decl.Label = name.Name
			log.GeneralLog("Analyser, NewVarDecl; Valuespec: %s\n\n", var_decl.Label)
			// get var type // also gets any values in data
			var_decl.Type = (*sm).NewVarType(node)
			var_decl.Token = tok

			if len(node_type.Values) == 1 {
				value_expr := node_type.Values[0]

				// get value
				_value, _var_type := (*sm).NewVarValue(value_expr, node_type.Pos())

				// add to values
				var_decl = *(var_decl).AddValue(_value)

				// if not found from decl, find in assignment
				if var_decl.Type.Type == VAR_DATA_TYPE_NONE {
					log.DebugLog("Analyser, NewVarDecl: type not found in decl")
					if _var_type.Type == VAR_DATA_TYPE_NONE {
						log.FailureLog("Analyser, NewVarDecl: unable to infer type from assignment")
					} else {
						var_decl.Type = _var_type
					}
				}

				// add to scope
				var _new_decl_id = (*sm).NewVarDeclID(&var_decl)

				if current_scope_id, ok := (*sm).PeekID(); ok {
					(*(*sm).ScopeMap)[current_scope_id] = (*(*sm).ScopeMap)[current_scope_id].AddDecl(_new_decl_id, var_decl.Label)
				} else {
					log.FailureLog("Analyser, NewVarDecl: Failed to add decl to scope list: ScopeID: %s, DeclID: %s", current_scope_id, _new_decl_id)
				}

				// add to ScopeManager
				(*sm.Decls)[_new_decl_id] = &var_decl

				// // if chan, add to var type info
				// if var_decl.Type.Type == VAR_DATA_TYPE_CHAN {
				// 	var_decl.Type.Info["ChanType"] = _info["ChanType"]
				// 	var_decl.Type.Info["BufferSize"] = _info["BufferSize"]

				// 	// update chan type
				// 	if var_decl.Type.Info["ChanType"] == "async" {
				// 		var_decl.Type.Type = VAR_DATA_TYPE_ASYNC_CHAN
				// 	} else if var_decl.Type.Info["ChanType"] == "sync" {
				// 		var_decl.Type.Type = VAR_DATA_TYPE_SYNC_CHAN
				// 	} else {
				// 		log.WarningLog("Analyser, NewVarDecl Chan info error: %s", var_decl.Type.Info["ChanType"])
				// 	}

				// 	log.DebugLog("Analyser, NewVarDecl Chan Details: %s, %s", var_decl.Type.Type, var_decl.Type.Info["BufferSize"])
				// }

			} else {
				log.WarningLog("Analyser, NewVarDecl: Values len greater than 1")
			}

		}
		return sm, true

	// variable assignment
	case *ast.AssignStmt:
		log.DebugLog("Analyser, NewVarDecl: *ast.AssignStmt")
		var var_decl VarDecl

		switch node_type.Tok {
		case token.DEFINE:
			// check decl
			// for each decl
			for index, expr := range node_type.Lhs {
				// ensure ident
				switch expr_ident := expr.(type) {
				case *ast.Ident:

					var_decl.Label = expr_ident.Name
					log.GeneralLog("Analyser, NewVarDecl; Assignstmt, Ident: %s\n\n", var_decl.Label)

					var_decl.Type = (*sm).NewVarType(node)
					var_decl.Token = token.DEFINE

					// get value
					_value, _var_type := (*sm).NewVarValue(node_type.Rhs[index], node_type.Pos())

					// add to values
					var_decl = *(var_decl).AddValue(_value)

					// if not found from decl, find in assignment
					if var_decl.Type.Type == VAR_DATA_TYPE_NONE {
						log.DebugLog("Analyser, NewVarDecl: type not found in decl")
						if _var_type.Type == VAR_DATA_TYPE_NONE {
							log.FailureLog("Analyser, NewVarDecl: unable to infer type from assignment")
						} else {
							var_decl.Type = _var_type
						}
					}

					// add to scope
					var _new_decl_id = (*sm).NewVarDeclID(&var_decl)

					if current_scope_id, ok := (*sm).PeekID(); ok {
						(*(*sm).ScopeMap)[current_scope_id] = (*(*sm).ScopeMap)[current_scope_id].AddDecl(_new_decl_id, var_decl.Label)
					} else {
						log.FailureLog("Analyser, NewVarDecl: Failed to add decl to scope list: ScopeID: %s, DeclID: %s", current_scope_id, _new_decl_id)
					}

					// add to ScopeManager
					(*sm.Decls)[_new_decl_id] = &var_decl
				}
			}
		}
		return sm, true

	// unnaccounted for
	default:
		return sm, false
	}
}

func (decl *VarDecl) AddValue(value VarValue) *VarDecl {
	decl.Values = append(decl.Values, value)
	return decl
}

// Returns the value of a var at a scope, and its value index in decl values.
// returns -1 if not found
func (decl *VarDecl) FindValue(scope_id ID) (int, VarValue) {
	for _index, _vars := range (*decl).Values {
		// found
		if _vars.ScopeID == scope_id {
			return _index, _vars
		}
	}
	// return nothing
	return -1, VarValue{}
}

func (var_decl *VarDecl) Pos() token.Pos {
	return (*var_decl.Node).Pos()
}

func (var_decl *VarDecl) End() token.Pos {
	return (*var_decl.Node).End()
}

func (var_decl *VarDecl) ID() token.Pos {
	return (*var_decl.Node).Pos()
}

// MapOfDecls
type MapOfDecls map[ID]*VarDecl

func NewMapOfDecls() *MapOfDecls {
	return &MapOfDecls{}
}

func (decls *MapOfDecls) Size() int {
	return len(*decls)
}
