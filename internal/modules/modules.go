// Package modules lists the building blocks a club or an event can switch on: their keys, levels and dependencies.
// Everything else about a module (its API operations, data, UI) refers to these keys.
package modules

import "slices"

// Level says whether a module belongs to a club or to an event.
type Level string

// Levels of modules.
const (
	Club  Level = "club"
	Event Level = "event"
)

// Keys of the modules.
const (
	Reservations    = "reservations"
	Kiosk           = "kiosk"
	Helpers         = "helpers"
	Classifications = "classifications"
	Attempts        = "attempts"
	Awards          = "awards"
	SpecialPrizes   = "special_prizes"
	ClubCompetition = "club_competition"
	LaneSheet       = "lane_sheet"
	ControlStrip    = "control_strip"
	Ammo            = "ammo"
	Forms           = "forms"
)

// Module is a building block; Requires names modules that must be on as well, NeedsRange is for event modules
// that only make sense with lanes, OffByDefault keeps it off when an API client creates an event without modules.
type Module struct {
	Key          string   `json:"key"`
	Level        Level    `json:"level" enum:"club,event"`
	Requires     []string `json:"requires"`
	NeedsRange   bool     `json:"needsRange"`
	OffByDefault bool     `json:"-"`
}

// All is every module, in the order the UI shows them.
var All = []Module{
	{Key: Reservations, Level: Event, Requires: []string{}, NeedsRange: true},
	{Key: LaneSheet, Level: Event, Requires: []string{}},
	{Key: ControlStrip, Level: Event, Requires: []string{}},
	{Key: Ammo, Level: Event, Requires: []string{}, OffByDefault: true},
	{Key: Attempts, Level: Event, Requires: []string{}},
	{Key: Awards, Level: Event, Requires: []string{}},
	{Key: SpecialPrizes, Level: Event, Requires: []string{Awards}},
	{Key: ClubCompetition, Level: Event, Requires: []string{}},
	{Key: Classifications, Level: Event, Requires: []string{}},
	{Key: Kiosk, Level: Event, Requires: []string{}},
	{Key: Helpers, Level: Event, Requires: []string{}},
	{Key: Forms, Level: Club, Requires: []string{}},
}

// EventEnum is the comma separated list of event module keys for huma enum tags.
const EventEnum = "reservations,lane_sheet,control_strip,ammo,attempts,awards,special_prizes,club_competition,classifications,kiosk,helpers"

// ClubEnum is the comma separated list of club module keys for huma enum tags.
const ClubEnum = "forms"

// Get returns the module of a key.
func Get(key string) (Module, bool) {
	i := slices.IndexFunc(All, func(m Module) bool { return m.Key == key })
	if i < 0 {
		return Module{}, false
	}
	return All[i], true
}

// Keys returns the keys of every module of a level.
func Keys(level Level) []string {
	out := []string{}
	for _, m := range All {
		if m.Level == level {
			out = append(out, m.Key)
		}
	}
	return out
}

// Defaults returns the modules of a level that are on unless chosen otherwise.
func Defaults(level Level) []string {
	out := []string{}
	for _, m := range All {
		if m.Level == level && !m.OffByDefault {
			out = append(out, m.Key)
		}
	}
	return out
}

// Effective returns the modules that are on: the chosen and the used ones with everything they require, in the
// order of All, without the modules that need a range when there is none. Unknown keys are dropped.
func Effective(level Level, chosen, used []string, hasRange bool) []string {
	on := map[string]bool{}
	var add func(string)
	add = func(k string) {
		m, ok := Get(k)
		if !ok || m.Level != level || on[k] {
			return
		}
		on[k] = true
		for _, r := range m.Requires {
			add(r)
		}
	}
	for _, k := range append(slices.Clone(chosen), used...) {
		add(k)
	}
	out := []string{}
	for _, m := range All {
		if on[m.Key] && (hasRange || !m.NeedsRange) {
			out = append(out, m.Key)
		}
	}
	return out
}

// Missing returns the modules of before that are not in after.
func Missing(before, after []string) []string {
	out := []string{}
	for _, k := range before {
		if !slices.Contains(after, k) {
			out = append(out, k)
		}
	}
	return out
}

// Metadata marks a huma operation as part of a module; the server refuses it while the module is off. An empty key
// is a core operation.
func Metadata(key string) map[string]any {
	if key == "" {
		return nil
	}
	return map[string]any{"module": key}
}

// Of returns the module of an operation's metadata, "" for core operations.
func Of(metadata map[string]any) string {
	k, _ := metadata["module"].(string)
	return k
}
