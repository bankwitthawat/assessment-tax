package admin

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type DeductionPersonalRequest struct {
	Amount float64 `json:"amount"`
}

type DeductionPersonalResponse struct {
	PersonalDeduction float64 `json:"personalDeduction"`
}

type Err struct {
	Message string `json:"message"`
}

type MasDeductions struct {
	ID     int
	Type   string
	Amount float64
}

func DeductionPersosal(c echo.Context) error {

	req := DeductionPersonalRequest{}
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	//check value from db insert or update if exist
	dp, err := GetDeductionPersonal()

	if err != nil {

		if err == sql.ErrNoRows {
			// insert
			fmt.Println("DeductionPersosal insert")
			ist, err := CreateDeductionPersonal(req.Amount)
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

	fmt.Println("DeductionPersosal update")
	err = UpdateDeductionPersonal(cur)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, DeductionPersonalResponse{PersonalDeduction: cur.Amount})
}
