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
)

func GenerateTypes(metaschema *Metaschema) error {
	t, err := template.New("types.tmpl").Funcs(template.FuncMap{
		"toCamel":      strcase.ToCamel,
		"toLowerCamel": strcase.ToLowerCamel,
		"getImports":   getImports,
	}).ParseFiles("types.tmpl")
	if err != nil {
		return err
	}

	packageName := metaschema.GoPackageName()
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

func getImports(metaschema Metaschema) string {
	var imports strings.Builder
	imports.WriteString("import (\n")
	if metaschema.ContainsRootElement() {
		imports.WriteString("\t\"encoding/xml\"\n")
	}

	for _, im := range metaschema.ImportedMetaschema {
		if im.GoPackageName() != "validation_common_root" {
			imports.WriteString(fmt.Sprintf("\n\t\"github.com/docker/oscalkit/types/oscal/%s\"\n", im.GoPackageName()))
		}
	}

	imports.WriteString(")")

	return imports.String()
}
