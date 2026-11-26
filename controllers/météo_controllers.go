	package controllers

import (
	"errors"
	"html/template"
	"log"
	"net/http"

	"meteo-go/helpers"
	"meteo-go/models"
	"meteo-go/services"
)

type MeteoController struct {
	service   *services.MeteoService
	templates *template.Template
}

func NouveauMeteoController(service *services.MeteoService) (*MeteoController, error) {
   
	tmpl, err := template.New("").Funcs(helpers.FuncMap()).ParseGlob("templates/*.html")
	if err != nil {
		return nil, err
	}

	ctrl := &MeteoController{
		service:   service,
		templates: tmpl,
	}
	return ctrl, nil
}

func (c *MeteoController) AfficherErreur(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)

	data := models.PageData{
		Titre:    "Erreur",
		Message:  message,
		CodeHTTP: status,
	}

	err := c.templates.ExecuteTemplate(w, "erreur.html", data)
	if err != nil {
		log.Println("erreur template erreur.html:", err)
	}
}

func (c *MeteoController) Accueil(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		c.AfficherErreur(w, http.StatusNotFound, "Cette page n'existe pas.")
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	data := models.PageData{
		Titre:    "Accueil - Météo",
		CodeHTTP: http.StatusOK,
	}

	err := c.templates.ExecuteTemplate(w, "accueil.html", data)
	if err != nil {
		log.Println("erreur template accueil.html:", err)
	}
}

func (c *MeteoController) Ville(w http.ResponseWriter, r *http.Request) {
	q := helpers.NormaliserSaisie(r.URL.Query().Get("q"))
	msgErreur := helpers.ValiderNomVille(q)

	if msgErreur != "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)

		data := models.PageData{
			Titre:    "Recherche de commune",
			Query:    q,
			Message:  msgErreur,
			CodeHTTP: http.StatusBadRequest,
		}

		err := c.templates.ExecuteTemplate(w, "accueil.html", data)
		if err != nil {
			log.Println("erreur template accueil.html:", err)
		}
		return
	}

	communes, err := c.service.ChercherCommunes(q)
	if err != nil {
		log.Println("erreur recherche commune:", err)
		c.AfficherErreur(w, http.StatusServiceUnavailable, "Service momentanément indisponible.")
		return
	}

	if len(communes) == 0 {
		c.AfficherErreur(w, http.StatusNotFound, "Aucune commune ne correspond à « "+q+" ».")
		return
	}

	if len(communes) == 1 {
		c.afficherMeteoCode(w, communes[0].Code)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	data := models.PageData{
		Titre:    "Choisir une commune",
		Query:    q,
		Communes: communes,
		CodeHTTP: http.StatusOK,
	}

	err = c.templates.ExecuteTemplate(w, "villes.html", data)
	if err != nil {
		log.Println("erreur template villes.html:", err)
	}
}

func (c *MeteoController) Meteo(w http.ResponseWriter, r *http.Request) {
	code := helpers.NormaliserSaisie(r.URL.Query().Get("code"))
	if code == "" {
		c.AfficherErreur(w, http.StatusNotFound, "Cette commune n'existe pas.")
		return
	}

	c.afficherMeteoCode(w, code)
}

func (c *MeteoController) afficherMeteoCode(w http.ResponseWriter, code string) {
	prev, err := c.service.ObtenirPrevisionsParCode(code)
	if err != nil {
		if errors.Is(err, services.ErrCommuneIntrouvable) {
			c.AfficherErreur(w, http.StatusNotFound, "Cette commune n'existe pas.")
			return
		}
		log.Println("erreur meteo code", code, ":", err)
		c.AfficherErreur(w, http.StatusServiceUnavailable, "Service momentanément indisponible.")
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	data := models.PageData{
		Titre:      "Météo - " + prev.Commune.Nom,
		Previsions: prev,
		CodeHTTP:   http.StatusOK,
	}

	err = c.templates.ExecuteTemplate(w, "meteo.html", data)
	if err != nil {
		log.Println("erreur template meteo.html:", err)
	}
}