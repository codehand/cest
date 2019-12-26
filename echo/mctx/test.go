package mctx

import (
	"fmt"
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

// WitMysqlMigration is func plugin run migrate data test
func (s *S) WitMysqlMigration(dbConn *gorm.DB, migrationsFolderLocation string) error {
	if s == nil {
		panic("Server with server nil")
	}

	if dbConn == nil {
		return fmt.Errorf("%v", "Connection nil")
	}
	s.DBConn = dbConn
	mi, err := runMigration(s.DBConn.DB(), migrationsFolderLocation)
	if err != nil {
		return err
	}
	s.Migration = mi
	return nil
}

// WithEcho is func setup echo
func (s *S) WithEcho(e *echo.Echo) *S {
	if s == nil {
		panic("Server with server nil")
	}
	if e == nil {
		s.Server = echo.New()
		return s
	}
	s.Server = e
	return s
}
