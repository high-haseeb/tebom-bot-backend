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

	req.Header.Set("Authorization", "Bearer EAAHEr4Ftd4ABO1rX7lOppkFZA3SuGaAAsjwe3rikyOH9ycIi3lgzFyUZCEZCtmELmUIvyOLaHuHZADwr6ch1wXoopkW56yEHsy8cdniKYX8PFrRC4YdhRqC4NrJvmOj4aZCEPHgMmCpK6F57jxq5RSPKFmVZCVoPfjLSSeoLVJZCwFpN554OQZAsBprMLrvno3CsrZBSMgmisMq9llG4DVfjPEiLIzAH1EI1RGDHTXs7Y")
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

func sendPDFToWhatsApp(mediaID, recipient string) error {
	data := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                recipient,
		"type":              "document",
		"document": map[string]string{
			"id":       mediaID,
			"caption":  "Here is your requested PDF",
			"filename": "output.pdf",
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

	req.Header.Set("Authorization", "Bearer EAAHEr4Ftd4ABO1rX7lOppkFZA3SuGaAAsjwe3rikyOH9ycIi3lgzFyUZCEZCtmELmUIvyOLaHuHZADwr6ch1wXoopkW56yEHsy8cdniKYX8PFrRC4YdhRqC4NrJvmOj4aZCEPHgMmCpK6F57jxq5RSPKFmVZCVoPfjLSSeoLVJZCwFpN554OQZAsBprMLrvno3CsrZBSMgmisMq9llG4DVfjPEiLIzAH1EI1RGDHTXs7Y")
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
