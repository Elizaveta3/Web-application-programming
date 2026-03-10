package main

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
)

type PageData struct {
	Current        string
	HighVoltage    string
	Time           string
	CalculatedLoad string
	Hours          string
	PowerKz        string
	RSn            string
	XSn            string
	RSnMin         string
	XSnMin         string

	ResultCurrentNormal        string
	ResultCurrentPostAccident  string
	ResultEconomicCrossSection string
	ResultThermalStability     string
	Xs                         string
	Xt                         string
	TotalResistance            string
	InitialCurrentValues       string
	ISh3                       string
	ISh2                       string
	IShMin3                    string
	IShMin2                    string
	IShn3                      string
	IShn2                      string
	IShnMin3                   string
	IShnMin2                   string
	Iln3                       string
	Iln2                       string
	IlnMin3                    string
	IlnMin2                    string
}

var tmpl = template.Must(template.New("index").Parse(`<!DOCTYPE html>
<html lang="uk">
<head>
    <meta charset="UTF-8">
    <title>Веб-калькулятор</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 1100px; margin: 20px auto; line-height: 1.5; }
        section { margin-bottom: 32px; padding: 16px; border: 1px solid #ddd; border-radius: 8px; }
        form { display: grid; gap: 10px; max-width: 420px; }
        input { padding: 8px; }
        button { padding: 10px 14px; cursor: pointer; }
        p { margin: 10px 0; }
    </style>
</head>
<body>
    <h1>Розрахунок струму трифазного КЗ, струму однофазного КЗ, та перевірки на термічну та динамічну стійкість.</h1>

    <section>
        <h2>Внесіть дані для вибору кабеля для живлення двотрансформаторної підстанції системи внутрішнього електропостачання підприємства:</h2>
        <form method="post" action="/task1">
            <label for="current">I<sub>к</sub>, кА:</label>
            <input type="number" id="current" name="current" placeholder="2,5" step="any" value="{{.Current}}">

            <label for="high-voltage">U<sub>ном</sub>, кВ:</label>
            <input type="number" id="high-voltage" name="high_voltage" placeholder="10" step="any" value="{{.HighVoltage}}">

            <label for="time">t<sub>ф</sub>, с:</label>
            <input type="number" id="time" name="time" placeholder="2,5" step="any" value="{{.Time}}">

            <label for="calculated-load">S<sub>M</sub>, кВ*А:</label>
            <input type="number" id="calculated-load" name="calculated_load" placeholder="1300" step="any" value="{{.CalculatedLoad}}">

            <label for="hours">Т<sub>M</sub>, год:</label>
            <input type="number" id="hours" name="hours" placeholder="4000" step="any" value="{{.Hours}}">

            <button type="submit">Submit</button>
        </form>
        <div>
            <p><span>Розрахунковий струм для нормального режима: </span><span>{{.ResultCurrentNormal}}</span><span> А</span></p>
            <p><span>Розрахунковий струм для післяаварійного режима: </span><span>{{.ResultCurrentPostAccident}}</span><span> А</span></p>
            <p><span>Для внутрізаводської мережі вибираємо броньовані кабелі з паперовою ізоляцією в алюмінієвій оболонці типу ААБ з економічним перерізом : </span><span>{{.ResultEconomicCrossSection}}</span><span> мм<sup>2</sup></span></p>
            <p><span>Вибираємо кабель ААБ 10 3×25 з допустимим струмом 90 А. Однак за термічною стійкістю до дії струмів КЗ:</span><span>{{.ResultThermalStability}}</span><span> мм<sup>2</sup></span></p>
        </div>
    </section>

    <section>
        <h2>Внесіть дані для визначення струму КЗ на шинах 10 кВ ГПП:</h2>
        <form method="post" action="/task2">
            <label for="power-kz">S<sub>к</sub>, МВ*А:</label>
            <input type="number" id="power-kz" name="power_kz" placeholder="200" step="any" value="{{.PowerKz}}">
            <button type="submit">Submit</button>
        </form>
        <div>
            <p><span>Опори елементів заступної схеми: Х<sub>с</sub>=</span><span>{{.Xs}}</span><span> Ом,</span><span>Х<sub>т</sub>=</span><span>{{.Xt}}</span><span> Ом</span></p>
            <p><span>Сумарний опір для точкі К1: </span><span>{{.TotalResistance}}</span><span> Ом</span></p>
            <p><span>Початкові значення струму трифазного КЗ: </span><span>{{.InitialCurrentValues}}</span><span> кА</span></p>
        </div>
    </section>

    <section>
        <h2>Внесіть дані для визначення струмів КЗ Хмельницьких північних електричних мереж (ХПнЕМ):</h2>
        <form method="post" action="/task3">
            <label for="r-sn">R<sub>с.н</sub>, Ом:</label>
            <input type="number" id="r-sn" name="r_sn" placeholder="10,65" step="any" value="{{.RSn}}">

            <label for="x-sn">X<sub>с.н</sub>, Ом:</label>
            <input type="number" id="x-sn" name="x_sn" placeholder="24,02" step="any" value="{{.XSn}}">

            <label for="r-s-min">R<sub>с.min</sub>, Ом:</label>
            <input type="number" id="r-s-min" name="r_s_min" placeholder="34,88" step="any" value="{{.RSnMin}}">

            <label for="x-s-min">X<sub>с.min</sub>, Ом:</label>
            <input type="number" id="x-s-min" name="x_s_min" placeholder="65,68" step="any" value="{{.XSnMin}}">

            <button type="submit">Submit</button>
        </form>
        <div>
            <p><span>Cтруми трифазного та двофазного КЗ на шинах 10 кВ в норм. та мін. режимах, приведені до напригу 110 кВ: I<sub>ш</sub><sup>(3)</sup>=</span><span>{{.ISh3}}</span><span> A, </span><span>I<sub>ш</sub><sup>(2)</sup>=</span><span>{{.ISh2}}</span><span> A, </span><span>I<sub>ш.min</sub><sup>(3)</sup>=</span><span>{{.IShMin3}}</span><span> A, </span><span>I<sub>ш.min</sub><sup>(2)</sup>=</span><span>{{.IShMin2}}</span><span> A</span></p>
            <p><span>Дійсні струми трифазного та двофазного КЗ на шинах 10 кВ в норм. та мін. режимах: I<sub>ш.н</sub><sup>(3)</sup>=</span><span>{{.IShn3}}</span><span> A, </span><span>I<sub>ш.н</sub><sup>(2)</sup>=</span><span>{{.IShn2}}</span><span> A, </span><span>I<sub>ш.н.min</sub><sup>(3)</sup>=</span><span>{{.IShnMin3}}</span><span> A, </span><span>I<sub>ш.н.min</sub><sup>(2)</sup>=</span><span>{{.IShnMin2}}</span><span> A</span></p>
            <p><span>Струми трифазного та двофазного КЗ в точці 10 в норм. та мін. режимах: I<sub>л.н</sub><sup>(3)</sup>=</span><span>{{.Iln3}}</span><span> A, </span><span>I<sub>л.н</sub><sup>(2)</sup>=</span><span>{{.Iln2}}</span><span> A, </span><span>I<sub>л.н.min</sub><sup>(3)</sup>=</span><span>{{.IlnMin3}}</span><span> A, </span><span>I<sub>л.н.min</sub><sup>(2)</sup>=</span><span>{{.IlnMin2}}</span><span> A</span></p>
        </div>
    </section>
</body>
</html>`))

