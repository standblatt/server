package web

import (
	"strings"
	"testing"
)

func TestParseShooterCSV(t *testing.T) {
	in := "\xef\xbb\xbflicense;Firstname;LASTNAME;birthday;postalcode;city;club\n" +
		"13756;Sandra;Schwitter;2.6.1969;7313;St. Margrethenberg;Sargans\n" +
		";Guest;Person;;;;\n" +
		"\n" +
		"1;;Nofirst;;;;\n" +
		"2;Bad;Date;31.02.1990;;;\n"
	rows, errs, err := parseShooterCSV(strings.NewReader(in))
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 {
		t.Fatalf("rows: %+v", rows)
	}
	r := rows[0]
	if r.LicenseNumber != "013756" || r.Birthday != "1969-06-02" || r.City != "St. Margrethenberg" || r.Club != "Sargans" {
		t.Fatalf("row 1: %+v", r)
	}
	if rows[1].LicenseNumber != "" || rows[1].Club != "" {
		t.Fatalf("row 2: %+v", rows[1])
	}
	if len(errs) != 2 || errs[0].Line != 5 || errs[1].Line != 6 {
		t.Fatalf("errors: %+v", errs)
	}

	rows, _, err = parseShooterCSV(strings.NewReader("firstname,lastname,birthday\nA,B,1999-12-31\n"))
	if err != nil || len(rows) != 1 || rows[0].Birthday != "1999-12-31" {
		t.Fatalf("comma file: %+v %v", rows, err)
	}
	if _, _, err = parseShooterCSV(strings.NewReader("name;city\nx;y\n")); err == nil {
		t.Fatal("missing columns accepted")
	}
	if _, _, err = parseShooterCSV(strings.NewReader("")); err == nil {
		t.Fatal("empty file accepted")
	}
}
