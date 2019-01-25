// oscalkit - OSCAL conversion utility
// Written in 2017 by Andrew Weiss <andrew.weiss@docker.com>

// To the extent possible under law, the author(s) have dedicated all copyright
// and related and neighboring rights to this software to the public domain worldwide.
// This software is distributed without any warranty.

// You should have received a copy of the CC0 Public Domain Dedication along with this software.
// If not, see <http://creativecommons.org/publicdomain/zero/1.0/>.

package catalog

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"regexp"
	"strings"
)

// Prose ...
type Prose struct {
	XMLName xml.Name
	raw     []string
	order   []string
	P       []P
	UL      []UL
	OL      []OL
	Pre     []Pre
}

// ReplaceInsertParams replaces insert parameters
func (p *Prose) ReplaceInsertParams(parameterID, parameterValue string) error {

	rs := fmt.Sprintf(`<insert param-id="%s">`, parameterID)
	regex, err := regexp.Compile(rs)
	if err != nil {
		return err
	}
	for i := range p.P {
		p.P[i].Raw = regex.ReplaceAllString(p.P[i].Raw, parameterValue)

	}
	for i := range p.OL {
		p.OL[i].Raw = regex.ReplaceAllString(p.OL[i].Raw, parameterValue)
	}
	for i := range p.Pre {
		p.Pre[i].Raw = regex.ReplaceAllString(p.Pre[i].Raw, parameterValue)
	}
	for i := range p.UL {
		p.UL[i].Raw = regex.ReplaceAllString(p.UL[i].Raw, parameterValue)
	}
	return nil
}

// MarshalXML ...
func (p *Prose) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	raw := strings.Join(p.raw, "")

	if raw != "" {
		if err := xml.Unmarshal([]byte(raw), &p); err != nil {
			return err
		}
	}

	p.XMLName = xml.Name{Local: "ul"}
	if err := e.Encode(p.UL); err != nil {
		return err
	}
	p.XMLName = xml.Name{Local: "ol"}
	if err := e.Encode(p.OL); err != nil {
		return err
	}
	p.XMLName = xml.Name{Local: "p"}
	if err := e.Encode(p.P); err != nil {
		return err
	}
	p.XMLName = xml.Name{Local: "pre"}
	if err := e.Encode(p.Pre); err != nil {
		return err
	}

	return nil
}

// UnmarshalXML ...
func (p *Prose) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	switch start.Name.Local {
	case "ul":
		var ul UL
		if err := d.DecodeElement(&ul, &start); err != nil {
			return err
		}

		p.UL = append(p.UL, ul)
		p.order = append(p.order, "ul")

	case "ol":
		var ol OL
		if err := d.DecodeElement(&ol, &start); err != nil {
			return err
		}

		p.OL = append(p.OL, ol)
		p.order = append(p.order, "ol")

	case "p":
		var para P
		if err := d.DecodeElement(&para, &start); err != nil {
			return err
		}

		p.P = append(p.P, para)
		p.order = append(p.order, "p")

	case "pre":
		var pre Pre
		if err := d.DecodeElement(&pre, &start); err != nil {
			return err
		}

		p.Pre = append(p.Pre, pre)
		p.order = append(p.order, "pre")
	}

	return nil
}

// MarshalJSON ...
func (p *Prose) MarshalJSON() ([]byte, error) {
	// If prose originates from OpenControl
	if p.order == nil {
		for _, para := range p.P {
			if para.Raw != "" {
				p.raw = append(p.raw, para.Raw)
			}
		}
		for _, ul := range p.UL {
			if ul.Raw != "" {
				p.raw = append(p.raw, ul.Raw)
			}
		}
		for _, ol := range p.OL {
			if ol.Raw != "" {
				p.raw = append(p.raw, ol.Raw)
			}
		}
		for _, pre := range p.Pre {
			if pre.Raw != "" {
				p.raw = append(p.raw, pre.Raw)
			}
		}

		return json.Marshal(p.raw)
	}

	// If prose originates from XML
	var ulIndex int
	var olIndex int
	var pIndex int
	var preIndex int

	for _, element := range p.order {
		switch element {
		case "ul":
			if ulIndex < len(p.UL) {
				raw, err := xml.Marshal(p.UL[ulIndex])
				if err != nil {
					return nil, err
				}

				p.raw = append(p.raw, formatRawProse(string(raw)))

				ulIndex++
			}

		case "ol":
			if olIndex < len(p.OL) {
				raw, err := xml.Marshal(p.OL[olIndex])
				if err != nil {
					return nil, err
				}

				p.raw = append(p.raw, formatRawProse(string(raw)))

				olIndex++
			}

		case "p":
			if pIndex < len(p.P) {
				raw, err := xml.Marshal(p.P[pIndex])
				if err != nil {
					return nil, err
				}

				p.raw = append(p.raw, formatRawProse(string(raw)))

				pIndex++
			}

		case "pre":
			if preIndex < len(p.Pre) {
				raw, err := xml.Marshal(p.Pre[preIndex])
				if err != nil {
					return nil, err
				}

				p.raw = append(p.raw, formatRawProse(string(raw)))

				preIndex++
			}
		}
	}

	return json.Marshal(p.raw)
}

