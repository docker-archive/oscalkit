package templates

import (
	"html/template"
)

//GetProfileTemplate GetProfileTemplate
func GetProfileTemplate() *template.Template {

	return template.Must(template.New("").Parse(profileTemplate))
}

const profileTemplate = `
package oscalkit

import (
	"net/url"

	"github.com/docker/oscalkit/types/oscal/catalog"
	"github.com/docker/oscalkit/types/oscal/profile"

)


var ApplicableProfileControls = profile.Profile{
	Imports: []profile.Import{
		profile.Import{
			{{range .Imports}}
				Exclude: &profile.Exclude{
					IdSelectors: []profile.Call{
					},
				},
				Include: &profile.Include{
					IdSelectors: []profile.Call{
						{{range .Include.IdSelectors}}
							profile.Call{
									ControlId: "{{.ControlId}}",
									SubcontrolId: "{{.SubcontrolId}}",
							},
						{{end}}

					},
				},
				Href: &catalog.Href{
					URL: &url.URL{
						RawPath: "{{.Href}}",
					},
				},
			{{end}}
			},
	},
	Merge: &profile.Merge{
		AsIs: profile.AsIs("{{.Merge.AsIs}}"), 
	},
	Modify: &profile.Modify{
		Alterations: []profile.Alter{
			{{range .Modify.Alterations}}
				profile.Alter{
					ControlId: "123",
					Additions: []profile.Add{
						{{range .Additions}}
							profile.Add{
								Title: "{{.Title}}",
								Position: "{{.Position}}",
								Props: []catalog.Prop{
									{{range .Props}}
										catalog.Prop{
											Class: "{{.Class}}",
											Id:    "{{.Id}}",
											Value: "{{.Value}}",
										},
									{{end}}
								},
		
							},
						{{end}}
					},
				},
			{{end}}
		},
	},
}
`
