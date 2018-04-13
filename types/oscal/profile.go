// oscalkit - OSCAL conversion utility
// Written in 2017 by Andrew Weiss <andrew.weiss@docker.com>

// To the extent possible under law, the author(s) have dedicated all copyright
// and related and neighboring rights to this software to the public domain worldwide.
// This software is distributed without any warranty.

// You should have received a copy of the CC0 Public Domain Dedication along with this software.
// If not, see <http://creativecommons.org/publicdomain/zero/1.0/>.

package oscal

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// ...
const (
	MethodUseFirst Method = "use-first"
	MethodMerge    Method = "merge"
	MethodKeep     Method = "keep"

	OrderByKeep       OrderBy = "keep"
	OrderByAscending  OrderBy = "ascending"
	OrderByDescending OrderBy = "descending"

	PositionBefore   Position = "before"
	PositionAfter    Position = "after"
	PositionStarting Position = "starting"
	PositionEnding   Position = "ending"
)

var (
	methods = []Method{
		MethodUseFirst,
		MethodMerge,
		MethodKeep,
	}

	orderings = []OrderBy{
		OrderByKeep,
		OrderByAscending,
		OrderByDescending,
	}

	positions = []Position{
		PositionBefore,
		PositionAfter,
		PositionStarting,
		PositionEnding,
	}
)

// Profile ...
type Profile struct {
	XMLName xml.Name `xml:"http://csrc.nist.gov/ns/oscal/1.0 profile" json:"-" yaml:"-"`
	ID      string   `xml:"id,attr,omitempty" json:"id,omitempty" yaml:"id,omitempty"`
	Title   *Raw     `xml:"title,omitempty" json:"title,omitempty" yaml:"title,omitempty"`
	Imports []Import `xml:"import" json:"imports" yaml:"imports"`
	Merge   *Merge   `xml:"merge,omitempty" json:"merge,omitempty" yaml:"merge,omitempty"`
	Modify  *Modify  `xml:"modify,omitempty" json:"modify,omitempty" yaml:"modify,omitempty"`
}

// Href ...
type Href string

// MarshalJSON ...
func (h *Href) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("%s.json", strings.TrimSuffix(string(*h), filepath.Ext(string(*h)))))
}

// MarshalYAML ...
func (h *Href) MarshalYAML() (interface{}, error) {
	return fmt.Sprintf("%s.yaml", strings.TrimSuffix(string(*h), filepath.Ext(string(*h)))), nil
}

// MarshalXMLAttr ...
func (h *Href) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{Name: name, Value: fmt.Sprintf("%s.xml", strings.TrimSuffix(string(*h), filepath.Ext(string(*h))))}, nil
}

// Import ...
type Import struct {
	Href    *Href    `xml:"href,attr" json:"href" yaml:"href"`
	Include *Include `xml:"include,omitempty" json:"include,omitempty" yaml:"include,omitempty"`
	Exclude *Exclude `xml:"exclude,omitempty" json:"exclude,omitempty" yaml:"exclude,omitempty"`
}

// Merge ...
type Merge struct {
	Combine *Combine `xml:"combine,omitempty" json:"combine,omitempty" yaml:"combine,omitempty"`
	Build   *Build   `xml:"as-is,omitempty" json:"asIs,omitempty" yaml:"asIs,omitempty"`
	Frame   *Frame   `xml:"custom,omitempty" json:"custom,omitempty" yaml:"custom,omitempty"`
}

// Combine ...
type Combine struct {
	Method Method `xml:"method,attr" json:"method" yaml:"method"`
}

// Method ...
type Method string

// UnmarshalXMLAttr ...
func (m *Method) UnmarshalXMLAttr(attr xml.Attr) error {
	method := Method(attr.Value)
	for _, mt := range methods {
		if method == mt {
			*m = mt
			return nil
		}
	}

	return fmt.Errorf("Method attribute value %s is not a valid method type", attr.Value)
}

// Build ...
type Build struct{}

// Frame ...
type Frame struct {
	Title         *Raw           `xml:"title,omitempty" json:"title,omitempty" yaml:"title,omitempty"`
	ProfileGroups []ProfileGroup `xml:"group" json:"groups" yaml:"groups"`
	Calls         []Call         `xml:"call,omitempty" json:"calls,omitempty" yaml:"calls,omitempty"`
	Matches       []Match        `xml:"match,omitempty" json:"matches,omitempty" yaml:"matches,omitempty"`
}

// ProfileGroup ...
type ProfileGroup struct {
	Title         *Raw           `xml:"title,omitempty" json:"title,omitempty" yaml:"title,omitempty"`
	ProfileGroups []ProfileGroup `xml:"group" json:"groups" yaml:"groups"`
	Calls         []Call         `xml:"call,omitempty" json:"calls,omitempty" yaml:"calls,omitempty"`
	Matches       []Match        `xml:"match,omitempty" json:"matches,omitempty" yaml:"matches,omitempty"`
}

