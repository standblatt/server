// Package permissions lists what a role in a club can be allowed to do. Clubs combine these keys into their own
// roles; API operations name the keys they need.
package permissions

import (
	"slices"

	"github.com/standblatt/server/internal/modules"
)

// Keys of the permissions.
const (
	EventsView      = "events.view"
	EventsOperate   = "events.operate"
	EventsConfigure = "events.configure"
	ReportsView     = "reports.view"
	ShootersView    = "shooters.view"
	ShootersEdit    = "shooters.edit"
	ShootersManage  = "shooters.manage"
	RangesManage    = "ranges.manage"
	FormsManage     = "forms.manage"
	ClubManage      = "club.manage"
)

// Permission is a right a role can have; Module names the club module it belongs to, if any.
type Permission struct {
	Key    string `json:"key"`
	Module string `json:"module,omitempty"`
}

// All is every permission, in the order the role matrix shows them.
var All = []Permission{
	{Key: EventsView}, {Key: EventsOperate}, {Key: EventsConfigure}, {Key: ReportsView},
	{Key: ShootersView}, {Key: ShootersEdit}, {Key: ShootersManage},
	{Key: RangesManage}, {Key: FormsManage, Module: modules.Forms}, {Key: ClubManage},
}

// Enum is the comma separated list of permission keys for huma enum tags.
const Enum = "events.view,events.operate,events.configure,reports.view,shooters.view,shooters.edit,shooters.manage,ranges.manage,forms.manage,club.manage"

// Scopes of API tokens beyond the permissions: the owner's own training log.
const (
	TrainingRead  = "training:read"
	TrainingWrite = "training:write"
)

// ScopeEnum is the comma separated list of API token scopes for huma enum tags: the permissions and the training scopes.
const ScopeEnum = Enum + ",training:read,training:write"

// Scopes returns every API token scope: the permission keys and the training scopes.
func Scopes() []string { return append(Keys(), TrainingRead, TrainingWrite) }

// Keys returns every permission key.
func Keys() []string {
	out := make([]string, len(All))
	for i, p := range All {
		out[i] = p.Key
	}
	return out
}

// Known reports whether key is a permission.
func Known(key string) bool {
	return slices.ContainsFunc(All, func(p Permission) bool { return p.Key == key })
}

// Keys of the roles every club starts with.
const (
	Admin  = "admin"
	Board  = "board"
	Staff  = "staff"
	Member = "member"
)

// Defaults are the roles every club starts with and their permissions; the admin role always has every permission.
var Defaults = []struct {
	Key         string
	Permissions []string
}{
	{Admin, nil},
	{Board, []string{EventsView, EventsOperate, EventsConfigure, ReportsView, ShootersView, ShootersEdit, ShootersManage, FormsManage}},
	{Staff, []string{EventsView, EventsOperate, ReportsView, ShootersView, ShootersEdit}},
	{Member, []string{EventsView}},
}

// Of returns the permissions of a role: every one for the admin role, else the known ones of the list in registry order.
func Of(roleKey string, list []string) []string {
	if roleKey == Admin {
		return Keys()
	}
	out := []string{}
	for _, p := range All {
		if slices.Contains(list, p.Key) {
			out = append(out, p.Key)
		}
	}
	return out
}

// Any reports whether have contains one of want; an empty want is always met.
func Any(have []string, want ...string) bool {
	if len(want) == 0 {
		return true
	}
	return slices.ContainsFunc(want, func(w string) bool { return slices.Contains(have, w) })
}
