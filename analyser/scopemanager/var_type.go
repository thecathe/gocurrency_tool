package scopemanager

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/thecathe/gocurrency_tool/analyser/log"
)

// VarType
// DataType is a list of Types
// Argument contains specific arguments like "BufferSize" for channels.
// Use ParseInt etc for extracting values from resulting string.
type VarType struct {
	Type GeneralVarType
	Data CompoundVarType
	Info map[string]string
}

func NewVarType() VarType {
	var var_type VarType
	var_type.Type = VAR_DATA_TYPE_NONE
	var_type.Data = NewCompoundVarType()
	var_type.Info = make(map[string]string, 0)
	return var_type
}

// Retruns VarType when node is:
// - *ast.AssignStmt
// - *ast.ValueSpec
// - *ast.FieldList
func (sm *ScopeManager) NewVarType(node ast.Node) VarType {

	var var_type VarType = NewVarType()

	switch node_type := (node).(type) {

	// Define (:=) assignment
	case *ast.AssignStmt:
		log.DebugLog("Analyser, NewVarType: *ast.AssignStmt\n")
		var_type.Info["NodeTrace"] = fmt.Sprintf("%s>%s", var_type.Info["NodeTrace"], "*ast.AssignStmt")

		switch node_type.Tok {

		case token.DEFINE:

			// take type from rhs
			switch rhs_expr := node_type.Rhs[0].(type) {

			// same type as predefined var
			case *ast.Ident:
				var_type.Info["NodeTrace"] = fmt.Sprintf("%s>%s", var_type.Info["NodeTrace"], "*ast.Ident")
				// find corresponding decl
				if _decl_id, _scope_id, _, _ := (*sm).FindDeclID(rhs_expr.Name); _scope_id != "" {
					var_type.Type = (*(*sm).Decls)[_decl_id].Type.Type
				} else {
					log.FailureLog("Analayser, NewVarType; *ast.Ident: label: \"%s\" did not yield a decl_id", rhs_expr.Name)
				}

			// int or string
			case *ast.BasicLit:
				var_type.Info["NodeTrace"] = fmt.Sprintf("%s>%s", var_type.Info["NodeTrace"], "*ast.BasicLit")
				var_type.Type = TokKindToVarType(rhs_expr.Kind)

			// Data received from channel
			case *ast.UnaryExpr:
				var_type.Info["NodeTrace"] = fmt.Sprintf("%s>%s", var_type.Info["NodeTrace"], "*ast.UnaryExpr")

				switch rhs_expr.Op {

				// received from channel
				case token.ARROW:
					log.GeneralLog("Analyser, NewVarType; data received from channel")
					var channel_name string = rhs_expr.X.(*ast.Ident).Name
					// search outwardly for first decl of this label
					if channel_decl_id, _, _, _ := (*sm).FindDeclID(channel_name); channel_decl_id != "" {
						// copy
						copy(var_type.Data, (*sm.Decls)[channel_decl_id].Type.Data[1:])
					}

				// unaccounted for
				default:
					log.DebugLog("Analyser, NewVarType; *ast.Unary: unaccounted for token: %s", rhs_expr.Op.String())
				}

			// Type from Function
			case *ast.CallExpr:
				var_type.Info["NodeTrace"] = fmt.Sprintf("%s>%s", var_type.Info["NodeTrace"], "*ast.CallExpr")
				var_type = *(*sm).CallExprVarType(rhs_expr)

			// some other func
			case *ast.CompositeLit:
				var_type.Info["NodeTrace"] = fmt.Sprintf("%s>%s", var_type.Info["NodeTrace"], "*ast.CompositeLit")
				var_type = *var_type.CompositeLit(rhs_expr)
			//
			default:
				var_type.Info["NodeTrace"] = fmt.Sprintf("%s>%s", var_type.Info["NodeTrace"], "*ast.Default")
				var_type.Type = VAR_DATA_TYPE_OTHER
			}
		}

	// Params
	case *ast.Field:
		log.DebugLog("Analyser, NewVarType: *ast.Field\n")
		var_type.Info["NodeTrace"] = fmt.Sprintf("%s>%s", var_type.Info["NodeTrace"], "*ast.Field")

	// Declaration
	case *ast.ValueSpec:
		log.DebugLog("Analyser, NewVarType: *ast.ValueSpec\n")
		var_type.Info["NodeTrace"] = fmt.Sprintf("%s>%s", var_type.Info["NodeTrace"], "*ast.ValueSpec")

		// look in type field
		if node_type.Type != nil {
			switch value_type := node_type.Type.(type) {

			// Pointer
			case *ast.Ident:
				var_type.Info["NodeTrace"] = fmt.Sprintf("%s>%s", var_type.Info["NodeTrace"], "*ast.Ident")
				var_type.Type = DataStringToVarType(value_type.Name)

			// Channel
			case *ast.ChanType:
				var_type.Info["NodeTrace"] = fmt.Sprintf("%s>%s", var_type.Info["NodeTrace"], "*ast.ChanType")
				log.DebugLog("Analyser, NewVarType: chan\n")
				var_type.Type = VAR_DATA_TYPE_CHAN
				var_type.Data = append(var_type.Data, VAR_DATA_TYPE_CHAN)

				switch rhs_expr := value_type.Value.(type) {

				// chan of chan
				case *ast.ChanType:
					var_type.Info["NodeTrace"] = fmt.Sprintf("%s>%s", var_type.Info["NodeTrace"], "*ast.ChanType")
					log.DebugLog("Analyser, NewVarType: chan chan\n")
					var_type.Data = append(var_type.Data, VAR_DATA_TYPE_CHAN)

					switch inner_chan_expr := rhs_expr.Value.(type) {
					// int or string
					case *ast.Ident:
						var_type.Info["NodeTrace"] = fmt.Sprintf("%s>%s", var_type.Info["NodeTrace"], "*ast.Ident")
						var_type.Data = append(var_type.Data, DataStringToVarType(inner_chan_expr.Name))

					// int or string
					case *ast.BasicLit:
						var_type.Info["NodeTrace"] = fmt.Sprintf("%s>%s", var_type.Info["NodeTrace"], "*ast.BasicLit")
						var_type.Data = append(var_type.Data, TokKindToVarType(inner_chan_expr.Kind))

					// Type from Function
					case *ast.CallExpr:
						var_type.Info["NodeTrace"] = fmt.Sprintf("%s>%s", var_type.Info["NodeTrace"], "*ast.CallExpr")
						var_type = *(*sm).CallExprVarType(inner_chan_expr)

					// some other func
					case *ast.CompositeLit:
						var_type.Info["NodeTrace"] = fmt.Sprintf("%s>%s", var_type.Info["NodeTrace"], "*ast.CompositeLit")
						var_type = *var_type.CompositeLit(inner_chan_expr)

					}

				// int or string
				case *ast.BasicLit:
					var_type.Info["NodeTrace"] = fmt.Sprintf("%s>%s", var_type.Info["NodeTrace"], "*ast.BasicLit")
					var_type.Type = TokKindToVarType(rhs_expr.Kind)

				// Type from Function
				case *ast.CallExpr:
					var_type.Info["NodeTrace"] = fmt.Sprintf("%s>%s", var_type.Info["NodeTrace"], "*ast.CallExpr")
					var_type = *(*sm).CallExprVarType(rhs_expr)

				// some other func
				case *ast.CompositeLit:
					var_type.Info["NodeTrace"] = fmt.Sprintf("%s>%s", var_type.Info["NodeTrace"], "*ast.CompositeLit")
					var_type = *var_type.CompositeLit(rhs_expr)

				// unnaccounted for
				default:
					var_type.Info["NodeTrace"] = fmt.Sprintf("%s>%s", var_type.Info["NodeTrace"], "*ast.Default")
					var_type.Type = VAR_DATA_TYPE_OTHER
				}
			}
		}

	// unnaccounted for
	default:
		var_type.Type = VAR_DATA_TYPE_OTHER
		var_type.Info["NodeTrace"] = fmt.Sprintf("%s>%s", var_type.Info["NodeTrace"], "*ast.Default")
	}

	log.DebugLog("Analyser, leaving NewVarType: %s\n", var_type.Type)
	return var_type
}

