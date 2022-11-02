package scopemanager

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/thecathe/gocurrency_tool/analyser/log"
)

const (
	log_output_dir = "debug_output_logs"
)

// map of file names for avoiding duplicate files / overwriting
var filenames map[string]int = make(map[string]int, 0)

// combines all other log functions into one json
func (sm *ScopeManager) LogAll() (string, string) {
	var build_string string = ""

	// scope stack
	_, _stack_build_string := (*sm).LogStack()
	build_string = fmt.Sprintf("%s%s", build_string, _stack_build_string)

	// scopes
	_, _all_scopes_build_string := (*sm).LogScopes()
	build_string = fmt.Sprintf("%s%s", build_string, _all_scopes_build_string)

	// decls
	_, _all_decls_build_string := (*sm).LogDecls()
	build_string = fmt.Sprintf("%s%s", build_string, _all_decls_build_string)

	// trim , and end
	build_string = build_string[:len(build_string)-1]

	return "all", build_string
}

// returns log type and json of the current stack
func (sm *ScopeManager) LogStack() (string, string) {
	var build_string = fmt.Sprintf("\"stack\" : {\"size\" : %d, \"scopes\" : [ ", (*sm).Stack.Size())

	build_string = fmt.Sprintf("%s%s,", build_string, (*sm).StringifyStack())

	build_string = fmt.Sprintf("%s ]},", build_string[:len(build_string)-1])

	return "scope_stack", build_string
}

// returns log type and json containing all scopes and their decls
func (sm *ScopeManager) LogScopes() (string, string) {
	var build_string string = fmt.Sprintf("\"scopemap\" : { \"count\" : %d, \"scopes\" : [ ", (*sm).ScopeMap.Size())

	build_string = fmt.Sprintf("%s%s,", build_string, (*sm).StringifyScopes())

	build_string = fmt.Sprintf("%s ]},", build_string[:len(build_string)-1])

	return "scope_decls", build_string
}

// returns log type and json of decls
func (sm *ScopeManager) LogDecls() (string, string) {
	var build_string string = fmt.Sprintf("\"declmap\" : { \"count\" : %d, \"decls\" : [ ", (*sm).Decls.Size())

	build_string = fmt.Sprintf("%s%s,", build_string, (*sm).StringifyDecls())

	build_string = fmt.Sprintf("%s ]},", build_string[:len(build_string)-1])

	return "decls", build_string
}

// returns stringified decls
func (sm *ScopeManager) StringifyDecls() string {
	var build_string string = ""

	for _decl_id := range *(*sm).Decls {
		if _, ok := (*(*sm).Decls)[_decl_id]; ok {
			build_string = fmt.Sprintf("%s%s,", build_string, (*sm).StringifyDecl(_decl_id))
		}
	}
	if len(build_string) > 0 {
		build_string = build_string[:len(build_string)-1]
	}

	return build_string
}

// returns stringified decl from id
func (sm *ScopeManager) StringifyDecl(decl_id ID) string {
	var build_string string = fmt.Sprintf("{\"decl_id\" : \"%v\",\"label\" : \"%s\", ", decl_id, (*(*sm).Decls)[decl_id].Label)

	build_string = fmt.Sprintf("%s\"data_type\" : {\"type\" : \"%v\", ", build_string, (*(*sm).Decls)[decl_id].Type.Type)

	if len((*(*sm).Decls)[decl_id].Type.Data) > 0 {
		build_string = fmt.Sprintf("%s\"data\" : [\"%s\"], \"info\" : [ ", build_string, strings.Join(*(*(*sm).Decls)[decl_id].Type.Data.Slice(), "\", \""))
	} else {
		build_string = fmt.Sprintf("%s\"data\" : [], \"info\" : [ ", build_string)
	}

	for _key, _value := range (*(*sm).Decls)[decl_id].Type.Info {
		build_string = fmt.Sprintf("%s{\"%s\" : \"%s\"},", build_string, _key, _value)
	}
	build_string = fmt.Sprintf("%s], \"values\": [ ", build_string[:len(build_string)-1])

	for _, _value := range (*(*sm).Decls)[decl_id].Values {
		build_string = fmt.Sprintf("%s{\"scope_id\": \"%s\", \"value\": \"%s\", \"pos\": \"%d\"},", build_string, _value.ScopeID, _value.Value, _value.Pos)
	}
	build_string = fmt.Sprintf("%s]}}", build_string[:len(build_string)-1])

	return build_string
}

// stringifies all scopes in stack
func (sm *ScopeManager) StringifyStack() string {
	var build_string string = ""

	for _, _scope_id := range *(*sm).Stack {
		if _, ok := (*(*sm).ScopeMap)[_scope_id]; ok {
			build_string = fmt.Sprintf("%s%s,", build_string, sm.StringifyScope(_scope_id))
		}
	}
	if len(build_string) > 0 {
		build_string = build_string[:len(build_string)-1]
	}

	return build_string
}

