package main

import (
	"github.com/pt-sinan-akbar/initializers"
	"github.com/pt-sinan-akbar/models"
	"log"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables")
	}
	initializers.ConnectDB(&config)
}
func main() {
	tx := initializers.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Fatal("? Transaction failed and rolled back")
		}
	}()

	billOwner := models.BillOwner{
		ID:          1,
		UserId:      1,
		Name:        "Owner 1",
		Contact:     "Contact 1",
		BankAccount: "Bank Account 1",
	}

	billData := []models.BillData{
		{
			ID:        1,
			BillId:    "BILL-001",
			StoreName: "Store 1",
			SubTotal:  100000,
			Discount:  10000,
			Tax:       1000,
			Service:   5000,
			Total:     95000,
			Misc:      "Misc 1",
		},
		{
			ID:        2,
			BillId:    "BILL-001",
			StoreName: "Store 2",
			SubTotal:  200000,
			Discount:  20000,
			Tax:       2000,
			Service:   10000,
			Total:     190000,
			Misc:      "Misc 2",
		},
	}

	bill := models.Bill{
		ID:        "BILL-001",
		BillData:  billData,
		BillOwner: billOwner,
		Name:      "Bill 1",
		RawImage:  "Raw Image 1",
	}

	if err := tx.Create(&bill).Error; err != nil {
		tx.Rollback()
		log.Fatal("? Could not seed bill")
	}

	tx.Commit()
}
