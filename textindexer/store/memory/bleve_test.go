package memory

import (
	gc "gopkg.in/check.v1"
	"links-r-us/textindexer/index/indextest"
	"testing"
)

var _ = gc.Suite(new(InMemoryBleveIndexer))

// Register our test-suit with go test.
func Test(t *testing.T) { gc.TestingT(t) }

type InMemoryBleveTestSuit struct {
	indextest.SuiteBase
	idx *InMemoryBleveIndexer
}

func (s *InMemoryBleveTestSuit) SetUpTest(c *gc.C) {
	idx, err := NewInMemoryBleveIndexer()
	c.Assert(err, gc.IsNil)
	s.SetIndexer(idx)
	s.idx = idx
}

func (s *InMemoryBleveIndexer) TearDownTest(c *gc.C) {
	c.Assert(s.idx.Close(), gc.IsNil)
}
