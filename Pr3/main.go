package main

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
)

type EPRow struct {
	Name string

	Eta  float64 // ηн (не використовується в формулах — як і в твоєму JS)
	Cos  float64 // cosφ (не використовується в формулах — як і в твоєму JS)
	U    float64 // Uн (беремо U тільки з першого рядка, як у JS)
	N    float64 // n
	Pn   float64 // Pн
	Kv   float64 // Kв
	Tg   float64 // tgφ
}

type Results struct {
	// Inputs (щоб зберігати введені значення після submit)
	EP1 EPRow
	EP2 EPRow
	EP3 EPRow

	Submitted bool

	// Вивід (форматований як у твоєму JS toFixed/ceil)
	GroupUtilizationFactor string
	EffectiveAmount        string
	ActivePowerFactor      string
	CalculatedActiveLoad   string
	CalculatedReactiveLoad string
	FullPower              string
	GroupCurrent           string

	UtilizationRatesWorkshop     string
	EffectiveNumberEPWorkshop    string
	CoefficientActiveApacityWork string

	ActiveTireLoad       string
	ReactiveTireLoad     string
	FullPowerTires       string
	GroupCurrentTires    string
}

var tpl = template.Must(template.ParseFiles("templates/index.html"))

func main() {
	mux := http.NewServeMux()

	// статика
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	mux.HandleFunc("/", handleIndex)

	addr := ":8080"
	log.Println("Server started on http://localhost" + addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		_ = tpl.Execute(w, Results{})
		return

	case http.MethodPost:
		res := Results{Submitted: true}

		res.EP1 = readEPRow(r, "")
		res.EP2 = readEPRow(r, "-2")
		res.EP3 = readEPRow(r, "-3")

		// Як у твоєму JS:
		// u = Number(inputLoadVoltage1.value);
		u := res.EP1.U

		activePowerFactor := 1.25
		coefficientActiveApacityWorkshop := 0.7

		// calculatedMultip = n * p
		calculatedMultip1 := multiplicationNandP(res.EP1.N, res.EP1.Pn)
		calculatedMultip2 := multiplicationNandP(res.EP2.N, res.EP2.Pn)
		calculatedMultip3 := multiplicationNandP(res.EP3.N, res.EP3.Pn)

		// calculatedMultipInSquare = n * p^2
		calculatedMultipInSquare1 := sumMultiplicationNandPInSquare(res.EP1.N, res.EP1.Pn)
		calculatedMultipInSquare2 := sumMultiplicationNandPInSquare(res.EP2.N, res.EP2.Pn)
		calculatedMultipInSquare3 := sumMultiplicationNandPInSquare(res.EP3.N, res.EP3.Pn)

		// calculatedMultipK = (n*p) * k
		calculatedMultipK1 := multiplicationK(calculatedMultip1, res.EP1.Kv)
		calculatedMultipK2 := multiplicationK(calculatedMultip2, res.EP2.Kv)
		calculatedMultipK3 := multiplicationK(calculatedMultip3, res.EP3.Kv)

		// multiplicationTg: tg * calculatedMultipK, rounded to 1 decimal like toFixed(1), then parsed back to Number
		calculatedMultipTg1 := multiplicationTgRounded1(res.EP1.Tg, calculatedMultipK1)
		calculatedMultipTg2 := multiplicationTgRounded1(res.EP2.Tg, calculatedMultipK2)
		calculatedMultipTg3 := multiplicationTgRounded1(res.EP3.Tg, calculatedMultipK3)

		calculatedGroupUtilizationFactor := getGroupUtilizationFactor(
			calculatedMultip1, calculatedMultip2, calculatedMultip3,
			calculatedMultipK1, calculatedMultipK2, calculatedMultipK3,
		)

		calculatedEffectiveAmount := getEffectiveAmount(
			calculatedMultip1, calculatedMultip2, calculatedMultip3,
			calculatedMultipInSquare1, calculatedMultipInSquare2, calculatedMultipInSquare3,
		)

		calculatedActiveLoad := getcalculatedActiveLoad(activePowerFactor, calculatedMultipK1, calculatedMultipK2, calculatedMultipK3)
		calculatedReactiveLoad := getcalculatedReactiveLoad(activePowerFactor, calculatedMultipTg1, calculatedMultipTg2, calculatedMultipTg3)
		calculatedFullPower := getFullPower(calculatedActiveLoad, calculatedReactiveLoad)
		calculatedGroupCurrent := getGroupCurrent(calculatedActiveLoad, u)

		calculatedUtilizationRatesWorkshop := getUtilizationRatesWorkshop(
			calculatedMultip1, calculatedMultip2, calculatedMultip3,
			calculatedMultipK1, calculatedMultipK2, calculatedMultipK3,
		)

		calculatedEffectiveNumberEPWorkshop := getEffectiveNumberEPWorkshop(
			calculatedMultip1, calculatedMultip2, calculatedMultip3,
			calculatedMultipInSquare1, calculatedMultipInSquare2, calculatedMultipInSquare3,
		)

		calculatedActiveTireLoad := getActiveTireLoad(coefficientActiveApacityWorkshop)
		calculatedReactiveLoadTires := getReactiveLoadTires(coefficientActiveApacityWorkshop)
		calculatedFullPowerTires := getFullPowerTires(calculatedActiveTireLoad, calculatedReactiveLoadTires)
		calculatedGroupCurrentTires := getGroupCurrentTires(calculatedActiveTireLoad, u)

		// Форматування як у твоєму updateResults()
		res.GroupUtilizationFactor = fmt.Sprintf("%.4f", calculatedGroupUtilizationFactor)
		res.EffectiveAmount = fmt.Sprintf("%.0f", math.Ceil(calculatedEffectiveAmount))
		res.ActivePowerFactor = fmt.Sprintf("%.2f", activePowerFactor) // у JS просто 1.25, але так читабельніше
		res.CalculatedActiveLoad = fmt.Sprintf("%.2f", calculatedActiveLoad)
		res.CalculatedReactiveLoad = fmt.Sprintf("%.2f", calculatedReactiveLoad)
		res.FullPower = fmt.Sprintf("%.3f", calculatedFullPower)
		res.GroupCurrent = fmt.Sprintf("%.2f", calculatedGroupCurrent)

		res.UtilizationRatesWorkshop = fmt.Sprintf("%.2f", calculatedUtilizationRatesWorkshop)
		res.EffectiveNumberEPWorkshop = fmt.Sprintf("%.2f", calculatedEffectiveNumberEPWorkshop)
		res.CoefficientActiveApacityWork = fmt.Sprintf("%.1f", coefficientActiveApacityWorkshop)

		res.ActiveTireLoad = fmt.Sprintf("%.1f", calculatedActiveTireLoad)
		res.ReactiveTireLoad = fmt.Sprintf("%.1f", calculatedReactiveLoadTires)
		res.FullPowerTires = fmt.Sprintf("%.1f", calculatedFullPowerTires)
		res.GroupCurrentTires = fmt.Sprintf("%.2f", calculatedGroupCurrentTires)

		_ = tpl.Execute(w, res)
		return

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// ----- читання інпутів -----

func readEPRow(r *http.Request, suffix string) EPRow {
	row := EPRow{
		Name: r.FormValue("name-of-EP" + suffix),
		Eta:  parseFloat(r.FormValue("nominal-value-efficiency-coefficient" + suffix)),
		Cos:  parseFloat(r.FormValue("load-power-factor" + suffix)),
		U:    parseFloat(r.FormValue("load-voltage" + suffix)),
		N:    parseFloat(r.FormValue("number-of-EP" + suffix)),
		Pn:   parseFloat(r.FormValue("nominal-power-of-EP" + suffix)),
		Kv:   parseFloat(r.FormValue("utilization-rate" + suffix)),
		Tg:   parseFloat(r.FormValue("reactive-power-factor" + suffix)),
	}
	return row
}

func parseFloat(s string) float64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	// щоб "0,92" теж парсилось
	s = strings.ReplaceAll(s, ",", ".")
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return v
}

// ----- ТІ Ж ФОРМУЛИ, що в JS -----

func multiplicationNandP(n, p float64) float64 { return n * p }

func sumMultiplicationNandPInSquare(n, p float64) float64 { return n * math.Pow(p, 2) }

func multiplicationK(nP, k float64) float64 { return nP * k }

// JS:
// const multiplicationTg = (tg, calculatedMultipK1) =>  {
//     const result = tg * calculatedMultipK1;
//     return result.toFixed(1);
// }
// потім Number(...) — тобто це реально число, округлене до 1 знака.
func multiplicationTgRounded1(tg, calculatedMultipK float64) float64 {
	result := tg * calculatedMultipK
	return roundTo(result, 1)
}

func getGroupUtilizationFactor(
	calculatedMultip1, calculatedMultip2, calculatedMultip3,
	calculatedMultipK1, calculatedMultipK2, calculatedMultipK3 float64,
) float64 {
	// JS:
	// sumMultiplicationNandP = calculatedMultip1 + calculatedMultip2 + calculatedMultip3 + 28 + 168 + 20 + 64 + 20;
	// sumMultiplicationK    = calculatedMultipK1 + calculatedMultipK2 + calculatedMultipK3 + 3.36 + 25.2 + 10 + 12.8 + 13;
	sumMultiplicationNandP := calculatedMultip1 + calculatedMultip2 + calculatedMultip3 + 28 + 168 + 20 + 64 + 20
	sumMultiplicationK := calculatedMultipK1 + calculatedMultipK2 + calculatedMultipK3 + 3.36 + 25.2 + 10 + 12.8 + 13
	if sumMultiplicationNandP == 0 {
		return 0
	}
	return sumMultiplicationK / sumMultiplicationNandP
}

func getEffectiveAmount(
	calculatedMultip1, calculatedMultip2, calculatedMultip3,
	calculatedMultipInSquare1, calculatedMultipInSquare2, calculatedMultipInSquare3 float64,
) float64 {
	// JS:
	// sumMultiplicationNandPInSquare = (calculatedMultip1 + calculatedMultip2 + calculatedMultip3 + 28 + 168 + 20 + 64 + 20) ** 2;
	// sumMultiplicationInSquare2 = calculatedMultipInSquare1 + calculatedMultipInSquare2 + calculatedMultipInSquare3
	//   + 14**2*2 + 42**2*4 + 20**2 + 32**2*2 + 20**2;
	sumBase := calculatedMultip1 + calculatedMultip2 + calculatedMultip3 + 28 + 168 + 20 + 64 + 20
	sumMultiplicationNandPInSquare := math.Pow(sumBase, 2)

	sumMultiplicationInSquare2 := calculatedMultipInSquare1 + calculatedMultipInSquare2 + calculatedMultipInSquare3 +
		math.Pow(14, 2)*2 + math.Pow(42, 2)*4 +
		math.Pow(20, 2) + math.Pow(32, 2)*2 + math.Pow(20, 2)

	if sumMultiplicationInSquare2 == 0 {
		return 0
	}
	return sumMultiplicationNandPInSquare / sumMultiplicationInSquare2
}

func getcalculatedActiveLoad(activePowerFactor, calculatedMultipK1, calculatedMultipK2, calculatedMultipK3 float64) float64 {
	// JS:
	// sumMultiplicationNandP = calculatedMultip1 + calculatedMultip2 + calculatedMultip3 + 3.36 + 25.2 + 10 + 12.8 + 13;
	// return activePowerFactor*sumMultiplicationNandP;
	sumMultiplicationNandP := calculatedMultipK1 + calculatedMultipK2 + calculatedMultipK3 + 3.36 + 25.2 + 10 + 12.8 + 13
	return activePowerFactor * sumMultiplicationNandP
}

func getcalculatedReactiveLoad(activePowerFactor, calculatedMultipTg1, calculatedMultipTg2, calculatedMultipTg3 float64) float64 {
	// JS:
	// sumMultiplicationTg = calculatedMultipTg1 + calculatedMultipTg2 + calculatedMultipTg3 + 3.36 + 33.5 + 12.8 + 7.5 + 9.5;
	// return activePowerFactor*sumMultiplicationTg;
	sumMultiplicationTg := calculatedMultipTg1 + calculatedMultipTg2 + calculatedMultipTg3 + 3.36 + 33.5 + 12.8 + 7.5 + 9.5
	return activePowerFactor * sumMultiplicationTg
}

func getFullPower(p, q float64) float64 { return math.Sqrt(math.Pow(p, 2) + math.Pow(q, 2)) }

func getGroupCurrent(p, u float64) float64 {
	if u == 0 {
		return 0
	}
	return p / u
}

func getUtilizationRatesWorkshop(
	calculatedMultip1, calculatedMultip2, calculatedMultip3,
	calculatedMultipK1, calculatedMultipK2, calculatedMultipK3 float64,
) float64 {
	// JS:
	// sumMultiplicationNandP = calculatedMultip1 + calculatedMultip2 + calculatedMultip3 + 28 + 168 + 20 + 64 + 20 + 456*2 + 465 + 200 +240;
	// sumMultiplicationK = calculatedMultipK1 + calculatedMultipK2 + calculatedMultipK3 + 3.36 + 25.2 + 10 + 12.8 + 13 + 95.1*3 + 40 + 192;
	sumMultiplicationNandP := calculatedMultip1 + calculatedMultip2 + calculatedMultip3 + 28 + 168 + 20 + 64 + 20 + 456*2 + 465 + 200 + 240
	sumMultiplicationK := calculatedMultipK1 + calculatedMultipK2 + calculatedMultipK3 + 3.36 + 25.2 + 10 + 12.8 + 13 + 95.1*3 + 40 + 192
	if sumMultiplicationNandP == 0 {
		return 0
	}
	return sumMultiplicationK / sumMultiplicationNandP
}

func getEffectiveNumberEPWorkshop(
	calculatedMultip1, calculatedMultip2, calculatedMultip3,
	calculatedMultipInSquare1, calculatedMultipInSquare2, calculatedMultipInSquare3 float64,
) float64 {
	// JS:
	// sumMultiplicationNandPInSquare = (calculatedMultip1 + calculatedMultip2 + calculatedMultip3 + 28 + 168 + 20 + 64 + 20 + 456*2 + 465 + 200 +240)**2;
	// sumMultiplicationInSquare2 = calculatedMultipInSquare1 + calculatedMultipInSquare2 + calculatedMultipInSquare3
	//   + 14**2*2 + 42**2*4 + 20**2 + 32**2*2 + 20**2 + 14792*3 + 20000 + 28800;
	sumBase := calculatedMultip1 + calculatedMultip2 + calculatedMultip3 + 28 + 168 + 20 + 64 + 20 + 456*2 + 465 + 200 + 240
	sumMultiplicationNandPInSquare := math.Pow(sumBase, 2)

	sumMultiplicationInSquare2 := calculatedMultipInSquare1 + calculatedMultipInSquare2 + calculatedMultipInSquare3 +
		math.Pow(14, 2)*2 + math.Pow(42, 2)*4 +
		math.Pow(20, 2) + math.Pow(32, 2)*2 + math.Pow(20, 2) +
		14792*3 + 20000 + 28800

	if sumMultiplicationInSquare2 == 0 {
		return 0
	}
	return sumMultiplicationNandPInSquare / sumMultiplicationInSquare2
}

func getActiveTireLoad(coefficientActiveApacityWorkshop float64) float64 {
	// JS: return 752*coefficientActiveApacityWorkshop;
	return 752 * coefficientActiveApacityWorkshop
}

func getReactiveLoadTires(coefficientActiveApacityWorkshop float64) float64 {
	// JS: return 657*coefficientActiveApacityWorkshop;
	return 657 * coefficientActiveApacityWorkshop
}

func getFullPowerTires(p, q float64) float64 { return math.Sqrt(math.Pow(p, 2) + math.Pow(q, 2)) }

func getGroupCurrentTires(p, u float64) float64 {
	if u == 0 {
		return 0
	}
	return p / u
}

// roundTo(x, decimals) — для імітації toFixed(decimals) як число (не рядок)
func roundTo(x float64, decimals int) float64 {
	pow := math.Pow(10, float64(decimals))
	return math.Round(x*pow) / pow
}