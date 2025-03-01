package main

import (
	"encoding/json"
	"math"
	"net/http"
)

type InputDistributionTire struct {
	Name string
	NomK float64
	Cos  float64
	Un   float64
	N    int
	Pn   int
	Kv   float64
	Tg   float64
}

type InputDistributionTireTotal struct {
	N     int
	NP    int
	NPK   int
	NPKtg int
	NP2   int
}

type OutputDistributionTire struct {
	Name  string
	NP    int
	NPK   float64
	NPKtg float64
	NP2   int
	Ip    float64
}

type OutputDistributionTireTotal struct {
	K   float64
	NEf float64
	Ka  float64
	Ra  float64
	Rr  float64
	Sr  float64
	Ip  float64
}

type TotalDistributionTire struct {
	Name  string
	N     int
	NP    int
	K     float64
	NPK   float64
	NPKtg float64
	NP2   int
	NEf   float64
	Ka    float64
	Ra    float64
	Rr    float64
	Sr    float64
	Ip    float64
}

func round(num float64) float64 {
	return math.Round(num*100) / 100
}

func calculateDistributionTire(input InputDistributionTire) OutputDistributionTire {
	nP := input.N * input.Pn
	nPK := float64(nP) * input.Kv
	nPKtg := nPK * input.Tg
	nP2 := input.N * input.Pn * input.Pn
	I := float64(nP) / (math.Sqrt(3) * input.Un * input.Cos * input.NomK)

	return OutputDistributionTire{
		Name:  input.Name,
		NP:    nP,
		NPK:   round(nPK),
		NPKtg: round(nPKtg),
		NP2:   nP2,
		Ip:    round(I),
	}
}

func calculateDistributionTireTotal(inputs []InputDistributionTire, outputs []OutputDistributionTire) TotalDistributionTire {
	nP := 0
	nPK := 0.0
	nPKtg := 0.0
	nP2 := 0

	for _, output := range outputs {
		nP += output.NP
		nPK += output.NPK
		nPKtg += output.NPKtg
		nP2 += output.NP2
	}

	K := nPK / float64(nP)
	nEf := math.Pow(float64(nP), 2) / float64(nP2)
	Ra := 1.25 * nPK
	S := math.Sqrt(math.Pow(Ra, 2) + math.Pow(nPKtg, 2))
	I := Ra / 0.38

	return TotalDistributionTire{
		Name:  "ШР",
		N:     len(inputs),
		NP:    nP,
		K:     round(K),
		NPK:   round(nPK),
		NPKtg: round(nPKtg),
		NP2:   nP2,
		NEf:   round(nEf),
		Ka:    1.25,
		Ra:    round(Ra),
		Rr:    round(nPKtg),
		Sr:    round(S),
		Ip:    round(I),
	}
}

func calculateTotal(input InputDistributionTireTotal) OutputDistributionTireTotal {
	K := float64(input.NPK) / float64(input.NP)
	nEf := math.Pow(float64(input.NP), 2) / float64(input.NP2)
	Ra := 0.7 * float64(input.NPK)
	Rr := 0.7 * float64(input.NPKtg)
	S := math.Sqrt(math.Pow(Ra, 2) + math.Pow(Rr, 2))
	I := Ra / 0.38

	return OutputDistributionTireTotal{
		K:   round(K),
		NEf: round(nEf),
		Ka:  0.7,
		Ra:  Ra,
		Rr:  Rr,
		Sr:  round(S),
		Ip:  round(I),
	}
}

func main() {
	http.HandleFunc("/results", func(w http.ResponseWriter, r *http.Request) {
		inputDistributionTires := []InputDistributionTire{
			{"Шліфувальний верстат", 0.92, 0.9, 0.38, 4, 21, 0.15, 1.33},
			{"Свердлильний верстат", 0.92, 0.9, 0.38, 2, 14, 0.12, 1.0},
			{"Фугувальний верстат", 0.92, 0.9, 0.38, 4, 42, 0.15, 1.33},
			{"Циркулярна пила", 0.92, 0.9, 0.38, 1, 36, 0.3, 1.56},
			{"Прес", 0.92, 0.9, 0.38, 1, 20, 0.5, 0.75},
			{"Полірувальний верстат", 0.92, 0.9, 0.38, 1, 40, 0.22, 1.0},
			{"Фрезерний верстат", 0.92, 0.9, 0.38, 2, 32, 0.2, 1.0},
			{"Вентилятор", 0.92, 0.9, 0.38, 1, 20, 0.65, 0.75},
		}

		largeEP1 := InputDistributionTire{"Зварювальний трансформатор", 0.92, 0.9, 0.38, 2, 100, 0.2, 3.0}
		largeEP2 := InputDistributionTire{"Сушильна шафа", 0.92, 0.9, 0.38, 2, 120, 0.8, 0.0}

		inputDistributionTireTotal := InputDistributionTireTotal{81, 2330, 752, 657, 96399}

		results := make([]OutputDistributionTire, 0)
		for _, input := range inputDistributionTires {
			results = append(results, calculateDistributionTire(input))
		}

		totalDistributionTiresResult := calculateDistributionTireTotal(inputDistributionTires, results)
		totalOutputResult := calculateTotal(inputDistributionTireTotal)
		result1 := calculateDistributionTire(largeEP1)
		result2 := calculateDistributionTire(largeEP2)

		response := map[string]interface{}{
			"individualResults":            results,
			"totalDistributionTiresResult": totalDistributionTiresResult,
			"result1":                      result1,
			"result2":                      result2,
			"totalOutputResult":            totalOutputResult,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.ListenAndServe(":8080", nil)
}
