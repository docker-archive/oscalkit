package generator

import (
	"net/url"
	"testing"
)

func TestIsHttp(t *testing.T) {

	httpRoute := "http://localhost:3000"
	expectedOutputForHTTP := true

	nonHTTPRoute := "NIST.GOV.JSON"
	expectedOutputForNonHTTP := false

	r, err := url.Parse(httpRoute)
	if err != nil {
		t.Error(err)
	}
	if isHTTPResource(r) != expectedOutputForHTTP {
		t.Error("Invalid output for http routes")
	}

	r, err = url.Parse(nonHTTPRoute)
	if err != nil {
		t.Error(err)
	}
	if isHTTPResource(r) != expectedOutputForNonHTTP {
		t.Error("Invalid output for non http routes")
	}

}
