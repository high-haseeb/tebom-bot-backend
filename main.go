package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// NOTE: The cookie will expire, make a system to refresh the cookie when it expires.
// For now we are replacing the cookie whenever it expires!
var COOKIE = "INGRESSCOOKIE=1740315502.326.3053.111819|a341718d8fdd83143d3341dcba51991c; .AspNetCore.Session=CfDJ8Hmpsi7h0nFGicsClGqI2ucqeksV8eOUCmDOGue%2Bx2yHcGSH99b5IfkdsZVYoq57Um1J5sk2qzEGj9zRw%2FbsQ5qdP6MWRoD%2Bx%2F816FnLqjKCP9qqr6EOZWoQVmlecXm3yco%2BqBQg4w7pAuPf%2BbGdd66HF6ihIqoLf7QWepJYKRlh; SessionId=2A4110F9291C8990E1163A83A493B3947DE937B5E7E6A06CB7CE851D5EB7F74F; QueryType=trafik; mainHeaderGuid="

func main() {
	const PORT = ":6969"
	mux := http.NewServeMux()
	mux.HandleFunc("/getTrafficInfo", getTrafficInformation)
	mux.HandleFunc("/getOffers", getOffers)

	fmt.Printf("INFO: Listening on port %s\n", PORT)
	if err := http.ListenAndServe(PORT, mux); err != nil {
		fmt.Printf("ERROR: Failed to start HTTP server: %s\n", err.Error())
	}
}

type TrafficInformation struct {
	Calisilanfirma          string `json:"Calisilanfirma"`
	Calisilansube           string `json:"Calisilansube"`
	Calisilanuser           string `json:"Calisilanuser"`
	IsYK                    bool   `json:"IsYK"`
	NationalNumber          string `json:"NationalNumber"`
	LicensePlateNumber      string `json:"LicensePlateNumber"`
	LicensePermitNumber     string `json:"LicensePermitNumber"`
	Phone                   string `json:"Phone"`
	EMail                   string `json:"EMail"`
	HaveLicensePermitNumber bool   `json:"HaveLicensePermitNumber"`
	IsSorgu                 bool   `json:"IsSorgu"`
	IsDisabled              bool   `json:"IsDisabled"`
	ProfessionCode          int    `json:"ProfessionCode"`
	MasterBranch            int    `json:"MasterBranch"`
	MortgageeType           string `json:"MortgageeType"`
	MortgageeBankCode       string `json:"MortgageeBankCode"`
	MortgageeBankBranchCode string `json:"MortgageeBankBranchCode"`
	MortgageeFinancerCode   string `json:"MortgageeFinancerCode"`
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func getTrafficInformation(w http.ResponseWriter, r *http.Request) {
	// CORS Bullshit
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var reqBody TrafficInformation
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	formData := url.Values{}
	formData.Set("input[Calisilanfirma]", reqBody.Calisilanfirma)
	formData.Set("input[Calisilansube]", reqBody.Calisilansube)
	formData.Set("input[Calisilanuser]", reqBody.Calisilanuser)
	formData.Set("input[IsYK]", boolToString(reqBody.IsYK))
	formData.Set("input[NationalNumber]", reqBody.NationalNumber)
	formData.Set("input[LicensePlateNumber]", reqBody.LicensePlateNumber)
	formData.Set("input[LicensePermitNumber]", reqBody.LicensePermitNumber)
	formData.Set("input[Phone]", reqBody.Phone)
	formData.Set("input[EMail]", reqBody.EMail)
	formData.Set("input[HaveLicensePermitNumber]", boolToString(reqBody.HaveLicensePermitNumber))
	formData.Set("input[IsSorgu]", boolToString(reqBody.IsSorgu))
	formData.Set("input[IsDisabled]", boolToString(reqBody.IsDisabled))
	formData.Set("input[ProfessionCode]", strconv.Itoa(reqBody.ProfessionCode))
	formData.Set("input[MasterBranch]", strconv.Itoa(reqBody.MasterBranch))
	formData.Set("input[MortgageeType]", reqBody.MortgageeType)
	formData.Set("input[MortgageeBankCode]", reqBody.MortgageeBankCode)
	formData.Set("input[MortgageeBankBranchCode]", reqBody.MortgageeBankBranchCode)
	formData.Set("input[MortgageeFinancerCode]", reqBody.MortgageeFinancerCode)

	encodedForm := formData.Encode()

	client := &http.Client{}
	externalURL := "https://portal.acente365.com/OfferNew/YeniTrafikBilgi"
	req, err := http.NewRequest("POST", externalURL, strings.NewReader(encodedForm))
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("cookie", COOKIE)

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to send request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

type RequestPayload struct {
	HeaderGuid string `json:"guid"`
}

func getOffers(w http.ResponseWriter, r *http.Request) {
	// CORS Headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle OPTIONS request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Ensure it's a POST request
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var requestData RequestPayload
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Ensure HeaderGuid is provided
	if requestData.HeaderGuid == "" {
		http.Error(w, "Missing HeaderGuid", http.StatusBadRequest)
		return
	}

	// Construct the external URL with the provided HeaderGuid
	externalURL := "https://portal.acente365.com/OfferNew/GetListOffer?HeaderGuid=" + requestData.HeaderGuid + "&QueryType=trafik"

	client := &http.Client{}
	req, err := http.NewRequest("GET", externalURL, nil)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("cookie", COOKIE)

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to send request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	// Set response headers and return the external API response
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
