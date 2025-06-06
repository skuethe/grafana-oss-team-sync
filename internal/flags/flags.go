package flags

import (
	"flag"
)

const (
	FlagVersion string = "version"
)

var Version bool

func Load() {
	flag.BoolVar(&Version, FlagVersion, false, "print the version, commit hash and build date and then exit")
	flag.Parse()
}
