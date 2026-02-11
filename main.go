package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type PageData struct {
	FuelInput        *FuelInput
	DryMass          *DryMassResult
	CombustibleMass  *CombustibleMassResult
	HeatCombustion   *HeatCombustionResult
	ShowFuelResult   bool

	FuelOilInput       *FuelOilInput
	FuelOilComposition *FuelOilCompositionResult
	FuelOilHeat        float64
	ShowFuelOilResult  bool
}

var tmpl *template.Template

func main() {
	funcMap := template.FuncMap{
		"f2": func(v float64) string {
			return fmt.Sprintf("%.2f", v)
		},
		"fg": func(v float64) string {
			return strconv.FormatFloat(v, 'f', -1, 64)
		},
	}

	tmpl = template.Must(
		template.New("index.html").Funcs(funcMap).ParseFiles("templates/index.html"),
	)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", handleIndex)

	fmt.Println("Сервер запущено на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	data := &PageData{}

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		switch r.FormValue("calculator") {
		case "fuel":
			input := &FuelInput{
				Hydrogen: parseFloat(r.FormValue("hydrogen")),
				Carbon:   parseFloat(r.FormValue("carbon")),
				Sulfur:   parseFloat(r.FormValue("sulfur")),
				Nitrogen: parseFloat(r.FormValue("nitrogen")),
				Oxygen:   parseFloat(r.FormValue("oxygen")),
				Moisture: parseFloat(r.FormValue("moisture")),
				Ash:      parseFloat(r.FormValue("ash")),
			}
			data.FuelInput = input
			data.ShowFuelResult = true
			data.DryMass = calculateDryMass(input)
			data.CombustibleMass = calculateCombustibleMass(input)
			data.HeatCombustion = calculateFuelHeatCombustion(input)

		case "fuel-oil":
			input := &FuelOilInput{
				Carbon:         parseFloat(r.FormValue("carbon-fuel-oil")),
				Hydrogen:       parseFloat(r.FormValue("hydrogen-fuel-oil")),
				Sulfur:         parseFloat(r.FormValue("sulfur-fuel-oil")),
				Vanadium:       parseFloat(r.FormValue("vanadi-fuel-oil")),
				Oxygen:         parseFloat(r.FormValue("oxygen-fuel-oil")),
				Moisture:       parseFloat(r.FormValue("moisture-fuel-oil")),
				Ash:            parseFloat(r.FormValue("ash-fuel-oil")),
				HeatCombustion: parseFloat(r.FormValue("lower-heat-combustion")),
			}
			data.FuelOilInput = input
			data.ShowFuelOilResult = true
			data.FuelOilComposition = calculateFuelOilComposition(input)
			data.FuelOilHeat = calculateFuelOilHeatCombustion(input)
		}
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Template error:", err)
	}
}

func parseFloat(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}
