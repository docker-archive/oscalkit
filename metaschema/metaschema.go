package metaschema

import (
	"encoding/xml"
	"fmt"
	"github.com/iancoleman/strcase"
	"net/url"
	"strings"
)

//go:generate go run generate.go

// ...
const (
	AsTypeBoolean         AsType = "boolean"
	AsTypeEmpty           AsType = "empty"
	AsTypeString          AsType = "string"
	AsTypeMixed           AsType = "mixed"
	AsTypeMarkupLine      AsType = "markup-line"
	AsTypeMarkupMultiLine AsType = "markup-multiline"
	AsTypeDate            AsType = "date"
	AsTypeDateTimeTz      AsType = "dateTime-with-timezone"
	AsTypeNcName          AsType = "NCName"
	AsTypeEmail           AsType = "email"
	AsTypeURI             AsType = "uri"
	AsTypeBase64          AsType = "base64Binary"

	ShowDocsXML     ShowDocs = "xml"
	ShowDocsJSON    ShowDocs = "json"
	ShowDocsXMLJSON ShowDocs = "xml json"
)

var ShowDocsOptions = []ShowDocs{
	ShowDocsXML,
	ShowDocsJSON,
	ShowDocsXMLJSON,
}

type GoType interface {
	GoName() string
	GetMetaschema() *Metaschema
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
	Dependencies       map[string]GoType
}

func (metaschema *Metaschema) registerDependency(name string, dependency GoType) {
	if dependency.GetMetaschema() != metaschema {
		if metaschema.Dependencies == nil {
			metaschema.Dependencies = make(map[string]GoType)
		}
		if _, ok := metaschema.Dependencies[name]; !ok {
			metaschema.Dependencies[name] = dependency
		}
	}
}

func (metaschema *Metaschema) linkAssemblies(list []Assembly) error {
	var err error
	for i, a := range list {
		if a.Ref != "" {
			a.Def, err = metaschema.GetDefineAssembly(a.Ref)
			if err != nil {
				return err
			}
			a.Metaschema = metaschema
			metaschema.registerDependency(a.Ref, a.Def)
			list[i] = a
		}
	}
	return nil
}

func (metaschema *Metaschema) linkFields(list []Field) error {
	var err error
	for i, f := range list {
		if f.Ref != "" {
			f.Def, err = metaschema.GetDefineField(f.Ref)
			if err != nil {
				return err
			}
			f.Metaschema = metaschema
			metaschema.registerDependency(f.Ref, f.Def)
			list[i] = f
		}
	}
	return nil
}

func (metaschema *Metaschema) linkFlags(list []Flag) error {
	var err error
	for i, f := range list {
		if f.Ref != "" {
			f.Def, err = metaschema.GetDefineFlag(f.Ref)
			if err != nil {
				return err
			}
			f.Metaschema = metaschema
			list[i] = f
		}
	}
	return nil
}

func (metaschema *Metaschema) LinkDefinitions() error {
	var err error
	for _, da := range metaschema.DefineAssembly {
		if err = metaschema.linkFlags(da.Flags); err != nil {
			return err
		}
		if err = metaschema.linkAssemblies(da.Model.Assembly); err != nil {
			return err
		}
		if err = metaschema.linkFields(da.Model.Field); err != nil {
			return err
		}
		for _, c := range da.Model.Choice {
			if err = metaschema.linkAssemblies(c.Assembly); err != nil {
				return err
			}
			if err = metaschema.linkFields(c.Field); err != nil {
				return err
			}
		}
	}

	for _, df := range metaschema.DefineField {
		if err = metaschema.linkFlags(df.Flags); err != nil {
			return err
		}
	}
	return nil
}

