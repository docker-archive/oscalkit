package impl

import (
	"regexp"
	"strings"

	"github.com/docker/oscalkit/types/oscal/implementation"
)

func detokenizeParameterString(paramStr string) (string, string) {
	tokens := strings.Split(paramStr, profileDelimiter)
	if len(tokens) < 2 {
		return "", ""
	}
	profileID := getProfileID(strings.TrimSpace(tokens[0]))
	return strings.TrimSpace(profileID), strings.TrimSpace(tokens[1])

}

func fillImplementsProfile(imp *implementation.Implementation, cmps []component, CSVS [][]string) *implementation.Implementation {
	for _, c := range cmps {
		if !c.hasParameterMapping {
			continue
		}
		for i := rowIndex; i < totalControlsInExcel; i++ {
			parameterID := strings.TrimSpace(CSVS[i][c.parameterIDIndex])
			parameterType := strings.TrimSpace(CSVS[i][c.parameterStringIndex])
			mappings := strings.Split(parameterType, delimiter)
			for _, mapping := range mappings {
				profileID, checkAndValue := detokenizeParameterString(mapping)
				if profileID == "" || checkAndValue == "" {
					continue
				}
				for i := range imp.ComponentDefinitions {
					mapImplementsProfile(&imp.ComponentDefinitions[i],
						parameterID,
						profileID,
						checkAndValue)
				}
			}
		}
	}
	return imp

}

func addParemeters(ctrlconf *implementation.ControlConfiguration, parameterID, parameterValue string) []implementation.Parameter {
	paramFound := false
	for k, param := range ctrlconf.Parameters {
		if param.ParameterID == parameterID {
			paramFound = true
			valueFound := false
			for _, possibleValue := range param.PossibleValues {
				if possibleValue == parameterValue {
					valueFound = true
					break
				}
			}
			if !valueFound {
				ctrlconf.Parameters[k].PossibleValues = append(
					ctrlconf.Parameters[k].PossibleValues,
					parameterValue,
				)
			}
		}
	}
	if !paramFound {
		ctrlconf.Parameters = append(ctrlconf.Parameters, implementation.Parameter{
			ParameterID:    parameterID,
			PossibleValues: []string{parameterValue},
		})
	}
	return ctrlconf.Parameters
}

func mapImplementsProfile(cd *implementation.ComponentDefinition, parameterID, profileID, checkAndValue string) {
	ip := findOrCreateImplementsProfile(cd, profileID)
	checkName, parameterValue := parseCheckAndValue(checkAndValue)
	componentConfigID, found := findConfigIDByName(cd.ComponentConfigurations, checkName)
	if !found {
		return
	}
	ctrlConf := findOrCreateControlConfig(ip, componentConfigID)
	parameters := addParemeters(ctrlConf, parameterID, parameterValue)
	for j, cc := range ip.ControlConfigurations {
		if cc.ConfigurationIDRef == ctrlConf.ConfigurationIDRef {
			ip.ControlConfigurations[j].Parameters = parameters
		}
	}
}

func findOrCreateControlConfig(ip *implementation.ImplementsProfile, configIDRef string) *implementation.ControlConfiguration {
	for _, cc := range ip.ControlConfigurations {
		if cc.ConfigurationIDRef == configIDRef {
			return &cc
		}
	}
	newControlConfig := implementation.ControlConfiguration{
		ConfigurationIDRef: configIDRef,
		Parameters:         []implementation.Parameter{},
	}
	ip.ControlConfigurations = append(ip.ControlConfigurations, newControlConfig)
	return &newControlConfig
}

func findOrCreateImplementsProfile(cd *implementation.ComponentDefinition, profileID string) *implementation.ImplementsProfile {
	for _, ip := range cd.ImplementsProfiles {
		if ip.ProfileID == profileID {
			return ip
		}
	}
	newProfile := &implementation.ImplementsProfile{
		ProfileID: profileID,
		ControlConfigurations: func() []implementation.ControlConfiguration {
			arr := []implementation.ControlConfiguration{}
			for _, cc := range cd.ComponentConfigurations {
				arr = append(arr, implementation.ControlConfiguration{
					ConfigurationIDRef: cc.ID,
				})
			}
			return arr
		}(),
	}
	cd.ImplementsProfiles = append(cd.ImplementsProfiles, newProfile)
	return newProfile
}

func getProfileID(s string) string {
	s = strings.TrimSpace(s)
	if v, ok := profileMap[s]; ok {
		return v
	}
	return s
}

//parseCheckAndValue  Something(<=2) will change to [Something, <=2]
func parseCheckAndValue(s string) (string, string) {
	reg := regexp.MustCompile(`(\w{1,})\((.{1,})\)`)
	matches := reg.FindAllStringSubmatch(s, -1)
	if len(matches) < 1 {
		return "", ""
	}
	if len(matches[0]) < 2 {
		return "", ""
	}
	return matches[0][1], matches[0][2]

}

func findConfigIDByName(ccfs []*implementation.ComponentConfiguration, name string) (string, bool) {
	for _, ccf := range ccfs {
		if ccf.Name == name {
			return ccf.ID, true
		}
	}
	return "", false
}
