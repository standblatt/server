package domain

import (
	"fmt"
	"math"
	"strings"
)

var (
	ones = []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten", "eleven",
		"twelve", "thirteen", "fourteen", "fifteen", "sixteen", "seventeen", "eighteen", "nineteen"}
	tens = []string{"", "", "twenty", "thirty", "forty", "fifty", "sixty", "seventy", "eighty", "ninety"}

	onesDE = []string{"null", "eins", "zwei", "drei", "vier", "fünf", "sechs", "sieben", "acht", "neun", "zehn", "elf",
		"zwölf", "dreizehn", "vierzehn", "fünfzehn", "sechzehn", "siebzehn", "achtzehn", "neunzehn"}
	tensDE = []string{"", "", "zwanzig", "dreissig", "vierzig", "fünfzig", "sechzig", "siebzig", "achtzig", "neunzig"}

	onesFR = []string{"zéro", "un", "deux", "trois", "quatre", "cinq", "six", "sept", "huit", "neuf", "dix", "onze",
		"douze", "treize", "quatorze", "quinze", "seize", "dix-sept", "dix-huit", "dix-neuf"}
	tensFR = []string{"", "", "vingt", "trente", "quarante", "cinquante", "soixante", "septante", "quatre-vingt", "nonante"}

	onesIT = []string{"zero", "uno", "due", "tre", "quattro", "cinque", "sei", "sette", "otto", "nove", "dieci", "undici",
		"dodici", "tredici", "quattordici", "quindici", "sedici", "diciassette", "diciotto", "diciannove"}
	tensIT = []string{"", "", "venti", "trenta", "quaranta", "cinquanta", "sessanta", "settanta", "ottanta", "novanta"}
)

// numberWords writes n below one million in English words.
func numberWords(n int) string {
	switch {
	case n < 20:
		return ones[n]
	case n < 100:
		if n%10 == 0 {
			return tens[n/10]
		}
		return tens[n/10] + "-" + ones[n%10]
	case n < 1000:
		if n%100 == 0 {
			return ones[n/100] + " hundred"
		}
		return ones[n/100] + " hundred and " + numberWords(n%100)
	default:
		if n%1000 == 0 {
			return numberWords(n/1000) + " thousand"
		}
		sep := " "
		if n%1000 < 100 {
			sep = " and "
		}
		return numberWords(n/1000) + " thousand" + sep + numberWords(n%1000)
	}
}

// numberWordsDE writes n below one million in German words, as one word.
func numberWordsDE(n int) string {
	prefix := func(n int) string { return einsToEin(numberWordsDE(n)) }
	rest := func(n int) string {
		if n == 0 {
			return ""
		}
		return numberWordsDE(n)
	}
	switch {
	case n < 20:
		return onesDE[n]
	case n < 100:
		if n%10 == 0 {
			return tensDE[n/10]
		}
		return prefix(n%10) + "und" + tensDE[n/10]
	case n < 1000:
		return prefix(n/100) + "hundert" + rest(n%100)
	default:
		return prefix(n/1000) + "tausend" + rest(n%1000)
	}
}

// einsToEin turns a trailing "eins" into "ein" for use before hundert, tausend or und.
func einsToEin(w string) string {
	if strings.HasSuffix(w, "eins") {
		return strings.TrimSuffix(w, "s")
	}
	return w
}

// numberWordsFR writes n below one million in Swiss French words (septante, nonante).
func numberWordsFR(n int) string {
	rest := func(n int) string {
		if n == 0 {
			return ""
		}
		return " " + numberWordsFR(n)
	}
	switch {
	case n < 20:
		return onesFR[n]
	case n < 100:
		switch {
		case n == 80:
			return "quatre-vingts"
		case n%10 == 0:
			return tensFR[n/10]
		case n%10 == 1 && n != 81:
			return tensFR[n/10] + " et un"
		}
		return tensFR[n/10] + "-" + onesFR[n%10]
	case n < 1000:
		if n < 200 {
			return "cent" + rest(n%100)
		}
		if n%100 == 0 {
			return onesFR[n/100] + " cents"
		}
		return onesFR[n/100] + " cent" + rest(n%100)
	default:
		if n < 2000 {
			return "mille" + rest(n%1000)
		}
		w := numberWordsFR(n / 1000)
		if strings.HasSuffix(w, "cents") || strings.HasSuffix(w, "vingts") {
			w = strings.TrimSuffix(w, "s")
		}
		return w + " mille" + rest(n%1000)
	}
}

