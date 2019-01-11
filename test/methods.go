package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/docker/oscalkit/types/oscal"
	"github.com/docker/oscalkit/types/oscal/catalog"
	"github.com/docker/oscalkit/types/oscal/profile"
)

// StructExaminer To Verify The Structure
func StructExaminer(t reflect.Type, depth int) {
	fmt.Println("\nType is", t.Name(), "and kind is", t.Kind())
	switch t.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Ptr, reflect.Slice:
		fmt.Println("Contained type:")
		StructExaminer(t.Elem(), depth+1)
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			fmt.Print(f.Name+" "+f.Type.Name(), f.Type.Kind())
			if f.Tag != "" {
				fmt.Println(" " + f.Tag)
			}
		}
	}
}

// ProtocolsMapping Method To Parse All The Controls From Catalog.go
func ProtocolsMapping(check []catalog.Catalog) map[string][]string {

	SecurityControls := make(map[string][]string)

	for CatalogCount := 0; CatalogCount < len(check); CatalogCount++ {
		for GroupsCount := 0; GroupsCount < len(check[CatalogCount].Groups); GroupsCount++ {
			for ControlsCount := 0; ControlsCount < len(check[CatalogCount].Groups[GroupsCount].Controls); ControlsCount++ {
				SecurityControls[check[CatalogCount].Groups[GroupsCount].Controls[ControlsCount].Id] = append(SecurityControls[check[CatalogCount].Groups[GroupsCount].Controls[ControlsCount].Id], check[CatalogCount].Groups[GroupsCount].Controls[ControlsCount].Class)
				SecurityControls[check[CatalogCount].Groups[GroupsCount].Controls[ControlsCount].Id] = append(SecurityControls[check[CatalogCount].Groups[GroupsCount].Controls[ControlsCount].Id], string(check[CatalogCount].Groups[GroupsCount].Controls[ControlsCount].Title))
				for SubControlsCount := 0; SubControlsCount < len(check[CatalogCount].Groups[GroupsCount].Controls[ControlsCount].Subcontrols); SubControlsCount++ {
					SecurityControls[check[CatalogCount].Groups[GroupsCount].Controls[ControlsCount].Subcontrols[SubControlsCount].Id] = append(SecurityControls[check[CatalogCount].Groups[GroupsCount].Controls[ControlsCount].Subcontrols[SubControlsCount].Id], check[CatalogCount].Groups[GroupsCount].Controls[ControlsCount].Subcontrols[SubControlsCount].Class)
					SecurityControls[check[CatalogCount].Groups[GroupsCount].Controls[ControlsCount].Subcontrols[SubControlsCount].Id] = append(SecurityControls[check[CatalogCount].Groups[GroupsCount].Controls[ControlsCount].Subcontrols[SubControlsCount].Id], string(check[CatalogCount].Groups[GroupsCount].Controls[ControlsCount].Subcontrols[SubControlsCount].Title))
				}
			}
		}
	}
	return SecurityControls
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

// controlInProfile checks if the control provided exists in the provided profile or not
func controlInProfile(controlID string, profile []string) bool {
	for _, value := range profile {
		if value == controlID {
			return true
		}
	}
	return false
}

// DownloadCatalog writes the JSON of the provided URL into a catalog.json file
func DownloadCatalog(url string) (string, error) {
	save := strings.Split(url, "/")
	tmpDir, err := ioutil.TempDir(".", "oscaltesttmp")
	if err != nil {
		log.Fatal(err)
	}
	filename := tmpDir + "/" + save[len(save)-1]
	println("Catalog will be downloaded to: " + filename)
	catalog, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer catalog.Close()
	if !strings.Contains(url, "http") {
		url = "https://raw.githubusercontent.com/usnistgov/OSCAL/master/content/nist.gov/SP800-53/rev4/" + save[len(save)-1]
	}
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

// ProfileParsing method to parse the profile and return the controls and subcontrols
func ProfileParsing(parsedProfile *profile.Profile) []string {

	SecurityControls := make([]string, 0)

	for i := 0; i < len(parsedProfile.Imports); i++ {
		for j := 0; j < len(parsedProfile.Imports[i].Include.IdSelectors); j++ {
			if parsedProfile.Imports[i].Include.IdSelectors[j].ControlId != "" {
				SecurityControls = append(SecurityControls, parsedProfile.Imports[i].Include.IdSelectors[j].ControlId)
			}
			if parsedProfile.Imports[i].Include.IdSelectors[j].SubcontrolId != "" {
				SecurityControls = append(SecurityControls, parsedProfile.Imports[i].Include.IdSelectors[j].SubcontrolId)
			}
		}
	}
	return SecurityControls
}

// ProfileProcessing is used to get to the catalog referenced in the profile and parse it into a map
func ProfileProcessing(parsedProfile *profile.Profile) map[string][]string {
	SecurityControlsDetails := make(map[string][]string)

	for l := 0; l < len(parsedProfile.Imports); l++ {
		println("Import:", parsedProfile.Imports[l].Href.Scheme+"://"+parsedProfile.Imports[l].Href.Host+parsedProfile.Imports[l].Href.Path)
		//println("Import:", parsedProfile.Imports[l].Href.Path)
		dirName, err := DownloadCatalog(parsedProfile.Imports[l].Href.Scheme + "://" + parsedProfile.Imports[l].Href.Host + parsedProfile.Imports[l].Href.Path)
		if err != nil {
			log.Fatal(err)
		}
		save := strings.Split(parsedProfile.Imports[l].Href.Path, "/")
		filename := dirName + "/" + save[len(save)-1]
		f, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		check, _ := ProfileCatalogCheck(f)
		if check == "Catalog" {

			ProfileControls := ImportParsing(parsedProfile, parsedProfile.Imports[l].Href.Path)

			catalog := dirName + "/" + save[len(save)-1]
			f, err := os.Open(catalog)
			if err != nil {
				log.Fatal(err)
			}

			parsedCatalog, err := GetCatalog(f)
			if err != nil {
				log.Fatal(err)
			}
			CatalogControlsDetails := make(map[string][]string)

			for i := 0; i < len(parsedCatalog.Groups); i++ {
				for j := 0; j < len(parsedCatalog.Groups[i].Controls); j++ {
					if controlInProfile(parsedCatalog.Groups[i].Controls[j].Id, ProfileControls) {
						CatalogControlsDetails[parsedCatalog.Groups[i].Controls[j].Id] = append(CatalogControlsDetails[parsedCatalog.Groups[i].Controls[j].Id], parsedCatalog.Groups[i].Controls[j].Class)
						CatalogControlsDetails[parsedCatalog.Groups[i].Controls[j].Id] = append(CatalogControlsDetails[parsedCatalog.Groups[i].Controls[j].Id], string(parsedCatalog.Groups[i].Controls[j].Title))
					}
					for k := 0; k < len(parsedCatalog.Groups[i].Controls[j].Subcontrols); k++ {
						if controlInProfile(parsedCatalog.Groups[i].Controls[j].Subcontrols[k].Id, ProfileControls) {
							CatalogControlsDetails[parsedCatalog.Groups[i].Controls[j].Subcontrols[k].Id] = append(CatalogControlsDetails[parsedCatalog.Groups[i].Controls[j].Subcontrols[k].Id], parsedCatalog.Groups[i].Controls[j].Subcontrols[k].Class)
							CatalogControlsDetails[parsedCatalog.Groups[i].Controls[j].Subcontrols[k].Id] = append(CatalogControlsDetails[parsedCatalog.Groups[i].Controls[j].Subcontrols[k].Id], string(parsedCatalog.Groups[i].Controls[j].Subcontrols[k].Title))
						}
					}
				}
			}
			println("Size of Catalog: ", len(CatalogControlsDetails))
			if len(SecurityControlsDetails) == 0 {
				SecurityControlsDetails = appendMaps(SecurityControlsDetails, CatalogControlsDetails)
			} else if len(SecurityControlsDetails) > 0 {
				SecurityControlsDetails = appendMaps(SecurityControlsDetails, CatalogControlsDetails)
				SecurityControlsDetails = uniqueMaps(SecurityControlsDetails, CatalogControlsDetails)
			}

			println("Size of SecurityControls: ", len(SecurityControlsDetails))
		} else if check == "Profile" {
			f, err := os.Open(save[len(save)-1])
			if err != nil {
				log.Fatal(err)
			}

			parsedProfile1, err := GetProfile(f)
			if err != nil {
				log.Fatal(err)
			}

			save := ProfileProcessing(parsedProfile1)
			save1 := ImportParsing(parsedProfile, parsedProfile.Imports[l].Href.Path)

			println(len(save))
			println(len(save1))

			SecurityControlsDetails = appendMaps(SecurityControlsDetails, CommonMap(save1, save))
			println(len(SecurityControlsDetails))
			// parsedProfile and parsedProfile1 common and save in SecurityControls
		}
	}

	return SecurityControlsDetails
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

// CommonMap returns the elements in Map that are also present in slice
func CommonMap(slice1 []string, CatalogControlsDetails map[string][]string) map[string][]string {
	Result := make(map[string][]string)
	for _, s1element := range slice1 {
		if _, ok := CatalogControlsDetails[s1element]; ok {
			save := CatalogControlsDetails[s1element]
			CatalogControlsDetails[s1element] = append(CatalogControlsDetails[s1element], save[0])
			CatalogControlsDetails[s1element] = append(CatalogControlsDetails[s1element], save[1])
		}
	}
	return Result
}

// RemoveDuplicateSlice returns the elements after removing duplicates
func RemoveDuplicateSlice(slice1 []string, slice2 []string) []string {
	result := make([]string, 0)
	count := 0
	for _, s2element := range slice2 {
		for _, s1element := range slice1 {
			if s2element != s1element {
				count++
			}
		}
		if count > 0 {
			result = append(result, s2element)
		}
		count = 0
	}
	return result
}

// ImportParsing method to parse the profile and return the controls and subcontrols
func ImportParsing(parsedProfile *profile.Profile, link string) []string {

	SecurityControls := make([]string, 0)

	for i := 0; i < len(parsedProfile.Imports); i++ {
		if parsedProfile.Imports[i].Href.Path == link {
			for j := 0; j < len(parsedProfile.Imports[i].Include.IdSelectors); j++ {
				if parsedProfile.Imports[i].Include.IdSelectors[j].ControlId != "" {
					SecurityControls = append(SecurityControls, parsedProfile.Imports[i].Include.IdSelectors[j].ControlId)
				}
				if parsedProfile.Imports[i].Include.IdSelectors[j].SubcontrolId != "" {
					SecurityControls = append(SecurityControls, parsedProfile.Imports[i].Include.IdSelectors[j].SubcontrolId)
				}
			}
		}
	}
	return SecurityControls
}

// AreMapsSame compares the values of two  same length maps and returns true if both the maps have the same key value pairs
func AreMapsSame(profileControlsDetails map[string][]string, codeGeneratedMapping map[string][]string) bool {
	for key := range profileControlsDetails {
		if !reflect.DeepEqual(profileControlsDetails[key], codeGeneratedMapping[key]) {
			println("Mapping for " + key + " incorrect.")
			return false
		}
	}
	return true
}

func unique(intSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func appendMaps(SecurityControlsDetails map[string][]string, CatalogControlsDetails map[string][]string) map[string][]string {

	for k, v := range CatalogControlsDetails {
		SecurityControlsDetails[k] = v
	}

	return SecurityControlsDetails
}

func uniqueMaps(SecurityControlsDetails map[string][]string, CatalogControlsDetails map[string][]string) map[string][]string {

	for k, v := range CatalogControlsDetails {
		if _, ok := SecurityControlsDetails[k]; !ok {
			SecurityControlsDetails[k] = v
		}
	}

	return SecurityControlsDetails
}
