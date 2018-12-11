package generator

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"

	"github.com/opencontrol/oscalkit/types/oscal/profile"

	"github.com/opencontrol/oscalkit/types/oscal"
	"github.com/opencontrol/oscalkit/types/oscal/catalog"
)

//ReadCatalog ReadCatalog
func ReadCatalog(r io.Reader) (*catalog.Catalog, error) {

	o, err := oscal.New(r)
	if err != nil {
		return nil, fmt.Errorf("cannot read oscal catalog from file %v,", err)
	}
	return o.Catalog, nil

}

//ReadProfile reads profile from byte array
func ReadProfile(r io.Reader) (*profile.Profile, error) {

	o, err := oscal.New(r)
	if err != nil {
		return nil, fmt.Errorf("cannot read oscal profile from file. err: %v,", err)
	}
	return o.Profile, nil
}

func isHTTPResource(url *url.URL) bool {
	return strings.Contains(url.Scheme, "http")
}

func getName(url *url.URL) string {
	fragments := strings.Split(url.Path, "/")
	return (fragments[len(fragments)-1])
}

func fetchFromHTTPResource(uri *url.URL) ([]byte, error) {
	c := http.Client{Timeout: 10 * time.Second}
	resp, err := c.Get(uri.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read response body %v", err)
	}
	return body, nil

}

//GetAbsolutePath gets absolute file path
func GetAbsolutePath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}
	return filepath.Abs(path)
}

//GetCatalogFilePath GetCatalogFilePath
func GetCatalogFilePath(urlString string) (string, error) {
	uri, err := url.Parse(urlString)
	if err != nil {
		return "", fmt.Errorf("invalid URL pattern %v", err)
	}

	if !isHTTPResource(uri) {
		if filepath.IsAbs(urlString) {
			return urlString, nil
		}
		p, err := filepath.Abs(urlString)
		spew.Dump(p)
		return p, err
	}
	body, err := fetchFromHTTPResource(uri)
	fileName := "/tmp/" + getName(uri)
	f, err := os.Create(fileName)
	if err != nil {
		return "", fmt.Errorf("cannot create json file %v", err)
	}
	defer f.Close()
	_, err = f.Write(body)

	if err != nil {
		return "", fmt.Errorf("cannot write on file %v", err)
	}
	return fileName, nil

}
