package helpers

import (
	"bytes"
	"fmt"
	"github.com/fogleman/gg"
	"github.com/go-resty/resty/v2"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/pt-sinan-akbar/initializers"
	"image/color"
	"image/png"
	"math/rand"
	"time"
)

func GenerateID() (string, error) {
	id, err := gonanoid.Generate("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 8)
	if err != nil {
		return "", fmt.Errorf("failed to generate id: %v", err)
	}
	return id, nil
}

func GenerateInitialsImage(name string) ([]byte, error) {
	const S = 256
	rand.Seed(time.Now().UnixNano())
	bgColor := color.RGBA{
		R: 138,
		G: 206,
		B: 0,
		A: 255,
	}

	dc := gg.NewContext(S, S)
	dc.SetColor(bgColor)
	dc.Clear()
	dc.SetRGB(1, 1, 1)

	//initial := string(name[0])
	dc.DrawStringAnchored("brat", S/2, S/2, 0.5, 0.5)
	var buf bytes.Buffer
	err := png.Encode(&buf, dc.Image())
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
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
