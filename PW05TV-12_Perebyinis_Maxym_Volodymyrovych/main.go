package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

type CalculationRequest struct {
	LossesEmergency     float64 `json:"lossesEmergency"`
	LossesPlanned       float64 `json:"lossesPlanned"`
	QuantityPl110kV     float64 `json:"quantityPl110kV"`
	QuantityAttachments float64 `json:"quantityAttachments"`
}

type CalculationResponse struct {
	LossExpectation float64 `json:"lossExpectation"`
	FailureRate1    float64 `json:"failureRate1"`
	FailureRate2    float64 `json:"failureRate2"`
	ComparisonText  string  `json:"comparisonText"`
}

func calculateLosses(lossesEmergency, lossesPlanned float64) float64 {
	mEmergency := 0.01 * 0.045 * 5120 * 6451
	mPlanned := 0.004 * 5120 * 6451
	return lossesEmergency*mEmergency + lossesPlanned*mPlanned
}

func calculateComparison(quantityPl110kV, quantityAttachments float64) (float64, float64, string) {
	failureRate1 := 0.01 + 0.007*quantityPl110kV + 0.015 + 0.02 + 0.03*quantityAttachments
	averageRecoveryDuration := (0.01*30 + 0.007*quantityPl110kV*10 + 0.015*100 + 0.02*15 + 0.03*quantityAttachments*2) / failureRate1
	ka := (failureRate1 * averageRecoveryDuration) / 8760
	kp := 1.2 * (43.0 / 8760.0)
	failureRate2 := 2*failureRate1*(ka+kp) + 0.02

	comparisonText := ""
	if failureRate1 < failureRate2 {
		comparisonText = "Одноколова система має вищу надійність."
	} else if failureRate1 > failureRate2 {
		comparisonText = "Двоколова система має вищу надійність."
	} else {
		comparisonText = "Надійність обох систем однакова."
	}

	return failureRate1, failureRate2, comparisonText
}

func handleCalculate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не дозволений", http.StatusMethodNotAllowed)
		return
	}

	var req CalculationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Помилка вхідних даних", http.StatusBadRequest)
		fmt.Println("JSON Decode Error:", err)
		return
	}

	var res CalculationResponse
	switch r.URL.Path {
	case "/api/losses":
		res.LossExpectation = calculateLosses(req.LossesEmergency, req.LossesPlanned)
	case "/api/comparing":
		res.FailureRate1, res.FailureRate2, res.ComparisonText = calculateComparison(req.QuantityPl110kV, req.QuantityAttachments)
	default:
		http.Error(w, "Невірний endpoint", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func handlePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/" + r.URL.Path[1:] + ".html")
	if err != nil {
		http.Error(w, "Помилка завантаження сторінки", http.StatusInternalServerError)
		fmt.Println("Помилка завантаження шаблону:", err)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Помилка рендерингу сторінки", http.StatusInternalServerError)
		fmt.Println("Помилка рендерингу шаблону:", err)
	}
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/comparing", handlePage)
	http.HandleFunc("/losses", handlePage)
	http.HandleFunc("/api/comparing", handleCalculate)
	http.HandleFunc("/api/losses", handleCalculate)

	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
