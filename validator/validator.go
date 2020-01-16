package validator

import (
	"path/filepath"

	"github.com/docker/oscalkit/pkg/json_validation"
	"github.com/docker/oscalkit/pkg/xml_validation"
	"github.com/sirupsen/logrus"
)

// Validator ...
type Validator interface {
	Validate(file ...string) error
}

type jsonValidator struct {
	SchemaFile string
}

type xmlValidator struct {
	SchemaFile string
}

// New creates a Validator based on the specified schema file
func New(schemaFile string) Validator {
	switch filepath.Ext(schemaFile) {
	case ".json":
		return jsonValidator{SchemaFile: schemaFile}

	case ".xsd":
		return xmlValidator{SchemaFile: schemaFile}
	}

	return nil
}

// Validate validates one or more JSON files against a specific
// JSON schema.
func (j jsonValidator) Validate(file ...string) error {
	for _, f := range file {
		err := json_validation.Validate(j.SchemaFile, f)
		if err != nil {
			logrus.Error(err)
		} else {
			logrus.Infof("%s is valid against JSON schema %s", f, j.SchemaFile)
		}
	}

	return nil
}

// Validate validates one or more XML files against a specific
// XML schema (.xsd). Wrapper around `xmllint`
func (x xmlValidator) Validate(file ...string) error {
	for _, f := range file {
		err := xml_validation.Validate(x.SchemaFile, f)
		if err != nil {
			logrus.Error(err)
		} else {
			logrus.Infof("%s is valid against XML schema %s", f, x.SchemaFile)
		}
	}
	return nil
}
