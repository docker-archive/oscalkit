// oscalkit - OSCAL conversion utility
// Written in 2017 by Andrew Weiss <andrew.weiss@docker.com>

// To the extent possible under law, the author(s) have dedicated all copyright
// and related and neighboring rights to this software to the public domain worldwide.
// This software is distributed without any warranty.

// You should have received a copy of the CC0 Public Domain Dedication along with this software.
// If not, see <http://creativecommons.org/publicdomain/zero/1.0/>.

package oscal

// Component ...
type Component struct {
	ID            string      `xml:"id,attr,omitempty" json:"id,omitempty" yaml:"id,omitempty"`
	OptionalClass string      `xml:"class,attr" json:"class,omitempty" yaml:"class,omitempty"`
	Title         *Raw        `xml:"title" json:"title" yaml:"title"`
	Params        []Param     `xml:"param,omitempty" json:"params,omitempty" yaml:"params,omitempty"`
	Props         []Prop      `xml:"prop,omitempty" json:"props,omitempty" yaml:"props,omitempty"`
	Links         []Link      `xml:"link,omitempty" json:"links,omitempty" yaml:"links,omitempty"`
	Prose         *Prose      `xml:",any" json:"prose,omitempty" yaml:"prose,omitempty"`
	Parts         []Part      `xml:"part,omitempty" json:"parts,omitempty" yaml:"parts,omitempty"`
	Components    []Component `xml:"component" json:"components" yaml:"components"`
}
