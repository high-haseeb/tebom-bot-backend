package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env");
	PORT := os.Getenv("PORT");
	// if err := WaSendFlow(); err != nil {
	// 	fmt.Printf(err.Error());
	// }

	mux := http.NewServeMux();
	// WaSendText("+923038023397", "Lütfen cep telefonunuzdaki giriş isteğini kabul edin");

	mux.Handle("/get/vehicleInfo", Middleware(GetVehicleInformation));
	mux.Handle("/startOffers", Middleware(StartOffer));
	mux.Handle("/get/offers", Middleware(GetOffers));
	mux.Handle("/get/PDF", Middleware(GetPDF));
	mux.HandleFunc("/webhook", WaHandleWebhooks);

	log.Printf("INFO: Listening on port %s\n", PORT);
	if err := http.ListenAndServe(PORT, mux); err != nil {
		log.Fatalf("ERROR: Failed to start HTTP server: %s\n", err.Error());
	}
}
