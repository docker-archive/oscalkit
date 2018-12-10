package generator

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/opencontrol/oscalkit/types/oscal/profile"

	"github.com/opencontrol/oscalkit/types/oscal/catalog"
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

func TestGenerateCatalogs(t *testing.T) {

	catalogs := []*catalog.Catalog{
		&catalog.Catalog{
			Title: "TestCatalog1",
		},
		&catalog.Catalog{
			Title: "TestCatalog2",
		},
	}
	f, err := os.Open(temporaryFilePathForCatalogsGoFile)
	if err != nil {
		f, err = os.Create(temporaryFilePathForCatalogsGoFile)
		if err != nil {
			t.Error("cannot create file")
		}
	}
	defer func() {
		err = os.Remove(temporaryFilePathForCatalogsGoFile)
		if err != nil {
			t.Error("canont remove file")
		}
	}()
	defer f.Close()
	err = GenerateCatalogs(f, catalogs)
	if err != nil {
		t.Error(err)
	}
	x, err := os.Open(temporaryFilePathForCatalogsGoFile)
	if err != nil {
		t.Error(err)
	}
	defer x.Close()
	b, err := ioutil.ReadAll(x)
	if err != nil {
		t.Error(err)
	}
	output := strings.TrimSpace(string(b))

	if !strings.Contains(output, "TestCatalog1") {
		t.Error("does not contain TestCatalog1 in file")
	}
	if !strings.Contains(output, "TestCatalog2") {
		t.Error("does not contain TestCatalog2 in file")
	}

}

func TestGenerateCatalogsWithEmptyCatalogs(t *testing.T) {

	dummyFile := os.File{}
	err := GenerateCatalogs(&dummyFile, []*catalog.Catalog{})
	if err == nil {
		t.Error("should return error")
	}

}

func TestReadCatalog(t *testing.T) {

	catalogTitle := "NIST SP800-53"
	bytesToWriteInCatalogJSONfile := []byte(string(
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
		}`, catalogTitle)))

	f, err := os.Create(temporaryFilePathForCatalogJSON)
	if err != nil {
		t.Error("cannot create file")
	}
	defer f.Close()
	defer func() {
		err := os.Remove(temporaryFilePathForCatalogJSON)
		if err != nil {
			t.Error("cannot delete file")
		}
	}()

	_, err = f.Write(bytesToWriteInCatalogJSONfile)
	if err != nil {
		t.Error(err)
	}
	x, err := os.Open(temporaryFilePathForCatalogJSON)
	if err != nil {
		t.Error(err)
	}
	defer x.Close()

	c, err := ReadCatalog(x)
	if err != nil {
		t.Error(err)
	}

	if c.Title != catalog.Title(catalogTitle) {
		t.Error("title not equal")
	}

}

func TestReadInvalidCatalog(t *testing.T) {

	invalidBytes := []byte(string(`{ "catalog": "some dummy bad json"}`))
	f, err := os.Create(temporaryFilePathForCatalogJSON)
	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	defer func() {
		err := os.Remove(temporaryFilePathForCatalogJSON)
		if err != nil {
			t.Error(err)
		}
	}()
	f.Write(invalidBytes)

	x, err := os.Open(temporaryFilePathForCatalogJSON)
	if err != nil {
		t.Error(err)
	}
	defer x.Close()
	_, err = ReadCatalog(x)
	if err == nil {
		t.Error("successfully parsed invalid catalog file")
	}
}

func TestIntersectProfile(t *testing.T) {

	href, _ := url.Parse("https://raw.githubusercontent.com/usnistgov/OSCAL/master/content/nist.gov/SP800-53/rev4/NIST_SP-800-53_rev4_catalog.json")
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
	x := IntersectProfile(&p)
	if len(x) != 1 {
		t.Error("there must be one catalog")
	}
	if x[0].Groups[0].Controls[0].Id != "ac-1" {
		t.Error("Invalid control Id")
	}

}

func TestIntersectionWithInvalidHref(t *testing.T) {

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
	catalogs := IntersectProfile(&p)
	if len(catalogs) > 0 {
		t.Error("nothing should be parsed due to bad url")
	}
}

func TestGetCatalogInvalidFilePath(t *testing.T) {

	url := "http://[::1]a"
	_, err := GetCatalogFilePath(url)
	if err == nil {
		t.Error("should fail")
	}
}

func failTest(err error, t *testing.T) {
	if err != nil {
		t.Error(t)
	}
}
