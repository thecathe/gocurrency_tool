package scopemanager

import (
	"fmt"
	"go/ast"
	"go/token"
)

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
func (sm *ScopeManager) NewVarType(node ast.Node) VarType {

	var var_type VarType
	var_type.Data = make([]string, 0)

	switch node_type := (node).(type) {

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