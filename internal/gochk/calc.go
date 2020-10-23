package gochk

import "strings"

func setResultType(results *[]CheckResult, dependencyOrders []string, path string) bool {
	_, currentLayer := include(dependencyOrders, path)
	dependencies, err := retrieveDependencies(dependencyOrders, path, currentLayer)
	if err != nil {
		*results = append([]CheckResult{newWarning(err.Error())}, *results...)
		return false
	}
	if len(dependencies) == 0 {
		*results = append([]CheckResult{newNone(path)}, *results...)
		return false
	}
	if violations := retrieveViolations(dependencyOrders, currentLayer, dependencies); len(violations) > 0 {
		*results = append(*results, violations...)
		return true
	}
	*results = append([]CheckResult{newVerified(path)}, *results...)
	return false
}

func retrieveViolations(dependencyOrders []string, currentLayer int, dependencies []dependency) []CheckResult {
	violations := make([]CheckResult, 0, len(dependencies))
	for _, d := range dependencies {
		if d.importLayer < currentLayer {
			message := d.filePath + " imports " + d.importPath + "\n => " + dependencyOrders[d.fileLayer] + " depends on " + dependencyOrders[d.importLayer]
			violations = append(violations, newViolated(message))
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
