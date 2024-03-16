package data

import (
	"path/filepath"
	"runtime"
)

// ----- globals
var basepath string

// ----- functions

// this function will run first when package is imported
func init() {
	_, currentFile, _, _ := runtime.Caller(0)
	basepath = filepath.Dir(currentFile)
}

// return the path of the file relative to this package
func Path(rel string) string {
	if filepath.IsAbs(rel) {
		return rel
	}

	return filepath.Join(basepath, rel)
}
