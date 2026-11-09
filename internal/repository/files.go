package repository

import (
	"os"
	"path/filepath"

	"github.com/skvdmt/auth-back/internal/model"
)

// File Репозиторий файлов.
type Files struct {
	documentation string
}

// NewFiles Конструктор.
func NewFiles() (*Files, error) {
	model.Log.Info.Info("repository files creating")
	f := &Files{}
	if err := f.updateDocumentation(); err != nil {
		return nil, err
	}
	return f, nil
}

// Documentation Документация.
func (f *Files) Documentation() string {
	return f.documentation
}

// updateDocumentation Обновление документации.
func (f *Files) updateDocumentation() error {
	model.Log.Info.Info("update documentation starting")
	const (
		name = "openapi.yaml"
		dev  = "."
		prod = "/usr/local/share/doc"
	)
	m, o := os.LookupEnv(model.MODE)
	dir := prod
	if o && m == model.MODE_DEV {
		dir = dev
	}
	doc, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		return err
	}
	f.documentation = string(doc)
	return nil
}
