package admin

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func DeductionPersosal(c echo.Context) error {

	return c.JSON(http.StatusCreated, "OK")
}
