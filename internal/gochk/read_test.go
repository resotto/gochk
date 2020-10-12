package gochk

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	pfmt     = "\"fmt\""
	pstrings = "\"strings\""

	testDirPath         = "../../test/data/test"
	underscoreTestPath  = "../../test/data/_test.go"
	blockCommentsPath   = "../../test/data/blockComments.go"
	lockfilePath        = "../../test/data/lockfile.txt"
	lockedPath          = "../../test/data/adapter/locked.go"
	multipleImportsPath = "../../test/data/multipleImports.go"
	oneImportPath       = "../../test/data/oneImport.go"
	singleCommentsPath  = "../../test/data/singleComments.go"
	fourthLayerPath     = "../../test/data/external/fourthLayer.go"

	adapterPkgPath     = "\"github.com/resotto/gochk/test/data/adapter\""
	applicationPkgPath = "\"github.com/resotto/gochk/test/data/application\""
	domainPkgPath      = "\"github.com/resotto/gochk/test/data/domain\""
)

var (
	ignorePaths  = []string{"test", "_test"}
	dependencies = []string{"external", "adapter", "application", "domain"}
)

func createFile(path string, contentsPath string, permission os.FileMode) string {
	contentsFilepath, _ := filepath.Abs(contentsPath)
	bytes, rerr := ioutil.ReadFile(contentsFilepath)
	if rerr != nil {
		panic(rerr)
	}
	filepath, _ := filepath.Abs(path)
	ioutil.WriteFile(filepath, bytes, 0700)
	os.Chmod(path, permission)
	return filepath
}

func removeFile(path string) {
	filepath, _ := filepath.Abs(path)
	err := os.Remove(filepath)
	if err != nil {
		panic(err)
	}
}

func mkdir(path string, permission os.FileMode) {
	filepath, _ := filepath.Abs(path)
	err := os.Mkdir(filepath, permission)
	if err != nil {
		panic(err)
	}
}

func rmdir(path string) {
	filepath, _ := filepath.Abs(path)
	err := os.Remove(filepath)
	if err != nil {
		panic(err)
	}
}

func setup() {
	mkdir(testDirPath, os.ModePerm)
	createFile(lockedPath, lockfilePath, 0300)
}

func teardown() {
	rmdir(testDirPath)
	removeFile(lockedPath)
}

func TestMain(m *testing.M) {
	setup()
	result := m.Run()
	teardown()
	os.Exit(result)
}

func TestMatchIgnore(t *testing.T) {
	type result struct {
		matched bool
		err     error
	}
	tests := []struct {
		name        string
		ignorePaths []string
		targetPath  string
		expected    result
	}{
		{
			"test dir",
			ignorePaths,
			testDirPath,
			result{matched: true, err: filepath.SkipDir},
		},
		{
			"_test.go file",
			ignorePaths,
			underscoreTestPath,
			result{matched: true, err: nil},
		},
		{
			"file which is not included in ignore and is not .go",
			ignorePaths,
			lockfilePath,
			result{matched: true, err: nil},
		},
		{
			".go file which is not included in ignore",
			ignorePaths,
			"./print.go",
			result{matched: false, err: nil},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := filepath.Walk(tt.targetPath, func(path string, info os.FileInfo, err error) error {
				matched, err := matchIgnore(tt.ignorePaths, path, info)
				if matched != tt.expected.matched {
					fmt.Println(tt.name, path) // debug
					t.Errorf("got %t, want %t", matched, tt.expected.matched)
				}
				if err != tt.expected.err {
					t.Errorf("got %s, want %s", err, tt.expected.err)
				}
				return nil
			})
			if err != nil {
				panic(err)
			}
		})
	}
}

