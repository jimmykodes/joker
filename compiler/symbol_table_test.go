package compiler

import "testing"

func TestDefine(t *testing.T) {
	expected := map[string]Symbol{
		"a": {"a", GlobalScope, 0},
		"b": {"b", GlobalScope, 1},
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
}

func TestResolveGlobal(t *testing.T) {
	global := NewSymbolTable()

	global.Define("a")
	global.Define("b")

	expected := []Symbol{
		{"a", GlobalScope, 0},
		{"b", GlobalScope, 1},
	}
	for _, sym := range expected {
		res, ok := global.Resolve(sym.Name)
		if !ok {
			t.Errorf("could not resolve name: %s", sym.Name)
			continue
		}
		if res != sym {
			t.Errorf("invalid value for name %s: got %+v - want %+v", sym.Name, res, sym)
		}
	}
}
