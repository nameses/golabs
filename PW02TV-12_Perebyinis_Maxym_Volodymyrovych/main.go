package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type PageData struct {
	HeaderTitle    string
	EmissionOutput template.HTML
}

type EmissionsResult struct {
	CoalEmission  float64
	MazutEmission float64
	GasEmission   float64
}

const (
	CoalFactor  = 2.86
	MazutFactor = 3.07
	GasFactor   = 2.02
)

func CalculateEmissions(coal, mazut, gas float64) EmissionsResult {
	const (
		coalFactor   = 150.0
		coalEmission = 20.47
		oilFactor    = 0.57
		oilEmission  = 40.40 * 0.98
	)

	return EmissionsResult{
		CoalEmission:  coalFactor * coal * coalEmission / 1e6,
		MazutEmission: oilFactor * mazut * oilEmission / 1e6,
		GasEmission:   0.0,
	}
}

func formatResult(result EmissionsResult) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Викид вугілля: %.2f тонн<br>", result.CoalEmission))
	sb.WriteString(fmt.Sprintf("Викид мазуту: %.2f тонн<br>", result.MazutEmission))
	sb.WriteString(fmt.Sprintf("Викид природного газу: %.2f тонн", result.GasEmission))
	return sb.String()
}

func PageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/main.html")
	if err != nil {
		http.Error(w, "Помилка завантаження сторінки", http.StatusInternalServerError)
		return
	}

	data := PageData{
		HeaderTitle: "Розрахунок шкідливих викидів",
	}

	if r.Method == http.MethodPost {
		coalInput := r.FormValue("coal_input")
		mazutInput := r.FormValue("mazut_input")
		gasInput := r.FormValue("gas_input")

		coal, _ := strconv.ParseFloat(coalInput, 64)
		mazut, _ := strconv.ParseFloat(mazutInput, 64)
		gas, _ := strconv.ParseFloat(gasInput, 64)

		emissionTotal := CalculateEmissions(coal, mazut, gas)
		result := formatResult(emissionTotal)
		data.EmissionOutput = template.HTML(result)
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Помилка відображення сторінки", http.StatusInternalServerError)
	}
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", PageHandler)

	fmt.Println("Сервер запущено на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
