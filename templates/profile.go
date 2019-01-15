package templates

import (
	"html/template"
)

//GetProfileTemplate GetProfileTemplate
func GetProfileTemplate() *template.Template {

	return template.Must(template.New("").Parse(profileTemplate))
}

const profileTemplate = `
package {{.PackageName}}

import (
	"net/url"

	"github.com/docker/oscalkit/types/oscal/catalog"
	"github.com/docker/oscalkit/types/oscal/profile"

)


var ApplicableProfileControls = profile.Profile{
	Imports: []profile.Import{
		profile.Import{
			{{range .Profile.Imports}}
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
		AsIs: profile.AsIs("{{.Profile.Merge.AsIs}}"), 
	},
	Modify: &profile.Modify{
		Alterations: []profile.Alter{
			{{range .Profile.Modify.Alterations}}
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
