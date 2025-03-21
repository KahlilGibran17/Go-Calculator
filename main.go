package main

import (
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Knetic/govaluate"
)

func main() {
	myApp := app.New()
	myApp.Settings().SetTheme(theme.DarkTheme()) // Bisa diganti dengan theme.LightTheme()
	myWindow := myApp.NewWindow("Stylish Calculator")

	display := widget.NewEntry()
	display.SetText("")
	display.Disable() // Agar tidak bisa diketik langsung

	history := widget.NewLabel("History:")

	// Tombol kalkulator dengan ikon
	buttons := []struct {
		label string
		icon  fyne.Resource
	}{
		{"7", nil}, {"8", nil}, {"9", nil}, {"/", theme.MediaFastForwardIcon()},
		{"4", nil}, {"5", nil}, {"6", nil}, {"*", theme.MediaPlayIcon()},
		{"1", nil}, {"2", nil}, {"3", nil}, {"-", theme.NavigateBackIcon()},
		{"C", nil}, {"0", nil}, {".", nil}, {"+", theme.NavigateNextIcon()},
		{"⌫", theme.CancelIcon()},
		{"=", theme.ConfirmIcon()},
	}

	grid := container.NewGridWithColumns(4)

	lastInput := ""

	// Membuat tombol dengan tampilan lebih modern
	for _, btn := range buttons {
		button := widget.NewButton(btn.label, func(b string) func() {
			return func() {
				text := display.Text // Ambil teks yang ada di layar kalkulator

				if b == "C" {
					display.SetText("")
					lastInput = ""
				} else if b == "⌫" { // Tambahkan kondisi untuk tombol Backspace
					if len(text) > 0 {
						display.SetText(text[:len(text)-1]) // Hapus karakter terakhir
					}
				} else if b == "=" {
					result, err := evaluateExpression(text)
					if err != nil {
						display.SetText("Error")
					} else {
						history.SetText(history.Text + "\n" + text + " = " + strconv.FormatFloat(result, 'f', 2, 64))
						display.SetText(strconv.FormatFloat(result, 'f', 2, 64))
					}
					lastInput = "="
				} else {
					if isOperator(b) && (text == "" || isOperator(lastInput)) {
						return // Hindari operator berturut-turut atau operator di awal
					}

					display.SetText(text + b) // Update display dengan input baru
					lastInput = b
				}
			}
		}(btn.label))

		// Jika tombol memiliki ikon, gunakan ikon
		if btn.icon != nil {
			button.SetIcon(btn.icon)
		}

		grid.Add(button)
	}

	// Menambahkan padding agar lebih rapi
	mainContainer := container.NewPadded(
		container.NewVBox(
			display,
			history,
			grid,
		),
	)

	myWindow.SetContent(mainContainer)
	myWindow.ShowAndRun()
}

// Evaluasi ekspresi matematika
func evaluateExpression(expression string) (float64, error) {
	expression = strings.TrimSpace(expression) // Hapus spasi yang tidak perlu

	// Cek jika terakhir adalah operator
	if strings.HasSuffix(expression, "+") ||
		strings.HasSuffix(expression, "-") ||
		strings.HasSuffix(expression, "*") ||
		strings.HasSuffix(expression, "/") {
		return 0, fmt.Errorf("invalid input")
	}

	exp, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		return 0, fmt.Errorf("invalid expression")
	}

	result, err := exp.Evaluate(nil)
	if err != nil {
		return 0, fmt.Errorf("error evaluating")
	}

	// Pastikan hasil berupa float64
	if res, ok := result.(float64); ok {
		return res, nil
	}
	return 0, fmt.Errorf("unexpected result type")
}

// Fungsi untuk mengecek apakah karakter adalah operator
func isOperator(s string) bool {
	return s == "+" || s == "-" || s == "*" || s == "/"
}
