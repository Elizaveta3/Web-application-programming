package main

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
)

type PageData struct {
	AverageDailyCapacity string
	MeanSquareDeviation  string
	Oversight            string
	CostElectricity      string

	HasResult bool

	ShareEnergy         int
	W1                  float64
	Profit1             float64
	W2                  float64
	Fine1               float64
	ImprovedShareEnergy int
	W3                  float64
	Profit2             float64
	W4                  float64
	Fine2               float64
	MainProfit          float64
	Error               string
}

var pageTmpl = template.Must(template.New("page").Parse(`
<!DOCTYPE html>
<html lang="uk">
<head>
	<meta charset="UTF-8">
	<title>Веб-калькулятор</title>
	<style>
		body{
			font-family: Arial, sans-serif;
			padding: 40px;
			background: #f2f2f2;
		}
		h1{
			max-width: 1000px;
			line-height: 1.4;
		}
		.form-section, .result-section{
			background: white;
			padding: 20px;
			margin-top: 20px;
			border-radius: 10px;
			box-shadow: 0 2px 8px rgba(0,0,0,0.08);
			max-width: 800px;
		}
		form{
			display: grid;
			grid-template-columns: 1fr;
			gap: 10px;
		}
		input{
			padding: 8px;
			font-size: 16px;
		}
		button{
			padding: 10px 18px;
			font-size: 16px;
			cursor: pointer;
			margin-top: 10px;
		}
		p{
			font-size: 16px;
		}
		.error{
			color: #b00020;
			font-weight: bold;
		}
	</style>
</head>
<body>
	<h1>ПРАКТИЧНА РОБОТА №6. Розрахунок прибутку від сонячних електростанцій з встановленою системою прогнозування сонячної потужності.</h1>

	<section class="form-section">
		<h2>Внесіть дані до розрахунку прибутку від сонячних електростанцій:</h2>

		{{if .Error}}
			<p class="error">{{.Error}}</p>
		{{end}}

		<form method="POST" action="/">
			<label for="average-daily-capacity">P<sub>c</sub>, МВт:</label>
			<input type="number" step="any" id="average-daily-capacity" name="average-daily-capacity" placeholder="5" value="{{.AverageDailyCapacity}}">

			<label for="mean-square-deviation">Сер. квад. відхилення, МВт:</label>
			<input type="number" step="any" id="mean-square-deviation" name="mean-square-deviation" placeholder="1" value="{{.MeanSquareDeviation}}">

			<label for="oversight">Зменшити похибку до:</label>
			<input type="number" step="any" id="oversight" name="oversight" placeholder="0.25" value="{{.Oversight}}">

			<label for="cost-electricity">В, грн/кВт*год:</label>
			<input type="number" step="any" id="cost-electricity" name="cost-electricity" placeholder="7" value="{{.CostElectricity}}">

			<button type="submit">Submit</button>
		</form>
	</section>

	{{if .HasResult}}
	<section class="result-section">
		<div>
			<p><span>Частка енергії, що генерується без небалансів: </span><span>{{.ShareEnergy}}</span><span> %</span></p>
			<p><span>W1: </span><span>{{printf "%.2f" .W1}}</span><span> МВт*год</span></p>
			<p><span>Прибуток 1: </span><span>{{printf "%.2f" .Profit1}}</span><span> тис. грн</span></p>
			<p><span>W2: </span><span>{{printf "%.2f" .W2}}</span><span> МВт*год</span></p>
			<p><span>Штраф 1: </span><span>{{printf "%.2f" .Fine1}}</span><span> тис. грн</span></p>
			<p><span>Після вдосконалення системи прогнозу частка енергії, що генерується без небалансів: </span><span>{{.ImprovedShareEnergy}}</span><span> %</span></p>
			<p><span>W3: </span><span>{{printf "%.2f" .W3}}</span><span> МВт*год</span></p>
			<p><span>Прибуток 2: </span><span>{{printf "%.1f" .Profit2}}</span><span> тис. грн</span></p>
			<p><span>W4: </span><span>{{printf "%.2f" .W4}}</span><span> МВт*год</span></p>
			<p><span>Штраф 2: </span><span>{{printf "%.2f" .Fine2}}</span><span> тис. грн</span></p>
			<p><span>Головний прибуток: </span><span>{{printf "%.1f" .MainProfit}}</span><span> тис. грн</span></p>
		</div>
	</section>
	{{end}}
</body>
</html>
`))

