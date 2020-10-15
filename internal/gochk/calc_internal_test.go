package gochk

import (
	"testing"
)

const (
	pathA       = "./a/test.go"
	pathB       = "./b/test.go"
	pathC       = "./c/test.go"
	currentPath = "./c/path.go"
)

func TestSetResultType(t *testing.T) {
	tests := []struct {
		name             string
		checkResults     []CheckResult
		dependencyOrders []string
		path             string
		expected         []CheckResult
	}{
		{
			"first layer file which violates dependencies",
			[]CheckResult{},
			dependencyOrders,
			firstLayerPath,
			[]CheckResult{
				CheckResult{resultType: violated, message: "not tested", color: red},
				CheckResult{resultType: violated, message: "not tested", color: red},
			},
		},
		{
			"file which cannot be opened",
			[]CheckResult{},
			dependencyOrders,
			lockedPath,
			[]CheckResult{
				CheckResult{resultType: warning, message: "not tested", color: purple},
			},
		},
		{
			"file which has no imports",
			[]CheckResult{},
			dependencyOrders,
			underscoreTestPath,
			[]CheckResult{
				CheckResult{resultType: none, message: "not tested", color: teal},
			},
		},
		{
			"file which has verified dependency",
			[]CheckResult{},
			dependencyOrders,
			fourthLayerPath,
			[]CheckResult{
				CheckResult{resultType: verified, message: "not tested", color: green},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			setResultType(&tt.checkResults, tt.dependencyOrders, tt.path)
			if len(tt.checkResults) != len(tt.expected) {
				t.Errorf("got %d, want %d", len(tt.checkResults), len(tt.expected))
			}
			for i, r := range tt.checkResults {
				if r.resultType != tt.expected[i].resultType {
					t.Errorf("got %s, want %s", r.resultType, tt.expected[i].resultType)
				}
				if r.color != tt.expected[i].color {
					t.Errorf("got %s, want %s", r.color, tt.expected[i].color)
				}
			}
		})
	}
}

func TestRetrieveViolations(t *testing.T) {
	tests := []struct {
		name             string
		dependencyOrders []string
		currentLayer     int
		dependencies     []dependency
		expected         []CheckResult
	}{
		{
			"two violations at first layer",
			dependencyOrders,
			2,
			[]dependency{
				dependency{filePath: currentPath, fileLayer: 2, importPath: pathC, importLayer: 0},
				dependency{filePath: currentPath, fileLayer: 2, importPath: pathB, importLayer: 1},
				dependency{filePath: currentPath, fileLayer: 2, importPath: pathA, importLayer: 2},
			},
			[]CheckResult{
				CheckResult{resultType: violated, message: "not tested", color: red},
				CheckResult{resultType: violated, message: "not tested", color: red},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			results := retrieveViolations(tt.dependencyOrders, tt.currentLayer, tt.dependencies)
			if len(results) != len(tt.expected) {
				t.Errorf("got %d, want %d", len(results), len(tt.expected))
			}
			for i, r := range results {
				if r.resultType != tt.expected[i].resultType {
					t.Errorf("got %s, want %s", r.resultType, tt.expected[i].resultType)
				}
				if r.color != tt.expected[i].color {
					t.Errorf("got %s, want %s", r.color, tt.expected[i].color)
				}
			}
		})
	}
}

func TestRetrieveImportPath(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected string
	}{
		{
			"import path exists",
			"import \"" + pathA + "\"",
			"\"" + pathA + "\"",
		},
		{
			"import path doesn't exist",
			"(nothing)",
			"",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := retrieveImportPath(tt.line)
			if result != tt.expected {
				t.Errorf("got %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestInclude(t *testing.T) {
	type result struct {
		found bool
		index int
	}
	tests := []struct {
		name     string
		strs     []string
		s        string
		expected result
	}{
		{
			"str included",
			[]string{"a", "b", "c"},
			pathC,
			result{found: true, index: 2},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			found, index := include(tt.strs, tt.s)
			if found != tt.expected.found {
				t.Errorf("got %t, want %t", found, tt.expected.found)
			}
			if index != tt.expected.index {
				t.Errorf("got %d, want %d", index, tt.expected.index)
			}
		})
	}
}
