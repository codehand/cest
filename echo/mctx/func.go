package mctx

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"

	"github.com/labstack/echo/v4"
)

const (
	version = "echo v4"
)

type Args interface {
	// GetVal() interface{}
}

type A1 int

func (a A1) GetVal() interface{} {
	return 1
}

var (
	// _ Args = (*string)(nil)
	// _ Args = (*int)(nil)
	_ Args = (*A1)(nil)
)

func setArgs(v int) A1 {
	return 1
}

type ResultFn struct {
	Data  interface{}
	Error error
	Kind  reflect.Kind
}
type TaskFn struct {
	Fn         echo.HandlerFunc
	URL        string
	Method     string
	Host       string
	Header     map[string]string
	Params     map[string]string
	Args       Args
	Payload    interface{}
	Result     interface{}
	StatusCode int
	WantErr    bool
}

type EchoFn struct {
	AfterFn  []*TaskFn
	BeforeFn []*TaskFn
	version  string
}

func newResultFn(obj interface{}, err error) *ResultFn {
	return &ResultFn{
		Data:  obj,
		Error: err,
		Kind:  reflect.Bool,
	}
}

func (r *ResultFn) HasError() bool {
	return r.Error != nil
}

func (r *ResultFn) ToString() string {
	return fmt.Sprintf("%v", reflect.ValueOf(r.Data).Elem())
}

func (r *ResultFn) Value(pointer interface{}) {
	p := reflect.ValueOf(pointer).Elem().Type()
	if p.Kind() == reflect.Ptr {
		fmt.Printf("not support type  %v \n", p)
		return
	}
	if p != r.kind() {
		fmt.Printf("not match between %v and %v\n", p, r.kind())
		return
	}
	switch q := pointer.(type) {
	case *string:
		{
			*q = reflect.ValueOf(r.Data).Elem().String()
		}
	case *int:
		{
			*q = int(reflect.ValueOf(r.Data).Elem().Int())
		}
	case *bool:
		{
			*q = reflect.ValueOf(r.Data).Elem().Bool()
		}
	default:
		{
			fmt.Printf("Not found type %v\n", reflect.ValueOf(pointer).Elem().Type())
		}
	}
	return
}

func (r *ResultFn) kind() reflect.Type {
	v := reflect.ValueOf(r.Data).Elem()
	return v.Type()
}

func (r *ResultFn) pointer() reflect.Type {
	v := reflect.ValueOf(r.Data).Elem()
	return reflect.PtrTo(v.Type())
}

// DoBefore ...
func (s *EchoFn) DoBefore(e *echo.Echo) []*ResultFn {
	results := make([]*ResultFn, 0)
	for i, hl := range s.BeforeFn {
		if hl.Result == nil {
			results = append(results, newResultFn(nil, fmt.Errorf("[E.0] hl.Fn() index %d failes because not match define result return", i)))
			break
		}
		ctx, _, res := NewContext(e, hl.Method, hl.URL, hl.Params, hl.Header, hl.Payload)
		if err := hl.Fn(ctx); err != nil {
			results = append(results, newResultFn(nil, fmt.Errorf("[E.1] hl.Fn() index %d error = %v", i, err)))
			break
		}
		if !hl.WantErr {
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				results = append(results, newResultFn(nil, fmt.Errorf("[E.2] hl.Fn() index %d read body response with error = %v", i, err)))
				break
			}
			if res.Code != hl.StatusCode {
				es := BaseStatus{}
				err = json.Unmarshal(body, &es)
				if err == nil {
					results = append(results, newResultFn(nil, fmt.Errorf("[E.3.0] hl.Fn() %d call failed because %s", i, es.GetMsg())))
				}
				results = append(results, newResultFn(nil, fmt.Errorf("[E.3] hl.Fn() %d call failed error = %v with `%s`", i, err, string(body))))
			}
			var a = reflect.TypeOf(hl.Result)
			v := reflect.New(a.Elem())
			newP := v.Interface()
			err = json.Unmarshal(body, newP)
			if err != nil {
				results = append(results, newResultFn(nil, fmt.Errorf("[E.4] hl.Fn() index %d read body response with error = %v", i, err)))
			}
			s1 := reflect.ValueOf(newP).Elem()
			switch s1.Kind() {
			case reflect.Struct:
				{
					for r := 0; r < s1.NumField(); r++ {
						fieldValue := s1.Field(r)
						fieldType := s1.Type().Field(r)
						if fieldType.Name == "AccessToken" {
							p := fieldValue.Elem()
							s.BeforeFn[i+1].Header = map[string]string{
								"Authorization": "Bearer " + p.String(),
							}
							break
						}
					}
				}
			case reflect.Bool:
				{
					results = append(results, newResultFn(newP, nil))
				}
			case reflect.String:
				{
					results = append(results, newResultFn(newP, nil))
				}
			case reflect.Int:
				{
					results = append(results, newResultFn(newP, nil))
				}
			}

		} else {
			// want to err
		}
	}
	return results
}

