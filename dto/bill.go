package dto

type BillSummary struct {
	BillId  string              `json:"bill_id"`
	Name    string              `json:"name"`
	Members []BillSummaryMember `json:"members"`
	Contact BillSummaryContact  `json:"contact"`
}

type BillSummaryContact struct {
	// models.BillOwner simplified
	Name        string `json:"name"`
	Contact     string `json:"contact"`
	BankAccount string `json:"bank_account"`
}

type BillSummaryMember struct {
	// models.BillMember simplified
	Name      string            `json:"name"`
	PriceOwe  *float64          `json:"price_owe"`
	BillItems []BillSummaryItem `json:"items"`
}

type BillSummaryItem struct {
	// models.BillItem simplified
	Name     string  `json:"name"`
	Qty      int64   `json:"qty"`
	Price    float64 `json:"price"`
	Tax      float64 `json:"tax"`
	Service  float64 `json:"service"`
	Discount float64 `json:"discount"`
	Total    float64 `json:"total"`
}
