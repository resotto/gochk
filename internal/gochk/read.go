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
	multipleImport := false
	imports := make([]string, 0, 10)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 6 && strings.EqualFold(line[:6], "import") {
			if strings.Contains(line, "(") {
				multipleImport = true
				continue
			}
			imports = append(imports, retrievePath(line))
			break
		}
		if multipleImport {
			if strings.Contains(line, ")") {
				break
			} else if strings.EqualFold(line, "") {
				continue
			}
			imports = append(imports, retrievePath(line))
		}
	}
	return imports
}

func retrievePath(line string) string {
	firstQuoIndex := strings.Index(line, "\"")
	return line[firstQuoIndex:]
}
