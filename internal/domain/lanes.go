package domain

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

// ParseLanes reads lane lists like "1, 5, 7-10" into sorted, unique lane numbers.
func ParseLanes(s string) ([]int, error) {
	var out []int
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		from, to, isRange := strings.Cut(part, "-")
		a, err := strconv.Atoi(strings.TrimSpace(from))
		if err != nil || a < 1 {
			return nil, fmt.Errorf("Invalid lane %q", part)
		}
		b := a
		if isRange {
			if b, err = strconv.Atoi(strings.TrimSpace(to)); err != nil || b < a {
				return nil, fmt.Errorf("Invalid lane range %q", part)
			}
		}
		for n := a; n <= b; n++ {
			out = append(out, n)
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("No lanes given")
	}
	slices.Sort(out)
	return slices.Compact(out), nil
}

// FormatLanes writes sorted lane numbers compactly: runs of three or more become ranges, "1-3, 5, 7, 8".
func FormatLanes(lanes []int) string {
	var parts []string
	for i := 0; i < len(lanes); {
		j := i
		for j+1 < len(lanes) && lanes[j+1] == lanes[j]+1 {
			j++
		}
		switch j - i {
		case 0:
			parts = append(parts, strconv.Itoa(lanes[i]))
		case 1:
			parts = append(parts, strconv.Itoa(lanes[i]), strconv.Itoa(lanes[j]))
		default:
			parts = append(parts, fmt.Sprintf("%d-%d", lanes[i], lanes[j]))
		}
		i = j + 1
	}
	return strings.Join(parts, ", ")
}