// stringifies all scopes
func (sm *ScopeManager) StringifyScopes() string {
	var build_string string = ""

	for _scope_id := range *(*sm).ScopeMap {
		if _, ok := (*(*sm).ScopeMap)[_scope_id]; ok {
			build_string = fmt.Sprintf("%s%s,", build_string, sm.StringifyScope(_scope_id))
		}
	}
	if len(build_string) > 0 {
		build_string = build_string[:len(build_string)-1]
	}

	return build_string
}

// stringifies all decls and elevated decls of given scope id
func (sm *ScopeManager) StringifyScope(scope_id ID) string {

	var _scope_decl_len = len((*(*sm).ScopeMap)[scope_id].Decls)
	var _elevated_decl_len = len(*(*(*sm).ScopeMap)[scope_id].ElevatedIDs)

	var var_list string = fmt.Sprintf("{\"scope_id\" : \"%s\", \"vars_count\" : %d, \"elevated_count\" : %d, \"vars\": [ ", scope_id, _scope_decl_len, _elevated_decl_len)

	// for all scope decls
	for _decl_id := range (*(*sm).ScopeMap)[scope_id].Decls {
		var_list = fmt.Sprintf("%s{ \"decl_id\" : \"%s\", \"values\" : [ ", var_list, _decl_id)
		if _index, _value := (*(*sm.Decls)[_decl_id]).FindValue(scope_id); _index >= 0 {
			var_list = fmt.Sprintf("%s\"%s\"", var_list, _value.Value)
		} else {
			log.WarningLog("StringifyScope: decl found, but value not found: %s,\t%s", _decl_id, scope_id)
		}
		var_list = fmt.Sprintf("%s ]},", var_list)
	}

	// trim last ,
	var_list = fmt.Sprintf("%s ], \"elevated_vars\" : [ ", var_list[:len(var_list)-1])

	// for all elevated decls
	if _elevated_ids, ok := (*sm).GetElevated(scope_id); ok {
		// for each id
		for _, _elevated_id := range *_elevated_ids {
			// for each decl
			for _decl_id := range (*(*sm).ScopeMap)[_elevated_id].Decls {
				var_list = fmt.Sprintf("%s{\"decl_id\" : \"%s\", \"elevated_id\":\"%s\",\"values\" : [ ", var_list, _decl_id, _elevated_id)
				if _index, _value := (*(*sm.Decls)[_decl_id]).FindValue(_elevated_id); _index >= 0 {
					var_list = fmt.Sprintf("%s\"%s\"", var_list, _value.Value)
				} else {
					log.WarningLog("Analyser, NewVarValue, decl found, but value not found: %s,\t%s", _decl_id, _elevated_id)
				}
				var_list = fmt.Sprintf("%s ]},", var_list)
			}
		}
	}
	// trim last ,
	var_list = fmt.Sprintf("%s ]}", var_list[:len(var_list)-1])

	return var_list
}

// prints all decls assigned to a scope, and those elevated to it
func (sm *ScopeManager) PrintScope(scope_id ID) {
	log.DebugLog("Scope Decls:\n%s\n\n\n\n", (*sm).StringifyScope(scope_id))
}

// returns file
func (sm *ScopeManager) ToFile(log_type string, file_content string) {
	if file, ok := CreateLog(log_type); ok {

		if _json, err := json.MarshalIndent(json.RawMessage(fmt.Sprintf("{%s}", file_content)), "", "    "); err != nil {
			log.FailureLog("Failed to marshalindent json:%v\n%s", err, file_content)
		} else {
			file.WriteString(string(_json))
		}

		file.Close()
	} else {
		log.FailureLog("Failed to create file for log: %s", log_type)
	}
}

// creates file and returns it
func CreateLog(log_type string) (*os.File, bool) {
	// make dir if not exists
	os.MkdirAll(log_output_dir, os.ModePerm)
	var _init_file_name string = fmt.Sprintf("%s %s", time.Now().Format("2006-01-02 150405"), log_type)

	if _, ok := filenames[_init_file_name]; ok {
		filenames[_init_file_name] += 1
	} else {
		filenames[_init_file_name] = 0
	}

	var file_name string = fmt.Sprintf("%s %s.json", _init_file_name, fmt.Sprintf("%02d", filenames[_init_file_name]))

	if f, err := os.Create(filepath.Join(log_output_dir, file_name)); err == nil {
		return f, true
	} else {
		log.FailureLog("CreateLog: Failed to create file: %s\n\terror: %v", file_name, err)
	}
	return nil, false
}
