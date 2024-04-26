package tax

import (
	"fmt"
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
	Tax      uint64             `json:"tax"`
	TaxLevel []TaxLevelResponse `json: "taxLevel"`
}

type TaxLevelResponse struct {
	Level string `json:"level"`
	Tax   uint64 `json:"tax"`
}

type PersonalIncomeTax struct {
	ID           int
	Level        int
	Description  string
	Percent_rate float64
	Min_Amount   float64
	Max_Amount   float64
}

type Err struct {
	Message string `json:"message"`
}

type TaxLevel struct {
	totalIncome          float64
	wht                  float64
	masPersonalIncomeTax []PersonalIncomeTax
}

func SumTotalIncomeWithAllowances(mas TaxRequest) float64 {
	result := mas.TotalIncome - 60000 //หักค่าลดหย่อนส่วนตัว

	for _, v := range mas.Allowances {
		if v.AllowanceType == "donation" || v.AllowanceType == "k-receipt" {
			result -= v.Amount
		}
	}

	return result
}

func SumTaxLevel(r TaxLevel) TaxResponse {
	result := TaxResponse{}
	totalIncome := r.totalIncome
	totalTax := uint64(0)
	result.Tax = uint64(0)

	for i := 0; i < len(r.masPersonalIncomeTax); i++ {

		taxPerLevel := 0.0

		if r.totalIncome < r.masPersonalIncomeTax[i].Min_Amount || len(r.masPersonalIncomeTax) == 0 {
			result.TaxLevel = append(result.TaxLevel, TaxLevelResponse{
				Level: r.masPersonalIncomeTax[i].Description,
				Tax:   uint64(taxPerLevel),
			})
			fmt.Println("result.Tax 1", result.Tax)
			continue
		}

		if totalIncome > r.masPersonalIncomeTax[i].Max_Amount {
			totalIncome = totalIncome - r.masPersonalIncomeTax[i].Max_Amount
		}

		if r.masPersonalIncomeTax[i].Percent_rate > 0 {
			totalIncome = totalIncome * (r.masPersonalIncomeTax[i].Percent_rate / 100)
			taxPerLevel = totalIncome
		}

		result.TaxLevel = append(result.TaxLevel, TaxLevelResponse{
			Level: r.masPersonalIncomeTax[i].Description,
			Tax:   uint64(taxPerLevel),
		})

		totalTax += uint64(taxPerLevel)
	}

	if r.wht > r.totalIncome {
		result.Tax = 0
		return result
	}

	result.Tax = totalTax - uint64(r.wht)
	return result
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
		wht:                  req.WHT,
		masPersonalIncomeTax: masTaxLevel,
	}
	tax := SumTaxLevel(sumTax)
	fmt.Println("SumTaxLevel: ", tax)

	return c.JSON(http.StatusOK, tax)
}