// Modify ...
type Modify struct {
	ParamSettings []ParamSetting `xml:"set-param,omitempty" json:"paramSettings,omitempty" yaml:"paramSettings,omitempty"`
	Alterations   []Alteration   `xml:"alter,omitempty" json:"alterations,omitempty" yaml:"alterations,omitempty"`
}

// Include ...
type Include struct {
	All     *All    `xml:"all,omitempty" json:"all,omitempty" yaml:"all,omitempty"`
	Calls   []Call  `xml:"call,omitempty" json:"calls,omitempty" yaml:"calls,omitempty"`
	Matches []Match `xml:"match,omitempty" json:"matches,omitempty" yaml:"matches,omitempty"`
}

// WithSubcontrols ...
type WithSubcontrols bool

// WithControl ...
type WithControl bool

// UnmarshalXMLAttr ...
func (w *WithSubcontrols) UnmarshalXMLAttr(attr xml.Attr) error {
	*w = attr.Value == "yes"

	return nil
}

// UnmarshalXMLAttr ...
func (w *WithControl) UnmarshalXMLAttr(attr xml.Attr) error {
	*w = attr.Value == "yes"

	return nil
}

// MarshalXMLAttr ...
func (w *WithSubcontrols) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	if *w {
		return xml.Attr{Name: name, Value: "yes"}, nil
	}

	return xml.Attr{}, nil
}

// All ...
type All struct {
	WithSubcontrols WithSubcontrols `xml:"with-subcontrols,attr,omitempty"`
}

type withSubcontrolsJSON struct {
	WithSubcontrols bool `json:"withSubcontrols"`
}

// Match ...
type Match struct {
	WithSubcontrols WithSubcontrols `xml:"with-subcontrols,attr,omitempty"`
	WithControl     WithControl     `xml:"with-control,attr,omitempty"`
	Pattern         string          `xml:"pattern,attr,omitempty"`
	OrderBy         OrderBy         `xml:"order-by,attr,omitempty"`
}

// OrderBy ...
type OrderBy string

// UnmarshalXMLAttr ...
func (o *OrderBy) UnmarshalXMLAttr(attr xml.Attr) error {
	ob := OrderBy(attr.Value)
	for _, ordering := range orderings {
		if ob == ordering {
			*o = ordering
			return nil
		}
	}

	return fmt.Errorf("order-by attribute value %s is not a valid ordering type", attr.Value)
}

// MarshalJSON ...
func (a *All) MarshalJSON() ([]byte, error) {
	if a.WithSubcontrols {
		return json.Marshal(withSubcontrolsJSON{true})
	}

	return json.Marshal(true)
}

// MarshalYAML ...
func (a *All) MarshalYAML() (interface{}, error) {
	if a.WithSubcontrols {
		return withSubcontrolsJSON{true}, nil
	}

	return true, nil
}

// UnmarshalJSON ...
func (a *All) UnmarshalJSON(data []byte) error {
	var withSubcontrolsObj withSubcontrolsJSON

	if err := json.Unmarshal(data, &withSubcontrolsObj); err != nil {
		var all bool
		if err := json.Unmarshal(data, &all); err != nil {
			return err
		}

		// Go lacks support for marshaling of self-closing tags
		*a = All{}
		return nil
	}

	a.WithSubcontrols = WithSubcontrols(withSubcontrolsObj.WithSubcontrols)
	return nil
}

// Exclude ...
type Exclude struct {
	Matches []Match       `xml:"match,omitempty" json:"matches,omitempty" yaml:"matches,omitempty"`
	Calls   []ExcludeCall `xml:"call,omitempty" json:"calls,omitempty" yaml:"calls,omitempty"`
}

// ExcludeCall ...
type ExcludeCall struct {
	ControlID    string `xml:"control-id,attr,omitempty" json:"controlId,omitempty" yaml:"controlId,omitempty"`
	SubcontrolID string `xml:"subcontrol-id,attr,omitempty" json:"subcontrolId,omitempty" yaml:"subcontrolId,omitempty"`
}

// Call ...
type Call struct {
	WithSubcontrols *WithSubcontrols `xml:"with-subcontrols,attr,omitempty" json:"withSubcontrols,omitempty" yaml:"withSubcontrols,omitempty"`
	ControlID       string           `xml:"control-id,attr,omitempty" json:"controlId,omitempty" yaml:"controlId,omitempty"`
	SubcontrolID    string           `xml:"subcontrol-id,attr,omitempty" json:"subcontrolId,omitempty" yaml:"subcontrolId,omitempty"`
}

