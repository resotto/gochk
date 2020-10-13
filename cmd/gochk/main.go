package main

import (
	"flag"

	"github.com/resotto/gochk/internal/gochk"
)

func main() {
	flag.Parse()
	config := gochk.ParseConfig()
	if flag.Arg(0) != "" {
		config.TargetPath = flag.Arg(0)
	}
	results := gochk.Check(config)
	gochk.Show(results)
}
