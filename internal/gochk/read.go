package gochk

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Config is converted data of config.json
type Config struct {
	TargetPath                 string
	DependencyOrders           []string
	Ignore                     []string
	PrintViolationsAtTheBottom bool
}

type resultType string

const (
	none     resultType = "None"
	verified            = "Verified"
	ignored             = "Ignored"
	warning             = "Warning"
	violated            = "Violated"
)

// CheckResult is the result of dependency checking
type CheckResult struct {
	resultType resultType
	message    string
	color      color
}

type dependency struct {
	filePath    string
	fileLayer   int
	importPath  string
	importLayer int
}

// ParseConfig parses config.json
func ParseConfig() Config {
	absPath, _ := filepath.Abs("configs/config.json") // NOTICE: from root directory
	bytes, err := ioutil.ReadFile(absPath)
	if err != nil {
		panic(err)
	}
	var config Config
	json.Unmarshal(bytes, &config)
	return config
}

// Check checks dependencies
func Check(cfg Config) ([]CheckResult, bool) {
	violated := false
	results := make([]CheckResult, 0, 1000)
	filepath.Walk(cfg.TargetPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			results = append([]CheckResult{CheckResult{resultType: warning, message: err.Error(), color: purple}}, results...)
			return nil
		}
		if matched, skipType := matchIgnore(cfg.Ignore, path, info); matched {
			results = append([]CheckResult{CheckResult{resultType: ignored, message: path, color: yellow}}, results...)
			return skipType
		}
		if info.IsDir() || !strings.Contains(info.Name(), ".go") {
			return nil
		}
		violated = violated || setResultType(&results, cfg.DependencyOrders, path)
		return nil
	})
	return results, violated
}

func matchIgnore(ignorePaths []string, path string, info os.FileInfo) (bool, error) {
	if included, _ := include(ignorePaths, path); included {
		if info.IsDir() {
			return true, filepath.SkipDir
		}
		return true, nil
	}
	return false, nil
}

func retrieveDependencies(dependencyOrders []string, path string, currentLayer int) ([]dependency, error) {
	filepath, _ := filepath.Abs(path)
	f, err := os.Open(filepath)
	defer f.Close()
	if err != nil {
		return []dependency{}, err
	}
	importPaths := readImports(f)
	dependencies := make([]dependency, 0, len(importPaths))
	for _, importPath := range importPaths {
		if included, i := include(dependencyOrders, importPath); included {
			dependencies = append(dependencies, dependency{filePath: path, fileLayer: currentLayer, importPath: importPath, importLayer: i})
		}
	}
	return dependencies, nil
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
