package compiler

import (
	"testing"
)

func TestDefine(t *testing.T) {
	expected := map[string]Symbol{
		"a": {"a", GlobalScope, 0},
		"b": {"b", GlobalScope, 1},
		"c": {"c", LocalScope, 0},
		"d": {"d", LocalScope, 1},
		"e": {"e", LocalScope, 0},
		"f": {"f", LocalScope, 1},
	}
	global := NewSymbolTable()

	a := global.Define("a")
	if a != expected["a"] {
		t.Errorf("invalid value for a: got %v - want %v", a, expected["a"])
	}

	b := global.Define("b")
	if b != expected["b"] {
		t.Errorf("invalid value for a: got %v - want %v", b, expected["b"])
	}

	local1 := NewSymbolTable(OuterSymbolTable(global))

	if c := local1.Define("c"); c != expected["c"] {
		t.Errorf("invalid value for c: got %v - want %v", c, expected["c"])
	}

	if d := local1.Define("d"); d != expected["d"] {
		t.Errorf("invalid value for d: got %v - want %v", d, expected["d"])
	}

	local2 := NewSymbolTable(OuterSymbolTable(local1))

	if e := local2.Define("e"); e != expected["e"] {
		t.Errorf("invalid value for e: got %v - want %v", e, expected["e"])
	}
	if f := local2.Define("f"); f != expected["f"] {
		t.Errorf("invalid value for f: got %v - want %v", f, expected["f"])
	}
}

func TestResolve(t *testing.T) {
	global := NewSymbolTable()
	local1 := NewSymbolTable(OuterSymbolTable(global))
	local2 := NewSymbolTable(OuterSymbolTable(local1))

	global.Define("a")
	global.Define("b")
	local1.Define("c")
	local1.Define("d")
	local2.Define("e")
	local2.Define("f")

	expected := []struct {
		st      *SymbolTable
		symbols []Symbol
	}{
		{
			st: global,
			symbols: []Symbol{
				{"a", GlobalScope, 0},
				{"b", GlobalScope, 1},
			},
		},
		{
			st: local1,
			symbols: []Symbol{
				{"a", GlobalScope, 0},
				{"b", GlobalScope, 1},
				{"c", LocalScope, 0},
				{"d", LocalScope, 1},
			},
		},
		{
			st: local2,
			symbols: []Symbol{
				{"a", GlobalScope, 0},
				{"b", GlobalScope, 1},
				{"c", FreeScope, 0},
				{"d", FreeScope, 1},
				{"e", LocalScope, 0},
				{"f", LocalScope, 1},
			},
		},
	}
	for _, tt := range expected {
		for _, sym := range tt.symbols {
			res, ok := tt.st.Resolve(sym.Name)
			if !ok {
				t.Errorf("could not resolve name: %s", sym.Name)
				continue
			}
			if res != sym {
				t.Errorf("invalid value for name %s: got %+v - want %+v", sym.Name, res, sym)
			}
		}
	}
}

func TestResolveFree(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")
	global.Define("b")

	local1 := NewSymbolTable(OuterSymbolTable(global))
	local1.Define("c")
	local1.Define("d")

	local2 := NewSymbolTable(OuterSymbolTable(local1))
	local2.Define("e")
	local2.Define("f")

	tests := []struct {
		table        *SymbolTable
		expected     []Symbol
		expectedFree []Symbol
	}{
		{
			table: local1,
			expected: []Symbol{
				{"a", GlobalScope, 0},
				{"b", GlobalScope, 1},
				{"c", LocalScope, 0},
				{"d", LocalScope, 1},
			},
			expectedFree: []Symbol{},
		},
		{
			table: local2,
			expected: []Symbol{
				{"a", GlobalScope, 0},
				{"b", GlobalScope, 1},
				{"c", FreeScope, 0},
				{"d", FreeScope, 1},
				{"e", LocalScope, 0},
				{"f", LocalScope, 1},
			},
			expectedFree: []Symbol{
				{"c", LocalScope, 0},
				{"d", LocalScope, 1},
			},
		},
	}
	for _, tt := range tests {
		for _, sym := range tt.expected {
			res, ok := tt.table.Resolve(sym.Name)
			if !ok {
				t.Errorf("could not resolve name: %s", sym.Name)
				continue
			}
			if res != sym {
				t.Errorf("invalid value for symbol %s: got %+v - want %+v", sym.Name, res, sym)
			}
		}
		if len(tt.table.FreeSymbols) != len(tt.expectedFree) {
			t.Errorf("invalid number of free symbols: got %d - want %d", len(tt.table.FreeSymbols), len(tt.expectedFree))
			continue
		}
		for i, sym := range tt.expectedFree {
			res := tt.table.FreeSymbols[i]
			if res != sym {
				t.Errorf("invalid value for free symbol %s: got %+v - want %+v", sym.Name, res, sym)
			}
		}
	}
}
