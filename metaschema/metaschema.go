package metaschema

import (
	"encoding/xml"
	"fmt"
	"github.com/iancoleman/strcase"
	"net/url"
)

//go:generate go run generate.go

// ...
const (
	AsBoolean As = "boolean"
	AsEmpty   As = "empty"
	AsString  As = "string"
	AsMixed   As = "mixed"

	ShowDocsXML     ShowDocs = "xml"
	ShowDocsJSON    ShowDocs = "json"
	ShowDocsXMLJSON ShowDocs = "xml json"
)

var FieldConstraints = []As{
	AsBoolean,
	AsEmpty,
	AsString,
	AsMixed,
}

var ShowDocsOptions = []ShowDocs{
	ShowDocsXML,
	ShowDocsJSON,
	ShowDocsXMLJSON,
}

// Metaschema is the root metaschema element
type Metaschema struct {
	XMLName xml.Name `xml:"http://csrc.nist.gov/ns/oscal/metaschema/1.0 METASCHEMA"`
	Top     string   `xml:"top,attr"`
	Root    string   `xml:"root,attr"`

	// SchemaName describes the scope of application of the data format. For
	// example "OSCAL Catalog"
	SchemaName *SchemaName `xml:"schema-name"`

	// ShortName is a coded version of the schema name for use when a string-safe
	// identifier is needed. For example "oscal-catalog"
	ShortName *ShortName `xml:"short-name"`

	// Remarks are paragraphs describing the metaschema
	Remarks *Remarks `xml:"remarks,omitempty"`

	// Import is a URI to an external metaschema
	Import []Import `xml:"import"`

	// DefineAssembly is one or more assembly definitions
	DefineAssembly []DefineAssembly `xml:"define-assembly"`

	// DefineField is one or more field definitions
	DefineField []DefineField `xml:"define-field"`

	// DefineFlag is one or more flag definitions
	DefineFlag []DefineFlag `xml:"define-flag"`

	ImportedMetaschema []Metaschema
}

func (metaschema *Metaschema) LinkDefinitions() error {
	var err error
	for _, da := range metaschema.DefineAssembly {
		for i, f := range da.Flags {
			if f.Ref != "" {
				f.Def, err = metaschema.GetDefineFlag(f.Ref)
				if err != nil {
					return err
				}
				da.Flags[i] = f
			}

		}
		for i, a := range da.Model.Assembly {
			if a.Ref != "" {
				a.Def, err = metaschema.GetDefineAssembly(a.Ref)
				if err != nil {
					return err
				}
				da.Model.Assembly[i] = a
			}
		}
		for i, f := range da.Model.Field {
			if f.Ref != "" {
				f.Def, err = metaschema.GetDefineField(f.Ref)
				if err != nil {
					return err
				}
				da.Model.Field[i] = f

			}
		}
	}

	for _, df := range metaschema.DefineField {
		for i, f := range df.Flags {
			if f.Ref != "" {
				f.Def, err = metaschema.GetDefineFlag(f.Ref)
				if err != nil {
					return err
				}
				df.Flags[i] = f
			}
		}
	}
	return nil
}

func (metaschema *Metaschema) GetDefineField(name string) (*DefineField, error) {
	for _, v := range metaschema.DefineField {
		if name == v.Name {
			return &v, nil
		}
	}
	for _, m := range metaschema.ImportedMetaschema {
		f, err := m.GetDefineField(name)
		if err == nil {
			return f, err
		}
	}
	return nil, fmt.Errorf("Could not find define-field element with name='%s'.", name)
}

func (metaschema *Metaschema) GetDefineAssembly(name string) (*DefineAssembly, error) {
	for _, v := range metaschema.DefineAssembly {
		if name == v.Name {
			return &v, nil
		}
	}
	for _, m := range metaschema.ImportedMetaschema {
		a, err := m.GetDefineAssembly(name)
		if err == nil {
			return a, err
		}
	}
	return nil, fmt.Errorf("Could not find define-assembly element with name='%s'.", name)
}

