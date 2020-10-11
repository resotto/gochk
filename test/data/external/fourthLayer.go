package external

import (
	"github.com/resotto/gochk/test/data/adapter"
	"github.com/resotto/gochk/test/data/application"
	"github.com/resotto/gochk/test/data/domain"
)

// Fourth is type for tests
type Fourth int

var (
	first  = domain.First{}
	second = application.Second{}
	third  = adapter.Third{}
)
