package handler

import (
	"fmt"
	"testing"

	"github.com/codehand/cest/echo/mctx"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// TestPackage is func test
// Created at 15-11-2019 by tam gh
func TestPackage(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		if got := Package(); got != tt.want {
			t.Errorf("%q. Package() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

// TestHealthCheck is func test
// Created at 15-11-2019 by tam gh
func TestHealthCheck(t *testing.T) {
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		scriptsFn *mctx.EchoFn
	}{
		// TODO: Add test cases.
	}
	var e *echo.Echo = echo.New()

	for _, tt := range tests {
		if tt.scriptsFn != nil && tt.scriptsFn.BeforeFn != nil {
			// TODO: Add to do something before script.
			rs := tt.scriptsFn.DoBefore(e)
			assert.NotNil(t, rs)
			for _, item := range rs {
				assert.NoError(t, item.Error)
				assert.NotNil(t, item.Data)
				fmt.Printf("get*ValueOf: %v\n", item.ToString())
			}
		}
		ctx, _, _ := mctx.NewContext(e, echo.GET, "", nil, nil)
		if err := HealthCheck(ctx); err != nil {
			t.Errorf("%q. HealthCheck() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
		if tt.scriptsFn != nil && tt.scriptsFn.AfterFn != nil {
			// TODO: Add to do something after script.
			rs := tt.scriptsFn.DoAfter(e)
			assert.NotNil(t, rs)
			for _, item := range rs {
				assert.NoError(t, item.Error)
				assert.NotNil(t, item.Data)
				fmt.Printf("get*ValueOf: %v\n", item.ToString())
			}
		}
	}
}
