package generator

import (
	"github.com/docker/oscalkit/types/oscal/profile"
)

//AppendAlterations appends alter attributes from import chain
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
				alt, err := FindAlter(call, p)
				if err != nil {
					return nil, err
				}
				p.Modify.Alterations = append(p.Modify.Alterations, *alt)
			}
		}
	}
	return p, nil
}
