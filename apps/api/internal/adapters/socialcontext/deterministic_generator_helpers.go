package socialcontext

import "strings"

func isObjectlessReversal(object string) bool {
	switch strings.TrimSpace(object) {
	case "", "了":
		return true
	default:
		return false
	}
}

func hasAnyPrefix(value string, prefixes ...string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(value, prefix) {
			return true
		}
	}
	return false
}

func hasAnySubstring(value string, needles ...string) bool {
	for _, needle := range needles {
		if strings.Contains(value, needle) {
			return true
		}
	}
	return false
}