func (metaschema *Metaschema) GetDefineFlag(name string) (*DefineFlag, error) {
	for _, v := range metaschema.DefineFlag {
		if name == v.Name {
			return &v, nil
		}
	}
	for _, m := range metaschema.ImportedMetaschema {
		f, err := m.GetDefineFlag(name)
		if err == nil {
			return f, err
		}
	}
	return nil, fmt.Errorf("Could not find define-flag element with name='%s'.", name)
}

func (Metaschema *Metaschema) ContainsRootElement() bool {
	for _, v := range Metaschema.DefineAssembly {
		if v.RepresentsRootElement() {
			return true
		}
	}
	return false
}

// DefineAssembly is a definition for for an object or element that contains
// structured content
type DefineAssembly struct {
	Name     string `xml:"name,attr"`
	GroupAs  string `xml:"group-as,attr"`
	ShowDocs string `xml:"show-docs,attr"`
	Address  string `xml:"address,attr"`

	Flags       []Flag    `xml:"flag"`
	FormalName  string    `xml:"formal-name"`
	Description string    `xml:"description"`
	Remarks     *Remarks  `xml:"remarks"`
	Model       *Model    `xml:"model"`
	Examples    []Example `xml:"example"`
}

func (da *DefineAssembly) RepresentsRootElement() bool {
	return da.Name == "catalog" || da.Name == "profile" || da.Name == "declarations"
}

type DefineField struct {
	Name     string `xml:"name,attr"`
	GroupAs  string `xml:"group-as,attr"`
	ShowDocs string `xml:"show-docs,attr"`

	Flags       []Flag    `xml:"flag"`
	FormalName  string    `xml:"formal-name"`
	Description string    `xml:"description"`
	Remarks     *Remarks  `xml:"remarks"`
	Examples    []Example `xml:"example"`
	As          As        `xml:"as"`
}

func (df *DefineField) RequiresPointer() bool {
	return len(df.Flags) > 0
}

type DefineFlag struct {
	Name     string   `xml:"name,attr"`
	AsType   datatype `xml:"as-type,attr"`
	ShowDocs ShowDocs `xml:"show-docs,attr"`

	FormalName  string    `xml:"formal-name"`
	Description string    `xml:"description"`
	Remarks     *Remarks  `xml:"remarks"`
	Examples    []Example `xml:"example"`
}

type Model struct {
	Assembly []Assembly `xml:"assembly"`
	Field    []Field    `xml:"field"`
	Choice   []Choice   `xml:"choice"`
	Prose    *struct{}  `xml:"prose"`
	Any      *struct{}  `xml:"any"`
}

type Assembly struct {
	Named string `xml:"named,attr"`

	Description string   `xml:"description"`
	Remarks     *Remarks `xml:"remarks"`
	Ref         string   `xml:"ref,attr"`
	GroupAs     *GroupAs `xml:"group-as"`
	Def         *DefineAssembly
}

func (a *Assembly) GoComment() string {
	if a.Description != "" {
		return a.Description
	}
	return a.Def.Description
}

type Field struct {
	Named    string `xml:"named,attr"`
	Required string `xml:"required,attr"`

	Description string   `xml:"description"`
	Remarks     *Remarks `xml:"remarks"`
	Ref         string   `xml:"ref,attr"`
	GroupAs     *GroupAs `xml:"group-as"`
	Def         *DefineField
}

func (f *Field) GoComment() string {
	if f.Description != "" {
		return f.Description
	}
	return f.Def.Description
}

type Flag struct {
	Name     string   `xml:"name,attr"`
	AsType   datatype `xml:"as-type,attr"`
	Required string   `xml:"required,attr"`

	Description string   `xml:"description"`
	Remarks     *Remarks `xml:"remarks"`
	Values      []Value  `xml:"value"`
	Ref         string   `xml:"ref,attr"`
	Def         *DefineFlag
}

func (f *Flag) GoComment() string {
	if f.Description != "" {
		return f.Description
	}
	return f.Def.Description
}

func (f *Flag) GoDatatype() (string, error) {
	dt := f.AsType
	if dt == "" {
		dt = f.Def.AsType
	}

	if goDatatypeMap[dt] == "" {
		return "", fmt.Errorf("Unknown as-type='%s' found.", dt)
	}
	return goDatatypeMap[dt], nil
}

