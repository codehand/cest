package mctx

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func timeoutMiddleware(next echo.HandlerFunc, timeout time.Duration) echo.HandlerFunc {
	return func(c echo.Context) error {
		// wrap the request context with a timeout
		ctx, cancel := context.WithTimeout(c.Request().Context(), timeout)
		defer func() {
			// check if context timeout was reached
			if ctx.Err() == context.DeadlineExceeded {
				// write response and abort the request
				c.Response().WriteHeader(http.StatusGatewayTimeout)

				// 	c.Abort()
			}

			//cancel to clear resources after finished
			cancel()
		}()
		// replace request with context wrapped request
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}
