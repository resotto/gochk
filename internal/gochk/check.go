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

// Check makes sure that the direction of dependencies is correct
func Check(cfg Config) {
	errorDeps, err := walkFiles(cfg)
	if err != nil {
		panic(err)
	}
	if len(errorDeps) > 0 {
		for _, d := range errorDeps {
			printError(d.filepath, d.path, cfg.DependencyOrders, d.currentLayer, d.index)
		}
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
		tempDeps := checkDependency(cfg.DependencyOrders, path)
		errorDeps = append(errorDeps, tempDeps...)
		return nil
	})
}

func matchIgnore(ignorePaths []string, path string, info os.FileInfo) (bool, error) {
	if include(ignorePaths, path) {
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

func include(strs []string, s string) bool {
	for _, v := range strs {
		if strings.Contains(s, v) {
			return true
		}
	}
	return false
}