// comsposite lit to extractexpr
func (vt *VarType) CompositeLit(node *ast.CompositeLit) *VarType {

	log.DebugLog("Analyser, NewVarType: CompositeLit\n")
	switch node.Type.(type) {
	// Get them all
	case *ast.SelectorExpr:

		// add to type
		vt = vt.ExtractExpr(node.Type)
		vt.Type = VAR_DATA_FUNC_RET

	default:
		vt.Info["NodeTrace"] = fmt.Sprintf("%s>%s", vt.Info["NodeTrace"], "*ast.Default")
		vt.Type = VAR_DATA_TYPE_OTHER
	}

	return vt
}

// Returns []string containing selectorexpor x, sel of compositelit in each element
// ast.Expr should be of type *ast.SelectorExpr
func (vt *VarType) ExtractExpr(current_sel_expr ast.Expr) *VarType {

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
				vt.Info["NodeTrace"] = fmt.Sprintf("%s>%s", vt.Info["NodeTrace"], "*ast.SelectorExpr")
				// make x selector
				current_sel_expr = inner_sel_type

			// Ident
			case *ast.Ident:
				vt.Info["NodeTrace"] = fmt.Sprintf("%s>%s", vt.Info["NodeTrace"], "*ast.Ident")
				// add x to beginning
				sel_expr = append([]string{inner_sel_type.Name}, sel_expr...)
				// end loop
				loop = false
			}

		default:
			vt.Info["NodeTrace"] = fmt.Sprintf("%s>%s", vt.Info["NodeTrace"], "*ast.Default")
			loop = false
		}
	}

	// context
	vt.Info["SelectorExpr"] = strings.Join(sel_expr, ";")
	vt.Info["Function"] = sel_expr[len(sel_expr)-1]

	return vt
}

