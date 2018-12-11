package generator

import (
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"

	"github.com/Sirupsen/logrus"

	"github.com/opencontrol/oscalkit/types/oscal/catalog"
	"github.com/opencontrol/oscalkit/types/oscal/profile"
)

//CreateCatalogsFromProfile maps profile controls to multiple catalogs
func CreateCatalogsFromProfile(profile *profile.Profile) []*catalog.Catalog {

	var outputCatalogs []*catalog.Catalog
	//Get first import of the profile (which is a catalog)
	for _, profileImport := range profile.Imports {

		//ForEach Import's Href, Fetch the Catalog JSON file
		catalogReference, err := GetCatalogFilePath(profileImport.Href.String())
		if err != nil {
			logrus.Errorf("invalid file path: %v", err)
			continue
		}

		f, err := os.Open(catalogReference)
		if err != nil {
			logrus.Errorf("cannot read file: %v", err)
			continue
		}
		defer func() {
			err := os.Remove(catalogReference)
			if err != nil {
				log.Fatal(err)
			}
		}()
		//Once fetched, Read the catalog JSON and Marshall it to Go struct.
		importedCatalog, err := ReadCatalog(f)
		if err != nil {
			logrus.Errorf("cannot parse catalog listed in import.href %v", err)
		}
		//Prepare a new catalog object to merge into the final List of OutputCatalogs
		newCatalog := getMappedCatalogControlsFromImport(importedCatalog, profileImport)
		outputCatalogs = append(outputCatalogs, &newCatalog)
	}
	return outputCatalogs
}

func getMappedCatalogControlsFromImport(importedCatalog *catalog.Catalog, profileImport profile.Import) catalog.Catalog {
	newCatalog := catalog.Catalog{Groups: []catalog.Group{}}
	for _, group := range importedCatalog.Groups {
		//Prepare a new group to append matching controls into.
		newGroup := catalog.Group{
			Title:    group.Title,
			Controls: []catalog.Control{},
		}
		//Append controls to the new group if matches
		for _, catalogControl := range group.Controls {
			for _, z := range profileImport.Include.IdSelectors {
				if catalogControl.Id == z.ControlId {
					newGroup.Controls = append(newGroup.Controls, catalog.Control{
						Id:          catalogControl.Id,
						Class:       catalogControl.Class,
						Title:       catalogControl.Title,
						Subcontrols: catalogControl.Subcontrols,
					})
					spew.Dump(newGroup.Controls[0].Subcontrols)
				}
			}
		}
		if len(newGroup.Controls) > 0 {
			newCatalog.Groups = append(newCatalog.Groups, newGroup)
		}
	}
	return newCatalog
}
