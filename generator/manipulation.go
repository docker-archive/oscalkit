package generator

import (
	"fmt"

	"github.com/docker/oscalkit/types/oscal/catalog"
	"github.com/docker/oscalkit/types/oscal/profile"
)

func AddPartInCatalog(alterations []profile.Alter, c *catalog.Catalog) *catalog.Catalog {
	for _, x := range alterations {
		for i, g := range c.Groups {
			for j, ctrl := range g.Controls {
				if ctrl.Id == x.ControlId {
					ctrlIncrement := 1
					for _, add := range x.Additions {
						for _, p := range add.Parts {
							appended := false
							for catalogCtrlPartIndex, catalogPart := range ctrl.Parts {
								if p.Class == catalogPart.Class {
									appended = true
									ctrl.Parts[catalogCtrlPartIndex].Id = p.Id + fmt.Sprintf("_%d", ctrlIncrement)
									ctrlIncrement++
									p.Id = p.Id + fmt.Sprintf("_%d", ctrlIncrement)
									position := (ctrlIncrement + catalogCtrlPartIndex) - 1
									ctrl.Parts = append(ctrl.Parts[:position], append([]catalog.Part{p}, ctrl.Parts[position:]...)...)
								}
							}
							if !appended {
								ctrl.Parts = append(ctrl.Parts, p)
							}
						}
					}
					c.Groups[i].Controls[j] = ctrl
				}
				for k, subctrl := range c.Groups[i].Controls[j].Subcontrols {
					subCtrlIncrement := 1
					if subctrl.Id == x.SubcontrolId {
						for _, add := range x.Additions {
							for _, p := range add.Parts {
								appended := false
								for catalogSubCtrlPartIndex, catalogPart := range subctrl.Parts {
									if p.Class == catalogPart.Class {
										appended = true
										subctrl.Parts[catalogSubCtrlPartIndex].Id = p.Id + fmt.Sprintf("_%d", subCtrlIncrement)
										subCtrlIncrement++
										position := (subCtrlIncrement + catalogSubCtrlPartIndex) - 1
										subctrl.Parts = append(subctrl.Parts[:position], append([]catalog.Part{p}, subctrl.Parts[:position]...)...)
									}
								}
								if !appended {
									subctrl.Parts = append(subctrl.Parts, p)
								}
							}

						}
					}
					c.Groups[i].Controls[j].Subcontrols[k] = subctrl

				}
			}
		}
	}
	return c
}
