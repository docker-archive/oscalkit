package impl

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/opencontrol/oscalkit/types/oscal/catalog"
	"github.com/opencontrol/oscalkit/types/oscal/profile"
)

const (
	temporaryProfileID = "oscaltestprofile"
	catalogRef         = "https://raw.githubusercontent.com/usnistgov/OSCAL/master/content/nist.gov/SP800-53/rev4/NIST_SP-800-53_rev4_catalog.xml"
)

func TestGenerateImplementation(t *testing.T) {

	components := []string{"CompA", "CompB", "CompC"}
	controls := []string{"ac-2", "ac-2.2", "ac-4"}
	ComponentDetails := [][]string{
		[]string{controls[0], fmt.Sprintf("%s|%s", components[0], components[1]), "2-Narrative"},
		[]string{controls[1], fmt.Sprintf("%s|%s", components[0], components[1]), "2.2-Narrative"},
		[]string{controls[2], "CompC", "4-Narrative"},
	}
	csvs := make([][]string, TotalControlsInExcel)
	for i := range csvs {
		csvs[i] = make([]string, 20)
	}
	for i, x := range ComponentDetails {
		csvs[i+RowIndex][ControlIndex] = x[0]
		csvs[i+RowIndex][ComponentNameIndex] = x[1]
		csvs[i+RowIndex][NarrativeIndex] = x[2]
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
	if len(i.ComponentDefinitions[0].ComponentConfigurations) != len(components) {
		t.Error("mismatch number of components")
	}
	if len(i.ComponentDefinitions[0].ControlImplementations[0].ControlIds) != len(controls) {
		t.Error("mismatch number of controls")
	}

}
