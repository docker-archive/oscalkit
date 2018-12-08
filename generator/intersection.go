package generator

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/opencontrol/oscalkit/types/oscal/catalog"
	"github.com/opencontrol/oscalkit/types/oscal/profile"
)

//IntersectProfile IntersectProfile
func IntersectProfile(profile *profile.Profile) []*catalog.Catalog {

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
			logrus.Errorf("cannot open file: %v", err)
			continue
		}
		//Once fetched, Read the catalog JSON and Marshall it to Go struct.
		importedCatalog, err := ReadCatalog(f)
		if err != nil {
			logrus.Errorf("cannot parse catalog listed in import.href %v", err)
		}
		//Prepare a new catalog object to merge into the final List of OutputCatalogs
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
						newGroup.Controls = append(newGroup.Controls, catalogControl)
					}
				}
			}
			if len(newGroup.Controls) > 0 {
				newCatalog.Groups = append(newCatalog.Groups, newGroup)
			}
		}
		outputCatalogs = append(outputCatalogs, &newCatalog)
	}
	return outputCatalogs

}
