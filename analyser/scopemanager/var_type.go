package scopemanager

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"

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
		log.DebugLog("NewVarType: *ast.AssignStmt\n")
		var_type.Info["NodeTrace"] = fmt.Sprintf("%s > %s", var_type.Info["NodeTrace"], "*ast.AssignStmt")

		switch node_type.Tok {

		case token.DEFINE:

			// take type from rhs
			switch rhs_expr := node_type.Rhs[0].(type) {

			// same type as predefined var
			case *ast.Ident:
				var_type.Info["NodeTrace"] = fmt.Sprintf("%s > %s", var_type.Info["NodeTrace"], "*ast.Ident")
				// check bool
				if rhs_expr.Name == "true" {
					var_type.Type = VAR_DATA_TYPE_BOOL
				} else if rhs_expr.Name == "false" {
					var_type.Type = VAR_DATA_TYPE_BOOL
				} else {
					// find corresponding decl
					if _decl_id, _scope_id, _, _ := (*sm).FindDeclID(rhs_expr.Name); _scope_id != "" {
						var_type.Type = (*(*sm).Decls)[_decl_id].Type.Type
					} else {
						log.FailureLog("NewVarType; *ast.Ident: label: \"%s\" did not yield a decl_id", rhs_expr.Name)
					}
				}

			// int or string
			case *ast.BasicLit:
				var_type.Info["NodeTrace"] = fmt.Sprintf("%s > %s", var_type.Info["NodeTrace"], "*ast.BasicLit")
				var_type.Type = TokKindToVarType(rhs_expr.Kind)

			// Data received from channel
			case *ast.UnaryExpr:
				var_type.Info["NodeTrace"] = fmt.Sprintf("%s > %s", var_type.Info["NodeTrace"], "*ast.UnaryExpr")

				switch rhs_expr.Op {

				// received from channel
				case token.ARROW:
					log.GeneralLog("NewVarType; data received from channel")
					var channel_name string = rhs_expr.X.(*ast.Ident).Name
					// search outwardly for first decl of this label
					if channel_decl_id, _scope_id, _elevated, _elevated_id := (*sm).FindDeclID(channel_name); channel_decl_id != "" {
						if _elevated {
							_scope_id = _elevated_id
						}
						// channels known to have value as most specific data type
						if _index, _var_value := (*(*sm).Decls)[channel_decl_id].FindValue(_scope_id); _index >= 0 {
							var_type.Type = GeneralVarType(_var_value.Value)
						}
					}

				// unaccounted for
				default:
					var_type.Info["NodeTrace"] = fmt.Sprintf("%s > %s", var_type.Info["NodeTrace"], "Default")
					log.DebugLog("NewVarType; *ast.Unary: unaccounted for token: %s, should recover in VarValue", rhs_expr.Op.String())
				}

			// Type from Function
			case *ast.CallExpr:
				var_type.Info["NodeTrace"] = fmt.Sprintf("%s > %s", var_type.Info["NodeTrace"], "*ast.CallExpr")
				var_type = *(*sm).CallExprVarType(&var_type, rhs_expr)

			// some other func
			case *ast.CompositeLit:
				var_type.Info["NodeTrace"] = fmt.Sprintf("%s > %s", var_type.Info["NodeTrace"], "*ast.CompositeLit")
				var_type = *var_type.CompositeLit(rhs_expr)
			//
			default:
				var_type.Info["NodeTrace"] = fmt.Sprintf("%s > %s", var_type.Info["NodeTrace"], "Default")
				var_type.Type = VAR_DATA_TYPE_OTHER
			}
		}

	// Params
	case *ast.Field:
		log.DebugLog("NewVarType: *ast.Field\n")
		var_type.Info["NodeTrace"] = fmt.Sprintf("%s > %s", var_type.Info["NodeTrace"], "*ast.Field")

	// Declaration
	case *ast.ValueSpec:
		log.DebugLog("NewVarType: *ast.ValueSpec\n")
		var_type.Info["NodeTrace"] = fmt.Sprintf("%s > %s", var_type.Info["NodeTrace"], "*ast.ValueSpec")

		// look in type field
		if node_type.Type != nil {
			switch value_type := node_type.Type.(type) {

			// Pointer
			case *ast.Ident:
				var_type.Info["NodeTrace"] = fmt.Sprintf("%s > %s", var_type.Info["NodeTrace"], "*ast.Ident")
				var_type.Type = DataStringToVarType(value_type.Name)

			// Channel
			case *ast.ChanType:
				var_type.Info["NodeTrace"] = fmt.Sprintf("%s > %s", var_type.Info["NodeTrace"], "*ast.ChanType")
				var_type.Type = VAR_DATA_TYPE_CHAN
				var_type = *(*sm).ChanTypeVarType(&var_type, value_type)

			//
			default:
				var_type.Info["NodeTrace"] = fmt.Sprintf("%s > %s", var_type.Info["NodeTrace"], "Default")
				var_type.Type = VAR_DATA_TYPE_OTHER
			}
		}

	// unnaccounted for
	default:
		var_type.Type = VAR_DATA_TYPE_OTHER
		var_type.Info["NodeTrace"] = fmt.Sprintf("%s > %s", var_type.Info["NodeTrace"], "Default")
	}

	log.DebugLog("leaving NewVarType: %s\n", var_type.Type)
	return var_type
}

