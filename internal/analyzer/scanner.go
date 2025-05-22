package analyzer

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type ScanConfig struct {
	IncludeTests   bool
	RemoveComments bool
	MinifyOutput   bool
}

type Scanner struct {
	config ScanConfig
	fset   *token.FileSet
}

type GoFile struct {
	Path         string
	Name         string
	Package      string
	Imports      []string
	Dependencies []string
	Content      string
	CleanContent string
	AST          *ast.File
	Size         int64
	LOC          int // Lines of Code
}

func NewScanner(config ScanConfig) *Scanner {
	return &Scanner{
		config: config,
		fset:   token.NewFileSet(),
	}
}

func (s *Scanner) ScanDirectory(dir string) ([]*GoFile, error) {
	var files []*GoFile

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if s.shouldSkipPath(path, d) {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		if !d.IsDir() && strings.HasSuffix(path, ".go") {
			goFile, err := s.parseGoFile(path)
			if err != nil {
				// Log error mas continua o processamento
				return nil
			}

			files = append(files, goFile)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Resolver dependências entre arquivos
	s.resolveDependencies(files, dir)

	return files, nil
}

func (s *Scanner) shouldSkipPath(path string, d fs.DirEntry) bool {
	name := d.Name()

	// Pular diretórios irrelevantes
	if d.IsDir() {
		skipDirs := []string{
			"vendor", ".git", ".svn", ".hg",
			"node_modules", ".vscode", ".idea",
			"bin", "build", "dist", "target",
			"tmp", "temp", ".tmp", "__pycache__",
			".DS_Store", "Thumbs.db",
		}

		for _, skipDir := range skipDirs {
			if name == skipDir || (strings.HasPrefix(name, ".") && name != ".") {
				return true
			}
		}
		return false
	}

	// Pular arquivos irrelevantes
	if !strings.HasSuffix(path, ".go") {
		return true
	}

	// Pular testes se não configurado para incluir
	if !s.config.IncludeTests && strings.HasSuffix(path, "_test.go") {
		return true
	}

	// Pular arquivos específicos
	skipFiles := []string{
		"doc.go", // arquivos só de documentação
	}

	baseName := filepath.Base(path)
	for _, skipFile := range skipFiles {
		if baseName == skipFile {
			return true
		}
	}

	return false
}

func (s *Scanner) parseGoFile(filePath string) (*GoFile, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Obter tamanho do arquivo
	stat, _ := os.Stat(filePath)

	node, err := parser.ParseFile(s.fset, filePath, content, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	goFile := &GoFile{
		Path:    filePath,
		Name:    filepath.Base(filePath),
		Package: node.Name.Name,
		Content: string(content),
		AST:     node,
		Size:    stat.Size(),
		LOC:     s.countLines(string(content)),
	}

	// Extrair imports
	for _, imp := range node.Imports {
		importPath := strings.Trim(imp.Path.Value, `"`)
		goFile.Imports = append(goFile.Imports, importPath)
	}

	// Limpar conteúdo para IA
	goFile.CleanContent = s.cleanContentForAI(string(content))

	return goFile, nil
}

func (s *Scanner) countLines(content string) int {
	return len(strings.Split(content, "\n"))
}

func (s *Scanner) cleanContentForAI(content string) string {
	if !s.config.RemoveComments && !s.config.MinifyOutput {
		return content
	}

	lines := strings.Split(content, "\n")
	var cleanLines []string

	inBlockComment := false
	inStringLiteral := false
	emptyLineCount := 0

	for _, line := range lines {
		originalLine := line
		trimmed := strings.TrimSpace(line)

		// Pular linhas vazias excessivas
		if trimmed == "" {
			if s.config.MinifyOutput {
				emptyLineCount++
				if emptyLineCount <= 1 { // Permitir apenas uma linha vazia consecutiva
					cleanLines = append(cleanLines, "")
				}
				continue
			} else {
				cleanLines = append(cleanLines, line)
				continue
			}
		}
		emptyLineCount = 0

		// Detectar comentários de bloco
		if strings.Contains(line, "/*") && !inStringLiteral {
			inBlockComment = true
			if s.config.RemoveComments && !s.isImportantComment(line) {
				continue
			}
		}

		if inBlockComment {
			if strings.Contains(line, "*/") {
				inBlockComment = false
			}
			if s.config.RemoveComments && !s.isImportantComment(line) {
				continue
			}
		}

		// Processar comentários de linha
		if strings.Contains(line, "//") && !inStringLiteral {
			if s.config.RemoveComments {
				line = s.removeLineComments(line)
				if strings.TrimSpace(line) == "" {
					continue
				}
			}
		}

		// Otimizar espaçamento se minificando
		if s.config.MinifyOutput {
			line = s.optimizeWhitespaceForAI(originalLine)
		}

		cleanLines = append(cleanLines, line)
	}

	// Remover linhas vazias no final
	for len(cleanLines) > 0 && strings.TrimSpace(cleanLines[len(cleanLines)-1]) == "" {
		cleanLines = cleanLines[:len(cleanLines)-1]
	}

	return strings.Join(cleanLines, "\n")
}

func (s *Scanner) removeLineComments(line string) string {
	// Procurar por // que não estão dentro de strings
	inString := false
	inRawString := false
	escaped := false

	for i, char := range line {
		if escaped {
			escaped = false
			continue
		}

		if char == '\\' && inString && !inRawString {
			escaped = true
			continue
		}

		if char == '`' {
			inRawString = !inRawString
			continue
		}

		if char == '"' && !inRawString {
			inString = !inString
			continue
		}

		if !inString && !inRawString && i < len(line)-1 {
			if line[i] == '/' && line[i+1] == '/' {
				comment := line[i:]
				if s.isImportantComment(comment) {
					return line // Manter comentário importante
				}
				// Remover comentário, mantendo código antes dele
				codePart := strings.TrimRight(line[:i], " \t")
				return codePart
			}
		}
	}

	return line
}

func (s *Scanner) isImportantComment(comment string) bool {
	lower := strings.ToLower(comment)

	// Comentários importantes para manter
	importantKeywords := []string{
		"//go:", "//go:generate", "//go:build", "//go:embed",
		"+build", "package ", "copyright", "license", "author",
		"todo", "fixme", "note", "warning", "deprecated",
		"bug", "hack", "important", "security", "performance",
		"api", "public", "exported", "interface",
	}

	for _, keyword := range importantKeywords {
		if strings.Contains(lower, keyword) {
			return true
		}
	}

	// Comentários de documentação (antes de declarações)
	if strings.HasPrefix(lower, "// ") && (strings.Contains(lower, "func ") ||
		strings.Contains(lower, "type ") ||
		strings.Contains(lower, "var ") ||
		strings.Contains(lower, "const ")) {
		return true
	}

	return false
}

func (s *Scanner) optimizeWhitespaceForAI(line string) string {
	if strings.TrimSpace(line) == "" {
		return ""
	}

	// Preservar indentação básica mas otimizar
	leadingSpaces := len(line) - len(strings.TrimLeft(line, " \t"))
	content := strings.TrimSpace(line)

	// Reduzir indentação excessiva
	indent := leadingSpaces
	if indent > 16 { // Máximo de 16 espaços de indentação
		indent = 16
	}

	// Usar tabs para economizar tokens
	tabCount := indent / 4
	spaceCount := indent % 4

	return strings.Repeat("\t", tabCount) + strings.Repeat(" ", spaceCount) + content
}

func (s *Scanner) resolveDependencies(files []*GoFile, projectDir string) {
	projectModule := s.getProjectModule(projectDir)
	fileMap := make(map[string]*GoFile)
	packageFiles := make(map[string][]*GoFile)

	// Criar mapas
	for _, file := range files {
		fileMap[file.Path] = file
		packageFiles[file.Package] = append(packageFiles[file.Package], file)
	}

	// Resolver dependências
	for _, file := range files {
		dependencies := make(map[string]bool)

		for _, imp := range file.Imports {
			// Dependências locais do projeto
			if strings.HasPrefix(imp, projectModule) {
				relatedFiles := s.findRelatedFiles(imp, projectModule, projectDir, fileMap)
				for _, related := range relatedFiles {
					if related != file.Path { // Não incluir o próprio arquivo
						dependencies[related] = true
					}
				}
			}
		}

		// Converter map para slice ordenado
		for dep := range dependencies {
			file.Dependencies = append(file.Dependencies, dep)
		}
	}
}

func (s *Scanner) getProjectModule(projectDir string) string {
	modPath := filepath.Join(projectDir, "go.mod")
	if content, err := os.ReadFile(modPath); err == nil {
		// Usar regex para extrair módulo de forma mais robusta
		re := regexp.MustCompile(`module\s+([^\s\n]+)`)
		matches := re.FindStringSubmatch(string(content))
		if len(matches) > 1 {
			return strings.TrimSpace(matches[1])
		}
	}
	return filepath.Base(projectDir)
}

func (s *Scanner) findRelatedFiles(importPath, projectModule, projectDir string, fileMap map[string]*GoFile) []string {
	var relatedFiles []string

	// Remover o módulo do projeto do import path
	relPath := strings.TrimPrefix(importPath, projectModule)
	if strings.HasPrefix(relPath, "/") {
		relPath = relPath[1:]
	}

	fullPath := filepath.Join(projectDir, relPath)

	// Verificar se é um diretório
	if stat, err := os.Stat(fullPath); err == nil && stat.IsDir() {
		// Procurar arquivos Go no diretório
		filepath.WalkDir(fullPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}

			if strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
				if _, exists := fileMap[path]; exists {
					relatedFiles = append(relatedFiles, path)
				}
			}
			return nil
		})
	} else {
		// Pode ser um arquivo específico
		goFile := fullPath + ".go"
		if _, exists := fileMap[goFile]; exists {
			relatedFiles = append(relatedFiles, goFile)
		}
	}

	return relatedFiles
}
