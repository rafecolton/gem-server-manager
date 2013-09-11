package gswat

import (
	"fmt"
	"os"
	"path"
)

var (
	Version  string
	Rev      string
	progName string
)

func init() {
	progName = path.Base(os.Args[0])
	if Version == "" {
		Version = "<unknown>"
	}
	if Rev == "" {
		Rev = "<unknown>"
	}
}

func ProgVersion() string {
	return fmt.Sprintf("%s %s", progName, Version)
}
