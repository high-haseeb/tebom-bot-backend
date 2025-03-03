package main

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

func WaSendText(to, message string) error {
	data := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                to,
		"type":              "text",
		"text": map[string]string{
			"body": message,
		},
	}

	if err := WaSendMessage(data); err != nil {
		return err
	}

	return nil
}

func WaSendFlow(to string) error {
	data := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                to,
		"type":              "interactive",
		"interactive": map[string]interface{}{
			"type": "flow",
			"header": map[string]interface{}{
				"type": "text",
				"text": "Flow message header",
			},
			"body": map[string]interface{}{
				"text": "Flow message body",
			},
			"footer": map[string]interface{}{
				"text": "Flow message footer",
			},
			"action": map[string]interface{}{
				"name": "flow",
				"parameters": map[string]interface{}{
					"flow_message_version": "3",
					// "flow_token":           "AQAAAAACS5FpgQ_cAAAAAD0QI3s.",
					// "flow_id":              "1",
					// 609237645413288

					"flow_cta":    "Book!",
					"flow_action": "navigate",
					"flow_action_payload": map[string]interface{}{
						"screen": "<SCREEN_NAME>",
						"data": map[string]interface{}{
							"product_name":        "name",
							"product_description": "description",
							"product_price":       100,
						},
					},
				},
			},
		},
	}

	if err := WaSendMessage(data); err != nil {
		return err
	}

	return nil
}

func WaSendInteractive(to string) error {
	data := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                to,
		"type":              "interactive",
		"interactive": map[string]interface{}{
			"type": "cta_url",
 			"header": map[string]interface{}{ 
      			"type": "text",
      			"text": "TebomNET",
    		},
			"body": map[string]interface{}{
				"text": "Trafik Sigortası teklifi için lütfen formu doldurunuz.",
			},

			"action": map[string]interface{}{
				"name": "cta_url",
      			"parameters": map[string]interface{}{
        			"display_text": "Gönder",
        			"url": "http://188.132.135.5:8080/",
      			},
			},
		},
	}
	if err := WaSendMessage(data); err != nil {
		return err
	}
	return nil;
}

func WaSendListMessage(to, header, body, footer string, buttonText string, sections []map[string]interface{}) error {
	data := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                to,
		"type":              "interactive",
		"interactive": map[string]interface{}{
			"type": "list",
			"body": map[string]string{
				"text": body,
			},
			"footer": map[string]string{
				"text": footer,
			},
			"action": map[string]interface{}{
				"button": buttonText,
				"sections": sections,
			},
		},
	}

	return WaSendMessage(data)
}

func WaSendMessage(data map[string]interface{}) error {
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

func WaSendPDF(filePath string) error {
	mediaId, err := WaUploadMedia(filePath)
	if err != nil {
		return err
	}

	err = WaSendDocument(mediaId, os.Getenv("WHATSAPP_CLOUD_API_TEST_RECIPIENT"), filePath)
	if err != nil {
		return err
	}
	return nil
}

func WaUploadMedia(filePath string) (string, error) {
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

func WaSendDocument(mediaID, recipient, filePath string) error {
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

	if err := WaSendMessage(data); err != nil {
		return err
	}

	return nil
}


type Status struct {
	ID         string `json:"id"`
	Status     string `json:"status"`
	Timestamp  string `json:"timestamp"`
	Recipient  string `json:"recipient_id"`
}

type Text struct {
	Body string `json:"body"` 
}

type Image struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type Audio struct {
	ID  string `json:"id"`  
	URL string `json:"url"`
}

type Message struct {
	ID        string  `json:"id"`        
	From      string  `json:"from"`     
	Timestamp string  `json:"timestamp"`
	Type      string  `json:"type"`     
	Text      *Text   `json:"text,omitempty"` 
	Image     *Image  `json:"image,omitempty"`
	Audio     *Audio  `json:"audio,omitempty"`
}

type WebhookRequest struct {
	Entry []struct {
		Changes []struct {
			Value struct {
				Statuses []Status `json:"statuses,omitempty"`
				Messages []Message `json:"messages,omitempty"`
			} `json:"value"`
		} `json:"changes"`
	} `json:"entry"`
}

func webhookVerify(w http.ResponseWriter, r *http.Request) {
	mode := r.URL.Query().Get("hub.mode");
	token := r.URL.Query().Get("hub.verify_token");
	challenge := r.URL.Query().Get("hub.challenge");

	verifyToken := os.Getenv("VERIFICATION_CODE");

	if mode == "subscribe" && token == verifyToken {
		fmt.Println("WEBHOOK_VERIFIED")
		w.WriteHeader(http.StatusOK);
		w.Write([]byte(challenge));
	} else {
		fmt.Println("VERIFICATION_FAILED");
		http.Error(w, `{"status":"error","message":"Verification failed"}`, http.StatusForbidden);
	}
}

func webhookPost(w http.ResponseWriter, r *http.Request) {
	var req WebhookRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading request body:", err)
		http.Error(w, `{"status":"error","message":"Failed to read request body"}`, http.StatusInternalServerError)
		return
	}

	// Decode JSON body
	err = json.Unmarshal(body, &req)
	if err != nil {
		fmt.Println("Failed to decode JSON:", err, "Raw Body:", string(body))
		http.Error(w, `{"status":"error","message":"Invalid JSON provided"}`, http.StatusBadRequest)
		return
	}

	// Check if there are entries
	if len(req.Entry) == 0 {
		fmt.Println("No 'entry' field in webhook request:", string(body))
		http.Error(w, `{"status":"error","message":"Invalid webhook structure"}`, http.StatusBadRequest)
		return
	}

	for _, entry := range req.Entry {
		for _, change := range entry.Changes {
			// Handle WhatsApp status updates
			if len(change.Value.Statuses) > 0 {
				status := change.Value.Statuses[0]
				fmt.Printf("Status Update Received - Status: %s\n", status.Status)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"status":"ok"}`))
				return
			}

			// Handle WhatsApp messages
			if len(change.Value.Messages) > 0 {
				message := change.Value.Messages[0]
				if message.Text != nil {
					fmt.Printf("Received Message - From: %s, Message: %s\n", message.From, message.Text.Body)
					handleMessage(change.Value.Messages)
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{"status":"ok"}`))
					return
				} else {
					fmt.Println("Received message but no text content:", message)
				}
			}
		}
	}

	// Not a valid WhatsApp API event
	http.Error(w, `{"status":"error","message":"Not a WhatsApp API event"}`, http.StatusNotFound);
}

func handleMessage(messages []Message) {
	for _, msg := range messages {

		from := msg.From;
		text := msg.Text.Body;

		if text == "" {
			fmt.Println("No text message found");
			continue
		}

		if err := WaSendInteractive(from); err != nil {
			fmt.Println(err.Error());
		}
	}
}

func WaHandleWebhooks(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		webhookVerify(w, r)
	} else if r.Method == http.MethodPost {
		webhookPost(w, r)
	} else {
		http.Error(w, `{"status":"error","message":"Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}
