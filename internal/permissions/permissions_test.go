package permissions

import (
	"strings"
	"testing"
)

// TestRegistry keeps the enum in step with All and the defaults within it.
func TestRegistry(t *testing.T) {
	if got := strings.Join(Keys(), ","); got != Enum {
		t.Fatalf("Enum = %q, want %q", Enum, got)
	}
	if got := strings.Join(Scopes(), ","); got != ScopeEnum {
		t.Fatalf("ScopeEnum = %q, want %q", ScopeEnum, got)
	}
	for _, d := range Defaults {
		for _, p := range d.Permissions {
			if !Known(p) {
				t.Errorf("default role %s has unknown %s", d.Key, p)
			}
		}
	}
}

// TestOf checks the admin role and the filtering of unknown keys.
func TestOf(t *testing.T) {
	if len(Of(Admin, nil)) != len(All) {
		t.Fatal("admin has every permission")
	}
	if got := Of(Member, []string{"nope", ReportsView, EventsView}); strings.Join(got, ",") != "events.view,reports.view" {
		t.Fatalf("got %v", got)
	}
	if !Any([]string{EventsView}) || Any([]string{EventsView}, ClubManage) || !Any([]string{EventsView}, ClubManage, EventsView) {
		t.Fatal("Any")
	}
}
