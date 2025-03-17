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

	entry := widget.NewEntry()
	entry.SetPlaceHolder("Enter expression (e.g. 3 + 4)")

	resultLabel := widget.NewLabel("Result: ")

	button := widget.NewButton("Calculate", func() {
		input := entry.Text
		result, err := evaluateExpression(input)
		if err != nil {
			resultLabel.SetText("Error: " + err.Error())
		} else {
			resultLabel.SetText("Result: " + strconv.FormatFloat(result, 'f', 2, 64))
		}
	})

	myWindow.SetContent(container.NewVBox(
		entry,
		button,
		resultLabel,
	))

	myWindow.ShowAndRun()
}

func evaluateExpression(expression string) (float64, error) {
	// Pisahkan ekspresi berdasarkan spasi
	tokens := strings.Fields(expression)
	if len(tokens) != 3 {
		return 0, fmt.Errorf("invalid expression format")
	}

	num1, err1 := strconv.ParseFloat(tokens[0], 64)
	if err1 != nil {
		return 0, fmt.Errorf("invalid number: %s", tokens[0])
	}

	num2, err2 := strconv.ParseFloat(tokens[2], 64)
	if err2 != nil {
		return 0, fmt.Errorf("invalid number: %s", tokens[2])
	}

	operator := tokens[1]
	switch operator {
	case "+":
		return num1 + num2, nil
	case "-":
		return num1 - num2, nil
	case "*":
		return num1 * num2, nil
	case "/":
		if num2 == 0 {
			return 0, fmt.Errorf("division by zero is not allowed")
		}
		return num1 / num2, nil
	default:
		return 0, fmt.Errorf("invalid operator: %s", operator)
	}
}
