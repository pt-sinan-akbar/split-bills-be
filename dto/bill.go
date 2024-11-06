package dto

type DynamicUpdateRequest struct {
	BillId string `json:"bill_id"`
	Item   struct {
		Id       int     `json:"id"`
		Price    float64 `json:"price"`
		Quantity int     `json:"quantity"`
	} `json:"item"`
	Tax     float64 `json:"tax"`
	Service float64 `json:"service"`
}
