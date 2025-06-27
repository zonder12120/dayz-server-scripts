package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
)

// 🔧 Пути к файлам
const (
	inputPath  = "C:\\Games\\DayZServer\\mpmissions\\dayzOffline.enoch\\expansion\\settings\\AIPatrolSettings.json"
	backupPath = "C:\\Games\\DayZServer\\mpmissions\\dayzOffline.enoch\\expansion\\settings\\AIPatrolSettings_backup.json"
)

// 🔧 Настройки генерации патрулей
const (
	minAI        = 10    // Минимальное количество ботов в патруле
	maxAI        = 15    // Максимальное количество ботов в патруле
	respawnTime  = 600.0 // Время респавна в секундах (10 минут)
	minWaypoints = 3     // Минимальное количество Waypoints на новый патруль
	maxWaypoints = 7     // Максимальное количество Waypoints на новый патруль
	mapMinCoord  = 250   // Минимальная координата X/Z на карте
	mapMaxCoord  = 15000 // Максимальная координата X/Z на карте
	// mapMaxCoord = 15000 // Для ChernarusPlus
	// mapMaxCoord = 12800 // Для Livonia (Enoch)
	// mapMaxCoord = 12800 // Для Namalsk
	// mapMaxCoord = 20480 // Для Deer Isle
	// mapMaxCoord = 10240 // Для Esseker
	// mapMaxCoord = 20480 // Для Takistan
	// mapMaxCoord = 12800 // Для Banov
	// mapMaxCoord = 10240 // Для Pripyat
	// mapMaxCoord = 10240 // Для Valning
	// mapMaxCoord = 12288 // Для Iztek
	patrolMultiplier = 2 // Во сколько раз увеличить количество патрулей основываясь на исходном файле. Всегда необходимо иметь backup со стандартными значениями.
)

// generateWaypoint создаёт уникальную точку Waypoint в пределах карты
func generateWaypoint(existing map[string]struct{}) []float64 {
	var x, z float64
	for {
		x = float64(rand.Intn(mapMaxCoord-mapMinCoord) + mapMinCoord)
		z = float64(rand.Intn(mapMaxCoord-mapMinCoord) + mapMinCoord)
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
	data, err := os.ReadFile(inputPath)
	if err != nil {
		panic("Ошибка чтения исходного JSON: " + err.Error())
	}

	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		panic("Ошибка записи бэкапа: " + err.Error())
	}

	var settings Settings
	if err := json.Unmarshal(data, &settings); err != nil {
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
		p.NumberOfAI = rand.Intn(maxAI-minAI+1) + minAI
		p.RespawnTime = respawnTime
		newPatrols = append(newPatrols, p)
	}

	originalCount := len(settings.Patrols)
	countToAdd := originalCount * patrolMultiplier
	fmt.Printf("📦 Исходных патрулей: %d. Будет добавлено: %d.\n", originalCount, countToAdd)

	// Генерация новых патрулей
	for i := 0; i < countToAdd; i++ {
		template := settings.Patrols[i%originalCount]
		newP := template
		newP.NumberOfAI = rand.Intn(maxAI-minAI+1) + minAI
		newP.RespawnTime = respawnTime
		wpCount := rand.Intn(maxWaypoints-minWaypoints+1) + minWaypoints
		newP.Waypoints = nil
		for j := 0; j < wpCount; j++ {
			newP.Waypoints = append(newP.Waypoints, generateWaypoint(usedCoords))
		}
		newPatrols = append(newPatrols, newP)
		fmt.Printf("✅ Новый патруль #%d: %d точек маршрута\n", i+1, wpCount)
	}

	settings.Patrols = newPatrols

	output, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		panic("Ошибка сериализации JSON: " + err.Error())
	}
	if err := os.WriteFile(inputPath, output, 0644); err != nil {
		panic("Ошибка сохранения нового JSON: " + err.Error())
	}

	fmt.Printf("🎉 Готово! Всего патрулей: %d. Бэкап: %s\n", len(settings.Patrols), backupPath)
}
