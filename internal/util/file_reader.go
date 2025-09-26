package util

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ledongthuc/pdf"
)

// FileReader handles reading and extracting text from various file types
type FileReader struct{}

// NewFileReader creates a new FileReader instance
func NewFileReader() *FileReader {
	return &FileReader{}
}

// ReadFile reads and extracts text content from a file based on its extension
func (fr *FileReader) ReadFile(filePath string) (string, error) {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".pdf":
		return fr.readPDF(filePath)
	case ".txt":
		return fr.readText(filePath)
	default:
		return "", fmt.Errorf("unsupported file type: %s", ext)
	}
}

// readPDF extracts text content from a PDF file
func (fr *FileReader) readPDF(filePath string) (string, error) {
	file, reader, err := pdf.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open PDF file: %w", err)
	}
	defer file.Close()

	var content strings.Builder
	totalPages := reader.NumPage()

	for i := 1; i <= totalPages; i++ {
		page := reader.Page(i)
		if page.V.IsNull() {
			continue
		}

		text, err := page.GetPlainText(nil)
		if err != nil {
			continue // Skip pages with extraction errors
		}

		content.WriteString(text)
		content.WriteString("\n")
	}

	return content.String(), nil
}

// readText reads content from a text file
func (fr *FileReader) readText(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open text file: %w", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read text file: %w", err)
	}

	return string(content), nil
}
