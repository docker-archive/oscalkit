package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/docker/oscalkit/types/oscal"
	"github.com/docker/oscalkit/types/oscal/catalog"
	"github.com/docker/oscalkit/types/oscal/profile"
)

// ProtocolsMapping Method To Parse The generated .go file and save the
// mapping of ID, Class & Titles
func ProtocolsMapping(check []catalog.Catalog) map[string][]string {

	securityControls := make(map[string][]string)
	for catalogCount := 0; catalogCount < len(check); catalogCount++ {
		for groupsCount := 0; groupsCount < len(check[catalogCount].Groups); groupsCount++ {
			for controlsCount := 0; controlsCount < len(check[catalogCount].Groups[groupsCount].Controls); controlsCount++ {
				if _, ok := securityControls[check[catalogCount].Groups[groupsCount].Controls[controlsCount].Id]; ok {
				} else {
					securityControls[check[catalogCount].Groups[groupsCount].Controls[controlsCount].Id] = append(securityControls[check[catalogCount].Groups[groupsCount].Controls[controlsCount].Id], check[catalogCount].Groups[groupsCount].Controls[controlsCount].Class)
					securityControls[check[catalogCount].Groups[groupsCount].Controls[controlsCount].Id] = append(securityControls[check[catalogCount].Groups[groupsCount].Controls[controlsCount].Id], string(check[catalogCount].Groups[groupsCount].Controls[controlsCount].Title))
				}

				for controlPartCount := 0; controlPartCount < len(check[catalogCount].Groups[groupsCount].Controls[controlsCount].Parts); controlPartCount++ {
					if _, ok := securityControls[check[catalogCount].Groups[groupsCount].Controls[controlsCount].Id+"?"+check[catalogCount].Groups[groupsCount].Controls[controlsCount].Parts[controlPartCount].Id]; ok {
					} else {
						if check[catalogCount].Groups[groupsCount].Controls[controlsCount].Parts[controlPartCount].Id != "" {
							securityControls[check[catalogCount].Groups[groupsCount].Controls[controlsCount].Id+"?"+check[catalogCount].Groups[groupsCount].Controls[controlsCount].Parts[controlPartCount].Id] = append(securityControls[check[catalogCount].Groups[groupsCount].Controls[controlsCount].Id+"?"+check[catalogCount].Groups[groupsCount].Controls[controlsCount].Parts[controlPartCount].Id], check[catalogCount].Groups[groupsCount].Controls[controlsCount].Parts[controlPartCount].Class)
							securityControls[check[catalogCount].Groups[groupsCount].Controls[controlsCount].Id+"?"+check[catalogCount].Groups[groupsCount].Controls[controlsCount].Parts[controlPartCount].Id] = append(securityControls[check[catalogCount].Groups[groupsCount].Controls[controlsCount].Id+"?"+check[catalogCount].Groups[groupsCount].Controls[controlsCount].Parts[controlPartCount].Id], string(check[catalogCount].Groups[groupsCount].Controls[controlsCount].Parts[controlPartCount].Title))
						} else if check[catalogCount].Groups[groupsCount].Controls[controlsCount].Parts[controlPartCount].Id == "" && check[catalogCount].Groups[groupsCount].Controls[controlsCount].Parts[controlPartCount].Class == "assessment" {
							securityControls[check[catalogCount].Groups[groupsCount].Controls[controlsCount].Id+"?"+check[catalogCount].Groups[groupsCount].Controls[controlsCount].Parts[controlPartCount].Id] = append(securityControls[check[catalogCount].Groups[groupsCount].Controls[controlsCount].Id+"?"+check[catalogCount].Groups[groupsCount].Controls[controlsCount].Parts[controlPartCount].Id], check[catalogCount].Groups[groupsCount].Controls[controlsCount].Parts[controlPartCount].Class)
							securityControls[check[catalogCount].Groups[groupsCount].Controls[controlsCount].Id+"?"+check[catalogCount].Groups[groupsCount].Controls[controlsCount].Parts[controlPartCount].Id] = append(securityControls[check[catalogCount].Groups[groupsCount].Controls[controlsCount].Id+"?"+check[catalogCount].Groups[groupsCount].Controls[controlsCount].Parts[controlPartCount].Id], string(check[catalogCount].Groups[groupsCount].Controls[controlsCount].Parts[controlPartCount].Title))
						}
					}
				}

				for subControlsCount := 0; subControlsCount < len(check[catalogCount].Groups[groupsCount].Controls[controlsCount].Controls); subControlsCount++ {
					if _, ok := securityControls[check[catalogCount].Groups[groupsCount].Controls[controlsCount].Controls[subControlsCount].Id]; ok {
					} else {
						securityControls[check[catalogCount].Groups[groupsCount].Controls[controlsCount].Controls[subControlsCount].Id] = append(securityControls[check[catalogCount].Groups[groupsCount].Controls[controlsCount].Controls[subControlsCount].Id], check[catalogCount].Groups[groupsCount].Controls[controlsCount].Controls[subControlsCount].Class)
						securityControls[check[catalogCount].Groups[groupsCount].Controls[controlsCount].Controls[subControlsCount].Id] = append(securityControls[check[catalogCount].Groups[groupsCount].Controls[controlsCount].Controls[subControlsCount].Id], string(check[catalogCount].Groups[groupsCount].Controls[controlsCount].Controls[subControlsCount].Title))
					}
					for subControlsPartCount := 0; subControlsPartCount < len(check[catalogCount].Groups[groupsCount].Controls[controlsCount].Controls[subControlsCount].Parts); subControlsPartCount++ {
						if _, ok := securityControls[check[catalogCount].Groups[groupsCount].Controls[controlsCount].Controls[subControlsCount].Id+"?"+check[catalogCount].Groups[groupsCount].Controls[controlsCount].Controls[subControlsCount].Parts[subControlsPartCount].Id]; ok {
						} else {
							if check[catalogCount].Groups[groupsCount].Controls[controlsCount].Controls[subControlsCount].Parts[subControlsPartCount].Id != "" {
								securityControls[check[catalogCount].Groups[groupsCount].Controls[controlsCount].Controls[subControlsCount].Id+"?"+check[catalogCount].Groups[groupsCount].Controls[controlsCount].Controls[subControlsCount].Parts[subControlsPartCount].Id] = append(securityControls[check[catalogCount].Groups[groupsCount].Controls[controlsCount].Controls[subControlsCount].Id+"?"+check[catalogCount].Groups[groupsCount].Controls[controlsCount].Controls[subControlsCount].Parts[subControlsPartCount].Id], check[catalogCount].Groups[groupsCount].Controls[controlsCount].Controls[subControlsCount].Parts[subControlsPartCount].Class)
								securityControls[check[catalogCount].Groups[groupsCount].Controls[controlsCount].Controls[subControlsCount].Id+"?"+check[catalogCount].Groups[groupsCount].Controls[controlsCount].Controls[subControlsCount].Parts[subControlsPartCount].Id] = append(securityControls[check[catalogCount].Groups[groupsCount].Controls[controlsCount].Controls[subControlsCount].Id+"?"+check[catalogCount].Groups[groupsCount].Controls[controlsCount].Controls[subControlsCount].Parts[subControlsPartCount].Id], string(check[catalogCount].Groups[groupsCount].Controls[controlsCount].Controls[subControlsCount].Parts[subControlsPartCount].Title))
							} else if check[catalogCount].Groups[groupsCount].Controls[controlsCount].Controls[subControlsCount].Parts[subControlsPartCount].Id == "" && check[catalogCount].Groups[groupsCount].Controls[controlsCount].Controls[subControlsCount].Parts[subControlsPartCount].Class == "assessment" {
								securityControls[check[catalogCount].Groups[groupsCount].Controls[controlsCount].Controls[subControlsCount].Id+"?"+check[catalogCount].Groups[groupsCount].Controls[controlsCount].Controls[subControlsCount].Parts[subControlsPartCount].Id] = append(securityControls[check[catalogCount].Groups[groupsCount].Controls[controlsCount].Controls[subControlsCount].Id+"?"+check[catalogCount].Groups[groupsCount].Controls[controlsCount].Controls[subControlsCount].Parts[subControlsPartCount].Id], check[catalogCount].Groups[groupsCount].Controls[controlsCount].Controls[subControlsCount].Parts[subControlsPartCount].Class)
								securityControls[check[catalogCount].Groups[groupsCount].Controls[controlsCount].Controls[subControlsCount].Id+"?"+check[catalogCount].Groups[groupsCount].Controls[controlsCount].Controls[subControlsCount].Parts[subControlsPartCount].Id] = append(securityControls[check[catalogCount].Groups[groupsCount].Controls[controlsCount].Controls[subControlsCount].Id+"?"+check[catalogCount].Groups[groupsCount].Controls[controlsCount].Controls[subControlsCount].Parts[subControlsPartCount].Id], string(check[catalogCount].Groups[groupsCount].Controls[controlsCount].Controls[subControlsCount].Parts[subControlsPartCount].Title))
							}
						}
					}
				}
			}
		}
	}
	return securityControls
}

