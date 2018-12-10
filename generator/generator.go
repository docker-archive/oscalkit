package generator

import (
	"errors"
	"os"

	"github.com/opencontrol/oscalkit/templates"
	"github.com/opencontrol/oscalkit/types/oscal/catalog"
)

//GenerateCatalogs GenerateCatalogs
func GenerateCatalogs(f *os.File, c []*catalog.Catalog) error {
	if len(c) == 0 {
		return errors.New("no Catalogs")
	}
	t := templates.GetCatalogTemplate()
	return t.Execute(f, struct {
		Catalogs []*catalog.Catalog
	}{c})

}
