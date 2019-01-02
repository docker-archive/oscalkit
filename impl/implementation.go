package impl

import (
	"regexp"
	"strings"

	"github.com/opencontrol/oscalkit/types/oscal/implementation"
	"github.com/opencontrol/oscalkit/types/oscal/profile"
	uuid "github.com/satori/go.uuid"
)

const (
	//TotalControlsInExcel the total number of controls in the excel sheet
	TotalControlsInExcel = 264
	//ComponentNameIndex The Column at which name of the component configuration is present
	ComponentNameIndex = 15
	//NarrativeIndex The Column at which narrative of the component configuration is present
	NarrativeIndex = 16
	//ControlIndex Column at which control is present in the excel sheet
	ControlIndex = 2
	//RowIndex Starting point for valid rows (neglects titles)
	RowIndex = 3
)

type guidMap map[string]uuid.UUID
type cdMap map[string]implementation.ComponentDefinition

//GenerateImplementation generates implementation from component excel sheet
func GenerateImplementation(CSVS [][]string, p *profile.Profile, c Catalog) implementation.Implementation {

	ComponentDefinitonMap := make(map[string]implementation.ComponentDefinition)
	checkAgainstGUUID := make(map[string]uuid.UUID)

	for i := RowIndex; i < TotalControlsInExcel; i++ {
		applicableControl := CSVS[i][ControlIndex]
		if applicableControl == "" {
			continue
		}
		applicableNarrative := CSVS[i][NarrativeIndex]
		ListOfComponentConfigName := strings.Split(CSVS[i][ComponentNameIndex], "|")
		for _, componentConfigName := range ListOfComponentConfigName {
			componentConfigName = strings.TrimSpace(componentConfigName)
			if componentConfigName == "" {
				continue
			}
			if _, ok := ComponentDefinitonMap[componentConfigName]; !ok {
				CreateComponentDefinition(checkAgainstGUUID, ComponentDefinitonMap, componentConfigName, p, c, applicableControl, applicableNarrative)
			} else {
				securityCheck := ComponentDefinitonMap[componentConfigName]
				guid := checkAgainstGUUID[componentConfigName]
				ComponentDefinitonMap[componentConfigName] = AppendParameterInImplementation(securityCheck, guid, p, c, applicableControl)
			}
		}
	}
	return CompileImplemenatation(ComponentDefinitonMap, CSVS, c, p)

}

//CreateComponentDefinition creates a component definition
func CreateComponentDefinition(gm guidMap, cdm cdMap, componentConfName string, p *profile.Profile, c Catalog, control, narrative string) {

	componentConfGUID := uuid.NewV4()
	gm[componentConfName] = componentConfGUID
	controlConfiguration := implementation.ControlConfiguration{
		ConfigurationIDRef: componentConfGUID.String(),
	}
	var parameters []implementation.Parameter
	if p.Modify != nil {
		for _, param := range p.Modify.ParamSettings {
			if param.Id == "" {
				continue
			}
			if c.GetControl(param.Id) == c.GetControl(control) {
				if existsInParams(param.Id, parameters) {
					continue
				}
				x := GenerateImplementationParamter(param)
				parameters = append(parameters, x)
			}
		}
	}

	controlConfiguration.Parameters = parameters
	cdm[componentConfName] = implementation.ComponentDefinition{
		ComponentConfigurations: []*implementation.ComponentConfiguration{
			CreateComponentConfiguration(componentConfGUID, componentConfName, narrative),
		},
		ImplementsProfiles: []*implementation.ImplementsProfile{
			&implementation.ImplementsProfile{
				ProfileID: p.ID,
				ControlConfigurations: []implementation.ControlConfiguration{
					controlConfiguration,
				},
			},
		},
	}

}

//CreateComponentConfiguration creates component configuration
func CreateComponentConfiguration(guid uuid.UUID, componentConfName, narrative string) *implementation.ComponentConfiguration {

	return &implementation.ComponentConfiguration{
		ID:          guid.String(),
		Name:        componentConfName,
		Description: narrative,
		ConfigurableValues: []implementation.ConfigurableValue{
			implementation.ConfigurableValue{
				ValueID: uuid.NewV4().String(),
				Value:   0,
			},
		},
	}
}