// GetCatalog gets a catalog
func GetCatalog(r io.Reader) (*catalog.Catalog, error) {
	o, err := oscal.New(r)
	if err != nil {
		return nil, err
	}
	if o.Catalog == nil {
		return nil, fmt.Errorf("cannot map profile")
	}
	return o.Catalog, nil
}

// GetProfile gets a profile
func GetProfile(r io.Reader) (*profile.Profile, error) {
	o, err := oscal.New(r)
	if err != nil {
		return nil, err
	}
	if o.Profile == nil {
		return nil, fmt.Errorf("cannot map profile")
	}
	return o.Profile, nil
}

// controlInProfile accepts a Control or ControlID and an array of all
// the controls & subcontrols present in the profile.
func controlInProfile(controlID string, profile []string) bool {
	for _, value := range profile {
		if value == controlID {
			return true
		}
	}
	return false
}

// ParentControlCheck checks if the subcontrol's parent controls exists
// in the provided array on parent controls
func ParentControlCheck(subcontrol string, parentcontrols []string) bool {

	subControlTrim := strings.Split(subcontrol, ".")

	for _, value := range parentcontrols {
		if value == subControlTrim[0] {
			return true
		}
	}
	return false
}

// DownloadCatalog writes the JSON of the provided URL into a catalog.json file
func DownloadCatalog(url string) (string, error) {
	urlSplit := strings.Split(url, "/")
	tmpDir, err := ioutil.TempDir(".", "oscaltesttmp")
	if err != nil {
		log.Fatal(err)
	}
	fileName := tmpDir + "/" + urlSplit[len(urlSplit)-1]
	println("Catalog will be downloaded to: " + fileName)
	catalog, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer catalog.Close()
	println("Downloading catalog from URL: " + url)
	data, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer data.Body.Close()
	_, err = io.Copy(catalog, data.Body)
	if err != nil {
		return "", err
	}
	return tmpDir, nil
}

