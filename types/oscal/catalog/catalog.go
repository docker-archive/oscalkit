package catalog

import (
	"encoding/xml"
)

// A collection of controls
type Catalog struct {
	XMLName xml.Name `xml:"http://csrc.nist.gov/ns/oscal/1.0 catalog" json:"-"`
	// Unique identifier
	Id string `xml:"id,attr,omitempty" json:"id,omitempty"`
	// Declares a major/minor version for this metaschema
	ModelVersion string        `xml:"model-version,attr,omitempty" json:"modelVersion,omitempty"`
	Title        Title         `xml:"title,omitempty" json:"title,omitempty"`
	Declarations *Declarations `xml:"declarations,omitempty" json:"declarations,omitempty"`
	References   *References   `xml:"references,omitempty" json:"references,omitempty"`
	Sections     []Section     `xml:"section,omitempty" json:"sections,omitempty"`
	Groups       []Group       `xml:"group,omitempty" json:"groups,omitempty"`
	Controls     []Control     `xml:"control,omitempty" json:"controls,omitempty"`
}

// Allows the inclusion of prose content within a Catalog.
type Section struct {

	// Unique identifier
	Id string `xml:"id,attr,omitempty" json:"id,omitempty"`
	// Identifies the property or object within the control; a semantic hint
	Class      string      `xml:"class,attr,omitempty" json:"class,omitempty"`
	Title      Title       `xml:"title,omitempty" json:"title,omitempty"`
	References *References `xml:"references,omitempty" json:"references,omitempty"`
	Sections   []Section   `xml:"section,omitempty" json:"sections,omitempty"`
	Prose      *Prose      `xml:",any" json:"prose,omitempty"`
}

// A group of controls, or of groups of controls.
type Group struct {

	// Unique identifier
	Id string `xml:"id,attr,omitempty" json:"id,omitempty"`
	// Identifies the property or object within the control; a semantic hint
	Class      string      `xml:"class,attr,omitempty" json:"class,omitempty"`
	Title      Title       `xml:"title,omitempty" json:"title,omitempty"`
	Props      []Prop      `xml:"prop,omitempty" json:"props,omitempty"`
	References *References `xml:"references,omitempty" json:"references,omitempty"`
	Params     []Param     `xml:"param,omitempty" json:"params,omitempty"`
	Parts      []Part      `xml:"part,omitempty" json:"parts,omitempty"`
	Groups     []Group     `xml:"group,omitempty" json:"groups,omitempty"`
	Controls   []Control   `xml:"control,omitempty" json:"controls,omitempty"`
}

// A structured information object representing a security or privacy control. Each
// security or privacy control within the Catalog is defined by a distinct control
// instance.
type Control struct {

	// Unique identifier
	Id string `xml:"id,attr,omitempty" json:"id,omitempty"`
	// Identifies the property or object within the control; a semantic hint
	Class       string       `xml:"class,attr,omitempty" json:"class,omitempty"`
	Title       Title        `xml:"title,omitempty" json:"title,omitempty"`
	Props       []Prop       `xml:"prop,omitempty" json:"props,omitempty"`
	Links       []Link       `xml:"link,omitempty" json:"links,omitempty"`
	References  *References  `xml:"references,omitempty" json:"references,omitempty"`
	Params      []Param      `xml:"param,omitempty" json:"params,omitempty"`
	Parts       []Part       `xml:"part,omitempty" json:"parts,omitempty"`
	Subcontrols []Subcontrol `xml:"subcontrol,omitempty" json:"subcontrols,omitempty"`
}

// A control extension or enhancement
type Subcontrol struct {

	// Unique identifier
	Id string `xml:"id,attr,omitempty" json:"id,omitempty"`
	// Identifies the property or object within the control; a semantic hint
	Class      string      `xml:"class,attr,omitempty" json:"class,omitempty"`
	Title      Title       `xml:"title,omitempty" json:"title,omitempty"`
	Props      []Prop      `xml:"prop,omitempty" json:"props,omitempty"`
	Links      []Link      `xml:"link,omitempty" json:"links,omitempty"`
	References *References `xml:"references,omitempty" json:"references,omitempty"`
	Params     []Param     `xml:"param,omitempty" json:"params,omitempty"`
	Parts      []Part      `xml:"part,omitempty" json:"parts,omitempty"`
}

