package generator

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/docker/oscalkit/types/oscal"
	"github.com/docker/oscalkit/types/oscal/catalog"
	"github.com/docker/oscalkit/types/oscal/profile"
)

// HTTPFilePath map of http resource against filepath to lessen downloads
type HTTPFilePath struct {
	sync.Mutex
	m map[string]string
}

var pathmap = HTTPFilePath{
	m: make(map[string]string),
}

func findAlter(p *profile.Profile, call profile.Call) (*profile.Alter, bool, error) {

	if p.Modify == nil {
		p.Modify = &profile.Modify{
			Alterations:   []profile.Alter{},
			ParamSettings: []profile.SetParam{},
		}
	}
	for _, alt := range p.Modify.Alterations {
		if EquateAlter(alt, call) {
			return &alt, true, nil
		}
	}
	for _, imp := range p.Imports {
		err := ValidateHref(imp.Href)
		if err != nil {
			return nil, false, err
		}
		path := imp.Href.String()
		if isHTTPResource(imp.Href.URL) {
			pathmap.Lock()
			if v, ok := pathmap.m[imp.Href.String()]; !ok {
				path, err = GetFilePath(imp.Href.String())
				if err != nil {
					return nil, false, err
				}
				pathmap.m[imp.Href.String()] = path
			} else {
				path = v
			}
			pathmap.Unlock()
		}
		f, err := os.Open(path)
		if err != nil {
			return nil, false, err
		}
		defer f.Close()

		o, err := oscal.New(f)
		if err != nil {
			return nil, false, err
		}
		if o.Profile == nil {
			continue
		}
		p, err = SetBasePath(o.Profile, imp.Href.String())
		if err != nil {
			return nil, false, err
		}
		o.Profile = p
		alt, found, err := findAlter(o.Profile, call)
		if err != nil {
			return nil, false, err
		}
		if !found {
			continue
		}
		return alt, true, nil
	}
	return nil, false, nil
}

// EquateAlter equates alter with call
func EquateAlter(alt profile.Alter, call profile.Call) bool {

	if alt.ControlId == "" && alt.SubcontrolId == call.SubcontrolId {
		return true
	}
	if alt.SubcontrolId == "" && alt.ControlId == call.ControlId {
		return true
	}
	return false
}

// GetAlters gets alter attributes from import chain
func GetAlters(p *profile.Profile) ([]profile.Alter, error) {

	var alterations []profile.Alter
	for _, i := range p.Imports {
		for _, call := range i.Include.IdSelectors {
			found := false
			if p.Modify == nil {
				p.Modify = &profile.Modify{
					Alterations:   []profile.Alter{},
					ParamSettings: []profile.SetParam{},
				}
			}
			for _, alt := range p.Modify.Alterations {
				if EquateAlter(alt, call) {
					alterations = append(alterations, alt)
					found = true
					break
				}
			}
			if !found {
				alt, found, err := findAlter(p, call)
				if err != nil {
					return nil, err
				}
				if !found {
					continue
				}
				alterations = append(alterations, *alt)
			}

		}
	}
	return alterations, nil

}

// SetBasePath sets up base paths for profiles
func SetBasePath(p *profile.Profile, parentPath string) (*profile.Profile, error) {
	for i, x := range p.Imports {
		err := ValidateHref(x.Href)
		if err != nil {
			return nil, err
		}
		parentURL, err := url.Parse(parentPath)
		if err != nil {
			return nil, err
		}
		// If the import href is http. Do nothing as it doesn't depend on the parent path
		if isHTTPResource(x.Href.URL) {
			continue
		}
		//if parent is HTTP, and imports are relative, modify imports to http
		if !isHTTPResource(x.Href.URL) && isHTTPResource(parentURL) {
			url, err := makeURL(parentURL, x.Href.URL)
			if err != nil {
				return nil, err
			}
			p.Imports[i].Href = &catalog.Href{URL: url}
			continue
		}
		path := fmt.Sprintf("%s/%s", path.Dir(parentPath), path.Base(x.Href.String()))
		path, err = filepath.Abs(path)
		if err != nil {
			return nil, err
		}
		uri, err := url.Parse(path)
		if err != nil {
			return nil, err
		}
		p.Imports[i].Href = &catalog.Href{URL: uri}
	}
	return p, nil
}

func makeURL(url, child *url.URL) (*url.URL, error) {
	newURL, err := url.Parse(fmt.Sprintf("%s://%s%s/%s", url.Scheme, url.Host, path.Dir(url.Path), child.String()))
	if err != nil {
		return nil, err
	}
	return newURL, nil
}
