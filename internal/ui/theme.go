package ui

import (
	"fmt"
	"image/color"
	"path/filepath"
	"strings" // Adicionado para a lógica de status com ícones

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type ColorRGBA = color.NRGBA

// Paleta de Cores Refinada
var (
	ColorPrimary        = color.NRGBA{R: 99, G: 102, B: 241, A: 255} // Indigo-500
	ColorPrimaryHover   = color.NRGBA{R: 80, G: 83, B: 220, A: 255}  // Um pouco mais escuro/saturado para hover
	ColorPrimaryPressed = color.NRGBA{R: 65, G: 68, B: 200, A: 255}  // Mais escuro para pressed

	ColorSuccess        = color.NRGBA{R: 16, G: 185, B: 129, A: 255} // Emerald-500 (Era Secondary)
	ColorSuccessHover   = color.NRGBA{R: 14, G: 165, B: 115, A: 255}
	ColorSuccessPressed = color.NRGBA{R: 12, G: 145, B: 100, A: 255}

	ColorDanger  = color.NRGBA{R: 239, G: 68, B: 68, A: 255}  // Red-500
	ColorWarning = color.NRGBA{R: 245, G: 158, B: 11, A: 255} // Amber-500

	ColorBackground  = color.NRGBA{R: 243, G: 244, B: 246, A: 255} // Gray-100 (um pouco mais escuro que Gray-50)
	ColorSurface     = color.NRGBA{R: 255, G: 255, B: 255, A: 255} // White
	ColorBorder      = color.NRGBA{R: 229, G: 231, B: 235, A: 255} // Gray-200
	ColorBorderLight = color.NRGBA{R: 240, G: 240, B: 240, A: 255} // Para bordas muito sutis

	ColorTextPrimary   = color.NRGBA{R: 17, G: 24, B: 39, A: 255}    // Gray-900
	ColorTextSecondary = color.NRGBA{R: 75, G: 85, B: 99, A: 255}    // Gray-600 (um pouco mais escuro que Gray-500)
	ColorTextMuted     = color.NRGBA{R: 156, G: 163, B: 175, A: 255} // Gray-400
	ColorTextOnPrimary = color.NRGBA{R: 255, G: 255, B: 255, A: 255} // White (para texto em botões primários)
	ColorTextOnSuccess = color.NRGBA{R: 255, G: 255, B: 255, A: 255} // White (para texto em botões de sucesso)

	ColorTransparent = color.NRGBA{A: 0}
)

// Constantes de UI (conceituais)
var (
	smallRadius  = unit.Dp(8)
	mediumRadius = unit.Dp(12)
	largeRadius  = unit.Dp(16)

	smallPadding  = unit.Dp(8)
	mediumPadding = unit.Dp(16)
	largePadding  = unit.Dp(24)
	xlargePadding = unit.Dp(32)
)

func (a *App) layout(gtx layout.Context) layout.Dimensions {
	paint.Fill(gtx.Ops, ColorBackground)

	return layout.UniformInset(largePadding).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(a.layoutHeader),
			layout.Rigid(layout.Spacer{Height: largePadding}.Layout),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				if a.showSettings {
					return a.layoutSettings(gtx)
				}
				return a.layoutMain(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: largePadding}.Layout),
			layout.Rigid(a.layoutFooter),
		)
	})
}

func (a *App) layoutHeader(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Alignment: layout.Middle, Spacing: layout.SpaceBetween}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			title := material.H4(a.theme, "🚀 Go Context Generator Pro")
			title.Color = ColorTextPrimary
			title.Font.Weight = font.Bold // Dar mais peso ao título
			return title.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btn := material.Button(a.theme, &a.settingsBtn, "⚙️ Configurações")
			btn.Background = ColorTransparent // Botão de texto/ícone
			btn.Color = ColorTextSecondary
			btn.Inset = layout.UniformInset(smallPadding)
			btn.CornerRadius = smallRadius

			// Feedback de hover sutil para botões de texto/ícone
			if a.settingsBtn.Hovered() {
				btn.Color = ColorPrimary
			}
			if a.showSettings { // Indicar estado ativo
				btn.Background = ColorBorderLight // Fundo sutil quando ativo
			}
			return btn.Layout(gtx)
		}),
	)
}

