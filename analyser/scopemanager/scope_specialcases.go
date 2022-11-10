package scopemanager

// scope id to special scope
// type SpecialCases map[ID]SpecialScope[any]

// func NewSpecialCases() *SpecialCases {
// 	var _nsc SpecialCases = make(SpecialCases)
// 	return &_nsc
// }

// func test() {
// 	_new_special_cases := *NewSpecialCases()

// 	_new_for_loop := NewSpecialCase[ForLoop]()
// 	__new_for_loop := (ForLoop).NewCase(ID(""))
// 	_new_select_recv := NewSpecialCase[SelectRecv]()

// 	_new_special_cases[ID("")] = SpecialScope[any](_new_for_loop)
// }

// func (_special *ForLoop) NewCase(_scope_id ID) *ForLoop {
// 	var ret ForLoop

// 	return &ret
// }

// // func (sc *SpecialCases) NewCase(_scope_id ID, _scope_kind SpecialScopeKind) (bool, *SpecialCases) {

// // 	switch _scope_kind {

// // 	// for loop
// // 	case SPECIAL_SCOPE_TYPE_FOR_LOOP:
// // 		var _new_special_scope SpecialScope[ForLoop] = *NewSpecialScope[ForLoop](_scope_id)

// // 		(*sc)[_scope_id] = _new_special_scope
// // 		return false, sc

// // 	// select recv
// // 	case SPECIAL_SCOPE_TYPE_SELECT:
// // 		var _new_special_scope SpecialScope[SelectRecv] = *NewSpecialScope[SelectRecv](_scope_id)
// // 		(*sc)[_scope_id] = _new_special_scope
// // 		return false, sc

// // 	// unaccounted for
// // 	default:
// // 		return false, sc
// // 	}
// // }

// type SpecSc[_kind SpecialS] string

// type SpecialS interface {
// 	SpecialScopeForLoop | SpecialScopeSelectRecv
// }

// type SpecialScopeForLoop SpecialScope[ForLoop]
// type SpecialScopeSelectRecv SpecialScope[SelectRecv]

// type SpecialScope[_kind SpecialScopes] struct {
// 	ScopeID ID
// 	Kind    _kind
// }

// // func NewSpecialScope[_scope_kind SpecialScopeKinds](_scope_id ID, __scope_kind _scope_kind) *SpecialScope[_scope_kind] {
// // 	_new_special_scope := *NewSpecialCase[_scope_kind]()

// // 	_new_special_scope = _scope_id

// // 	return &_new_special_scope
// // }

// func NewSpecialCase[_scope_kind SpecialScopes]() *SpecialScope[_scope_kind] {
// 	var _new_special_case SpecialScope[_scope_kind]
// 	return &_new_special_case
// }

// type SpecialScopes interface {
// 	ForLoop | SelectRecv | any
// }

type ForLoop struct {
	Init string
	Cond string
	Post string
}

// // func NewSpecialScope[_scope_kind ForLoop](_scope_id ID) *SpecialScope[ForLoop] {

// // 	var _new_special_scope SpecialScope[ForLoop] = SpecialScope[ForLoop]{}

// // 	_new_special_scope.ScopeID = _scope_id

// // 	return &_new_special_scope
// // }

type SelectRecv struct {
	Has struct {
		Default bool
		Timeout bool
	}
	Recv struct {
		Sync  *[]ID
		Async *[]ID
	}
	Branch struct {
		Default  map[string]string
		Timeouts []map[string]string
	}
}

// // type SpecialScopeKind string

// // const (
// // 	SPECIAL_SCOPE_TYPE_NONE    SpecialScopeKind = "None"
// // 	SPECIAL_SCOPE_TYPE_UNKNOWN SpecialScopeKind = "Unknown"

// // 	SPECIAL_SCOPE_TYPE_FOR_LOOP  SpecialScopeKind = "For Loop"
// // 	SPECIAL_SCOPE_TYPE_SELECT    SpecialScopeKind = "Select"
// // 	SPECIAL_SCOPE_TYPE_RECURSION SpecialScopeKind = "Recusion"
// // )
