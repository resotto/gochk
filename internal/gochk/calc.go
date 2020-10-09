package gochk

import (
	"os"
	"path/filepath"
	"strings"
)

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
