// Laravel equivalent: tests/Feature/IndexTest.php
package feature

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"goravel/tests"
)

type IndexTestSuite struct {
	suite.Suite
	tests.TestCase
}

func TestIndexTestSuite(t *testing.T) {
	suite.Run(t, new(IndexTestSuite))
}

func (s *IndexTestSuite) TestRootListsEndpoints() {
	response, err := s.Http(s.T()).Get("/")
	s.Require().NoError(err)
	response.AssertOk()

	var body struct {
		Endpoints []struct{ Method, Path string }
	}
	s.Require().NoError(response.Bind(&body))

	s.Contains(body.Endpoints, struct{ Method, Path string }{"GET", "/health"})
	s.Contains(body.Endpoints, struct{ Method, Path string }{"DELETE", "/api/clubs/{id}"})
	s.Contains(body.Endpoints, struct{ Method, Path string }{"POST", "/api/fixtures/{id}/result"})
	s.NotContains(body.Endpoints, struct{ Method, Path string }{"STATIC", "/public"})
}

func (s *IndexTestSuite) TestBrowsersGetAnHtmlPageWithTheLogo() {
	response, err := s.Http(s.T()).WithHeader("Accept", "text/html").Get("/")
	s.Require().NoError(err)

	response.AssertOk()
	response.AssertSee([]string{"goravel.dev/logo.png", "/api/league/table"})
}
