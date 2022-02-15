package scopemanager

type ScopeType string

const (
	SCOPE_TYPE_NONE         ScopeType = "None"
	SCOPE_TYPE_ROOT         ScopeType = "Root"
	SCOPE_TYPE_FILE         ScopeType = "File"
	SCOPE_TYPE_GEN_DECL     ScopeType = "Decl"
	SCOPE_TYPE_FUNC_CALL    ScopeType = "Function Call"
	SCOPE_TYPE_FUNC_DECL    ScopeType = "Function Declaration"
	SCOPE_TYPE_IF           ScopeType = "If"
	SCOPE_TYPE_SELECT       ScopeType = "Select"
	SCOPE_TYPE_SWITCH       ScopeType = "Switch"
	SCOPE_TYPE_TYPE_SWITCH  ScopeType = "Type Switch"
	SCOPE_TYPE_FOR          ScopeType = "For Loop"
	SCOPE_TYPE_RANGE        ScopeType = "Ranged For Loop"
	SCOPE_TYPE_GO_NAMED     ScopeType = "Goroutine (Named)"
	SCOPE_TYPE_GO_ANONYMOUS ScopeType = "Goroutine (Anonymous)"
)
