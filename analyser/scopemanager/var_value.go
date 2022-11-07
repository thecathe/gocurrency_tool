package scopemanager

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"

	"github.com/thecathe/gocurrency_tool/analyser/log"
)

// VarValue
type VarValue struct {
	Value   string
	Pos     token.Pos
	ScopeID ID
	Trace   []string
}

// Returns VarValue using ast.ValueSpec .Values[]ast.Expr and .Pos
func (sm *ScopeManager) NewVarValue(expr ast.Expr, pos token.Pos) (VarValue, VarType) {
	var value VarValue
	value.Value = "unknown"
	value.Trace = make([]string, 0)

	var _var_type VarType = NewVarType()
	_var_type.Type = VAR_DATA_TYPE_NONE

	value.Pos = pos
	if scope_id, ok := (*sm).PeekID(); ok {
		value.ScopeID = scope_id
	}

	switch value_expr := expr.(type) {

	// simple value
	case *ast.BasicLit:
		value.Trace = append(value.Trace, "*ast.BasicLit")
		value.Value = fmt.Sprintf("%v", value_expr.Value)

	case *ast.Ident:
		value.Trace = append(value.Trace, "*ast.Ident")
		_temp_value, _temp_var_type := (*sm).IdentVarValue(value_expr, &value, &_var_type)
		value = *_temp_value
		_var_type = *_temp_var_type

	// add as is
	default:
		value.Trace = append(value.Trace, "Default")

		switch inner_expr := value_expr.(type) {

		// Type from Function
		case *ast.CallExpr:
			value.Trace = append(value.Trace, "*ast.CallExpr")
			_var_type = *(*sm).CallExprVarType(&_var_type, inner_expr)
			// patch var value
			if _var_type.Type == VAR_DATA_TYPE_ASYNC_CHAN || _var_type.Type == VAR_DATA_TYPE_SYNC_CHAN || _var_type.Type == VAR_DATA_TYPE_CHAN {
				value.Value = string(_var_type.Data.Get(len(_var_type.Data) - 1))
			}

		//
		case *ast.BinaryExpr:
			value.Trace = append(value.Trace, "*ast.BinaryExpr")
			_temp_value, _temp_type, _result_expr := (*sm).BinaryExprVarValue(inner_expr, &value, &_var_type)
			value = *_temp_value
			_var_type = *_temp_type

			if ok, _temp_value, _temp_type := _result_expr.Evaluate(); ok {
				value.Value = _temp_value
				_var_type.Type = _temp_type
			} else {
				value.Value = _temp_value
				log.FailureLog("BinaryExprVarValue, Evaluation failed yielding: \"%s\"", _temp_value)
			}

		//
		case *ast.UnaryExpr:
			value.Trace = append(value.Trace, "*ast.Unary")
			_temp_value, _temp_type := (*sm).UnaryExprVarValue(inner_expr, &value, &_var_type)
			value = *_temp_value
			_var_type = *_temp_type

			// add op to beginning
			// value.Value = fmt.Sprintf("%s%s", inner_expr.Op.String(), value.Value)

		//
		default:
			value.Trace = append(value.Trace, "Default")
			value.Value = "Other"
		}

	}

	// remove surrounding quotes
	if len(value.Value) >= 2 {
		if value.Value[:1] == "\"" {
			value.Value = value.Value[1:]
		}

		if value.Value[len(value.Value)-1:] == "\"" {
			value.Value = value.Value[:len(value.Value)-1]
		}
	}

	log.DebugLog("NewVarValue: %s", value.Value)
	return value, _var_type
}

type OpVarValue struct {
	Type  GeneralVarType
	Value string
}

type BinaryVarValue struct {
	Lhs OpVarValue
	Op  token.Token
	Rhs OpVarValue
}

