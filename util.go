package wrappy

import (
	"os"
	"path/filepath"
	"strings"
)

// combine does a join() on all the non-empty terms
func combine(sep string, terms ...string) string {
	var nonemptyTerms []string
	for _, s := range terms {
		if s != "" {
			nonemptyTerms = append(nonemptyTerms, s)
		}
	}
	return strings.Join(nonemptyTerms, sep)
}

func packagePath(directory string) string {
	fullPath, err := filepath.Abs(directory)
	if err != nil {
		panic(err)
	}

	gopath := os.Getenv("GOPATH")
	path := strings.TrimPrefix(fullPath, gopath+"/src/")

	vendorIdx := strings.Index(path, "/vendor/")
	if vendorIdx != -1 {
		path = path[vendorIdx+len("/vendor/"):]
	}

	return path
}
