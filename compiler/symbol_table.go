package compiler

type SymbolScope string

const (
	GlobalScope = "GLOBAL"
	LocalScope  = "LOCAL"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	outer          *SymbolTable
	store          map[string]Symbol
	numDefinitions int
}

func NewSymbolTable(options ...SymbolTableOption) *SymbolTable {
	st := &SymbolTable{
		store: make(map[string]Symbol),
	}

	for _, opt := range options {
		opt.Apply(st)
	}

	return st
}

func (s *SymbolTable) Define(name string) Symbol {
	sym := Symbol{Name: name, Scope: GlobalScope, Index: s.numDefinitions}
	if s.outer != nil {
		sym.Scope = LocalScope
	}

	s.store[name] = sym
	s.numDefinitions++
	return sym
}

func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	sym, ok := s.store[name]
	if !ok && s.outer != nil {
		return s.outer.Resolve(name)
	}
	return sym, ok
}

type SymbolTableOption interface {
	Apply(*SymbolTable)
}

type SymbolTableOptionFunc func(*SymbolTable)

func (f SymbolTableOptionFunc) Apply(s *SymbolTable) {
	f(s)
}

func OuterSymbolTable(outer *SymbolTable) SymbolTableOptionFunc {
	return func(s *SymbolTable) {
		s.outer = outer
	}
}
