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

			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSumTaxLevel(t *testing.T) {
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
				totalIncome: 0,
				wht:         100,
				masPersonalIncomeTax: []PersonalIncomeTax{
					{ID: 1, Level: 1, Description: ""},
				},
			},
			want: 0,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {

			got := SumTaxLevel(tt.input)

			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
