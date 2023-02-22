package cdb

import (
	"database/sql"
	"github.com/joho/godotenv"
	gc "gopkg.in/check.v1"
	"links-r-us/linkgraph/graph/graphtest"
	"os"
	"testing"
)

var _ = gc.Suite(new(CockroachDbGraphTestSuite))

type CockroachDbGraphTestSuite struct {
	graphtest.SuiteBase
	db *sql.DB
}

// Register our test-suite with go test.
func Test(t *testing.T) { gc.TestingT(t) }

func (s *CockroachDbGraphTestSuite) SetUpSuite(c *gc.C) {
	err := godotenv.Load("../../../.env")
	dsn := os.Getenv("CDB_DSN")

	if dsn == "" || err != nil {
		c.Skip("Missing CDB_DSN env var; skipping cockroachdb-backed graph test suite")
	}

	g, err := NewCockroachDbGraph(dsn)
	c.Assert(err, gc.IsNil)
	s.SetGraph(g)
	s.db = g.db
}

func (s *CockroachDbGraphTestSuite) SetUpTest(c *gc.C) { s.flushDB(c) }

func (s *CockroachDbGraphTestSuite) TearDownSuite(c *gc.C) {
	if s.db != nil {
		s.flushDB(c)
		c.Assert(s.db.Close(), gc.IsNil)
	}
}

func (s *CockroachDbGraphTestSuite) flushDB(c *gc.C) {
	_, err := s.db.Exec("DELETE FROM links")
	c.Assert(err, gc.IsNil)
	_, err = s.db.Exec("DELETE FROM edges")
	c.Assert(err, gc.IsNil)
}