func main() {
	http.HandleFunc("/", homePage)

	fmt.Println("Server started at http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	data := PageData{}

	if r.Method == http.MethodPost {
		data.AverageDailyCapacity = r.FormValue("average-daily-capacity")
		data.MeanSquareDeviation = r.FormValue("mean-square-deviation")
		data.Oversight = r.FormValue("oversight")
		data.CostElectricity = r.FormValue("cost-electricity")

		averageCapacity, err1 := strconv.ParseFloat(data.AverageDailyCapacity, 64)
		meanSquareDev, err2 := strconv.ParseFloat(data.MeanSquareDeviation, 64)
		oversight, err3 := strconv.ParseFloat(data.Oversight, 64)
		costElectricity, err4 := strconv.ParseFloat(data.CostElectricity, 64)

		if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
			data.Error = "Будь ласка, введіть коректні числові значення."
		} else {
			calcShareEnergy := getShareEnergy(averageCapacity, meanSquareDev)
			calcW1 := getW1(averageCapacity, calcShareEnergy)
			calcProfit1 := getProfit(calcW1, costElectricity)
			calcW2 := getW2(averageCapacity, calcShareEnergy)
			calcFine1 := getFine(calcW2, costElectricity)

			calcImprovedShareEnergy := getShareEnergy(averageCapacity, oversight)
			calcW3 := getW1(averageCapacity, calcImprovedShareEnergy)
			calcProfit2 := roundTo1(getProfit(calcW3, costElectricity))
			calcW4 := getW2(averageCapacity, calcImprovedShareEnergy)
			calcFine2 := getFine(calcW4, costElectricity)
			calcMainProfit := roundTo1(getMainProfit(calcProfit2, calcFine2))

			data.HasResult = true
			data.ShareEnergy = calcShareEnergy
			data.W1 = calcW1
			data.Profit1 = calcProfit1
			data.W2 = calcW2
			data.Fine1 = calcFine1
			data.ImprovedShareEnergy = calcImprovedShareEnergy
			data.W3 = calcW3
			data.Profit2 = calcProfit2
			data.W4 = calcW4
			data.Fine2 = calcFine2
			data.MainProfit = calcMainProfit
		}
	}

	err := pageTmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
	}
}

func getShareEnergy(averageCapacity, meanSquareDev float64) int {
	lowerBound := 4.75
	upperBound := 5.25

	a := erf((upperBound - averageCapacity) / (math.Sqrt(2) * meanSquareDev))
	b := erf((lowerBound - averageCapacity) / (math.Sqrt(2) * meanSquareDev))

	result := 0.5 * (a - b)

	return int(math.Round(result))
}

func erf(x float64) float64 {
	p := 0.3275911
	a1 := 0.254829592
	a2 := -0.284496736
	a3 := 1.421413741
	a4 := -1.453152027
	a5 := 1.061405429

	sign := 1.0
	if x < 0 {
		sign = -1.0
	}
	x = math.Abs(x)

	t := 1.0 / (1.0 + p*x)
	y := 1.0 - (((((a5*t+a4)*t)+a3)*t+a2)*t+a1)*t*math.Exp(-x*x)

	return sign * y * 100
}

func getW1(averageCapacity float64, calcShareEnergy int) float64 {
	return (averageCapacity * 24 * float64(calcShareEnergy)) / 100
}

func getProfit(w, costElectricity float64) float64 {
	return w * costElectricity
}

func getW2(averageCapacity float64, calcShareEnergy int) float64 {
	return (averageCapacity * 24 * float64(100-calcShareEnergy)) / 100
}

func getFine(w, costElectricity float64) float64 {
	return w * costElectricity
}

func getMainProfit(profit, fine float64) float64 {
	return profit - fine
}

func roundTo1(x float64) float64 {
	return math.Round(x*10) / 10
}
