package gochk

import (
	"fmt"
	"log"
)

// Color is ANSI escape code
type color string

const (
	teal   color = "\033[1;36m"
	green        = "\033[1;32m"
	yellow       = "\033[1;33m"
	purple       = "\033[1;35m"
	red          = "\033[1;31m"
	reset        = "\033[0m"
)

// Show prints results
func Show(results []CheckResult) {
	violatesIncluded := false
	for _, r := range results {
		printColorMessage(r)
		if !violatesIncluded && r.resultType == violated {
			violatesIncluded = true
		}
	}
	if violatesIncluded {
		log.Fatal("Dependencies which violate dependency orders found!")
	} else {
		log.Print("No violations")
		printAA()
	}
}

func printColorMessage(cr CheckResult) {
	fmt.Printf("%s%-11s%s\n%s", cr.color, "["+cr.resultType+"]", cr.message, reset)
}

func printAA() {
	aa := []string{
		"    ________     _______       ______    __     __    __   _ _",
		"   /  ______\\   /  ___  \\     /  ____\\  |  |   |  |  |  | /   /",
		"  /  /  ____   /  /   \\  \\   /  /       |  |___|  |  |  |/   /",
		" /  /  |_   | |  |     |  | |  |        |   ___   |  |      /",
		" \\  \\    \\  | |  |     |  | |  |        |  |   |  |  |  |\\  \\",
		"  \\  \\___/  /  \\  \\___/  /   \\  \\_____  |  |   |  |  |  | \\  \\",
		"   \\_______/    \\_______/     \\_______\\ |__|   |__|  |__|  \\__\\",
	}
	for _, s := range aa {
		fmt.Println(s)
	}
}
