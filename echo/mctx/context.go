package mctx

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
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
	fmt.Println("options.query:", options.query)
	if len(options.query) > 0 {
		q := make(url.Values)
		for k, v := range options.query {
			q.Set(k, v)
		}
		options.path = options.path + "/?" + q.Encode()
		fmt.Println("???")
	}

	req := httptest.NewRequest(options.method, options.path, strings.NewReader(string(obj)))
	res := httptest.NewRecorder()
	ctx := e.NewContext(req, res)

	fmt.Println("options.path: ", options.path)
	ctx.SetPath(options.path)

	fmt.Println("ctx.Path(): ", ctx.Path())
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

// ConvertURL is func convert struct (interface) to map
func ConvertURL(i interface{}) (values url.Values) {
	values = url.Values{}
	if i == nil {
		return
	}
	values = url.Values{}
	iVal := reflect.ValueOf(i).Elem()
	typ := iVal.Type()
	for i := 0; i < iVal.NumField(); i++ {
		f := iVal.Field(i)
		tag := typ.Field(i).Tag.Get("json")
		if tag != "" {
			var v string
			switch f.Interface().(type) {
			case int, int8, int16, int32, int64:
				v = strconv.FormatInt(f.Int(), 10)
			case uint, uint8, uint16, uint32, uint64:
				v = strconv.FormatUint(f.Uint(), 10)
			case float32:
				v = strconv.FormatFloat(f.Float(), 'f', 4, 32)
			case float64:
				v = strconv.FormatFloat(f.Float(), 'f', 4, 64)
			case []byte:
				v = string(f.Bytes())
			case string:
				v = f.String()
			}
			values.Set(tag, v)
		}
	}
	return
}
