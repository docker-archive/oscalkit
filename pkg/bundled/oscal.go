package bundled

import (
	"github.com/markbates/pkger"
	"io"
	"io/ioutil"
	"os"
)

type BundledFile struct {
	Path string
}

func HtmlXslt() (*BundledFile, error) {
	in, err := pkger.Open("/OSCAL/src/utils/util/publish/XSLT/oscal-browser-display.xsl")
	if err != nil {
		return nil, err
	}
	defer in.Close()

	out, err := ioutil.TempFile("/tmp", "oscal_xslt")
	if err != nil {
		return nil, err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return nil, err
	}
	return &BundledFile{Path: out.Name()}, nil
}

func (f *BundledFile) Cleanup() {
	os.Remove(f.Path)
}
