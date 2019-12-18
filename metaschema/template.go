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
)

func GenerateTypes(metaschema *Metaschema) error {
	t, err := template.New("types.tmpl").Funcs(template.FuncMap{
		"toLower":       strings.ToLower,
		"toCamel":       strcase.ToCamel,
		"toLowerCamel":  strcase.ToLowerCamel,
		"plural":        inflection.Plural,
		"packageImport": packageImport,
		"getImports":    getImports,
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
