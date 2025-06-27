package main

import (
	"encoding/xml"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	// 🔧 Пути к файлам
	path = "C:\\Games\\DayZServer\\mpmissions\\dayzOffline.chernarusplus\\db\\types.xml"

	// Коэффициент лута (на него умножается исходное значение, нужно всегда иметь backup файл с ванильными значениями, либо перегенерировать его
	scaleFactor = 0.5
)

type Types struct {
	XMLName xml.Name `xml:"types"`
	Types   []Type   `xml:"type"`
}

type Type struct {
	XMLName  xml.Name `xml:"type"`
	Name     string   `xml:"name,attr"`
	InnerXML string   `xml:",innerxml"`
}

func scaleNominal(innerXML string) string {
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

func main() {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Ошибка чтения файла: %v\n", err)
		return
	}

	var types Types
	if err := xml.Unmarshal(data, &types); err != nil {
		fmt.Printf("Ошибка парсинга XML: %v\n", err)
		return
	}

	for i := range types.Types {
		types.Types[i].InnerXML = scaleNominal(types.Types[i].InnerXML)
	}

	output, err := xml.MarshalIndent(types, "", "  ")
	if err != nil {
		fmt.Printf("Ошибка сериализации XML: %v\n", err)
		return
	}
	output = append([]byte(xml.Header), output...)

	backupPath := filepath.Join(filepath.Dir(path), "types_backup.xml")
	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		fmt.Printf("Ошибка создания бэкапа: %v\n", err)
		return
	}

	if err := os.WriteFile(path, output, 0644); err != nil {
		fmt.Printf("Ошибка записи файла: %v\n", err)
		return
	}

	fmt.Println("Файл успешно обновлён. Бэкап сохранён как:", backupPath)
}
