package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
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

// generateWaypoint создаёт уникальную точку Waypoint в пределах карты
func generateWaypoint(existing map[string]struct{}, cfg Config) []float64 {
	var x, z float64
	for {
		x = float64(rand.Intn(cfg.MapMaxCoord-cfg.MapMinCoord) + cfg.MapMinCoord)
		z = float64(rand.Intn(cfg.MapMaxCoord-cfg.MapMinCoord) + cfg.MapMinCoord)
		key := fmt.Sprintf("%.0f_%.0f", x, z)
		if _, exists := existing[key]; !exists {
			existing[key] = struct{}{}
			break
		}
	}
	return []float64{x, 0.0, z}
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
	cfg, err := loadConfig("config.yml")
	if err != nil {
		fmt.Printf("Ошибка загрузки конфигурации: %v\n", err)
		return
	}

	inputPath := filepath.Join(filepath.Dir(cfg.Path), "AIPatrolSettings.json")
	
	backupBase := filepath.Join(filepath.Dir(cfg.Path), "backups", "AIPatrolSettings_backup.json")
	backupPath := getBackupPathWithIndex(backupBase)

	data, err := os.ReadFile(inputPath)
	if err != nil {
		panic("Ошибка чтения исходного JSON: " + err.Error())
	}

	if err = os.WriteFile(backupPath, data, 0644); err != nil {
		panic("Ошибка записи бэкапа: " + err.Error())
	}

	var settings Settings
	if err = json.Unmarshal(data, &settings); err != nil {
		panic("Ошибка парсинга JSON: " + err.Error())
	}

	// Уникальные координаты Waypoint’ов
	usedCoords := make(map[string]struct{})
	for _, p := range settings.Patrols {
		for _, wp := range p.Waypoints {
			key := fmt.Sprintf("%.0f_%.0f", wp[0], wp[2])
			usedCoords[key] = struct{}{}
		}
	}

	// Обновление существующих патрулей
	var newPatrols []Patrol
	for _, p := range settings.Patrols {
		p.NumberOfAI = rand.Intn(cfg.MaxAI-cfg.MinAI+1) + cfg.MinAI
		p.RespawnTime = cfg.RespawnTime
		newPatrols = append(newPatrols, p)
	}

	originalCount := len(settings.Patrols)
	countToAdd := originalCount * cfg.PatrolMultiplier
	fmt.Printf("📦 Исходных патрулей: %d. Будет добавлено: %d.\n", originalCount, countToAdd)

	// Генерация новых патрулей
	for i := 0; i < countToAdd; i++ {
		template := settings.Patrols[i%originalCount]
		newP := template
		newP.NumberOfAI = rand.Intn(cfg.MaxAI-cfg.MinAI+1) + cfg.MinAI
		newP.RespawnTime = cfg.RespawnTime
		wpCount := rand.Intn(cfg.MaxWaypoints-cfg.MinWaypoints+1) + cfg.MinWaypoints
		newP.Waypoints = nil
		for j := 0; j < wpCount; j++ {
			newP.Waypoints = append(newP.Waypoints, generateWaypoint(usedCoords, cfg))
		}
		newPatrols = append(newPatrols, newP)
		fmt.Printf("✅ Новый патруль #%d: %d точек маршрута\n", i+1, wpCount)
	}

	settings.Patrols = newPatrols

	output, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		panic("Ошибка сериализации JSON: " + err.Error())
	}
	if err = os.WriteFile(inputPath, output, 0644); err != nil {
		panic("Ошибка сохранения нового JSON: " + err.Error())
	}

	fmt.Printf("🎉 Готово! Всего патрулей: %d. Бэкап: %s\n", len(settings.Patrols), backupPath)
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
