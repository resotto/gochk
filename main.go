package main

import (
	"flag"

	"github.com/resotto/gochk/cmd/gochk"
)

func main() {
	flag.Parse()
	config := gochk.Parse()
	if flag.Arg(0) != "" { // flag.Arg(0) equals to "" if not provided
		config.TargetPath = flag.Arg(0)
	}
	gochk.Check(config)
}
