package main

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
)

type Results struct {
	Coal       string
	OilFuel    string
	NaturalGas string

	EmissionCoal      int
	GrossEmissionCoal int

	EmissionOilFuel      string
	GrossEmissionOilFuel string

	EmissionNaturalGas      int
	GrossEmissionNaturalGas int

	Submitted bool
}

func emissionCoal() int {
	return int(math.Round((math.Pow(10, 6) / 20.47) * 0.8 * (25.2 / (100 - 1.5)) * (1 - 0.985)))
}

func grossEmissionCoal(k int, coal float64) int {
	return int(math.Round(math.Pow(10, -6) * float64(k) * 20.47 * coal))
}

func emissionOilFuel() float64 {
	return (math.Pow(10, 6) / 39.48) * 1 * (0.15 / (100 - 0)) * (1 - 0.985)
}

func grossEmissionOilFuel(k float64, oilFuel float64) float64 {
	return math.Pow(10, -6) * k * 39.48 * oilFuel
}

func grossEmissionNaturalGas(k int, gas float64) int {
	return int(math.Round(math.Pow(10, -6) * float64(k) * 33.08 * gas))
}

func formatFloat2(f float64) string {
	return fmt.Sprintf("%.2f", f)
}

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		log.Println("Template parse error:", err)
		return
	}

	data := Results{}

	if r.Method == http.MethodPost {
		coalStr := r.FormValue("coal")
		oilFuelStr := r.FormValue("oil-fuel")
		naturalGasStr := r.FormValue("natural-gas")

		coal, _ := strconv.ParseFloat(coalStr, 64)
		oilFuel, _ := strconv.ParseFloat(oilFuelStr, 64)
		naturalGas, _ := strconv.ParseFloat(naturalGasStr, 64)

		calcEmissionCoal := emissionCoal()
		calcGrossEmissionCoal := grossEmissionCoal(calcEmissionCoal, coal)

		calcEmissionOilFuel := emissionOilFuel()
		calcEmissionOilFuelRounded, _ := strconv.ParseFloat(formatFloat2(calcEmissionOilFuel), 64)
		calcGrossEmissionOilFuel := grossEmissionOilFuel(calcEmissionOilFuelRounded, oilFuel)

		valueEmissionNaturalGas := 0
		calcGrossEmissionNaturalGas := grossEmissionNaturalGas(valueEmissionNaturalGas, naturalGas)

		data = Results{
			Coal:       coalStr,
			OilFuel:    oilFuelStr,
			NaturalGas: naturalGasStr,

			EmissionCoal:      calcEmissionCoal,
			GrossEmissionCoal: calcGrossEmissionCoal,

			EmissionOilFuel:      formatFloat2(calcEmissionOilFuel),
			GrossEmissionOilFuel: formatFloat2(calcGrossEmissionOilFuel),

			EmissionNaturalGas:      valueEmissionNaturalGas,
			GrossEmissionNaturalGas: calcGrossEmissionNaturalGas,

			Submitted: true,
		}
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Println("Template execute error:", err)
	}
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Сервер запущено на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
