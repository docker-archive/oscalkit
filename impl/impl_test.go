package impl

import (
	"fmt"
	"testing"

	"github.com/docker/oscalkit/types/oscal/implementation"
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
	comps := []component{
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

	i := GenerateImplementation(csvs, &NISTCatalog{"NISTSP80053"})

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

func TestFindOrCreateImplementsProfile(t *testing.T) {
	cd := implementation.ComponentDefinition{}
	profileID := "123"
	implementsProfile := findOrCreateImplementsProfile(&cd, profileID)
	if implementsProfile == nil {
		t.Error("nothing created")
		return
	}
	for _, x := range cd.ImplementsProfiles {
		if x.ProfileID == profileID {
			return
		}
	}
	t.Error("profile did not get appended")
}

func TestFindOrCreateControlConfig(t *testing.T) {
	cd := implementation.ComponentDefinition{
		ImplementsProfiles: []*implementation.ImplementsProfile{
			&implementation.ImplementsProfile{
				ProfileID: "123",
			},
		},
	}
	configIDRef := "some-guid"
	controlConfig := findOrCreateControlConfig(cd.ImplementsProfiles[0], configIDRef)
	if controlConfig == nil {
		t.Error("empty ctrl conf")
	}
	for _, x := range cd.ImplementsProfiles[0].ControlConfigurations {
		if x.ConfigurationIDRef == configIDRef {
			return
		}
	}
	t.Error("coudnt find conf id")
}

func TestMapImplementsProfile(t *testing.T) {

	parameterID := "ac-10_prm_3"
	profileID := "uuid-fedramp-high-20180806-195540"
	configurationIDRef := "random-guid-id"
	Value := "<=2"
	checkAndValue := fmt.Sprintf("MakeItRight(%s)", Value)
	cd := implementation.ComponentDefinition{
		ComponentConfigurations: []*implementation.ComponentConfiguration{
			&implementation.ComponentConfiguration{
				ID:   configurationIDRef,
				Name: "MakeItRight",
			},
		},
	}
	mapImplementsProfile(&cd, parameterID, profileID, checkAndValue)
	if len(cd.ImplementsProfiles) < 1 {
		t.Error("implements profile array should have one element")
	}
	if cd.ImplementsProfiles[0].ProfileID != profileID {
		t.Errorf("profile id should be %s", profileID)
	}
	if len(cd.ImplementsProfiles[0].ControlConfigurations) < 1 {
		t.Error("control configurations array should have one element")
	}
	if cd.ImplementsProfiles[0].ControlConfigurations[0].ConfigurationIDRef != configurationIDRef {
		t.Errorf("config ref id should be %s", configurationIDRef)
	}
	if len(cd.ImplementsProfiles[0].ControlConfigurations[0].Parameters) < 1 {
		t.Error("parameters array should have one element")
	}
	if cd.ImplementsProfiles[0].ControlConfigurations[0].Parameters[0].ParameterID != parameterID {
		t.Errorf("parameter id should be %s", parameterID)
	}
	if len(cd.ImplementsProfiles[0].ControlConfigurations[0].Parameters[0].PossibleValues) < 1 {
		t.Error("possible value array should have one element")
	}
	if cd.ImplementsProfiles[0].ControlConfigurations[0].Parameters[0].PossibleValues[0] != Value {
		t.Errorf("possible value should be %s", Value)
	}
}

func TestMapImplementsProfileWithMultiplePossibleValues(t *testing.T) {

	parameterID := "ac-10_prm_3"
	profileID := "uuid-fedramp-high-20180806-195540"
	configurationIDRef := "random-guid-id"
	Value := "<=2"
	Value2 := "<=3"
	checkAndValue := fmt.Sprintf("MakeItRight(%s)", Value)
	checkAndValue2 := fmt.Sprintf("MakeItRight(%s)", Value2)

	cd := implementation.ComponentDefinition{
		ComponentConfigurations: []*implementation.ComponentConfiguration{
			&implementation.ComponentConfiguration{
				ID:   configurationIDRef,
				Name: "MakeItRight",
			},
		},
	}
	mapImplementsProfile(&cd, parameterID, profileID, checkAndValue)
	mapImplementsProfile(&cd, parameterID, profileID, checkAndValue2)
	if len(cd.ImplementsProfiles) < 1 {
		t.Error("implements profile array should have one element")
	}
	if cd.ImplementsProfiles[0].ProfileID != profileID {
		t.Errorf("profile id should be %s", profileID)
	}
	if len(cd.ImplementsProfiles[0].ControlConfigurations) < 1 {
		t.Error("control configurations array should have one element")
	}
	if cd.ImplementsProfiles[0].ControlConfigurations[0].ConfigurationIDRef != configurationIDRef {
		t.Errorf("config ref id should be %s", configurationIDRef)
	}
	if len(cd.ImplementsProfiles[0].ControlConfigurations[0].Parameters) < 1 {
		t.Error("parameters array should have one element")
	}
	if cd.ImplementsProfiles[0].ControlConfigurations[0].Parameters[0].ParameterID != parameterID {
		t.Errorf("parameter id should be %s", parameterID)
	}
	if len(cd.ImplementsProfiles[0].ControlConfigurations[0].Parameters[0].PossibleValues) < 2 {
		t.Error("possible value array should have one element")
	}
	if cd.ImplementsProfiles[0].ControlConfigurations[0].Parameters[0].PossibleValues[0] != Value {
		t.Errorf("possible value should be %s", Value)
	}
	if cd.ImplementsProfiles[0].ControlConfigurations[0].Parameters[0].PossibleValues[1] != Value2 {
		t.Errorf("possible value should be %s", Value)
	}
}

func TestMapImplementsProfileWithMultipleProfiles(t *testing.T) {

	parameterID := "ac-10_prm_3"
	profileID := "uuid-fedramp-high-20180806-195540"
	profileID2 := "uuid-fedramp-moderate-20180806-195540"
	configurationIDRef := "random-guid-id"
	Value := "<=2"
	checkAndValue := fmt.Sprintf("MakeItRight(%s)", Value)

	cd := implementation.ComponentDefinition{
		ComponentConfigurations: []*implementation.ComponentConfiguration{
			&implementation.ComponentConfiguration{
				ID:   configurationIDRef,
				Name: "MakeItRight",
			},
		},
	}
	mapImplementsProfile(&cd, parameterID, profileID, checkAndValue)
	mapImplementsProfile(&cd, parameterID, profileID2, checkAndValue)
	if len(cd.ImplementsProfiles) < 2 {
		t.Error("implements profile array should have one element")
	}
	if cd.ImplementsProfiles[0].ProfileID != profileID {
		t.Errorf("profile id should be %s", profileID)
	}
	if cd.ImplementsProfiles[1].ProfileID != profileID2 {
		t.Errorf("profile id should be %s", profileID2)
	}
	if len(cd.ImplementsProfiles[0].ControlConfigurations) < 1 {
		t.Error("control configurations array should have one element")
	}
	if cd.ImplementsProfiles[0].ControlConfigurations[0].ConfigurationIDRef != configurationIDRef {
		t.Errorf("config ref id should be %s", configurationIDRef)
	}
	if len(cd.ImplementsProfiles[0].ControlConfigurations[0].Parameters) < 1 {
		t.Error("parameters array should have one element")
	}
	if cd.ImplementsProfiles[0].ControlConfigurations[0].Parameters[0].ParameterID != parameterID {
		t.Errorf("parameter id should be %s", parameterID)
	}
	if len(cd.ImplementsProfiles[0].ControlConfigurations[0].Parameters[0].PossibleValues) < 1 {
		t.Error("possible value array should have one element")
	}
	if cd.ImplementsProfiles[0].ControlConfigurations[0].Parameters[0].PossibleValues[0] != Value {
		t.Errorf("possible value should be %s", Value)
	}
}

func TestMapImplementsProfileWithMultipleParameters(t *testing.T) {

	parameterID := "ac-10_prm_3"
	parameterID2 := "ac-10_prm_2"
	profileID := "uuid-fedramp-high-20180806-195540"
	configurationIDRef := "random-guid-id"
	Value := "<=2"
	checkAndValue := fmt.Sprintf("MakeItRight(%s)", Value)

	cd := implementation.ComponentDefinition{
		ComponentConfigurations: []*implementation.ComponentConfiguration{
			&implementation.ComponentConfiguration{
				ID:   configurationIDRef,
				Name: "MakeItRight",
			},
		},
	}
	mapImplementsProfile(&cd, parameterID, profileID, checkAndValue)
	mapImplementsProfile(&cd, parameterID2, profileID, checkAndValue)
	if len(cd.ImplementsProfiles) < 1 {
		t.Error("implements profile array should have one element")
	}
	if cd.ImplementsProfiles[0].ProfileID != profileID {
		t.Errorf("profile id should be %s", profileID)
	}

	if len(cd.ImplementsProfiles[0].ControlConfigurations) < 1 {
		t.Error("control configurations array should have one element")
	}
	if cd.ImplementsProfiles[0].ControlConfigurations[0].ConfigurationIDRef != configurationIDRef {
		t.Errorf("config ref id should be %s", configurationIDRef)
	}
	if len(cd.ImplementsProfiles[0].ControlConfigurations[0].Parameters) < 2 {
		t.Error("parameters array should have one element")
	}
	if cd.ImplementsProfiles[0].ControlConfigurations[0].Parameters[0].ParameterID != parameterID {
		t.Errorf("parameter id should be %s", parameterID)
	}
	if cd.ImplementsProfiles[0].ControlConfigurations[0].Parameters[1].ParameterID != parameterID2 {
		t.Errorf("parameter id should be %s", parameterID2)
	}
	if len(cd.ImplementsProfiles[0].ControlConfigurations[0].Parameters[0].PossibleValues) < 1 {
		t.Error("possible value array should have one element")
	}
	if cd.ImplementsProfiles[0].ControlConfigurations[0].Parameters[0].PossibleValues[0] != Value {
		t.Errorf("possible value should be %s", Value)
	}
}

func TestGetProfileIDWithValidProfile(t *testing.T) {
	x := "FedRAMP_High"
	o := getProfileID(x)
	if o != profileMap[x] {
		t.Error("failed to map profile id")
	}
}
func TestGetProfileIDWithInvalidProfile(t *testing.T) {
	x := "123"
	o := getProfileID(x)
	if o == profileMap[x] {
		t.Error("mapped invalid profile id")
	}
}

func TestDetokenizeParameterString(t *testing.T) {

	x := "FedRAMP_High->SetParam(5)"
	p := "FedRAMP_High"
	c := "SetParam(5)"
	profileID, checkAndValue := detokenizeParameterString(x)
	if profileMap[p] != profileID || c != checkAndValue {
		t.Errorf("failed to tokenize parameter string %s| output %s:%s", x, profileMap[p], checkAndValue)
	}
}
