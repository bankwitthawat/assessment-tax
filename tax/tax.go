package tax

import (
	"net/http"
	"sort"

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

type PersonalIncomeTax struct {
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

type TaxLevel struct {
	totalIncome          uint64
	wht                  uint64
	masPersonalIncomeTax []PersonalIncomeTax
}

func SumTotalIncomeWithAllowances(mas TaxRequest) uint64 {
	result := mas.TotalIncome - 60000

	for _, v := range mas.Allowances {
		if v.AllowanceType == "donation" || v.AllowanceType == "k-receipt" {
			result -= v.Amount
		}
	}

	return uint64(result)
}

func SumTaxLevel(r TaxLevel) uint64 {
	if len(r.masPersonalIncomeTax) == 0 {
		return 0
	}
	return 0
}

func Calculatation(c echo.Context) error {
	req := TaxRequest{}
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	if req.TotalIncome < 0 {
		req.TotalIncome = 0
	}

	if req.WHT < 0 {
		req.WHT = 0
	}

	if req.TotalIncome < req.WHT {
		return c.JSON(http.StatusOK, TaxResponse{Tax: 0})
	}

	totalIncome := SumTotalIncomeWithAllowances(req)

	masTaxLevel := GetMasPersonalIncomTax()
	sort.SliceStable(masTaxLevel, func(i, j int) bool {
		return masTaxLevel[i].Level < masTaxLevel[j].Level
	})

	sumTax := TaxLevel{
		totalIncome:          totalIncome,
		wht:                  uint64(req.WHT),
		masPersonalIncomeTax: masTaxLevel,
	}

	tax := SumTaxLevel(sumTax)

	return c.JSON(http.StatusOK, TaxResponse{Tax: tax})
}
