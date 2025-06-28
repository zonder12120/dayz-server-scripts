package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/zonder12120/dayz-server-scripts/pkg"
)

type Config struct {
	Path             string  `yaml:"path"`
	MinAI            int     `yaml:"min_ai"`
	MaxAI            int     `yaml:"max_ai"`
	RespawnTime      float64 `yaml:"respawn_time"`
	MinWaypoints     int     `yaml:"min_waypoints"`
	MaxWaypoints     int     `yaml:"max_waypoints"`
	MapMinCoord      int     `yaml:"map_min_coord"`
	MapMaxCoord      int     `yaml:"map_max_coord"`
	PatrolMultiplier int     `yaml:"patrol_multiplier"`
}

func generateWaypoint(usedCoords map[string]struct{}, cfg Config) ([]float64, error) {
	const maxAttempts = 1000
	for attempts := 0; attempts < maxAttempts; attempts++ {
		x := float64(rand.Intn(cfg.MapMaxCoord-cfg.MapMinCoord) + cfg.MapMinCoord)
		z := float64(rand.Intn(cfg.MapMaxCoord-cfg.MapMinCoord) + cfg.MapMinCoord)
		key := fmt.Sprintf("%.0f_%.0f", x, z)
		if _, exists := usedCoords[key]; !exists {
			usedCoords[key] = struct{}{}
			return []float64{x, 0.0, z}, nil
		}
	}
	return nil, fmt.Errorf("failed to generate a unique waypoint after %d attempts. Map may be saturated with points", maxAttempts)
}

type Settings struct {
	MVersion                        int      `json:"m_Version"`
	Enabled                         int      `json:"Enabled"`
	DespawnTime                     float64  `json:"DespawnTime"`
	RespawnTime                     float64  `json:"RespawnTime"`
	MinDistRadius                   float64  `json:"MinDistRadius"`
	MaxDistRadius                   float64  `json:"MaxDistRadius"`
	DespawnRadius                   float64  `json:"DespawnRadius"`
	AccuracyMin                     float64  `json:"AccuracyMin"`
	AccuracyMax                     float64  `json:"AccuracyMax"`
	ThreatDistanceLimit             float64  `json:"ThreatDistanceLimit"`
	NoiseInvestigationDistanceLimit float64  `json:"NoiseInvestigationDistanceLimit"`
	DamageMultiplier                float64  `json:"DamageMultiplier"`
	DamageReceivedMultiplier        float64  `json:"DamageReceivedMultiplier"`
	ObjectPatrols                   []any    `json:"ObjectPatrols"`
	Patrols                         []Patrol `json:"Patrols"`
}

