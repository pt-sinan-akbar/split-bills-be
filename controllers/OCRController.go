package controllers

import (
	"gorm.io/gorm"
)

type OCRController struct {
	DB *gorm.DB
}

func NewOCRController(DB *gorm.DB) OCRController {
	return OCRController{DB}
}
