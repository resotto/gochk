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
		if include(cfg.Ignore, path) {
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
}

func include(strs []string, elm string) bool {
	for _, v := range strs {
		if strings.Contains(elm, v) {
			return true
		}
	}
	return false
}
