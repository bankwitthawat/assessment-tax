package admin

import (
	"fmt"
	"log"

	dbConfig "github.com/bankwitthawat/assessment-tax/pkg/db"
)

func GetDeduction(types string) (MasDeductions, error) {
	result := MasDeductions{}
	stmt, err := dbConfig.DB.Prepare("SELECT * FROM mas_deductions WHERE type=$1; ")
	if err != nil {
		log.Fatal("can't prepare query one row statment", err)
	}

	rowType := types
	row := stmt.QueryRow(rowType)
	err = row.Scan(&result.ID, &result.Type, &result.Amount)
	if err != nil {
		// log.Fatal("can't Scan row into variables", err)
		return MasDeductions{}, err
	}

	return result, nil
}

func CreateDeduction(types string, amount float64) (MasDeductions, error) {
	result := MasDeductions{}

	row := dbConfig.DB.QueryRow("INSERT INTO mas_deductions (type, amount) values ($1, $2)  RETURNING id, type, amount", types, amount)
	err := row.Scan(&result.ID, &result.Type, &result.Amount)
	if err != nil {
		fmt.Println("can't scan id", err)
		return MasDeductions{}, err
	}

	return result, nil
}

func UpdateDeduction(u MasDeductions) error {

	stmt, err := dbConfig.DB.Prepare("UPDATE mas_deductions SET amount=$2 WHERE id=$1 RETURNING id, type, amount;")

	if err != nil {
		// log.Fatal("can't prepare statment update", err)
		return err
	}

	if _, err := stmt.Exec(u.ID, u.Amount); err != nil {
		// log.Fatal("error execute update ", err)
		return err
	}

	return nil
}
