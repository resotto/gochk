package domain

import (
	"github.com/resotto/gochk/test/testdata/adapter"
	"github.com/resotto/gochk/test/testdata/application"
)

// First is type for tests
type First struct{}

var (
	second = application.Second{}
	third  = adapter.Third{}
)
