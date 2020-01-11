package impl

import (
	"regexp"
	"strings"

	"github.com/docker/oscalkit/types/oscal/catalog"
	"github.com/docker/oscalkit/types/oscal/implementation"
	"github.com/docker/oscalkit/types/oscal/profile"
	uuid "github.com/satori/go.uuid"
)

const (
	// totalControlsInExcel the total number of controls in the excel sheet
	totalControlsInExcel = 264
	// controlIndex Column at which control is present in the excel sheet
	controlIndex = 2
	// rowIndex Starting point for valid rows (neglects titles)
	rowIndex         = 3
	delimiter        = "|"
	profileDelimiter = "->"
	componentIDRegex = `cpe:[0-9].[0-9]:[a-z]:docker:[a-z-]*:(\d+\.)?(\d+\.)?(\*|\d+)`
	// componentNameRow is the index for getting component name
	componentNameRow = 1
	// NIST subControlID Regex
	subControlIDRegex = "[a-zA-Z]{2}-\\d{1,2}\\.\\d{1,2}"
)

type guidMap map[string]uuid.UUID
type cdMap map[string]implementation.ComponentDefinition
type component struct {
	id                   string
	compNameIndex        int
	name                 string
	parameterIDIndex     int
	parameterStringIndex int
	uuidIndex            int
	narrativeIndex       int
	definition           cdMap
	hasParameterMapping  bool
}

var profileMap = map[string]string{
	"FedRAMP_High":     "uuid-fedramp-high-20180806-195540",
	"FedRAMP_HIGH":     "uuid-fedramp-high-20180806-195540",
	"FedRAMP_Moderate": "uuid-fedramp-moderate-20180806-195542",
	"FedRAMP_moderate": "uuid-fedramp-moderate-20180806-195542",
}

// GenerateImplementation generates implementation from component excel sheet
func GenerateImplementation(CSVS [][]string, c Catalog) implementation.Implementation {
	var cdMapList = make([]cdMap, 0)
	components := []component{
		{
			id:                   getComponentID(CSVS[componentNameRow][17]),
			name:                 "UCP",
			compNameIndex:        17,
			parameterIDIndex:     18,
			parameterStringIndex: 19,
			uuidIndex:            20,
			narrativeIndex:       21,
			definition:           make(cdMap),
			hasParameterMapping:  true,
		},
		{
			id:                  getComponentID(CSVS[componentNameRow][22]),
			name:                "DTR",
			compNameIndex:       22,
			uuidIndex:           23,
			narrativeIndex:      24,
			definition:          make(cdMap),
			hasParameterMapping: false,
		},
		{
			id:                  getComponentID(CSVS[componentNameRow][14]),
			name:                "Engine",
			compNameIndex:       14,
			uuidIndex:           15,
			narrativeIndex:      16,
			definition:          make(cdMap),
			hasParameterMapping: false,
		},
	}
	for _, comp := range components {
		//comp.compNameIndex, comp.definition, comp.uuidIndex, comp.narrativeIndex, comp.parameterIDIndex, comp.parameterStringIndex, comp.id
		cdMapList = append(cdMapList, fillCDMap(CSVS, comp, c))
	}
	return CompileImplementation(cdMapList, CSVS, c, components)
}

func fillCDMap(CSVS [][]string, comp component, c Catalog) cdMap {
	checkAgainstGUID := make(map[string]uuid.UUID)
	for i := rowIndex; i < totalControlsInExcel; i++ {
		applicableControl := CSVS[i][controlIndex]
		if applicableControl == "" {
			continue
		}
		applicableNarrative := CSVS[i][comp.narrativeIndex]
		parameterID := CSVS[i][comp.parameterIDIndex]
		parameterString := CSVS[i][comp.parameterStringIndex]
		ListOfComponentConfigName := strings.Split(CSVS[i][comp.compNameIndex], delimiter)
		for compIndex, componentConfigName := range ListOfComponentConfigName {
			componentConfigName = strings.TrimSpace(componentConfigName)
			if componentConfigName == "" {
				continue
			}
			if _, ok := comp.definition[componentConfigName]; !ok {

				guid := strings.Split(CSVS[i][comp.uuidIndex], delimiter)[compIndex]
				guid = strings.TrimSpace(guid)
				CreateComponentDefinition(checkAgainstGUID, comp.definition, componentConfigName, c, applicableControl, applicableNarrative, guid, comp.id, parameterID, parameterString)
			} else {
				securityCheck := comp.definition[componentConfigName]
				guid := checkAgainstGUID[componentConfigName]
				temp := AppendControlInImplementation(securityCheck, guid, c, applicableControl)
				comp.definition[componentConfigName] = temp
			}
		}
	}
	return comp.definition
}

