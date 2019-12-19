// +build ignore

package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/docker/oscalkit/metaschema"
)

const (
	oscalRepo         = "https://github.com/usnistgov/OSCAL.git"
	metaschemaBaseDir = "OSCAL/src/metaschema/%s"
)

var (
	pkgName = map[string]string{
		"catalog":              "catalog",
		"profile":              "profile",
		"system-security-plan": "ssp",
	}
)

func main() {
	rmCmd := exec.Command("rm", "-rf", "OSCAL/")
	if err := rmCmd.Run(); err != nil {
		log.Fatal(err)
	}

	cloneCmd := exec.Command("git", "clone", "--depth", "1", oscalRepo)
	if err := cloneCmd.Run(); err != nil {
		log.Fatal(err)
	}

	metaschemaPaths := map[string]string{
		"catalog": "oscal_catalog_metaschema.xml",
		"profile": "oscal_profile_metaschema.xml",
		"ssp":     "oscal_ssp_metaschema.xml",
	}

	for _, metaschemaPath := range metaschemaPaths {
		f, err := os.Open(fmt.Sprintf(metaschemaBaseDir, metaschemaPath))
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		meta, err := decode(f)
		if err != nil {
			log.Fatal(err)
		}

		if err := metaschema.GenerateTypes(meta); err != nil {
			log.Fatalf("Error generating go types for metaschema: %s", err)
		}
	}

	rmCmd = exec.Command("rm", "-rf", "OSCAL/")
	if err := rmCmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func decode(r io.Reader) (*metaschema.Metaschema, error) {
	var meta metaschema.Metaschema

	d := xml.NewDecoder(r)

	if err := d.Decode(&meta); err != nil {
		return nil, fmt.Errorf("Error decoding metaschema: %s", err)
	}

	for _, imported := range meta.Import {
		if imported.Href == nil {
			return nil, fmt.Errorf("import element in %s is missing 'href' attribute", r)
		}
		imf, err := os.Open(fmt.Sprintf(metaschemaBaseDir, imported.Href.URL.String()))
		if err != nil {
			return nil, err
		}
		defer imf.Close()

		importedMeta, err := decode(imf)
		if err != nil {
			return nil, err
		}

		meta.ImportedMetaschema = append(meta.ImportedMetaschema, *importedMeta)
	}
	err := meta.LinkDefinitions()

	return &meta, err
}
