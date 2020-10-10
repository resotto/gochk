package gochk

import "strings"

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

func retrieveViolations(currentLayer int, importLayers []dependency) []dependency {
	violations := make([]dependency, 0, len(importLayers))
	for _, d := range importLayers {
		if d.index < currentLayer {
			violations = append(violations, d)
		}
	}
	return violations
}

func retrieveImportPath(line string) string {
	firstQuoIndex := strings.Index(line, "\"")
	if firstQuoIndex == -1 {
		return ""
	}
	return line[firstQuoIndex:]
}

func include(strs []string, s string) (bool, int) {
	for i, v := range strs {
		if strings.Contains(s, v) {
			return true, i
		}
	}
	return false, -1
}
