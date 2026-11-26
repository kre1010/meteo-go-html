package routes

import (
	"net/http"

	"meteo-go/controllers"
)

func EnregistrerRoutes(mux *http.ServeMux, ctrl *controllers.MeteoController) {
	dossierAssets := http.Dir("assets")
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(dossierAssets)))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		chemin := r.URL.Path

		if chemin == "/" {
			ctrl.Accueil(w, r)
			return
		}
		if chemin == "/ville" {
			ctrl.Ville(w, r)
			return
		}
		if chemin == "/meteo" {
			ctrl.Meteo(w, r)
			return
		}

		ctrl.AfficherErreur(w, http.StatusNotFound, "Cette page n'existe pas.")
	})
}