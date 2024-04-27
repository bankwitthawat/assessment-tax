package tax

import (
	"encoding/csv"
	"fmt"
	"io"
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
	result := mas.TotalIncome
	maxp, err := GetDeduction("personal")
	if err != nil {
		fmt.Println("Error:", err)
		maxp = 60000
	}

	maxk, err := GetDeduction("k-receipt")
	if err != nil {
		maxk = 50000
	}
	if maxk > 100000 {
		maxk = 100000
	}

	result = result - maxp

	for _, v := range mas.Allowances {
		if v.AllowanceType == "donation" {
			result -= v.Amount
		}

		if v.AllowanceType == "k-receipt" {
			if v.Amount > 100000 {
				v.Amount = maxk
			}
			if v.Amount < 0 {
				v.Amount = 0
			}
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

	if uint64(r.wht) > totalTax {
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

type TaxData struct {
	TotalIncome float64 `json:"totalIncome"`
	WHT         float64 `json:"wht"`
	Donation    float64 `json:"donation"`
}

type Taxes struct {
	TotalIncome float64 `json:"totalIncome"`
	Tax         float64 `json:"tax"`
}

type TaxDataResponse struct {
	Taxes []Taxes `json:"taxes"`
}

func UploadCSV(c echo.Context) error {

	// Read form file
	file, err := c.FormFile("taxFile")
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Create a new CSV reader
	r := csv.NewReader(src)

	// Skip header row
	_, err = r.Read()
	if err != nil && err != io.EOF {
		return err
	}

	var data []TaxData

	// Read CSV records
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Parse CSV record into TaxData struct
		td := TaxData{}
		fmt.Sscanf(record[0], "%f", &td.TotalIncome)
		fmt.Sscanf(record[1], "%f", &td.WHT)
		fmt.Sscanf(record[2], "%f", &td.Donation)

		fmt.Println(td.TotalIncome)
		fmt.Println(td.WHT)
		fmt.Println(td.Donation)

		// Append TaxData to data slice
		data = append(data, td)
	}

	//get tax level master
	masTaxLevel := GetMasPersonalIncomTax()
	sort.SliceStable(masTaxLevel, func(i, j int) bool {
		return masTaxLevel[i].Level < masTaxLevel[j].Level
	})

	result := TaxDataResponse{}
	taxes := []Taxes{}

	for _, t := range data {
		tt := TaxRequest{
			TotalIncome: t.TotalIncome,
			WHT:         t.WHT,
			Allowances: []Allowance{
				{AllowanceType: "donation", Amount: t.Donation},
			},
		}

		totalIncome := SumTotalIncomeWithAllowances(tt)

		sumTax := TaxLevel{
			totalIncome:          totalIncome,
			wht:                  tt.WHT,
			masPersonalIncomeTax: masTaxLevel,
		}
		st := SumTaxLevel(sumTax)

		taxes = append(taxes, Taxes{TotalIncome: sumTax.totalIncome, Tax: float64(st.Tax)})
	}

	result.Taxes = taxes

	return c.JSON(http.StatusOK, result)
}
