package mctx

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"

	. "gopkg.in/check.v1"
)

// DBSuite ...
type DBSuite struct {
	DBConn    *gorm.DB
	Migration *Migration
}

// S ...
type S struct {
	Server *echo.Echo
	DBSuite
	Token string
}

var _ = Suite(&S{})

// PerformRequest is func call http to server
func (s *S) PerformRequest(method string, path string, params url.Values) (*httptest.ResponseRecorder, error) {
	paramsEncoded := params.Encode()
	reader := strings.NewReader(paramsEncoded)
	if method == "GET" || method == "HEAD" {
		path += "?" + paramsEncoded
	}

	request, err := http.NewRequest(method, path, reader)
	if err != nil {
		return nil, err
	}

	if method == "POST" || method == "PUT" {
		request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}
	if s.Token != "" {
		request.Header.Add("Authorization", s.Token)
	}

	response := httptest.NewRecorder()

	s.Server.ServeHTTP(response, request)
	return response, nil
}

// TestHealthCheck is func test for handler HealthCheck
func (s *S) TestHealthCheck(c *C) {
	_, err := s.PerformRequest("GET", "/healthcheck", url.Values{})
	c.Assert(err, Equals, nil)
}

// BenchmarkHealthCheck is func test benchmark of gocheck. You can use that go test -check.b -check.bmem
func (s *S) BenchmarkHealthCheck(c *C) {
	for i := 0; i < c.N; i++ {
		s.PerformRequest("GET", "/healthcheck", url.Values{})
	}
}
