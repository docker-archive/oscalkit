package oscal_source

import (
	"bytes"
	"errors"
	"github.com/docker/oscalkit/pkg/bundled"
	"github.com/docker/oscalkit/pkg/xslt"
	"github.com/docker/oscalkit/types/oscal"
)

func (s *OSCALSource) HTML() (*bytes.Buffer, error) {
	if s.OSCAL().DocumentType() != oscal.CatalogDocument {
		return nil, errors.New("HTML is supported only for OSCAL Catalog")
	}
	transformation, err := bundled.HtmlXslt()
	if err != nil {
		return nil, err
	}
	defer transformation.Cleanup()

	return xslt.Transform(transformation.Path, s.UserPath)
}
