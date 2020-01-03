package mctx

import "time"

type options struct {
	timeout time.Duration
	caching bool
	method  string
	path    string
	params  map[string]string
	headers map[string]string
	query   map[string]string
	body    interface{}
}

// Option is model define
type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) {
	f(o)
}

// WithTimeout is func set apply timeout
func WithTimeout(t time.Duration) Option {
	return optionFunc(func(o *options) {
		o.timeout = t
	})
}

// WithCaching is func set apply caching
func WithCaching(cache bool) Option {
	return optionFunc(func(o *options) {
		o.caching = cache
	})
}

// WithMethod is func set apply method
func WithMethod(method string) Option {
	return optionFunc(func(o *options) {
		o.method = method
	})
}

// WithPath is func set apply path
func WithPath(path string) Option {
	return optionFunc(func(o *options) {
		o.path = path
	})
}

// WithBody is func set apply body
func WithBody(body interface{}) Option {
	return optionFunc(func(o *options) {
		o.body = body
	})
}

// WithParams is func set apply params
func WithParams(params map[string]string) Option {
	return optionFunc(func(o *options) {
		o.params = params
	})
}

// WithHeaders is func set apply headers
func WithHeaders(headers map[string]string) Option {
	return optionFunc(func(o *options) {
		o.headers = headers
	})
}

// WithQuery is func set apply query
func WithQuery(queries map[string]string) Option {
	return optionFunc(func(o *options) {
		o.query = queries
	})
}
