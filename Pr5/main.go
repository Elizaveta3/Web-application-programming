package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type PageData struct {
	Task1Result string
	Task2Result string
}

func main() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/task1", calculateTask1)
	http.HandleFunc("/task2", calculateTask2)

	http.Handle("/styles.css", http.FileServer(http.Dir(".")))

	fmt.Println("Server started at http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}

func renderTemplate(w http.ResponseWriter, data PageData) {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, "Не вдалося завантажити index.html", http.StatusInternalServerError)
		fmt.Println("Template error:", err)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Помилка відображення сторінки", http.StatusInternalServerError)
		fmt.Println("Execute error:", err)
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, PageData{})
}

func calculateTask1(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	sum := 0.0

	if value := r.FormValue("v110"); value != "" {
		v, _ := strconv.ParseFloat(value, 64)
		sum += v
	}

	if value := r.FormValue("t110"); value != "" {
		v, _ := strconv.ParseFloat(value, 64)
		sum += v
	}

	if value := r.FormValue("bus10"); value != "" {
		v, _ := strconv.ParseFloat(value, 64)
		qty, _ := strconv.ParseFloat(r.FormValue("bus10_quantity"), 64)
		if qty <= 0 {
			qty = 1
		}
		sum += v * qty
	}

	if value := r.FormValue("pl110"); value != "" {
		v, _ := strconv.ParseFloat(value, 64)
		sum += v
	}

	if value := r.FormValue("pl10"); value != "" {
		v, _ := strconv.ParseFloat(value, 64)
		sum += v
	}

	result := fmt.Sprintf("Показник надійності одноколової системи: %.4f", sum)
	renderTemplate(w, PageData{Task1Result: result})
}

func calculateTask2(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	lossEmergency, _ := strconv.ParseFloat(r.FormValue("loss_emergency"), 64)
	lossPlanned, _ := strconv.ParseFloat(r.FormValue("loss_planned"), 64)

	totalLoss := lossEmergency + lossPlanned

	result := fmt.Sprintf("Загальні збитки від перерв електропостачання: %.2f грн/кВт*год", totalLoss)
	renderTemplate(w, PageData{Task2Result: result})
}