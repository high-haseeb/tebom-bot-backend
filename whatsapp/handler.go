package whatsapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
)

func WaSendPDF(filePath string) error {
	mediaId, err := uploadPDFToWhatsApp(filePath)
	if err != nil {
		return err
	}

	err = sendPDFToWhatsApp(mediaId, os.Getenv("WHATSAPP_CLOUD_API_TEST_RECIPIENT"), filePath);
	if err != nil {
		return err
	}
	return nil
}

func uploadPDFToWhatsApp(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	fileName := filepath.Base(filePath)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	partHeaders := make(textproto.MIMEHeader)
	partHeaders.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, fileName))
	partHeaders.Set("Content-Type", "application/pdf")

	part, err := writer.CreatePart(partHeaders)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return "", err
	}
	_ = writer.WriteField("messaging_product", "whatsapp")

	writer.Close()

	req, err := http.NewRequest("POST", "https://graph.facebook.com/v22.0/505601559313159/media", body)
	if err != nil {
		return "", err
	}

	token := os.Getenv("WHATSAPP_CLOUD_API_ACCESS_TOKEN")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var responseData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		return "", err
	}
	fmt.Println(responseData)
	mediaID, ok := responseData["id"].(string)
	if !ok {
		return "", fmt.Errorf("invalid media response")
	}

	return mediaID, nil
}

func sendPDFToWhatsApp(mediaID, recipient, filePath string) error {
	data := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                recipient,
		"type":              "document",
		"document": map[string]string{
			"id":       mediaID,
			"caption":  filepath.Base(filePath),
			"filename": filepath.Base(filePath),
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://graph.facebook.com/v18.0/505601559313159/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	token := os.Getenv("WHATSAPP_CLOUD_API_ACCESS_TOKEN")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to send WhatsApp message: %s", string(body))
	}

	return nil
}
