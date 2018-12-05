package catalog

import (
	"encoding/json"
	"encoding/xml"
	"net/url"
)

type Href struct {
	*url.URL
}

// UnmarshalXMLAttr unmarshals an href to a url.URL
func (h *Href) UnmarshalXMLAttr(attr xml.Attr) error {
	url, err := url.Parse(attr.Value)
	if err != nil {
		return err
	}

	*h = Href{url}

	return nil
}

func (h *Href) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = s[1 : len(s)-1]
	url, err := url.Parse(s)
	if err != nil {
		return err
	}
	h.URL = url
	return nil
}

func (h *Href) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String())
}

func (h *Href) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	if h.URL != nil {
		rawURI := h.URL.String()

		return xml.Attr{Name: name, Value: rawURI}, nil
	}

	return xml.Attr{Name: name}, nil
}
