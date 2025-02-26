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
var COOKIE = "INGRESSCOOKIE=1740608801.557.1096.240344|a341718d8fdd83143d3341dcba51991c; .AspNetCore.Session=CfDJ8Jqe8KNpVVFPg2qq9%2FOnQc6I1cnJPlnq5sX1VGAtTQUMKXe%2BzyP3jZOHEfaQF9DGr8EtzbxDYLxs6PfztzOw%2BTJoR69V2PHEp9ESXQ9HeJkQh8N1qFuumNlldVd7qjiByLhL90O313Z5uFUlqU7kcJDHXp6W4dmutgsbI7GPXcSK; QueryType=trafik; mainHeaderGuid="

func main() {
	const PORT = ":4040"
	mux := http.NewServeMux()
	mux.HandleFunc("/getTrafficInfo", getTrafficInformation)
	mux.HandleFunc("/getOffers", getOffers)
	mux.HandleFunc("/startOfferList", sendOfferRequest)
	mux.HandleFunc("/getPDF", getPDF)

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

type GetOfferPayload struct {
	BrandCodeFull                string `json:"BrandCodeFull"`
	UsingType                    string `json:"UsingType"`
	BuildYear                    string `json:"BuildYear"`
	FuelType                     string `json:"FuelType"`
	ColorCode                    string `json:"ColorCode"`
	IsRenewalPeriodTraffic       string `json:"IsRenewalPeriodTraffic"`
	ContinueWithoutOldPolicyInfo string `json:"ContinueWithoutOldPolicyInfo"`
	HeaderGuid                   string `json:"HeaderGuid"`
	VehicleTypeTraffic           string `json:"VehicleTypeTraffic"`
	Branch                       string `json:"Branch"`
	QueryType                    string `json:"QueryType"`
	TestQueue                    string `json:"TestQueue"`
}

func sendOfferRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var requestData GetOfferPayload
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if requestData.HeaderGuid == "" {
		http.Error(w, "Missing HeaderGuid", http.StatusBadRequest)
		return
	}

	externalURL := "https://portal.acente365.com/OfferNew/YeniTrafikGuncelle"

	formData := url.Values{}
	formData.Set("input[BrandCodeFull]", requestData.BrandCodeFull)
	formData.Set("input[UsingType]", requestData.UsingType)
	formData.Set("input[BuildYear]", requestData.BuildYear)
	formData.Set("input[FuelType]", requestData.FuelType)
	formData.Set("input[ColorCode]", requestData.ColorCode)
	formData.Set("input[IsRenewalPeriodTraffic]", requestData.IsRenewalPeriodTraffic)
	formData.Set("input[ContinueWithoutOldPolicyInfo]", requestData.ContinueWithoutOldPolicyInfo)
	formData.Set("input[HeaderGuid]", requestData.HeaderGuid)
	formData.Set("input[VehicleTypeTraffic]", requestData.VehicleTypeTraffic)
	formData.Set("input[Branch]", requestData.Branch)
	formData.Set("input[QueryType]", requestData.QueryType)
	formData.Set("input[TestQueue]", requestData.TestQueue)

	client := &http.Client{}
	req, err := http.NewRequest("POST", externalURL, strings.NewReader(formData.Encode()))
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("cookie", COOKIE)
	req.Header.Set("origin", "https://portal.acente365.com")
	req.Header.Set("referer", "https://portal.acente365.com/offer/TrafficOffer?HeaderGuid="+requestData.HeaderGuid+"&listOffer=true")
	req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")

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


type PDFRequest struct {
	Type        string `json:"type"`
	SelectedItem string `json:"selectedItem"`
	HeaderGuid  string `json:"headerGuid"`
	QueryType   string `json:"queryType"`
	GroupType   string `json:"groupType"`
}

func getPDF(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var requestData PDFRequest
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if requestData.HeaderGuid == "" {
		http.Error(w, "Missing HeaderGuid", http.StatusBadRequest)
		return
	}

	externalURL := "https://portal.acente365.com/Offer/DownloadAllOffersPdf"

	formData := url.Values{}
	formData.Set("type", requestData.Type)
	formData.Set("selectedItem", requestData.SelectedItem)
	formData.Set("headerGuid", requestData.HeaderGuid)
	formData.Set("queryType", requestData.QueryType)
	formData.Set("groupType", requestData.GroupType)

	client := &http.Client{}
	req, err := http.NewRequest("POST", externalURL, strings.NewReader(formData.Encode()))
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("accept", "*/*")
	req.Header.Set("origin", "https://portal.acente365.com")
	req.Header.Set("referer", "https://portal.acente365.com/offer/TrafficOffer?HeaderGuid="+requestData.HeaderGuid+"&listOffer=true")
	req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")
	req.Header.Set("x-requested-with", "XMLHttpRequest")
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

	w.Header().Set("Content-Type", "application/pdf")
	w.Write(body)
}

