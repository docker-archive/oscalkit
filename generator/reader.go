package generator

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/opencontrol/oscalkit/types/oscal"
	"github.com/opencontrol/oscalkit/types/oscal/catalog"
)

func readOscal(f *os.File) (*oscal.OSCAL, error) {
	r := bufio.NewReader(f)
	o, err := oscal.New(r)
	if err != nil {
		return nil, err
	}
	return o, nil
}

//ReadCatalog ReadCatalog
func ReadCatalog(f *os.File) (*catalog.Catalog, error) {

	o, err := readOscal(f)
	if err != nil {
		return nil, fmt.Errorf("cannot read oscal catalog from file %v,", err)
	}
	return o.Catalog, nil

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

//GetCatalogFilePath GetCatalogFilePath
func GetCatalogFilePath(urlString string) (string, error) {
	uri, err := url.Parse(urlString)
	if err != nil {
		return "", fmt.Errorf("invalid URL pattern %v", err)
	}

	if !isHTTPResource(uri) {
		return urlString, nil
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
