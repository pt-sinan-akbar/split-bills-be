package main

import (
	"fmt"
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
	initializers.DB.AutoMigrate(&models.Bill{}, &models.BillData{}, &models.BillOwner{}, &models.BillMember{}, &models.BillItem{})
	fmt.Println("? Migration complete")
}
