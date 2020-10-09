package gochk

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type dependency struct {
	filepath     string
	currentLayer int
	path         string
	index        int
}

// Check makes sure the direction of dependency is correct
func Check(cfg Config) {
	errorDeps := make([]dependency, 0, 0)
	err := filepath.Walk(cfg.TargetPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}
		if include(cfg.Ignore, path) { // todo
			if info.IsDir() {
				printIgnored(path)
				return filepath.SkipDir
			}
			printIgnored(path)
			return nil
		}
		if info.IsDir() || !strings.Contains(info.Name(), ".go") {
			return nil
		}
		tempDeps := checkDependency(cfg.DependencyOrders, path)
		errorDeps = append(errorDeps, tempDeps...)
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	if len(errorDeps) > 0 {
		for _, d := range errorDeps {
			printError(d.filepath, d.path, cfg.DependencyOrders, d.currentLayer, d.index)
		}
	}
}

func checkDependency(dependencies []string, path string) []dependency {
	currentLayer := search(dependencies, path)
	importLayers := retrieveLayers(dependencies, path, currentLayer)

	if len(importLayers) == 0 {
		printNone(path)
		return nil
	}
	redDeps := make([]dependency, 0, len(importLayers))

	for _, d := range importLayers {
		if d.index < currentLayer {
			redDeps = append(redDeps, d)
			continue
		}
	}
	if len(redDeps) > 0 {
		return redDeps
	}
	printVerified(path)
	return nil
}

func retrieveLayers(dependencies []string, path string, currentLayer int) []dependency {
	layers := make([]dependency, 0, 10)
	filepath, _ := filepath.Abs(path)
	f, err := os.Open(filepath)
	defer f.Close()
	if err != nil {
		printWarning(filepath)
		return layers
	}
	imports := readImports(f)

	for _, v := range imports {
		l := search(dependencies, v)
		if l != -1 {
			layers = append(layers, dependency{
				filepath:     path,
				currentLayer: currentLayer,
				path:         v,
				index:        l,
			})
		}
	}
	return layers
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