// ProfileParsing method to parse the profile and return the controls and subcontrols ID's
func ProfileParsing(parsedProfile *profile.Profile) []string {

	securityControls := make([]string, 0)

	for importCount := 0; importCount < len(parsedProfile.Imports); importCount++ {
		for idSelectorCount := 0; idSelectorCount < len(parsedProfile.Imports[importCount].Include.IdSelectors); idSelectorCount++ {
			if parsedProfile.Imports[importCount].Include.IdSelectors[idSelectorCount].ControlId != "" {
				securityControls = append(securityControls, parsedProfile.Imports[importCount].Include.IdSelectors[idSelectorCount].ControlId)
			}
			if parsedProfile.Imports[importCount].Include.IdSelectors[idSelectorCount].ControlId != "" {
				securityControls = append(securityControls, parsedProfile.Imports[importCount].Include.IdSelectors[idSelectorCount].ControlId)
			}
		}
	}
	return securityControls
}

// ParentControls to get the list of all parent controls in the profile
func ParentControls(parsedProfile *profile.Profile) []string {
	parentControlsList := make([]string, 0)

	for importCount := 0; importCount < len(parsedProfile.Imports); importCount++ {
		temp := ParseImport(parsedProfile, parsedProfile.Imports[importCount].Href.Path, "Parent")
		parentControlsList = appendslice(parentControlsList, temp)
	}

	parentControlsList = unique(parentControlsList)

	return parentControlsList
}

