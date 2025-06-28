package main

import (
	"encoding/xml"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/zonder12120/dayz-server-scripts/pkg"
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
	XMLName    xml.Name `xml:"type"`
	Name       string   `xml:"name,attr"`
	NominalTag *Nominal `xml:"nominal"`
	Rest       string   `xml:",innerxml"`
}

type Nominal struct {
	Value int `xml:",chardata"`
}

// scaleNominalValue вычисляет новое значение nominal, масштабируя его
func scaleNominalValue(value int, scaleFactor float64) int {
	newValue := int(math.Ceil(float64(value) * scaleFactor))
	if newValue < 1 {
		return 1
	}
	return newValue
}

func main() {
	// Загрузка конфигурации
	var cfg Config
	err := pkg.LoadConfig("config.yml", &cfg)
	if err != nil {
		fmt.Printf("Ошибка загрузки конфигурации: %v\n", err)
		return
	}

	// Определение путей к файлам
	inputPath := filepath.Join(filepath.Dir(cfg.Path), "types.xml")
	backupBase := filepath.Join(filepath.Dir(cfg.Path), "backups", "types_backup.xml")

	// Чтение XML-файла
	data, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Printf("Ошибка чтения исходного XML-файла: %v\n", err)
		return
	}

	// Создание бэкапа
	backupPath := pkg.GetBackupPathWithIndex(backupBase)
	if err = os.MkdirAll(filepath.Dir(backupPath), 0755); err != nil {
		fmt.Printf("Ошибка создания директории для бэкапа: %v\n", err)
		return
	}
	if err = os.WriteFile(backupPath, data, 0644); err != nil {
		fmt.Printf("Ошибка создания бэкапа: %v\n", err)
		return
	}
	fmt.Println("💾 Бэкап сохранён как:", backupPath)

	var types Types
	decoder := xml.NewDecoder(strings.NewReader(string(data)))
	decoder.Entity = xml.HTMLEntity
	decoder.Strict = false
	if err = decoder.Decode(&types); err != nil {
		fmt.Printf("Ошибка парсинга XML: %v\n", err)
		return
	}

	// Масштабирование номинальных значений
	for i := range types.Types {
		if types.Types[i].NominalTag != nil {
			types.Types[i].NominalTag.Value = scaleNominalValue(types.Types[i].NominalTag.Value, cfg.ScaleFactor)
		}
	}

	output, err := xml.MarshalIndent(types, "", "  ")
	if err != nil {
		fmt.Printf("Ошибка сериализации XML: %v\n", err)
		return
	}

	output = append([]byte(xml.Header), output...)

	if err = os.WriteFile(inputPath, output, 0644); err != nil {
		fmt.Printf("Ошибка записи файла: %v\n", err)
		return
	}

	fmt.Println("🎉 Готово! Значения nominal успешно обновлены.")
}
