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
