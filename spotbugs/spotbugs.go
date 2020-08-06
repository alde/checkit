package spotbugs

import (
	"encoding/xml"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

// SourceLine entry
type SourceLine struct {
	Classname  string `xml:"classname,attr"`
	Start      int    `xml:"start,attr"`
	End        int    `xml:"end,attr"`
	Sourcepath string `xml:"sourcepath,attr"`
	Sourcefile string `xml:"sourcefile,attr"`
	Message    string `xml:"Message"`
}

// BugInstance entry
type BugInstance struct {
	Abbreviation string `xml:"abbrev,attr"`
	Category     string `xml:"category,attr"`
	Priority     int    `xml:"priority,attr"`
	Type         string `xml:"type,attr"`
	ShortMessage string `xml:"ShortMessage"`
	LongMessage  string `xml:"LongMessage"`
	Class        struct {
		Name       string     `xml:"classname,attr"`
		SourceLine SourceLine `xml:"SourceLine"`
		Message    string     `xml:"message"`
	} `xml:"Class"`
	Method struct {
		Name       string     `xml:"name,attr"`
		ClassName  string     `xml:"classname,attr"`
		SourceLine SourceLine `xml:"SourceLine"`
		Message    string     `xml:"message"`
	} `xml:"Method"`
	Field struct {
		Name       string     `xml:"name,attr"`
		ClassName  string     `xml:"classname,attr"`
		SourceLine SourceLine `xml:"SourceLine"`
		Message    string     `xml:"message"`
	} `xml:"Field"`
	SourceLine []SourceLine `xml:"SourceLine"`
}

// The Spotbugs struct represents the spotbugs xml format
type Spotbugs struct {
	AnalysisTimestamp int64  `xml:"analysisTimestamp,attr"`
	Version           string `xml:"version,attr"`
	Timestamp         int64  `xml:"timestamp,attr"`
	Project           struct {
		Name string `xml:"projectName,attr"`
	} `xml:"Project"`
	BugInstance []BugInstance `xml:"BugInstance"`
	Errors      struct {
		MissingClasses int64 `xml:"missingClasses,attr"`
		Errors         int64 `xml:"errors,attr"`
	} `xml:"Errors"`
	FindBugsSummary struct {
		NumPackages int    `xml:"num_packages,attr"`
		TotalBugs   int    `xml:"total_bugs,attr"`
		Timestamp   string `xml:"timestamp,attr"`
		FileStats   []struct {
			Path     string `xml:"path,attr"`
			Size     int    `xml:"size,attr"`
			BugCount int    `xml:"bug_count,attr`
			BugHash  string `xml:"bugHash,attr"`
		} `xml:"FileStats"`
		PackageStats []struct {
			Name       string `xml:"package,attr"`
			TotalBugs  int    `xml:"total_bugs,attr"`
			TotalSize  int    `xml:"total_size,attr"`
			TotalTypes int    `xml:"total_types,attr"`
			ClassStats []struct {
				Name       string `xml:"class,attr"`
				Bugs       string `xml:"bugs,attr"`
				Size       string `xml:"size,attr"`
				Interface  string `xml:"interface,attr"`
				SourceFile string `xml:"sourceFile,attr"`
			} `xml:"ClassStats"`
		} `xml:"PackageStats"`
	} `xml:"FindBugsSummary"`
	ClassFeatures []struct{} `xml:"ClassFeatures"`
	History       []struct{} `xml:"History"`
}

// Process a list of file paths, loading them as Spotbugs structs.
func Process(files []string) (bugs []*Spotbugs) {
	for _, file := range files {
		v := &Spotbugs{}
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
				Debug("unable to parse file as a spotbugs xml compatible file")
			continue
		}
		bugs = append(bugs, v)
	}
	return
}

// ExtractMessage will try its best to find a message to associate with the Spotbugs violation.
func ExtractMessage(bi BugInstance) string {
	if bi.LongMessage != "" {
		return bi.LongMessage
	}
	if bi.ShortMessage != "" {
		return bi.ShortMessage
	}
	firstSourceLineMessage := func(sl []SourceLine) string {
		for _, s := range sl {
			if s.Message != "" {
				return s.Message
			}
		}
		return ""
	}(bi.SourceLine)
	if firstSourceLineMessage != "" {
		return firstSourceLineMessage
	}
	return "no error message found in Spotbugs XML report"
}
