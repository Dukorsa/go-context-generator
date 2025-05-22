package ui

import (
	"fmt"
	"image/color"
	"path/filepath"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// Definição de cores modernas
type ColorRGBA = color.NRGBA

var (
	// Cores primárias
	ColorPrimary   = color.NRGBA{R: 99, G: 102, B: 241, A: 255} // Indigo
	ColorSecondary = color.NRGBA{R: 16, G: 185, B: 129, A: 255} // Emerald
	ColorDanger    = color.NRGBA{R: 239, G: 68, B: 68, A: 255}  // Red
	ColorWarning   = color.NRGBA{R: 245, G: 158, B: 11, A: 255} // Amber
	ColorSuccess   = color.NRGBA{R: 34, G: 197, B: 94, A: 255}  // Green

	// Cores de fundo
	ColorBackground = color.NRGBA{R: 249, G: 250, B: 251, A: 255} // Gray-50
	ColorSurface    = color.NRGBA{R: 255, G: 255, B: 255, A: 255} // White
	ColorBorder     = color.NRGBA{R: 229, G: 231, B: 235, A: 255} // Gray-200

	// Cores de texto
	ColorTextPrimary   = color.NRGBA{R: 17, G: 24, B: 39, A: 255}    // Gray-900
	ColorTextSecondary = color.NRGBA{R: 107, G: 114, B: 128, A: 255} // Gray-500
	ColorTextMuted     = color.NRGBA{R: 156, G: 163, B: 175, A: 255} // Gray-400
)

func (a *App) layout(gtx layout.Context) layout.Dimensions {
	// Background
	paint.Fill(gtx.Ops, ColorBackground)

	return layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceBetween}.Layout(gtx,
			// Header
			layout.Rigid(a.layoutHeader),

			// Spacer
			layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),

			// Main Content
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				if a.showSettings {
					return a.layoutSettings(gtx)
				}
				return a.layoutMain(gtx)
			}),

			// Spacer
			layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),

			// Footer
			layout.Rigid(a.layoutFooter),
		)
	})
}

func (a *App) layoutHeader(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			title := material.H4(a.theme, "🚀 Go Context Generator Pro")
			title.Color = ColorTextPrimary
			return title.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btn := material.Button(a.theme, &a.settingsBtn, "⚙️")
			btn.Background = ColorSurface
			btn.Color = ColorTextSecondary
			return btn.Layout(gtx)
		}),
	)
}

func (a *App) layoutMain(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceEvenly}.Layout(gtx,
		// Cards de seleção
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Spacing: layout.SpaceEvenly}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return a.layoutPathCard(gtx, "📂 Pasta de Origem", a.srcPath, &a.selectSrcBtn, "Selecionar Origem")
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(16)}.Layout),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return a.layoutPathCard(gtx, "💾 Pasta de Destino", a.destPath, &a.selectDestBtn, "Selecionar Destino")
				}),
			)
		}),

		// Spacer
		layout.Rigid(layout.Spacer{Height: unit.Dp(24)}.Layout),

		// Status e Progress
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return a.layoutStatusSection(gtx)
		}),

		// Spacer
		layout.Rigid(layout.Spacer{Height: unit.Dp(24)}.Layout),

		// Generate Button
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return a.layoutGenerateButton(gtx)
			})
		}),
	)
}

func (a *App) layoutPathCard(gtx layout.Context, title, path string, btn *widget.Clickable, btnText string) layout.Dimensions {
	return Card{
		Color:        ColorSurface,
		CornerRadius: unit.Dp(12),
		Elevation:    unit.Dp(2),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceEvenly}.Layout(gtx,
				// Título
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					titleLabel := material.H6(a.theme, title)
					titleLabel.Color = ColorTextPrimary
					return titleLabel.Layout(gtx)
				}),

				// Spacer
				layout.Rigid(layout.Spacer{Height: unit.Dp(12)}.Layout),

				// Path ou placeholder
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if path != "" {
						pathLabel := material.Body2(a.theme, filepath.Base(path))
						pathLabel.Color = ColorTextSecondary
						return pathLabel.Layout(gtx)
					}

					emptyLabel := material.Body2(a.theme, "Nenhuma pasta selecionada")
					emptyLabel.Color = ColorTextMuted
					return emptyLabel.Layout(gtx)
				}),

				// Spacer
				layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),

				// Botão
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					button := material.Button(a.theme, btn, btnText)
					button.Background = ColorPrimary
					button.Color = ColorSurface
					return button.Layout(gtx)
				}),
			)
		})
	})
}

