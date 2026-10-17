package store_test

import (
	"context"
	"errors"
	"testing"

	"github.com/standblatt/server/internal/store"
	"github.com/standblatt/server/internal/testdb"
)

func names(list []store.Shooter) []string {
	out := []string{}
	for _, s := range list {
		out = append(out, s.Firstname+" "+s.Lastname)
	}
	return out
}

func TestShooterSearch(t *testing.T) {
	st := testdb.New(t)
	ctx := context.Background()
	tid, _ := st.CreateTenant(ctx, "Club")
	other, _ := st.CreateTenant(ctx, "Other")
	club, _ := st.CreateClub(ctx, tid, store.ClubInput{Name: "Sargans Pistol"})

	for _, in := range []store.ShooterInput{
		{LicenseNumber: "013756", Firstname: "Sandra", Lastname: "Schwitter", Birthday: "1969-06-02"},
		{LicenseNumber: "173128", Firstname: "Attilio", Lastname: "Abbiati"},
		{Firstname: "Max", Lastname: "Muster_mann"},
	} {
		clubs := []int64{}
		if in.Firstname == "Sandra" {
			clubs = append(clubs, club.ID)
		}
		if _, err := st.CreateShooter(ctx, tid, in, clubs); err != nil {
			t.Fatal(err)
		}
	}
	st.CreateShooter(ctx, other, store.ShooterInput{LicenseNumber: "013756", Firstname: "Sandra", Lastname: "Other"}, nil)

	cases := []struct {
		q    string
		club int64
		want []string
	}{
		{"13756", 0, []string{"Sandra Schwitter"}},
		{"013756", 0, []string{"Sandra Schwitter"}},
		{"1375", 0, []string{}},
		{"sandra schw", 0, []string{"Sandra Schwitter"}},
		{"schwitter sandra", 0, []string{"Sandra Schwitter"}},
		{"abbi", 0, []string{"Attilio Abbiati"}},
		{"_", 0, []string{"Max Muster_mann"}},
		{"%", 0, []string{}},
		{"", club.ID, []string{"Sandra Schwitter"}},
		{"", 0, []string{"Attilio Abbiati", "Max Muster_mann", "Sandra Schwitter"}},
	}
	for _, c := range cases {
		got, err := st.SearchShooters(ctx, tid, c.q, c.club, 50, 0)
		if err != nil {
			t.Fatal(err)
		}
		if g := names(got); len(g) != len(c.want) || (len(g) > 0 && g[0] != c.want[0]) {
			t.Errorf("search %q club %d: got %v, want %v", c.q, c.club, g, c.want)
		}
	}

	hits, _ := st.SearchShooters(ctx, tid, "013756", 0, 50, 0)
	if len(hits) != 1 || len(hits[0].Clubs) != 1 || hits[0].Clubs[0].Name != "Sargans Pistol" || hits[0].Birthday != "1969-06-02" {
		t.Fatalf("unexpected hit: %+v", hits)
	}
	if _, err := st.CreateShooter(ctx, tid, store.ShooterInput{LicenseNumber: "013756", Firstname: "x", Lastname: "y"}, nil); !errors.Is(err, store.ErrExists) {
		t.Fatalf("duplicate license: %v", err)
	}
	if err := st.AddShooterClub(ctx, other, hits[0].ID, club.ID); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("club link across tenants: %v", err)
	}
	if _, err := st.GetShooter(ctx, other, hits[0].ID); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("shooter across tenants: %v", err)
	}
}

func TestImportShooters(t *testing.T) {
	st := testdb.New(t)
	ctx := context.Background()
	tid, _ := st.CreateTenant(ctx, "Club")

	rows := []store.ImportRow{
		{ShooterInput: store.ShooterInput{LicenseNumber: "000001", Firstname: "Anna", Lastname: "Alpha", Birthday: "1980-01-01", City: "Bern"}, Club: "Club One"},
		{ShooterInput: store.ShooterInput{LicenseNumber: "000001", Firstname: "Anna", Lastname: "Alpha"}, Club: "Club Two"},
		{ShooterInput: store.ShooterInput{Firstname: "Guest", Lastname: "Person", Birthday: "1990-05-05"}, Club: "Club One"},
	}
	res, err := st.ImportShooters(ctx, tid, rows)
	if err != nil {
		t.Fatal(err)
	}
	if res != (store.ImportResult{Created: 2, Updated: 0, ClubsCreated: 2, Links: 3}) {
		t.Fatalf("first import: %+v", res)
	}

	rows[0].City = ""
	rows[0].Lastname = "Alpha-Beta"
	rows[1].Lastname = "Alpha-Beta"
	res, err = st.ImportShooters(ctx, tid, rows)
	if err != nil {
		t.Fatal(err)
	}
	if res != (store.ImportResult{Created: 0, Updated: 2, ClubsCreated: 0, Links: 0}) {
		t.Fatalf("second import: %+v", res)
	}
	got, _ := st.SearchShooters(ctx, tid, "1", 0, 50, 0)
	if len(got) != 1 || got[0].Lastname != "Alpha-Beta" || got[0].City != "Bern" || got[0].Birthday != "1980-01-01" || len(got[0].Clubs) != 2 {
		t.Fatalf("merged shooter: %+v", got)
	}
	clubs, _ := st.ListClubs(ctx, tid)
	if len(clubs) != 2 || clubs[0].ShooterCount != 2 {
		t.Fatalf("clubs: %+v", clubs)
	}
}

func TestSearchShootersPages(t *testing.T) {
	st := testdb.New(t)
	ctx := context.Background()
	tid, _ := st.CreateTenant(ctx, "Club")
	for _, n := range []string{"Amsel", "Buch", "Dachs", "Eiche", "Fink"} {
		if _, err := st.CreateShooter(ctx, tid, store.ShooterInput{Firstname: "Hans", Lastname: n}, nil); err != nil {
			t.Fatal(err)
		}
	}
	page, _ := st.SearchShooters(ctx, tid, "", 0, 2, 2)
	if len(page) != 2 || page[0].Lastname != "Dachs" || page[1].Lastname != "Eiche" {
		t.Fatalf("second page: %+v", page)
	}
	if n, err := st.CountShooters(ctx, tid, "", 0); err != nil || n != 5 {
		t.Fatalf("count: %d %v", n, err)
	}
	if n, _ := st.CountShooters(ctx, tid, "ch", 0); n != 3 {
		t.Fatalf("count with search: %d", n)
	}
}
