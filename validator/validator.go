package validator

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/docker/oscalkit/pkg/xml_validation"
	"github.com/santhosh-tekuri/jsonschema"
	"github.com/santhosh-tekuri/jsonschema/loader"
	"github.com/sirupsen/logrus"
)

// Workaround for unpublished schemas referenced by http://csrc.nist.gov/ns/oscal
var basePath string

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

type oscalLoader struct{}

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

func (oscalLoader) Load(url string) (io.ReadCloser, error) {
	return os.Open(filepath.Join(basePath, filepath.Base(url)))
}

// Validate validates one or more JSON files against a specific
// JSON schema.
func (j jsonValidator) Validate(file ...string) error {
	basePath = filepath.Dir(j.SchemaFile)
	schema, err := jsonschema.Compile(j.SchemaFile)
	if err != nil {
		return fmt.Errorf("Error compiling OSCAL schema: %v", err)
	}

	logrus.Debugf("Validating %s against OSCAL schema", file)

	for _, f := range file {
		rawFile, err := os.Open(f)
		if err != nil {
			return fmt.Errorf("Error opening file: %s, %v", f, err)
		}
		defer rawFile.Close()

		if err = schema.Validate(rawFile); err != nil {
			return err
		}

		logrus.Infof("%s is valid against JSON schema %s", f, j.SchemaFile)
	}

	return nil
}

// Validate validates one or more XML files against a specific
// XML schema (.xsd). Wrapper around `xmllint`
func (x xmlValidator) Validate(file ...string) error {
	for _, f := range file {
		rawFile, err := os.Open(f)
		if err != nil {
			return err
		}
		defer rawFile.Close()

		err = xml_validation.Validate(x.SchemaFile, f)
		if err != nil {
			logrus.Error(err)
		} else {
			logrus.Infof("%s is valid against XML schema %s", f, x.SchemaFile)
		}
	}
	return nil
}

func init() {
	loader.Register("http", oscalLoader{})
}
