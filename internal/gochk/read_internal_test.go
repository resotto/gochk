package gochk

import (
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	pfmt     = "\"fmt\""
	pstrings = "\"strings\""

	testConfigPath      = "../../test/testdata/test.json"
	testDirPath         = "../../test/testdata/test"
	underscoreTestPath  = "../../test/testdata/_test.go"
	blockCommentsPath   = "../../test/testdata/blockComments.go"
	lockfilePath        = "../../test/testdata/lockfile.txt"
	lockedPath          = "../../test/testdata/adapter/locked.go"
	violatefilePath     = "../../test/testdata/violatefile.txt"
	violatePath         = "../../test/testdata/application/violate.go"
	multipleImportsPath = "../../test/testdata/multipleImports.go"
	oneImportPath       = "../../test/testdata/oneImport.go"
	singleCommentsPath  = "../../test/testdata/singleComments.go"
	firstLayerPath      = "../../test/testdata/domain/firstLayer.go"
	secondLayerPath     = "../../test/testdata/application/secondLayer.go"
	fourthLayerPath     = "../../test/testdata/external/fourthLayer.go"

	adapterPkgPath     = "\"github.com/resotto/gochk/test/testdata/adapter\""
	applicationPkgPath = "\"github.com/resotto/gochk/test/testdata/application\""
	domainPkgPath      = "\"github.com/resotto/gochk/test/testdata/domain\""
)

var (
	ignorePaths      = []string{"test", "_test"}
	dependencyOrders = []string{"external", "adapter", "application", "domain"}
)

func createFile(path string, contentsPath string, permission os.FileMode) string {
	contentsFilepath, _ := filepath.Abs(contentsPath)
	bytes, rerr := ioutil.ReadFile(contentsFilepath)
	if rerr != nil {
		panic(rerr)
	}
	filepath, _ := filepath.Abs(path)
	err := ioutil.WriteFile(filepath, bytes, permission)
	if err != nil {
		panic(err)
	}
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
	createFile(violatePath, violatefilePath, os.ModePerm)
}

func teardown() {
	rmdir(testDirPath)
	removeFile(lockedPath)
	removeFile(violatePath)
}

func TestMain(m *testing.M) {
	setup()
	result := m.Run()
	teardown()
	os.Exit(result)
}

func TestParseConfig(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected Config
	}{
		{
			"succeeded in parsing config",
			testConfigPath,
			Config{
				DependencyOrders:           []string{"external", "adapter", "application", "domain"},
				Ignore:                     []string{"test", ".git"},
				PrintViolationsAtTheBottom: false,
			},
		},
		{
			"failed at parsing config",
			testDirPath + "/test.json",
			Config{
				DependencyOrders:           []string{"external", "adapter", "application", "domain"},
				Ignore:                     []string{"test", ".git"},
				PrintViolationsAtTheBottom: false,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			defer func() {
				err := recover()
				if err != nil && !strings.Contains(tt.name, "failed at parsing config") {
					t.Errorf("got %s, want %s", tt.name, "\"failed at parsing config\" should create panic")
				}
			}()
			result := ParseConfig(tt.path)
			for i, r := range result.DependencyOrders {
				if r != tt.expected.DependencyOrders[i] {
					t.Errorf("got %s, want %s", r, tt.expected.DependencyOrders[i])
				}
			}
			for i, r := range result.Ignore {
				if r != tt.expected.Ignore[i] {
					t.Errorf("got %s, want %s", r, tt.expected.Ignore[i])
				}
			}
			if result.PrintViolationsAtTheBottom != tt.expected.PrintViolationsAtTheBottom {
				t.Errorf("got %t, want %t", result.PrintViolationsAtTheBottom, tt.expected.PrintViolationsAtTheBottom)
			}
		})
	}
}

func TestCheck(t *testing.T) {
	tests := []struct {
		name       string
		targetPath string
		cfg        Config
		expected   []CheckResult
	}{
		{
			"violation found",
			violatePath,
			Config{
				DependencyOrders:           dependencyOrders,
				Ignore:                     []string{},
				PrintViolationsAtTheBottom: true,
			},
			[]CheckResult{newViolated("this message is not tested")},
		},
		{
			"no results",
			violatefilePath,
			Config{
				DependencyOrders:           dependencyOrders,
				Ignore:                     []string{},
				PrintViolationsAtTheBottom: true,
			},
			[]CheckResult{},
		},
		{
			"ignored",
			testDirPath,
			Config{
				DependencyOrders:           dependencyOrders,
				Ignore:                     []string{"test"},
				PrintViolationsAtTheBottom: true,
			},
			[]CheckResult{newIgnored("this message is not tested")},
		},
		{
			"warning for the path which doesn't exist",
			testDirPath + "/none.go",
			Config{
				DependencyOrders:           dependencyOrders,
				Ignore:                     []string{},
				PrintViolationsAtTheBottom: true,
			},
			[]CheckResult{newWarning("this message is not tested")},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			results, _ := Check(tt.targetPath, tt.cfg)
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

func TestRetrieveDependencies(t *testing.T) {
	tests := []struct {
		name         string
		dependencies []string
		path         string
		currentLayer int
		expected     []dependency
	}{
		{
			"one layer retrieval",
			dependencyOrders,
			fourthLayerPath,
			0,
			[]dependency{dependency{fourthLayerPath, 0, domainPkgPath, 3}},
		},
		{
			"include file which couldn't be opened",
			dependencyOrders,
			lockedPath,
			1,
			[]dependency{},
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
			results, _ := retrieveDependencies(tt.dependencies, tt.path, tt.currentLayer, n.Imports)
			for i, r := range results {
				if r.filePath != tt.expected[i].filePath {
					t.Errorf("got %s, want %s", r.filePath, tt.expected[i].filePath)
				}
				if r.fileLayer != tt.expected[i].fileLayer {
					t.Errorf("got %d, want %d", r.fileLayer, tt.expected[i].fileLayer)
				}
				if r.importPath != tt.expected[i].importPath {
					t.Errorf("got %s, want %s", r.importPath, tt.expected[i].importPath)
				}
				if r.importLayer != tt.expected[i].importLayer {
					t.Errorf("got %d, want %d", r.importLayer, tt.expected[i].importLayer)
				}
			}
		})
	}
}
