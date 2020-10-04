package main

import (
	"flag"

	"github.com/resotto/gochk/cmd/gochk"
)

func main() {
	flag.Parse()
	config := gochk.Parse()
	argPath := flag.Arg(0) // if not provided, argPath is ""
	if argPath != "" {
		config.TargetPath = argPath
	}
	gochk.Check(config)
}