func TestRetrieveLayers(t *testing.T) {
	tests := []struct {
		name         string
		dependencies []string
		path         string
		currentLayer int
		expected     []dependency
	}{
		{
			"four layers retrieval",
			dependencies,
			fourthLayerPath,
			0,
			[]dependency{
				dependency{fourthLayerPath, 0, adapterPkgPath, 1},
				dependency{fourthLayerPath, 0, applicationPkgPath, 2},
				dependency{fourthLayerPath, 0, domainPkgPath, 3},
			},
		},
		{
			"include file which couldn't be opened",
			dependencies,
			lockedPath,
			1,
			[]dependency{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			results := retrieveLayers(tt.dependencies, tt.path, tt.currentLayer)
			fmt.Println(results)
			for i, r := range results {
				if r.filepath != tt.expected[i].filepath {
					t.Errorf("got %s, want %s", r.filepath, tt.expected[i].filepath)
				}
				if r.currentLayer != tt.expected[i].currentLayer {
					t.Errorf("got %d, want %d", r.currentLayer, tt.expected[i].currentLayer)
				}
				if r.path != tt.expected[i].path {
					t.Errorf("got %s, want %s", r.path, tt.expected[i].path)
				}
				if r.index != tt.expected[i].index {
					t.Errorf("got %d, want %d", r.index, tt.expected[i].index)
				}
			}
		})
	}
}

func TestReadImports(t *testing.T) {
	tests := []struct {
		name     string
		filepath string
		expected []string
	}{
		{
			"single import",
			oneImportPath,
			[]string{pfmt},
		},
		{
			"multiple import",
			multipleImportsPath,
			[]string{pfmt, pstrings},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			filepath, _ := filepath.Abs(tt.filepath)
			f, err := os.Open(filepath)
			defer f.Close()
			if err != nil {
				t.Errorf("couldn't open file: %s", filepath)
			}
			results := readImports(f)
			for i, importPath := range results {
				if importPath != tt.expected[i] {
					t.Errorf("got %s, want %s", importPath, tt.expected[i])
				}
			}
		})
	}
}

func TestSkipToImportStatement(t *testing.T) {
	tests := []struct {
		name     string
		filepath string
		expected bool
	}{
		{
			"with block comments",
			blockCommentsPath,
			true,
		},
		{
			"without block comments",
			singleCommentsPath,
			true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			filepath, _ := filepath.Abs(tt.filepath)
			f, err := os.Open(filepath)
			defer f.Close()
			if err != nil {
				t.Errorf("couldn't open file: %s", filepath)
				return
			}
			scanner := bufio.NewScanner(f)
			skipToImportStatement(scanner)
			scanner.Scan()
			if line := scanner.Text(); len(line) <= 6 || !strings.EqualFold(line[:6], "import") {
				t.Errorf("didn't skip to import statement: %s", line)
			}
		})
	}
}

func TestSkipBlockComments(t *testing.T) {
	tests := []struct {
		name     string
		filepath string
		expected bool
	}{
		{
			"block comments included and scanner.Scan() should return true",
			blockCommentsPath,
			true,
		},
		{
			"no block comments and scanner.Scan() should return true",
			oneImportPath,
			true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			filepath, _ := filepath.Abs(tt.filepath)
			f, err := os.Open(filepath)
			defer f.Close()
			if err != nil {
				t.Errorf("couldn't open file: %s", filepath)
				return
			}
			scanner := bufio.NewScanner(f)
			scanner.Scan()
			line := scanner.Text()
			if skipBlockComments(line, scanner); scanner.Scan() != tt.expected {
				t.Errorf("got %t, want %t", !tt.expected, tt.expected)
			}
		})
	}
}

func TestRetrieveMultipleImportPath(t *testing.T) {
	tests := []struct {
		name     string
		filepath string
		expected []string
	}{
		{
			"multiple import path",
			multipleImportsPath,
			[]string{pfmt, pstrings},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			filepath, _ := filepath.Abs(tt.filepath)
			f, err := os.Open(filepath)
			defer f.Close()
			if err != nil {
				t.Errorf("couldn't open file: %s", filepath)
				return
			}
			scanner := bufio.NewScanner(f)
			var line string
			for scanner.Scan() {
				if line = scanner.Text(); len(line) > 8 && strings.EqualFold(line[:8], "import (") {
					break
				}
			}
			results := retrieveMultipleImportPath(scanner, line)
			for i, r := range results {
				if r != tt.expected[i] {
					t.Errorf("got %s, want %s", r, tt.expected[i])
				}
			}
		})
	}
}
