# 🚀 Go Context Generator Pro v2.0

Uma ferramenta moderna e otimizada para gerar arquivos de contexto de projetos Go, especialmente projetada para uso com **Inteligência Artificial**.

## 🎯 Objetivo

Transformar projetos Go complexos em arquivos de contexto organizados e otimizados, onde cada arquivo contém:
- O código principal
- Todas as suas dependências locais
- Estrutura completa do projeto
- Metadados relevantes
- **Tokens otimizados para IA**

## ✨ Características Principais

- 🔍 **Análise Inteligente**: Detecta automaticamente dependências entre arquivos
- 🧹 **Otimização para IA**: Remove comentários desnecessários e otimiza tokens
- 🎨 **Interface Moderna**: UI nativa usando Gio UI
- 🌐 **Multiplataforma**: Windows, macOS e Linux
- ⚙️ **Configurável**: Opções personalizáveis de processamento
- 📊 **Estatísticas**: Análise completa do projeto

## 🔧 Instalação e Uso

### Pré-requisitos

- Go 1.21 ou superior
- Sistema operacional: Windows, macOS ou Linux

### Compilação

```bash
# Clone ou extraia o projeto
cd go-context-generator

# Baixar dependências
go mod tidy

# Compilar
go build -o go-context-generator

# Executar
./go-context-generator  # Linux/macOS
# ou
go-context-generator.exe  # Windows
```

### Como Usar

1. **Executar a aplicação**
   - Interface gráfica moderna será aberta

2. **Selecionar Pasta de Origem**
   - Clique em "Selecionar Origem"
   - Navegue até a pasta do seu projeto Go
   - Selecione a pasta raiz do projeto

3. **Selecionar Pasta de Destino**
   - Clique em "Selecionar Destino"
   - Escolha onde salvar os arquivos de contexto
   - Por padrão, cria pasta "go-contexts" no Desktop

4. **Configurar Opções** (⚙️)
   - **🧹 Remover comentários**: Remove comentários desnecessários
   - **🧪 Incluir testes**: Processa arquivos *_test.go
   - **⚡ Otimizar para IA**: Minimiza tokens extras

5. **Gerar Contextos**
   - Clique em "🚀 Gerar Arquivos de Contexto"
   - Acompanhe o progresso na interface
   - Arquivos serão salvos na pasta de destino

## 📁 Estrutura de Saída

### Arquivos Gerados

```
pasta-destino/
├── 00_PROJECT_OVERVIEW.txt           # Visão geral completa
├── main_CONTEXT.txt                  # Contexto do main.go
├── internal_ui_app_CONTEXT.txt       # Contexto de app.go
├── internal_config_settings_CONTEXT.txt
└── ...                               # Um arquivo por .go
```

### Formato dos Arquivos de Contexto

```
🎯 AI CONTEXT FILE
==================

📋 FILE METADATA
----------------
File: main.go
Package: main
Lines of Code: 45
Generated: 2025-05-22 14:30:00

📥 DEPENDENCIES
---------------
Standard Library:
  • fmt
  • os
External Packages:
  • gioui.org/app
Local Project:
  • go-context-generator/internal/ui

🏗️ PROJECT CONTEXT
------------------
main: main.go
ui: app.go, components.go, theme.go
analyzer: scanner.go
config: settings.go
generator: context.go

💻 SOURCE CODE
===============

[código limpo e otimizado aqui]

🔗 RELATED CODE
===============

--- DEPENDENCY 1: internal/ui/app.go ---
Package: ui | LOC: 150

[código das dependências aqui]
```

## 🎛️ Configurações Avançadas

### Otimizações para IA

- **Remoção Inteligente de Comentários**: Mantém documentação importante, remove comentários triviais
- **Compressão de Whitespace**: Reduz espaços desnecessários mantendo legibilidade
- **Organização Hierárquica**: Estrutura clara de dependências
- **Metadados Contextuais**: Informações essenciais para IA entender o código

### Arquivos Ignorados Automaticamente

- Diretórios: `vendor/`, `.git/`, `node_modules/`, `.vscode/`, etc.
- Arquivos: `*_test.go` (opcional), `doc.go`, etc.
- Extensions: Apenas `.go` são processados

## 🔍 Exemplo de Uso com IA

```
Prompt para IA:
"Analise este contexto de código Go e sugira melhorias na arquitetura"

[Cole o conteúdo do arquivo *_CONTEXT.txt gerado]
```

A IA receberá:
- ✅ Código completo com dependências
- ✅ Estrutura do projeto
- ✅ Metadados relevantes
- ✅ Tokens otimizados (menos custos)

## 📊 Benefícios para IA

| Aspecto | Antes | Depois |
|---------|-------|--------|
| **Tokens** | ~100% | ~60% (40% redução) |
| **Contexto** | Parcial | Completo |
| **Dependências** | Manual | Automático |
| **Estrutura** | Invisível | Clara |
| **Otimização** | Nenhuma | IA-específica |

## 🛠️ Estrutura do Código

```
go-context-generator/
├── main.go                    # Ponto de entrada
├── internal/
│   ├── analyzer/
│   │   └── scanner.go         # Análise de código Go
│   ├── config/
│   │   └── settings.go        # Configurações
│   ├── generator/
│   │   └── context.go         # Geração de contextos
│   └── ui/
│       ├── app.go            # Interface principal
│       ├── components.go     # Componentes UI
│       ├── dialogs.go        # Diálogos nativos
│       └── theme.go          # Tema e layout
└── go.mod                    # Dependências
```

## 🚀 Recursos Avançados

### Análise de Dependências
- Detecta imports locais automaticamente
- Mapeia relações entre arquivos
- Inclui código das dependências no contexto

### Interface Multiplataforma
- **Windows**: Diálogos nativos do Windows
- **macOS**: AppleScript para seleção de pastas
- **Linux**: Zenity, KDialog, Yad (com fallbacks)

### Otimizações de Performance
- Processamento paralelo de arquivos
- Cache de análise AST
- Progress feedback em tempo real

## 🔧 Solução de Problemas

### Problema: Diálogo de pasta não abre
**Solução**: 
- **Linux**: Instale `zenity` (`sudo apt install zenity`)
- **Windows**: Execute como administrador se necessário
- **macOS**: Permita acesso ao Finder nas configurações

### Problema: Arquivos não encontrados
**Solução**:
- Verifique se a pasta contém arquivos `.go`
- Certifique-se de que não é apenas uma pasta `vendor/