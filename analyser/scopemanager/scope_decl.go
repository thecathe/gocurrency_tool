package scopemanager

// ScopeDeclMap
// Label => NewVarDeclID().ID
type ScopeDeclMap map[ID]ID

func NewScopeDeclMap() *ScopeDeclMap {
	return &ScopeDeclMap{}
}
