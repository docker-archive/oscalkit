// oscalkit - OSCAL conversion utility
// Written in 2017 by Andrew Weiss <andrew.weiss@docker.com>

// To the extent possible under law, the author(s) have dedicated all copyright
// and related and neighboring rights to this software to the public domain worldwide.
// This software is distributed without any warranty.

// You should have received a copy of the CC0 Public Domain Dedication along with this software.
// If not, see <http://creativecommons.org/publicdomain/zero/1.0/>.

package oscal

// Section ...
type Section struct {
	ID            string      `xml:"id,attr,omitempty" json:"id,omitempty" yaml:"id,omitempty"`
	OptionalClass string      `xml:"class,attr,omitempty" json:"class,omitempty" yaml:"class,omitempty"`
	Title         *Raw        `xml:"title" json:"title" yaml:"title"`
	Prose         *Prose      `xml:",any" json:"prose,omitempty" yaml:"prose,omitempty"`
	Sections      []Section   `xml:"section,omitempty" json:"sections,omitempty" yaml:"sections,omitempty"`
	References    *References `xml:"references,omitempty" json:"references,omitempty" yaml:"references,omitempty"`
}