func parseFloat(val string) float64 {
	f, _ := strconv.ParseFloat(val, 64)
	return f
}

func format1(v float64) string { return fmt.Sprintf("%.1f", v) }
func format2(v float64) string { return fmt.Sprintf("%.2f", v) }
func format3(v float64) string { return fmt.Sprintf("%.3f", v) }
func round(v float64) string   { return strconv.Itoa(int(math.Round(v))) }

func getCalculatedCurrentNormal(s, u float64) string {
	result := (s / 2) / (math.Sqrt(3) * u)
	return format1(result)
}

func getCalculatedCurrentPostAccident(i float64) string {
	return round(2 * i)
}

func getEconomicCrossSection(calculatedCurrentNormal, j float64) string {
	result := calculatedCurrentNormal / j
	return format1(result)
}

func determineJ(hours float64) float64 {
	if hours > 3000 && hours < 5000 {
		return 1.4
	} else if hours >= 1000 && hours < 3000 {
		return 1.6
	}
	return 1.2
}

func getThermalStability(current, time float64) string {
	result := (current * 1000 * math.Sqrt(time)) / 92
	return format1(result)
}

func getXs(u, s float64) string {
	result := math.Pow(u, 2) / s
	return format2(result)
}

func getXt(ucn, uk, s float64) string {
	result := (uk / 100) * (math.Pow(ucn, 2) / s)
	return format2(result)
}

func getTotalResistance(calculatedXs, calculatedXt float64) float64 {
	return calculatedXs + calculatedXt
}

