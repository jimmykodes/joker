package compiler

type SymbolScope string

const (
	GlobalScope = "GLOBAL"
	LocalScope  = "LOCAL"
	FreeScope   = "FREE"
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

	FreeSymbols []Symbol
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

func (s *SymbolTable) defineFree(orig Symbol) Symbol {
	sym := Symbol{
		Name:  orig.Name,
		Scope: FreeScope,
		Index: len(s.FreeSymbols),
	}
	s.FreeSymbols = append(s.FreeSymbols, orig)
	s.store[orig.Name] = sym
	return sym
}

func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	sym, ok := s.store[name]
	if !ok && s.outer != nil {
		sym, ok = s.outer.Resolve(name)
		if !ok {
			return sym, ok
		}
		if sym.Scope == GlobalScope {
			return sym, ok
		}
		free := s.defineFree(sym)
		return free, true
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