// ParamSetting ...
type ParamSetting struct {
	ParamID       string `xml:"param-id,attr,omitempty" json:"paramId,omitempty" yaml:"paramId,omitempty"`
	OptionalClass string `xml:"class,attr,omitempty" json:"class,omitempty" yaml:"class,omitempty"`
	Desc          *Raw   `xml:"desc,omitempty" json:"desc,omitempty" yaml:"desc,omitempty"`
	Label         string `xml:"label,omitempty" json:"label,omitempty" yaml:"label,omitempty"`
	ParamValue    string `xml:"value,omitempty" json:"value,omitempty" yaml:"value,omitempty"`
}

// Alteration ...
type Alteration struct {
	ControlID    string    `xml:"control-id,attr,omitempty" json:"controlId,omitempty" yaml:"controlId,omitempty"`
	SubcontrolID string    `xml:"subcontrol-id,attr,omitempty" json:"subcontrolId,omitempty" yaml:"subcontrolId,omitempty"`
	Remove       []Remove  `xml:"remove,omitempty" json:"remove,omitempty" yaml:"remove,omitempty"`
	Augment      []Augment `xml:"add,omitempty" json:"add,omitempty" yaml:"add,omitempty"`
}

// Remove ...
type Remove struct {
	ClassRef string `xml:"class-ref,attr,omitempty" json:"classRef,omitempty" yaml:"classRef,omitempty"`
	IDRef    string `xml:"id-ref,attr,omitempty" json:"idRef,omitempty" yaml:"idRef,omitempty"`
	ItemName string `xml:"item-name,attr,omitempty" json:"itemName,omitempty" yaml:"itemName,omitempty"`
}

// Augment ...
type Augment struct {
	Position *Position `xml:"position,attr" json:"position" yaml:"position"`
	Title    *Raw      `xml:"title,omitempty" json:"title,omitempty" yaml:"title,omitempty"`
	Params   []Param   `xml:"param,omitempty" json:"params,omitempty" yaml:"params,omitempty"`
	Props    []Prop    `xml:"prop,omitempty" json:"props,omitempty" yaml:"props,omitempty"`
	Parts    []Part    `xml:"part,omitempty" json:"parts,omitempty" yaml:"parts,omitempty"`
}

// Position ...
type Position string

// UnmarshalXMLAttr ...
func (p *Position) UnmarshalXMLAttr(attr xml.Attr) error {
	position := Position(attr.Value)
	for _, pt := range positions {
		if position == pt {
			*p = pt
			return nil
		}
	}

	return fmt.Errorf("Position attribute value %s is not a valid position type", attr.Value)
}

// Component ...
func (p *Profile) Component() string {
	return "profile"
}

// RawXML ...
func (p *Profile) RawXML(prettify bool) ([]byte, error) {
	if prettify {
		return xml.MarshalIndent(p, "", "  ")
	}
	return xml.Marshal(p)
}

// RawJSON ...
func (p *Profile) RawJSON(prettify bool) ([]byte, error) {
	if prettify {
		return json.MarshalIndent(p, "", "  ")
	}
	return json.Marshal(p)

}

// RawYAML ...
func (p *Profile) RawYAML() ([]byte, error) {
	return yaml.Marshal(p)
}

// ScaffoldImplementation ...
func (p *Profile) ScaffoldImplementation() *Implementation {
	implementation := &Implementation{
		Title: fmt.Sprintf("%s Implementation Scaffolding", p.Title),
	}

	var components []Item
	for _, impt := range p.Imports {
		for _, call := range impt.Include.Calls {
			parts := strings.Split(call.ControlID, ".")
			components = addComponent(components, call, parts)
		}
	}

	implementation.Components.Items = components

	return implementation
}

func addComponent(items []Item, call Call, suffix []string) []Item {
	var id string
	if call.ControlID != "" {
		id = call.ControlID
	} else if call.SubcontrolID != "" {
		id = call.SubcontrolID
	}

	if len(suffix) == 2 {
		var id string
		if call.ControlID != "" {
			id = call.ControlID
		} else if call.SubcontrolID != "" {
			id = call.SubcontrolID
		}

		items = append(items, Item{
			ID:    id,
			Title: &Raw{Value: fmt.Sprintf("%s Title", id)},
			Parts: []Part{
				Part{
					OptionalClass: "satisfies",
					Prose: &Prose{
						P: []P{
							P{
								OptionalClass: "narrative",
								Raw:           fmt.Sprintf("%s Narrative", id),
							},
						},
					},
					Props: []Prop{
						Prop{
							OptionalClass: "status",
							Value:         "none",
						},
					},
				},
			},
		})
	} else {
		for i, item := range items {
			if strings.HasPrefix(id, item.ID) {
				items[i].Items = addComponent(items[i].Items, call, suffix[1:])
			}
		}
	}

	return items
}
