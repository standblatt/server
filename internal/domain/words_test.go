package domain

import "testing"

func TestAmountInWords(t *testing.T) {
	cases := map[string]map[float64]string{
		"en": {
			0:       "Zero francs",
			1:       "One franc",
			25:      "Twenty-five francs",
			21.5:    "Twenty-one francs and fifty centimes",
			100:     "One hundred francs",
			105.05:  "One hundred and five francs and five centimes",
			1000:    "One thousand francs",
			1234:    "One thousand two hundred and thirty-four francs",
			2005:    "Two thousand and five francs",
			0.01:    "Zero francs and one centime",
			1234567: "CHF 1234567.00",
		},
		"de": {
			1:      "Ein Franken",
			16:     "Sechzehn Franken",
			25:     "Fünfundzwanzig Franken",
			30.5:   "Dreissig Franken und fünfzig Rappen",
			101:    "Einhundertein Franken",
			166:    "Einhundertsechsundsechzig Franken",
			1000:   "Eintausend Franken",
			21300:  "Einundzwanzigtausenddreihundert Franken",
			100.01: "Einhundert Franken und ein Rappen",
		},
		"fr": {
			1:      "Un franc",
			21:     "Vingt et un francs",
			71:     "Septante et un francs",
			80:     "Quatre-vingts francs",
			81:     "Quatre-vingt-un francs",
			95.2:   "Nonante-cinq francs et vingt centimes",
			100:    "Cent francs",
			200:    "Deux cents francs",
			250:    "Deux cent cinquante francs",
			1000:   "Mille francs",
			200000: "Deux cent mille francs",
			80000:  "Quatre-vingt mille francs",
		},
		"it": {
			1:    "Un franco",
			21:   "Ventuno franchi",
			23:   "Ventitré franchi",
			38:   "Trentotto franchi",
			100:  "Cento franchi",
			250:  "Duecentocinquanta franchi",
			108:  "Centotto franchi",
			180:  "Centottanta franchi",
			380:  "Trecentottanta franchi",
			1000: "Mille franchi",
			2005: "Duemilacinque franchi",
			15.5: "Quindici franchi e cinquanta centesimi",
		},
	}
	for locale, values := range cases {
		for v, want := range values {
			if got := AmountInWords(locale, v); got != want {
				t.Errorf("%s %.2f: got %q, want %q", locale, v, got, want)
			}
		}
	}
	if got := AmountInWords("xx", 2); got != "Two francs" {
		t.Errorf("unknown locale: %q", got)
	}
}

func TestPassesText(t *testing.T) {
	main := int64(1)
	got := PassesText([]PassesEntry{
		{ProgrammeID: 3, Series: 2},
		{ProgrammeID: 2, Series: 2, CombinableWithID: &main, CombinedSeries: []int{1, 3}},
		{ProgrammeID: 1, Series: 4},
	})
	if want := "P1 P2 + P0 + P1* P2 P3* P4"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
