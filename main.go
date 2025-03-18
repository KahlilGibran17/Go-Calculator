package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Knetic/govaluate"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myApp.Settings().SetTheme(theme.DarkTheme()) // Mode gelap
	myWindow := myApp.NewWindow("Calculator")

	display := widget.NewEntry()
	display.SetText("")
	history := widget.NewLabel("History:")

	buttons := []string{
		"7", "8", "9", "/",
		"4", "5", "6", "*",
		"1", "2", "3", "-",
		"C", "0", ".", "+",
		"=",
	}

	grid := container.NewGridWithColumns(4)

	lastInput := ""

	for _, btn := range buttons {
		button := widget.NewButton(btn, func(b string) func() {
			return func() {
				if b == "C" {
					display.SetText("")
					lastInput = ""
				} else if b == "=" {
					result, err := evaluateExpression(display.Text)
					if err != nil {
						display.SetText("Error")
					} else {
						history.SetText(history.Text + "\n" + display.Text + " = " + strconv.FormatFloat(result, 'f', 2, 64))
						display.SetText(strconv.FormatFloat(result, 'f', 2, 64))
					}
					lastInput = "="
				} else {
					if isOperator(b) && isOperator(lastInput) {
						return // Hindari menambahkan operator berturut-turut
					}
					display.SetText(display.Text + b)
					lastInput = b
				}
			}
		}(btn))
		grid.Add(button)
	}

	myWindow.Canvas().SetOnTypedKey(func(event *fyne.KeyEvent) {
		key := string(event.Name)
		if (key >= "0" && key <= "9") || key == "+" || key == "-" || key == "*" || key == "/" || key == "." {
			if isOperator(key) && isOperator(lastInput) {
				return // Hindari operator berulang
			}
			display.SetText(display.Text + string(key))
			lastInput = key
		} else if key == "Backspace" {
			if len(display.Text) > 0 {
				display.SetText(display.Text[:len(display.Text)-1])
			}
		} else if key == "Return" {
			result, err := evaluateExpression(display.Text)
			if err != nil {
				display.SetText("Error")
			} else {
				history.SetText(history.Text + "\n" + display.Text + " = " + strconv.FormatFloat(result, 'f', 2, 64))
				display.SetText(strconv.FormatFloat(result, 'f', 2, 64))
			}
		}
	})

	myWindow.SetContent(container.NewVBox(
		display,
		history,
		grid,
	))

	myWindow.ShowAndRun()
}

func evaluateExpression(expression string) (float64, error) {
	// Hindari kesalahan jika operator terakhir tidak valid
	expression = strings.TrimRight(expression, "+-*/")

	exp, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		return 0, fmt.Errorf("invalid expression")
	}

	result, err := exp.Evaluate(nil)
	if err != nil {
		return 0, fmt.Errorf("error evaluating")
	}

	return result.(float64), nil
}

func isOperator(s string) bool {
	return s == "+" || s == "-" || s == "*" || s == "/"
}
