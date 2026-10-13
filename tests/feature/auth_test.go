// Laravel equivalent: tests/Feature/AuthTest.php - see docs/chapters/19-authentication.md
package feature

import (
	"strings"
	"testing"

	testinghttp "github.com/goravel/framework/contracts/testing/http"
	"github.com/stretchr/testify/suite"

	"goravel/app/facades"
	"goravel/app/models"
	"goravel/tests"
)

type AuthTestSuite struct {
	suite.Suite
	tests.TestCase
	user models.User
}

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}

func (s *AuthTestSuite) SetupTest() {
	s.Require().NoError(s.TruncateTables())
	s.Require().NoError(facades.Orm().Factory().Create(&s.user, map[string]any{
		"Name":  "Alex Admin",
		"Email": "alex@example.com",
	}))
}

func (s *AuthTestSuite) login(email, password string) testinghttp.Response {
	response, err := s.Http(s.T()).WithHeader("Content-Type", "application/json").
		Post("/api/auth/login", strings.NewReader(`{"email":"`+email+`","password":"`+password+`"}`))
	s.Require().NoError(err)

	return response
}

// #region test
func (s *AuthTestSuite) TestLoginReturnsAToken() {
	response := s.login("alex@example.com", "password")
	response.AssertOk()

	var body struct {
		TokenType   string `json:"token_type"`
		AccessToken string `json:"access_token"`
	}
	s.Require().NoError(response.Bind(&body))
	s.Equal("Bearer", body.TokenType)
	s.NotEmpty(body.AccessToken)
}

func (s *AuthTestSuite) TestWrongPasswordAndUnknownEmailLookTheSame() {
	message := map[string]any{"message": "These credentials do not match our records."}

	s.login("alex@example.com", "wrong").AssertUnauthorized().AssertExactJson(message)
	s.login("nobody@example.com", "password").AssertUnauthorized().AssertExactJson(message)
}

func (s *AuthTestSuite) TestMeReturnsTheTokensUser() {
	response, err := s.LoginAs(s.T(), s.user).Get("/api/auth/me")
	s.Require().NoError(err)
	response.AssertOk()

	var body struct{ Data models.User }
	s.Require().NoError(response.Bind(&body))
	s.Equal(s.user.ID, body.Data.ID)
	s.Equal("alex@example.com", body.Data.Email)

	// json:"-" on the model keeps the hash out of every response
	content, err := response.Content()
	s.Require().NoError(err)
	s.NotContains(content, "password")
}

func (s *AuthTestSuite) TestWritesNeedAToken() {
	response, err := s.Http(s.T()).WithHeader("Content-Type", "application/json").
		Post("/api/clubs", strings.NewReader(`{"name":"Anon FC","ground":"Park"}`))
	s.Require().NoError(err)

	response.AssertUnauthorized()
}

func (s *AuthTestSuite) TestReadsDoNot() {
	response, err := s.Http(s.T()).Get("/api/clubs")
	s.Require().NoError(err)

	response.AssertOk()
}

func (s *AuthTestSuite) TestLoggedOutTokenIsRejected() {
	// A separate user, on purpose. Goravel's JWTs carry no unique ID, so the
	// same user logging in twice in the same second gets the identical token.
	// Logging this one out would also disable s.user's token in other tests.
	var leaver models.User
	s.Require().NoError(facades.Orm().Factory().Create(&leaver))
	request := s.LoginAs(s.T(), leaver)

	response, err := request.Post("/api/auth/logout", nil)
	s.Require().NoError(err)
	response.AssertNoContent()

	response, err = request.Get("/api/auth/me")
	s.Require().NoError(err)
	response.AssertUnauthorized()
}

// #endregion test

func (s *AuthTestSuite) TestGarbageTokenIsRejected() {
	response, err := s.Http(s.T()).WithToken("not-a-jwt").Get("/api/auth/me")
	s.Require().NoError(err)

	response.AssertUnauthorized()
}
