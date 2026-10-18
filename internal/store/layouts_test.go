package store_test

import (
	"context"
	"errors"
	"testing"

	"github.com/standblatt/server/internal/domain"
	"github.com/standblatt/server/internal/store"
	"github.com/standblatt/server/internal/testdb"
)

func TestLayouts(t *testing.T) {
	st := testdb.New(t)
	ctx := context.Background()
	a, _ := st.CreateTenant(ctx, "A")
	b, _ := st.CreateTenant(ctx, "B")
	in := store.LayoutInput{Name: "Army sheet", Layout: domain.Layout{Width: 148, Height: 210, Elements: []domain.Element{{Type: "text", X: 1, Y: 2, Value: "{{total}}"}}}}
	id, err := st.CreateLayout(ctx, a, in)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := st.CreateLayout(ctx, a, in); !errors.Is(err, store.ErrExists) {
		t.Fatalf("same name twice: %v", err)
	}
	if err := st.SetLayoutBackground(ctx, a, id, []byte{1, 2}); err != nil {
		t.Fatal(err)
	}
	l, err := st.GetLayout(ctx, a, id)
	bg, _ := st.LayoutBackground(ctx, a, id)
	if err != nil || !l.HasBackground || len(bg) != 2 || l.Elements[0].Value != "{{total}}" {
		t.Fatalf("layout: %+v %v %v", l, bg, err)
	}
	if _, err := st.GetLayout(ctx, b, id); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("layout of another tenant: %v", err)
	}
	if err := st.UpdateLayout(ctx, b, id, in); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("update across tenants: %v", err)
	}
	if err := st.DeleteLayout(ctx, b, id); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("delete across tenants: %v", err)
	}
	if list, _ := st.ListLayouts(ctx, b); len(list) != 0 {
		t.Fatalf("tenant B sees layouts of A: %+v", list)
	}
	rng, _ := st.CreateRange(ctx, b, store.RangeInput{Name: "10m", LaneCount: 2})
	ev, _ := st.CreateEvent(ctx, b, store.EventInput{RangeIDs: []int64{rng.ID}, Name: "B event", StartsOn: "2027-01-01", Timezone: "UTC"})
	p, _ := st.CreateProgramme(ctx, b, ev, fieldMatch())
	if err := st.SetProgrammeLayout(ctx, b, p, &id); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("form of another tenant on a programme: %v", err)
	}
}
