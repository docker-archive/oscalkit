package implementation

type ComponentDefinition struct {
	ID                      string                    `xml:"id,attr,omitempty" json:"id"`
	Name                    string                    `xml:"name,attr,omitempty" json:"name,omitempty"`
	ComponentType           string                    `xml:"component-type,attr,omitempty" json:"componentType,omitempty"`
	Version                 string                    `xml:"version,attr,omitempty" json:"version,omitempty"`
	ComponentConfigurations []*ComponentConfiguration `xml:"component-configurations,attr,omitempty" json:"componentConfigurations,omitempty"`
	ImplementsProfiles      []*ImplementsProfile      `xml:"implements-profile,attr,omitempty" json:"implementsProfiles,omitempty"`
	ControlImplementations  []*ControlImplementation  `xml:"control-implementations,attr,omitempty" json:"controlImplementations,omitempty"`
	Relationships           []*Relationship           `xml:"relationships,attr,omitempty" json:"relationships,omitempty"`
}
type Mechanism struct {
	ID   string `xml:"id,attr,omitempty" json:"id"`
	Type string `xml:"type,attr,omitempty" json:"type"`
	Data string `xml:"data,attr,omitempty" json:"data"`
}

type Label struct {
	AdminSetting string `xml:"admin-setting,attr,omitempty" json:"adminSetting"`
}

type ConfigurableValue struct {
	ValueID string `xml:"value-id,attr,omitempty" json:"valueId"`
	Value   int    `xml:"value,attr,omitempty" json:"value"`
}

type ComponentConfiguration struct {
	ID                     string              `xml:"id,attr,omitempty" json:"id"`
	Labels                 Label               `xml:"labels,attr,omitempty" json:"labels"`
	Name                   string              `xml:"name,attr,omitempty" json:"name"`
	Description            string              `xml:"description,attr,omitempty" json:"description"`
	ProvisioningMechanisms []Mechanism         `xml:"provision-mechanisms,attr,omitempty" json:"provisioningMechanisms"`
	ValidationMechanisms   []Mechanism         `xml:"validation-mechanisms,attr,omitempty" json:"validationMechanisms"`
	ConfigurableValues     []ConfigurableValue `xml:"configurable-values,attr,omitempty" json:"configurableValues"`
}

type Parameter struct {
	ParameterID string `xml:"parameter-id,attr,omitempty" json:"parameterId"`
	ValueID     string `xml:"value-id,attr,omitempty" json:"valueId"`
	Value       string `xml:"value,attr,omitempty" json:"value"`
	Guidance    string `xml:"guidance,attr,omitempty" json:"guidance"`
}

type ImplementsProfile struct {
	ProfileID             string                 `xml:"profile-id,attr,omitempty" json:"profileId"`
	ControlConfigurations []ControlConfiguration `xml:"control-configurations,attr,omitempty" json:"controlConfigurations"`
}

type ControlConfiguration struct {
	ConfigurationIDRef string      `xml:"configuration-id-ref,attr,omitempty" json:"configurationIdRef"`
	Parameters         []Parameter `xml:"parameters,attr,omitempty" json:"parameters"`
}

type ControlId struct {
	CatalogIDRef string `xml:"control-id-ref,attr,omitempty" json:"catalogIdRef,omitempty"`
	ControlID    string `xml:"control-id,attr,omitempty" json:"controlId,omitempty"`
	ItemID       string `xml:"item-id,attr,omitempty" json:"itemId,omitempty"`
}

type Relationship struct {
	IDRef       string `xml:"id-ref,attr,omitempty" json:"idRef"`
	Type        string `xml:"type,attr,omitempty" json:"type"`
	Cardinality string `xml:"cadinality,attr,omitempty" json:"cardinality"`
}

type ValidationMechanism struct {
	ValidationMechanismRefIds []string    `xml:"validation-mechanism-ref-ids,attr,omitempty" json:"validationMechanismRefIds"`
	ValidatedControls         []ControlId `xml:"validated-controls,attr,omitempty" json:"validatedControls"`
}

type ControlImplementation struct {
	ID                       string                 `xml:"id,attr,omitempty" json:"id"`
	ControlIds               []ControlId            `xml:"control-ids,attr,omitempty" json:"controlIds"`
	SatisfactionRequirements string                 `xml:"satisfaction-requirements,attr,omitempty" json:"satisfactionRequirements"`
	Guidance                 string                 `xml:"guidance,attr,omitempty" json:"guidance"`
	ControlConfigurations    []ControlConfiguration `xml:"control-configurations,attr,omitempty" json:"controlConfigurations"`
	ValidationMechanisms     []ValidationMechanism  `xml:"validation-mechanisms,attr,omitempty" json:"validationMechanisms"`
	Parameters               []Parameter            `xml:"parameters,attr,omitempty" json:"parameters"`
}

type Implementation struct {
	ComponentDefinitions    []ComponentDefinition   `xml:"component-definitions,attr,omitempty" json:"componentDefinitions"`
	Capabilities            Capabilities            `xml:"capabilities,attr,omitempty" json:"capabilities,omitempty"`
	ComponentSpecifications ComponentSpecifications `xml:"component-specificiations,attr,omitempty" json:"component-specifications,omitempty"`
}

type Capabilities struct{}

type ComponentSpecifications struct{}