// ProfileProcessing is used to generate the mapping of ID Class & Title of
// all the controls subcontrols and parts
func ProfileProcessing(parsedProfile *profile.Profile, ListParentControls []string) map[string][]string {
	securityControlsDetails := make(map[string][]string)

	for importCounts := 0; importCounts < len(parsedProfile.Imports); importCounts++ {
		println("Import:", parsedProfile.Imports[importCounts].Href.String())
		dirName := "test_util/artifacts/"
		var err error
		if strings.Contains(parsedProfile.Imports[importCounts].Href.String(), "http") {
			dirName, err = DownloadCatalog(parsedProfile.Imports[importCounts].Href.String())
			if err != nil {
				log.Fatal(err)
			}
		}
		urlSplit := strings.Split(parsedProfile.Imports[importCounts].Href.Path, "/")
		fileName := dirName + "/" + urlSplit[len(urlSplit)-1]
		f, err := os.Open(fileName)
		if err != nil {
			log.Fatal(err)
		}
		check, _ := ProfileCatalogCheck(f)
		if check == "Catalog" {

			profileControls := ParseImport(parsedProfile, parsedProfile.Imports[importCounts].Href.Path, "all")

			catalogPath := dirName + "/" + urlSplit[len(urlSplit)-1]
			f, err := os.Open(catalogPath)
			if err != nil {
				log.Fatal(err)
			}

			parsedCatalog, err := GetCatalog(f)
			if err != nil {
				log.Fatal(err)
			}

			catalogControlsDetails := ParseCatalog(parsedCatalog, profileControls, ListParentControls)

			partsProfileControls := ProfileParsing(parsedProfile)

			parts := ParseParts(parsedProfile, partsProfileControls)

			catalogControlsDetails = appendAlterations(catalogControlsDetails, parts)

			println("Size of Catalog: ", len(catalogControlsDetails))
			if len(securityControlsDetails) == 0 {
				securityControlsDetails = appendMaps(securityControlsDetails, catalogControlsDetails)
			} else if len(securityControlsDetails) > 0 {
				securityControlsDetails = appendMaps(securityControlsDetails, catalogControlsDetails)
				securityControlsDetails = uniqueMaps(securityControlsDetails, catalogControlsDetails)
			}
			println("Size of securityControls: ", len(securityControlsDetails))

		} else if check == "Profile" {

			fmt.Println("profile path: " + urlSplit[len(urlSplit)-1])
			f, err := os.Open(dirName + "/" + urlSplit[len(urlSplit)-1])
			if err != nil {
				log.Fatal(err)
			}

			ProfileHref, err := GetProfile(f)
			if err != nil {
				log.Fatal(err)
			}

			ParsedProfile := ProfileProcessing(ProfileHref, ListParentControls)
			ParsedProfileControls := ParseImport(parsedProfile, parsedProfile.Imports[importCounts].Href.Path, "all")

			partsProfileControls := ProfileParsing(parsedProfile)

			parts := ParseParts(parsedProfile, partsProfileControls)

			println("Recursive count = ", len(ParsedProfile))
			println("Count of profile = ", len(ParsedProfileControls))

			println("Common = ", len(CommonMap(ParsedProfileControls, ParsedProfile)))
			securityControlsDetails = appendMaps(securityControlsDetails, CommonMap(ParsedProfileControls, ParsedProfile))

			securityControlsDetails = appendAlterations(securityControlsDetails, parts)

			println("Final Count = ", len(securityControlsDetails))
		}
	}

	return securityControlsDetails
}

