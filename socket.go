package main

import (
	"crypto/rand"
	"crypto/sha256"
    "encoding/json"
	"encoding/hex"
	"log"
    "strings"
	"net/url"
    "os"
	"time"
    "fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
    wsURL  = "wss://sock5s.acente365.com/socket.io/?EIO=4&transport=websocket"
    loginURL = "https://portal.acente365.com/account/login"
    token  = "9A2754A2886AF16C2849582998F19EB5DB84E22DEC2161A8679A5CF2928G43DE|"
    user   = "rizacenkercivelek"
    pass   = "Cenk.2025"
    pingInterval = 25 * time.Second
    cookieFilePath = "cookies.txt"
)

type LoginResponse struct {
    StatusCode int  `json:"StatusCode"`;
    Success    bool `json:"Success"`;
}

func GetCookie() string {
	data, err := os.ReadFile("cookies.txt")
	if err != nil {
		log.Println("ERROR: Can not read cookies:", err.Error())
		return ""
	}

	cookies := strings.Split(strings.TrimSpace(string(data)), "\n")
	var ingressCookie, sessionCookie string

	for _, cookie := range cookies {
		cookie = strings.TrimSpace(cookie)
		if strings.HasPrefix(cookie, "INGRESSCOOKIE=") {
			ingressCookie = strings.Split(cookie, ";")[0]
		} else if strings.HasPrefix(cookie, ".AspNetCore.Session=") {
			sessionCookie = strings.Split(cookie, ";")[0]
		}
	}

	if ingressCookie == "" || sessionCookie == "" {
		log.Println("ERROR: Required cookies not found")
		return ""
	}

	return fmt.Sprintf("%s; %s", ingressCookie, sessionCookie)
}

func SaveCookie(responseHeaders http.Header) {
	cookies := responseHeaders.Values("Set-Cookie")
	cookieData := strings.Join(cookies, "\n")
	err := os.WriteFile(cookieFilePath, []byte(cookieData), 0644)
	if err != nil {
        log.Println("ERROR: can not save cookies:", err.Error());
		return;
	}

    log.Printf("INFO: Cookies saved successfully to %", cookieFilePath);
}

