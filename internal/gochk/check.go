package gochk

type dependency struct {
	filepath     string
	currentLayer int
	path         string
	index        int
}

// Check ensures that the direction of dependencies is correct
func Check(cfg Config) {
	violations, err := walkFiles(cfg)
	if err != nil {
		panic(err)
	}
	if len(violations) > 0 {
		for _, d := range violations {
			printError(d.filepath, d.path, cfg.DependencyOrders, d.currentLayer, d.index)
		}
		panic("Dependencies which violate dependency orders found!")
	}
}

func checkDependency(dependencies []string, path string) []dependency {
	_, currentLayer := include(dependencies, path)
	importLayers := retrieveLayers(dependencies, path, currentLayer)
	if len(importLayers) == 0 {
		printNone(path)
		return nil
	}
	if violations := retrieveViolations(currentLayer, importLayers); len(violations) > 0 {
		return violations
	}
	printVerified(path)
	return nil
}
