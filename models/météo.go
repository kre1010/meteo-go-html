package models

import "time"

type Commune struct {
	Code         string    `json:"code"`
	Nom          string    `json:"nom"`
	CodesPostaux []string  `json:"codesPostaux"`
	Population   int       `json:"population"`
	Centre       CentreGeo `json:"centre"`
}

type CentreGeo struct {
	Coordinates []float64 `json:"coordinates"`
}

type JourMeteo struct {
	Date         time.Time
	TempMin      float64
	TempMax      float64
	WeatherCode  int
	LeftPercent  float64
	WidthPercent float64
}

type PrevisionsMeteo struct {
	Commune    Commune
	Jours      []JourMeteo
	MinSemaine float64
	MaxSemaine float64
	MapURL     string
}

type PageData struct {
	Titre      string
	Query      string
	Message    string
	Communes   []Commune
	Previsions *PrevisionsMeteo
	CodeHTTP   int
}