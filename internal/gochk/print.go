package gochk

import "fmt"

const (
	red    = "\033[1;31m%-11s%s\n\033[0m"
	yellow = "\033[1;33m%-11s%s\n\033[0m"
	green  = "\033[1;32m%-11s%s\n\033[0m"
	teal   = "\033[1;36m%-11s%s\n\033[0m"
)

func printNone(path string) {
	fmt.Printf(teal, "[None]", path)
}

func printVerified(path string) {
	fmt.Printf(green, "[Verified]", path)
}

func printIgnored(path string) {
	fmt.Printf(yellow, "[Ignored]", path)
}

func printError(filepath string, path string, dependencyOrders []string, currentLayer int, index int) {
	fmt.Printf(red, "[Error]", filepath+" imports "+path)
	fmt.Printf(red, "", dependencyOrders[currentLayer]+" depends on "+dependencyOrders[index])
}
