package main

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/bankwitthawat/assessment-tax/tax"
)

func main() {

	tax.InitDB()

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})

	t := e.Group("/tax")
	t.POST("/calculations", tax.Calculatation)

	// e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
	// 	if username == "apidesign" || password == "45678" {
	// 		return true, nil
	// 	}
	// 	return false, nil
	// }))

	e.Logger.Fatal(e.Start(":8080"))
}
