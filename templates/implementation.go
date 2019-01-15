package templates

import "html/template"

//GetImplementationTemplate gets implementation template for implementation go struct file
func GetImplementationTemplate() (*template.Template, error) {
	return template.New("").Parse(implementationTemplate)
}

const implementationTemplate = `package oscalkit

import (
	"github.com/docker/oscalkit/types/oscal/implementation"
)

var ImplementationGenerated = implementation.Implementation{
	Capabilities: implementation.Capabilities{},
	ComponentDefinitions: []implementation.ComponentDefinition{
		{{range .ComponentDefinitions}}
		implementation.ComponentDefinition{
			ComponentConfigurations: []*implementation.ComponentConfiguration{
					{{range .ComponentConfigurations}}
					&implementation.ComponentConfiguration{
							ID:                     ` + "`{{.ID}}`" + `,
							Name:                   ` + "`{{.Name}}`" + `,
							Description:            ` + "`{{.Description}}`" + `,
							ProvisioningMechanisms: []implementation.Mechanism{},
							ValidationMechanisms:   []implementation.Mechanism{},
							ConfigurableValues: []implementation.ConfigurableValue{
								{{range .ConfigurableValues}}
									implementation.ConfigurableValue{
											Value:   "{{.Value}}",
											ValueID: "{{.ValueID}}",
									},
								{{end}}
							},				
						},
					{{end}}
				},
				ImplementsProfiles: []*implementation.ImplementsProfile{
					{{range .ImplementsProfiles}}
						&implementation.ImplementsProfile{
							ProfileID: "{{.ProfileID}}",
							ControlConfigurations: []implementation.ControlConfiguration{
								{{range .ControlConfigurations}}
									implementation.ControlConfiguration{
										ConfigurationIDRef: "{{.ConfigurationIDRef}}",
										Parameters:         []implementation.Parameter{
											{{range .Parameters}}
											implementation.Parameter{
												Guidance: "{{.Guidance}}",
												ParameterID: "{{.ParameterID}}",
												ValueID: "{{.ValueID}}",
												PossibleValues: []string{
													{{range .PossibleValues}}
														"{{.}}",
													{{end}}
												},												
											},
											{{end}}
										},
									},
								{{end}}
							},
						},
					{{end}}
				},
				ControlImplementations: []*implementation.ControlImplementation{
					{{range .ControlImplementations}}
						&implementation.ControlImplementation{
							ID: "{{.ID}}",
							ControlIds: []implementation.ControlId{
								{{range .ControlIds}}
								implementation.ControlId{
										CatalogIDRef: "{{.CatalogIDRef}}",
										ControlID:	 ` + "`{{.ControlID}}`" + `,
										ItemID: 	 ` + "`{{.ItemID}}`" + `,
									},
								{{end}}
							},
							ControlConfigurations: []implementation.ControlConfiguration{
								{{range .ControlConfigurations}}
									implementation.ControlConfiguration{
										ConfigurationIDRef: "{{.ConfigurationIDRef}}",
										ProvisioningMechanisms: []implementation.ProvisioningMechanism{
											{{range .ProvisioningMechanisms}}
												implementation.ProvisioningMechanism{
													ProvisionedControls: []implementation.ControlId{
														{{range .ProvisionedControls}}
															implementation.ControlId{
																ControlID:    "{{.ControlID}}",
																CatalogIDRef: "{{.CatalogIDRef}}",
																ItemID:       "{{.ItemID}}",
															},
														{{end}}
													},
												},
							
											{{end}}
										},
										Parameters:         []implementation.Parameter{
											{{range .Parameters}}
											implementation.Parameter{
												Guidance: "{{.Guidance}}",
												ParameterID: "{{.ParameterID}}",
												ValueID: "{{.ValueID}}",
												DefaultValue: "{{.DefaultValue}}",
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
		{{end}}
	},
}`