// DoAfter ...
func (s *EchoFn) DoAfter(e *echo.Echo) []*ResultFn {
	results := make([]*ResultFn, 0)
	for i, hl := range s.AfterFn {
		if hl.Result == nil {
			results = append(results, newResultFn(nil, fmt.Errorf("[E.0] hl.Fn() index %d failes because not match define result return", i)))
			break
		}
		ctx, _, res := NewContext(e, hl.Method, hl.URL, hl.Params, hl.Header, hl.Payload)
		if err := hl.Fn(ctx); err != nil {
			results = append(results, newResultFn(nil, fmt.Errorf("[E.1] hl.Fn() index %d error = %v", i, err)))
			break
		}
		if !hl.WantErr {
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				results = append(results, newResultFn(nil, fmt.Errorf("[E.2] hl.Fn() index %d read body response with error = %v", i, err)))
				break
			}
			if res.Code != hl.StatusCode {
				es := BaseStatus{}
				err = json.Unmarshal(body, &es)
				if err == nil {
					results = append(results, newResultFn(nil, fmt.Errorf("[E.3.0] hl.Fn() %d call failed because %s", i, es.GetMsg())))
				}
				results = append(results, newResultFn(nil, fmt.Errorf("[E.3] hl.Fn() %d call failed error = %v with `%s`", i, err, string(body))))
			}
			var a = reflect.TypeOf(hl.Result)
			v := reflect.New(a.Elem())
			newP := v.Interface()
			err = json.Unmarshal(body, newP)
			if err != nil {
				results = append(results, newResultFn(nil, fmt.Errorf("[E.4] hl.Fn() index %d read body response with error = %v", i, err)))
			}
			s1 := reflect.ValueOf(newP).Elem()
			switch s1.Kind() {
			case reflect.Struct:
				{
					for r := 0; r < s1.NumField(); r++ {
						fieldValue := s1.Field(r)
						fieldType := s1.Type().Field(r)
						if fieldType.Name == "AccessToken" {
							p := fieldValue.Elem()
							s.BeforeFn[i+1].Header = map[string]string{
								"Authorization": "Bearer " + p.String(),
							}
							break
						}
					}
				}
			case reflect.Bool:
				{
					results = append(results, newResultFn(newP, nil))
				}
			case reflect.String:
				{
					results = append(results, newResultFn(newP, nil))
				}
			case reflect.Int:
				{
					results = append(results, newResultFn(newP, nil))
				}
			}

		} else {
			// want to err
		}
	}
	return results
}

func (s *EchoFn) AfterFnIndex(pos int) *TaskFn {
	if s == nil || s.AfterFn == nil || pos >= len(s.AfterFn) || pos < 0 {
		return nil
	}
	return s.AfterFn[pos]
}

func (s *EchoFn) BeforeFnIndex(pos int) *TaskFn {
	if s == nil || s.BeforeFn == nil || pos >= len(s.BeforeFn) || pos < 0 {
		return nil
	}
	return s.BeforeFn[pos]
}

func (s *EchoFn) Version() string {
	if s.version == "" {
		return version
	}
	return s.version
}

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
