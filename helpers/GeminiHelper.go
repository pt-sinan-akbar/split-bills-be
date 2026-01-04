package helpers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/go-resty/resty/v2"
	"github.com/pt-sinan-akbar/initializers"
)

type GeminiHelper struct{}

func NewGeminiHelper() GeminiHelper {
	return GeminiHelper{}
}

type GeminiRequest struct {
	Contents         []Content        `json:"contents"`
	GenerationConfig GenerationConfig `json:"generation_config"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text       string      `json:"text,omitempty"`
	InlineData *InlineData `json:"inline_data,omitempty"`
}

type InlineData struct {
	MimeType string `json:"mime_type"`
	Data     string `json:"data"`
}

type GenerationConfig struct {
	ResponseMimeType string `json:"response_mime_type"`
}

type GeminiResponse struct {
	Candidates []Candidate `json:"candidates"`
}

type Candidate struct {
	Content Content `json:"content"`
}

type ExtractedBillData struct {
	StoreName string          `json:"Store Name"`
	Date      string          `json:"Date"`
	Items     []ExtractedItem `json:"Items"`
	Subtotal  float64         `json:"Subtotal"`
	Tax       float64         `json:"Tax amount"`
	Service   float64         `json:"Service charge amount"`
	Discount  float64         `json:"Discount amount"`
	Total     float64         `json:"Total amount"`
}

type ExtractedItem struct {
	Name     string  `json:"Name"`
	Qty      int64   `json:"Quantity"`
	Price    float64 `json:"Price"`
	Subtotal float64 `json:"Subtotal"`
}

func (gh GeminiHelper) ExtractBillData(file *multipart.FileHeader) (ExtractedBillData, error) {
	// Read file content
	src, err := file.Open()
	if err != nil {
		return ExtractedBillData{}, fmt.Errorf("failed to open file: %v", err)
	}
	defer src.Close()

	fileBytes, err := io.ReadAll(src)
	if err != nil {
		return ExtractedBillData{}, fmt.Errorf("failed to read file: %v", err)
	}

	// Convert to base64
	base64Data := base64.StdEncoding.EncodeToString(fileBytes)

	// Prepare request
	client := resty.New()
	apiKey := initializers.ConfigSetting.GeminiApiKey
	if apiKey == "" {
		return ExtractedBillData{}, fmt.Errorf("GEMINI_API_KEY is not set")
	}

	prompt := `Extract the following data from the bill image:
	- Store Name
	- Date (YYYY-MM-DD format if possible)
	- Items (Name, Quantity, Price (unit price), Subtotal (qty * price))
	- Subtotal (before tax/service)
	- Tax amount
	- Service charge amount
	- Discount amount
	- Total amount
	
	Ensure numeric values are floats/integers. If a field is missing, use 0 or empty string. Return strictly JSON.`

	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=" + apiKey

	requestBody := GeminiRequest{
		Contents: []Content{
			{
				Parts: []Part{
					{Text: prompt},
					{
						InlineData: &InlineData{
							MimeType: file.Header.Get("Content-Type"),
							Data:     base64Data,
						},
					},
				},
			},
		},
		GenerationConfig: GenerationConfig{
			ResponseMimeType: "application/json",
		},
	}

	var response GeminiResponse
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(requestBody).
		SetResult(&response).
		Post(url)

	if err != nil {
		return ExtractedBillData{}, fmt.Errorf("gemini api request failed: %v", err)
	}

	if resp.IsError() {
		return ExtractedBillData{}, fmt.Errorf("gemini api returned error: %s", resp.String())
	}

	if len(response.Candidates) == 0 || len(response.Candidates[0].Content.Parts) == 0 {
		return ExtractedBillData{}, fmt.Errorf("empty response from gemini")
	}

	jsonText := response.Candidates[0].Content.Parts[0].Text

	var extractedData ExtractedBillData
	if err := json.Unmarshal([]byte(jsonText), &extractedData); err != nil {
		return ExtractedBillData{}, fmt.Errorf("failed to parse gemini response: %v, raw: %s", err, jsonText)
	}

	return extractedData, nil
}
