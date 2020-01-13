package generator

import (
	"fmt"

	"github.com/docker/oscalkit/impl"
	"github.com/docker/oscalkit/types/oscal/catalog"
	"github.com/docker/oscalkit/types/oscal/profile"
)

// ProcessAddition processes additions of a profile
func ProcessAddition(alt profile.Alter, controls []catalog.Control) []catalog.Control {
	for j, ctrl := range controls {
		if ctrl.Id == alt.ControlId {
			for _, add := range alt.Additions {
				for _, p := range add.Parts {
					appended := false
					for _, catalogPart := range ctrl.Parts {
						if p.Class == catalogPart.Class {
							appended = true
							// append with all the parts with matching classes
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
		for k, subctrl := range controls[j].Controls {
			if subctrl.Id == alt.ControlId {
				for _, add := range alt.Additions {
					for _, p := range add.Parts {
						appended := false
						for _, catalogPart := range subctrl.Parts {
							if p.Class == catalogPart.Class {
								appended = true
								// append with all the parts
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
			controls[j].Controls[k] = subctrl
		}
	}
	return controls
}

// ProcessAlterations processes alteration section of a profile
func ProcessAlterations(alterations []profile.Alter, c *catalog.Catalog) *catalog.Catalog {
	for _, alt := range alterations {
		for i, g := range c.Groups {
			c.Groups[i].Controls = ProcessAddition(alt, g.Controls)
		}
	}
	return c
}

// ProcessSetParam processes set-param of a profile
func ProcessSetParam(setParams []profile.Set, c *catalog.Catalog, catalogHelper impl.Catalog) *catalog.Catalog {
	for _, sp := range setParams {
		ctrlID := catalogHelper.GetControl(sp.ParamId)
		for i, g := range c.Groups {
			for j, catalogCtrl := range g.Controls {
				if ctrlID == catalogCtrl.Id {
					for k := range catalogCtrl.Parts {
						if len(sp.Constraints) == 0 {
							continue
						}
						c.Groups[i].Controls[j].Parts[k].ModifyProse(sp.ParamId, sp.Constraints[0].Value)
					}
				}
			}
		}
	}
	return c
}

// ModifyParts modifies parts
func ModifyParts(p catalog.Part, controlParts []catalog.Part) []catalog.Part {

	// append with all the parts
	var parts []catalog.Part
	for i, part := range controlParts {
		if p.Class != part.Class {
			parts = append(parts, part)
			continue
		}
		id := part.Id
		part.Id = fmt.Sprintf("%s_%d", id, i)
		parts = append(parts, part)
		part.Id = fmt.Sprintf("%s_%d", id, i+1)
		parts = append(parts, part)
	}
	return parts
}
