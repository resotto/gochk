package gochk

import (
	"bufio"
	"os"
	"strings"
)

func readImports(f *os.File) []string {
	scanner := bufio.NewScanner(f)
	skipToImportStatement(scanner)
	scanner.Scan()
	if line := scanner.Text(); len(line) > 6 && strings.EqualFold(line[:6], "import") {
		if strings.Contains(line, "(") {
			return retrieveMultipleImportPath(scanner, line)
		}
		return []string{retrieveImportPath(line)}
	}
	return []string{}
}

func skipToImportStatement(scanner *bufio.Scanner) {
	scanner.Scan()
	line := scanner.Text()
	skipBlockComments(line, scanner)
	for true {
		if line := scanner.Text(); len(line) > 7 && strings.EqualFold(line[:7], "package") {
			scanner.Scan() // Points to two lines below the "package" declaration
			return
		}
		scanner.Scan()
	}
}

func skipBlockComments(line string, scanner *bufio.Scanner) {
	if strings.EqualFold(line, "/*") {
		for scanner.Scan() {
			if line := scanner.Text(); strings.EqualFold(line, "*/") {
				return
			}
		}
	}
}

func retrieveMultipleImportPath(scanner *bufio.Scanner, line string) []string {
	imports := make([]string, 0, 10)
	for scanner.Scan() {
		line = scanner.Text()
		if strings.EqualFold(line, ")") {
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