func (a *App) layoutStatusSection(gtx layout.Context) layout.Dimensions {
	return Card{
		Color:        ColorSurface,
		CornerRadius: unit.Dp(12),
		Elevation:    unit.Dp(1),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceEvenly}.Layout(gtx,
				// Status
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					status := material.Body1(a.theme, a.status)
					status.Color = a.getStatusColor()
					status.Alignment = text.Middle
					return status.Layout(gtx)
				}),

				// Progress bar se processando
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if a.isProcessing {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(layout.Spacer{Height: unit.Dp(12)}.Layout),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return ProgressBar{
									Progress: a.progress,
									Color:    ColorPrimary,
									Height:   unit.Dp(6),
									Radius:   unit.Dp(3),
								}.Layout(gtx)
							}),
							layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								percentage := int(a.progress * 100)
								label := material.Caption(a.theme, fmt.Sprintf("%d%% concluído", percentage))
								label.Color = ColorTextSecondary
								label.Alignment = text.Middle
								return label.Layout(gtx)
							}),
						)
					}
					return layout.Dimensions{}
				}),

				// Stats se houver arquivos processados
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if a.filesFound > 0 || a.filesGenerated > 0 {
						return layout.Flex{}.Layout(gtx,
							layout.Rigid(layout.Spacer{Height: unit.Dp(12)}.Layout),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return a.layoutStats(gtx)
							}),
						)
					}
					return layout.Dimensions{}
				}),
			)
		})
	})
}

func (a *App) layoutStats(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Spacing: layout.SpaceEvenly}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return a.layoutStatCard(gtx, "Encontrados", fmt.Sprintf("%d", a.filesFound), ColorPrimary)
		}),
		layout.Rigid(layout.Spacer{Width: unit.Dp(12)}.Layout),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return a.layoutStatCard(gtx, "Processados", fmt.Sprintf("%d", a.filesGenerated), ColorSuccess)
		}),
	)
}

func (a *App) layoutStatCard(gtx layout.Context, label, value string, color color.NRGBA) layout.Dimensions {
	return Card{
		Color:        ColorBackground,
		CornerRadius: unit.Dp(8),
		Elevation:    unit.Dp(0),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(12)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					valueLabel := material.H5(a.theme, value)
					valueLabel.Color = color
					valueLabel.Alignment = text.Middle
					return valueLabel.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					labelText := material.Caption(a.theme, label)
					labelText.Color = ColorTextSecondary
					labelText.Alignment = text.Middle
					return labelText.Layout(gtx)
				}),
			)
		})
	})
}

func (a *App) layoutGenerateButton(gtx layout.Context) layout.Dimensions {
	if a.canGenerate() && !a.isProcessing {
		btn := material.Button(a.theme, &a.generateBtn, "🚀 Gerar Arquivos de Contexto")
		btn.Background = ColorSuccess
		btn.Color = ColorSurface
		btn.Inset = layout.UniformInset(unit.Dp(16))
		return btn.Layout(gtx)
	}

	var btnText string
	var btnColor color.NRGBA

	if a.isProcessing {
		btnText = "⏳ Processando..."
		btnColor = ColorBorder
	} else {
		btnText = "Selecione as pastas primeiro"
		btnColor = ColorBorder
	}

	btn := material.Button(a.theme, &widget.Clickable{}, btnText)
	btn.Background = btnColor
	btn.Color = ColorTextMuted
	btn.Inset = layout.UniformInset(unit.Dp(16))
	return btn.Layout(gtx)
}

func (a *App) layoutSettings(gtx layout.Context) layout.Dimensions {
	return Card{
		Color:        ColorSurface,
		CornerRadius: unit.Dp(16),
		Elevation:    unit.Dp(4),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(24)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceEvenly}.Layout(gtx,
				// Título
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					title := material.H5(a.theme, "⚙️ Configurações")
					title.Color = ColorTextPrimary
					return title.Layout(gtx)
				}),

				layout.Rigid(layout.Spacer{Height: unit.Dp(24)}.Layout),

				// Opções
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.layoutCheckbox(gtx, &a.removeComments,
						"🧹 Remover comentários desnecessários",
						"Remove comentários que não são documentação importante")
				}),

				layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),

				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.layoutCheckbox(gtx, &a.includeTests,
						"🧪 Incluir arquivos de teste",
						"Processa arquivos *_test.go junto com o código")
				}),

				layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),

				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.layoutCheckbox(gtx, &a.minifyOutput,
						"⚡ Otimizar para IA",
						"Remove espaços extras e otimiza tokens para IA")
				}),

				layout.Rigid(layout.Spacer{Height: unit.Dp(24)}.Layout),

				// Botão voltar
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					btn := material.Button(a.theme, &a.settingsBtn, "← Voltar")
					btn.Background = ColorPrimary
					btn.Color = ColorSurface
					return btn.Layout(gtx)
				}),
			)
		})
	})
}

func (a *App) layoutCheckbox(gtx layout.Context, checkbox *widget.Bool, title, description string) layout.Dimensions {
	return layout.Flex{Alignment: layout.Start}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			cb := material.CheckBox(a.theme, checkbox, "")
			cb.Color = ColorPrimary
			return cb.Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Width: unit.Dp(12)}.Layout),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					titleLabel := material.Body1(a.theme, title)
					titleLabel.Color = ColorTextPrimary
					return titleLabel.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					desc := material.Caption(a.theme, description)
					desc.Color = ColorTextSecondary
					return desc.Layout(gtx)
				}),
			)
		}),
	)
}

func (a *App) layoutFooter(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			if a.lastRun != "" {
				label := material.Caption(a.theme, "Última execução: "+a.lastRun)
				label.Color = ColorTextMuted
				return label.Layout(gtx)
			}
			return layout.Dimensions{}
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			version := material.Caption(a.theme, "v2.0")
			version.Color = ColorTextMuted
			return version.Layout(gtx)
		}),
	)
}
