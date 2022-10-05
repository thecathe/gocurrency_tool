package scopemanager

// GeneralVarType
//
type GeneralVarType string

const (
	VAR_DATA_TYPE_OTHER         GeneralVarType = "Other"
	VAR_DATA_TYPE_INT           GeneralVarType = "Int"
	VAR_DATA_TYPE_STRING        GeneralVarType = "String"
	VAR_DATA_TYPE_STRUCT        GeneralVarType = "Struct"
	VAR_DATA_TYPE_ASYNC_CHANNEL GeneralVarType = "Async. Channel"
	VAR_DATA_TYPE_SYNC_CHANNEL  GeneralVarType = "Sync. Channel"
	VAR_DATA_FUNC_RET           GeneralVarType = "Function Return"
)
