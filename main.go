package main

import (
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Calculator")

	display := widget.NewEntry()
	display.Disable()
	display.SetText("")

	buttons := []string{
		"7", "8", "9", "/",
		"4", "5", "6", "*",
		"1", "2", "3", "-",
		"C", "0", "=", "+",
	}

	grid := container.NewGridWithColumns(4)

	for _, btn := range buttons {
		button := widget.NewButton(btn, func(b string) func() {
			return func() {
				text := display.Text

				switch b {
				case "C":
					display.SetText("")
				case "=":
					result, err := evaluateExpression(text)
					if err != nil {
						display.SetText("Error")
					} else {
						display.SetText(strconv.FormatFloat(result, 'f', 2, 64))
					}
				default:
					display.SetText(text + b)
				}
			}
		}(btn))
		grid.Add(button)
	}

	myWindow.SetContent(container.NewVBox(
		display,
		grid,
	))

	myWindow.ShowAndRun()
}

// ✅ Fungsi untuk mengevaluasi ekspresi matematika
func evaluateExpression(expression string) (float64, error) {
	expression = strings.ReplaceAll(expression, " ", "") // Hapus spasi

	// Mencari operator yang digunakan
	var operator string
	if strings.Contains(expression, "+") {
		operator = "+"
	} else if strings.Contains(expression, "-") {
		operator = "-"
	} else if strings.Contains(expression, "*") {
		operator = "*"
	} else if strings.Contains(expression, "/") {
		operator = "/"
	} else {
		return 0, fmt.Errorf("invalid expression")
	}

	// Memisahkan angka berdasarkan operator
	parts := strings.Split(expression, operator)
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid expression")
	}

	num1, err1 := strconv.ParseFloat(parts[0], 64)
	num2, err2 := strconv.ParseFloat(parts[1], 64)

	if err1 != nil || err2 != nil {
		return 0, fmt.Errorf("invalid number")
	}

	// Operasi matematika
	switch operator {
	case "+":
		return num1 + num2, nil
	case "-":
		return num1 - num2, nil
	case "*":
		return num1 * num2, nil
	case "/":
		if num2 == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return num1 / num2, nil
	}

	return 0, fmt.Errorf("invalid operation")
}
