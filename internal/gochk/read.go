package gochk

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Config is data from config.json
type Config struct {
	TargetPath       string
	DependencyOrders []string
	Ignore           []string
}

// Parse parses config.json
func Parse() Config {
	absPath, _ := filepath.Abs("configs/config.json") // NOTICE: from root directory
	bytes, err := ioutil.ReadFile(absPath)
	if err != nil {
		panic(err)
	}
	var config Config
	json.Unmarshal(bytes, &config)
	return config
}

func walkFiles(cfg Config) ([]dependency, error) {
	violations := make([]dependency, 0, 0)
	return violations, filepath.Walk(cfg.TargetPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}
		if matched, skipType := matchIgnore(cfg.Ignore, path, info); matched {
			return skipType
		}
		violations = append(violations, (checkDependency(cfg.DependencyOrders, path))...)
		return nil
	})
}

func matchIgnore(ignorePaths []string, path string, info os.FileInfo) (bool, error) {
	if included, _ := include(ignorePaths, path); included {
		printIgnored(path)
		if info.IsDir() {
			return true, filepath.SkipDir
		}
		return true, nil
	}
	if info.IsDir() || !strings.Contains(info.Name(), ".go") {
		return true, nil
	}
	return false, nil
}

func retrieveLayers(dependencies []string, path string, currentLayer int) []dependency {
	filepath, _ := filepath.Abs(path)
	f, err := os.Open(filepath)
	defer f.Close()
	if err != nil {
		printWarning(filepath)
		return []dependency{}
	}
	importPaths := readImports(f)
	return retrieveIndices(importPaths, dependencies, path, currentLayer)
}

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