// CreateComponentDefinition creates a component definition
func CreateComponentDefinition(gm guidMap, cdm cdMap, componentConfName string, c Catalog, control, narrative, guid string, cdID string, parameterID, parameterString string) {

	componentConfGUID, _ := uuid.FromString(guid)
	gm[componentConfName] = componentConfGUID
	controlConfiguration := implementation.ControlConfiguration{
		ConfigurationIDRef: componentConfGUID.String(),
	}
	controlConfiguration.ProvisioningMechanisms = []implementation.ProvisioningMechanism{
		implementation.ProvisioningMechanism{
			ProvisionedControls: []implementation.ControlId{
				implementation.ControlId{
					ControlID:    c.GetControl(control),
					CatalogIDRef: c.GetID(),
					ItemID:       "",
				},
			},
		},
	}
	cdm[componentConfName] = implementation.ComponentDefinition{
		ID: cdID,
		ComponentConfigurations: []*implementation.ComponentConfiguration{
			CreateComponentConfiguration(componentConfGUID, componentConfName, narrative),
		},
		ImplementsProfiles: []*implementation.ImplementsProfile{},
		ControlImplementations: []*implementation.ControlImplementation{
			&implementation.ControlImplementation{
				ControlConfigurations: []implementation.ControlConfiguration{
					controlConfiguration,
				},
			},
		},
	}

}

// CreateComponentConfiguration creates component configuration
func CreateComponentConfiguration(guid uuid.UUID, componentConfName, narrative string) *implementation.ComponentConfiguration {

	return &implementation.ComponentConfiguration{
		ID:          guid.String(),
		Name:        componentConfName,
		Description: narrative,
		ConfigurableValues: []implementation.ConfigurableValue{
			implementation.ConfigurableValue{
				ValueID: uuid.NewV4().String(),
				Value:   "0",
			},
		},
	}
}

// AppendControlInImplementation appends a control in the implementation
func AppendControlInImplementation(cd implementation.ComponentDefinition, guid uuid.UUID, c Catalog, control string) implementation.ComponentDefinition {
	for i := range cd.ControlImplementations {
		for j := range cd.ControlImplementations[i].ControlConfigurations {
			if cd.ControlImplementations[i].ControlConfigurations[j].ConfigurationIDRef == guid.String() {
				ctrl := c.GetControl(control)
				pControls := cd.ControlImplementations[i].ControlConfigurations[j].ProvisioningMechanisms[0].ProvisionedControls
				if existsInControls(ctrl, pControls) {
					continue
				}
				cd.ControlImplementations[i].ControlConfigurations[j].ProvisioningMechanisms[0].ProvisionedControls = append(
					cd.ControlImplementations[i].ControlConfigurations[j].ProvisioningMechanisms[0].ProvisionedControls,
					implementation.ControlId{ControlID: ctrl, CatalogIDRef: c.GetID(), ItemID: ""},
				)
			}
		}

	}
	return cd
}

