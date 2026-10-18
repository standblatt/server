package store

import (
	"testing"

	"github.com/standblatt/server/internal/modules"
)

func TestEveryModuleHasUsage(t *testing.T) {
	for _, k := range modules.Keys(modules.Event) {
		if _, ok := eventUsage[k]; !ok {
			t.Errorf("eventUsage has no entry for %s", k)
		}
	}
	if len(eventUsage) != len(modules.Keys(modules.Event)) {
		t.Errorf("eventUsage has unknown modules")
	}
	for _, k := range modules.Keys(modules.Club) {
		if _, ok := clubUsage[k]; !ok {
			t.Errorf("clubUsage has no entry for %s", k)
		}
	}
}
