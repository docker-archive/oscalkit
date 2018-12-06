package generator

import (
	"errors"
	"os"

	"github.com/opencontrol/oscalkit/templates"
	"github.com/opencontrol/oscalkit/types/oscal/catalog"
	"github.com/opencontrol/oscalkit/types/oscal/profile"
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

//GenerateProfile Generates profile.go on executer level directory with exported populated struct
func GenerateProfile(f *os.File, p *profile.Profile) error {
	if p == nil {
		return errors.New("Nil Profile")
	}
	t := templates.GetProfileTemplate()
	return t.Execute(f, p)

}