type Patrol struct {
	Name                            string      `json:"Name"`
	Persist                         int         `json:"Persist"`
	Faction                         string      `json:"Faction"`
	Formation                       string      `json:"Formation"`
	FormationLooseness              float64     `json:"FormationLooseness"`
	LoadoutFile                     string      `json:"LoadoutFile"`
	Units                           []any       `json:"Units"`
	NumberOfAI                      int         `json:"NumberOfAI"`
	Behaviour                       string      `json:"Behaviour"`
	Speed                           string      `json:"Speed"`
	UnderThreatSpeed                string      `json:"UnderThreatSpeed"`
	CanBeLooted                     int         `json:"CanBeLooted"`
	UnlimitedReload                 int         `json:"UnlimitedReload"`
	SniperProneDistanceThreshold    float64     `json:"SniperProneDistanceThreshold"`
	AccuracyMin                     float64     `json:"AccuracyMin"`
	AccuracyMax                     float64     `json:"AccuracyMax"`
	ThreatDistanceLimit             float64     `json:"ThreatDistanceLimit"`
	NoiseInvestigationDistanceLimit float64     `json:"NoiseInvestigationDistanceLimit"`
	DamageMultiplier                float64     `json:"DamageMultiplier"`
	DamageReceivedMultiplier        float64     `json:"DamageReceivedMultiplier"`
	CanBeTriggeredByAI              int         `json:"CanBeTriggeredByAI"`
	MinDistRadius                   float64     `json:"MinDistRadius"`
	MaxDistRadius                   float64     `json:"MaxDistRadius"`
	DespawnRadius                   float64     `json:"DespawnRadius"`
	MinSpreadRadius                 float64     `json:"MinSpreadRadius"`
	MaxSpreadRadius                 float64     `json:"MaxSpreadRadius"`
	Chance                          float64     `json:"Chance"`
	WaypointInterpolation           string      `json:"WaypointInterpolation"`
	DespawnTime                     float64     `json:"DespawnTime"`
	RespawnTime                     float64     `json:"RespawnTime"`
	UseRandomWaypointAsStartPoint   int         `json:"UseRandomWaypointAsStartPoint"`
	Waypoints                       [][]float64 `json:"Waypoints"`
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
	inputPath := filepath.Join(filepath.Dir(cfg.Path), "AIPatrolSettings.json")
	backupBase := filepath.Join(filepath.Dir(cfg.Path), "backups", "AIPatrolSettings_backup.json")

	// Чтение JSON-файла
	data, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Printf("Ошибка чтения исходного JSON-файла: %v\n", err)
		return
	}

	// Создание бэкапа
	backupPath := pkg.GetBackupPathWithIndex(backupBase)
	if err = os.MkdirAll(filepath.Dir(backupPath), 0755); err != nil {
		fmt.Printf("Ошибка создания директории для бэкапа: %v\n", err)
		return
	}
	if err = os.WriteFile(backupPath, data, 0644); err != nil {
		fmt.Printf("Ошибка записи бэкапа: %v\n", err)
		return
	}
	fmt.Println("💾 Бэкап сохранён как:", backupPath)

	var settings Settings
	if err = json.Unmarshal(data, &settings); err != nil {
		fmt.Printf("Ошибка парсинга JSON: %v\n", err)
		return
	}

	usedCoords := make(map[string]struct{})
	for _, p := range settings.Patrols {
		for _, wp := range p.Waypoints {
			key := fmt.Sprintf("%.0f_%.0f", wp[0], wp[2])
			usedCoords[key] = struct{}{}
		}
	}

	originalCount := len(settings.Patrols)
	if originalCount == 0 {
		fmt.Println("❌ В исходном файле нет патрулей. Добавление новых патрулей невозможно.")
		return
	}

	countToAdd := originalCount * cfg.PatrolMultiplier
	fmt.Printf("📦 Исходных патрулей: %d. Будет добавлено: %d.\n", originalCount, countToAdd)

	newPatrols := make([]Patrol, 0, originalCount+countToAdd)

	// Обновление существующих патрулей
	for _, p := range settings.Patrols {
		p.NumberOfAI = rand.Intn(cfg.MaxAI-cfg.MinAI+1) + cfg.MinAI
		p.RespawnTime = cfg.RespawnTime
		newPatrols = append(newPatrols, p)
	}

	// Генерация новыйх патрулей
	for i := 0; i < countToAdd; i++ {
		template := settings.Patrols[i%originalCount]
		newP := template
		newP.Name = fmt.Sprintf("%s_cloned_%d", strings.TrimSuffix(template.Name, "_cloned_"), i+1)
		newP.NumberOfAI = rand.Intn(cfg.MaxAI-cfg.MinAI+1) + cfg.MinAI
		newP.RespawnTime = cfg.RespawnTime
		wpCount := rand.Intn(cfg.MaxWaypoints-cfg.MinWaypoints+1) + cfg.MinWaypoints
		newP.Waypoints = nil
		for j := 0; j < wpCount; j++ {
			wp, err := generateWaypoint(usedCoords, cfg)
			if err != nil {
				fmt.Printf("⚠️ Предупреждение: Не удалось сгенерировать Waypoint для патруля #%d. %v\n", i+1, err)
				break
			}
			newP.Waypoints = append(newP.Waypoints, wp)
		}
		newPatrols = append(newPatrols, newP)
		fmt.Printf("✅ Новый патруль #%d: %d точек маршрута\n", i+1, wpCount)
	}

	settings.Patrols = newPatrols

	output, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		fmt.Printf("Ошибка сериализации JSON: %v\n", err)
		return
	}

	if err = os.WriteFile(inputPath, output, 0644); err != nil {
		fmt.Printf("Ошибка сохранения нового JSON: %v\n", err)
		return
	}

	fmt.Printf("🎉 Готово! Всего патрулей: %d. Файл обновлён: %s\n", len(settings.Patrols), inputPath)
}
