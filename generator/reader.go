package generator

import (
	"bufio"
	"os"

	"github.com/opencontrol/oscalkit/types/oscal"
	"github.com/opencontrol/oscalkit/types/oscal/catalog"
	"github.com/opencontrol/oscalkit/types/oscal/profile"
)

func readOscal(f *os.File) (*oscal.OSCAL, error) {
	r := bufio.NewReader(f)
	o, err := oscal.New(r)
	if err != nil {
		return nil, err
	}
	return o, nil
}

//ReadCatalog ReadCatalog
func ReadCatalog(f *os.File) (*catalog.Catalog, error) {

	o, err := readOscal(f)
	if err != nil {
		return nil, err
	}
	return o.Catalog, nil

}

//ReadProfile ReadProfile
func ReadProfile(f *os.File) (*profile.Profile, error) {

	o, err := readOscal(f)
	if err != nil {
		return nil, err
	}
	return o.Profile, nil
}
