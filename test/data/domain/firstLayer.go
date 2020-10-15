package domain

import (
	"github.com/resotto/gochk/test/data/adapter"
	"github.com/resotto/gochk/test/data/application"
)

// First is type for tests
type First struct{}

var (
	second = application.Second{}
	third  = adapter.Third{}
)
