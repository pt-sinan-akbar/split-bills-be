package manager

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/pt-sinan-akbar/helper"
	"github.com/pt-sinan-akbar/initializers"
	"gorm.io/gorm"
	"io"
	"mime/multipart"
)

type BillManager struct {
	DB     *gorm.DB
	Config initializers.Config
}

func NewBillManager(DB *gorm.DB) (*BillManager, error) {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		return nil, fmt.Errorf("could not load config: %v", err)
	}
	return &BillManager{DB, config}, nil
}

func (bm BillManager) SaveImage(image *multipart.FileHeader, id string) error {
	imageBytes, err := image.Open()
	if err != nil {
		return fmt.Errorf("failed to open image: %v", err)
	}
	defer imageBytes.Close()

	fileBytes, err := io.ReadAll(imageBytes)
	if err != nil {
		return fmt.Errorf("failed to read image: %v", err)
	}

	client := resty.New()
	supabaseURL := bm.Config.SupabaseStorageURL
	bucketName := bm.Config.SupabaseBucket
	folderName := bm.Config.SupbaseFolder
	fileName := id
	jwtToken := bm.Config.SupabaseSecretKey
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+jwtToken).
		SetHeader("Content-Type", "image/jpeg").
		SetBody(fileBytes).
		Post(fmt.Sprintf("%s/storage/v1/object/%s/%s/%s", supabaseURL, bucketName, folderName, fileName))
	fmt.Printf("##########: %s/storage/v1/object/%s/%s/%s\n", supabaseURL, bucketName, folderName, fileName)
	if err != nil {
		return fmt.Errorf("failed to upload image: %v", err)
	}
	if resp.StatusCode() != 200 && resp.StatusCode() != 201 {
		return fmt.Errorf("failed to upload image: %v", resp.String())
	}

	return nil
}

func (bm BillManager) UploadImageToPython(image *multipart.FileHeader, id string) error {
	//var bill models.Bill
	//bill.ID = id
	//bill.RawImage = id
	//
	//// send image to python
	//
	//// if response.Error != nil
	//// save bill to db
	//result := bm.DB.Create(&bill)
	//if result.Error != nil {
	//	return fmt.Errorf("failed to save bill: %v", result.Error)
	//}
	//// iterate over python response to create BillData and BillItem
	//
	return nil
}

func (bm BillManager) UploadBill(image *multipart.FileHeader) error {
	// generate id
	id, err := helper.GenerateID()
	if err != nil {
		return err
	}
	fmt.Println("DONE GENERATING ID")

	// upload image to python
	if err := bm.UploadImageToPython(image, id); err != nil {
		return err
	}
	fmt.Println("DONE UPLOADING IMAGE TO PYTHON")

	// upload image to supabase
	if err := bm.SaveImage(image, id); err != nil {
		fmt.Println("ERROR UPLOADING IMAGE TO SUPABASE")
		fmt.Println("ID: ", id)
		return err
	}
	fmt.Println("DONE UPLOADING IMAGE TO SUPABASE")

	return nil
}