func getInitialCurrentValues(ucn, totalResistance float64) string {
	result := ucn / (math.Sqrt(3) * totalResistance)
	return format1(result)
}

func getReactance(ukMax float64) string {
	result := (ukMax * math.Pow(115, 2)) / (100 * 6.3)
	return round(result)
}

func getXSh(XSn, calculatedReactance float64) float64 {
	return XSn + calculatedReactance
}

func getZSh(calculatedXSh float64) string {
	result := math.Sqrt(math.Pow(10.65, 2) + math.Pow(calculatedXSh, 2))
	return format1(result)
}

func getXShMin(XSnMin, calculatedReactance float64) float64 {
	return XSnMin + calculatedReactance
}

func getZShMin(RSnMin, calculatedXShMin float64) string {
	result := math.Sqrt(math.Pow(RSnMin, 2) + math.Pow(calculatedXShMin, 2))
	return format1(result)
}

func getISh3(calculatedZSh float64) string {
	result := (115 * 1000) / (1.73 * calculatedZSh)
	return round(result)
}

func getISh2(calculatedISh3 float64) string {
	result := calculatedISh3 * (1.73 / 2)
	return round(result)
}

func getIShMin3(calculatedZShMin float64) string {
	result := (115 * 1000) / (1.73 * calculatedZShMin)
	return round(result)
}

func getIShMin2(calculatedIShMin3 float64) string {
	result := calculatedIShMin3 * (1.73 / 2)
	return round(result)
}

func getCoef(unn, uvn float64) string {
	result := math.Pow(unn, 2) / math.Pow(uvn, 2)
	return format3(result)
}

func getRAndXShn(a, calculatedCoef float64) string {
	result := a * calculatedCoef
	return format2(result)
}

func getZShn(RSnMin, calculatedXShMin float64) string {
	result := math.Sqrt(math.Pow(RSnMin, 2) + math.Pow(calculatedXShMin, 2))
	return format2(result)
}

func getIShn3(calculatedZShn float64) string {
	result := (11 * 1000) / (1.73 * calculatedZShn)
	return round(result)
}

func getResistance(length, a float64) string {
	result := length * a
	return format2(result)
}

func render(w http.ResponseWriter, data PageData) {
	_ = tmpl.Execute(w, data)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	render(w, PageData{})
}

func task1Handler(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	currentStr := r.FormValue("current")
	highVoltageStr := r.FormValue("high_voltage")
	timeStr := r.FormValue("time")
	calculatedLoadStr := r.FormValue("calculated_load")
	hoursStr := r.FormValue("hours")

	current := parseFloat(currentStr)
	highVoltage := parseFloat(highVoltageStr)
	timeVal := parseFloat(timeStr)
	calculatedLoad := parseFloat(calculatedLoadStr)
	hours := parseFloat(hoursStr)

	j := determineJ(hours)
	calculatedCurrentNormal := getCalculatedCurrentNormal(calculatedLoad, highVoltage)
	calculatedCurrentPostAccident := getCalculatedCurrentPostAccident(parseFloat(calculatedCurrentNormal))
	calculatedEconomicCrossSection := getEconomicCrossSection(parseFloat(calculatedCurrentNormal), j)
	calculatedThermalStability := getThermalStability(current, timeVal)

	render(w, PageData{
		Current:                    currentStr,
		HighVoltage:                highVoltageStr,
		Time:                       timeStr,
		CalculatedLoad:             calculatedLoadStr,
		Hours:                      hoursStr,
		ResultCurrentNormal:        calculatedCurrentNormal,
		ResultCurrentPostAccident:  calculatedCurrentPostAccident,
		ResultEconomicCrossSection: calculatedEconomicCrossSection,
		ResultThermalStability:     calculatedThermalStability,
	})
}

func task2Handler(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	powerKzStr := r.FormValue("power_kz")
	powerKz := parseFloat(powerKzStr)

	averageNominalPointVoltage := 10.5
	shortCircuitVoltage := 10.5
	ratedPowerTransformer := 6.3

	calculatedXs := getXs(averageNominalPointVoltage, powerKz)
	calculatedXt := getXt(averageNominalPointVoltage, shortCircuitVoltage, ratedPowerTransformer)
	calculatedTotalResistance := getTotalResistance(parseFloat(calculatedXs), parseFloat(calculatedXt))
	calculatedInitialCurrentValues := getInitialCurrentValues(averageNominalPointVoltage, calculatedTotalResistance)

	render(w, PageData{
		PowerKz:              powerKzStr,
		Xs:                   calculatedXs,
		Xt:                   calculatedXt,
		TotalResistance:      strconv.FormatFloat(calculatedTotalResistance, 'f', -1, 64),
		InitialCurrentValues: calculatedInitialCurrentValues,
	})
}

