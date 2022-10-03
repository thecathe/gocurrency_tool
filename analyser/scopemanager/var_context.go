package scopemanager

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