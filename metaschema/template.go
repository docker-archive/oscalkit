package metaschema

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"os"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
	wordwrap "github.com/mitchellh/go-wordwrap"
)

const (
	DatatypeString  Datatype = "string"
	DatatypeIDRef   Datatype = "IDREF"
	DatatypeNCName  Datatype = "NCName"
	DatatypeNMToken Datatype = "NMTOKEN"
	DatatypeID      Datatype = "ID"
	DatatypeAnyURI  Datatype = "anyURI"
)

type Datatype string

var datatypes = map[Datatype]string{
	DatatypeString:  "string",
	DatatypeIDRef:   "string",
	DatatypeNCName:  "string",
	DatatypeNMToken: "string",
	DatatypeID:      "string",
	DatatypeAnyURI:  "Href",
}

func GenerateTypes(metaschema *Metaschema) error {
	t, err := template.New("types.tmpl").Funcs(template.FuncMap{
		"toLower":         strings.ToLower,
		"toCamel":         strcase.ToCamel,
		"toLowerCamel":    strcase.ToLowerCamel,
		"plural":          inflection.Plural,
		"wrapString":      wrapString,
		"parseDatatype":   parseDatatype,
		"commentFlag":     commentFlag,
		"packageImport":   packageImport,
		"getImports":      getImports,
		"requiresPointer": requiresPointer,
	}).ParseFiles("types.tmpl")
	if err != nil {
		return err
	}

	packageName := strings.ToLower(metaschema.Root)
	f, err := os.Create(fmt.Sprintf("../types/oscal/%s/%s.go", packageName, packageName))
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "types.tmpl", metaschema); err != nil {
		return err
	}

	p, err := format.Source(buf.Bytes())
	if err != nil {
		return errors.New(err.Error() + " in following file:\n" + string(buf.Bytes()))
	}

	_, err = f.Write(p)
	if err != nil {
		return err
	}

	return nil
}

func wrapString(text string) []string {
	text = strings.Join(strings.Fields(text), " ")

	return strings.Split(wordwrap.WrapString(text, 80), "\n")
}

func parseDatatype(datatype string, packageName string) string {
	if packageName != "catalog" {
		if dt, ok := datatypes[Datatype(datatype)]; ok && dt != "string" {
			return fmt.Sprintf("*catalog.%s", dt)
		}
	}
	return datatypes[Datatype(datatype)]
}

func commentFlag(flagName string, flagDefs []DefineFlag) []string {
	for _, df := range flagDefs {
		if flagName == df.Name {
			return wrapString(df.Description)
		}
	}

	return nil
}

func packageImport(named string, metaschema Metaschema) string {
	for _, df := range metaschema.DefineFlag {
		if df.Name == named {
			return ""
		}
	}

	for _, da := range metaschema.DefineAssembly {
		if da.Name == named {
			return ""
		}
	}

	for _, df := range metaschema.DefineField {
		if df.Name == named {
			return ""
		}
	}

	for _, im := range metaschema.ImportedMetaschema {
		return im.Root + "."
	}

	return ""
}

func getImports(metaschema Metaschema) string {
	var imports strings.Builder
	imports.WriteString("import (\n")
	if metaschema.ContainsRootElement() {
		imports.WriteString("\t\"encoding/xml\"\n")
	}

	for _, im := range metaschema.ImportedMetaschema {
		imports.WriteString(fmt.Sprintf("\n\t\"github.com/docker/oscalkit/types/oscal/%s\"\n", strings.ReplaceAll(strings.ToLower(im.Root), "-", "_")))
	}

	imports.WriteString(")")

	return imports.String()
}

func requiresPointer(fieldName string, metaschema Metaschema) bool {
	df, err := metaschema.GetDefineField(fieldName)
	if err == nil {
		return len(df.Flags) > 0
	}
	return false
}
