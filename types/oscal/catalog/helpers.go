package catalog

// ControlOpts to generate controls
type ControlOpts struct {
	Params      []Param
	Parts       []Part
	Subcontrols []Subcontrol
}

// NewPart creates a new part
func NewPart(id, title, narrative string) Part {
	return Part{
		Id:    id,
		Title: Title(title),
		Prose: &Prose{
			P: []P{P{Raw: narrative}},
		},
	}
}

// NewControl creates a new control
func NewControl(id, title string, opts *ControlOpts) Control {
	ctrl := Control{
		Id:    id,
		Title: Title(title),
	}
	if opts != nil {
		ctrl.Subcontrols = opts.Subcontrols
		ctrl.Parts = opts.Parts
		ctrl.Params = opts.Params
	}
	return ctrl
}
