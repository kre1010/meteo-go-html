package services

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"meteo-go/api"
	"meteo-go/models"
)

var (
	ErrCommuneIntrouvable = errors.New("commune introuvable")
	ErrDonneesIncompletes = errors.New("données météo incomplètes")
)

type MeteoService struct {
	client *api.ClientMeteo
}

func NouveauMeteoService(client *api.ClientMeteo) *MeteoService {
	return &MeteoService{client: client}
}

func (s *MeteoService) ChercherCommunes(nom string) ([]models.Commune, error) {
	return s.client.RechercherCommunesParNom(nom)
}

func (s *MeteoService) ObtenirPrevisionsParCode(codeINSEE string) (*models.PrevisionsMeteo, error) {
	commune, err := s.client.RechercherCommuneParCode(codeINSEE)
	if err != nil {
		return nil, err
	}
	if commune == nil || commune.Code == "" || len(commune.Centre.Coordinates) < 2 {
		return nil, ErrCommuneIntrouvable
	}

	lon := commune.Centre.Coordinates[0]
	lat := commune.Centre.Coordinates[1]

	raw, err := s.client.RecupererPrevisions7Jours(lat, lon)
	if err != nil {
		return nil, err
	}

	nbJours := len(raw.Daily.Time)
	if nbJours == 0 ||
		len(raw.Daily.TemperatureMin) != nbJours ||
		len(raw.Daily.TemperatureMax) != nbJours ||
		len(raw.Daily.WeatherCode) != nbJours {
		return nil, ErrDonneesIncompletes
	}

	var jours []models.JourMeteo
	minSemaine := raw.Daily.TemperatureMin[0]
	maxSemaine := raw.Daily.TemperatureMax[0]

	for i := 0; i < nbJours; i++ {
		dateParsed, err := time.Parse("2006-01-02", raw.Daily.Time[i])
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDonneesIncompletes, err)
		}

		tMin := raw.Daily.TemperatureMin[i]
		tMax := raw.Daily.TemperatureMax[i]

		if tMin < minSemaine {
			minSemaine = tMin
		}
		if tMax > maxSemaine {
			maxSemaine = tMax
		}

		jour := models.JourMeteo{
			Date:        dateParsed,
			TempMin:     tMin,
			TempMax:     tMax,
			WeatherCode: raw.Daily.WeatherCode[i],
		}
		jours = append(jours, jour)
	}

	amplitude := maxSemaine - minSemaine
	for i := 0; i < len(jours); i++ {
		if amplitude == 0 {
			jours[i].LeftPercent = 0
			jours[i].WidthPercent = 100
		} else {
			jours[i].LeftPercent = ((jours[i].TempMin - minSemaine) / amplitude) * 100
			jours[i].WidthPercent = ((jours[i].TempMax - jours[i].TempMin) / amplitude) * 100
		}
	}

	d := 0.05
	bbox := fmt.Sprintf("%f,%f,%f,%f", lon-d, lat-d, lon+d, lat+d)
	marker := fmt.Sprintf("%f,%f", lat, lon)
	mapURL := "https://www.openstreetmap.org/export/embed.html?bbox=" + url.QueryEscape(bbox) + "&layer=mapnik&marker=" + url.QueryEscape(marker)

	res := &models.PrevisionsMeteo{
		Commune:    *commune,
		Jours:      jours,
		MinSemaine: minSemaine,
		MaxSemaine: maxSemaine,
		MapURL:     mapURL,
	}
	return res, nil
}