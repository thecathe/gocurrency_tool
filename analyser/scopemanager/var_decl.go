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
			// get var type // also gets any values in data
			var_decl.Type = (*sm).NewVarType(node)
			var_decl.Token = tok

			if node_type.Values != nil {
				for _, value_expr := range node_type.Values {
					// add to values
					var_decl = *(var_decl).AddValue((*sm).NewVarValue(value_expr, node_type.Pos()))
				}
			}

			// add to ScopeManager
			(*sm.Decls)[(*sm).NewVarDeclID(&var_decl)] = &var_decl

			log.DebugLog("Analyser, NewVarDecl %d: %s, %s", len(*(*sm).Decls), var_decl.Type.Type, var_decl.Label)
		}
		return sm, true

	// variable assignment
	case *ast.AssignStmt:
		// find pre-existing decl
		// todo
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
					var_decl.Type = (*sm).NewVarType(node)
					var_decl.Token = token.DEFINE

					// add to values
					var_decl = *(var_decl).AddValue((*sm).NewVarValue(node_type.Rhs[index], node_type.Pos()))

					// add to ScopeManager
					(*sm.Decls)[(*sm).NewVarDeclID(&var_decl)] = &var_decl
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

// finds the value of a var in a given scope.
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
