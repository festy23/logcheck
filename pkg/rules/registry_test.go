package rules

import "testing"

func TestNewRegistry(t *testing.T) {
	reg := NewRegistry()
	if got := len(reg.Rules()); got != 4 {
		t.Errorf("NewRegistry() has %d rules, want 4", got)
	}
}

func TestRegistry_Disable(t *testing.T) {
	reg := NewRegistry()
	if !reg.Disable("lowercase") {
		t.Error("Disable(lowercase) returned false, want true")
	}
	if got := len(reg.Rules()); got != 3 {
		t.Errorf("after Disable got %d rules, want 3", got)
	}
	for _, r := range reg.Rules() {
		if r.Name() == "lowercase" {
			t.Error("lowercase rule still present after Disable")
		}
	}
}

func TestRegistry_Disable_NotFound(t *testing.T) {
	reg := NewRegistry()
	if reg.Disable("nonexistent") {
		t.Error("Disable(nonexistent) returned true, want false")
	}
	if got := len(reg.Rules()); got != 4 {
		t.Errorf("rules count changed to %d, want 4", got)
	}
}

func TestRegistry_Rules_IsCopy(t *testing.T) {
	reg := NewRegistry()
	rules := reg.Rules()
	rules[0] = nil
	if reg.Rules()[0] == nil {
		t.Error("Rules() returned the internal slice, not a copy")
	}
}

func TestDefaultRules_Names(t *testing.T) {
	reg := NewRegistry()
	want := map[string]bool{
		"lowercase":    true,
		"english":      true,
		"specialchars": true,
		"sensitive":    true,
	}
	for _, r := range reg.Rules() {
		if !want[r.Name()] {
			t.Errorf("unexpected rule %q", r.Name())
		}
		delete(want, r.Name())
	}
	for name := range want {
		t.Errorf("missing rule %q", name)
	}
}
