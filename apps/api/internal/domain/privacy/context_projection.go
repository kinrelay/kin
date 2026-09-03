package privacy

import "strings"

// ContextProjection is the minimal friend-visible privacy contract for one Social Context.
// It intentionally excludes raw Activity data, provider payloads, and provenance details.
type ContextProjection struct {
	Meaning string
}

// ProjectContext applies the owner's relationship-specific disclosure decision.
// Absence of a decision and hidden decisions are both default-deny.
func ProjectContext(meaning string, decision *DisclosureDecision) (ContextProjection, bool) {
	if !AllowsDisclosure(decision) {
		return ContextProjection{}, false
	}

	return ContextProjection{Meaning: strings.Join(strings.Fields(meaning), " ")}, true
}