func (a *App) layoutMain(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx, // SpaceBetween para ocupar largura
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return a.layoutPathCard(gtx, "📂 Pasta de Origem", a.srcPath, &a.selectSrcBtn, "Selecionar Origem")
				}),
				layout.Rigid(layout.Spacer{Width: mediumPadding}.Layout),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return a.layoutPathCard(gtx, "💾 Pasta de Destino", a.destPath, &a.selectDestBtn, "Selecionar Destino")
				}),
			)
		}),
		layout.Rigid(layout.Spacer{Height: largePadding}.Layout),
		layout.Rigid(a.layoutStatusSection),
		// layout.Flexed(1, layout.Spacer{}.Layout), // Empurra o botão de gerar para baixo se necessário
		layout.Rigid(layout.Spacer{Height: largePadding}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, a.layoutGenerateButton) // Centraliza o botão de gerar
		}),
	)
}

func (a *App) layoutPathCard(gtx layout.Context, title, path string, btnWidget *widget.Clickable, btnText string) layout.Dimensions {
	return Card{
		Color:        ColorSurface,
		CornerRadius: mediumRadius,
		Elevation:    unit.Dp(2), // Sutil elevação
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(mediumPadding).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical, Alignment: layout.Start, Spacing: layout.SpaceBetween}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					cardTitle := material.H6(a.theme, title)
					cardTitle.Color = ColorTextPrimary
					return cardTitle.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: smallPadding}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					// Exibição do caminho com truncamento elegante
					pathDisplay := "Nenhuma pasta selecionada."
					pathColor := ColorTextMuted
					if path != "" {
						dir := filepath.Dir(path)
						base := filepath.Base(path)
						// Tentar mostrar "pasta_pai/pasta_selecionada"
						// ou "...pasta_longa_demais/pasta_selecionada"
						displayDir := filepath.Base(dir)
						if displayDir == "." || displayDir == "" { // Se for raiz ou algo similar
							pathDisplay = base
						} else {
							// Limitar o comprimento do diretório pai para exibição
							maxDirLen := 25
							if len(dir) > maxDirLen {
								displayDir = "..." + dir[len(dir)-(maxDirLen-3):]
							} else {
								displayDir = dir
							}
							pathDisplay = filepath.Join(displayDir, base)
						}
						pathColor = ColorTextSecondary
					}
					pathLabel := material.Body2(a.theme, pathDisplay)
					pathLabel.Color = pathColor
					return pathLabel.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: mediumPadding}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					btn := material.Button(a.theme, btnWidget, btnText)
					btn.Background = ColorPrimary
					if btnWidget.Pressed() {
						btn.Background = ColorPrimaryPressed
					} else if btnWidget.Hovered() {
						btn.Background = ColorPrimaryHover
					}
					btn.Color = ColorTextOnPrimary
					btn.CornerRadius = smallRadius
					btn.Inset = layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(10), Left: mediumPadding, Right: mediumPadding}
					return btn.Layout(gtx)
				}),
			)
		})
	})
}

func (a *App) layoutStatusSection(gtx layout.Context) layout.Dimensions {
	return Card{
		Color:        ColorSurface,
		CornerRadius: mediumRadius,
		Elevation:    unit.Dp(1), // Menor elevação para status
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(mediumPadding).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			a.mu.RLock()
			status := a.status
			isProcessing := a.isProcessing
			progress := a.progress
			filesFound := a.filesFound
			filesGenerated := a.filesGenerated
			a.mu.RUnlock()

			statusColor := a.getStatusColor() // Sua lógica para cor de status é boa

			// Adicionar ícones ao status de forma mais sistemática
			var statusIcon string
			if strings.Contains(status, "✅") {
				statusIcon = "✅ "
				status = strings.Replace(status, "✅", "", 1)
			}
			if strings.Contains(status, "❌") {
				statusIcon = "❌ "
				status = strings.Replace(status, "❌", "", 1)
			}
			if strings.Contains(status, "⚠️") {
				statusIcon = "⚠️ "
				status = strings.Replace(status, "⚠️", "", 1)
			}
			if strings.Contains(status, "🔍") {
				statusIcon = "🔍 "
				status = strings.Replace(status, "🔍", "", 1)
			}
			if strings.Contains(status, "🔄") {
				statusIcon = "🔄 "
				status = strings.Replace(status, "🔄", "", 1)
			}
			if strings.Contains(status, "⚡") {
				statusIcon = "⚡ "
				status = strings.Replace(status, "⚡", "", 1)
			}
			if isProcessing && statusIcon == "" {
				statusIcon = "⏳ "
			}
			status = strings.TrimSpace(status)

			return layout.Flex{Axis: layout.Vertical, Alignment: layout.Start}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if statusIcon != "" {
								iconLabel := material.Label(a.theme, unit.Sp(16), statusIcon) // Tamanho do ícone um pouco maior
								iconLabel.Color = statusColor
								return layout.Inset{Right: smallPadding}.Layout(gtx, iconLabel.Layout)
							}
							return layout.Dimensions{}
						}),
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							statusLabel := material.Body1(a.theme, status)
							statusLabel.Color = statusColor
							return statusLabel.Layout(gtx)
						}),
					)
				}),

				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if isProcessing {
						return layout.Inset{Top: mediumPadding, Bottom: smallPadding}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{Axis: layout.Vertical, Alignment: layout.Start}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return ProgressBar{
										Progress: progress,
										Color:    ColorPrimary,
										Height:   unit.Dp(10), // Barra de progresso mais proeminente
										Radius:   unit.Dp(5),
									}.Layout(gtx)
								}),
								layout.Rigid(layout.Spacer{Height: smallPadding}.Layout),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									percentage := int(progress * 100)
									progressText := fmt.Sprintf("%d%% concluído", percentage)
									if filesFound > 0 && filesGenerated < filesFound {
										progressText = fmt.Sprintf("%d / %d arquivos (%d%%)", filesGenerated, filesFound, percentage)
									}
									label := material.Caption(a.theme, progressText)
									label.Color = ColorTextSecondary
									label.Alignment = text.End // Alinhar à direita sob a barra
									return label.Layout(gtx)
								}),
							)
						})
					}
					return layout.Dimensions{}
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if !isProcessing && (filesFound > 0 || filesGenerated > 0) {
						return layout.Inset{Top: mediumPadding}.Layout(gtx, a.layoutStats)
					}
					return layout.Dimensions{}
				}),
			)
		})
	})
}

