package tax

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Allowance struct {
	AllowanceType string  `json:"allowanceType"`
	Amount        float64 `json:"amount"`
}

type TaxRequest struct {
	TotalIncome float64     `json:"totalIncome"`
	WHT         float64     `json:"wht"`
	Allowances  []Allowance `json:"allowances"`
}

type TaxResponse struct {
	Tax uint64 `json:"tax"`
}

type PersonalIncomTax struct {
	ID           int
	Level        int
	Description  string
	Percent_rate uint64
	Min_Amount   *int
	Max_Amount   *int
}

type Err struct {
	Message string `json:"message"`
}

func SumTotalIncomeWithAllowances(mas TaxRequest) uint64 {
	result := mas.TotalIncome

	for _, v := range mas.Allowances {
		result -= v.Amount
	}

	return uint64(result)
}

func Calculatation(c echo.Context) error {
	req := TaxRequest{}
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	// if req.TotalIncome < 0 {
	// 	req.TotalIncome = 0
	// }

	// if req.WHT < 0 {
	// 	req.WHT = 0
	// }

	// if req.TotalIncome < req.WHT {
	// 	return c.JSON(http.StatusOK, TaxResponse{Tax: 0})
	// }

	totalIncome := SumTotalIncomeWithAllowances(req)

	// tax := uint64(0)

	// mas_tax := GetMasPersonalIncomTax()

	return c.JSON(http.StatusOK, TaxResponse{Tax: totalIncome})
}
