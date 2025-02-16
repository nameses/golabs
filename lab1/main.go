package main

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

type PageData struct {
	Title  string
	Result template.HTML
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/task1", task1Handler)
	http.HandleFunc("/task2", task2Handler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/index.html")
	tmpl.Execute(w, nil)
}

func task1Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		hp, _ := strconv.ParseFloat(r.FormValue("hp"), 64)
		cp, _ := strconv.ParseFloat(r.FormValue("cp"), 64)
		sp, _ := strconv.ParseFloat(r.FormValue("sp"), 64)
		np, _ := strconv.ParseFloat(r.FormValue("np"), 64)
		op, _ := strconv.ParseFloat(r.FormValue("op"), 64)
		wp, _ := strconv.ParseFloat(r.FormValue("wp"), 64)
		ap, _ := strconv.ParseFloat(r.FormValue("ap"), 64)

		coefRC := 100.0 / (100.0 - wp)
		coefRG := 100.0 / (100.0 - wp - ap)
		total := hp + cp + sp + np + op + wp + ap

		percentagesDryMass := map[string]float64{
			"HP, %: ": hp * coefRC,
			"CP, %: ": cp * coefRC,
			"SP, %: ": sp * coefRC,
			"NP, %: ": np * coefRC,
			"OP, %: ": op * coefRC,
			"AP, %: ": ap * coefRC,
		}

		percentagesHotMass := map[string]float64{
			"HP, %: ": hp * coefRG,
			"CP, %: ": cp * coefRG,
			"SP, %: ": sp * coefRG,
			"NP, %: ": np * coefRG,
			"OP, %: ": op * coefRG,
		}

		qr := 339.0*cp + 1030.0*hp - 108.8*(op-sp) - 25.0*wp
		//з кДж/кг в МДж/кг
		qr /= 1000.0

		qc := (qr + 0.025*wp) * (100.0 / (100.0 - wp))
		qg := (qr + 0.025*wp) * (100.0 / (100.0 - wp - ap))

		var result string

		if total != 100.0 {
			result = "Сума компонентів повинна бути 100%"
		} else {
			result = "Результат обчислення:\n"
			result += "Суха маса:\n"

			for k, v := range percentagesDryMass {
				result += k + strconv.FormatFloat(v, 'f', 2, 64) + ", "
			}

			result += "\nГорюча маса:\n"

			for k, v := range percentagesHotMass {
				result += k + strconv.FormatFloat(v, 'f', 2, 64) + ", "
			}
			result += "\n"
			result += "Нижча теплота згоряння для робочої маси: " + strconv.FormatFloat(qr, 'f', 2, 64) + "\n"
			result += "Нижча теплота згоряння для сухої маси: " + strconv.FormatFloat(qc, 'f', 2, 64) + "\n"
			result += "Нижча теплота згоряння для горючої маси: " + strconv.FormatFloat(qg, 'f', 2, 64) + "\n"

			result = strings.ReplaceAll(result, "\n", "<br>")
		}

		tmpl, _ := template.ParseFiles("templates/task1.html")
		tmpl.Execute(w, PageData{Title: "Розрахунок складу", Result: template.HTML(result)})
	} else {
		tmpl, _ := template.ParseFiles("templates/task1.html")
		tmpl.Execute(w, PageData{Title: "Розрахунок складу", Result: template.HTML("")})
	}
}

func task2Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		carbon, _ := strconv.ParseFloat(r.FormValue("carbon"), 64)
		hydrogen, _ := strconv.ParseFloat(r.FormValue("hydrogen"), 64)
		oxygen, _ := strconv.ParseFloat(r.FormValue("oxygen"), 64)
		sulfur, _ := strconv.ParseFloat(r.FormValue("sulfur"), 64)
		calorific_value, _ := strconv.ParseFloat(r.FormValue("calorific_value"), 64)
		moisture, _ := strconv.ParseFloat(r.FormValue("moisture"), 64)
		ash, _ := strconv.ParseFloat(r.FormValue("ash"), 64)
		vanadium, _ := strconv.ParseFloat(r.FormValue("vanadium"), 64)

		carbonPR := carbon * (100 - moisture - ash) / 100
		hydrogenPR := hydrogen * (100 - moisture - ash) / 100
		oxygenPR := oxygen * (100 - moisture - ash) / 100
		sulfurPR := sulfur * (100 - moisture - ash) / 100
		ashPR := ash * (100 - moisture) / 100
		vanadiumPR := vanadium * (100 - moisture) / 100

		lowerHeatingValuePR := calorific_value*(100-moisture-ash)/100 - 0.025*moisture

		result := "Склад робочої маси мазуту:\n " +
			"Вуглець: " + strconv.FormatFloat(carbonPR, 'f', 2, 64) + "%\n" +
			"Водень: " + strconv.FormatFloat(hydrogenPR, 'f', 2, 64) + "%\n" +
			"Кисень: " + strconv.FormatFloat(oxygenPR, 'f', 2, 64) + "%\n" +
			"Сірка: " + strconv.FormatFloat(sulfurPR, 'f', 2, 64) + "%\n" +
			"Зола: " + strconv.FormatFloat(ashPR, 'f', 2, 64) + "%\n" +
			"Ванадій: " + strconv.FormatFloat(vanadiumPR, 'f', 2, 64) + " (мг/кг)\n" +
			"Нижча теплота згоряння: " + strconv.FormatFloat(lowerHeatingValuePR, 'f', 2, 64)

		result = strings.ReplaceAll(result, "\n", "<br>")

		tmpl, _ := template.ParseFiles("templates/task2.html")
		tmpl.Execute(w, PageData{Title: "Розрахунок перерахунку елементарного складу", Result: template.HTML(result)})
	} else {
		tmpl, _ := template.ParseFiles("templates/task2.html")
		tmpl.Execute(w, PageData{Title: "Розрахунок перерахунку елементарного складу", Result: template.HTML("")})
	}
}
