package performance

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/resotto/gochk/internal/gochk"
)

const (
	externalDirPath           = "./external"
	adapterDirPath            = "./adapter"
	postgresqlDirPath         = "./adapter/postgresql"
	modelDirPath              = "./adapter/postgresql/model"
	adapterRepositoryDirPath  = "./adapter/repository"
	adapterServiceDirPath     = "./adapter/service"
	viewDirPath               = "./adapter/view"
	applicationDirPath        = "./application"
	applicationServiceDirPath = "./application/service"
	usecaseDirPath            = "./application/usecase"
	domainDirPath             = "./domain"
	factoryDirPath            = "./domain/factory"
	domainRepositoryDirPath   = "./domain/repository"
	valueobjectDirPath        = "./domain/valueobject"

	externalTxtPath    = "../testdata/external.txt"
	adapterTxtPath     = "../testdata/adapter.txt"
	applicationTxtPath = "../testdata/application.txt"
	domainTxtPath      = "../testdata/domain.txt"
)

var dependencyOrders = []string{"external", "adapter", "application", "domain"}

func createFile(path string, contentsPath string) {
	contentsFilepath, _ := filepath.Abs(contentsPath)
	bytes, rerr := ioutil.ReadFile(contentsFilepath)
	if rerr != nil {
		panic(rerr)
	}
	filepath, _ := filepath.Abs(path)
	err := ioutil.WriteFile(filepath, bytes, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func mkdir(path string) {
	filepath, _ := filepath.Abs(path)
	err := os.Mkdir(filepath, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func rmdir(path string) {
	filepath, _ := filepath.Abs(path)
	err := os.RemoveAll(filepath)
	if err != nil {
		panic(err)
	}
}

func createDirs() {
	mkdir(externalDirPath)
	mkdir(adapterDirPath)
	mkdir(postgresqlDirPath)
	mkdir(modelDirPath)
	mkdir(adapterRepositoryDirPath)
	mkdir(adapterServiceDirPath)
	mkdir(viewDirPath)
	mkdir(applicationDirPath)
	mkdir(applicationServiceDirPath)
	mkdir(usecaseDirPath)
	mkdir(domainDirPath)
	mkdir(factoryDirPath)
	mkdir(domainRepositoryDirPath)
	mkdir(valueobjectDirPath)
}

func setup() {
	start := time.Now()
	createDirs()
	dirAndContents := [][]string{
		{externalDirPath, externalTxtPath},
		{adapterDirPath, adapterTxtPath},
		{applicationDirPath, applicationTxtPath},
		{domainDirPath, domainTxtPath},
	}
	for _, dc := range dirAndContents {
		for i := 0; i < 10000; i++ {
			createFile(dc[0]+"/g"+strconv.Itoa(i)+".go", dc[1])
		}
	}
	end := time.Now()
	fmt.Println("created all files for performance test in", end.Sub(start))
}

func teardown() {
	start := time.Now()
	rmdir(externalDirPath)
	rmdir(adapterDirPath)
	rmdir(applicationDirPath)
	rmdir(domainDirPath)
	end := time.Now()
	fmt.Println("deleted all files for performance test recursively in", end.Sub(start))
}

func TestMain(m *testing.M) {
	setup()
	result := m.Run()
	teardown()
	os.Exit(result)
}

func TestCheckPerformance(t *testing.T) {
	tests := []struct {
		name       string
		targetPath string
		cfg        gochk.Config
		expected   string
	}{
		{
			"Check() performance test",
			"../performance/",
			gochk.Config{
				DependencyOrders:           dependencyOrders,
				Ignore:                     []string{"test"},
				PrintViolationsAtTheBottom: true,
			},
			"99m",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			start := time.Now()
			gochk.Check(tt.targetPath, tt.cfg)
			end := time.Now()
			diff := end.Sub(start)
			expected, _ := time.ParseDuration(tt.expected)
			if diff > expected {
				t.Errorf("got %s, want %s", diff.String(), expected.String())
			}
		})
	}
}
