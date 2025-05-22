package ui

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"go-context-generator/internal/analyzer"
	"go-context-generator/internal/config"
	"go-context-generator/internal/generator"
)

type App struct {
	theme    *material.Theme
	settings *config.Settings

	// UI Components
	selectSrcBtn  widget.Clickable
	selectDestBtn widget.Clickable
	generateBtn   widget.Clickable
	settingsBtn   widget.Clickable

	// UI State
	srcPath        string
	destPath       string
	lastRun        string
	status         string
	isProcessing   bool
	progress       float32
	filesFound     int
	filesGenerated int

	// Settings UI
	showSettings   bool
	removeComments widget.Bool
	includeTests   widget.Bool
	minifyOutput   widget.Bool

	// Background processing
	ctx    context.Context
	cancel context.CancelFunc
	mu     sync.RWMutex
}

func NewApp(theme *material.Theme) *App {
	ctx, cancel := context.WithCancel(context.Background())

	settings := config.LoadSettings()

	app := &App{
		theme:    theme,
		settings: settings,
		ctx:      ctx,
		cancel:   cancel,
		status:   "Pronto para começar! Selecione as pastas de origem e destino.",
	}

	// Aplicar configurações salvas
	app.removeComments.Value = settings.RemoveComments
	app.includeTests.Value = settings.IncludeTests
	app.minifyOutput.Value = settings.MinifyOutput

	// Restaurar caminhos salvos se existirem
	if settings.LastSrcPath != "" {
		app.srcPath = settings.LastSrcPath
	}
	if settings.LastDestPath != "" {
		app.destPath = settings.LastDestPath
	} else {
		app.destPath = getDefaultDestPath()
	}

	return app
}

func (a *App) Run(w *app.Window) error {
	defer a.cancel()
	var ops op.Ops

	for {
		e := w.NextEvent()
		switch e := e.(type) {
		case system.DestroyEvent:
			a.saveSettings()
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			a.handleEvents(gtx)
			a.layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
}

func (a *App) handleEvents(gtx layout.Context) {
	// Botão de seleção de origem
	if a.selectSrcBtn.Clicked(gtx) {
		go func() {
			if path := selectFolder("Selecione a pasta com o código"); path != "" {
				a.mu.Lock()
				a.srcPath = path
				a.status = fmt.Sprintf("✓ Origem selecionada: %s", filepath.Base(path))
				a.mu.Unlock()
			}
		}()
	}

	// Botão de seleção de destino
	if a.selectDestBtn.Clicked(gtx) {
		go func() {
			if path := selectFolder("Selecione onde salvar os arquivos de contexto"); path != "" {
				a.mu.Lock()
				a.destPath = path
				a.status = fmt.Sprintf("✓ Destino selecionado: %s", filepath.Base(path))
				a.mu.Unlock()
			}
		}()
	}

	// Botão de configurações
	if a.settingsBtn.Clicked(gtx) {
		a.showSettings = !a.showSettings
	}

	// Botão de geração
	if a.generateBtn.Clicked(gtx) && a.canGenerate() {
		go a.generateContextFiles()
	}

	// Atualizar configurações quando mudarem
	oldRemoveComments := a.removeComments.Value
	oldIncludeTests := a.includeTests.Value
	oldMinifyOutput := a.minifyOutput.Value

	if a.removeComments.Value != oldRemoveComments {
		a.settings.RemoveComments = a.removeComments.Value
	}
	if a.includeTests.Value != oldIncludeTests {
		a.settings.IncludeTests = a.includeTests.Value
	}
	if a.minifyOutput.Value != oldMinifyOutput {
		a.settings.MinifyOutput = a.minifyOutput.Value
	}
}

func (a *App) canGenerate() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.srcPath != "" && a.destPath != "" && !a.isProcessing
}

func (a *App) generateContextFiles() {
	a.mu.Lock()
	a.isProcessing = true
	a.progress = 0
	a.filesFound = 0
	a.filesGenerated = 0
	a.status = "🔍 Escaneando arquivos Go..."
	a.mu.Unlock()

	defer func() {
		a.mu.Lock()
		a.isProcessing = false
		a.lastRun = time.Now().Format("15:04 - 02/01/2006")

		if a.filesGenerated > 0 {
			a.progress = 1.0
			a.status = fmt.Sprintf("✅ Concluído! %d arquivos de contexto gerados", a.filesGenerated)
		}
		a.mu.Unlock()
	}()

	// Salvar configurações atuais
	a.saveSettings()

	// Criar scanner com configurações
	scanner := analyzer.NewScanner(analyzer.ScanConfig{
		IncludeTests:   a.settings.IncludeTests,
		RemoveComments: a.settings.RemoveComments,
		MinifyOutput:   a.settings.MinifyOutput,
	})

	// Escanear arquivos
	files, err := scanner.ScanDirectory(a.srcPath)
	if err != nil {
		a.mu.Lock()
		a.status = "❌ Erro ao escanear arquivos: " + err.Error()
		a.mu.Unlock()
		return
	}

	if len(files) == 0 {
		a.mu.Lock()
		a.status = "⚠️ Nenhum arquivo Go encontrado na pasta selecionada"
		a.mu.Unlock()
		return
	}

	a.mu.Lock()
	a.filesFound = len(files)
	a.status = fmt.Sprintf("🔄 Encontrados %d arquivos Go. Gerando contextos...", len(files))
	a.mu.Unlock()

	// Gerar arquivos de contexto
	gen := generator.NewGenerator(generator.Config{
		OutputDir:      a.destPath,
		SourceDir:      a.srcPath,
		RemoveComments: a.settings.RemoveComments,
		MinifyOutput:   a.settings.MinifyOutput,
	})

	gen.SetProgressCallback(func(current, total int) {
		a.mu.Lock()
		if total > 0 {
			a.progress = float32(current) / float32(total)
		}
		a.filesGenerated = current
		if current < total {
			a.status = fmt.Sprintf("⚡ Processando... %d/%d arquivos", current, total)
		}
		a.mu.Unlock()
	})

	if err := gen.GenerateContextFiles(files); err != nil {
		a.mu.Lock()
		a.status = "❌ Erro na geração: " + err.Error()
		a.mu.Unlock()
		return
	}
}

func (a *App) saveSettings() {
	// Atualizar paths nas configurações
	a.settings.LastSrcPath = a.srcPath
	a.settings.LastDestPath = a.destPath

	// Salvar configurações
	a.settings.Save()
}

func (a *App) getStatusColor() ColorRGBA {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.isProcessing {
		return ColorPrimary
	}

	if strings.Contains(a.status, "❌") {
		return ColorDanger
	}

	if strings.Contains(a.status, "✅") {
		return ColorSuccess
	}

	if strings.Contains(a.status, "⚠️") {
		return ColorWarning
	}

	return ColorTextSecondary
}
