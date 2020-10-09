package gochk

import "fmt"

const (
	red    = "\033[1;31m%-11s%s\n\033[0m"
	yellow = "\033[1;33m%-11s%s\n\033[0m"
	green  = "\033[1;32m%-11s%s\n\033[0m"
	teal   = "\033[1;36m%-11s%s\n\033[0m"
)

const (
	none     = "[None]"
	verified = "[Verified]"
	ignored  = "[Ignored]"
	err      = "[Error]"
)

func printNone(path string) {
	fmt.Printf(teal, none, path)
}

func printVerified(path string) {
	fmt.Printf(green, verified, path)
}

func printIgnored(path string) {
	fmt.Printf(yellow, ignored, path)
}

func printError(filepath string, path string, dependencyOrders []string, currentLayer int, index int) {
	fmt.Printf(red, err, filepath+" imports "+path)
	fmt.Printf(red, "", dependencyOrders[currentLayer]+" depends on "+dependencyOrders[index])
}
