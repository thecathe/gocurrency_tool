package scopemanager

// ParseType
//
type ParseType string

const (
	// general
	PARSE_NONE       ParseType = "None"
	PARSE_SCOPE_EXIT ParseType = "Scope_Exit"
	PARSE_SKIPPED    ParseType = "Skipped"

	// scopes
	PARSE_PACKAGE          ParseType = "Package"
	PARSE_PACKAGE_CONST    ParseType = "Package_Const"
	PARSE_PACKAGE_VAR      ParseType = "Package_Var"
	PARSE_FILE             ParseType = "File"
	PARSE_FILE_IMPORT      ParseType = "File_Import"
	PARSE_FUNC_DECL        ParseType = "Func_Decl"
	PARSE_FUNC_LIT         ParseType = "Func_Lit"
	PARSE_GO_STMT          ParseType = "Go_Stmt"
	PARSE_FOR_STMT         ParseType = "For_Stmt"
	PARSE_RANGE_STMT       ParseType = "Range_Stmt"
	PARSE_IF_STMT          ParseType = "If_Stmt"
	PARSE_SELECT_STMT      ParseType = "Select_Stmt"
	PARSE_SWTICH_STMT      ParseType = "Switch_Stmt"
	PARSE_TYPE_SWITCH_STMT ParseType = "T-Swtich_Stmt"

	// vars
	PARSE_ASSIGN ParseType = "Assign"
	PARSE_DECL   ParseType = "Decl"
	PARSE_DEFINE ParseType = "Define"
	// PARSE_FUNC_CALL_PARAMS ParseType = "None"
	PARSE_FUNC_DECL_PARAMS ParseType = "Func_Decl_Params"

	// fails
	PARSE_FAIL              ParseType = "Fail"
	PARSE_FAIL_ASSIGN_TOKEN ParseType = "Fail_Assign_Tok"
	PARSE_FAIL_DECL_TOKEN   ParseType = "Fail_Decl_Tok"
	PARSE_FAIL_DEFAULT      ParseType = "Fail_Default"
	PARSE_FAIL_FIELD_LIST   ParseType = "Fail_Field_List"
	PARSE_FAIL_GEN_DECL     ParseType = "Fail_Gen_Decl"
	PARSE_FAIL_STACK_PEEK   ParseType = "Fail_Stack_Peek"
	PARSE_FAIL_STACK_POP    ParseType = "Fail_Stack_Pop"
	PARSE_FAIL_VALUE_SPEC   ParseType = "Fail_Value_Spec"
)
