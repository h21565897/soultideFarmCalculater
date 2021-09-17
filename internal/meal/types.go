package meal

import (
	"soultide/internal/blocks"
	"soultide/internal/food"
)

type meal struct {
	Name       string
	Attraction int
	Foods      map[string]int
}

var (
	// SimplifiedMeals TODO
	SimplifiedMeals = make(map[string]int) // 名字查询菜的索引

	// ParsedMeals TODO
	ParsedMeals = make([]Meal, len(primitiveMeals))
)

// Meal TODO
type Meal struct {
	Name       string
	Attraction int
	FoodCnt    int
	Foods      []int // 每个槽是对应的食材的需求数量
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
	for k, v := range meals {
		var m Meal
		m.Attraction = v.Attraction
		m.Name = v.Name
		m.Foods = blocks.NewFoodSlice()
		for k, v := range v.Foods {
			m.Foods[food.GetFoodIdByName(k)] = v
		}
		ParsedMeals[k] = m
		SimplifiedMeals[v.Name] = k
	}

}

// GetMealById TODO
func GetMealById(id int) Meal {
	return ParsedMeals[id]
}

// GetMealIdByName TODO
func GetMealIdByName(name string) int {
	return SimplifiedMeals[name]
}
