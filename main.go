package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/alde/checkit/checkstyle"
	"github.com/alde/checkit/spotbugs"
	"github.com/sirupsen/logrus"
)

func main() {
	pwd, _ := os.Getwd()
	dir := flag.String("dir", pwd, "directory to find compatible files in")
	excludes := flag.String("exclude", "", "paths to exclude (comma-separated list strings)")
	output := flag.String("output", "STDOUT", "file to write to (STDOUT for special case)")
	flag.Parse()

	files := findFiles(*dir, strings.Split(*excludes, ","))

	checkits := []*Checkit{}
	for _, sb := range spotbugs.Process(files) {
		checkits = append(checkits, fromSpotbugs(sb))
	}
	for _, cs := range checkstyle.Process(files) {
		checkits = append(checkits, fromCheckstyle(cs))
	}

	squashed := squash(checkits...)

	if *output == "STDOUT" {
		fmt.Printf("%s\n", squashed.toJSON())
	} else {
		err := ioutil.WriteFile(*output, []byte(squashed.toJSON()), 0644)
		if err != nil {
			logrus.WithError(err).Error("unable to write file")
		}

	}
}

func findFiles(dir string, excludes []string) []string {
	files := []string{}

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		for _, excl := range excludes {
			if strings.Contains(path, excl) {
				continue
			}
		}
		files = append(files, path)
		return nil
	})

	return files
}