// MarshalYAML ...
func (p *Prose) MarshalYAML() (interface{}, error) {
	return p.raw, nil
}

// UnmarshalJSON ...
func (p *Prose) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &p.raw)
}

// Raw ...
type Raw struct {
	Value string `xml:",innerxml"`
}

// MarshalJSON ...
func (r *Raw) MarshalJSON() ([]byte, error) {
	return json.Marshal(formatRawProse(r.Value))
}

// MarshalYAML ...
func (r *Raw) MarshalYAML() (interface{}, error) {
	return r.Value, nil
}

// UnmarshalJSON ...
func (r *Raw) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	r.Value = raw

	return nil
}

// P ...
type P struct {
	XMLName     xml.Name `xml:"p" json:"-" yaml:"-"`
	Q           []Q      `xml:"q" json:"-" yaml:"-"`
	Code        []Code   `xml:"code" json:"-" yaml:"-"`
	EM          []EM     `xml:"em" json:"-" yaml:"-"`
	Strong      []Strong `xml:"strong" json:"-" yaml:"-"`
	B           []B      `xml:"b" json:"-" yaml:"-"`
	I           []I      `xml:"i" json:"-" yaml:"-"`
	Sub         []Sub    `xml:"sub" json:"-" yaml:"-"`
	Sup         []Sup    `xml:"sup" json:"-" yaml:"-"`
	Xref        []Xref   `xml:"a" json:"-" yaml:"-"`
	Assignments []Assign `xml:"-" json:"-" yaml:"-"`
	Selection   []Select `xml:"-" json:"-" yaml:"-"`
	Raw         string   `xml:",innerxml" json:"raw,omitempty" yaml:"raw,omitempty"`
}

// Pre ...
type Pre struct {
	XMLName xml.Name `xml:"pre"`
	ID      string   `xml:"id,attr,omitempty"`
	Raw     string   `xml:",innerxml"`
}

// OL ...
type OL struct {
	XMLName xml.Name `xml:"ol"`
	Raw     string   `xml:",innerxml"`
}

// UL ...
type UL struct {
	XMLName xml.Name `xml:"ul"`
	Raw     string   `xml:",innerxml"`
}

// Q ...
type Q struct {
	Code   []Code   `xml:"code" json:"-" yaml:"-"`
	EM     []EM     `xml:"em" json:"-" yaml:"-"`
	I      []I      `xml:"i" json:"-" yaml:"-"`
	Strong []Strong `xml:"strong" json:"-" yaml:"-"`
	Sub    []Sub    `xml:"sub" json:"-" yaml:"-"`
	Sup    []Sup    `xml:"sup" json:"-" yaml:"-"`
	Value  string   `xml:",chardata" json:"-" yaml:"-"`
}

// Code ...
type Code struct {
	Q      []Q      `xml:"q" json:"-" yaml:"-"`
	Code   []Code   `xml:"code" json:"-" yaml:"-"`
	EM     []EM     `xml:"em" json:"-" yaml:"-"`
	Strong []Strong `xml:"strong" json:"-" yaml:"-"`
	B      []B      `xml:"b" json:"-" yaml:"-"`
	I      []I      `xml:"i" json:"-" yaml:"-"`
	Sub    []Sub    `xml:"sub" json:"-" yaml:"-"`
	Sup    []Sup    `xml:"sup" json:"-" yaml:"-"`
	Value  string   `xml:",chardata" json:"-" yaml:"-"`
}

// EM ...
type EM struct {
	Q      []Q      `xml:"q" json:"-" yaml:"-"`
	Code   []Code   `xml:"code" json:"-" yaml:"-"`
	EM     []EM     `xml:"em" json:"-" yaml:"-"`
	Strong []Strong `xml:"strong" json:"-" yaml:"-"`
	B      []B      `xml:"b" json:"-" yaml:"-"`
	I      []I      `xml:"i" json:"-" yaml:"-"`
	Sub    []Sub    `xml:"sub" json:"-" yaml:"-"`
	Sup    []Sup    `xml:"sup" json:"-" yaml:"-"`
	Value  string   `xml:",chardata" json:"-" yaml:"-"`
	Xref   []Xref   `xml:"a" json:"-" yaml:"-"`
}

