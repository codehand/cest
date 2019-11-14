package handler

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/labstack/echo/v4"
)

func Package() string {
	pc, _, _, _ := runtime.Caller(1)
	parts := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	pl := len(parts)
	pkage := ""
	funcName := parts[pl-1]
	if parts[pl-2][0] == '(' {
		funcName = parts[pl-2] + "." + funcName
		pkage = strings.Join(parts[0:pl-2], ".")
	} else {
		pkage = strings.Join(parts[0:pl-1], ".")
	}
	return pkage
}

func HealthCheck(c echo.Context) error {
	// To do something
	fmt.Println(Package())
	return c.JSON(http.StatusOK, 1)
}

// golang
