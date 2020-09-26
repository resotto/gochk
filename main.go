package main

import (
	"flag"

	"github.com/resotto/gochk/cmd/gochk"
)

func main() {
	flag.Parse()
	argPath := flag.Arg(0) // if not provided, argPath is ""
	config := gochk.Parse()

	path := config.DefaultTargetPath
	orders := config.DependencyOrders

	if argPath != "" {
		path = argPath
	}

	gochk.Check(path, orders)
}
