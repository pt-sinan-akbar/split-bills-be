package initializers

import (
	"fmt"
	"github.com/pt-sinan-akbar/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func ConnectDB(config *Config) {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", config.DBHost, config.DBUserName, config.DBUserPassword, config.DBName, config.DBPort)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		log.Fatal("Failed to connect to the Database")
	}
	fmt.Println("? Connected Successfully to the Database")
	err = DB.SetupJoinTable(&models.BillMember{}, "BillItem", &models.BillMemberItem{})
	if err != nil {
		log.Fatalf("Failed to setup join table: %v", err)
	} else {
		fmt.Println("? Join table setup successfully")
	}
	err = DB.SetupJoinTable(&models.BillItem{}, "BillMember", &models.BillMemberItem{})
	if err != nil {
		log.Fatalf("Failed to setup join table: %v", err)
	} else {
		fmt.Println("? Join table setup successfully")
	}
}
