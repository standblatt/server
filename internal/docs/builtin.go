package docs

import "github.com/standblatt/server/internal/domain"

// Built-in print forms: the documents Standblatt prints unless an event chooses a form of the club instead.
const (
	BuiltinWreathCard   = "wreath_card"
	BuiltinControlStrip = "control_strip"
)

// Builtin returns a built-in print form by key; clubs copy them to change them.
func Builtin(key string) (domain.Layout, bool) {
	switch key {
	case BuiltinWreathCard:
		return domain.Layout{Width: 105, Height: 148, Source: domain.SourceResults, Elements: []domain.Element{
			{Type: "text", X: 30, Y: 94.5, W: 65, H: 6, Size: 10, Bold: true, Value: "{{shooter.lastname}}"},
			{Type: "text", X: 30, Y: 102.7, W: 50, H: 6, Size: 10, Value: "{{shooter.firstname}}"},
			{Type: "text", X: 85, Y: 102.7, W: 20, H: 6, Size: 10, Value: "{{shooter.birthyear}}"},
			{Type: "text", X: 30, Y: 111, W: 65, H: 5, Size: 10, Value: "{{club}}"},
			{Type: "text", X: 30, Y: 119.3, W: 65, H: 5, Size: 10, Value: "{{event.date}}"},
			{Type: "text", X: 15, Y: 131, W: 80, H: 4, Size: 8, Align: "R", Value: "{{event.organizer}}"},
			{Type: "text", X: 30, Y: 135.3, W: 65, H: 5, Size: 12, Bold: true, Value: "{{award.value}}"},
			{Type: "text", X: 30, Y: 140.3, W: 65, H: 4, Size: 8, Value: "{{award.value.words}}"},
		}}, true
	case BuiltinControlStrip:
		return domain.Layout{Width: 90.3, Height: 29, Source: domain.SourceRegistration, Elements: []domain.Element{
			{Type: "text", X: 2, Y: 1.5, W: 57, H: 3, Size: 7, Italic: true, Value: "{{event.name}}"},
			{Type: "text", X: 2, Y: 5, W: 57, H: 4, Size: 9, Bold: true, Value: "{{shooter.lastname}} {{shooter.firstname}}, {{shooter.birthday}}"},
			{Type: "text", X: 2, Y: 9.5, W: 57, H: 3.5, Size: 8, Value: "{{shooter.license}}"},
			{Type: "text", X: 2, Y: 13.5, W: 57, H: 3.5, Size: 8, Value: "{{club}}"},
			{Type: "text", X: 2, Y: 20, W: 12, H: 6, Size: 13, Bold: true, Value: "{{category.short}}"},
			{Type: "text", X: 16, Y: 20, W: 43, H: 3, Size: 7, Wrap: true, Value: "{{passes}}"},
			{Type: "qr", X: 63, Y: 1.5, W: 24, Value: "{{registration.code}}"},
			{Type: "text", X: 62, Y: 25.5, W: 26, H: 2.5, Size: 5, Bold: true, Align: "C", Value: "ID {{registration}}"},
		}}, true
	}
	return domain.Layout{}, false
}

// BuiltinKeys lists the built-in print forms.
var BuiltinKeys = []string{BuiltinWreathCard, BuiltinControlStrip}
