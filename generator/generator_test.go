package generator

import (
	"bytes"
	"fmt"
	"net/url"
	"testing"

	"github.com/docker/oscalkit/types/oscal/catalog"
	"github.com/docker/oscalkit/types/oscal/profile"
)

const (
	temporaryFilePathForCatalogJSON    = "/tmp/catalog.json"
	temporaryFilePathForProfileJSON    = "/tmp/profile.json"
	temporaryFilePathForCatalogsGoFile = "/tmp/catalogs.go"
)

func TestIsHttp(t *testing.T) {

	httpRoute := "http://localhost:3000"
	expectedOutputForHTTP := true

	nonHTTPRoute := "NIST.GOV.JSON"
	expectedOutputForNonHTTP := false

	r, err := url.Parse(httpRoute)
	if err != nil {
		t.Error(err)
	}
	if isHTTPResource(r) != expectedOutputForHTTP {
		t.Error("Invalid output for http routes")
	}

	r, err = url.Parse(nonHTTPRoute)
	if err != nil {
		t.Error(err)
	}
	if isHTTPResource(r) != expectedOutputForNonHTTP {
		t.Error("Invalid output for non http routes")
	}

}

func TestReadCatalog(t *testing.T) {

	catalogTitle := "NIST SP800-53"
	r := bytes.NewReader([]byte(string(
		fmt.Sprintf(`
		{
			"catalog": {
				"title": "%s",
				"declarations": {
					"href": "NIST_SP-800-53_rev4_declarations.xml"
				},
				"groups": [
					{
						"controls": [
							{
								"id": "at-1",
								"class": "SP800-53",
								"title": "Security Awareness and Training Policy and Procedures",
								"params": [
									{
										"id": "at-1_prm_1",
										"label": "organization-defined personnel or roles"
									},
									{
										"id": "at-1_prm_2",
										"label": "organization-defined frequency"
									},
									{
										"id": "at-1_prm_3",
										"label": "organization-defined frequency"
									}
								]
							}
						]
					}
				]
			}
		}`, catalogTitle))))

	c, err := ReadCatalog(r)
	if err != nil {
		t.Error(err)
	}

	if c.Title != catalog.Title(catalogTitle) {
		t.Error("title not equal")
	}

}

func TestReadInvalidCatalog(t *testing.T) {

	r := bytes.NewReader([]byte(string(`{ "catalog": "some dummy bad json"}`)))
	_, err := ReadCatalog(r)
	if err == nil {
		t.Error("successfully parsed invalid catalog file")
	}
}

func TestCreateCatalogsFromProfile(t *testing.T) {

	href, _ := url.Parse("https://raw.githubusercontent.com/usnistgov/OSCAL/master/content/nist.gov/SP800-53/rev4/NIST_SP-800-53_rev4_catalog.xml")
	p := profile.Profile{
		Imports: []profile.Import{
			profile.Import{
				Href: &catalog.Href{
					URL: href,
				},
				Include: &profile.Include{
					IdSelectors: []profile.Call{
						profile.Call{
							ControlId: "ac-1",
						},
					},
				},
			},
		},
	}
	x, err := CreateCatalogsFromProfile(&p)
	if err != nil {
		t.Errorf("error should be null")
	}
	if len(x) != 1 {
		t.Error("there must be one catalog")
	}
	if x[0].Groups[0].Controls[0].Id != "ac-1" {
		t.Error("Invalid control Id")
	}

}

func TestCreateCatalogsFromProfileWithBadHref(t *testing.T) {

	href, _ := url.Parse("this is a bad url")
	p := profile.Profile{
		Imports: []profile.Import{
			profile.Import{
				Href: &catalog.Href{
					URL: href,
				},
				Include: &profile.Include{
					IdSelectors: []profile.Call{
						profile.Call{
							ControlId: "ac-1",
						},
					},
				},
			},
		},
	}
	catalogs, err := CreateCatalogsFromProfile(&p)
	if err != nil {
		t.Error("error should  be nil")
	}
	if len(catalogs) > 0 {
		t.Error("nothing should be parsed due to bad url")
	}
}

func TestSubControlsMapping(t *testing.T) {

	profile := profile.Profile{
		Imports: []profile.Import{
			profile.Import{
				Href: &catalog.Href{
					URL: func() *url.URL {
						url, _ := url.Parse("https://raw.githubusercontent.com/usnistgov/OSCAL/master/content/nist.gov/SP800-53/rev4/NIST_SP-800-53_rev4_catalog.xml")
						return url
					}(),
				},
				Include: &profile.Include{
					IdSelectors: []profile.Call{
						profile.Call{
							ControlId: "ac-1",
						},
						profile.Call{
							ControlId: "ac-2",
						},
						profile.Call{
							SubcontrolId: "ac-2.1",
						},
						profile.Call{
							SubcontrolId: "ac-2.2",
						},
					},
				},
			},
		},
	}

	c, err := CreateCatalogsFromProfile(&profile)
	if err != nil {
		t.Error("error should be nil")
	}
	if c[0].Groups[0].Controls[1].Subcontrols[0].Id != "ac-2.1" {
		t.Errorf("does not contain ac-2.1 in subcontrols")
	}

}

func TestGetCatalogInvalidFilePath(t *testing.T) {

	url := "http://[::1]a"
	_, err := GetFilePath(url)
	if err == nil {
		t.Error("should fail")
	}
}

func failTest(err error, t *testing.T) {
	if err != nil {
		t.Error(t)
	}
}