// retrieves vartype data, and determins vartype type
func (vt *VarType) DetermineType(node ast.Node) *VarType {
	vt.Data = CompoundVarTypeBuilder(node)

	// if only one, must be this
	if len(vt.Data) == 1 {
		vt.Type = vt.Data[0]
		return vt
	}

	return vt
}

// returns the var type derived from a call expr
func (sm *ScopeManager) CallExprVarType(node *ast.CallExpr) *VarType {

	var _var_type VarType = NewVarType()

	log.DebugLog("Analyser, NewVarValue: CallExpr\n")

	var call_name string = node.Fun.(*ast.Ident).Name
	_var_type.Info["Function"] = call_name

	switch call_name {

	// Channel or Slice
	case "make":
		log.DebugLog("Analyser, NewVarValue: make\n")

		// extract chan type
		_var_type = *_var_type.DetermineType(node.Args[0])

		switch node.Args[0].(type) {

		// Channel
		case *ast.ChanType:

			// If Async. Channel
			if len(node.Args) > 1 {
				_var_type.Type = VAR_DATA_TYPE_ASYNC_CHAN

				// get buffer size
				switch _buffer_expr := node.Args[1].(type) {

				// Buffer inline
				case *ast.BasicLit:
					_var_type.Info["BufferSize"] = fmt.Sprintf("%v", _buffer_expr.Value)

				// Buffer from var
				case *ast.Ident:
					// search outwardly for first decl of this label
					if _decl_id, _scope_id, _elevated, _elevated_id := (*sm).FindDeclID(_buffer_expr.Name); _scope_id != "" {
						if _elevated {
							// switch out scope id with elevated
							_scope_id = _elevated_id
						}
						// get value
						if _index, _buffer_value := (*(*sm.Decls)[_decl_id]).FindValue(_scope_id); _index >= 0 {
							_var_type.Info["BufferSize"] = _buffer_value.Value
						} else {
							log.WarningLog("Analyser, NewVarValue, decl found, but value not found: %s,\t%s", _buffer_expr.Name, _scope_id)
						}

					} else {
						log.WarningLog("Analyser, NewVarValue, decl not found: %s", _buffer_expr.Name)
					}

				default:
					_var_type.Info["BufferSize"] = fmt.Sprintf("unknown: [%v]", _buffer_expr)
				}
			} else {
				// sync channel
				_var_type.Type = VAR_DATA_TYPE_SYNC_CHAN
				_var_type.Info["BufferSize"] = "0"
			}

		default:
			_var_type.Type = VAR_DATA_TYPE_OTHER
		}

	}

	return &_var_type
}
