package es

import (
	"github.com/joho/godotenv"
	"links-r-us/textindexer/index/indextest"
	"os"
	"strings"
	"testing"

	gc "gopkg.in/check.v1"
)

var _ = gc.Suite(new(ElasticSearchTestSuite))

func Test(t *testing.T) { gc.TestingT(t) }

type ElasticSearchTestSuite struct {
	indextest.SuiteBase
	idx *ElasticSearchIndexer
}

func (s *ElasticSearchTestSuite) SetUpSuite(c *gc.C) {
	err := godotenv.Load("../../../.env")

	certPath := os.Getenv("ES_CERT_PATH")
	nodeList := os.Getenv("ES_NODES")
	apiKey := os.Getenv("ES_API_KEY")

	if nodeList == "" || certPath == "" || apiKey == "" {
		c.Skip("Missing one of envvar; skipping elasticsearch-backed index test suite")
	}

	idx, err := NewElasticSearchIndexer(strings.Split(nodeList, ","), certPath, apiKey, true)
	c.Assert(err, gc.IsNil)
	s.SetIndexer(idx)
	s.idx = idx
}

func (s *ElasticSearchTestSuite) SetUpTest(c *gc.C) {
	if s.idx.es != nil {
		_, err := s.idx.es.Indices.Delete([]string{indexName})
		c.Assert(err, gc.IsNil)
		err = ensureIndex(s.idx.es)
		c.Assert(err, gc.IsNil)
	}
}
