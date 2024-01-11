package compiler

type SymbolScope string

const (
	GlobalScope = "GLOBAL"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	store          map[string]Symbol
	numDefinitions int
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		store: make(map[string]Symbol),
	}
}

func (s *SymbolTable) Define(name string) Symbol {
	sym := Symbol{Name: name, Scope: GlobalScope, Index: s.numDefinitions}
	s.store[name] = sym
	s.numDefinitions++
	return sym
}

func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	sym, ok := s.store[name]
	return sym, ok
}
