package validation_root

// Markup ...
type Markup struct {
	Raw string `xml:",innerxml" json:"raw,omitempty" yaml:"raw,omitempty"`
}