// Parameters provide a mechanism for the dynamic assignment of value(s) in a
// control.
type Param struct {

	// Unique identifier
	Id string `xml:"id,attr,omitempty" json:"id,omitempty"`
	// Identifies the property or object within the control; a semantic hint
	Class string `xml:"class,attr,omitempty" json:"class,omitempty"`
	// Another parameter invoking this one
	DependsOn    string       `xml:"depends-on,attr,omitempty" json:"dependsOn,omitempty"`
	Label        Label        `xml:"label,omitempty" json:"label,omitempty"`
	Descriptions []Desc       `xml:"desc,omitempty" json:"descs,omitempty"`
	Constraints  []Constraint `xml:"constraint,omitempty" json:"constraints,omitempty"`
	Links        []Link       `xml:"link,omitempty" json:"links,omitempty"`
	Guidance     []Guideline  `xml:"guideline,omitempty" json:"guidelines,omitempty"`
	Value        Value        `xml:"value,omitempty" json:"value,omitempty"`
	Select       *Select      `xml:"select,omitempty" json:"select,omitempty"`
}

// A prose statement that provides a recommendation for the use of a parameter.
type Guideline struct {
	Prose *Prose `xml:",any" json:"prose,omitempty"`
}

// Presenting a choice among alternatives
type Select struct {

	// When selecting, a requirement such as one or more
	HowMany      string   `xml:"how-many,attr,omitempty" json:"howMany,omitempty"`
	Alternatives []Choice `xml:"choice,omitempty" json:"choices,omitempty"`
}

// A partition or component of a control, subcontrol or part
type Part struct {

	// Unique identifier
	Id string `xml:"id,attr,omitempty" json:"id,omitempty"`
	// Identifies the property or object within the control; a semantic hint
	Class string `xml:"class,attr,omitempty" json:"class,omitempty"`
	Title Title  `xml:"title,omitempty" json:"title,omitempty"`
	Props []Prop `xml:"prop,omitempty" json:"props,omitempty"`
	Links []Link `xml:"link,omitempty" json:"links,omitempty"`
	Parts []Part `xml:"part,omitempty" json:"parts,omitempty"`
	Prose *Prose `xml:",any" json:"prose,omitempty"`
}

// A group of reference descriptions
type References struct {

	// Unique identifier
	Id    string `xml:"id,attr,omitempty" json:"id,omitempty"`
	Links []Link `xml:"link,omitempty" json:"links,omitempty"`
	Refs  []Ref  `xml:"ref,omitempty" json:"refs,omitempty"`
}

// A reference, with one or more citations to standards, related documents, or
// other resources
type Ref struct {

	// Unique identifier
	Id        string     `xml:"id,attr,omitempty" json:"id,omitempty"`
	Citations []Citation `xml:"citation,omitempty" json:"citations,omitempty"`
	Prose     *Prose     `xml:",any" json:"prose,omitempty"`
}

// Either a reference to a declarations file, or a set of declarations
type Declarations struct {
	// A link to a document or document fragment (actual, nominal or projected)
	Href  Href   `xml:"href,attr,omitempty" json:"href,omitempty"`
	Value string `xml:",chardata" json:"value,omitempty"`
}

// A title for display and navigation, exclusive of more specific properties
type Title string

// A value with a name, attributed to the containing control, subcontrol, part, or// group.
type Prop struct {
	// Unique identifier
	Id string `xml:"id,attr,omitempty" json:"id,omitempty"`

	// Identifies the property or object within the control; a semantic hint
	Class string `xml:"class,attr,omitempty" json:"class,omitempty"`
	Value string `xml:",chardata" json:"value,omitempty"`
}

// A placeholder for a missing value, in display.
type Label string

// Indicates and explains the purpose and use of a parameter
type Desc struct {
	// Unique identifier
	Id    string `xml:"id,attr,omitempty" json:"id,omitempty"`
	Value string `xml:",chardata" json:"value,omitempty"`
}

// A formal or informal expression of a constraint or test
type Constraint struct {
	// A formal (executable) expression of a constraint
	Test  string `xml:"test,attr,omitempty" json:"test,omitempty"`
	Value string `xml:",chardata" json:"value,omitempty"`
}

// Indicates a permissible value for a parameter or property
type Value string

// A value selection among several such options
type Choice string

// A line or paragraph with a hypertext link
type Link struct {
	// A link to a document or document fragment (actual, nominal or projected)
	Href Href `xml:"href,attr,omitempty" json:"href,omitempty"`

	// Purpose of the link
	Rel   string `xml:"rel,attr,omitempty" json:"rel,omitempty"`
	Value string `xml:",chardata" json:"value,omitempty"`
}

// Citation of a resource
type Citation struct {
	// Unique identifier
	Id string `xml:"id,attr,omitempty" json:"id,omitempty"`

	// A link to a document or document fragment (actual, nominal or projected)
	Href  Href   `xml:"href,attr,omitempty" json:"href,omitempty"`
	Value string `xml:",chardata" json:"value,omitempty"`
}
