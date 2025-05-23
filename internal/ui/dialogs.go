package ui

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/sqweek/dialog"
)

// selectFolder usa a biblioteca github.com/sqweek/dialog para abrir um diálogo nativo de seleção de pasta.
// Retorna o caminho da pasta selecionada ou uma string vazia se o usuário cancelar ou ocorrer um erro.
func selectFolder(title string) string {
	// dialog.Directory().Browse() é uma chamada bloqueante.
	// Deve ser chamada em uma goroutine para não bloquear a UI thread.
	path, err := dialog.Directory().Title(title).Browse()

	if err != nil {
		if err == dialog.ErrCancelled {
			// Usuário cancelou a seleção, não é um erro crítico.
			log.Println("Seleção de pasta cancelada pelo usuário.")
			return ""
		}
		// Outro erro ocorreu durante a seleção da pasta.
		log.Printf("Erro ao selecionar pasta: %v", err)
		// Poderia-se adicionar uma notificação ao usuário aqui, se desejado.
		return ""
	}

	// A biblioteca geralmente garante que o caminho é válido e existe,
	// mas uma verificação extra com dirExists pode ser mantida por segurança adicional.
	if path != "" && !dirExists(path) {
		log.Printf("Caminho selecionado pela biblioteca de diálogo não é um diretório válido: %s", path)
		return ""
	}

	return path
}

// dirExists verifica se um caminho existe e é um diretório.
func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		// Se os.Stat retorna um erro, o caminho não existe ou não é acessível.
		return false
	}
	return info.IsDir()
}

// getDefaultDestPath retorna um caminho padrão para a pasta de destino dos arquivos de contexto.
// Cria uma pasta "go-contexts" no Desktop ou Documentos do usuário, dependendo do SO.
func getDefaultDestPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		// Se não conseguir o diretório home, usa um subdiretório local.
		destPath := "go-contexts"
		if err := os.MkdirAll(destPath, 0755); err != nil {
			// Se mesmo assim falhar, retorna o diretório atual relativo.
			return "."
		}
		return destPath
	}

	var contextDir string
	switch runtime.GOOS {
	case "windows":
		// Tenta criar no Desktop primeiro
		desktopPath := filepath.Join(home, "Desktop", "go-contexts")
		if err := os.MkdirAll(desktopPath, 0755); err == nil {
			contextDir = desktopPath
		} else {
			// Fallback para Documentos se o Desktop falhar ou não for acessível
			documentsPath := filepath.Join(home, "Documents", "go-contexts")
			if err := os.MkdirAll(documentsPath, 0755); err == nil {
				contextDir = documentsPath
			} else {
				// Fallback extremo se ambos falharem
				contextDir = filepath.Join(home, "go-contexts")
			}
		}
	case "darwin":
		// No macOS, geralmente é comum usar o Desktop
		contextDir = filepath.Join(home, "Desktop", "go-contexts")
	default: // linux e outros
		// Para Linux e outros, um diretório diretamente no home é comum
		contextDir = filepath.Join(home, "go-contexts")
	}

	// Tenta criar o diretório de destino final
	if err := os.MkdirAll(contextDir, 0755); err != nil {
		// Se a criação falhar (ex: permissões), tenta um fallback no diretório de trabalho atual
		cwd, cwdErr := os.Getwd()
		if cwdErr == nil {
			fallbackPath := filepath.Join(cwd, "go-contexts")
			if err := os.MkdirAll(fallbackPath, 0755); err == nil {
				return fallbackPath
			}
		}
		// Último recurso, apenas o nome da pasta (será relativo ao local de execução)
		return "go-contexts"
	}

	return contextDir
}