func (metaschema *Metaschema) GetDefineField(name string) (*DefineField, error) {
	for _, v := range metaschema.DefineField {
		if name == v.Name {
			if v.Metaschema == nil {
				v.Metaschema = metaschema
			}
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
			if v.Metaschema == nil {
				v.Metaschema = metaschema
			}
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
			if v.Metaschema == nil {
				v.Metaschema = metaschema
			}
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

func (metaschema *Metaschema) GoPackageName() string {
	return strings.ReplaceAll(strings.ToLower(metaschema.Root), "-", "_")
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
	Metaschema  *Metaschema
}

func (da *DefineAssembly) GoName() string {
	return strcase.ToCamel(da.Name)
}

func (da *DefineAssembly) RepresentsRootElement() bool {
	return da.Name == "catalog" || da.Name == "profile" || da.Name == "declarations"
}

func (a *DefineAssembly) GoComment() string {
	return handleMultiline(a.Description)
}

func (a *DefineAssembly) GetMetaschema() *Metaschema {
	return a.Metaschema
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
	AsType      AsType    `xml:"as-type,attr"`
	Metaschema  *Metaschema
}

func (df *DefineField) GoName() string {
	return strcase.ToCamel(df.Name)
}

func (df *DefineField) RequiresPointer() bool {
	return len(df.Flags) > 0
}

func (f *DefineField) GoComment() string {
	return handleMultiline(f.Description)
}

func (df *DefineField) GetMetaschema() *Metaschema {
	return df.Metaschema
}

func (df *DefineField) IsMarkup() bool {
	return df.AsType == AsTypeMarkupMultiLine
}

type DefineFlag struct {
	Name     string   `xml:"name,attr"`
	AsType   datatype `xml:"as-type,attr"`
	ShowDocs ShowDocs `xml:"show-docs,attr"`

	FormalName  string    `xml:"formal-name"`
	Description string    `xml:"description"`
	Remarks     *Remarks  `xml:"remarks"`
	Examples    []Example `xml:"example"`
	Metaschema  *Metaschema
}

func (df *DefineFlag) GoName() string {
	return strcase.ToCamel(df.Name)
}

func (df *DefineFlag) GetMetaschema() *Metaschema {
	return df.Metaschema
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
	Metaschema  *Metaschema
}

func (a *Assembly) GoComment() string {
	if a.Description != "" {
		return handleMultiline(a.Description)
	}
	return a.Def.GoComment()
}

func (a *Assembly) GoName() string {
	if a.Named != "" {
		return strcase.ToCamel(a.Named)
	}
	return a.Def.GoName()
}

func (a *Assembly) GoMemLayout() string {
	if a.GroupAs != nil {
		return "[]"
	}
	return "*"
}

func (a *Assembly) XmlName() string {
	if a.GroupAs != nil {
		return a.GroupAs.Name
	} else if a.Named != "" {
		return a.Named
	} else {
		return a.Def.Name
	}
}

func (a *Assembly) GoPackageName() string {
	if a.Ref == "" {
		return ""
	} else if a.Def.Metaschema == a.Metaschema {
		return ""
	} else {
		return a.Def.Metaschema.GoPackageName() + "."
	}
}

type Field struct {
	Named    string `xml:"named,attr"`
	Required string `xml:"required,attr"`

	Description string   `xml:"description"`
	Remarks     *Remarks `xml:"remarks"`
	Ref         string   `xml:"ref,attr"`
	GroupAs     *GroupAs `xml:"group-as"`
	Def         *DefineField
	Metaschema  *Metaschema
}

func (f *Field) GoComment() string {
	if f.Description != "" {
		return handleMultiline(f.Description)
	}
	return f.Def.GoComment()
}

func (f *Field) RequiresPointer() bool {
	return f.Def.RequiresPointer()
}

func (f *Field) GoName() string {
	if f.Named != "" {
		return strcase.ToCamel(f.Named)
	}
	return f.Def.GoName()
}

func (f *Field) GoPackageName() string {
	if f.Ref == "" {
		return ""
	} else if f.Def.Metaschema == f.Metaschema {
		return ""
	} else {
		return f.Def.Metaschema.GoPackageName() + "."
	}
}

func (f *Field) GoMemLayout() string {
	if f.GroupAs != nil {
		return "[]"
	} else if f.RequiresPointer() {
		return "*"
	}
	return ""
}

func (f *Field) XmlName() string {
	if f.GroupAs != nil {
		return f.GroupAs.Name
	} else if f.Named != "" {
		return f.Named
	} else {
		return f.Def.Name
	}
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
	Metaschema  *Metaschema
}

func (f *Flag) GoComment() string {
	if f.Description != "" {
		return handleMultiline(f.Description)
	}
	return handleMultiline(f.Def.Description)
}

func (f *Flag) GoDatatype() (string, error) {
	dt := f.AsType
	if dt == "" {
		if f.Ref == "" && (f.Name == "position" || f.Name == "asset-id" || f.Name == "use" || f.Name == "system") {
			// workaround bug: inline definition without type hint https://github.com/usnistgov/OSCAL/pull/570
			return "string", nil
		}
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
	return f.Def.GoName()
}

func (f *Flag) XmlName() string {
	if f.Name != "" {
		return f.Name
	}
	return f.Def.Name
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

type AsType string

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
	datatypeString             datatype = "string"
	datatypeIDRef              datatype = "IDREF"
	datatypeNCName             datatype = "NCName"
	datatypeNMToken            datatype = "NMTOKEN"
	datatypeID                 datatype = "ID"
	datatypeAnyURI             datatype = "anyURI"
	datatypeURIRef             datatype = "uri-reference"
	datatypeURI                datatype = "uri"
	datatypeNonNegativeInteger datatype = "nonNegativeInteger"
)

var goDatatypeMap = map[datatype]string{
	datatypeString:             "string",
	datatypeIDRef:              "string",
	datatypeNCName:             "string",
	datatypeNMToken:            "string",
	datatypeID:                 "string",
	datatypeURIRef:             "string",
	datatypeURI:                "string",
	datatypeNonNegativeInteger: "uint64",
}

func handleMultiline(comment string) string {
	return strings.ReplaceAll(comment, "\n", "\n // ")
}
