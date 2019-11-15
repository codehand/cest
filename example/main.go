package main

import (
	"fmt"
	"go/parser"
	"go/token"

	"github.com/codehand/cest/example/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func packageName(file string) (string, error) {
	fset := token.NewFileSet()

	// parse the go soure file, but only the package clause
	astFile, err := parser.ParseFile(fset, file, nil, parser.PackageClauseOnly)
	if err != nil {
		return "", err
	}

	if astFile.Name == nil {
		return "", fmt.Errorf("no package name found")
	}

	return astFile.Name.Name, nil
}

func main() {
	e := echo.New()
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
	}))
	e.Use(middleware.Logger())
	e.Use(middleware.RequestID())
	e.GET("/healthcheck", handler.HealthCheck)

	// var b bytes.Buffer

	err := e.Start(":9090")
	if err != nil {
		panic(err)
	}
}

func ABC() {

}
