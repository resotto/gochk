package gochk

import "fmt"

const (
	red    = "\033[1;31m%s\033[0m"
	yellow = "\033[1;33m%s\033[0m"
	green  = "\033[1;32m%s\033[0m"
	teal   = "\033[1;36m%s\033[0m"
)

func print(color string, message string) {
	fmt.Printf(color, message)
	fmt.Println()
}

func printNone(path string) {
	print(teal, "[None]     "+path)
}

func printVerified(path string) {
	print(green, "[Verified] "+path)
}

func printIgnored(path string) {
	print(yellow, "[Ignored]  "+path)
}

func printError(filepath string, path string, dependencyOrders []string, currentLayer int, index int) {
	print(red, "[Error]    "+filepath+" imports "+path)
	print(red, "           \""+dependencyOrders[currentLayer]+"\" depends on \""+dependencyOrders[index]+"\"")
}
