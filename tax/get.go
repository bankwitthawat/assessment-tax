package tax

import (
	"fmt"
	"log"

	dbConfig "github.com/bankwitthawat/assessment-tax/pkg/db"
)

func GetMasPersonalIncomTax() []PersonalIncomeTax {

	stmt, err := dbConfig.DB.Prepare("SELECT * FROM mas_personal_income_tax")
	if err != nil {
		log.Fatal("can't prepare query all ", err)
	}

	rows, err := stmt.Query()
	if err != nil {
		log.Fatal("can't query all ", err)
	}

	tax_level := []PersonalIncomeTax{}

	for rows.Next() {
		p := PersonalIncomeTax{}
		fmt.Println(rows)
		err := rows.Scan(&p.ID, &p.Level, &p.Description, &p.Percent_rate, &p.Min_Amount, &p.Max_Amount)
		if err != nil {
			log.Fatal("can't Scan row ", err)
		}
		tax_level = append(tax_level, p)
	}

	return tax_level

}

func GetDeduction(types string) (float64, error) {
	result := float64(0)
	row := dbConfig.DB.QueryRow("SELECT amount FROM mas_deductions WHERE type=$1", types)
	err := row.Scan(&result)
	if err != nil {
		//return default value
		return 60000, err
	}
	return result, nil
}
