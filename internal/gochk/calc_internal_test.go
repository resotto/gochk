package gochk

import (
	"go/parser"
	"go/token"
	"path/filepath"
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
		CheckResults     []CheckResult
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
				newViolated("this message is not tested"),
				newViolated("this message is not tested"),
			},
		},
		{
			"file which cannot be opened",
			[]CheckResult{},
			dependencyOrders,
			lockedPath,
			[]CheckResult{newWarning("this message is not tested")},
		},
		{
			"file which has no imports",
			[]CheckResult{},
			dependencyOrders,
			underscoreTestPath,
			[]CheckResult{newNone("this message is not tested")},
		},
		{
			"file which has verified dependency",
			[]CheckResult{},
			dependencyOrders,
			fourthLayerPath,
			[]CheckResult{newVerified("this message is not tested")},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			filepath, err := filepath.Abs(tt.path)
			if err != nil {
				return
			}
			n, err := parser.ParseFile(token.NewFileSet(), filepath, nil, parser.ImportsOnly)
			if err != nil {
				return
			}
			setResultType(&tt.CheckResults, tt.dependencyOrders, tt.path, n.Imports)
			if len(tt.CheckResults) != len(tt.expected) {
				t.Errorf("got %d, want %d", len(tt.CheckResults), len(tt.expected))
			}
			for i, r := range tt.CheckResults {
				if r.resultType != tt.expected[i].resultType {
					t.Errorf("got %s, want %s", r.resultType, tt.expected[i].resultType)
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
				newViolated("this message is not tested"),
				newViolated("this message is not tested"),
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
