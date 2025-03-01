package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"math"
	"net/http"
)

type CalculationRequest struct {
	Sm     float64 `json:"sm"`
	Ik     float64 `json:"ik"`
	Tf     float64 `json:"tf"`
	Power  float64 `json:"power"`
	Rsn    float64 `json:"rsn"`
	Xsn    float64 `json:"xsn"`
	RsnMin float64 `json:"rsnMin"`
	XsnMin float64 `json:"xsnMin"`
}

type CalculationResponse struct {
	Im   float64 `json:"im"`
	ImPa float64 `json:"imPa"`
	SEk  float64 `json:"sEk"`
	SVsS float64 `json:"sVsS"`
	Xc   float64 `json:"xc"`
	Xt   float64 `json:"xt"`
	Ip0  float64 `json:"ip0"`
	Ish3 float64 `json:"ish3"`
	Ish2 float64 `json:"ish2"`
}

func calculateCables(sm, ik, tf float64) CalculationResponse {
	im := sm / (2.0 * 1.732 * 10.0)
	imPa := 2.0 * im
	sEk := im / 1.4
	sVsS := ik * 1.732 * tf / 92
	return CalculationResponse{Im: im, ImPa: imPa, SEk: sEk, SVsS: sVsS}
}

func calculateBusbar(power float64) CalculationResponse {
	xc := math.Pow(10.5, 2) / power
	xt := (10.5 / 100) * (math.Pow(10.5, 2) / 6.3)
	ip0 := 10.5 / (math.Sqrt(3.0) * (xc + xt))
	return CalculationResponse{Xc: xc, Xt: xt, Ip0: ip0}
}

func calculateSubstation(rsn, xsn, rsnMin, xsnMin float64) CalculationResponse {
	ish3 := (115.0 * 1000) / (math.Sqrt(3.0) * math.Sqrt(rsn*rsn+xsn*xsn))
	ish2 := ish3 * (math.Sqrt(3.0) / 2)
	return CalculationResponse{Ish3: ish3, Ish2: ish2}
}

func handleCalculate(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/api/calc1":
		var req CalculationRequest
		json.NewDecoder(r.Body).Decode(&req)
		res := calculateCables(req.Sm, req.Ik, req.Tf)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	case "/api/calc2":
		var req CalculationRequest
		json.NewDecoder(r.Body).Decode(&req)
		res := calculateBusbar(req.Power)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	case "/api/calc3":
		var req CalculationRequest
		json.NewDecoder(r.Body).Decode(&req)
		res := calculateSubstation(req.Rsn, req.Xsn, req.RsnMin, req.XsnMin)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	}
}

func handlePage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/" + r.URL.Path[1:] + ".html")
	tmpl.Execute(w, nil)
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/calc1", handlePage)
	http.HandleFunc("/calc2", handlePage)
	http.HandleFunc("/calc3", handlePage)
	http.HandleFunc("/api/calc1", handleCalculate)
	http.HandleFunc("/api/calc2", handleCalculate)
	http.HandleFunc("/api/calc3", handleCalculate)

	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
