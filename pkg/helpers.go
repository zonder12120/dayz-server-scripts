package pkg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// LoadConfig загружает конфигурацию из файла YAML
func LoadConfig(configPath string, out any) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, out)
	return err
}

// GetBackupPathWithIndex генерирует путь и имя для бэкапа, добавляя индекс, если файл уже существует
func GetBackupPathWithIndex(basePath string) string {
	ext := filepath.Ext(basePath)
	name := strings.TrimSuffix(filepath.Base(basePath), ext)
	dir := filepath.Dir(basePath)

	backupPath := filepath.Join(dir, name+ext)
	i := 1
	for {
		if _, err := os.Stat(backupPath); os.IsNotExist(err) {
			break
		}
		backupPath = filepath.Join(dir, fmt.Sprintf("%s (%d)%s", name, i, ext))
		i++
	}
	return backupPath
}
