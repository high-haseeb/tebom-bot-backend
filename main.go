package main

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	godotenv.Load(".env");
	PORT := os.Getenv("PORT");

	mux := http.NewServeMux();

	mux.Handle("/get/vehicleInfo", Middleware(GetVehicleInformation));
	mux.Handle("/startOffers", Middleware(StartOffer));
	mux.Handle("/get/offers", Middleware(GetOffers));
	mux.Handle("/get/PDF", Middleware(GetPDF));

	log.Printf("INFO: Listening on port %s\n", PORT);
	if err := http.ListenAndServe(PORT, mux); err != nil {
		log.Fatalf("ERROR: Failed to start HTTP server: %s\n", err.Error());
	}
}
