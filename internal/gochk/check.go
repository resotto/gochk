package gochk

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	red    = "\033[1;31m%s\033[0m"
	yellow = "\033[1;33m%s\033[0m"
	green  = "\033[1;32m%s\033[0m"
	teal   = "\033[1;36m%s\033[0m"
)

type dependency struct {
	path  string
	index int
}

// Check makes sure the direction of dependency is correct
func Check(cfg Config) {
	err := filepath.Walk(cfg.TargetPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}
		if include(cfg.Ignore, path) { // todo
			if info.IsDir() {
				print(yellow, "[Ignored]  "+path)
				return filepath.SkipDir
			}
			print(yellow, "[Ignored]  "+path)
			return nil
		}
		if info.IsDir() || !strings.Contains(info.Name(), ".go") {
			return nil
		}
		checkDependency(cfg.DependencyOrders, path)
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}

func checkDependency(dependencies []string, path string) {
	currentLayer := search(dependencies, path)
	importLayers := retrieveLayers(dependencies, path)

	if len(importLayers) == 0 {
		print(teal, "[None]     "+path)
		return
	}
	redDeps := make([]dependency, 0, len(importLayers))

	for _, d := range importLayers {
		if d.index < currentLayer {
			redDeps = append(redDeps, d)
			continue
		}
	}
	if len(redDeps) > 0 {
		for _, d := range redDeps {
			print(red, "[Error]    "+path+" imports "+d.path)
			print(red, "           \""+dependencies[currentLayer]+"\" depends on \""+dependencies[d.index]+"\"")
		}
	} else {
		print(green, "[Verified] "+path)
	}
}

func retrieveLayers(dependencies []string, path string) []dependency {
	filepath, _ := filepath.Abs(path)
	imports := readImports(filepath)
	layers := make([]dependency, 0, len(imports))

	for _, v := range imports {
		l := search(dependencies, v)
		if l != -1 {
			layers = append(layers, dependency{
				path:  v,
				index: l,
			})
		}
	}
	return layers
}

func readImports(filepath string) []string {
	f, _ := os.Open(filepath)
	defer f.Close()
	scanner := bufio.NewScanner(f)
	multipleImport := false
	imports := make([]string, 0, 10)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "import") {
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

func search(strs []string, elm string) int {
	for i, v := range strs {
		if strings.Contains(elm, v) {
			return i
		}
	}
	return -1
}

func include(strs []string, elm string) bool {
	for _, v := range strs {
		if strings.Contains(elm, v) {
			return true
		}
	}
	return false
}

func print(color string, message string) {
	fmt.Printf(color, message)
	fmt.Println()
}
