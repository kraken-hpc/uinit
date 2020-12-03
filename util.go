package uinit

import (
	"strings"
	"unicode"
)

// misc utility functions

// SplitCommandLine strings on spaces except when a space is within a quoted, bracketed, or braced string.
// Supports nesting multiple brackets or braces.
func SplitCommandLine(s string) []string {
	// func that ranges across the provided map[rune]int returning true if any
	// values are greater than the provided int.
	mapGt := func(runes map[rune]int, g int) bool {
		for _, i := range runes {
			if i > g {
				return true
			}
		}
		return false
	}

	lastRune := map[rune]int{}
	f := func(c rune) bool {
		switch {
		case lastRune[c] > 0:
			lastRune[c]--
			return false
		case unicode.In(c, unicode.Quotation_Mark):
			lastRune[c]++
			return false
		case c == '[':
			lastRune[']']++
			return false
		case c == '{':
			lastRune['}']++
			return false
		case mapGt(lastRune, 0):
			return false
		default:
			return c == ' '
		}
	}
	return strings.FieldsFunc(s, f)
}
