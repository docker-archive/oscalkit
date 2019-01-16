package profile

import (
	"encoding/xml"

	"github.com/docker/oscalkit/types/oscal/catalog"
)

// Each OSCAL profile is defined by a Profile element
type Profile struct {
	XMLName xml.Name `xml:"http://csrc.nist.gov/ns/oscal/1.0 profile" json:"-"`
	ID      string   `xml:"id,attr,omitempty" json:"id,omitempty"`
	Merge   *Merge   `xml:"merge,omitempty" json:"merge,omitempty"`
	Modify  *Modify  `xml:"modify,omitempty" json:"modify,omitempty"`
	Imports []Import `xml:"import,omitempty" json:"imports,omitempty"`
}

// An Import element designates a catalog, profile, or other resource to be
// included (referenced and potentially modified) by this profile.
type Import struct {
	Href    *catalog.Href `xml:"href,attr,omitempty" json:"href,omitempty"`
	Include *Include      `xml:"include,omitempty" json:"include,omitempty"`
	Exclude *Exclude      `xml:"exclude,omitempty" json:"exclude,omitempty"`
}

// A Merge element merges controls in resolution.
type Merge struct {
	Combine *Combine `xml:"combine,omitempty" json:"combine,omitempty"`
	AsIs    AsIs     `xml:"as-is,omitempty" json:"asIs,omitempty"`
	Custom  *Custom  `xml:"custom,omitempty" json:"custom,omitempty"`
}

// A Custom element frames a structure for embedding represented controls in
// resolution.
type Custom struct {
	IdSelectors      []Call  `xml:"call,omitempty" json:"calls,omitempty"`
	PatternSelectors []Match `xml:"match,omitempty" json:"matches,omitempty"`
	Groups           []Group `xml:"group,omitempty" json:"groups,omitempty"`
}

// As in catalogs, a group of (selected) controls or of groups of controls
type Group struct {
	Groups           []Group `xml:"group,omitempty" json:"groups,omitempty"`
	IdSelectors      []Call  `xml:"call,omitempty" json:"calls,omitempty"`
	PatternSelectors []Match `xml:"match,omitempty" json:"matches,omitempty"`
}

// Set parameters or amend controls in resolution
type Modify struct {
	ParamSettings []SetParam `xml:"set-param,omitempty" json:"set-params,omitempty"`
	Alterations   []Alter    `xml:"alter,omitempty" json:"alters,omitempty"`
}

// Specifies which controls and subcontrols to include from the resource (source
// catalog) being imported
type Include struct {
	All              *All    `xml:"all,omitempty" json:"all,omitempty"`
	IdSelectors      []Call  `xml:"call,omitempty" json:"calls,omitempty"`
	PatternSelectors []Match `xml:"match,omitempty" json:"matches,omitempty"`
}

// Which controls and subcontrols to exclude from the resource (source catalog)
// being imported
type Exclude struct {
	IdSelectors      []Call  `xml:"call,omitempty" json:"calls,omitempty"`
	PatternSelectors []Match `xml:"match,omitempty" json:"matches,omitempty"`
}

// A parameter setting, to be propagated to points of insertion
type SetParam struct {
	Id           string               `xml:"param-id,attr,omitempty" json:"id,omitempty"`
	Class        string               `xml:"class,attr,omitempty" json:"class,omitempty"`
	DependsOn    string               `xml:"depends-on,attr,omitempty" json:"dependsOn,omitempty"`
	Label        catalog.Label        `xml:"label,omitempty" json:"label,omitempty"`
	Descriptions []catalog.Desc       `xml:"desc,omitempty" json:"descs,omitempty"`
	Constraints  []catalog.Constraint `xml:"constraint,omitempty" json:"constraints,omitempty"`
	Links        []catalog.Link       `xml:"link,omitempty" json:"links,omitempty"`
	Parts        []catalog.Part       `xml:"part,omitempty" json:"parts,omitempty"`
	Value        catalog.Value        `xml:"value,omitempty" json:"value,omitempty"`
	Select       *catalog.Select      `xml:"select,omitempty" json:"select,omitempty"`
}

