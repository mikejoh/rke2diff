package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mikejoh/rke2diff/internal/buildinfo"
)

type rke2diffOptions struct {
	version bool
}

func main() {
	var rke2diffOpts rke2diffOptions
	flag.BoolVar(&rke2diffOpts.version, "version", false, "Print the version number.")
	flag.Parse()

	if rke2diffOpts.version {
		fmt.Println(buildinfo.Get())
		os.Exit(0)
	}
}
