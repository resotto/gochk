/*

gochk checks whether .go files violate Clean Architecture The Dependency Rule or not, and prints its results.

https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html#the-dependency-rule

Usage:
	gochk [flag]

The flags are:
	-t
		Target path. Default value is ".".
	-c
		The path of config.json. Default value is "configs/config.json".

Example:
	gochk -t=../../../goilerplate -c=../../configs/config.json

*/
package main

import (
	"flag"

	"github.com/resotto/gochk/internal/gochk"
)

func main() {
	targetPath := flag.String("t", ".", "target path")
	configPath := flag.String("c", "configs/config.json", "configuration file path")
	flag.Parse()
	config := gochk.ParseConfig(*configPath)
	results, violated := gochk.Check(*targetPath, config)
	gochk.Show(results, violated, config.PrintViolationsAtTheBottom)
}
