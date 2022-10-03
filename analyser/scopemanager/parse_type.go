package scopemanager

// ParseType
//
type ParseType string

const (
	// general
	PARSE_NONE       ParseType = "None"
	PARSE_SCOPE_EXIT ParseType = "None"

	// scopes
	PARSE_PACKAGE          ParseType = "None"
	PARSE_PACKAGE_CONST    ParseType = "None"
	PARSE_PACKAGE_VAR      ParseType = "None"
	PARSE_FILE             ParseType = "None"
	PARSE_FILE_IMPORT      ParseType = "None"
	PARSE_FUNC_DECL        ParseType = "None"
	PARSE_FUNC_LIT         ParseType = "None"
	PARSE_GO_STMT          ParseType = "None"
	PARSE_FOR_STMT         ParseType = "None"
	PARSE_RANGE_STMT       ParseType = "None"
	PARSE_IF_STMT          ParseType = "None"
	PARSE_SELECT_STMT      ParseType = "None"
	PARSE_SWTICH_STMT      ParseType = "None"
	PARSE_TYPE_SWITCH_STMT ParseType = "None"

	// vars
	PARSE_ASSIGN ParseType = "None"
	PARSE_DECL   ParseType = "None"
	PARSE_DEFINE ParseType = "None"
	// PARSE_FUNC_CALL_PARAMS ParseType = "None"
	PARSE_FUNC_DECL_PARAMS ParseType = "None"

	// fails
	PARSE_FAIL              ParseType = "None"
	PARSE_FAIL_ASSIGN_TOKEN ParseType = "None"
	PARSE_FAIL_DECL_TOKEN   ParseType = "None"
	PARSE_FAIL_DEFAULT      ParseType = "None"
	PARSE_FAIL_FIELD_LIST   ParseType = "None"
	PARSE_FAIL_GEN_DECL     ParseType = "None"
	PARSE_FAIL_STACK_PEEK   ParseType = "None"
	PARSE_FAIL_STACK_POP    ParseType = "None"
	PARSE_FAIL_VALUE_SPEC   ParseType = "None"
)
