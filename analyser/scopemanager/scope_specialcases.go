package scopemanager

type SpecialScope[S SpecialScopes] struct {
	ScopeID ID
	Case    S
}

type SpecialScopes interface {
	ForLoop | ~string
}

type SpecialCases map[ID]SpecialScope

type ForLoop struct {
	Init string
	Cond string
	Post string
}