// extract unary operation
func (sm *ScopeManager) UnaryExprVarValue(_node *ast.UnaryExpr, _value *VarValue, _type *VarType) (*VarValue, *VarType) {

	_value.Value = "UnaryExpr"

	// get type
	switch inner_type := _node.X.(type) {

	case *ast.Ident:
		_type.Info["ValueTrace"] = fmt.Sprintf("%s > %s", _type.Info["ValueTrace"], "*ast.Ident")
		_temp_value, _temp_var_type := (*sm).IdentVarValue(inner_type, _value, _type)
		_value = _temp_value
		_type = _temp_var_type

		_value.Trace = append(_value.Trace, fmt.Sprintf("*ast.Ident (%s: %s)", _value.Value, _type.Type))

	case *ast.BasicLit:
		_type.Info["ValueTrace"] = fmt.Sprintf("%s > %s", _type.Info["ValueTrace"], "*ast.BasicLit")
		_value.Value = inner_type.Value
		_type.Type = TokKindToVarType(inner_type.Kind)

		// add op back
		_value.Value = fmt.Sprintf("%s%s", _node.Op.String(), _value.Value)

		_value.Trace = append(_value.Trace, fmt.Sprintf("*ast.BasicLit (%s: %s)", _value.Value, _type.Type))

	// Type from Function
	case *ast.CallExpr:
		_type.Info["ValueTrace"] = fmt.Sprintf("%s > %s", _type.Info["ValueTrace"], "*ast.CallExpr")
		_type = (*sm).CallExprVarType(_type, inner_type)
		// patch var value
		if _type.Type == VAR_DATA_TYPE_ASYNC_CHAN || _type.Type == VAR_DATA_TYPE_SYNC_CHAN || _type.Type == VAR_DATA_TYPE_CHAN {
			_value.Value = string(_type.Data.Get(len(_type.Data) - 1))
		}

		// add op back
		_value.Value = fmt.Sprintf("%s%s", _node.Op.String(), _value.Value)

		_value.Trace = append(_value.Trace, fmt.Sprintf("*ast.CallExpr (%s: %s)", _value.Value, _type.Type))

	//
	case *ast.BinaryExpr:
		_type.Info["ValueTrace"] = fmt.Sprintf("%s > %s", _type.Info["ValueTrace"], "*ast.BinaryExpr")
		_temp_value, _temp_type, _result_expr := (*sm).BinaryExprVarValue(inner_type, _value, _type)
		_value = _temp_value
		_type = _temp_type

		if ok, _temp_value, _temp_type := _result_expr.Evaluate(); ok {
			_value.Value = _temp_value
			_type.Type = _temp_type
		} else {
			_value.Value = _temp_value
			log.FailureLog("BinaryExprVarValue, Evaluation failed yielding: \"%s\"", _temp_value)
		}

		// add op back
		_value.Value = fmt.Sprintf("%s%s", _node.Op.String(), _value.Value)

		_value.Trace = append(_value.Trace, fmt.Sprintf("*ast.BinaryExpr (%s: %s)", _value.Value, _type.Type))
	}
	return _value, _type
}

func (expr *BinaryVarValue) Evaluate() (bool, string, GeneralVarType) {
	if _ok, _result, _gvt := expr.DevEvaluate(); _ok {
		log.DebugLog("[%s(%s) %s %s(%s)] = %s", expr.Lhs.Type, expr.Lhs.Value, expr.Op.String(), expr.Rhs.Type, expr.Rhs.Value, _result)

		return _ok, _result, _gvt

	} else {
		return _ok, _result, _gvt
	}
}

