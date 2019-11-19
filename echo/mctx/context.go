package mctx

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo/v4"
)

// NewContext is func new context echo
func NewContext(e *echo.Echo, method, path string, params map[string]string, headers map[string]string, body interface{}) (echo.Context, *http.Request, *httptest.ResponseRecorder) {
	// var payload *strings.Reader = nil
	var obj []byte
	if body != nil && method != "GET" && method != "" {
		obj, _ = json.Marshal(body)
		// payload = strings.NewReader(string(obj))
	}
	req := httptest.NewRequest(method, "/", strings.NewReader(string(obj)))
	res := httptest.NewRecorder()
	ctx := e.NewContext(req, res)
	ctx.SetPath(path)
	for k, v := range headers {
		ctx.Request().Header.Add(k, v)
	}
	var ks []string
	var vs []string
	for k, v := range params {
		ks = append(ks, k)
		vs = append(vs, v)
	}
	if len(ks) > 0 {
		ctx.SetParamNames(ks...)
		ctx.SetParamValues(vs...)
	}

	ctx.Request().Header.Add("Content-Type", "application/json")
	return ctx, req, res
}

// NewCustomContext is func new context echo with opts
func NewCustomContext(e *echo.Echo, opts ...Option) (echo.Context, *http.Request, *httptest.ResponseRecorder) {
	// default configs
	options := options{}
	for _, o := range opts {
		o.apply(&options)
	}
	var obj []byte
	if options.body != nil && options.method != "GET" && options.method != "" {
		obj, _ = json.Marshal(options.body)
		// payload = strings.NewReader(string(obj))
	}
	req := httptest.NewRequest(options.method, "/", strings.NewReader(string(obj)))
	res := httptest.NewRecorder()
	ctx := e.NewContext(req, res)
	ctx.SetPath(options.path)
	for k, v := range options.headers {
		ctx.Request().Header.Add(k, v)
	}
	var ks []string
	var vs []string
	for k, v := range options.params {
		ks = append(ks, k)
		vs = append(vs, v)
	}

	// set varible vadiric
	if len(ks) > 0 {
		ctx.SetParamNames(ks...)
		ctx.SetParamValues(vs...)
	}

	ctx.Request().Header.Add("Content-Type", "application/json")
	return ctx, req, res
}
