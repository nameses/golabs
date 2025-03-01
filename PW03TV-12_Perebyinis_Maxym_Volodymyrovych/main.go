package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
)

type CalculationRequest struct {
	Power          float64 `json:"power"`
	DeviationCurr  float64 `json:"deviation_current"`
	DeviationFinal float64 `json:"deviation_final"`
	Price          float64 `json:"price"`
}

type Results struct {
	IncomeCurrent  float64 `json:"income_current"`
	PenaltyCurrent float64 `json:"penalty_current"`
	IncomeAfter    float64 `json:"income_after"`
	PenaltyAfter   float64 `json:"penalty_after"`
	IncomeFinal    float64 `json:"income_final"`
}

func calculateResults(req CalculationRequest) Results {
	deltaWCurrent := integrate(req.Power, req.DeviationCurr)
	incomeCurrent := (req.Power * 24 * deltaWCurrent) * req.Price
	fineCurrent := (req.Power * 24 * (1 - deltaWCurrent)) * req.Price

	deltaWAfter := integrate(req.Power, req.DeviationFinal)
	incomeAfter := (req.Power * 24 * deltaWAfter) * req.Price
	fineAfter := (req.Power * 24 * (1 - deltaWAfter)) * req.Price
	incomeFinal := incomeAfter - fineAfter

	return Results{incomeCurrent, fineCurrent, incomeAfter, fineAfter, incomeFinal}
}

func integrate(averagePower, deviation float64) float64 {
	lowerLimit := 4.75
	upperLimit := 5.25
	steps := 10000
	stepSize := (upperLimit - lowerLimit) / float64(steps)
	sum := 0.0
	for i := 0; i < steps; i++ {
		x1 := lowerLimit + float64(i)*stepSize
		x2 := x1 + stepSize
		sum += 0.5 * (calculatePd(x1, averagePower, deviation) + calculatePd(x2, averagePower, deviation)) * stepSize
	}

	return sum
}

func calculatePd(p, averagePower, deviation float64) float64 {
	return (1 / (deviation * math.Sqrt(2*math.Pi))) * math.Exp(-math.Pow(p-averagePower, 2)/(2*math.Pow(deviation, 2)))
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req CalculationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	result := calculateResults(req)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/calculate", handler)

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}
