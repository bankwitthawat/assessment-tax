package tax

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Allowance struct {
	AllowanceType string `json:"allowanceType"`
	Amount        uint64 `json:"amount"`
}

type IncomeDetails struct {
	TotalIncome uint64      `json:"totalIncome"`
	WHT         uint64      `json:"wht"`
	Allowances  []Allowance `json:"allowances"`
}

type TaxResponse struct {
	Tax uint64 `json:"tax"`
}

func Calculatation(c echo.Context) error {

	return c.JSON(http.StatusCreated, "")
}
