package main

import (
	"log"
	"net/http"
	"os"
    "time"

	"github.com/joho/godotenv"
)

func main() {
    godotenv.Load(".env");
    PORT := os.Getenv("PORT");

    go HandleWebsocket();
    mux := http.NewServeMux();

    if err := GetToken(); err != nil {
        log.Println("ERROR: can not read the token", err.Error());
    }

    mux.Handle("/get/vehicleInfo", Middleware(GetVehicleInformation));
    mux.Handle("/startOffers", Middleware(StartOffer));
    mux.Handle("/get/offers", Middleware(GetOffers));
    mux.Handle("/get/PDF", Middleware(GetPDF));
    mux.HandleFunc("/webhook", WaHandleWebhooks);

    go func() {
        ticker := time.NewTicker(30 * time.Minute)
        defer ticker.Stop()

        for {
            responseData := send_MFA_Request()
            log.Println("INFO: sent get MFA request, server response status:", responseData)
            <-ticker.C 
        }
    }()


    log.Printf("INFO: Listening on port %s\n", PORT);
    if err := http.ListenAndServe(PORT, mux); err != nil {
    	log.Fatalf("ERROR: Failed to start HTTP server: %s\n", err.Error());
    }
}
