package calculate

import (
	"soultide/internal/blocks"
	doll2 "soultide/internal/doll"
	"soultide/internal/food"
	"soultide/internal/meal"
)

type solver struct {
	meals     []pMeal
	userFoods map[string]int // 当前拥有的数量

	threshold           float64
	acceptableDeviation int
	targetAttraction    int

	res Result

	p currParams
}
type pMeal struct {
	maxCount int
	meal     meal.Meal
}

type currParams struct {
	currAttraction int
	currTimeCost   float64 // 种地的消耗时间
	currTimeWait   float64 // 商店可能的等待时间
	currTimeTrue   float64
	currCoinCost   float64
	currRes        map[string]int
	currNeededFood map[string]int
	currBlocks     map[string]int
}

// Result TODO
type Result struct {
	Result           map[string]int
	ResTimeCost      float64
	ResTimeCostFarm  float64
	ResTimeCostStore float64
	ResCoinCost      float64
	ResAttraction    int
	ResNeededFood    map[string]int
	ResBlocks        map[string]int
}

func (s *solver) checkParams() bool {
	if s.p.currAttraction > s.targetAttraction+s.acceptableDeviation ||
		s.p.currAttraction < s.targetAttraction-s.acceptableDeviation {
		return false
	}
	// 统计当前的需要的食物总量
	for k, v := range s.p.currRes {
		if v > 0 {
			for n, m := range meal.PrimitiveMeals[k].Foods {
				s.p.currNeededFood[n] += m * v
			}
		}
	}
	// 减去已经有的
	for k, v := range s.p.currNeededFood {
		if n, ok := s.userFoods[k]; ok {
			if v > n {
				s.p.currNeededFood[k] = v - n
			} else {
				s.p.currNeededFood[k] = 0
			}
		}
		// 先直接计算一波花费
		s.p.currCoinCost += food.SimpifiedFood[k].CoinCost * float64(s.p.currNeededFood[k])
	}
	// 花费已经超过了
	if s.p.currCoinCost/float64(s.p.currAttraction) > s.threshold {
		return false
	}
	// 尝试街区
	blockSolver := blocks.NewBlockSlover(s.p.currNeededFood)
	res := blockSolver.Solve()
	s.p.currTimeTrue = res.ResTimeCost
	s.p.currNeededFood = res.ResFoodNeeded
	s.p.currTimeWait = res.StoreTimeCost
	s.p.currTimeCost = res.FarmTimeCost
	s.p.currBlocks = res.ResResult
	return true
}

func (s *solver) solve(depth int) {
	if depth == len(s.meals) {
		s.p.reset2()
		if !s.checkParams() {
			return
		}
		if s.p.currTimeTrue < s.res.ResTimeCost {
			s.res.Result = make(map[string]int)
			for k, v := range s.p.currRes {
				s.res.Result[k] = v
			}
			s.res.ResTimeCost = s.p.currTimeTrue
			s.res.ResCoinCost = s.p.currCoinCost
			s.res.ResAttraction = s.p.currAttraction
			s.res.ResNeededFood = s.p.currNeededFood
			s.res.ResTimeCostFarm = s.p.currTimeCost
			s.res.ResTimeCostStore = s.p.currTimeWait
			s.res.ResBlocks = s.p.currBlocks
			//fmt.Println(s.res.ResTimeCostStore)
			//fmt.Println(s.res.ResNeededFood)
		}
		return
	}
	currMeal := s.meals[depth]
	// 取当时的
	currsolver := *s

	for i := 0; i <= currMeal.maxCount; i++ {
		s.p.currAttraction = currsolver.p.currAttraction + currMeal.meal.Attraction*i
		s.p.currRes[currMeal.meal.Name] = i
		s.solve(depth + 1)
	}

}
func (s *currParams) reset2() {
	s.currCoinCost = 0
	s.currTimeCost = 0
	s.currNeededFood = make(map[string]int)
	s.currTimeWait = 0
	s.currTimeTrue = 0
}

func (s *currParams) reset() {
	s.currCoinCost = 0
	s.currTimeCost = 0
	s.currAttraction = 0
	s.currRes = make(map[string]int)
	s.currNeededFood = make(map[string]int)
	s.currTimeWait = 0
	s.currTimeTrue = 0
}

// NewSolver TODO
func NewSolver(dollName string, targetAttraction int, threshold float64, acceptableDeviation int, userFoods map[string]int) *solver {
	s := &solver{
		meals:               nil,
		threshold:           threshold,
		acceptableDeviation: acceptableDeviation,
		targetAttraction:    targetAttraction,
		res:                 Result{},
		p:                   currParams{},
		userFoods:           userFoods,
	}
	currDoll := doll2.SimplifiedDolls[dollName]
	s.meals = make([]pMeal, 0, 100)
	for _, v := range currDoll.Favorites {
		pm := pMeal{
			maxCount: targetAttraction / meal.SimplifiedMeals[v].Attraction,
			meal:     meal.SimplifiedMeals[v],
		}
		s.meals = append(s.meals, pm)
	}
	s.res.ResCoinCost = 999999
	s.res.ResTimeCost = 999999
	s.p.reset()
	return s
}

// Solve TODO
func Solve(name string, target int, threshold float64, acceptableDeviation int, userFoods map[string]int) Result {
	solver := NewSolver(name, target, threshold, acceptableDeviation, userFoods)
	solver.solve(0)
	return solver.res
}
