package helpers

import(
	"bytes"
	"fmt"
	"github.com/fogleman/gg"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"image/color"
	"image/png"
	"math/rand"
	"time"
)
type ErrResponse struct {
	Message string `json:"message"`
}

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