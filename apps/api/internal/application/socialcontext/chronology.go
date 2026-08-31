package socialcontext

import (
	"strings"
	"time"
)

// ActivityChronologyLess defines one deterministic total order for normalized Activity projections.
// Missing timestamps sort before present timestamps; Activity ID is the final stable tie-breaker.
func ActivityChronologyLess(leftID string, leftOccurredAt, leftContributedAt time.Time, rightID string, rightOccurredAt, rightContributedAt time.Time) bool {
	if comparison := compareOptionalTime(leftOccurredAt, rightOccurredAt); comparison != 0 {
		return comparison < 0
	}
	if comparison := compareOptionalTime(leftContributedAt, rightContributedAt); comparison != 0 {
		return comparison < 0
	}
	return strings.TrimSpace(leftID) < strings.TrimSpace(rightID)
}

func compareOptionalTime(left, right time.Time) int {
	switch {
	case left.IsZero() && !right.IsZero():
		return -1
	case !left.IsZero() && right.IsZero():
		return 1
	case left.IsZero() && right.IsZero():
		return 0
	case left.Before(right):
		return -1
	case right.Before(left):
		return 1
	default:
		return 0
	}
}
