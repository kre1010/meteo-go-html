package main

import (
	"log"
	"net/http"

	"meteo-go/api"
	"meteo-go/controllers"
	"meteo-go/routes"
	"meteo-go/services"
)

func main() {
	client := api.NouveauClientMeteo()
	service := services.NouveauMeteoService(client)

	ctrl, err := controllers.NouveauMeteoController(service)
	if err != nil {
		log.Println("erreur chargement templates:", err)
		return
	}

	mux := http.NewServeMux()
	routes.EnregistrerRoutes(mux, ctrl)

	log.Println("Serveur lance sur http://localhost:8080")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Println("erreur serveur:", err)
	}
}