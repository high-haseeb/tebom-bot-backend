package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func GetCookie() string {
	//BUG: the cookie is invalidated after some time.
	// Make a system to check invalid cookies and send a request to refresh.
	return os.Getenv("TEBOM_COOKIE")
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func RespondWithError(w http.ResponseWriter, errorType, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json");
	w.WriteHeader(statusCode);
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:   errorType,
		Message: message,
	});
}

type GetVehicleInformationReq struct {
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

func GenerateVehicleInformationForm(request GetVehicleInformationReq) string {
	formData := url.Values{};
	formData.Set("input[Calisilanfirma]", request.Calisilanfirma);
	formData.Set("input[Calisilansube]", request.Calisilansube);
	formData.Set("input[Calisilanuser]", request.Calisilanuser);
	formData.Set("input[IsYK]", boolToString(request.IsYK));
	formData.Set("input[NationalNumber]", request.NationalNumber);
	formData.Set("input[LicensePlateNumber]", request.LicensePlateNumber);
	formData.Set("input[LicensePermitNumber]", request.LicensePermitNumber);
	formData.Set("input[Phone]", request.Phone);
	formData.Set("input[EMail]", request.EMail);
	formData.Set("input[HaveLicensePermitNumber]", boolToString(request.HaveLicensePermitNumber));
	formData.Set("input[IsSorgu]", boolToString(request.IsSorgu));
	formData.Set("input[IsDisabled]", boolToString(request.IsDisabled));
	formData.Set("input[ProfessionCode]", strconv.Itoa(request.ProfessionCode));
	formData.Set("input[MasterBranch]", strconv.Itoa(request.MasterBranch));
	formData.Set("input[MortgageeType]", request.MortgageeType);
	formData.Set("input[MortgageeBankCode]", request.MortgageeBankCode);
	formData.Set("input[MortgageeBankBranchCode]", request.MortgageeBankBranchCode);
	formData.Set("input[MortgageeFinancerCode]", request.MortgageeFinancerCode);
	return formData.Encode();
}

func GetVehicleInformation(w http.ResponseWriter, r *http.Request) {
	var request GetVehicleInformationReq;
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		RespondWithError(w, "INVALID_REQ", err.Error(), http.StatusBadRequest);
		return;
	}

	form := GenerateVehicleInformationForm(request);
	response, err := SendExternalFormRequest("https://portal.acente365.com/OfferNew/YeniTrafikBilgi", form);
	if err != nil {
		RespondWithError(w, "Can not send Request", err.Error(), http.StatusInternalServerError);
	}

	w.Header().Set("Content-Type", "application/json");
	body, err := io.ReadAll(response.Body);
	if err != nil {
	    RespondWithError(w, "Can not read response body", err.Error(), http.StatusInternalServerError);
	    return;
	}
	defer response.Body.Close();
	w.Write(body);
}

type StartOfferRequest struct {
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

func GenerateStartOfferForm(request StartOfferRequest) string {
	formData := url.Values{};
	formData.Set("input[BrandCodeFull]", request.BrandCodeFull);
	formData.Set("input[UsingType]", request.UsingType);
	formData.Set("input[BuildYear]", request.BuildYear);
	formData.Set("input[FuelType]", request.FuelType);
	formData.Set("input[ColorCode]", request.ColorCode);
	formData.Set("input[IsRenewalPeriodTraffic]", request.IsRenewalPeriodTraffic);
	formData.Set("input[ContinueWithoutOldPolicyInfo]", request.ContinueWithoutOldPolicyInfo);
	formData.Set("input[HeaderGuid]", request.HeaderGuid);
	formData.Set("input[VehicleTypeTraffic]", request.VehicleTypeTraffic);
	formData.Set("input[Branch]", request.Branch);
	formData.Set("input[QueryType]", request.QueryType);
	formData.Set("input[TestQueue]", request.TestQueue);
	return formData.Encode();
}

func StartOffer(w http.ResponseWriter, r *http.Request) {
	var request StartOfferRequest;
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		RespondWithError(w, "INVALID REQUEST", err.Error(), http.StatusBadRequest);
		return
	}

	if request.HeaderGuid == "" {
		RespondWithError(w, "Missing GUID", "the guid must be present", http.StatusBadRequest);
		return
	}

	URL := "https://portal.acente365.com/OfferNew/YeniTrafikGuncelle"
	form := GenerateStartOfferForm(request);
	response, err := SendExternalFormRequest(URL, form);
	if err != nil {
		RespondWithError(w, "Something went wrong", err.Error(), http.StatusInternalServerError);
		return;
	}

	w.Header().Set("Content-Type", "application/json");
	body, err := io.ReadAll(response.Body);
	if err != nil {
	    RespondWithError(w, "Can not read response body", err.Error(), http.StatusInternalServerError);
	    return;
	}
	w.Write(body);
}