func (f *Flag) GoName() string {
	if f.Name != "" {
		return strcase.ToCamel(f.Name)
	}
	return strcase.ToCamel(f.Def.Name)
}

type Choice struct {
	Field    []Field    `xml:"field"`
	Assembly []Assembly `xml:"assembly"`
}

type GroupAs struct {
	Name string `xml:"name,attr"`
}

type Import struct {
	Href *Href `xml:"href,attr"`
}

// Remarks are descriptions for a particular metaschema, assembly, field, flag
// or example
type Remarks struct {
	P []P `xml:"p"`

	InnerXML string `xml:",innerxml"`
}

type Value struct {
	InnerXML string `xml:",innerxml"`
}

type Title struct {
	Code []string `xml:"code"`
	Q    []string `xml:"q"`

	InnerXML string `xml:",innerxml"`
}

type ShortName struct {
	Code []string `xml:"code"`
	Q    []string `xml:"q"`

	InnerXML string `xml:",innerxml"`
}

type SchemaName struct {
	Code []string `xml:"code"`
	Q    []string `xml:"q"`

	InnerXML string `xml:",innerxml"`
}

type Example struct {
	Href *Href  `xml:"href,attr"`
	Path string `xml:"path,attr"`

	Description string   `xml:"description"`
	Remarks     *Remarks `xml:"remarks"`

	InnerXML string `xml:",innerxml"`
}

type P struct {
	A      []A      `xml:"a"`
	Code   []string `xml:"code"`
	Q      []string `xml:"q"`
	EM     []string `xml:"em"`
	Strong []string `xml:"strong"`

	CharData string `xml:",chardata"`
}

type A struct {
	XMLName xml.Name `xml:"a"`
	Href    *Href    `xml:"href,attr"`

	CharData      string `xml:",chardata"`
	ProcessedLink string `xml:"-"`
}

func (a *A) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type anchor A

	if err := d.DecodeElement((*anchor)(a), &start); err != nil {
		return err
	}

	if a.Href != nil {
		a.ProcessedLink = fmt.Sprintf("%s (%s)", a.CharData, a.Href.URL.String())
	}

	return nil
}

type Href struct {
	URL *url.URL
}

func (h *Href) UnmarshalXMLAttr(attr xml.Attr) error {
	URL, err := url.Parse(attr.Value)
	if err != nil {
		return err
	}

	h.URL = URL

	return nil
}

func (h *Href) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	if h.URL != nil {
		rawURI := h.URL.String()

		return xml.Attr{Name: name, Value: rawURI}, nil
	}

	return xml.Attr{Name: name}, nil
}

type As string

func (a As) UnmarshalXMLAttr(attr xml.Attr) error {
	as := As(attr.Value)

	for _, fieldConstraint := range FieldConstraints {
		if as == fieldConstraint {
			a = as
			return nil
		}
	}

	return fmt.Errorf("Field constraint \"%s\" is not a valid constraint", attr.Value)
}

type ShowDocs string

func (sd ShowDocs) UnmarshalXMLAttr(attr xml.Attr) error {
	showDocs := ShowDocs(attr.Value)

	for _, showDocsOption := range ShowDocsOptions {
		if showDocs == showDocsOption {
			sd = showDocs
			return nil
		}
	}

	return fmt.Errorf("Show docs option \"%s\" is not a valid option", attr.Value)
}

type datatype string

const (
	datatypeString  datatype = "string"
	datatypeIDRef   datatype = "IDREF"
	datatypeNCName  datatype = "NCName"
	datatypeNMToken datatype = "NMTOKEN"
	datatypeID      datatype = "ID"
	datatypeAnyURI  datatype = "anyURI"
	datatypeURIRef  datatype = "uri-reference"
)

var goDatatypeMap = map[datatype]string{
	datatypeString:  "string",
	datatypeIDRef:   "string",
	datatypeNCName:  "string",
	datatypeNMToken: "string",
	datatypeID:      "string",
	datatypeURIRef:  "string",
}
