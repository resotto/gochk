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
	exitMode := flag.Bool("e", false, "flag whether exits with 1 or not when violations occur. (false is default)")
	targetPath := flag.String("t", ".", "target path (\".\" is default)")
	configPath := flag.String("c", "configs/config.json", "configuration file path (\"configs/config.json\" is default)")
	flag.Parse()
	config := gochk.ParseConfig(*configPath)
	results, violated := gochk.Check(*targetPath, config)
	gochk.Show(results, violated, config.PrintViolationsAtTheBottom, *exitMode)
}