// evaluates a binary expression
// bool true means string has a value
// bool false means string only contains working out, eg: "int(5) + string(three)"
func (expr *BinaryVarValue) DevEvaluate() (bool, string, GeneralVarType) {

	// log.DebugLog("[%s(%s) %s %s(%s)]", expr.Lhs.Type, expr.Lhs.Value, expr.Op.String(), expr.Rhs.Type, expr.Rhs.Value)

	var result string

	// if both same data type
	if expr.Lhs.Type == expr.Rhs.Type {
		var result string

		// for resulting in bool
		switch expr.Op {

		// (==)
		case token.EQL:
			// compare
			if expr.Lhs.Value == expr.Rhs.Value {
				return true, "true", VAR_DATA_TYPE_BOOL
			} else {
				return true, "false", VAR_DATA_TYPE_BOOL
			}

		// (!=)
		case token.NEQ:
			// compare
			if expr.Lhs.Value == expr.Rhs.Value {
				return true, "false", VAR_DATA_TYPE_BOOL
			} else {
				return true, "true", VAR_DATA_TYPE_BOOL
			}

		// non trivial bool
		default:

			switch expr.Lhs.Type {

			// bool
			case VAR_DATA_TYPE_BOOL:
				switch expr.Op {
				//
				case token.LAND:
					// compare
					if expr.Lhs.Value == expr.Rhs.Value {
						return true, expr.Lhs.Value, VAR_DATA_TYPE_BOOL
					} else {
						return true, "false", VAR_DATA_TYPE_BOOL
					}

				//
				case token.LOR:
					// compare
					if expr.Lhs.Value == "true" || expr.Rhs.Value == "true" {
						return true, "true", VAR_DATA_TYPE_BOOL
					} else {
						return true, "false", VAR_DATA_TYPE_BOOL
					}

				// unaccounted for
				default:
					result = fmt.Sprintf("[%s(%s) %s %s(%s)]", expr.Lhs.Type, expr.Lhs.Value, expr.Op.String(), expr.Rhs.Type, expr.Rhs.Value)
					return false, result, VAR_DATA_TYPE_OTHER
				}

			// string
			case VAR_DATA_TYPE_STRING:

				switch expr.Op {
				//
				case token.ADD:
					return true, fmt.Sprintf("%s%s", expr.Lhs.Value, expr.Rhs.Value), VAR_DATA_TYPE_STRING

				// unaccounted for
				default:
					result = fmt.Sprintf("[%s(%s) %s %s(%s)]", expr.Lhs.Type, expr.Lhs.Value, expr.Op.String(), expr.Rhs.Type, expr.Rhs.Value)
					return false, result, VAR_DATA_TYPE_OTHER
				}

			// floats, ints & unaccounted for
			default:
				// int or float
				if expr.Lhs.Type == VAR_DATA_TYPE_INT || expr.Lhs.Type == VAR_DATA_TYPE_FLOAT {

					_i_lhs, _l_err := strconv.Atoi(expr.Lhs.Value)
					_i_rhs, _r_err := strconv.Atoi(expr.Rhs.Value)

					if _l_err == nil && _r_err == nil {
						// use op
						switch expr.Op {
						//
						case token.ADD:
							return true, fmt.Sprintf("%d", _i_lhs+_i_rhs), expr.Lhs.Type

						//
						case token.SUB:
							return true, fmt.Sprintf("%d", _i_lhs-_i_rhs), expr.Lhs.Type

						//
						case token.MUL:
							return true, fmt.Sprintf("%d", _i_lhs*_i_rhs), expr.Lhs.Type

						//
						case token.QUO:
							return true, fmt.Sprintf("%d", _i_lhs/_i_rhs), VAR_DATA_TYPE_FLOAT

						//
						case token.REM:
							return true, fmt.Sprintf("%d", _i_lhs%_i_rhs), expr.Lhs.Type

						// (<)
						case token.LSS:
							return true, fmt.Sprintf("%t", _i_lhs < _i_rhs), VAR_DATA_TYPE_BOOL

						// (>)
						case token.GTR:
							return true, fmt.Sprintf("%t", _i_lhs > _i_rhs), VAR_DATA_TYPE_BOOL

						// (<=)
						case token.LEQ:
							return true, fmt.Sprintf("%t", _i_lhs <= _i_rhs), VAR_DATA_TYPE_BOOL

						// (>=)
						case token.GEQ:
							return true, fmt.Sprintf("%t", _i_lhs >= _i_rhs), VAR_DATA_TYPE_BOOL

						//
						case token.EQL:
							return true, fmt.Sprintf("%t", _i_lhs == _i_rhs), VAR_DATA_TYPE_BOOL

						//
						case token.NEQ:
							return true, fmt.Sprintf("%t", _i_lhs != _i_rhs), VAR_DATA_TYPE_BOOL

						// unaccounted for
						default:
							result = fmt.Sprintf("[%s(%s) %s %s(%s)]", expr.Lhs.Type, expr.Lhs.Value, expr.Op.String(), expr.Rhs.Type, expr.Rhs.Value)
							return false, result, VAR_DATA_TYPE_OTHER
						}

					} else {
						result = fmt.Sprintf("[%s(%s) %s %s(%s)]", expr.Lhs.Type, expr.Lhs.Value, expr.Op.String(), expr.Rhs.Type, expr.Rhs.Value)
						return false, result, VAR_DATA_TYPE_OTHER
					}
				} else {
					// unaccounted for
					result = fmt.Sprintf("[%s(%s) %s %s(%s)]", expr.Lhs.Type, expr.Lhs.Value, expr.Op.String(), expr.Rhs.Type, expr.Rhs.Value)
					return false, result, VAR_DATA_TYPE_OTHER
				}
			}

		}
	} else {
		// format unknown
		result = fmt.Sprintf("[%s(%s) %s %s(%s)]", expr.Lhs.Type, expr.Lhs.Value, expr.Op.String(), expr.Rhs.Type, expr.Rhs.Value)
		return false, result, VAR_DATA_TYPE_OTHER
	}
}