// An Alter element specifies changes to be made to an included control or
// subcontrol when a profile is resolved.
type Alter struct {

	// Value of the 'id' flag on a target control
	ControlId string `xml:"control-id,attr,omitempty" json:"controlId,omitempty"`
	// Value of the 'id' flag on a target subcontrol
	SubcontrolId string   `xml:"subcontrol-id,attr,omitempty" json:"subcontrolId,omitempty"`
	Removals     []Remove `xml:"remove,omitempty" json:"removes,omitempty"`
	Additions    []Add    `xml:"add,omitempty" json:"adds,omitempty"`
}

// Specifies contents to be added into controls or subcontrols, in resolution
type Add struct {

	// Where to add the new content with respect to the targeted element (beside it or
	// inside it)
	Position   string              `xml:"position,attr,omitempty" json:"position,omitempty"`
	Title      catalog.Title       `xml:"title,omitempty" json:"title,omitempty"`
	Props      []catalog.Prop      `xml:"prop,omitempty" json:"props,omitempty"`
	Links      []catalog.Link      `xml:"link,omitempty" json:"links,omitempty"`
	References *catalog.References `xml:"references,omitempty" json:"references,omitempty"`
	Params     []catalog.Param     `xml:"param,omitempty" json:"params,omitempty"`
	Parts      []catalog.Part      `xml:"part,omitempty" json:"parts,omitempty"`
}

// A Combine element defines whether and how to combine multiple (competing)// versions of the same control
type Combine struct {
	// How clashing controls and subcontrols should be handled
	Method string `xml:"method,attr,omitempty" json:"method,omitempty"`
	Value  string `xml:",chardata" json:"value,omitempty"`
}

// An As-is element indicates that the controls should be structured in resolution
// as they are structured in their source catalogs. It does not contain any
// elements or attributes.
type AsIs string

// Include all controls from the imported resource (catalog)
type All struct {
	// Whether subcontrols should be implicitly included, if not called.
	WithSubcontrols string `xml:"with-subcontrols,attr,omitempty" json:"withSubcontrols,omitempty"`
	Value           string `xml:",chardata" json:"value,omitempty"`
}

// Call a control or subcontrol by its ID
type Call struct {
	// Value of the 'id' flag on a target control
	ControlId string `xml:"control-id,attr,omitempty" json:"controlId,omitempty"`

	// Value of the 'id' flag on a target subcontrol
	SubcontrolId string `xml:"subcontrol-id,attr,omitempty" json:"subcontrolId,omitempty"`

	// Whether a control should be implicitly included, if not called.
	WithControl string `xml:"with-control,attr,omitempty" json:"withControl,omitempty"`

	// Whether subcontrols should be implicitly included, if not called.
	WithSubcontrols string `xml:"with-subcontrols,attr,omitempty" json:"withSubcontrols,omitempty"`
	Value           string `xml:",chardata" json:"value,omitempty"`
}

// Select controls by (regular expression) match on ID
type Match struct {
	// A regular expression matching the IDs of one or more controls or subcontrols to

	// be selected
	Pattern string `xml:"pattern,attr,omitempty" json:"pattern,omitempty"`

	// A regular expression matching the IDs of one or more controls or subcontrols to

	// be selected
	Order string `xml:"order,attr,omitempty" json:"order,omitempty"`

	// Whether a control should be implicitly included, if not called.
	WithControl string `xml:"with-control,attr,omitempty" json:"withControl,omitempty"`

	// Whether subcontrols should be implicitly included, if not called.
	WithSubcontrols string `xml:"with-subcontrols,attr,omitempty" json:"withSubcontrols,omitempty"`
	Value           string `xml:",chardata" json:"value,omitempty"`
}

// Specifies elements to be removed from a control or subcontrol, in resolution
type Remove struct {
	// Items to remove, by class
	ClassRef string `xml:"class-ref,attr,omitempty" json:"classRef,omitempty"`

	// Items to remove, by ID
	IdRef string `xml:"id-ref,attr,omitempty" json:"idRef,omitempty"`

	// Items to remove, by item type (name), e.g. title or prop
	ItemName string `xml:"item-name,attr,omitempty" json:"itemName,omitempty"`
	Value    string `xml:",chardata" json:"value,omitempty"`
}
