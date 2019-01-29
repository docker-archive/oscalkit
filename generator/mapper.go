package generator

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/docker/oscalkit/impl"
	"github.com/docker/oscalkit/types/oscal"
	"github.com/docker/oscalkit/types/oscal/catalog"
	"github.com/docker/oscalkit/types/oscal/profile"
	"github.com/sirupsen/logrus"
)

// CreateCatalogsFromProfile maps profile controls to multiple catalogs
func CreateCatalogsFromProfile(profileArg *profile.Profile) ([]*catalog.Catalog, error) {

	t := time.Now()
	done := 0
	errChan := make(chan error)
	catalogChan := make(chan *catalog.Catalog)
	var outputCatalogs []*catalog.Catalog
	logrus.Info("fetching alterations...")
	alterations, err := GetAlters(profileArg)
	if err != nil {
		return nil, err
	}
	logrus.Info("fetching alterations from import chain complete")

	logrus.Info("processing alteration and parameters... \nmapping to controls...")
	// Get first import of the profile (which is a catalog)
	for _, profileImport := range profileArg.Imports {
		err := ValidateHref(profileImport.Href)
		if err != nil {
			return nil, err
		}
		go func(profileImport profile.Import) {
			c := make(chan *catalog.Catalog)
			e := make(chan error)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			// ForEach Import's Href, Fetch the Catalog JSON file
			getCatalogForImport(ctx, profileImport, c, e, profileImport.Href.String())
			select {
			case importedCatalog := <-c:
				// Prepare a new catalog object to merge into the final List of OutputCatalogs
				if profileArg.Modify != nil {
					nc := impl.NISTCatalog{}
					importedCatalog = ProcessAlterations(alterations, importedCatalog)
					importedCatalog = ProcessSetParam(profileArg.Modify.ParamSettings, importedCatalog, &nc)
				}
				newCatalog, err := GetMappedCatalogControlsFromImport(importedCatalog, profileImport)
				if err != nil {
					errChan <- err
					return
				}
				catalogChan <- &newCatalog

			case err := <-e:
				errChan <- err
				return
			}

		}(profileImport)

	}
	for {
		select {
		case err := <-errChan:
			return nil, err
		case newCatalog := <-catalogChan:
			done++
			if newCatalog != nil {
				outputCatalogs = append(outputCatalogs, newCatalog)
			}
			if done == len(profileArg.Imports) {
				logrus.Infof("successfully mapped controls in %f seconds", time.Since(t).Seconds())
				return outputCatalogs, nil
			}
		}
	}
}

// GetMappedCatalogControlsFromImport gets mapped controls in catalog per profile import
func GetMappedCatalogControlsFromImport(importedCatalog *catalog.Catalog, profileImport profile.Import) (catalog.Catalog, error) {
	newCatalog := catalog.Catalog{
		Title:  importedCatalog.Title,
		Groups: []catalog.Group{},
	}
	for _, group := range importedCatalog.Groups {
		// Prepare a new group to append matching controls into.
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
							Params:      catalogControl.Params,
							Parts:       catalogControl.Parts,
						}
						// For subcontrols, find again in entire profile
						for _, catalogSubControl := range catalogControl.Subcontrols {
							for subcontrolIndex, subcontrol := range include.IdSelectors {
								// If found append the subcontrol into the control attribute
								if strings.ToLower(catalogSubControl.Id) == strings.ToLower(subcontrol.SubcontrolId) {
									newControl.Subcontrols = append(newControl.Subcontrols, catalog.Subcontrol{
										Id:    catalogSubControl.Id,
										Class: catalogSubControl.Class,
										Title: catalogSubControl.Title,
										Parts: catalogSubControl.Parts,
									})
									// Remove that subcontrol from profile. (for less computation)
									include.IdSelectors = append(include.IdSelectors[:subcontrolIndex], include.IdSelectors[subcontrolIndex+1:]...)
									break
								}
							}
						}
						// Finally append the control in the group.
						newGroup.Controls = append(newGroup.Controls, newControl)
						// Remove controlId from profile as well. (for less computation)
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

func getCatalogForImport(ctx context.Context, i profile.Import, c chan *catalog.Catalog, e chan error, basePath string) {
	go func(i profile.Import) {
		err := ValidateHref(i.Href)
		if err != nil {
			e <- fmt.Errorf("href cannot be nil")
			return
		}
		path, err := GetFilePath(i.Href.String())
		if err != nil {
			e <- err
			return
		}
		f, err := os.Open(path)
		if err != nil {
			e <- err
			return
		}
		defer f.Close()
		o, err := oscal.New(f)
		if err != nil {
			e <- err
			return
		}
		if o.Catalog != nil {
			c <- o.Catalog
			return
		}
		newP, err := SetBasePath(o.Profile, basePath)
		if err != nil {
			e <- err
			return
		}
		o.Profile = newP
		for _, p := range o.Profile.Imports {
			go func(p profile.Import) {
				getCatalogForImport(ctx, p, c, e, basePath)
			}(p)
		}
	}(i)
}
