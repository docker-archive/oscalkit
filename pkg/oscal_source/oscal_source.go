package oscal_source

import (
	"fmt"
	"github.com/docker/oscalkit/types/oscal"
	"os"
	"path/filepath"
	"strings"
)

type DocumentFormat int

const (
	UnknownFormat DocumentFormat = iota
	XmlFormat
	JsonFormat
	YamlFormat
)

// OSCALSource is intermediary that handles IO and low-level common operations consistently for oscalkit
type OSCALSource struct {
	UserPath string
	file     *os.File
	oscal    *oscal.OSCAL
}

// Open creates new OSCALSource and load it up
func Open(path string) (*OSCALSource, error) {
	result := OSCALSource{UserPath: path}
	return &result, result.open()
}

func (s *OSCALSource) open() error {
	var err error
	path := s.UserPath
	if !filepath.IsAbs(path) {
		if path, err = filepath.Abs(path); err != nil {
			return fmt.Errorf("Cannot get absolute path: %v", err)
		}
	}
	if _, err = os.Stat(path); err != nil {
		return fmt.Errorf("Cannot stat %s, %v", path, err)
	}
	if s.file, err = os.Open(path); err != nil {
		return fmt.Errorf("Cannot open file %s: %v", path, err)
	}
	if s.oscal, err = oscal.New(s.file); err != nil {
		return fmt.Errorf("Cannot parse file: %v", err)
	}
	return nil
}

func (s *OSCALSource) OSCAL() *oscal.OSCAL {
	return s.oscal
}

func (s *OSCALSource) DocumentFormat() DocumentFormat {
	if strings.HasSuffix(s.UserPath, ".xml") {
		return XmlFormat
	} else if strings.HasSuffix(s.UserPath, ".json") {
		return JsonFormat
	} else {
		return UnknownFormat
	}

}

// Close the OSCALSource
func (s *OSCALSource) Close() {
	if s.file != nil {
		s.file.Close()
		s.file = nil
	}
}
