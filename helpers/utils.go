package helpers

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/pt-sinan-akbar/initializers"
)

type ErrResponse struct {
	Message string `json:"message"`
}

func UploadFile(filename string, fileBytes []byte) error {
	client := resty.New()
	supabaseURL := initializers.ConfigSetting.SupabaseStorageURL
	bucketName := initializers.ConfigSetting.SupabaseBucket
	folderName := initializers.ConfigSetting.SupbaseFolder
	fileName := filename
	jwtToken := initializers.ConfigSetting.SupabaseSecretKey
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
