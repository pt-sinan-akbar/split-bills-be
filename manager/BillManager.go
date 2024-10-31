package manager

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/pt-sinan-akbar/initializers"
	"gorm.io/gorm"
	"io"
	"mime/multipart"
)

type BillManager struct {
	DB     *gorm.DB
	config initializers.Config
}

func NewBillManager(DB *gorm.DB, config initializers.Config) *BillManager {
	return &BillManager{DB, config}
}

func (bm BillManager) SaveImage(image *multipart.FileHeader) (string, error) {
	imageBytes, err := image.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open image: %v", err)
	}
	defer imageBytes.Close()

	fileBytes, err := io.ReadAll(imageBytes)
	if err != nil {
		return "", fmt.Errorf("failed to read image: %v", err)
	}

	client := resty.New()
	supabaseURL := bm.config.SupabaseStorageURLL
	bucketName := bm.config.SupabaseBucket
	fileName := image.Filename
	jwtToken := bm.config.SupabaseSecretKey
	resp, err := client.R().
		SetHeader("Authorization", "Bearer"+jwtToken).
		SetHeader("Content-Type", "image/jpeg").
		SetBody(fileBytes).
		Post(fmt.Sprintf("%s/storage/v1/upload/%s/%s", supabaseURL, bucketName, fileName))
	if err != nil {
		return "", fmt.Errorf("failed to upload image: %v", err)
	}
	if resp.StatusCode() != 200 && resp.StatusCode() != 201 {
		return "", fmt.Errorf("failed to upload image: %v", resp.String())
	}

	return fileName, nil
}
