package main

import (
	"fmt"
	"github.com/pt-sinan-akbar/initializers"
	"github.com/pt-sinan-akbar/models"
	"gorm.io/gorm"
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
	initializers.DB.AutoMigrate(
		&models.Bill{},
		&models.BillData{},
		&models.BillOwner{},
		&models.BillMember{},
		&models.BillItem{},
	)
	Session := initializers.DB.Session(&gorm.Session{
		PrepareStmt: true,
	})
	if Session != nil {
		fmt.Println("? Migration complete")
	}
}
