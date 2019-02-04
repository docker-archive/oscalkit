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

type Label map[string]string

type ConfigurableValue struct {
	ValueID string `xml:"value-id,attr,omitempty" json:"valueId"`
	Value   string `xml:"value,attr,omitempty" json:"value"`
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
	ParameterID    string        `xml:"parameter-id,attr,omitempty" json:"parameterId,omitempty"`
	ValueID        string        `xml:"value-id,attr,omitempty" json:"valueId,omitempty"`
	Guidance       []string      `xml:"guidance,attr,omitempty" json:"guidance,omitempty"`
	AssessedValue  AssessedValue `xml:"assessed-value,attr,omitempty" json:"assessedValue,omitempty"`
	PossibleValues []string      `xml:"possbile-values,attr,omitempty" json:"possibleValues,omitempty"`
	DefaultValue   string        `xml:"default-value,attr,omitempty" json:"defaultValue,omitempty"`
}

type ImplementsProfile struct {
	ProfileID             string                 `xml:"profile-id,attr,omitempty" json:"profileId"`
	ControlConfigurations []ControlConfiguration `xml:"control-configurations,attr,omitempty" json:"controlConfigurations"`
}

type ControlConfiguration struct {
	ConfigurationIDRef     string                  `xml:"configuration-id-ref,attr,omitempty" json:"configurationIdRef"`
	Parameters             []Parameter             `xml:"parameters,attr,omitempty" json:"parameters"`
	ProvisioningMechanisms []ProvisioningMechanism `xml:"provisioning-mechanisms,attr,omitempty" json:"provisioningMechanisms,omitempty"`
}

type ControlId struct {
	CatalogIDRef   string         `xml:"control-id-ref,attr,omitempty" json:"catalogIdRef,omitempty"`
	ControlID      string         `xml:"control-id,attr,omitempty" json:"controlId,omitempty"`
	ItemID         string         `xml:"item-id,attr,omitempty" json:"itemId,omitempty"`
	AssessmentData AssessmentData `xml:"assesment-data,attr,omitempty" json:"assesmentData,attr,omitempty"`
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

type AssessmentData struct {
	AssessmentID      string             `xml:"assessment-id,attr,omitempty" json:"assessmentId,omitempty"`
	Guidance          string             `xml:"guidance,attr,omitempty" json:"guidance,omitempty"`
	ValidationResults []ValidationResult `xml:"validation-results,attr,omitempty" json:"validationResults,omitempty"`
}

type ValidationResult struct {
	ValidationMechanismRefID string `xml:"validation-mechanism-refid,attr,omitempty" json:"validationMechanismRefId,omitempty"`
	Output                   string `xml:"output,attr,omitempty" json:"output,omitempty"`
	Compliant                bool   `xml:"compliant,attr,omitempty" json:"compliant,omitempty"`
}

//AssessedValue AssessedValue
type AssessedValue struct {
	AssessmentID string `xml:"assessment-id,attr,omitempty" json:"assessmentId,omitempty"`
	Output       string `xml:"output,attr,omitempty" json:"output,omitempty"`
	Compliant    bool   `xml:"compliant,attr,omitempty" json:"compliant,omitempty"`
	Value        string `xml:"value,attr,omitempty" json:"value"`
}

type ProvisioningMechanism struct {
	ProvisioningMechanismRefIds []string
	ProvisionedControls         []ControlId `xml:"provisioned-controls,attr,omitempty" json:"provisionedControls,omitempty"`
}