// Strong ...
type Strong struct {
	Q      []Q      `xml:"q" json:"-" yaml:"-"`
	Code   []Code   `xml:"code" json:"-" yaml:"-"`
	EM     []EM     `xml:"em" json:"-" yaml:"-"`
	Strong []Strong `xml:"strong" json:"-" yaml:"-"`
	B      []B      `xml:"b" json:"-" yaml:"-"`
	I      []I      `xml:"i" json:"-" yaml:"-"`
	Sub    []Sub    `xml:"sub" json:"-" yaml:"-"`
	Sup    []Sup    `xml:"sup" json:"-" yaml:"-"`
	Value  string   `xml:",chardata" json:"-" yaml:"-"`
	Xref   []Xref   `xml:"a" json:"-" yaml:"-"`
}

// I ...
type I struct {
	Q      []Q      `xml:"q" json:"-" yaml:"-"`
	Code   []Code   `xml:"code" json:"-" yaml:"-"`
	EM     []EM     `xml:"em" json:"-" yaml:"-"`
	Strong []Strong `xml:"strong" json:"-" yaml:"-"`
	B      []B      `xml:"b" json:"-" yaml:"-"`
	I      []I      `xml:"i" json:"-" yaml:"-"`
	Sub    []Sub    `xml:"sub" json:"-" yaml:"-"`
	Sup    []Sup    `xml:"sup" json:"-" yaml:"-"`
	Value  string   `xml:",chardata" json:"-" yaml:"-"`
	Xref   []Xref   `xml:"a" json:"-" yaml:"-"`
}

// B ...
type B struct {
	Q      []Q      `xml:"q" json:"-" yaml:"-"`
	Code   []Code   `xml:"code" json:"-" yaml:"-"`
	EM     []EM     `xml:"em" json:"-" yaml:"-"`
	Strong []Strong `xml:"strong" json:"-" yaml:"-"`
	B      []B      `xml:"b" json:"-" yaml:"-"`
	I      []I      `xml:"i" json:"-" yaml:"-"`
	Sub    []Sub    `xml:"sub" json:"-" yaml:"-"`
	Sup    []Sup    `xml:"sup" json:"-" yaml:"-"`
	Value  string   `xml:",chardata" json:"-" yaml:"-"`
	Xref   []Xref   `xml:"a" json:"-" yaml:"-"`
}

// Sub ...
type Sub struct {
	OptionalClass string `xml:"class,attr"`
	Value         string `xml:",chardata"`
}

// Sup ...
type Sup struct {
	OptionalClass string `xml:"class,attr"`
	Value         string `xml:",chardata"`
}

// Span ...
type Span struct {
	OptionalClass string   `xml:"class,attr"`
	Q             []Q      `xml:"q" json:"-" yaml:"-"`
	Code          []Code   `xml:"code" json:"-" yaml:"-"`
	EM            []EM     `xml:"em" json:"-" yaml:"-"`
	Strong        []Strong `xml:"strong" json:"-" yaml:"-"`
	B             []B      `xml:"b" json:"-" yaml:"-"`
	I             []I      `xml:"i" json:"-" yaml:"-"`
	Sub           []Sub    `xml:"sub" json:"-" yaml:"-"`
	Sup           []Sup    `xml:"sup" json:"-" yaml:"-"`
	Value         string   `xml:",chardata"`
	Xref          []Xref   `xml:"a" json:"xrefs" yaml:"xrefs"`
}

// Xref ...
type Xref struct {
	Href *Href  `xml:"href,attr" json:"href" yaml:"href"`
	Q    []Q    `xml:"q" json:"-" yaml:"-"`
	Code []Code `xml:"code" json:"-" yaml:"-"`
	EM   []struct {
		OptionalClass string `xml:"class,attr"`
		Value         string `xml:",chardata"`
	} `xml:"em" json:"-" yaml:"-"`
	Value string `xml:",chardata" json:"value" yaml:"value"`
}

// Assign ...
type Assign struct {
	ID      string `xml:"id,attr"`
	ParamID string `xml:"param-id"`
}

func formatRawProse(raw string) string {
	lines := strings.Split(raw, "\n")

	value := []string{}

	for _, line := range lines {
		value = append(value, strings.TrimSpace(line))
	}

	return strings.Join(value, " ")
}

func traverseParts(part *Part, parameterID, parameterVal string) {
	if part == nil {
		return
	}
	if part.Prose == nil {
		return
	}

	part.Prose.ReplaceInsertParams(parameterID, parameterVal)
	if len(part.Parts) == 0 {
		return
	}
	for i := range part.Parts {
		traverseParts(&part.Parts[i], parameterID, parameterVal)
	}
	return
}

// ModifyProse modifies prose insert parameter template
func (part *Part) ModifyProse(parameterID, parameterVal string) {
	traverseParts(part, parameterID, parameterVal)
}
