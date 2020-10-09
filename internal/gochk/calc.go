package gochk

import (
	"os"
	"path/filepath"
)

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

func retrieveIndices(importPaths []string, dependencies []string, path string, currentLayer int) []dependency {
	layers := make([]dependency, 0, 10)
	for _, importPath := range importPaths {
		if included, i := include(dependencies, importPath); included {
			layers = append(layers, dependency{
				filepath:     path,
				currentLayer: currentLayer,
				path:         importPath,
				index:        i,
			})
		}
	}
	return layers
}

func retrieveViolates(currentLayer int, importLayers []dependency) []dependency {
	violates := make([]dependency, 0, len(importLayers))
	for _, d := range importLayers {
		if d.index < currentLayer {
			violates = append(violates, d)
			continue
		}
	}
	return violates
}