func (a *App) layoutStats(gtx layout.Context) layout.Dimensions {
	a.mu.RLock()
	filesFound := a.filesFound
	filesGenerated := a.filesGenerated
	a.mu.RUnlock()

	return layout.Flex{Spacing: layout.SpaceAround}.Layout(gtx, // SpaceAround para melhor distribuição
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return a.layoutStatCard(gtx, "Arquivos Encontrados", fmt.Sprintf("%d", filesFound), ColorPrimary)
		}),
		layout.Rigid(layout.Spacer{Width: mediumPadding}.Layout),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return a.layoutStatCard(gtx, "Contextos Gerados", fmt.Sprintf("%d", filesGenerated), ColorSuccess)
		}),
	)
}

func (a *App) layoutStatCard(gtx layout.Context, label, value string, valueColor color.NRGBA) layout.Dimensions {
	return Card{ // Usar Card para stat cards também, para consistência
		Color:        ColorBackground, // Fundo sutilmente diferente da Surface
		CornerRadius: smallRadius,
		Elevation:    0, // Sem elevação, para parecer "embutido" no card de status
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(mediumPadding).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					valueLabel := material.H5(a.theme, value) // Valor maior
					valueLabel.Color = valueColor
					valueLabel.Alignment = text.Middle
					return valueLabel.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					labelText := material.Body2(a.theme, label) // Label um pouco maior que Caption
					labelText.Color = ColorTextSecondary
					labelText.Alignment = text.Middle
					return labelText.Layout(gtx)
				}),
			)
		})
	})
}

func (a *App) layoutGenerateButton(gtx layout.Context) layout.Dimensions {
	a.mu.RLock()
	canGenerate := a.srcPath != "" && a.destPath != "" && !a.isProcessing
	isProcessing := a.isProcessing
	a.mu.RUnlock()

	btnText := "🚀 Gerar Arquivos de Contexto"
	btnBgColor := ColorSuccess
	btnFgColor := ColorTextOnSuccess
	clickable := &a.generateBtn

	if isProcessing {
		btnText = "⏳ Processando..."
		btnBgColor = ColorWarning
		btnFgColor = ColorTextPrimary
		clickable = &widget.Clickable{} // Desabilitar clique
	} else if !canGenerate {
		btnText = "⚠️ Selecione as pastas"
		btnBgColor = ColorBorder
		btnFgColor = ColorTextMuted
		clickable = &widget.Clickable{} // Desabilitar clique
	}

	btn := material.Button(a.theme, clickable, btnText)
	btn.Color = btnFgColor
	btn.CornerRadius = smallRadius
	// Padding maior para o botão de ação principal
	btn.Inset = layout.Inset{Top: unit.Dp(14), Bottom: unit.Dp(14), Left: largePadding, Right: largePadding}
	btn.TextSize = unit.Sp(16) // Fonte maior para o botão principal

	// Aplicar cores de hover/pressed se o botão estiver ativo
	if canGenerate && !isProcessing {
		btn.Background = ColorSuccess
		if clickable.Pressed() {
			btn.Background = ColorSuccessPressed
		} else if clickable.Hovered() {
			btn.Background = ColorSuccessHover
		}
		// Envolver em um Card para dar elevação quando ativo
		return Card{
			Color:        btn.Background, // O Card em si será da cor do botão
			CornerRadius: btn.CornerRadius,
			Elevation:    unit.Dp(4), // Elevação para o botão principal
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			btn.Background = ColorTransparent // Botão interno transparente, Card faz o fundo
			return btn.Layout(gtx)
		})
	}
	// Para estados desabilitado/processando, o material.Button já lida com isso
	btn.Background = btnBgColor // Definir cor de fundo para estados não ativos
	return btn.Layout(gtx)
}

