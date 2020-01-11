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

	logrus.Debug("processing alteration and parameters... \nmapping to controls...")
	// Get first import of the profile (which is a catalog)
	for _, profileImport := range profileArg.Imports {
		err := ValidateHref(profileImport.Href)
		if err != nil {
			return nil, err
		}
		go func(profileImport profile.Import) {
			catalogHelper := impl.NISTCatalog{}
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
				newCatalog, err := GetMappedCatalogControlsFromImport(importedCatalog, profileImport, &catalogHelper)
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
func getSubControl(call profile.Call, ctrls []catalog.Control, helper impl.Catalog) (catalog.Control, error) {
	for _, ctrl := range ctrls {
		if ctrl.Id == helper.GetControl(call.ControlId) {
			for _, subctrl := range ctrl.Controls {
				if subctrl.Id == call.ControlId {
					return subctrl, nil
				}
			}
		}
	}
	return catalog.Control{}, fmt.Errorf("could not find subcontrol %s in catalog", call.ControlId)
}

// GetMappedCatalogControlsFromImport gets mapped controls in catalog per profile import
func GetMappedCatalogControlsFromImport(importedCatalog *catalog.Catalog, profileImport profile.Import, catalogHelper impl.Catalog) (catalog.Catalog, error) {
	newCatalog := catalog.Catalog{
		Title:  importedCatalog.Title,
		Groups: []catalog.Group{},
	}

	for _, group := range importedCatalog.Groups {
		newGroup := catalog.Group{
			Title:    group.Title,
			Controls: []catalog.Control{},
		}
		for _, ctrl := range group.Controls {
			for _, call := range profileImport.Include.IdSelectors {
				if call.ControlId == "" {
					if strings.ToLower(ctrl.Id) == strings.ToLower(catalogHelper.GetControl(call.ControlId)) {
						ctrlExistsInGroup := false
						sc, err := getSubControl(call, group.Controls, &impl.NISTCatalog{})
						if err != nil {
							return catalog.Catalog{}, err
						}
						for i, mappedCtrl := range newGroup.Controls {
							if mappedCtrl.Id == strings.ToLower(catalogHelper.GetControl(call.ControlId)) {
								ctrlExistsInGroup = true
								newGroup.Controls[i].Controls = append(newGroup.Controls[i].Controls, sc)
							}
						}
						if !ctrlExistsInGroup {
							newGroup.Controls = append(newGroup.Controls,
								catalog.Control{
									Id:         ctrl.Id,
									Class:      ctrl.Class,
									Title:      ctrl.Title,
									Parameters: ctrl.Parameters,
									Parts:      ctrl.Parts,
									Controls:   []catalog.Control{sc},
								})
						}
					}
				}
				if strings.ToLower(call.ControlId) == strings.ToLower(ctrl.Id) {
					ctrlExists := false
					for _, x := range newGroup.Controls {
						if x.Id == ctrl.Id {
							ctrlExists = true
							continue
						}
					}
					if !ctrlExists {
						newGroup.Controls = append(newGroup.Controls,
							catalog.Control{
								Id:         ctrl.Id,
								Class:      ctrl.Class,
								Title:      ctrl.Title,
								Controls:   []catalog.Control{},
								Parameters: ctrl.Parameters,
								Parts:      ctrl.Parts,
							},
						)
					}
				}
			}
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