// ParseCatalog accepts a catalog struct and return the mapping of Control,
// Controls & Parts. ID, Class & Titles
func ParseCatalog(parsedCatalog *catalog.Catalog, profileControls []string, ListParentControls []string) map[string][]string {
	catalogControlsDetails := make(map[string][]string)

	for groupCount := 0; groupCount < len(parsedCatalog.Groups); groupCount++ {
		for controlCount := 0; controlCount < len(parsedCatalog.Groups[groupCount].Controls); controlCount++ {
			if controlInProfile(parsedCatalog.Groups[groupCount].Controls[controlCount].Id, profileControls) {
				catalogControlsDetails[parsedCatalog.Groups[groupCount].Controls[controlCount].Id] = append(catalogControlsDetails[parsedCatalog.Groups[groupCount].Controls[controlCount].Id], parsedCatalog.Groups[groupCount].Controls[controlCount].Class)
				catalogControlsDetails[parsedCatalog.Groups[groupCount].Controls[controlCount].Id] = append(catalogControlsDetails[parsedCatalog.Groups[groupCount].Controls[controlCount].Id], string(parsedCatalog.Groups[groupCount].Controls[controlCount].Title))
				for controlPartCount := 0; controlPartCount < len(parsedCatalog.Groups[groupCount].Controls[controlCount].Parts); controlPartCount++ {
					if parsedCatalog.Groups[groupCount].Controls[controlCount].Parts[controlPartCount].Id != "" {
						catalogControlsDetails[parsedCatalog.Groups[groupCount].Controls[controlCount].Id+"?"+parsedCatalog.Groups[groupCount].Controls[controlCount].Parts[controlPartCount].Id] = append(catalogControlsDetails[parsedCatalog.Groups[groupCount].Controls[controlCount].Id+"?"+parsedCatalog.Groups[groupCount].Controls[controlCount].Parts[controlPartCount].Id], parsedCatalog.Groups[groupCount].Controls[controlCount].Parts[controlPartCount].Class)
						catalogControlsDetails[parsedCatalog.Groups[groupCount].Controls[controlCount].Id+"?"+parsedCatalog.Groups[groupCount].Controls[controlCount].Parts[controlPartCount].Id] = append(catalogControlsDetails[parsedCatalog.Groups[groupCount].Controls[controlCount].Id+"?"+parsedCatalog.Groups[groupCount].Controls[controlCount].Parts[controlPartCount].Id], string(parsedCatalog.Groups[groupCount].Controls[controlCount].Parts[controlPartCount].Title))
					} else if parsedCatalog.Groups[groupCount].Controls[controlCount].Parts[controlPartCount].Id == "" && parsedCatalog.Groups[groupCount].Controls[controlCount].Parts[controlPartCount].Class == "assessment" {
						catalogControlsDetails[parsedCatalog.Groups[groupCount].Controls[controlCount].Id+"?"+parsedCatalog.Groups[groupCount].Controls[controlCount].Parts[controlPartCount].Id] = append(catalogControlsDetails[parsedCatalog.Groups[groupCount].Controls[controlCount].Id+"?"+parsedCatalog.Groups[groupCount].Controls[controlCount].Parts[controlPartCount].Id], parsedCatalog.Groups[groupCount].Controls[controlCount].Parts[controlPartCount].Class)
						catalogControlsDetails[parsedCatalog.Groups[groupCount].Controls[controlCount].Id+"?"+parsedCatalog.Groups[groupCount].Controls[controlCount].Parts[controlPartCount].Id] = append(catalogControlsDetails[parsedCatalog.Groups[groupCount].Controls[controlCount].Id+"?"+parsedCatalog.Groups[groupCount].Controls[controlCount].Parts[controlPartCount].Id], string(parsedCatalog.Groups[groupCount].Controls[controlCount].Parts[controlPartCount].Title))
					}
				}
			}

			for subControlCount := 0; subControlCount < len(parsedCatalog.Groups[groupCount].Controls[controlCount].Controls); subControlCount++ {
				if controlInProfile(parsedCatalog.Groups[groupCount].Controls[controlCount].Controls[subControlCount].Id, profileControls) && ParentControlCheck(parsedCatalog.Groups[groupCount].Controls[controlCount].Controls[subControlCount].Id, ListParentControls) {
					catalogControlsDetails[parsedCatalog.Groups[groupCount].Controls[controlCount].Controls[subControlCount].Id] = append(catalogControlsDetails[parsedCatalog.Groups[groupCount].Controls[controlCount].Controls[subControlCount].Id], parsedCatalog.Groups[groupCount].Controls[controlCount].Controls[subControlCount].Class)
					catalogControlsDetails[parsedCatalog.Groups[groupCount].Controls[controlCount].Controls[subControlCount].Id] = append(catalogControlsDetails[parsedCatalog.Groups[groupCount].Controls[controlCount].Controls[subControlCount].Id], string(parsedCatalog.Groups[groupCount].Controls[controlCount].Controls[subControlCount].Title))
					for subControlPartCount := 0; subControlPartCount < len(parsedCatalog.Groups[groupCount].Controls[controlCount].Controls[subControlCount].Parts); subControlPartCount++ {
						if parsedCatalog.Groups[groupCount].Controls[controlCount].Controls[subControlCount].Parts[subControlPartCount].Id != "" {
							catalogControlsDetails[parsedCatalog.Groups[groupCount].Controls[controlCount].Controls[subControlCount].Id+"?"+parsedCatalog.Groups[groupCount].Controls[controlCount].Controls[subControlCount].Parts[subControlPartCount].Id] = append(catalogControlsDetails[parsedCatalog.Groups[groupCount].Controls[controlCount].Controls[subControlCount].Id+"?"+parsedCatalog.Groups[groupCount].Controls[controlCount].Controls[subControlCount].Parts[subControlPartCount].Id], parsedCatalog.Groups[groupCount].Controls[controlCount].Controls[subControlCount].Parts[subControlPartCount].Class)
							catalogControlsDetails[parsedCatalog.Groups[groupCount].Controls[controlCount].Controls[subControlCount].Id+"?"+parsedCatalog.Groups[groupCount].Controls[controlCount].Controls[subControlCount].Parts[subControlPartCount].Id] = append(catalogControlsDetails[parsedCatalog.Groups[groupCount].Controls[controlCount].Controls[subControlCount].Id+"?"+parsedCatalog.Groups[groupCount].Controls[controlCount].Controls[subControlCount].Parts[subControlPartCount].Id], string(parsedCatalog.Groups[groupCount].Controls[controlCount].Controls[subControlCount].Parts[subControlPartCount].Title))
						} else if parsedCatalog.Groups[groupCount].Controls[controlCount].Controls[subControlCount].Parts[subControlPartCount].Id == "" && parsedCatalog.Groups[groupCount].Controls[controlCount].Controls[subControlCount].Parts[subControlPartCount].Class == "assessment" {

							catalogControlsDetails[parsedCatalog.Groups[groupCount].Controls[controlCount].Controls[subControlCount].Id+"?"+parsedCatalog.Groups[groupCount].Controls[controlCount].Controls[subControlCount].Parts[subControlPartCount].Id] = append(catalogControlsDetails[parsedCatalog.Groups[groupCount].Controls[controlCount].Controls[subControlCount].Id+"?"+parsedCatalog.Groups[groupCount].Controls[controlCount].Controls[subControlCount].Parts[subControlPartCount].Id], parsedCatalog.Groups[groupCount].Controls[controlCount].Controls[subControlCount].Parts[subControlPartCount].Class)
							catalogControlsDetails[parsedCatalog.Groups[groupCount].Controls[controlCount].Controls[subControlCount].Id+"?"+parsedCatalog.Groups[groupCount].Controls[controlCount].Controls[subControlCount].Parts[subControlPartCount].Id] = append(catalogControlsDetails[parsedCatalog.Groups[groupCount].Controls[controlCount].Controls[subControlCount].Id+"?"+parsedCatalog.Groups[groupCount].Controls[controlCount].Controls[subControlCount].Parts[subControlPartCount].Id], string(parsedCatalog.Groups[groupCount].Controls[controlCount].Controls[subControlCount].Parts[subControlPartCount].Title))
						}
					}
				}
			}
		}
	}
	return catalogControlsDetails
}

