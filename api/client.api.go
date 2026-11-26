package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"meteo-go/models"
)

var ErrAPIIndisponible = errors.New("service momentané indisponible")

type ClientMeteo struct {
	httpClient *http.Client
}

func NouveauClientMeteo() *ClientMeteo {
	return &ClientMeteo{
		httpClient: &http.Client{
			Timeout: 8 * time.Second,
		},
	}
}

func (c *ClientMeteo) RechercherCommunesParNom(nom string) ([]models.Commune, error) {
	
	nomEnc := url.QueryEscape(nom)
	fullURL := "https://geo.api.gouv.fr/communes?nom=" + nomEnc + "&fields=code,nom,codesPostaux,population,centre&boost=population&limit=10"

	resp, err := c.httpClient.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrAPIIndisponible, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: code %d", ErrAPIIndisponible, resp.StatusCode)
	}

	var communes []models.Commune
	err = json.NewDecoder(resp.Body).Decode(&communes)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrAPIIndisponible, err)
	}

	return communes, nil
}

func (c *ClientMeteo) RechercherCommuneParCode(codeINSEE string) (*models.Commune, error) {
	fullURL := "https://geo.api.gouv.fr/communes/" + url.PathEscape(codeINSEE) + "?fields=code,nom,codesPostaux,population,centre"

	resp, err := c.httpClient.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrAPIIndisponible, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: code %d", ErrAPIIndisponible, resp.StatusCode)
	}

	var commune models.Commune
	err = json.NewDecoder(resp.Body).Decode(&commune)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrAPIIndisponible, err)
	}

	return &commune, nil
}

type openMeteoResponse struct {
	Daily struct {
		Time           []string  `json:"time"`
		TemperatureMin []float64 `json:"temperature_2m_min"`
		TemperatureMax []float64 `json:"temperature_2m_max"`
		WeatherCode    []int     `json:"weather_code"`
	} `json:"daily"`
}

func (c *ClientMeteo) RecupererPrevisions7Jours(lat, lon float64) (*openMeteoResponse, error) {
	
	fullURL := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&daily=temperature_2m_min,temperature_2m_max,weather_code&forecast_days=7&timezone=Europe/Paris",
		lat,
		lon,
	)

	resp, err := c.httpClient.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrAPIIndisponible, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: code %d", ErrAPIIndisponible, resp.StatusCode)
	}

	var data openMeteoResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrAPIIndisponible, err)
	}

	return &data, nil
}