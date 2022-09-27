package scopemanager

type ScopeType string

const (
	SCOPE_TYPE_NONE          ScopeType = "None"
	SCOPE_TYPE_PACKAGE       ScopeType = "Package"              // *ast.Package
	SCOPE_TYPE_FILE          ScopeType = "File"                 // *ast.File
	SCOPE_TYPE_PACKAGE_VAR   ScopeType = "Package Var"          // Peek().Type == File && *ast.GenDecl.tok == var
	SCOPE_TYPE_PACKAGE_CONST ScopeType = "Package Const"        // Peek().Type == File && *ast.GenDecl.tok == const
	SCOPE_TYPE_FILE_IMPORT   ScopeType = "File Import"          // Peek().Type == File && *ast.GenDecl.tok == import
	SCOPE_TYPE_FUNC_CALL     ScopeType = "Function Call"        // *ast.CallExpr
	SCOPE_TYPE_FUNC_DECL     ScopeType = "Function Declaration" // *ast.FuncDecl
	SCOPE_TYPE_IF            ScopeType = "If"                   // *ast.
	SCOPE_TYPE_SELECT        ScopeType = "Select"               // *ast.
	SCOPE_TYPE_SWITCH        ScopeType = "Switch"               // *ast.
	SCOPE_TYPE_TYPE_SWITCH   ScopeType = "Type Switch"          // *ast.
	SCOPE_TYPE_FOR           ScopeType = "For Loop"             // *ast.
	SCOPE_TYPE_RANGE         ScopeType = "Ranged For Loop"      // *ast.
	SCOPE_TYPE_GOROUTINE     ScopeType = "Goroutine"            // *ast.GoStmt
	// SCOPE_TYPE_GO_NAMED      ScopeType = "Goroutine (Named)"     // *ast.
	// SCOPE_TYPE_GO_ANONYMOUS  ScopeType = "Goroutine (Anonymous)" // *ast.
)
