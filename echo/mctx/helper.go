package mctx

import (
	"net/url"
	"reflect"
	"strconv"
)

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
