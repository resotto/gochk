package main

import (
	"flag"

	"github.com/resotto/gochk/internal/gochk"
)

func main() {
	flag.Parse()
	config := gochk.Parse()
	if flag.Arg(0) != "" {
		config.TargetPath = flag.Arg(0)
	}
	gochk.Check(config)
}
