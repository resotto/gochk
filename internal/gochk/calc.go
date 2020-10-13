package gochk

import "strings"

// judge path as none, verified or violated
func judgeResultType(dependencyOrders []string, path string) []CheckResult {
	results := make([]CheckResult, 0, 10)
	_, currentLayer := include(dependencyOrders, path)
	dependencies := retrieveDependencies(dependencyOrders, path, currentLayer)
	if len(dependencies) == 0 {
		results = append(results, CheckResult{resultType: none, message: path, color: teal})
		return results
	}
	if violations := retrieveViolations(dependencyOrders, currentLayer, dependencies); len(violations) > 0 {
		results = append(results, violations...)
		return results
	}
	results = append(results, CheckResult{resultType: verified, message: path, color: green})
	return results
}

func retrieveViolations(dependencyOrders []string, currentLayer int, dependencies []dependency) []CheckResult {
	violations := make([]CheckResult, 0, len(dependencies))
	for _, d := range dependencies {
		if d.importLayer < currentLayer {
			message := d.filePath + " imports " + d.importPath + "\n => " + dependencyOrders[d.fileLayer] + " depends on " + dependencyOrders[d.importLayer]
			violations = append(violations, CheckResult{resultType: violated, message: message, color: red})
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
