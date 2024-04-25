package tax

import (
	"fmt"
	"log"
)

func GetMasPersonalIncomTax() []PersonalIncomTax {

	stmt, err := DB.Prepare("SELECT * FROM mas_personal_income_tax")
	if err != nil {
		log.Fatal("can't prepare query all ", err)
	}

	rows, err := stmt.Query()
	if err != nil {
		log.Fatal("can't query all ", err)
	}

	tax_level := []PersonalIncomTax{}

	for rows.Next() {
		p := PersonalIncomTax{}
		fmt.Println(rows)
		err := rows.Scan(&p.ID, &p.Level, &p.Description, &p.Percent_rate, &p.Min_Amount, &p.Max_Amount)
		if err != nil {
			log.Fatal("can't Scan row ", err)
		}
		tax_level = append(tax_level, p)
	}

	return tax_level

}