// CompileImplementation compiles all checks from maps to implementation json
func CompileImplementation(cdList []cdMap, CSVS [][]string, cat Catalog, components []component) implementation.Implementation {

	x := implementation.Implementation{
		ComponentDefinitions: func() []implementation.ComponentDefinition {
			var cds []implementation.ComponentDefinition
			for _, cd := range cdList {
				compD := implementation.ComponentDefinition{
					ID: func() string {
						for _, c := range cd {
							for _, comp := range components {
								if c.ID == comp.id {
									return c.ID
								}
							}
						}
						return ""
					}(),

					ComponentConfigurations: func() []*implementation.ComponentConfiguration {
						var arr []*implementation.ComponentConfiguration
						for _, v := range cd {
							for _, x := range v.ComponentConfigurations {
								arr = append(arr, x)
							}
						}

						return arr
					}(),
					ControlImplementations: func() []*implementation.ControlImplementation {
						arr := []*implementation.ControlImplementation{
							&implementation.ControlImplementation{
								ControlIds:            []implementation.ControlId{},
								ControlConfigurations: []implementation.ControlConfiguration{},
							},
						}
						for i := 3; i < totalControlsInExcel; i++ {
							if CSVS[i][controlIndex] == "" {
								continue
							}
							c := strings.ToLower(CSVS[i][controlIndex])
							if cat.isSubControl(c) {
								arr[0].ControlIds = append(arr[0].ControlIds, implementation.ControlId{
									ControlID:    cat.GetControl(c),
									ItemID:       c,
									CatalogIDRef: cat.GetID(),
								})
								continue
							}
							arr[0].ControlIds = append(arr[0].ControlIds, implementation.ControlId{
								ControlID:    cat.GetControl(CSVS[i][controlIndex]),
								ItemID:       "",
								CatalogIDRef: cat.GetID(),
							})
						}

						for _, v := range cd {
							for _, x := range v.ControlImplementations {
								arr[0].ControlConfigurations = append(
									arr[0].ControlConfigurations,
									x.ControlConfigurations...,
								)
							}
						}
						return arr
					}(),
				}

				cds = append(cds, compD)
			}
			return cds
		}(),
	}
	i := fillImplementsProfile(&x, components, CSVS)
	return *i
}

//Catalog catalog interface to determine control id pattern
type Catalog interface {
	GetControl(p string) string
	isSubControl(s string) bool
	GetID() string
}

// NISTCatalog NIST80053 catalog
type NISTCatalog struct {
	ID string
}

// GetID returns the NIST catalogID
func (n *NISTCatalog) GetID() string {
	return n.ID
}

// GetControl GetControl
func (*NISTCatalog) GetControl(p string) string {

	p = strings.ToLower(p)
	x := strings.Split(p, ".")
	y := strings.Split(x[0], " ")
	z := strings.Split(y[0], "_")
	control := z[0]
	controlLen := len(control)
	isControl, _ := regexp.MatchString("([a-z][a-z]-[0-9]*)$", control)
	if !isControl {
		control = control[:controlLen-1]
	}
	return control
}

func (*NISTCatalog) isSubControl(s string) bool {
	substrings := []string{" ", "(", "."}
	for _, substr := range substrings {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}

// GenerateImplementationParameter GenerateImplementationParameter
func GenerateImplementationParameter(param profile.SetParam, guidance []string) implementation.Parameter {
	return implementation.Parameter{
		ParameterID: param.Id,
		PossibleValues: func() []string {
			values := []string{}
			for _, x := range param.Constraints {
				values = append(values, x.Value)
			}
			return values
		}(),
		Guidance:     guidance,
		ValueID:      "{{valueId}}",
		DefaultValue: "{{defaultValue}}",
	}
}

func existsInParams(pID string, p []implementation.Parameter) bool {
	for _, x := range p {
		if x.ParameterID == pID {
			return true
		}
	}
	return false

}

func existsInControls(cID string, controls []implementation.ControlId) bool {
	for _, x := range controls {
		if x.ControlID == cID {
			return true
		}
	}
	return false
}

func getComponentID(componentName string) string {
	re := regexp.MustCompile(componentIDRegex)
	return re.FindString(componentName)
}

func getGuidance(alterations []profile.Alter, paramID string) []string {
	subControlID := getSubControlIDFromParam(paramID)
	for _, alter := range alterations {
		if alter.ControlId == subControlID {
			for _, addition := range alter.Additions {
				for _, part := range addition.Parts {
					if part.Class == "guidance" {
						guidance := getGuidanceFromPart(part.Prose)
						return guidance
					}

				}
			}
		}
	}
	return []string{}
}

func getSubControlIDFromParam(paramID string) string {
	re := regexp.MustCompile(subControlIDRegex)
	return re.FindString(paramID)
}

func getGuidanceFromPart(part *catalog.Prose) []string {
	var guidance []string
	for _, p := range part.P {
		guidance = append(guidance, p.Raw)
	}
	return guidance
}
