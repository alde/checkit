package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/alde/checkit/checkstyle"
	"github.com/alde/checkit/spotbugs"
)

// Checkit base struct
type Checkit struct {
	Files []File `json:"files"`
}

// File information
type File struct {
	Name       string      `json:"name"`
	Violations []Violation `json:"violations"`
}

// Violation data
type Violation struct {
	Line     int    `json:"line"`
	Column   int    `json:"column"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
	Rule     string `json:"rule"`
}

func fromSpotbugs(sb *spotbugs.Spotbugs) *Checkit {
	files := []File{}
	for _, bi := range sb.BugInstance {
		violations := []Violation{}
		var sourcepath = ""
		for _, sl := range bi.SourceLine {
			sourcepath = sl.Sourcepath
			message := spotbugs.ExtractMessage(bi)
			v := Violation{
				Line:     sl.Start,
				Message:  message,
				Severity: "warning",
				Rule:     bi.Type,
			}
			if (Violation{} != v) {
				violations = append(violations, v)
			}
		}
		cf := File{
			Name:       sourcepath,
			Violations: violations,
		}
		files = append(files, cf)
	}

	cs := &Checkit{
		Files: files,
	}

	return cs
}

func fromCheckstyle(cs *checkstyle.Checkstyle) *Checkit {
	files := []File{}
	for _, f := range cs.File {
		if len(f.Error) > 0 {
			violations := []Violation{}
			for _, e := range f.Error {
				v := Violation{
					Column:   e.Column,
					Line:     e.Line,
					Message:  e.Message,
					Severity: e.Severity,
					Rule:     e.Source,
				}
				if (Violation{} != v) {
					violations = append(violations, v)
				}
			}
			if len(violations) > 0 {
				files = append(files, File{
					Name:       f.Name,
					Violations: violations,
				})
			}
		}
	}
	return &Checkit{
		Files: files,
	}
}

func squash(checkits ...*Checkit) *Checkit {
	files := []File{}

	for _, checkit := range checkits {
		for _, f := range checkit.Files {
			files = append(files, f)
		}
	}

	files = unique(files)

	return &Checkit{
		Files: files,
	}
}

func unique(files []File) []File {
	uniqueFiles := []File{}
	seen := make(map[string]bool)
	for _, file := range files {
		filerepr := fmt.Sprintf("%v", file)
		if _, ok := seen[filerepr]; !ok {
			seen[filerepr] = true
			uniqueFiles = append(uniqueFiles, file)
		}
	}
	return uniqueFiles
}

func (c *Checkit) toJSON() string {
	json, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	return string(json)
}
