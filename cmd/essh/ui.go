package main

import (
	"fmt"
	"os"

	"github.com/pterm/pterm"
)

func showtitle() {
	fmt.Println()
	title, err := pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle("E", pterm.NewStyle(pterm.FgLightCyan)),
		pterm.NewLettersFromStringWithStyle("SSH", pterm.NewStyle(pterm.FgLightRed))).Srender()
	check(err)
	subtext := pterm.DefaultBasicText.Sprintln("EC2 Instance Connect Helper")

	center := pterm.DefaultCenter.WithCenterEachLineSeparately(true)
	center.Print(title)
	center.Print(subtext)
}

func checkSpinnerError(spinner *pterm.SpinnerPrinter, err error) {
	if err != nil {
		spinner.Fail(err)
		os.Exit(1)
		// panic(err)
	}
}

func checkFatalError(err error) {
	if err != nil {
		pterm.Fatal.WithShowLineNumber(true).Print(err)
		os.Exit(1)
	}
}

// func clear() {
// 	print("\033[H\033[2J")
// }
