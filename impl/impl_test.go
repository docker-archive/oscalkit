package impl

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/docker/oscalkit/types/oscal/catalog"
	"github.com/docker/oscalkit/types/oscal/profile"
)

const (
	temporaryProfileID = "oscaltestprofile"
	catalogRef         = "https://raw.githubusercontent.com/usnistgov/OSCAL/master/content/nist.gov/SP800-53/rev4/NIST_SP-800-53_rev4_catalog.xml"
)

func TestGenerateImplementation(t *testing.T) {
	ucpCompID := "(component ID: cpe:2.3:a:docker:ucp:3.2.0:*:*:*:*:*:*:*)"
	dtrCompID := "(component ID: cpe:2.3:a:docker:dtr:2.7.0:*:*:*:*:*:*:*)"
	engineCompID := "(component ID: cpe:2.3:a:docker:engine-enterprise:18.09:*:*:*:*:*:*:*)"
	components := []string{"CompA", "CompB", "CompC"}
	comps := []componenet{
		{
			id:             getComponentID(ucpCompID),
			name:           "UCP",
			compNameIndex:  17,
			uuidIndex:      18,
			narrativeIndex: 19,
			definition:     make(cdMap),
		},
		{
			id:             getComponentID(dtrCompID),
			name:           "DTR",
			compNameIndex:  20,
			uuidIndex:      21,
			narrativeIndex: 22,
			definition:     make(cdMap),
		},
		{
			id:             getComponentID(engineCompID),
			name:           "Engine",
			compNameIndex:  14,
			uuidIndex:      15,
			narrativeIndex: 16,
			definition:     make(cdMap),
		},
	}
	controls := []string{"ac-2", "ac-2.2", "ac-4", "bc-1.1", "hk-1.2", "as-3.2", "af-1.23", "ar-5.2", "fp-8.5"}
	ComponentDetails := [][]string{
		[]string{controls[0], fmt.Sprintf("%s|%s", components[0], components[1]), "2-Narrative", "123|321"},
		[]string{controls[1], fmt.Sprintf("%s|%s", components[0], components[1]), "2.2-Narrative", "456|654"},
		[]string{controls[2], "CompC", "4-Narrative", "789|987"},
		[]string{controls[0], fmt.Sprintf("%s|%s", components[0], components[1]), "3-Narrative", "123|321"},
		[]string{controls[1], fmt.Sprintf("%s|%s", components[0], components[1]), "3.4-Narrative", "567|1231"},
		[]string{controls[2], "CompC", "5-Narrative", "789|987"},
		[]string{controls[0], fmt.Sprintf("%s|%s", components[0], components[1]), "4-Narrative", "123|321"},
		[]string{controls[1], fmt.Sprintf("%s|%s", components[0], components[1]), "4.1-Narrative", "111|222"},
		[]string{controls[2], "CompC", "6-Narrative", "789|987"},
	}
	csvs := make([][]string, totalControlsInExcel)
	for i := range csvs {
		csvs[i] = make([]string, 25)
	}

	for _, comp := range comps {
		for i, x := range ComponentDetails {
			csvs[i+rowIndex][controlIndex] = x[0]
			csvs[i+rowIndex][comp.compNameIndex] = x[1]
			csvs[i+rowIndex][comp.narrativeIndex] = x[2]
			csvs[i+rowIndex][comp.uuidIndex] = x[3]
		}
	}

	p := profile.Profile{
		ID: temporaryProfileID,
		Imports: []profile.Import{
			profile.Import{
				Href: &catalog.Href{
					URL: func() *url.URL {
						uri, _ := url.Parse(catalogRef)
						return uri
					}(),
				},
				Include: &profile.Include{
					IdSelectors: []profile.Call{
						profile.Call{
							ControlId: "ac-2",
						},
						profile.Call{
							ControlId: "ac-4",
						},
						profile.Call{
							SubcontrolId: "ac-4.2",
						},
					},
				},
			},
		},
		Modify: &profile.Modify{
			ParamSettings: []profile.SetParam{
				profile.SetParam{
					Id:          "ac-2_prm",
					Constraints: []catalog.Constraint{catalog.Constraint{Value: "some constraint"}},
				},
				profile.SetParam{
					Id:          "ac-2_prm_obj",
					Constraints: []catalog.Constraint{catalog.Constraint{Value: "some constraint"}},
				},
				profile.SetParam{
					Id:          "",
					Constraints: []catalog.Constraint{},
				},
				profile.SetParam{
					Id:          "ac-4_prm",
					Constraints: []catalog.Constraint{},
				},
			},
		},
	}
	i := GenerateImplementation(csvs, &p, &NISTCatalog{"NISTSP80053"})

	if len(i.ComponentDefinitions) != len(comps) {
		t.Error("mismatch number of component definitions")
	}
	if len(i.ComponentDefinitions[0].ComponentConfigurations) != len(components) {
		t.Error("mismatch number of components")
	}
	if len(i.ComponentDefinitions[0].ControlImplementations[0].ControlIds) != len(controls) {
		t.Error("mismatch number of controls")
	}

}