type GetOffersRequest struct {
	HeaderGuid string `json:"guid"`
}

func GetOffers(w http.ResponseWriter, r *http.Request) {
	var request GetOffersRequest;
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		RespondWithError(w, "INVALID REQUEST", err.Error(), http.StatusBadRequest);
		return
	}

	if request.HeaderGuid == "" {
		RespondWithError(w, "Missing GUID", "the guid must be present", http.StatusBadRequest);
		return
	}

	URL := "https://portal.acente365.com/OfferNew/GetListOffer?HeaderGuid=" + request.HeaderGuid + "&QueryType=trafik"
	response, err := SendExternalRequest(URL);
	if err != nil {
		RespondWithError(w, "Something went wrong", err.Error(), http.StatusInternalServerError);
		return;
	}

	w.Header().Set("Content-Type", "application/json");
	w.Write(response);
}

type GetPDFRequest struct {
	Type         string `json:"type"`
	SelectedItem string `json:"selectedItem"`
	HeaderGuid   string `json:"headerGuid"`
	QueryType    string `json:"queryType"`
	GroupType    string `json:"groupType"`
}

type GetPDFResponse struct {
	Success bool   `json:"Success"`
	Message string `json:"Message"`
	File    struct {
		FileContents     string `json:"FileContents"`
		FileDownloadName string `json:"FileDownloadName"`
	} `json:"file"`
}

func GetPDF(w http.ResponseWriter, r *http.Request) {
	var request GetPDFRequest;
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		RespondWithError(w, "INVALID REQUEST", err.Error(), http.StatusBadRequest)
		return;
	}

	if request.HeaderGuid == "" {
		RespondWithError(w, "Missing GUID", "the guid must be present", http.StatusBadRequest)
		return
	}

	URL := "https://portal.acente365.com/Offer/DownloadAllOffersPdf"
	form := GenerateGetPDFFrom(request)

	response, err := SendExternalFormRequest(URL, form);
	if err != nil {
		RespondWithError(w, "Something went wrong", err.Error(), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close();

	var responseData GetPDFResponse;
	if err := json.NewDecoder(response.Body).Decode(&responseData); err != nil {
		RespondWithError(w, "can not decode response form server", err.Error(), http.StatusInternalServerError);
		return;
	}

	fileName := responseData.File.FileDownloadName;
	fmt.Println(fileName);

	filePath := "./output/" + fileName;
	dataPDF, err := base64.StdEncoding.DecodeString(responseData.File.FileContents)
	if err != nil {
		RespondWithError(w, "failed to decode pdf", err.Error(), http.StatusInternalServerError);
		return;
	}

	err = os.WriteFile(filePath, dataPDF, 0644);
	if err != nil {
		RespondWithError(w, "failed to created file", err.Error(), http.StatusInternalServerError);
		return;
	}
    WaSendPDF(filePath);

	w.Header().Set("Content-Type", "application/pdf")

	json.NewEncoder(w).Encode(responseData);
}

func GenerateGetPDFFrom(request GetPDFRequest) string {
	formData := url.Values{};
	formData.Set("type", request.Type);
	formData.Set("selectedItem", request.SelectedItem);
	formData.Set("headerGuid", request.HeaderGuid);
	formData.Set("queryType", request.QueryType);
	formData.Set("groupType", request.GroupType);
	return formData.Encode();
}
 
func SendExternalRequest(URL string) ([]byte, error) {
	client := &http.Client{};
	req, err := http.NewRequest("POST", URL, nil);
	if err != nil {
		return nil, err;
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", GetCookie());

	resp, err := client.Do(req);
	if err != nil {
		return nil, err;
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body);
	if err != nil {
		return nil, err;
	}
    return body, nil;
}

func SendExternalFormRequest(URL, form string) (*http.Response, error) {
	client := &http.Client{};
	req, err := http.NewRequest("POST", URL, strings.NewReader(form));
	if err != nil {
		return nil, err;
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded");
	req.Header.Set("Cookie", GetCookie());

	resp, err := client.Do(req);
	if err != nil {
		return nil, err;
	}

    return resp, nil;
}

func Middleware(handler func(http.ResponseWriter, *http.Request)) http.Handler {
	return CORSMiddleware(http.HandlerFunc(handler));
} 

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*");
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS");
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization");

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return;
		}

		next.ServeHTTP(w, r);
	})
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