func send_MFA_Request() string {
	data := url.Values{}
	data.Set("loginModel[Username]", "rizacenkercivelek")
	data.Set("loginModel[Password]", "9R9pKrkI5dRhgJEyND0LSg==")
	data.Set("loginModel[RedirectUrl]", "")
	data.Set("loginModel[MfaControl]", "false")
	data.Set("loginModel[MfaCode]", "")
	data.Set("loginModel[Guid]", "186b5417-793b-4358-85b1-a9ad5b7d5471")

	client := &http.Client{}
	req, err := http.NewRequest("POST", loginURL, strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Referer", "https://portal.acente365.com/account")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	log.Println("GET_MFA: Response Status:", resp.Status);
	return resp.Status;
}

func sendLoginRequestWithMFA(mfaCode string) {
	data := url.Values{}
	data.Set("loginModel[Username]", "rizacenkercivelek");
	data.Set("loginModel[Password]", "9R9pKrkI5dRhgJEyND0LSg==");
	data.Set("loginModel[RedirectUrl]", "");
	data.Set("loginModel[MfaControl]", "true");
	data.Set("loginModel[MfaCode]", mfaCode);
	data.Set("loginModel[Guid]", "186b5417-793b-4358-85b1-a9ad5b7d5471");

	client := &http.Client{};
	req, err := http.NewRequest("POST", loginURL, strings.NewReader(data.Encode()));
	if err != nil {
        log.Fatalf("ERROR: can not request to %s:%s\n", loginURL, err.Error());
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Referer", "https://portal.acente365.com/account")

	resp, err := client.Do(req);
	if err != nil {
		log.Fatal(err);
	}
	defer resp.Body.Close();

    var response LoginResponse;
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        log.Println("ERROR: can not decode response from server", err.Error());
	}
    if !response.Success {
        log.Println("ERROR: did not receive cookie from server:", response.StatusCode);
    } else {
        log.Println("INFO: got cookie:", response.Success);
        SaveCookie(resp.Header);
    }
}

func GenerateUnique() string {
	// Generate 8 random 32-bit (4-byte) numbers
	var bytes [32]byte;
	_, err := rand.Read(bytes[:])
	if err != nil {
		panic(err)
	}
	hexString := hex.EncodeToString(bytes[:16])

	return fmt.Sprintf("%s%s-%s%s-%s%s-%s%s",
		hexString[0:4],   hexString[4:8],
        hexString[8:12],  hexString[12:16],
        hexString[16:20], hexString[20:24],
		hexString[24:28], hexString[28:32]);
}

func GenerateChannel(user, password string) string {
	normalizedUser := strings.ReplaceAll(strings.ToLower(user), " ", "")
	combined := normalizedUser + password;
	hash := sha256.Sum256([]byte(combined));
	return hex.EncodeToString(hash[:]);
}

func keepConnectionAlive(conn *websocket.Conn) {
	for {
		time.Sleep(pingInterval)
		err := conn.WriteMessage(websocket.TextMessage, []byte("3"));
		if err != nil {
			log.Println("Ping failed, reconnecting...");
			return
		}
        log.Println("INFO: ping");
	}
}

func HandleWebsocket() {
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil);
	if err != nil {
        log.Fatalf("ERROR: Failed to connect: %v", err);
	}
	defer conn.Close();
    log.Println("INFO: Connected to WebSocket");

	// Read initial message (should contain sid)
	_, message, err := conn.ReadMessage();
	if err != nil {
        log.Fatalf("ERROR: Read error: %v", err);
	}
    log.Printf("INFO: SID message: %s", message);

	authMessage := fmt.Sprintf(`40{"token":"%s%s"}`, token, GenerateUnique());
	err = conn.WriteMessage(websocket.TextMessage, []byte(authMessage));
	if err != nil {
        log.Fatalf("ERROR: Write error: %v", err)
	}
    log.Println("INFO: Sent authentication message");

    _, msg, err := conn.ReadMessage();
    fmt.Println(msg);
    if err != nil {
        log.Printf("ERROR: WS can not read message: %v", err);
        return;
    }

    channel := GenerateChannel("rizacenkercivelek", "Cenk.2025");
    joinMessage := fmt.Sprintf(`42["join-request","%s"]`, channel);
	conn.WriteMessage(websocket.TextMessage, []byte(joinMessage));
    log.Println("INFO: sent join-request");

	go keepConnectionAlive(conn);

    responseData := send_MFA_Request();
    log.Println("INFO: sent get MFA request, server response status:", responseData);

	for {
		_, msg, err := conn.ReadMessage();
		if err != nil {
            log.Printf("ERROR: WS can not read message: %v", err);
			break;
		}
        messageStr := string(msg);
        log.Printf("INFO: WS Received: %s", messageStr);

        if strings.Contains(messageStr, `"loginRequest"`) {
            mfaCode := extractLoginCode(messageStr);
            log.Printf("INFO: WS got MFA code %s", mfaCode);
            sendLoginRequestWithMFA(mfaCode);
        }
	}
}

func extractLoginCode(message string) string {
	var parsedMessage []any
	err := json.Unmarshal([]byte(message[2:]), &parsedMessage) // Skip first 2 chars "42"
	if err != nil {
        log.Println("ERROR: JSON Parse Error:", err.Error());
		return ""
	}

	if len(parsedMessage) > 1 {
		dataMap, ok := parsedMessage[1].(map[string]any)
		if ok {
			data, ok := dataMap["data"].(map[string]any)
			if ok {
				loginRequestData, ok := data["loginRequest"].(bool)
				if ok && !loginRequestData {
					if loginCode, exists := data["loginCode"].(string); exists {
						return loginCode
					}
				}
			}
		}
	}
	return ""
}
