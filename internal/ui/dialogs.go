package ui

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// selectFolder abre um diálogo nativo para seleção de pasta
func selectFolder(title string) string {
	switch runtime.GOOS {
	case "windows":
		return selectFolderWindows(title)
	case "darwin":
		return selectFolderMacOS(title)
	case "linux":
		return selectFolderLinux(title)
	default:
		return selectFolderLinux(title) // fallback
	}
}

// selectFolderWindows usa PowerShell para seleção de pasta no Windows
func selectFolderWindows(title string) string {
	// Tentar método moderno primeiro (Windows 10+)
	script := `
Add-Type -AssemblyName System.Windows.Forms
$FolderBrowser = New-Object System.Windows.Forms.FolderBrowserDialog
$FolderBrowser.Description = "` + title + `"
$FolderBrowser.ShowNewFolderButton = $true
$FolderBrowser.RootFolder = [System.Environment+SpecialFolder]::MyComputer
if($FolderBrowser.ShowDialog() -eq [System.Windows.Forms.DialogResult]::OK) {
    Write-Output $FolderBrowser.SelectedPath
}
`
	cmd := exec.Command("powershell", "-WindowStyle", "Hidden", "-Command", script)
	output, err := cmd.Output()
	if err == nil {
		path := strings.TrimSpace(string(output))
		if path != "" && dirExists(path) {
			return path
		}
	}

	// Fallback para método mais antigo
	script2 := `
$shell = New-Object -ComObject Shell.Application
$folder = $shell.BrowseForFolder(0, "` + title + `", 0, 0)
if ($folder -ne $null) {
    Write-Output $folder.Self.Path
}
`
	cmd2 := exec.Command("powershell", "-WindowStyle", "Hidden", "-Command", script2)
	if output2, err2 := cmd2.Output(); err2 == nil {
		path := strings.TrimSpace(string(output2))
		if path != "" && dirExists(path) {
			return path
		}
	}

	return ""
}

// selectFolderMacOS usa osascript para seleção de pasta no macOS
func selectFolderMacOS(title string) string {
	script := `tell application "Finder"
	set folderPath to choose folder with prompt "` + title + `" default location (path to desktop)
	return POSIX path of folderPath
end tell`

	cmd := exec.Command("osascript", "-e", script)
	output, err := cmd.Output()
	if err != nil {
		// Fallback usando dialog do sistema
		script2 := `tell application "System Events"
	set folderPath to choose folder with prompt "` + title + `"
	return POSIX path of folderPath
end tell`
		cmd2 := exec.Command("osascript", "-e", script2)
		if output2, err2 := cmd2.Output(); err2 == nil {
			output = output2
		} else {
			return ""
		}
	}

	path := strings.TrimSpace(string(output))
	// Remover quebras de linha extras
	path = strings.ReplaceAll(path, "\n", "")
	path = strings.ReplaceAll(path, "\r", "")

	if path != "" && dirExists(path) {
		return path
	}
	return ""
}

// selectFolderLinux tenta usar vários diálogos disponíveis no Linux
func selectFolderLinux(title string) string {
	// Lista de diálogos para tentar, em ordem de preferência
	dialogs := []struct {
		command string
		args    []string
	}{
		// Zenity (GNOME)
		{"zenity", []string{"--file-selection", "--directory", "--title=" + title}},
		// KDialog (KDE)
		{"kdialog", []string{"--getexistingdirectory", ".", "--title", title}},
		// Yad (alternativa ao zenity)
		{"yad", []string{"--file-selection", "--directory", "--title=" + title}},
		// Dialog (terminal)
		{"dialog", []string{"--stdout", "--title", title, "--dselect", ".", "20", "60"}},
	}

	for _, d := range dialogs {
		if isCommandAvailable(d.command) {
			cmd := exec.Command(d.command, d.args...)
			output, err := cmd.Output()
			if err == nil {
				path := strings.TrimSpace(string(output))
				if path != "" && dirExists(path) {
					return path
				}
			}
		}
	}

	// Fallback final: usar diálogo de entrada simples
	return selectFolderLinuxFallback(title)
}

// selectFolderLinuxFallback usa entrada de texto como último recurso
func selectFolderLinuxFallback(title string) string {
	// Tentar zenity com entrada de texto
	if isCommandAvailable("zenity") {
		defaultPath := getDefaultSourcePath()
		cmd := exec.Command("zenity", "--entry",
			"--title="+title,
			"--text=Digite o caminho da pasta:",
			"--entry-text="+defaultPath)

		output, err := cmd.Output()
		if err == nil {
			path := strings.TrimSpace(string(output))
			if path != "" && dirExists(path) {
				return path
			}
		}
	}

	// Última tentativa: retornar diretório atual
	return getDefaultSourcePath()
}

// isCommandAvailable verifica se um comando está disponível no sistema
func isCommandAvailable(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// dirExists verifica se um diretório existe
func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// getDefaultSourcePath retorna um caminho padrão para começar a busca
func getDefaultSourcePath() string {
	// Tentar diretório atual primeiro
	if cwd, err := os.Getwd(); err == nil {
		return cwd
	}

	// Tentar diretório home
	if home, err := os.UserHomeDir(); err == nil {
		return home
	}

	// Último recurso
	return "."
}

// getDefaultDestPath retorna um caminho padrão para destino
func getDefaultDestPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}

	var contextDir string
	switch runtime.GOOS {
	case "windows":
		// Tentar Desktop primeiro, depois Documents
		desktop := filepath.Join(home, "Desktop", "go-contexts")
		if err := os.MkdirAll(desktop, 0755); err == nil {
			contextDir = desktop
		} else {
			contextDir = filepath.Join(home, "Documents", "go-contexts")
		}
	case "darwin":
		// macOS: Desktop
		contextDir = filepath.Join(home, "Desktop", "go-contexts")
	default:
		// Linux: home directory
		contextDir = filepath.Join(home, "go-contexts")
	}

	// Criar o diretório se não existir
	if err := os.MkdirAll(contextDir, 0755); err != nil {
		// Se falhar, usar diretório atual
		if cwd, err := os.Getwd(); err == nil {
			return filepath.Join(cwd, "go-contexts")
		}
		return "go-contexts"
	}

	return contextDir
}

// getSafePath sanitiza um caminho para uso em comandos
func getSafePath(path string) string {
	// Remover caracteres perigosos
	path = strings.ReplaceAll(path, `"`, `\"`)
	path = strings.ReplaceAll(path, `'`, `\'`)
	return path
}

// expandPath expande ~ para o diretório home
func expandPath(path string) string {
	if strings.HasPrefix(path, "~") {
		if home, err := os.UserHomeDir(); err == nil {
			return strings.Replace(path, "~", home, 1)
		}
	}
	return path
}
