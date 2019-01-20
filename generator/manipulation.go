package generator

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/docker/oscalkit/types/oscal"
	"github.com/docker/oscalkit/types/oscal/catalog"
	"github.com/docker/oscalkit/types/oscal/profile"
)

//ProcessAddition processess additions of a profile
func ProcessAddition(alt profile.Alter, controls []catalog.Control) []catalog.Control {
	for j, ctrl := range controls {
		if ctrl.Id == alt.ControlId {
			for _, add := range alt.Additions {
				for _, p := range add.Parts {
					appended := false
					for _, catalogPart := range ctrl.Parts {
						if p.Class == catalogPart.Class {
							appended = true
							//append with all the parts with matching classes
							parts := ModifyParts(p, ctrl.Parts)
							ctrl.Parts = parts
						}
					}
					if !appended {
						ctrl.Parts = append(ctrl.Parts, p)
					}
				}
			}
			controls[j] = ctrl
		}
		for k, subctrl := range controls[j].Subcontrols {
			if subctrl.Id == alt.SubcontrolId {
				for _, add := range alt.Additions {
					for _, p := range add.Parts {
						appended := false
						for _, catalogPart := range subctrl.Parts {
							if p.Class == catalogPart.Class {
								appended = true
								//append with all the parts
								parts := ModifyParts(p, subctrl.Parts)
								subctrl.Parts = parts
							}
						}
						if !appended {
							subctrl.Parts = append(subctrl.Parts, p)
						}
					}

				}
			}
			controls[j].Subcontrols[k] = subctrl
		}
	}
	return controls
}

//ProcessAlteration processess alteration section of a profile
func ProcessAlteration(alterations []profile.Alter, c *catalog.Catalog) *catalog.Catalog {
	for _, alt := range alterations {
		for i, g := range c.Groups {
			c.Groups[i].Controls = ProcessAddition(alt, g.Controls)
		}
	}
	return c
}

//ModifyParts modifies parts
func ModifyParts(p catalog.Part, controlParts []catalog.Part) []catalog.Part {

	//append with all the parts
	var parts []catalog.Part
	for i, part := range controlParts {
		if p.Class != part.Class {
			parts = append(parts, part)
			continue
		}
		id := part.Id
		part.Id = fmt.Sprintf("%s_%d", id, i+1)
		parts = append(parts, part)
		part.Id = fmt.Sprintf("%s_%d", id, i+2)
		parts = append(parts, part)
	}
	return parts
}

//FindAlter finds alter manipulation attribute in the profile import chain
func FindAlter(call profile.Call, p *profile.Profile) (*profile.Alter, error) {

	ec := make(chan error)
	altCh := make(chan *profile.Alter)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	for _, i := range p.Imports {
		go func(i profile.Import) {
			traverseProfile(ctx, call, p, altCh, ec)
		}(i)
	}

	select {
	case alt := <-altCh:
		return alt, nil
	case err := <-ec:
		return nil, err
	case <-ctx.Done():
		return nil, timeoutErr(call)
	}

}

func traverseProfile(ctx context.Context, call profile.Call, p *profile.Profile, altCh chan *profile.Alter, errCh chan error) {

	if p == nil {
		errCh <- fmt.Errorf("profile cannot be nil")
		return
	}
	if p.Modify == nil {
		errCh <- fmt.Errorf("modify is nil")
		return
	}
	for _, alt := range p.Modify.Alterations {
		if alt.ControlId == call.ControlId {
			altCh <- &alt
			return
		}
		if alt.SubcontrolId == call.SubcontrolId {
			altCh <- &alt
			return
		}
	}

	for _, imp := range p.Imports {
		go func(imp profile.Import) {
			if imp.Href == nil {
				errCh <- fmt.Errorf("import href cannot be nil")
				return
			}
			path, err := GetFilePath(imp.Href.String())
			if err != nil {
				errCh <- err
				return
			}
			f, err := os.Open(path)
			if err != nil {
				errCh <- err
				return
			}
			o, err := oscal.New(f)
			if err != nil {
				errCh <- err
				return
			}
			if o.Profile == nil {
				logrus.Warn("catalog found")
				return
			}
			traverseProfile(ctx, call, o.Profile, altCh, errCh)
		}(imp)

	}

}

func timeoutErr(call profile.Call) error {
	if call.ControlId != "" {
		return fmt.Errorf("unable to find control id %s in the profile import chain, needs review", call.ControlId)
	}
	if call.SubcontrolId != "" {
		return fmt.Errorf("unable to find sub-control id %s in the profile import chain, needs review", call.SubcontrolId)
	}
	return nil
}