// numberWordsIT writes n below one million in Italian words, as one word.
func numberWordsIT(n int) string {
	rest := func(n int) string {
		if n == 0 {
			return ""
		}
		return numberWordsIT(n)
	}
	switch {
	case n < 20:
		return onesIT[n]
	case n < 100:
		t, u := tensIT[n/10], n%10
		switch u {
		case 0:
			return t
		case 1, 8:
			return t[:len(t)-1] + onesIT[u]
		case 3:
			return t + "tré"
		}
		return t + onesIT[u]
	case n < 1000:
		// cento drops its o before otto and ottanta: centotto, centottanta.
		r := rest(n % 100)
		hundred := "cento"
		if strings.HasPrefix(r, "o") {
			hundred = "cent"
		}
		if n < 200 {
			return hundred + r
		}
		return onesIT[n/100] + hundred + r
	default:
		if n < 2000 {
			return "mille" + rest(n%1000)
		}
		return numberWordsIT(n/1000) + "mila" + rest(n%1000)
	}
}

type amountLang struct {
	words                           func(int) string
	franc, francs, cent, cents, and string
	one                             string
}

var amountLangs = map[string]amountLang{
	"en": {numberWords, "franc", "francs", "centime", "centimes", "and", "one"},
	"de": {numberWordsDE, "Franken", "Franken", "Rappen", "Rappen", "und", "ein"},
	"fr": {numberWordsFR, "franc", "francs", "centime", "centimes", "et", "un"},
	"it": {numberWordsIT, "franco", "franchi", "centesimo", "centesimi", "e", "un"},
}

// AmountInWords writes a CHF amount below one million in words in the given locale (en, de, fr, it; others
// get English), e.g. "Twenty-five francs and fifty centimes" or "Fünfundzwanzig Franken und fünfzig Rappen".
func AmountInWords(locale string, chf float64) string {
	l, ok := amountLangs[locale]
	if !ok {
		l = amountLangs["en"]
	}
	cents := int(math.Round(chf * 100))
	if cents < 0 || cents >= 100_000_000 {
		return fmt.Sprintf("CHF %.2f", chf)
	}
	count := func(n int, one, many string) string {
		if n == 1 {
			return l.one + " " + one
		}
		w := l.words(n)
		if locale == "de" {
			w = einsToEin(w)
		}
		return w + " " + many
	}
	s := count(cents/100, l.franc, l.francs)
	if rest := cents % 100; rest > 0 {
		s += " " + l.and + " " + count(rest, l.cent, l.cents)
	}
	r := []rune(s)
	return strings.ToUpper(string(r[0])) + string(r[1:])
}

// PassesEntry describes one programme of a registration for the control strip.
type PassesEntry struct {
	ProgrammeID      int64
	Series           int
	CombinableWithID *int64
	CombinedSeries   []int
}

// PassesText summarises the series a shooter shoots: "P1 P2" per programme, "P0" for a programme that takes
// its series from another one, and "*" on the series of the other programme that are taken.
func PassesText(entries []PassesEntry) string {
	taken := map[int64]map[int]bool{}
	for _, e := range entries {
		if len(e.CombinedSeries) > 0 && e.CombinableWithID != nil {
			if taken[*e.CombinableWithID] == nil {
				taken[*e.CombinableWithID] = map[int]bool{}
			}
			for _, n := range e.CombinedSeries {
				taken[*e.CombinableWithID][n] = true
			}
		}
	}
	parts := make([]string, 0, len(entries))
	for _, e := range entries {
		if len(e.CombinedSeries) > 0 {
			parts = append(parts, "P0")
			continue
		}
		series := make([]string, e.Series)
		for i := range series {
			series[i] = fmt.Sprintf("P%d", i+1)
			if taken[e.ProgrammeID][i+1] {
				series[i] += "*"
			}
		}
		parts = append(parts, strings.Join(series, " "))
	}
	return strings.Join(parts, " + ")
}