func (a *App) layoutSettings(gtx layout.Context) layout.Dimensions {
	return Card{
		Color:        ColorSurface,
		CornerRadius: largeRadius, // Mais arredondado para o painel "modal"
		Elevation:    unit.Dp(8),  // Mais elevação para destacar
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(xlargePadding).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							// Ícone de Configurações maior
							iconLabel := material.H5(a.theme, "⚙️")
							iconLabel.Color = ColorTextPrimary
							return iconLabel.Layout(gtx)
						}),
						layout.Rigid(layout.Spacer{Width: smallPadding}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							title := material.H5(a.theme, "Configurações Avançadas")
							title.Color = ColorTextPrimary
							return title.Layout(gtx)
						}),
					)
				}),
				layout.Rigid(layout.Spacer{Height: xlargePadding}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.layoutCheckboxItem(gtx, &a.removeComments, "Remover Comentários", "Exclui a maioria dos comentários, preservando a documentação essencial (godoc).")
				}),
				layout.Rigid(layout.Spacer{Height: largePadding}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.layoutCheckboxItem(gtx, &a.includeTests, "Incluir Arquivos de Teste", "Processa arquivos de teste (ex: *_test.*, *.spec.*) juntamente com o código fonte.")
				}),
				layout.Rigid(layout.Spacer{Height: largePadding}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.layoutCheckboxItem(gtx, &a.minifyOutput, "Otimizar Saída para IA (Minify)", "Remove espaços em branco e quebras de linha desnecessários. A eficácia varia por linguagem.")
				}),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions { // Espaço flexível para empurrar o botão para baixo
					return layout.Spacer{Height: xlargePadding}.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					// Botão Voltar/Concluído alinhado à direita
					return layout.E.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(a.theme, &a.settingsBtn, "✓ Aplicar e Voltar")
						btn.Background = ColorPrimary
						if a.settingsBtn.Pressed() {
							btn.Background = ColorPrimaryPressed
						} else if a.settingsBtn.Hovered() {
							btn.Background = ColorPrimaryHover
						}
						btn.Color = ColorTextOnPrimary
						btn.CornerRadius = smallRadius
						btn.Inset = layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(10), Left: mediumPadding, Right: mediumPadding}
						return btn.Layout(gtx)
					})
				}),
			)
		})
	})
}

// layoutCheckboxItem é um helper para consistência nos itens de checkbox
func (a *App) layoutCheckboxItem(gtx layout.Context, checkbox *widget.Bool, title, description string) layout.Dimensions {
	return layout.Flex{Alignment: layout.Start}.Layout(gtx, // Alinhar checkbox com o topo do texto
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			cb := material.CheckBox(a.theme, checkbox, "") // Label vazia no widget CheckBox
			cb.Color = ColorPrimary
			cb.IconColor = ColorSurface // Cor do "check" mark
			// cb.Size = unit.Dp(22) // Tamanho do checkbox
			return layout.Inset{Right: mediumPadding}.Layout(gtx, cb.Layout)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					titleLabel := material.Label(a.theme, unit.Sp(16), title) // Usar material.Label para texto
					titleLabel.Color = ColorTextPrimary
					// titleLabel.Font.Weight = text.Medium // Se precisar de mais destaque
					return titleLabel.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					descLabel := material.Caption(a.theme, description)
					descLabel.Color = ColorTextSecondary
					return descLabel.Layout(gtx)
				}),
			)
		}),
	)
}

func (a *App) layoutFooter(gtx layout.Context) layout.Dimensions {
	return layout.Inset{Top: smallPadding}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Alignment: layout.Middle, Spacing: layout.SpaceBetween}.Layout(gtx,
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				a.mu.RLock()
				lastRun := a.lastRun
				a.mu.RUnlock()
				if lastRun != "" {
					label := material.Caption(a.theme, "Última execução: "+lastRun)
					label.Color = ColorTextMuted
					return label.Layout(gtx)
				}
				return layout.Dimensions{}
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				version := material.Caption(a.theme, "v2.1 Modern") // Atualizar versão se quiser
				version.Color = ColorTextMuted
				return version.Layout(gtx)
			}),
		)
	})
}
