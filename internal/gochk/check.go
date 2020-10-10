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

// Check ensures that the direction of dependencies is correct
func Check(cfg Config) {
	errorDeps, err := walkFiles(cfg)
	if err != nil {
		panic(err)
	}
	if len(errorDeps) > 0 {
		for _, d := range errorDeps {
			printError(d.filepath, d.path, cfg.DependencyOrders, d.currentLayer, d.index)
		}
		panic("Dependencies which violate dependency orders found!")
	}
}

func walkFiles(cfg Config) ([]dependency, error) {
	errorDeps := make([]dependency, 0, 0)
	return errorDeps, filepath.Walk(cfg.TargetPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}
		if ignored, ignoreError := matchIgnore(cfg.Ignore, path, info); ignored {
			return ignoreError
		}
		errorDeps = append(errorDeps, (checkDependency(cfg.DependencyOrders, path))...)
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

func checkDependency(dependencies []string, path string) []dependency {
	_, currentLayer := include(dependencies, path)
	importLayers := retrieveLayers(dependencies, path, currentLayer)
	if len(importLayers) == 0 {
		printNone(path)
		return nil
	}
	if violates := retrieveViolates(currentLayer, importLayers); len(violates) > 0 {
		return violates
	}
	printVerified(path)
	return nil
}

func include(strs []string, s string) (bool, int) {
	for i, v := range strs {
		if strings.Contains(s, v) {
			return true, i
		}
	}
	return false, -1
}
