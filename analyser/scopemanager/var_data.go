package scopemanager

import (
	"go/ast"
	"go/token"
)

// GeneralVarType
type GeneralVarType string

const (
	VAR_DATA_TYPE_NONE       GeneralVarType = "None"
	VAR_DATA_TYPE_OTHER      GeneralVarType = "Other"
	VAR_DATA_TYPE_INT        GeneralVarType = "Int"
	VAR_DATA_TYPE_STRING     GeneralVarType = "String"
	VAR_DATA_TYPE_BOOL       GeneralVarType = "Bool"
	VAR_DATA_TYPE_STRUCT     GeneralVarType = "Struct"
	VAR_DATA_TYPE_CHAN       GeneralVarType = "Chan"
	VAR_DATA_TYPE_ASYNC_CHAN GeneralVarType = "Async. Channel"
	VAR_DATA_TYPE_SYNC_CHAN  GeneralVarType = "Sync. Channel"
	VAR_DATA_FUNC_RET        GeneralVarType = "Function Return"
)

func DataStringToVarType(_data_type string) GeneralVarType {
	switch _data_type {
	case "int":
		return VAR_DATA_TYPE_INT

	case "string":
		return VAR_DATA_TYPE_STRING

	case "bool":
		return VAR_DATA_TYPE_BOOL

	default:
		return VAR_DATA_TYPE_OTHER
	}
}

func TokKindToVarType(_tok_kind token.Token) GeneralVarType {
	switch _tok_kind {
	case token.INT:
		return VAR_DATA_TYPE_INT

	case token.STRING:
		return VAR_DATA_TYPE_STRING

	case token.STRUCT:
		return VAR_DATA_TYPE_STRUCT

	default:
		return VAR_DATA_TYPE_OTHER
	}
}

func (sm *ScopeManager) NodeToVarType(_node ast.Node) GeneralVarType {
	switch _n := _node.(type) {

	case *ast.ChanType:
		return VAR_DATA_TYPE_CHAN

	case *ast.Ident:
		return DataStringToVarType(_n.Name)

	default:
		return VAR_DATA_TYPE_OTHER
	}
}

type CompoundVarType []GeneralVarType

func NewCompoundVarType() CompoundVarType {
	var cvt CompoundVarType = make(CompoundVarType, 0)
	return cvt
}

func (cvt *CompoundVarType) Slice() *[]string {
	var slice []string = make([]string, 0)
	for _, _vt := range *cvt {
		slice = append(slice, string(_vt))
	}
	return &slice
}

func (cvt *CompoundVarType) Get(index int) GeneralVarType {
	for _index, _vt := range *cvt {
		if index == _index {
			return _vt
		}
	}
	return GeneralVarType("")
}

func CompoundVarTypeBuilder(node ast.Node) CompoundVarType {
	var cvt = make(CompoundVarType, 0)

	switch _var_type := node.(type) {

	case *ast.ChanType:
		cvt = append(cvt, VAR_DATA_TYPE_CHAN)
		cvt = append(cvt, CompoundVarTypeBuilder(_var_type.Value)...)
		return cvt

	case *ast.Ident:
		cvt = append(cvt, DataStringToVarType(_var_type.Name))
		return cvt

	case *ast.BasicLit:
		cvt = append(cvt, TokKindToVarType(_var_type.Kind))
		return cvt

	default:
		return []GeneralVarType{VAR_DATA_TYPE_OTHER}

	}
}
