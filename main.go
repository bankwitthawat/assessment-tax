package main

import (
	"net/http"
	"os"

	"github.com/bankwitthawat/assessment-tax/admin"
	dbConfig "github.com/bankwitthawat/assessment-tax/pkg/db"
	"github.com/bankwitthawat/assessment-tax/tax"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func BasicAuthMiddleware(username, password string, c echo.Context) (bool, error) {
	if username == "adminTax" && password == "admin!" {
		return true, nil
	}

	au := os.Getenv("ADMIN_USERNAME")
	ap := os.Getenv("ADMIN_PASSWORD")

	if au == "adminTax" && ap == "admin!" {
		return true, nil
	}

	return false, nil
}

func EnvAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		au := os.Getenv("ADMIN_USERNAME")
		ap := os.Getenv("ADMIN_PASSWORD")

		if au == "adminTax" && ap == "admin!" {
			return next(c)
		}

		return echo.ErrUnauthorized
	}
}

func main() {

	port := os.Getenv("PORT")

	dbConfig.InitDB()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})

	t := e.Group("/tax")
	t.POST("/calculations", tax.Calculatation)

	a := e.Group("/admin")
	// a.Use(middleware.BasicAuth(BasicAuthMiddleware))
	a.Use(EnvAuthMiddleware)
	a.POST("/personal", admin.DeductionPersosal)

	e.Logger.Fatal(e.Start(":" + port))
}