// ProfileCatalogCheck checks if the path provided is for a profile or a catolog
func ProfileCatalogCheck(r io.Reader) (string, error) {
	o, err := oscal.New(r)
	if err != nil {
		return "Invalid File", err
	}
	if o.Profile == nil {
		return "Catalog", nil
	}
	if o.Catalog == nil {
		return "Profile", nil
	}
	return "Invalid File", nil
}

// CommonMap returns the elements in Map that are also present in profile
func CommonMap(profile []string, catalogControlsDetails map[string][]string) map[string][]string {

	commonMapping := make(map[string][]string)

	for key, mapValue := range catalogControlsDetails {
		for _, sliceValue := range profile {
			subControlTrim := strings.Split(key, "?")

			if sliceValue == key {
				commonMapping[key] = append(commonMapping[key], mapValue[0])
				commonMapping[key] = append(commonMapping[key], mapValue[1])
			} else if sliceValue == subControlTrim[0] {
				commonMapping[key] = append(commonMapping[key], mapValue[0])
				commonMapping[key] = append(commonMapping[key], mapValue[1])
			}
		}
	}
	return commonMapping
}

// ParseImport method to parse the profile and return the controls and subcontrols or only controls
func ParseImport(parsedProfile *profile.Profile, link string, token string) []string {

	securityControls := make([]string, 0)
	for importCount := 0; importCount < len(parsedProfile.Imports); importCount++ {
		if parsedProfile.Imports[importCount].Href.Path == link {
			for idSelectorCount := 0; idSelectorCount < len(parsedProfile.Imports[importCount].Include.IdSelectors); idSelectorCount++ {
				if parsedProfile.Imports[importCount].Include.IdSelectors[idSelectorCount].ControlId != "" {
					securityControls = append(securityControls, parsedProfile.Imports[importCount].Include.IdSelectors[idSelectorCount].ControlId)
				}
				if parsedProfile.Imports[importCount].Include.IdSelectors[idSelectorCount].ControlId != "" && token != "Parent" {
					securityControls = append(securityControls, parsedProfile.Imports[importCount].Include.IdSelectors[idSelectorCount].ControlId)
				}
			}
		}
	}

	return securityControls
}