// recursively expands a binary expression to a list of values to be evaluated
func (sm *ScopeManager) BinaryExprVarValue(_node *ast.BinaryExpr, _value *VarValue, _type *VarType) (*VarValue, *VarType, *BinaryVarValue) {

	var _binary_expr BinaryVarValue

	_binary_expr.Op = _node.Op

	// for either branch
	for _index, _branches := range []ast.Expr{_node.X, _node.Y} {

		var _temp_op_var OpVarValue

		// evaluate branch
		switch _branch := _branches.(type) {

		// continue expanding
		case *ast.ParenExpr:
			_value.Trace = append(_value.Trace, fmt.Sprintf("%d : *ast.ParenExpr", _index))

			// switch for paren
			switch _branch_p := _branch.X.(type) {

			case *ast.BinaryExpr:
				_value.Trace = append(_value.Trace, fmt.Sprintf("%d : *ast.BinaryExpr", _index))
				return (*sm).BinaryExprVarValue(_branch_p, _value, _type)

			default:
				_value.Trace = append(_value.Trace, fmt.Sprintf("%d : Default", _index))

				_temp_op_var.Value = ""
				_temp_op_var.Type = VAR_DATA_TYPE_OTHER

			}
		// split x y
		case *ast.BinaryExpr:
			// _value.Trace = append(_value.Trace, fmt.Sprintf("%d : *ast.BinaryExpr", _index))

			_temp_value, _temp_type, _temp_binary_expr := (*sm).BinaryExprVarValue(_branch, _value, _type)

			_value = _temp_value
			_type = _temp_type

			if ok, _evaluated_binary_expr, _resulting_type := _temp_binary_expr.Evaluate(); ok {
				_temp_op_var.Value = _evaluated_binary_expr
				_temp_op_var.Type = _resulting_type

				_value.Trace = append(_value.Trace, fmt.Sprintf("%d : *ast.BinaryExpr | %s", _index, _evaluated_binary_expr))
			} else {
				_temp_op_var.Value = _evaluated_binary_expr
				_temp_op_var.Type = _resulting_type
				_value.Trace = append(_value.Trace, fmt.Sprintf("%d : *ast.BinaryExpr | %s", _index, _evaluated_binary_expr))
				log.FailureLog("BinaryExprVarValue, Evaluation failed yielding: \"%s\"", _evaluated_binary_expr)
			}

		//
		case *ast.UnaryExpr:
			_temp_value, _temp_type := (*sm).UnaryExprVarValue(_branch, _value, _type)
			_temp_op_var.Value = _temp_value.Value
			_temp_op_var.Type = _temp_type.Type

			_value.Trace = append(_value.Trace, fmt.Sprintf("%d : *ast.UnaryExpr (%s: %s)", _index, _temp_op_var.Type, _temp_op_var.Value))

		// basic lit
		case *ast.Ident:
			_temp_value, _temp_type := (*sm).IdentVarValue(_branch, _value, _type)
			_temp_op_var.Value = _temp_value.Value
			_temp_op_var.Type = _temp_type.Type

			_value.Trace = append(_value.Trace, fmt.Sprintf("%d : *ast.Ident (%s: %s)", _index, _temp_op_var.Type, _temp_op_var.Value))

		// get from defined var
		case *ast.BasicLit:
			_temp_op_var.Value = _branch.Value
			_temp_op_var.Type = TokKindToVarType(_branch.Kind)

			_value.Trace = append(_value.Trace, fmt.Sprintf("%d : *ast.BasicLit (%s: %s)", _index, _temp_op_var.Type, _temp_op_var.Value))

		default:
			_value.Trace = append(_value.Trace, fmt.Sprintf("%d : Default", _index))

			_temp_op_var.Value = ""
			_temp_op_var.Type = VAR_DATA_TYPE_OTHER

		}

		// add to correct branch
		switch _index {
		case 0:
			_binary_expr.Lhs = _temp_op_var
		case 1:
			_binary_expr.Rhs = _temp_op_var
		default:
			log.FailureLog("BinaryExprVarValue: unexpected branch: %d", _index)
		}

	}

	return _value, _type, &_binary_expr
}

// returns the value and type retrieved from a given ident when used to define the value of a var assignment
func (sm *ScopeManager) IdentVarValue(_node *ast.Ident, _value *VarValue, _type *VarType) (*VarValue, *VarType) {

	_value.Value = fmt.Sprintf("%v", _node.Name)
	// check bool
	if _node.Name == "true" {
		_type.Type = VAR_DATA_TYPE_BOOL
	} else if _node.Name == "false" {
		_type.Type = VAR_DATA_TYPE_BOOL
	} else {
		// find corresponding decl
		if _decl_id, _scope_id, _elevated, _elevated_id := (*sm).FindDeclID(_node.Name); _scope_id != "" {
			_type.Type = (*(*sm).Decls)[_decl_id].Type.Type

			if _elevated {
				_scope_id = _elevated_id
			}

			if _index, _var_value := (*(*sm).Decls)[_decl_id].FindValue(_scope_id); _index >= 0 {
				_value.Value = fmt.Sprintf("%v", _var_value.Value)
			} else {
				_value.Value = fmt.Sprintf("ident: %v", _node.Name)
				log.FailureLog("NewVarValue; *ast.Ident: label: \"%s\" did not yield a decl_id", _node.Name)
			}
		} else {
			_value.Value = fmt.Sprintf("ident: %v", _node.Name)
			log.FailureLog("NewVarValue; *ast.Ident: label: \"%s\" did not yield a decl_id", _node.Name)
		}
	}

	return _value, _type
}