func task3Handler(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	rSnStr := r.FormValue("r_sn")
	xSnStr := r.FormValue("x_sn")
	rSnMinStr := r.FormValue("r_s_min")
	xSnMinStr := r.FormValue("x_s_min")

	RSn := parseFloat(rSnStr)
	XSn := parseFloat(xSnStr)
	RSnMin := parseFloat(rSnMinStr)
	XSnMin := parseFloat(xSnMinStr)

	ukMax := 11.1
	unn := 11.0
	uvn := 115.0
	length := 12.37

	calculatedReactance := getReactance(ukMax)
	calculatedXSh := getXSh(XSn, parseFloat(calculatedReactance))
	calculatedZSh := getZSh(calculatedXSh)
	calculatedXShMin := getXShMin(XSnMin, parseFloat(calculatedReactance))
	calculatedZShMin := getZShMin(RSnMin, calculatedXShMin)
	calculatedISh3 := getISh3(parseFloat(calculatedZSh))
	calculatedISh2 := getISh2(parseFloat(calculatedISh3))
	calculatedIShMin3 := getIShMin3(parseFloat(calculatedZShMin))
	calculatedIShMin2 := getIShMin2(parseFloat(calculatedIShMin3))
	calculatedCoef := getCoef(unn, uvn)
	calculatedRShn := getRAndXShn(RSn, parseFloat(calculatedCoef))
	calculatedXShn := getRAndXShn(calculatedXSh, parseFloat(calculatedCoef))
	calculatedZShn := getZShn(parseFloat(calculatedRShn), parseFloat(calculatedXShn))
	calculatedRShnMin := getRAndXShn(RSnMin, parseFloat(calculatedCoef))
	calculatedXShnMin := getRAndXShn(calculatedXShMin, parseFloat(calculatedCoef))
	calculatedZShnMin := getZShMin(parseFloat(calculatedRShnMin), parseFloat(calculatedXShnMin))
	calculatedIShn3 := getIShn3(parseFloat(calculatedZShn))
	calculatedIShn2 := getIShMin2(parseFloat(calculatedIShn3))
	calculatedIShn3Min := getIShn3(parseFloat(calculatedZShnMin))
	calculatedIShn2Min := getIShMin2(parseFloat(calculatedIShn3Min))
	calculateResistance := getResistance(length, 0.64)
	calculateReactanceX := getResistance(length, 0.363)
	calculateRSumN := getXSh(parseFloat(calculateResistance), parseFloat(calculatedRShn))
	calculateXSumN := math.Round(getXSh(parseFloat(calculateReactanceX), parseFloat(calculatedXShn))*10) / 10
	calculatedZSumN := getZShn(calculateRSumN, calculateXSumN)
	calculateRSumNMin := getXSh(parseFloat(calculateResistance), parseFloat(calculatedRShnMin))
	calculateXSumNMin := getXSh(parseFloat(calculateReactanceX), parseFloat(calculatedXShnMin))
	calculatedZSumNMin := getZShn(calculateRSumNMin, calculateXSumNMin)
	calculatedIln3 := getIShn3(parseFloat(calculatedZSumN))
	calculatedIln2 := getIShMin2(parseFloat(calculatedIln3))
	calculatedIlnMin3 := getIShn3(parseFloat(calculatedZSumNMin))
	calculatedIlnMin2 := getIShMin2(parseFloat(calculatedIlnMin3))

	render(w, PageData{
		RSn:      rSnStr,
		XSn:      xSnStr,
		RSnMin:   rSnMinStr,
		XSnMin:   xSnMinStr,
		ISh3:     calculatedISh3,
		ISh2:     calculatedISh2,
		IShMin3:  calculatedIShMin3,
		IShMin2:  calculatedIShMin2,
		IShn3:    calculatedIShn3,
		IShn2:    calculatedIShn2,
		IShnMin3: calculatedIShn3Min,
		IShnMin2: calculatedIShn2Min,
		Iln3:     calculatedIln3,
		Iln2:     calculatedIln2,
		IlnMin3:  calculatedIlnMin3,
		IlnMin2:  calculatedIlnMin2,
	})
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/task1", task1Handler)
	http.HandleFunc("/task2", task2Handler)
	http.HandleFunc("/task3", task3Handler)

	fmt.Println("Server started at http://localhost:8080")
	_ = http.ListenAndServe(":8080", nil)
}
