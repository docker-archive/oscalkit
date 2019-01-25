package generator

import (
	"fmt"
	"net/url"
	"path"
	"path/filepath"

	"github.com/docker/oscalkit/types/oscal/catalog"

	"github.com/docker/oscalkit/types/oscal/profile"
)

// AppendAlterations appends alter attributes from import chain
func AppendAlterations(p *profile.Profile) (*profile.Profile, error) {
	if p.Modify == nil {
		p.Modify = &profile.Modify{
			Alterations:   []profile.Alter{},
			ParamSettings: []profile.SetParam{},
		}
	}
	for _, i := range p.Imports {
		for _, call := range i.Include.IdSelectors {
			alterFound := false
			for _, alt := range p.Modify.Alterations {
				if call.ControlId == alt.ControlId {
					alterFound = true
					break
				}
			}
			if !alterFound {
				alt, found, err := FindAlter(call, p)
				if err != nil {
					return nil, err
				}
				if found {
					p.Modify.Alterations = append(p.Modify.Alterations, *alt)
				}
			}
		}
	}
	return p, nil
}

//SetBasePath sets up base paths for profiles
func SetBasePath(p *profile.Profile, parentPath string) (*profile.Profile, error) {

	for i, x := range p.Imports {

		if x.Href == nil {
			return nil, fmt.Errorf("href cannot be nil")
		}
		path := fmt.Sprintf("%s/%s", path.Dir(parentPath), path.Base(x.Href.String()))
		path, err := filepath.Abs(path)
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