// ParseParts method to parse the profile and return the mapping of all the parts
func ParseParts(parsedProfile *profile.Profile, list []string) map[string][]string {

	securityControls := make(map[string][]string)

	for modifyCount := 0; modifyCount < len(parsedProfile.Modify.Alterations); modifyCount++ {
		for alterCount := 0; alterCount < len(parsedProfile.Modify.Alterations[modifyCount].Additions); alterCount++ {
			for partCount := 0; partCount < len(parsedProfile.Modify.Alterations[modifyCount].Additions[alterCount].Parts); partCount++ {
				for _, s1Element := range list {
					if parsedProfile.Modify.Alterations[modifyCount].ControlId == s1Element {
						if parsedProfile.Modify.Alterations[modifyCount].ControlId != "" && parsedProfile.Modify.Alterations[modifyCount].Additions[alterCount].Parts[partCount].Class == "guidance" {
							securityControls[parsedProfile.Modify.Alterations[modifyCount].ControlId+"?"+parsedProfile.Modify.Alterations[modifyCount].ControlId+"_gdn"] = append(securityControls[parsedProfile.Modify.Alterations[modifyCount].ControlId+"?"+parsedProfile.Modify.Alterations[modifyCount].ControlId+"_gdn"], parsedProfile.Modify.Alterations[modifyCount].Additions[alterCount].Parts[partCount].Class)
						}
					} else if parsedProfile.Modify.Alterations[modifyCount].ControlId == s1Element {
						if parsedProfile.Modify.Alterations[modifyCount].ControlId != "" && parsedProfile.Modify.Alterations[modifyCount].Additions[alterCount].Parts[partCount].Class == "guidance" {
							securityControls[parsedProfile.Modify.Alterations[modifyCount].ControlId+"?"+parsedProfile.Modify.Alterations[modifyCount].ControlId+"_gdn"] = append(securityControls[parsedProfile.Modify.Alterations[modifyCount].ControlId+"?"+parsedProfile.Modify.Alterations[modifyCount].ControlId+"_gdn"], parsedProfile.Modify.Alterations[modifyCount].Additions[alterCount].Parts[partCount].Class)
						}
					}
				}
			}
		}
	}

	return securityControls
}

