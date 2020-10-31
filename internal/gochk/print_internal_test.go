package gochk

import (
	"strings"
	"testing"
)

func TestShow(t *testing.T) {
	testNames := []string{
		"violated and don't printViolationsAtTheBottom",
		"not violated and printViolationsAtTheBottom",
	}
	tests := []struct {
		name                       string
		results                    []CheckResult
		violated                   bool
		printViolationsAtTheBottom bool
	}{
		{
			testNames[1],
			[]CheckResult{
				newNone("test message"),
				newNone("test message"),
			},
			false,
			true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			defer func() {
				err := recover()
				if err != nil && !strings.EqualFold(tt.name, testNames[0]) {
					t.Errorf("%s shouldn't create panic", tt.name)
				}
			}()
			Show(tt.results, tt.violated, tt.printViolationsAtTheBottom)
		})
	}
}

func TestPrintConcurrently(t *testing.T) {
	testName := "printConcurrently() test"
	tests := []struct {
		name string
		crs  []CheckResult
	}{
		{
			testName,
			[]CheckResult{
				newViolated("test message"),
				newViolated("test message"),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			defer func() {
				err := recover()
				if err != nil && strings.EqualFold(tt.name, testName) {
					t.Errorf("%s shouldn't create panic", tt.name)
				}
			}()
			printSequentially(tt.crs)
		})
	}
}

func TestPrintSequentially(t *testing.T) {
	testName := "printSequentially() test"
	tests := []struct {
		name string
		crs  []CheckResult
	}{
		{
			testName,
			[]CheckResult{
				newViolated("test message"),
				newViolated("test message"),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			defer func() {
				err := recover()
				if err != nil && strings.EqualFold(tt.name, testName) {
					t.Errorf("%s shouldn't create panic", tt.name)
				}
			}()
			printSequentially(tt.crs)
		})
	}
}

func TestPrintColorMessage(t *testing.T) {
	testName := "printColorMessage() test"
	tests := []struct {
		name string
		cr   CheckResult
	}{
		{
			testName,
			newNone("test message"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			defer func() {
				err := recover()
				if err != nil && strings.EqualFold(tt.name, testName) {
					t.Errorf("%s shouldn't create panic", tt.name)
				}
			}()
			printColorMessage(tt.cr)
		})
	}
}

func TestPrintAA(t *testing.T) {
	testName := "printAA() test"
	tests := []struct {
		name string
	}{
		{
			testName,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			defer func() {
				err := recover()
				if err != nil && strings.EqualFold(tt.name, testName) {
					t.Errorf("%s shouldn't create panic", tt.name)
				}
			}()
			printAA()
		})
	}
}
