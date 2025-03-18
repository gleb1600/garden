package plantsfunc

import (
	"fmt"
	"time"
)

type Plant struct {
	ID             int
	Name           string
	Species        string
	Health         int
	Water          int
	Fertilizer     int
	PlantedAt      time.Time
	LastWatered    time.Time
	LastFertilized time.Time
	Stage          string
	Status         string
}

func UpdatePlantState(plant *Plant) {
	plant.Water -= 5
	plant.Fertilizer -= 2

	switch {
	case plant.Water < 20 && plant.Fertilizer < 10:
		plant.Health -= 15
		plant.Status = "Hungry & Thirsty"
	case plant.Water < 20:
		plant.Health -= 10
		plant.Status = "Thirsty"
	case plant.Fertilizer < 10:
		plant.Health -= 5
		plant.Status = "Hungry"
	default:
		plant.Health += 5
		plant.Status = "Healthy"
	}

	if plant.Health > 100 {
		plant.Health = 100
	}
	if plant.Health < 0 {
		plant.Health = 0
	}
	if plant.Water < 0 {
		plant.Water = 0
	}
	if plant.Fertilizer < 0 {
		plant.Fertilizer = 0
	}
}

func PrintGardenState(plants []Plant) {
	fmt.Printf("\n====================================================================\n")
	fmt.Println("ID | Вид     | Название      | Здоровье | Вода | Удобрение | Статус")
	fmt.Println("---+-------------------------+----------+------+-----------+---------")
	for _, p := range plants {
		healthColor := "\033[32m" // green
		if p.Health < 50 {
			healthColor = "\033[31m"
		} // red

		fmt.Printf(
			"%2d | %-7s | %-13s | %s%3d%%\033[0m     | %3d%% | %3d%%      | %s\n",
			p.ID,
			p.Species,
			p.Name,
			healthColor,
			p.Health,
			p.Water,
			p.Fertilizer,
			p.Status,
		)
	}
}
