package tax

import "testing"

func TestSumTotalIncomeWithAllowances(t *testing.T) {
	cases := []struct {
		name  string
		input TaxRequest
		want  uint64
	}{
		{
			name: "input income minus should be 0",
			input: TaxRequest{
				TotalIncome: -1,
				WHT:         0,
				Allowances: []Allowance{
					{AllowanceType: "donation", Amount: 0},
				},
			},
			want: 0,
		},
		{
			name: "input income 0 should be 0",
			input: TaxRequest{
				TotalIncome: 0,
				WHT:         0,
				Allowances: []Allowance{
					{AllowanceType: "donation", Amount: 0},
				},
			},
			want: 0,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {

			got := SumTotalIncomeWithAllowances(tt.input)

			if uint64(got) != tt.want {
				t.Errorf("got %d, want %d", uint64(got), tt.want)
			}
		})
	}
}

func TestSumTaxLevel(t *testing.T) {
	masTaxLevel := []PersonalIncomeTax{
		{ID: 1, Level: 1, Description: "0-150,000", Percent_rate: 0, Min_Amount: 0, Max_Amount: 150000},
		{ID: 2, Level: 2, Description: "150,001-500,000", Percent_rate: 10, Min_Amount: 150001, Max_Amount: 500000},
		{ID: 3, Level: 3, Description: "500,001-1,000,000", Percent_rate: 0, Min_Amount: 500001, Max_Amount: 1000000},
		{ID: 4, Level: 4, Description: "1,000,000-2,000,000", Percent_rate: 0, Min_Amount: 1000001, Max_Amount: 2000000},
		{ID: 5, Level: 5, Description: "2,000,000 ขึ้นไป", Percent_rate: 0, Min_Amount: 2000001, Max_Amount: 2000001},
	}
	cases := []struct {
		name  string
		input TaxLevel
		want  uint64
	}{
		{
			name: "if master taxLevel empty, tax should be 0",
			input: TaxLevel{
				totalIncome:          0,
				wht:                  0,
				masPersonalIncomeTax: []PersonalIncomeTax{},
			},
			want: 0,
		},

		{
			name: "if wht greater than income, tax should be 0",
			input: TaxLevel{
				totalIncome:          0,
				wht:                  100,
				masPersonalIncomeTax: masTaxLevel,
			},
			want: 0,
		},

		//
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {

			got := SumTaxLevel(tt.input)

			if got.Tax != tt.want {
				t.Errorf("got %d, want %d", got.Tax, tt.want)
			}
		})
	}
}
