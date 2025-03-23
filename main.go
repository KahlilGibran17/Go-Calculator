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

	// Simpan preferensi tema terakhir
	pref := myApp.Preferences()
	isDarkMode := pref.BoolWithFallback("dark_mode", true)
	if isDarkMode {
		myApp.Settings().SetTheme(theme.DarkTheme())
	} else {
		myApp.Settings().SetTheme(theme.LightTheme())
	}

	myWindow := myApp.NewWindow("Stylish Calculator")

	display := widget.NewEntry()
	display.SetText("")
	display.Disable()

	history := widget.NewLabel("History:")

	// **🔹 Deklarasikan variabel tombol dulu**
	var toggleThemeButton *widget.Button

	// Tombol kalkulator
	buttons := []struct {
		label string
		icon  fyne.Resource
	}{
		{"7", nil}, {"8", nil}, {"9", nil}, {"/", theme.MediaFastForwardIcon()},
		{"4", nil}, {"5", nil}, {"6", nil}, {"*", theme.MediaPlayIcon()},
		{"1", nil}, {"2", nil}, {"3", nil}, {"-", theme.NavigateBackIcon()},
		{"C", nil}, {"0", nil}, {".", nil}, {"+", theme.NavigateNextIcon()},
		{"⌫", theme.NavigateBackIcon()}, {"=", theme.ConfirmIcon()},
	}

	grid := container.NewGridWithColumns(4)
	lastInput := ""

	for _, btn := range buttons {
		button := widget.NewButton(btn.label, func(b string) func() {
			return func() {
				text := display.Text
				if b == "C" {
					display.SetText("")
					lastInput = ""
				} else if b == "⌫" { // Backspace Function
					if len(text) > 0 {
						display.SetText(text[:len(text)-1])
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
						return
					}
					display.SetText(text + b)
					lastInput = b
				}
			}
		}(btn.label))

		if btn.icon != nil {
			button.SetIcon(btn.icon)
		}

		grid.Add(button)
	}

	// **🔹 Isi variabel toggleThemeButton setelah dideklarasikan**
	toggleThemeButton = widget.NewButtonWithIcon("Dark Mode", theme.ViewRefreshIcon(), func() {
		isDarkMode = !isDarkMode
		pref.SetBool("dark_mode", isDarkMode)

		if isDarkMode {
			myApp.Settings().SetTheme(theme.DarkTheme())
			toggleThemeButton.SetText("Dark Mode")
		} else {
			myApp.Settings().SetTheme(theme.LightTheme())
			toggleThemeButton.SetText("Light Mode")
		}
	})

	// **Gunakan tombol di dalam container setelah dideklarasikan**
	mainContainer := container.NewVBox(
		toggleThemeButton,
		display,
		history,
		grid,
	)

	myWindow.SetContent(container.NewPadded(mainContainer))
	myWindow.ShowAndRun()
}

// Evaluasi ekspresi matematika
func evaluateExpression(expression string) (float64, error) {
	expression = strings.TrimSpace(expression)

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

	if res, ok := result.(float64); ok {
		return res, nil
	}
	return 0, fmt.Errorf("unexpected result type")
}

// Fungsi untuk mengecek apakah karakter adalah operator
func isOperator(s string) bool {
	return s == "+" || s == "-" || s == "*" || s == "/"
}
