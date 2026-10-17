package domain

import (
	"fmt"
	"strings"
)

// Problem is a validation message with its format and arguments kept apart, so the API can translate the format.
type Problem struct {
	Format string
	Args   []any
}

// Problemf makes a Problem like fmt.Sprintf would make its text.
func Problemf(format string, a ...any) Problem { return Problem{Format: format, Args: a} }

// String is the English text.
func (p Problem) String() string { return fmt.Sprintf(p.Format, p.Args...) }

// Join writes the English texts separated by sep.
func Join(ps []Problem, sep string) string {
	s := make([]string, len(ps))
	for i, p := range ps {
		s[i] = p.String()
	}
	return strings.Join(s, sep)
}