// comsposite lit to extractexpr
func (vt *VarType) CompositeLit(node *ast.CompositeLit) *VarType {

	log.DebugLog("NewVarType: CompositeLit\n")
	switch node.Type.(type) {
	// Get them all
	case *ast.SelectorExpr:
		vt.Info["NodeTrace"] = fmt.Sprintf("%s > %s", vt.Info["NodeTrace"], "*ast.SelectorExpr")

		// add to type
		vt = vt.ExtractExpr(node.Type)
		vt.Type = VAR_DATA_FUNC_RET

	default:
		vt.Info["NodeTrace"] = fmt.Sprintf("%s > %s", vt.Info["NodeTrace"], "Default")
		vt.Type = VAR_DATA_TYPE_OTHER
	}

	return vt
}

// Returns []string containing selectorexpor x, sel of compositelit in each element
// ast.Expr should be of type *ast.SelectorExpr
func (vt *VarType) ExtractExpr(current_sel_expr ast.Expr) *VarType {

	// For recursion on X,
	switch outer_sel_type := current_sel_expr.(type) {
	case *ast.SelectorExpr:

		// extracting from x
		switch inner_sel_type := outer_sel_type.X.(type) {

		// Selector
		case *ast.SelectorExpr:
			vt.Info["NodeTrace"] = fmt.Sprintf("%s > %s", vt.Info["NodeTrace"], "*ast.SelectorExpr")
			// make x selector
			vt = vt.ExtractExpr(inner_sel_type)
			vt.Info["SelectorExpr"] = fmt.Sprintf("%s;%s", vt.Info["SelectorExpr"], inner_sel_type.Sel.Name)
			return vt

		// Ident
		case *ast.Ident:
			vt.Info["NodeTrace"] = fmt.Sprintf("%s > %s", vt.Info["NodeTrace"], "*ast.Ident")
			vt.Info["Function"] = inner_sel_type.Name
			// add x to beginning
			vt.Info["SelectorExpr"] = fmt.Sprintf("%s;%s", vt.Info["SelectorExpr"], inner_sel_type.Name)
			return vt

		default:
			vt.Info["NodeTrace"] = fmt.Sprintf("%s > %s", vt.Info["NodeTrace"], "Default")
			vt.Info["Function"] = string(VAR_DATA_TYPE_OTHER)
			// add x to beginning
			vt.Info["SelectorExpr"] = fmt.Sprintf("%s;%s", vt.Info["SelectorExpr"], VAR_DATA_TYPE_OTHER)
			return vt
		}

	default:
		log.FailureLog("ExtractExpr; called on non-*ast.SelectorExpr: %v", outer_sel_type)
		vt.Info["NodeTrace"] = fmt.Sprintf("%s > %s", vt.Info["NodeTrace"], "Default")
		vt.Info["Function"] = string(VAR_DATA_TYPE_OTHER)
		// add x to beginning
		vt.Info["SelectorExpr"] = fmt.Sprintf("%s;%s", vt.Info["SelectorExpr"], VAR_DATA_TYPE_NONE)
		return vt
	}

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

// returns var type for chans
func (sm *ScopeManager) ChanTypeVarType(var_type *VarType, node *ast.ChanType) *VarType {

	var_type.Data = append(var_type.Data, VAR_DATA_TYPE_CHAN)

	switch value_type := node.Value.(type) {

	// chan of chan
	case *ast.ChanType:
		var_type.Info["NodeTrace"] = fmt.Sprintf("%s > %s", var_type.Info["NodeTrace"], "*ast.ChanType")
		var_type = (*sm).ChanTypeVarType(var_type, value_type)
		return var_type

	// int or string
	case *ast.Ident:
		var_type.Info["NodeTrace"] = fmt.Sprintf("%s > %s", var_type.Info["NodeTrace"], "*ast.Ident")
		var_type.Data = append(var_type.Data, DataStringToVarType(value_type.Name))
		return var_type

	// int or string
	case *ast.BasicLit:
		var_type.Info["NodeTrace"] = fmt.Sprintf("%s > %s", var_type.Info["NodeTrace"], "*ast.BasicLit")
		var_type.Type = TokKindToVarType(value_type.Kind)
		return var_type

	// Type from Function
	case *ast.CallExpr:
		var_type.Info["NodeTrace"] = fmt.Sprintf("%s > %s", var_type.Info["NodeTrace"], "*ast.CallExpr")
		var_type = (*sm).CallExprVarType(var_type, value_type)
		return var_type

	// some other func
	case *ast.CompositeLit:
		var_type.Info["NodeTrace"] = fmt.Sprintf("%s > %s", var_type.Info["NodeTrace"], "*ast.CompositeLit")
		var_type = var_type.CompositeLit(value_type)
		return var_type

	// unnaccounted for
	default:
		var_type.Info["NodeTrace"] = fmt.Sprintf("%s > %s", var_type.Info["NodeTrace"], "Default")
		var_type.Type = VAR_DATA_TYPE_OTHER
		return var_type
	}
}

// returns the var type derived from a call expr
func (sm *ScopeManager) CallExprVarType(var_type *VarType, node *ast.CallExpr) *VarType {

	var call_name string = node.Fun.(*ast.Ident).Name
	var_type.Info["Function"] = call_name

	switch call_name {

	// Channel or Slice
	case "make":

		// extract chan type
		var_type = var_type.DetermineType(node.Args[0])

		switch node.Args[0].(type) {

		// Channel
		case *ast.ChanType:
			var_type.Info["CallTrace"] = fmt.Sprintf("%s > %s", var_type.Info["CallTrace"], "*ast.ChanType")

			// If Async. Channel
			if len(node.Args) > 1 {

				// check if value of buffer is 0
				var chan_buffer_size int = -1
				// get buffer size
				switch _buffer_expr := node.Args[1].(type) {

				// Buffer inline
				case *ast.BasicLit:
					var_type.Info["CallTrace"] = fmt.Sprintf("%s > %s", var_type.Info["CallTrace"], "*ast.BasicLit")
					// get size of buffer
					if _size, err := strconv.Atoi(fmt.Sprintf("%v", _buffer_expr.Value)); err == nil {
						chan_buffer_size = _size
					} else {
						log.FailureLog("CallExprVarType: unable to determine size of channel buffer: %v", _buffer_expr.Value)
					}

				// Buffer from var
				case *ast.Ident:
					var_type.Info["CallTrace"] = fmt.Sprintf("%s > %s", var_type.Info["CallTrace"], "*ast.Ident")
					// search outwardly for first decl of this label
					if _decl_id, _scope_id, _elevated, _elevated_id := (*sm).FindDeclID(_buffer_expr.Name); _scope_id != "" {
						if _elevated {
							// switch out scope id with elevated
							_scope_id = _elevated_id
						}
						// get value
						if _index, _buffer_value := (*(*sm.Decls)[_decl_id]).FindValue(_scope_id); _index >= 0 {
							// get size of buffer
							if _size, err := strconv.Atoi(fmt.Sprintf("%v", _buffer_value.Value)); err == nil {
								chan_buffer_size = _size
							} else {
								log.FailureLog("CallExprVarType: unable to determine size of channel buffer: %v", _buffer_value.Value)
							}
						} else {
							log.WarningLog("CallExprVarType: decl found, but value not found: %s,\t%s", _buffer_expr.Name, _scope_id)
						}

					} else {
						log.WarningLog("CallExprVarType: decl not found: %s", _buffer_expr.Name)
					}

				default:
					var_type.Info["CallTrace"] = fmt.Sprintf("%s > %s", var_type.Info["CallTrace"], "Default")
					var_type.Info["BufferSize"] = fmt.Sprintf("unknown: [%v]", _buffer_expr)
					return var_type
				}

				// check buffer size
				if chan_buffer_size == -1 {
					// failed
					log.WarningLog("CallExprVarType: buffer size not determined: %d", chan_buffer_size)

				} else if chan_buffer_size == 0 {
					var_type.Type = VAR_DATA_TYPE_SYNC_CHAN
					var_type.Info["BufferSize"] = fmt.Sprintf("%d", chan_buffer_size)

				} else if chan_buffer_size > 0 {
					var_type.Type = VAR_DATA_TYPE_ASYNC_CHAN
					var_type.Info["BufferSize"] = fmt.Sprintf("%d", chan_buffer_size)

				} else {
					// failed
					log.WarningLog("CallExprVarType: buffer size not understood: %d", chan_buffer_size)
				}
				return var_type

			} else {
				// sync channel
				var_type.Type = VAR_DATA_TYPE_SYNC_CHAN
				var_type.Info["BufferSize"] = "0"
				return var_type
			}

		default:
			var_type.Info["CallTrace"] = fmt.Sprintf("%s > %s", var_type.Info["CallTrace"], "Default")
			var_type.Type = VAR_DATA_TYPE_OTHER
			return var_type
		}

	default:
		var_type.Info["CallTrace"] = fmt.Sprintf("%s > %s", var_type.Info["CallTrace"], "Default")
		var_type.Type = VAR_DATA_TYPE_OTHER
		return var_type
	}
}
