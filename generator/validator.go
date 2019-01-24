package generator

import (
	"fmt"
	"net/url"

	"github.com/docker/oscalkit/types/oscal/catalog"
)

func ValidateHref(href *catalog.Href) error {
	if href == nil {
		return fmt.Errorf("Href cannot be empty")
	}

	_, err := url.Parse(href.String())
	if err != nil {
		return err
	}
	return nil
}
