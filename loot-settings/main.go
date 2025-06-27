package main

import (
	"encoding/xml"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Path        string  `yaml:"path"`
	ScaleFactor float64 `yaml:"scale_factor"`
}

type Types struct {
	XMLName xml.Name `xml:"types"`
	Types   []Type   `xml:"type"`
}

type Type struct {
	XMLName  xml.Name `xml:"type"`
	Name     string   `xml:"name,attr"`
	InnerXML string   `xml:",innerxml"`
}

func scaleNominal(innerXML string, scaleFactor float64) string {
	lines := strings.Split(innerXML, "\n")
	for i, line := range lines {
		lineTrimmed := strings.TrimSpace(line)
		if strings.HasPrefix(lineTrimmed, "<nominal>") && strings.HasSuffix(lineTrimmed, "</nominal>") {
			valStr := strings.TrimPrefix(lineTrimmed, "<nominal>")
			valStr = strings.TrimSuffix(valStr, "</nominal>")
			valStr = strings.TrimSpace(valStr)

			if val, err := strconv.Atoi(valStr); err == nil {
				newVal := int(math.Ceil(float64(val) * scaleFactor))
				if newVal < 1 {
					newVal = 1
				}
				lines[i] = fmt.Sprintf("  <nominal>%d</nominal>", newVal)
			}
		}
	}
	return strings.Join(lines, "\n")
}

func loadConfig(configPath string) (Config, error) {
	var cfg Config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return cfg, err
	}
	err = yaml.Unmarshal(data, &cfg)
	return cfg, err
}

func main() {
	cfg, err := loadConfig("config.yml")
	if err != nil {
		fmt.Printf("Ошибка загрузки конфигурации: %v\n", err)
		return
	}

	inputPath := filepath.Join(filepath.Dir(cfg.Path), "types.xml")

	backupBase := filepath.Join(filepath.Dir(cfg.Path), "backups", "types_backup.xml")
	backupPath := getBackupPathWithIndex(backupBase)

	data, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Printf("Ошибка чтения файла: %v\n", err)
		return
	}

	var types Types
	if err = xml.Unmarshal(data, &types); err != nil {
		fmt.Printf("Ошибка парсинга XML: %v\n", err)
		return
	}

	for i := range types.Types {
		types.Types[i].InnerXML = scaleNominal(types.Types[i].InnerXML, cfg.ScaleFactor)
	}

	output, err := xml.MarshalIndent(types, "", "  ")
	if err != nil {
		fmt.Printf("Ошибка сериализации XML: %v\n", err)
		return
	}
	output = append([]byte(xml.Header), output...)

	if err = os.WriteFile(backupPath, data, 0644); err != nil {
		fmt.Printf("Ошибка создания бэкапа: %v\n", err)
		return
	}

	if err = os.WriteFile(inputPath, output, 0644); err != nil {
		fmt.Printf("Ошибка записи файла: %v\n", err)
		return
	}

	fmt.Println("🎉 Готово! Бэкап сохранён как:", backupPath)
}

func getBackupPathWithIndex(basePath string) string {
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
