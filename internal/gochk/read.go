package gochk

import (
	"bufio"
	"os"
	"strings"
)

func readImports(filepath string) []string {
	f, _ := os.Open(filepath)
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 6 && strings.EqualFold(line[:6], "import") {
			if strings.Contains(line, "(") {
				return retrieveMultipleImportPath(scanner, line)
			}
			return []string{retrieveImportPath(line)}
		}
	}
	return []string{}
}

func retrieveMultipleImportPath(scanner *bufio.Scanner, line string) []string {
	imports := make([]string, 0, 10)
	for scanner.Scan() {
		line = scanner.Text()
		if strings.Contains(line, ")") {
			break
		} else if strings.EqualFold(line, "") {
			continue
		}
		imports = append(imports, retrieveImportPath(line))
	}
	return imports
}

func retrieveImportPath(line string) string {
	firstQuoIndex := strings.Index(line, "\"")
	return line[firstQuoIndex:]
}
