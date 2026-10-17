package agent

import (
	"reflect"
	"testing"
)

func TestParseOptions(t *testing.T) {
	out := "PageSize/Media Size: A4 *A5 A6 Custom.WIDTHxHEIGHT\nInputSlot/Media Source: *Auto Tray1 Tray2\nnot an option line\nDuplex/2-Sided Printing: *None DuplexNoTumble\n"
	got := ParseOptions(out)
	want := []Option{
		{Name: "PageSize", Label: "Media Size", Values: []string{"A4", "A5", "A6", "Custom.WIDTHxHEIGHT"}, Default: "A5"},
		{Name: "InputSlot", Label: "Media Source", Values: []string{"Auto", "Tray1", "Tray2"}, Default: "Auto"},
		{Name: "Duplex", Label: "2-Sided Printing", Values: []string{"None", "DuplexNoTumble"}, Default: "None"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v", got)
	}
}