// appendslice appends two slices
func appendslice(slice []string, slice1 []string) []string {

	for sliceCount := 0; sliceCount < len(slice1); sliceCount++ {
		slice = append(slice, slice1[sliceCount])
	}

	return slice
}

// AreMapsSame compares the values of two  same length maps and returns true if both the maps have the same key value pairs
func AreMapsSame(profileControlsDetails map[string][]string, codeGeneratedMapping map[string][]string, token string) bool {
	for key := range profileControlsDetails {
		if !strings.Contains(key, "?") && token == "controls" {
			if profileControlsDetails[key][0] != codeGeneratedMapping[key][0] && profileControlsDetails[key][1] != codeGeneratedMapping[key][1] {
				println("Mapping for " + key + " incorrect.")
				return false
			}
		} else if strings.Contains(key, "?") && token == "parts" {
			if profileControlsDetails[key][0] != codeGeneratedMapping[key][0] && profileControlsDetails[key][1] != codeGeneratedMapping[key][1] {
				println("Mapping for " + key + " incorrect.")
				return false
			}
		}
	}
	return true
}

// unique returns unique values in the slice
func unique(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// appendMaps appends two maps
func appendMaps(securityControlsDetails map[string][]string, catalogControlsDetails map[string][]string) map[string][]string {

	for key, value := range catalogControlsDetails {
		securityControlsDetails[key] = value
	}

	return securityControlsDetails
}

func appendAlterations(securityControlsDetails map[string][]string, PartsDetails map[string][]string) map[string][]string {

	for key, value := range PartsDetails {
		if _, ok := securityControlsDetails[key]; ok {
			delete(securityControlsDetails, key)
			securityControlsDetails[key+"_1"] = value
			securityControlsDetails[key+"_2"] = value
		}
	}

	return securityControlsDetails
}

func uniqueMaps(securityControlsDetails map[string][]string, catalogControlsDetails map[string][]string) map[string][]string {

	for key, value := range catalogControlsDetails {
		if _, ok := securityControlsDetails[key]; !ok {
			securityControlsDetails[key] = value
		}
	}

	return securityControlsDetails
}

// Count to take count of either parts of controls & subcontrols
func Count(securityControlsDetails map[string][]string, token string) int {

	count := 0

	for key := range securityControlsDetails {
		if token == "parts" {
			count++
		} else if token == "controls" {
			if !strings.Contains(key, "?") {
				count++
			}
		}
	}

	return count
}
