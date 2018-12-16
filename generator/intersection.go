package generator

import (
	"fmt"
	"os"
	"strings"

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

func mapSubControls(mappedCatalogs *catalog.Catalog, calls []profile.Call) {
	for i, g := range mappedCatalogs.Groups {
		for j, control := range g.Controls {
			for _, call := range calls {
				if control.Id == getControlIDFromSubControl(call.SubcontrolId) {
					mappedCatalogs.Groups[i].Controls[j].Subcontrols = append(mappedCatalogs.Groups[i].Controls[j].Subcontrols, catalog.Subcontrol{
						Id: call.SubcontrolId,
					})
				}
			}
		}
	}

}
func isSubControl(s string) bool {
	substrings := []string{" ", "(", "."}
	for _, substr := range substrings {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}

func getControlIDFromSubControl(sc string) string {
	if len(sc) >= 4 {
		subc := strings.Split(sc, "-")
		if isSUbControl(subc[1][0:2]) {
			return fmt.Sprintf("%s-%s", subc[0], subc[1][0:1])
		}
		return fmt.Sprintf("%s-%s", subc[0], subc[1][0:2])
	}
	return sc

}

func getMappedCatalogControlsFromImport(importedCatalog *catalog.Catalog, profileImport profile.Import) catalog.Catalog {
	newCatalog := catalog.Catalog{
		Title:  importedCatalog.Title,
		Groups: []catalog.Group{},
	}
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
						Subcontrols: []catalog.Subcontrol{},
					})
				}
			}
		}
		if len(newGroup.Controls) > 0 {
			newCatalog.Groups = append(newCatalog.Groups, newGroup)
		}
	}
	mapSubControls(&newCatalog, profileImport.Include.IdSelectors)
	return newCatalog
}
