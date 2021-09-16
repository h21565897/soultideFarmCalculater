package meal

import (
	"soultide/internal/food"
)

type meal struct {
	name       string
	attraction int
	Foods      map[string]int
}

var (
	// SimplifiedMeals TODO
	SimplifiedMeals = make(map[string]Meal)
	// PrimitiveMeals TODO
	PrimitiveMeals = make(map[string]meal)
)

// Meal TODO
type Meal struct {
	Name       string
	TimeCost   float64
	CoinCost   float64
	Attraction int
}

func addNewMeal(m meal, ms []meal) {
	ms = append(ms, m)
}

// InitMeals TODO
func InitMeals() {
	initMeals(primitiveMeals)
}

// InitMeals TODO
func initMeals(meals []meal) {
	for _, v := range meals {
		var m Meal
		m.Attraction = v.attraction
		m.Name = v.name
		for name, f := range v.Foods {
			m.TimeCost += food.SimpifiedFood[name].TimeCost * float64(f)
			m.CoinCost += food.SimpifiedFood[name].CoinCost * float64(f)
		}
		SimplifiedMeals[m.Name] = m
		PrimitiveMeals[m.Name] = v
	}
}
