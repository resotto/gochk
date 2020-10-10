package gochk

import "testing"

const (
	pathA       = "./a/test.go"
	pathB       = "./b/test.go"
	pathC       = "./c/test.go"
	currentPath = "./c/path.go"
)

func TestRetrieveIndices(t *testing.T) {
	tests := []struct {
		name         string
		importPaths  []string
		dependencies []string
		path         string
		currentLayer int
		expected     []dependency
	}{
		{
			"third layer file in three layers",
			[]string{pathA, pathB, pathC},
			[]string{"c", "b", "a"},
			currentPath,
			0,
			[]dependency{
				dependency{filepath: currentPath, currentLayer: 0, path: pathA, index: 2},
				dependency{filepath: currentPath, currentLayer: 0, path: pathB, index: 1},
				dependency{filepath: currentPath, currentLayer: 0, path: pathC, index: 0},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			results := retrieveIndices(tt.importPaths, tt.dependencies, tt.path, tt.currentLayer)
			for i, d := range results {
				if d.filepath != tt.path {
					t.Errorf("got %s, want %s", d.filepath, tt.path)
				}
				if d.currentLayer != tt.currentLayer {
					t.Errorf("got %d, want %d", d.currentLayer, tt.currentLayer)
				}
				if d.path != tt.expected[i].path {
					t.Errorf("got %s, want %s", d.path, tt.expected[i].path)
				}
				if d.index != tt.expected[i].index {
					t.Errorf("got %d, want %d", d.index, tt.expected[i].index)
				}
			}
		})
	}
}

func TestRetrieveViolations(t *testing.T) {
	tests := []struct {
		name         string
		currentLayer int
		importLayers []dependency
		expected     []dependency
	}{
		{
			"two violations at first layer",
			2,
			[]dependency{
				dependency{filepath: currentPath, currentLayer: 2, path: pathC, index: 0},
				dependency{filepath: currentPath, currentLayer: 2, path: pathB, index: 1},
				dependency{filepath: currentPath, currentLayer: 2, path: pathA, index: 2},
			},
			[]dependency{
				dependency{filepath: currentPath, currentLayer: 2, path: pathC, index: 0},
				dependency{filepath: currentPath, currentLayer: 2, path: pathB, index: 1},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			results := retrieveViolations(tt.currentLayer, tt.importLayers)
			if len(results) != 2 {
				t.Errorf("got %d, want %d", len(results), len(tt.expected))
			}
			for i, d := range results {
				if d.currentLayer != tt.currentLayer {
					t.Errorf("got %d, want %d", d.currentLayer, tt.currentLayer)
				}
				if d.path != tt.expected[i].path {
					t.Errorf("got %s, want %s", d.path, tt.expected[i].path)
				}
				if d.index != tt.expected[i].index {
					t.Errorf("got %d, want %d", d.index, tt.expected[i].index)
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
