package checkstyle

import (
	"encoding/xml"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

type Checkstyle struct {
	Version string   `xml:"version,attr"`
	File    []CSFile `xml:"file"`
}
type CSFile struct {
	Name  string    `xml:"name,attr"`
	Error []CSError `xml:"error"`
}
type CSError struct {
	Line     int    `xml:"line,attr"`
	Column   int    `xml:"column,attr"`
	Severity string `xml:"severity,attr"`
	Message  string `xml:"message,attr"`
	Source   string `xml:"source,attr"`
}

// Process a list of file paths, loading them as Checkstyle structs.
func Process(files []string) (bugs []*Checkstyle) {
	for _, file := range files {
		v := &Checkstyle{}
		f, err := os.Open(file)
		if err != nil {
			logrus.WithError(err).Error("unable to open file for reading")
			continue
		}
		content, _ := ioutil.ReadAll(f)
		err = xml.Unmarshal(content, &v)
		if err != nil {
			logrus.
				WithField("file", file).
				WithError(err).
				Debug("unable to parse file as a Checkstyle xml compatible file")
			continue
		}
		bugs = append(bugs, v)
	}
	return
}
