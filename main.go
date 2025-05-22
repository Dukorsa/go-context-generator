package main

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"

	"go-context-generator/internal/ui"
)

func main() {
	go func() {
		w := app.NewWindow(
			app.Title("Go Context Generator Pro"),
			app.Size(unit.Dp(1000), unit.Dp(700)),
			app.MinSize(unit.Dp(800), unit.Dp(600)),
		)

		th := material.NewTheme()
		th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
		app := ui.NewApp(th)

		if err := app.Run(w); err != nil {
			log.Printf("Erro na aplicação: %v", err)
			os.Exit(1)
		}
	}()

	app.Main()
}
