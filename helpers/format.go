package helpers

import (
	"fmt"
	"html/template"
	"strings"
	"time"
	"unicode/utf8"
)

func NormaliserSaisie(s string) string {
	return strings.TrimSpace(s)
}

func ValiderNomVille(nom string) string {
	if nom == "" {
		return "Veuillez saisir le nom d'une ville."
	}
	
	nbLettres := utf8.RuneCountInString(nom)
	if nbLettres < 2 {
		return "Veuillez saisir au moins 2 caractères."
	}
	if nbLettres > 60 {
		return "Le nom de la ville est trop long (60 caractères maximum)."
	}
	return ""
}

func formatDateFR(t time.Time) string {
	jours := []string{"dim.", "lun.", "mar.", "mer.", "jeu.", "ven.", "sam."}
	mois := []string{"janv.", "févr.", "mars", "avr.", "mai", "juin", "juil.", "août", "sept.", "oct.", "nov.", "déc."}

	j := jours[t.Weekday()]
	m := mois[t.Month()-1]
	return fmt.Sprintf("%s %d %s", j, t.Day(), m)
}

func formatTemp(v float64) string {
	txt := fmt.Sprintf("%.1f", v)
	txt = strings.ReplaceAll(txt, ".", ",")
	return txt + " °C"
}

func labelMeteo(code int) string {
	if code == 0 {
		return "Ciel dégagé"
	}
	if code >= 1 && code <= 3 {
		return "Peu nuageux, partiellement nuageux, couvert"
	}
	if code == 45 || code == 48 {
		return "Brouillard"
	}
	if code >= 51 && code <= 57 {
		return "Bruine"
	}
	if code >= 61 && code <= 67 {
		return "Pluie"
	}
	if code >= 71 && code <= 77 {
		return "Neige"
	}
	if code >= 80 && code <= 82 {
		return "Averses"
	}
	if code == 85 || code == 86 {
		return "Averses de neige"
	}
	if code >= 95 && code <= 99 {
		return "Orage"
	}
	return "Inconnu"
}

func iconeMeteo(code int) string {
	if code == 0 {
		return "☀️"
	}
	if code >= 1 && code <= 3 {
		return "⛅"
	}
	if code == 45 || code == 48 {
		return "🌫️"
	}
	if code >= 51 && code <= 57 {
		return "🌦️"
	}
	if code >= 61 && code <= 67 {
		return "🌧️"
	}
	if code >= 71 && code <= 77 {
		return "❄️"
	}
	if code >= 80 && code <= 82 {
		return "🌧️"
	}
	if code == 85 || code == 86 {
		return "🌨️"
	}
	if code >= 95 && code <= 99 {
		return "⛈️"
	}
	return "🌡️"
}

func styleBarre(left, width float64) template.HTMLAttr {
	g := strings.ReplaceAll(fmt.Sprintf("%.2f", left), ",", ".")
	l := strings.ReplaceAll(fmt.Sprintf("%.2f", width), ",", ".")
	return template.HTMLAttr("style=\"left: " + g + "%; width: " + l + "%;\"")
}

func FuncMap() template.FuncMap {
	fm := template.FuncMap{}
	fm["frDate"] = formatDateFR
	fm["temp"] = formatTemp
	fm["weatherLabel"] = labelMeteo
	fm["weatherIcon"] = iconeMeteo
	fm["barAttr"] = styleBarre
	return fm
}