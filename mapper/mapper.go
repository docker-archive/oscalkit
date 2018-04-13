package mapper

import "github.com/opencontrol/oscalkit/types/oscal"

// Mapper ...
type Mapper interface {
	// Map returns an OSCAL-formatted mapping between two OSCAL components
	Map(to oscal.OSCAL) oscal.OSCAL
}
