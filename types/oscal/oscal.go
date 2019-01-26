package oscal

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"

	"github.com/docker/oscalkit/types/oscal/catalog"
	"github.com/docker/oscalkit/types/oscal/profile"
	yaml "gopkg.in/yaml.v2"
)

// OSCAL contains specific OSCAL components
type OSCAL struct {
	XMLName xml.Name         `json:"-" yaml:"-"`
	Catalog *catalog.Catalog `json:"catalog,omitempty" yaml:"catalog,omitempty"`
	// Declarations *Declarations `json:"declarations,omitempty" yaml:"declarations,omitempty"`
	Profile *profile.Profile `json:"profile,omitempty" yaml:"profile,omitempty"`
}

// MarshalXML marshals either a catalog or a profile
func (o *OSCAL) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if o.Catalog != nil {
		o.XMLName = o.Catalog.XMLName
		if err := e.Encode(o.Catalog); err != nil {
			return err
		}
	} else if o.Profile != nil {
		o.XMLName = o.Profile.XMLName
		if err := e.Encode(o.Profile); err != nil {
			return err
		}
	}

	return nil
}

// dockerOptions ...
// type dockerOptions struct {
// 	dockerYAMLFilepath string
// 	dockersDir         string
// }

// NewFromOC initializes an OSCAL type from raw docker data
// func NewFromOC(options dockerOptions) (*OSCAL, error) {
// 	ocFile, err := os.Open(options.dockerYAMLFilepath)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer ocFile.Close()

// 	rawOC, err := ioutil.ReadAll(ocFile)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var oc docker.docker
// 	if err := yaml.Unmarshal(rawOC, &oc); err != nil {
// 		return nil, err
// 	}

// 	ocComponentFileList := []string{}
// 	filepath.Walk(filepath.Join(options.dockersDir, "components/"), func(path string, f os.FileInfo, err error) error {
// 		if !f.IsDir() && (filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml") {
// 			absPath, err := filepath.Abs(path)
// 			if err != nil {
// 				return err
// 			}
// 			ocComponentFileList = append(ocComponentFileList, absPath)
// 		}

// 		return nil
// 	})

// 	ocComponents := []docker.Component{}
// 	for _, ocComponentFilepath := range ocComponentFileList {
// 		ocComponentFile, err := os.Open(ocComponentFilepath)
// 		if err != nil {
// 			return nil, err
// 		}
// 		defer ocComponentFile.Close()

// 		rawOCComponentFile, err := ioutil.ReadAll(ocComponentFile)

// 		var ocComponent docker.Component
// 		if err := yaml.Unmarshal(rawOCComponentFile, &ocComponent); err != nil {
// 			return nil, err
// 		}

// 		ocComponents = append(ocComponents, ocComponent)
// 	}

// 	return convertOC(oc, ocComponents)
// }

// New returns a concrete OSCAL type from a reader
func New(r io.Reader) (*OSCAL, error) {
	oscalBytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	d := xml.NewDecoder(bytes.NewReader(oscalBytes))
	for {
		token, err := d.Token()
		if err != nil || token == nil {
			break
		}
		switch startElement := token.(type) {
		case xml.StartElement:
			switch startElement.Name.Local {
			case "catalog":
				var catalog catalog.Catalog
				if err := d.DecodeElement(&catalog, &startElement); err != nil {
					return nil, err
				}
				return &OSCAL{Catalog: &catalog}, nil

			case "profile":
				var profile profile.Profile
				if err := d.DecodeElement(&profile, &startElement); err != nil {
					return nil, err
				}
				return &OSCAL{Profile: &profile}, nil
			}
		}
	}

	var oscalT map[string]json.RawMessage
	if err := json.Unmarshal(oscalBytes, &oscalT); err == nil {
		for k, v := range oscalT {
			switch k {
			case "catalog":
				var catalog catalog.Catalog
				if err := json.Unmarshal(v, &catalog); err != nil {
					return nil, err
				}
				return &OSCAL{Catalog: &catalog}, nil

			case "profile":
				var profile profile.Profile
				if err := json.Unmarshal(v, &profile); err != nil {
					return nil, err
				}
				return &OSCAL{Profile: &profile}, nil
			}
		}
	}

	return nil, errors.New("Malformed OSCAL. Must be XML or JSON")
}

// XML writes the OSCAL object as XML to the given writer
func (o *OSCAL) XML(w io.Writer, prettify bool) error {
	e := xml.NewEncoder(w)
	if prettify {
		e.Indent("", "  ")
		return e.Encode(o)
	}

	return e.Encode(o)
}

// JSON writes the OSCAL object as JSON to the given writer
func (o *OSCAL) JSON(w io.Writer, prettify bool) error {
	e := json.NewEncoder(w)
	if prettify {
		e.SetIndent("", "  ")
		return e.Encode(o)
	}

	return e.Encode(o)
}

// YAML writes the OSCAL object as YAML to the given writer
func (o *OSCAL) YAML(w io.Writer) error {
	return yaml.NewEncoder(w).Encode(o)
}
