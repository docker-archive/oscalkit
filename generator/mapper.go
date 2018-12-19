package generator

import (
	"os"
	"strings"

	"github.com/Sirupsen/logrus"

	"github.com/opencontrol/oscalkit/types/oscal/catalog"
	"github.com/opencontrol/oscalkit/types/oscal/profile"
)

//CreateCatalogsFromProfile maps profile controls to multiple catalogs
func CreateCatalogsFromProfile(profileArg *profile.Profile) ([]*catalog.Catalog, error) {

	done := 0
	errChan := make(chan error)
	doneChan := make(chan *catalog.Catalog)
	var outputCatalogs []*catalog.Catalog
	//Get first import of the profile (which is a catalog)
	for _, profileImport := range profileArg.Imports {
		go func(profileImport profile.Import) {
			//ForEach Import's Href, Fetch the Catalog JSON file
			catalogReference, err := GetCatalogFilePath(profileImport.Href.String())
			if err != nil {
				logrus.Errorf("invalid file path: %v", err)
				doneChan <- nil
				return
			}

			f, err := os.Open(catalogReference)
			if err != nil {
				logrus.Errorf("cannot read file: %v", err)
				doneChan <- nil
				return
			}

			//Once fetched, Read the catalog JSON and Marshall it to Go struct.
			importedCatalog, err := ReadCatalog(f)
			if err != nil {
				logrus.Errorf("cannot parse catalog listed in import.href %v", err)
				errChan <- err
				return

			}
			//Prepare a new catalog object to merge into the final List of OutputCatalogs
			newCatalog, err := getMappedCatalogControlsFromImport(importedCatalog, profileImport)
			if err != nil {
				errChan <- err
				return
			}

			doneChan <- &newCatalog
		}(profileImport)

	}
	for {
		select {
		case err := <-errChan:
			return nil, err
		case newCatalog := <-doneChan:
			done++
			if newCatalog != nil {
				outputCatalogs = append(outputCatalogs, newCatalog)
			}
			if done == len(profileArg.Imports) {
				return outputCatalogs, nil
			}
		}
	}
}

func getMappedCatalogControlsFromImport(importedCatalog *catalog.Catalog, profileImport profile.Import) (catalog.Catalog, error) {
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
			/**
			wrapped it around a function to achieve immutability
			profile.Include is a pointer hence, dereferencing
			*/
			func(include profile.Include) {
				for controlIndex, call := range include.IdSelectors {
					if strings.ToLower(catalogControl.Id) == strings.ToLower(call.ControlId) {
						newControl := catalog.Control{
							Id:          catalogControl.Id,
							Class:       catalogControl.Class,
							Title:       catalogControl.Title,
							Subcontrols: []catalog.Subcontrol{},
						}
						//For subcontrols, find again in entire profile
						for _, catalogSubControl := range catalogControl.Subcontrols {
							for subcontrolIndex, subcontrol := range include.IdSelectors {
								//if found append the subcontrol into the control attribute
								if strings.ToLower(catalogSubControl.Id) == strings.ToLower(subcontrol.SubcontrolId) {
									newControl.Subcontrols = append(newControl.Subcontrols, catalogSubControl)
									//remove that subcontrol from profile. (for less computation)
									include.IdSelectors = append(include.IdSelectors[:subcontrolIndex], include.IdSelectors[subcontrolIndex+1:]...)
									break
								}
							}
						}
						//finally append the control in the group.
						newGroup.Controls = append(newGroup.Controls, newControl)
						//remove controlId from profile as well. (for less computation)
						include.IdSelectors = append(include.IdSelectors[:controlIndex], profileImport.Include.IdSelectors[controlIndex+1:]...)
						break
					}
				}

			}(*profileImport.Include)
		}
		if len(newGroup.Controls) > 0 {
			newCatalog.Groups = append(newCatalog.Groups, newGroup)
		}
	}
	return newCatalog, nil
}
