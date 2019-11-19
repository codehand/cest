package mctx

import "time"

type options struct {
	timeout time.Duration
	caching bool
	method  string
	path    string
	params  map[string]string
	headers map[string]string
	body    interface{}
}

type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) {
	f(o)
}

func WithTimeout(t time.Duration) Option {
	return optionFunc(func(o *options) {
		o.timeout = t
	})
}

func WithCaching(cache bool) Option {
	return optionFunc(func(o *options) {
		o.caching = cache
	})
}

func WithMethod(method string) Option {
	return optionFunc(func(o *options) {
		o.method = method
	})
}

func WithPath(path string) Option {
	return optionFunc(func(o *options) {
		o.path = path
	})
}

func WithBody(body interface{}) Option {
	return optionFunc(func(o *options) {
		o.body = body
	})
}

func WithParams(params map[string]string) Option {
	return optionFunc(func(o *options) {
		o.params = params
	})
}

func WithHeaders(headers map[string]string) Option {
	return optionFunc(func(o *options) {
		o.headers = headers
	})
}