//AppendParameterInImplementation Appends parameter in the relative guid
func AppendParameterInImplementation(cd implementation.ComponentDefinition, guid uuid.UUID, p *profile.Profile, c Catalog, control string) implementation.ComponentDefinition {
	for i := range cd.ImplementsProfiles {
		for j := range cd.ImplementsProfiles[i].ControlConfigurations {
			if guid.String() == cd.ImplementsProfiles[i].ControlConfigurations[j].ConfigurationIDRef {
				for _, param := range p.Modify.ParamSettings {
					if param.Id == "" {
						continue
					}
					if existsInParams(param.Id, cd.ImplementsProfiles[i].ControlConfigurations[j].Parameters) {
						continue
					}
					if c.GetControl(param.Id) == c.GetControl(control) {
						x := GenerateImplementationParamter(param)
						cd.ImplementsProfiles[i].ControlConfigurations[j].Parameters = append(cd.ImplementsProfiles[i].ControlConfigurations[j].Parameters, x)
					}
				}
			}
		}
	}
	return cd

}

//CompileImplemenatation compiles all checks from maps to implementation json
func CompileImplemenatation(cd cdMap, CSVS [][]string, cat Catalog, p *profile.Profile) implementation.Implementation {
	return implementation.Implementation{
		ComponentDefinitions: []implementation.ComponentDefinition{
			implementation.ComponentDefinition{
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
							ControlIds: []implementation.ControlId{},
						},
					}
					for i := 3; i < TotalControlsInExcel; i++ {
						if CSVS[i][ControlIndex] == "" {
							continue
						}
						c := strings.ToLower(CSVS[i][ControlIndex])
						if cat.isSubControl(c) {
							arr[0].ControlIds = append(arr[0].ControlIds, implementation.ControlId{
								ControlID:    "",
								ItemID:       c,
								CatalogIDRef: cat.GetID(),
							})
							continue
						}
						arr[0].ControlIds = append(arr[0].ControlIds, implementation.ControlId{
							ControlID:    cat.GetControl(CSVS[i][ControlIndex]),
							ItemID:       "",
							CatalogIDRef: cat.GetID(),
						})
					}
					for _, def := range cd {
						for _, ci := range def.ImplementsProfiles {
							for _, cc := range ci.ControlConfigurations {
								arr[0].ControlConfigurations = append(arr[0].ControlConfigurations, cc)
							}
						}
					}
					return arr
				}(),
				ImplementsProfiles: []*implementation.ImplementsProfile{
					&implementation.ImplementsProfile{
						ProfileID: p.ID,
						ControlConfigurations: func() []implementation.ControlConfiguration {
							var arr []implementation.ControlConfiguration
							for _, v := range cd {
								for _, x := range v.ImplementsProfiles {
									for _, y := range x.ControlConfigurations {
										arr = append(arr, y)
									}
								}
							}
							return arr
						}(),
					},
				},
			},
		},
	}
}

//Catalog catalog interface to determine control id pattern
type Catalog interface {
	GetControl(p string) string
	isSubControl(s string) bool
	GetID() string
}

//NISTCatalog NIST80053 catalog
type NISTCatalog struct {
	ID string
}

func (n *NISTCatalog) GetID() string {
	return n.ID
}

//GetControl GetControl
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

//GenerateImplementationParamter GenerateImplementationParamter
func GenerateImplementationParamter(param profile.SetParam) implementation.Parameter {
	return implementation.Parameter{
		ParameterID: param.Id,
		Value: func() string {
			str := ""
			for _, x := range param.Constraints {
				str += x.Value + " "
			}
			return str
		}(),
		Guidance: "{{guidance}}",
		ValueID:  "{{valueId}}",
	}
}

func existsInParams(pId string, p []implementation.Parameter) bool {
	for _, x := range p {
		if x.ParameterID == pId {
			return true
		}
	}
	return false

}
