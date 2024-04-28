package admin

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
)

type DeductionRequest struct {
	Amount float64 `json:"amount"`
}

type DeductionPersonalResponse struct {
	PersonalDeduction float64 `json:"personalDeduction"`
}

type DeductionKReceiptResponse struct {
	KReceipt float64 `json:"kReceipt"`
}

type Err struct {
	Message string `json:"message"`
}

type MasDeductions struct {
	ID     int
	Type   string
	Amount float64
}

func SetDeductionPersosal(c echo.Context) error {

	req := DeductionRequest{}
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	//check value from db insert or update if exist
	dp, err := GetDeduction("personal")

	if err != nil {

		if err == sql.ErrNoRows {
			// insert
			ist, err := CreateDeduction("personal", req.Amount)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
			}
			c.JSON(http.StatusCreated, ist.Amount)

		} else {
			return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
		}
	}

	//update
	cur := MasDeductions{
		ID:     dp.ID,
		Type:   dp.Type,
		Amount: req.Amount,
	}

	err = UpdateDeduction(cur)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, DeductionPersonalResponse{PersonalDeduction: cur.Amount})
}

func SetDeductionKReceipt(c echo.Context) error {

	req := DeductionRequest{}
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	if req.Amount < 0 {
		return c.JSON(http.StatusBadRequest, Err{Message: "k-recept can't set less than 0"})
	}

	if req.Amount > 100000 {
		return c.JSON(http.StatusBadRequest, Err{Message: "k-recept can't set more than 100,000"})
	}

	dk, err := GetDeduction("k-receipt")
	if err != nil {

		if err == sql.ErrNoRows {
			// insert
			ist, err := CreateDeduction("k-receipt", req.Amount)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
			}
			c.JSON(http.StatusCreated, ist.Amount)

		} else {
			return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
		}
	}

	//update
	cur := MasDeductions{
		ID:     dk.ID,
		Type:   dk.Type,
		Amount: req.Amount,
	}

	err = UpdateDeduction(cur)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, DeductionKReceiptResponse{KReceipt: cur.Amount})
}
